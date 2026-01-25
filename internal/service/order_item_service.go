package service

import (
	"context"

	"github.com/imkarthi24/sf-backend/internal/mapper"
	requestModel "github.com/imkarthi24/sf-backend/internal/model/request"
	responseModel "github.com/imkarthi24/sf-backend/internal/model/response"
	"github.com/imkarthi24/sf-backend/internal/repository"
	"github.com/loop-kar/pixie/errs"
)

type OrderItemService interface {
	SaveOrderItem(*context.Context, requestModel.OrderItem) *errs.XError
	UpdateOrderItem(*context.Context, requestModel.OrderItem, uint) *errs.XError
	Get(*context.Context, uint) (*responseModel.OrderItem, *errs.XError)
	GetAll(*context.Context, string) ([]responseModel.OrderItem, *errs.XError)
	Delete(*context.Context, uint) *errs.XError
}

type orderItemService struct {
	orderItemRepo repository.OrderItemRepository
	mapper        mapper.Mapper
	respMapper    mapper.ResponseMapper
}

func ProvideOrderItemService(repo repository.OrderItemRepository, mapper mapper.Mapper, respMapper mapper.ResponseMapper) OrderItemService {
	return orderItemService{
		orderItemRepo: repo,
		mapper:        mapper,
		respMapper:    respMapper,
	}
}

func (svc orderItemService) SaveOrderItem(ctx *context.Context, orderItem requestModel.OrderItem) *errs.XError {
	dbOrderItem, err := svc.mapper.OrderItem(orderItem)
	if err != nil {
		return errs.NewXError(errs.INVALID_REQUEST, "Unable to save order item", err)
	}

	errr := svc.orderItemRepo.Create(ctx, dbOrderItem)
	if errr != nil {
		return errr
	}

	return nil
}

func (svc orderItemService) UpdateOrderItem(ctx *context.Context, orderItem requestModel.OrderItem, id uint) *errs.XError {
	dbOrderItem, err := svc.mapper.OrderItem(orderItem)
	if err != nil {
		return errs.NewXError(errs.INVALID_REQUEST, "Unable to update order item", err)
	}

	dbOrderItem.ID = id
	errr := svc.orderItemRepo.Update(ctx, dbOrderItem)
	if errr != nil {
		return errr
	}
	return nil
}

func (svc orderItemService) Get(ctx *context.Context, id uint) (*responseModel.OrderItem, *errs.XError) {
	orderItem, err := svc.orderItemRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	mappedOrderItem, mapErr := svc.respMapper.OrderItem(orderItem)
	if mapErr != nil {
		return nil, errs.NewXError(errs.MAPPING_ERROR, "Failed to map OrderItem data", mapErr)
	}

	return mappedOrderItem, nil
}

func (svc orderItemService) GetAll(ctx *context.Context, search string) ([]responseModel.OrderItem, *errs.XError) {
	orderItems, err := svc.orderItemRepo.GetAll(ctx, search)
	if err != nil {
		return nil, err
	}

	mappedOrderItems, mapErr := svc.respMapper.OrderItems(orderItems)
	if mapErr != nil {
		return nil, errs.NewXError(errs.MAPPING_ERROR, "Failed to map OrderItem data", mapErr)
	}

	return mappedOrderItems, nil
}

func (svc orderItemService) Delete(ctx *context.Context, id uint) *errs.XError {
	err := svc.orderItemRepo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
