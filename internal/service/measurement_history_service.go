package service

import (
	"context"

	"github.com/imkarthi24/sf-backend/internal/mapper"
	requestModel "github.com/imkarthi24/sf-backend/internal/model/request"
	responseModel "github.com/imkarthi24/sf-backend/internal/model/response"
	"github.com/imkarthi24/sf-backend/internal/repository"
	"github.com/loop-kar/pixie/errs"
)

type MeasurementHistoryService interface {
	SaveMeasurementHistory(*context.Context, requestModel.MeasurementHistory) *errs.XError
	Get(*context.Context, uint) (*responseModel.MeasurementHistory, *errs.XError)
	GetAll(*context.Context, string) ([]responseModel.MeasurementHistory, *errs.XError)
	GetByMeasurementId(*context.Context, uint) ([]responseModel.MeasurementHistory, *errs.XError)
}

type measurementHistoryService struct {
	measurementHistoryRepo repository.MeasurementHistoryRepository
	mapper                 mapper.Mapper
	respMapper             mapper.ResponseMapper
}

func ProvideMeasurementHistoryService(repo repository.MeasurementHistoryRepository, mapper mapper.Mapper, respMapper mapper.ResponseMapper) MeasurementHistoryService {
	return measurementHistoryService{
		measurementHistoryRepo: repo,
		mapper:                 mapper,
		respMapper:             respMapper,
	}
}

func (svc measurementHistoryService) SaveMeasurementHistory(ctx *context.Context, measurementHistory requestModel.MeasurementHistory) *errs.XError {
	dbMeasurementHistory, err := svc.mapper.MeasurementHistory(measurementHistory)
	if err != nil {
		return errs.NewXError(errs.INVALID_REQUEST, "Unable to save measurement history", err)
	}

	errr := svc.measurementHistoryRepo.Create(ctx, dbMeasurementHistory)
	if errr != nil {
		return errr
	}

	return nil
}

func (svc measurementHistoryService) Get(ctx *context.Context, id uint) (*responseModel.MeasurementHistory, *errs.XError) {
	measurementHistory, err := svc.measurementHistoryRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	mappedMeasurementHistory, mapErr := svc.respMapper.MeasurementHistory(measurementHistory)
	if mapErr != nil {
		return nil, errs.NewXError(errs.MAPPING_ERROR, "Failed to map MeasurementHistory data", mapErr)
	}

	return mappedMeasurementHistory, nil
}

func (svc measurementHistoryService) GetAll(ctx *context.Context, search string) ([]responseModel.MeasurementHistory, *errs.XError) {
	measurementHistories, err := svc.measurementHistoryRepo.GetAll(ctx, search)
	if err != nil {
		return nil, err
	}

	mappedMeasurementHistories, mapErr := svc.respMapper.MeasurementHistories(measurementHistories)
	if mapErr != nil {
		return nil, errs.NewXError(errs.MAPPING_ERROR, "Failed to map MeasurementHistory data", mapErr)
	}

	return mappedMeasurementHistories, nil
}

func (svc measurementHistoryService) GetByMeasurementId(ctx *context.Context, measurementId uint) ([]responseModel.MeasurementHistory, *errs.XError) {
	measurementHistories, err := svc.measurementHistoryRepo.GetByMeasurementId(ctx, measurementId)
	if err != nil {
		return nil, err
	}

	mappedMeasurementHistories, mapErr := svc.respMapper.MeasurementHistories(measurementHistories)
	if mapErr != nil {
		return nil, errs.NewXError(errs.MAPPING_ERROR, "Failed to map MeasurementHistory data", mapErr)
	}

	return mappedMeasurementHistories, nil
}
