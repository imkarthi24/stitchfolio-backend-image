package service

import (
	"context"

	"github.com/imkarthi24/sf-backend/internal/mapper"
	requestModel "github.com/imkarthi24/sf-backend/internal/model/request"
	responseModel "github.com/imkarthi24/sf-backend/internal/model/response"
	"github.com/imkarthi24/sf-backend/internal/repository"
	"github.com/loop-kar/pixie/errs"
)

type OrderService interface {
	SaveOrder(*context.Context, requestModel.Order) *errs.XError
	UpdateOrder(*context.Context, requestModel.Order, uint) *errs.XError
	Get(*context.Context, uint) (*responseModel.Order, *errs.XError)
	GetAll(*context.Context, string) ([]responseModel.Order, *errs.XError)
	Delete(*context.Context, uint) *errs.XError
}

type orderService struct {
	orderRepo  repository.OrderRepository
	mapper     mapper.Mapper
	respMapper mapper.ResponseMapper
}

func ProvideOrderService(repo repository.OrderRepository, mapper mapper.Mapper, respMapper mapper.ResponseMapper) OrderService {
	return orderService{
		orderRepo:  repo,
		mapper:     mapper,
		respMapper: respMapper,
	}
}

func (svc orderService) SaveOrder(ctx *context.Context, order requestModel.Order) *errs.XError {
	dbOrder, err := svc.mapper.Order(order)
	if err != nil {
		return errs.NewXError(errs.INVALID_REQUEST, "Unable to save order", err)
	}

	errr := svc.orderRepo.Create(ctx, dbOrder)
	if errr != nil {
		return errr
	}

	return nil
}

func (svc orderService) UpdateOrder(ctx *context.Context, order requestModel.Order, id uint) *errs.XError {
	dbOrder, err := svc.mapper.Order(order)
	if err != nil {
		return errs.NewXError(errs.INVALID_REQUEST, "Unable to update order", err)
	}

	dbOrder.ID = id
	errr := svc.orderRepo.Update(ctx, dbOrder)
	if errr != nil {
		return errr
	}
	return nil
}

func (svc orderService) Get(ctx *context.Context, id uint) (*responseModel.Order, *errs.XError) {
	order, err := svc.orderRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	mappedOrder, mapErr := svc.respMapper.Order(order)
	if mapErr != nil {
		return nil, errs.NewXError(errs.MAPPING_ERROR, "Failed to map Order data", mapErr)
	}

	return mappedOrder, nil
}

func (svc orderService) GetAll(ctx *context.Context, search string) ([]responseModel.Order, *errs.XError) {
	orders, err := svc.orderRepo.GetAll(ctx, search)
	if err != nil {
		return nil, err
	}

	mappedOrders, mapErr := svc.respMapper.Orders(orders)
	if mapErr != nil {
		return nil, errs.NewXError(errs.MAPPING_ERROR, "Failed to map Order data", mapErr)
	}

	return mappedOrders, nil
}

func (svc orderService) Delete(ctx *context.Context, id uint) *errs.XError {
	err := svc.orderRepo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
