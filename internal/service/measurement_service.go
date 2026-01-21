package service

import (
	"context"

	"github.com/imkarthi24/sf-backend/internal/entities"
	entitiy_types "github.com/imkarthi24/sf-backend/internal/entities/types"
	"github.com/imkarthi24/sf-backend/internal/mapper"
	requestModel "github.com/imkarthi24/sf-backend/internal/model/request"
	responseModel "github.com/imkarthi24/sf-backend/internal/model/response"
	"github.com/imkarthi24/sf-backend/internal/repository"
	"github.com/imkarthi24/sf-backend/internal/utils"
	"github.com/imkarthi24/sf-backend/pkg/errs"
	"github.com/imkarthi24/sf-backend/pkg/util"
)

type MeasurementService interface {
	SaveMeasurement(*context.Context, requestModel.Measurement) *errs.XError
	UpdateMeasurement(*context.Context, requestModel.Measurement, uint) *errs.XError
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

func (svc measurementService) UpdateMeasurement(ctx *context.Context, measurement requestModel.Measurement, id uint) *errs.XError {
	// Get the old measurement values before updating
	oldMeasurement, err := svc.measurementRepo.Get(ctx, id)
	if err != nil {
		return err
	}

	dbMeasurement, mapErr := svc.mapper.Measurement(measurement)
	if mapErr != nil {
		return errs.NewXError(errs.INVALID_REQUEST, "Unable to update measurement", mapErr)
	}

	dbMeasurement.ID = id
	errr := svc.measurementRepo.Update(ctx, dbMeasurement)
	if errr != nil {
		return errr
	}

	// Record measurement history for UPDATED action with old values
	errr = svc.recordMeasurementHistory(ctx, id, entities.MeasurementHistoryActionUpdated, &oldMeasurement.Value)
	if errr != nil {
		return errr
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
