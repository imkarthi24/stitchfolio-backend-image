package service

import (
	"context"

	"github.com/imkarthi24/sf-backend/internal/entities"
	entitiy_types "github.com/imkarthi24/sf-backend/internal/entities/types"
	"github.com/imkarthi24/sf-backend/internal/mapper"
	requestModel "github.com/imkarthi24/sf-backend/internal/model/request"
	responseModel "github.com/imkarthi24/sf-backend/internal/model/response"
	"github.com/imkarthi24/sf-backend/internal/repository"
	"github.com/loop-kar/pixie/errs"
	"github.com/imkarthi24/sf-backend/internal/utils"	
	"github.com/loop-kar/pixie/util"
)

type MeasurementService interface {
	SaveMeasurement(*context.Context, requestModel.Measurement) *errs.XError
	SaveBulkMeasurements(*context.Context, []requestModel.BulkMeasurementRequest) *errs.XError
	UpdateBulkMeasurements(*context.Context, []requestModel.Measurement) *errs.XError
	UpdateBulkMeasurementsByPerson(*context.Context, requestModel.BulkUpdateMeasurementRequest) *errs.XError
	Get(*context.Context, uint) (*responseModel.Measurement, *errs.XError)
	GetAll(*context.Context, string) ([]responseModel.Measurement, *errs.XError)
	Delete(*context.Context, uint) *errs.XError
}

type measurementService struct {
	measurementRepo        repository.MeasurementRepository
	measurementHistoryRepo repository.MeasurementHistoryRepository
	mapper                 mapper.Mapper
	respMapper             mapper.ResponseMapper
}

func ProvideMeasurementService(repo repository.MeasurementRepository, measurementHistoryRepo repository.MeasurementHistoryRepository, mapper mapper.Mapper, respMapper mapper.ResponseMapper) MeasurementService {
	return measurementService{
		measurementRepo:        repo,
		measurementHistoryRepo: measurementHistoryRepo,
		mapper:                 mapper,
		respMapper:             respMapper,
	}
}

func (svc measurementService) SaveMeasurement(ctx *context.Context, measurement requestModel.Measurement) *errs.XError {
	dbMeasurement, err := svc.mapper.Measurement(measurement)
	if err != nil {
		return errs.NewXError(errs.INVALID_REQUEST, "Unable to save measurement", err)
	}

	errr := svc.measurementRepo.Create(ctx, dbMeasurement)
	if errr != nil {
		return errr
	}

	// Record measurement history for CREATED action
	errr = svc.recordMeasurementHistory(ctx, dbMeasurement.ID, entities.MeasurementHistoryActionCreated, nil)
	if errr != nil {
		return errr
	}

	return nil
}

func (svc measurementService) SaveBulkMeasurements(ctx *context.Context, bulkRequests []requestModel.BulkMeasurementRequest) *errs.XError {
	var measurementsToCreate []*entities.Measurement
	userID := utils.GetUserId(ctx)

	for _, bulkRequest := range bulkRequests {
		for _, measurementItem := range bulkRequest.Measurements {
			var valuesJSON entitiy_types.JSON
			if len(measurementItem.Values) > 0 {
				valuesJSON = entitiy_types.JSON(measurementItem.Values)
			}

			measurement := &entities.Measurement{
				Model: &entities.Model{
					IsActive: true,
				},
				Value:       valuesJSON,
				PersonId:    bulkRequest.PersonId,
				DressTypeId: measurementItem.DressTypeId,
				TakenById:   &userID,
			}

			measurementsToCreate = append(measurementsToCreate, measurement)
		}
	}

	// Batch create all measurements
	errr := svc.measurementRepo.BatchCreate(ctx, measurementsToCreate)
	if errr != nil {
		return errr
	}

	// Record measurement history for each created measurement
	for _, measurement := range measurementsToCreate {
		errr = svc.recordMeasurementHistory(ctx, measurement.ID, entities.MeasurementHistoryActionCreated, nil)
		if errr != nil {
			return errr
		}
	}

	return nil
}

func (svc measurementService) UpdateBulkMeasurements(ctx *context.Context, measurements []requestModel.Measurement) *errs.XError {
	var measurementsToUpdate []*entities.Measurement
	var oldMeasurementsMap = make(map[uint]*entities.Measurement)

	// First, collect all IDs and fetch old values
	var ids []uint
	for _, measurement := range measurements {
		if measurement.ID == 0 {
			continue // Skip measurements without IDs
		}
		ids = append(ids, measurement.ID)
	}

	// Fetch old measurements for history tracking
	for _, id := range ids {
		oldMeasurement, err := svc.measurementRepo.Get(ctx, id)
		if err != nil {
			return err
		}
		oldMeasurementsMap[id] = oldMeasurement
	}

	// Map request models to entity models
	for _, measurement := range measurements {
		if measurement.ID == 0 {
			continue // Skip measurements without IDs
		}

		dbMeasurement, mapErr := svc.mapper.Measurement(measurement)
		if mapErr != nil {
			return errs.NewXError(errs.INVALID_REQUEST, "Unable to update measurement", mapErr)
		}

		// Ensure ID is set
		dbMeasurement.ID = measurement.ID
		measurementsToUpdate = append(measurementsToUpdate, dbMeasurement)
	}

	// Batch update all measurements
	errr := svc.measurementRepo.BatchUpdate(ctx, measurementsToUpdate)
	if errr != nil {
		return errr
	}

	// Record measurement history for each updated measurement
	for _, measurement := range measurementsToUpdate {
		oldMeasurement, exists := oldMeasurementsMap[measurement.ID]
		if exists {
			errr = svc.recordMeasurementHistory(ctx, measurement.ID, entities.MeasurementHistoryActionUpdated, &oldMeasurement.Value)
			if errr != nil {
				return errr
			}
		}
	}

	return nil
}

