package service

import (
	"context"

	"github.com/imkarthi24/sf-backend/internal/mapper"
	requestModel "github.com/imkarthi24/sf-backend/internal/model/request"
	responseModel "github.com/imkarthi24/sf-backend/internal/model/response"
	"github.com/imkarthi24/sf-backend/internal/repository"
	"github.com/imkarthi24/sf-backend/pkg/errs"
)

type EnquiryHistoryService interface {
	SaveEnquiryHistory(*context.Context, requestModel.EnquiryHistory) *errs.XError
	Get(*context.Context, uint) (*responseModel.EnquiryHistory, *errs.XError)
	GetAll(*context.Context, string) ([]responseModel.EnquiryHistory, *errs.XError)
	GetByEnquiryId(*context.Context, uint) ([]responseModel.EnquiryHistory, *errs.XError)
}

type enquiryHistoryService struct {
	enquiryHistoryRepo repository.EnquiryHistoryRepository
	mapper             mapper.Mapper
	respMapper         mapper.ResponseMapper
}

func ProvideEnquiryHistoryService(repo repository.EnquiryHistoryRepository, mapper mapper.Mapper, respMapper mapper.ResponseMapper) EnquiryHistoryService {
	return enquiryHistoryService{
		enquiryHistoryRepo: repo,
		mapper:             mapper,
		respMapper:         respMapper,
	}
}

func (svc enquiryHistoryService) SaveEnquiryHistory(ctx *context.Context, enquiryHistory requestModel.EnquiryHistory) *errs.XError {
	dbEnquiryHistory, err := svc.mapper.EnquiryHistory(enquiryHistory)
	if err != nil {
		return errs.NewXError(errs.INVALID_REQUEST, "Unable to save enquiry history", err)
	}

	errr := svc.enquiryHistoryRepo.Create(ctx, dbEnquiryHistory)
	if errr != nil {
		return errr
	}

	return nil
}

func (svc enquiryHistoryService) Get(ctx *context.Context, id uint) (*responseModel.EnquiryHistory, *errs.XError) {
	enquiryHistory, err := svc.enquiryHistoryRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	mappedEnquiryHistory, mapErr := svc.respMapper.EnquiryHistory(enquiryHistory)
	if mapErr != nil {
		return nil, errs.NewXError(errs.MAPPING_ERROR, "Failed to map EnquiryHistory data", mapErr)
	}

	return mappedEnquiryHistory, nil
}

func (svc enquiryHistoryService) GetAll(ctx *context.Context, search string) ([]responseModel.EnquiryHistory, *errs.XError) {
	enquiryHistories, err := svc.enquiryHistoryRepo.GetAll(ctx, search)
	if err != nil {
		return nil, err
	}

	mappedEnquiryHistories, mapErr := svc.respMapper.EnquiryHistories(enquiryHistories)
	if mapErr != nil {
		return nil, errs.NewXError(errs.MAPPING_ERROR, "Failed to map EnquiryHistory data", mapErr)
	}

	return mappedEnquiryHistories, nil
}

func (svc enquiryHistoryService) GetByEnquiryId(ctx *context.Context, enquiryId uint) ([]responseModel.EnquiryHistory, *errs.XError) {
	enquiryHistories, err := svc.enquiryHistoryRepo.GetByEnquiryId(ctx, enquiryId)
	if err != nil {
		return nil, err
	}

	mappedEnquiryHistories, mapErr := svc.respMapper.EnquiryHistories(enquiryHistories)
	if mapErr != nil {
		return nil, errs.NewXError(errs.MAPPING_ERROR, "Failed to map EnquiryHistory data", mapErr)
	}

	return mappedEnquiryHistories, nil
}
