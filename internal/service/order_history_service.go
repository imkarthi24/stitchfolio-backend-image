package service

import (
	"context"

	"github.com/imkarthi24/sf-backend/internal/mapper"
	requestModel "github.com/imkarthi24/sf-backend/internal/model/request"
	responseModel "github.com/imkarthi24/sf-backend/internal/model/response"
	"github.com/imkarthi24/sf-backend/internal/repository"
	"github.com/loop-kar/pixie/errs"
)

type OrderHistoryService interface {
	SaveOrderHistory(*context.Context, requestModel.OrderHistory) *errs.XError
	Get(*context.Context, uint) (*responseModel.OrderHistory, *errs.XError)
	GetAll(*context.Context, string) ([]responseModel.OrderHistory, *errs.XError)
	GetByOrderId(*context.Context, uint) ([]responseModel.OrderHistory, *errs.XError)
}

type orderHistoryService struct {
	orderHistoryRepo repository.OrderHistoryRepository
	mapper           mapper.Mapper
	respMapper       mapper.ResponseMapper
}

func ProvideOrderHistoryService(repo repository.OrderHistoryRepository, mapper mapper.Mapper, respMapper mapper.ResponseMapper) OrderHistoryService {
	return orderHistoryService{
		orderHistoryRepo: repo,
		mapper:           mapper,
		respMapper:       respMapper,
	}
}

func (svc orderHistoryService) SaveOrderHistory(ctx *context.Context, orderHistory requestModel.OrderHistory) *errs.XError {
	dbOrderHistory, err := svc.mapper.OrderHistory(orderHistory)
	if err != nil {
		return errs.NewXError(errs.INVALID_REQUEST, "Unable to save order history", err)
	}

	errr := svc.orderHistoryRepo.Create(ctx, dbOrderHistory)
	if errr != nil {
		return errr
	}

	return nil
}

func (svc orderHistoryService) Get(ctx *context.Context, id uint) (*responseModel.OrderHistory, *errs.XError) {
	orderHistory, err := svc.orderHistoryRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	mappedOrderHistory, mapErr := svc.respMapper.OrderHistory(orderHistory)
	if mapErr != nil {
		return nil, errs.NewXError(errs.MAPPING_ERROR, "Failed to map OrderHistory data", mapErr)
	}

	return mappedOrderHistory, nil
}

func (svc orderHistoryService) GetAll(ctx *context.Context, search string) ([]responseModel.OrderHistory, *errs.XError) {
	orderHistories, err := svc.orderHistoryRepo.GetAll(ctx, search)
	if err != nil {
		return nil, err
	}

	mappedOrderHistories, mapErr := svc.respMapper.OrderHistories(orderHistories)
	if mapErr != nil {
		return nil, errs.NewXError(errs.MAPPING_ERROR, "Failed to map OrderHistory data", mapErr)
	}

	return mappedOrderHistories, nil
}

func (svc orderHistoryService) GetByOrderId(ctx *context.Context, orderId uint) ([]responseModel.OrderHistory, *errs.XError) {
	orderHistories, err := svc.orderHistoryRepo.GetByOrderId(ctx, orderId)
	if err != nil {
		return nil, err
	}

	mappedOrderHistories, mapErr := svc.respMapper.OrderHistories(orderHistories)
	if mapErr != nil {
		return nil, errs.NewXError(errs.MAPPING_ERROR, "Failed to map OrderHistory data", mapErr)
	}

	return mappedOrderHistories, nil
}