func (svc measurementService) UpdateBulkMeasurementsByPerson(ctx *context.Context, request requestModel.BulkUpdateMeasurementRequest) *errs.XError {
	userID := utils.GetUserId(ctx)
	var measurementsToUpdate []*entities.Measurement
	var measurementsToCreate []*entities.Measurement
	var oldMeasurementsMap = make(map[uint]*entities.Measurement)

	// Process each person's measurements
	for _, personData := range request.Persons {
		for _, measurementItem := range personData.Measurements {
			// Try to find existing measurement by personId + dressTypeId
			existingMeasurement, err := svc.measurementRepo.GetByPersonIdAndDressTypeId(ctx, personData.PersonId, measurementItem.DressTypeId)
			if err != nil {
				return err
			}

			var valuesJSON entitiy_types.JSON
			if len(measurementItem.Values) > 0 {
				valuesJSON = entitiy_types.JSON(measurementItem.Values)
			}

			// Determine isActive value: use provided value or default to true
			isActive := true
			if personData.IsActive != nil {
				isActive = *personData.IsActive
			}

			if existingMeasurement != nil {
				// Store old value for history BEFORE updating
				oldValue := existingMeasurement.Value
				oldMeasurementsMap[existingMeasurement.ID] = &entities.Measurement{
					Value: oldValue,
				}
				// Update existing measurement
				existingMeasurement.Value = valuesJSON
				existingMeasurement.IsActive = isActive
				measurementsToUpdate = append(measurementsToUpdate, existingMeasurement)
			} else {
				// Create new measurement
				newMeasurement := &entities.Measurement{
					Model: &entities.Model{
						IsActive: isActive,
					},
					Value:       valuesJSON,
					PersonId:    personData.PersonId,
					DressTypeId: measurementItem.DressTypeId,
					TakenById:   &userID,
				}
				measurementsToCreate = append(measurementsToCreate, newMeasurement)
			}
		}
	}

	// Batch create new measurements
	if len(measurementsToCreate) > 0 {
		errr := svc.measurementRepo.BatchCreate(ctx, measurementsToCreate)
		if errr != nil {
			return errr
		}

		// Record history for created measurements
		for _, measurement := range measurementsToCreate {
			errr = svc.recordMeasurementHistory(ctx, measurement.ID, entities.MeasurementHistoryActionCreated, nil)
			if errr != nil {
				return errr
			}
		}
	}

	// Batch update existing measurements
	if len(measurementsToUpdate) > 0 {
		errr := svc.measurementRepo.BatchUpdate(ctx, measurementsToUpdate)
		if errr != nil {
			return errr
		}

		// Record history for updated measurements
		for _, measurement := range measurementsToUpdate {
			oldMeasurement, exists := oldMeasurementsMap[measurement.ID]
			if exists {
				errr = svc.recordMeasurementHistory(ctx, measurement.ID, entities.MeasurementHistoryActionUpdated, &oldMeasurement.Value)
				if errr != nil {
					return errr
				}
			}
		}
	}

	return nil
}

func (svc measurementService) Get(ctx *context.Context, id uint) (*responseModel.Measurement, *errs.XError) {
	measurement, err := svc.measurementRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	mappedMeasurement, mapErr := svc.respMapper.Measurement(measurement)
	if mapErr != nil {
		return nil, errs.NewXError(errs.MAPPING_ERROR, "Failed to map Measurement data", mapErr)
	}

	return mappedMeasurement, nil
}

func (svc measurementService) GetAll(ctx *context.Context, search string) ([]responseModel.Measurement, *errs.XError) {
	measurements, err := svc.measurementRepo.GetAll(ctx, search)
	if err != nil {
		return nil, err
	}

	mappedMeasurements, mapErr := svc.respMapper.Measurements(measurements)
	if mapErr != nil {
		return nil, errs.NewXError(errs.MAPPING_ERROR, "Failed to map Measurement data", mapErr)
	}

	return mappedMeasurements, nil
}

func (svc measurementService) Delete(ctx *context.Context, id uint) *errs.XError {
	// Get the measurement values before deleting
	oldMeasurement, err := svc.measurementRepo.Get(ctx, id)
	if err != nil {
		return err
	}

	err = svc.measurementRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	// Record measurement history for DELETED action with old values
	err = svc.recordMeasurementHistory(ctx, id, entities.MeasurementHistoryActionDeleted, &oldMeasurement.Value)
	if err != nil {
		return err
	}

	return nil
}

// recordMeasurementHistory creates a measurement history record
func (svc measurementService) recordMeasurementHistory(ctx *context.Context, measurementId uint, action entities.MeasurementHistoryAction, oldValues *entitiy_types.JSON) *errs.XError {
	userID := utils.GetUserId(ctx)
	performedAt := util.GetLocalTime()

	var oldValuesJSON entitiy_types.JSON
	if oldValues != nil {
		oldValuesJSON = *oldValues
	}

	history := &entities.MeasurementHistory{
		Model:         &entities.Model{IsActive: true},
		Action:        action,
		OldValues:     oldValuesJSON,
		MeasurementId: measurementId,
		PerformedAt:   performedAt,
		PerformedById: userID,
	}

	return svc.measurementHistoryRepo.Create(ctx, history)
}
