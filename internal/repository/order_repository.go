package repository

import (
	"context"

	"github.com/imkarthi24/sf-backend/internal/entities"
	"github.com/imkarthi24/sf-backend/internal/repository/scopes"
	"github.com/imkarthi24/sf-backend/pkg/db"
	"github.com/imkarthi24/sf-backend/pkg/errs"
)

type OrderRepository interface {
	Create(*context.Context, *entities.Order) *errs.XError
	Update(*context.Context, *entities.Order) *errs.XError
	Get(*context.Context, uint) (*entities.Order, *errs.XError)
	GetAll(*context.Context, string) ([]entities.Order, *errs.XError)
	Delete(*context.Context, uint) *errs.XError
}

type orderRepository struct {
	GormDAL
}

func ProvideOrderRepository(customDB GormDAL) OrderRepository {
	return &orderRepository{GormDAL: customDB}
}

func (or *orderRepository) Create(ctx *context.Context, order *entities.Order) *errs.XError {
	res := or.WithDB(ctx).Create(&order)
	if res.Error != nil {
		return errs.NewXError(errs.DATABASE, "Unable to save order", res.Error)
	}
	return nil
}

func (or *orderRepository) Update(ctx *context.Context, order *entities.Order) *errs.XError {
	return or.GormDAL.Update(ctx, *order)
}

func (or *orderRepository) Get(ctx *context.Context, id uint) (*entities.Order, *errs.XError) {
	order := entities.Order{}
	res := or.WithDB(ctx).
		Preload("Customer").
		Preload("OrderTakenBy", scopes.SelectFields("first_name", "last_name")).
		Preload("OrderItems.Person").
		Preload("OrderItems.DressType").
		Preload("OrderItems.Measurement").
		Find(&order, id)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find order", res.Error)
	}
	return &order, nil
}

func (or *orderRepository) GetAll(ctx *context.Context, search string) ([]entities.Order, *errs.XError) {
	var orders []entities.Order

	res := or.WithDB(ctx).Model(&entities.Order{}).
		Select(`"stitch"."Orders".*,
			(SELECT COALESCE(SUM(quantity), 0) FROM "stitch"."OrderItems" 
			 WHERE "stitch"."OrderItems".order_id = "stitch"."Orders".id) as order_quantity,
			(SELECT COALESCE(SUM(total), 0) FROM "stitch"."OrderItems" 
			 WHERE "stitch"."OrderItems".order_id = "stitch"."Orders".id) as order_value`).
		Scopes(db.Paginate(ctx)).
		Preload("Customer", scopes.SelectFields("first_name", "last_name")).
		Preload("OrderTakenBy", scopes.SelectFields("first_name", "last_name")).
		Preload("OrderItems.Person").
		Preload("OrderItems.DressType").
		Preload("OrderItems.Measurement").
		Find(&orders)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find orders", res.Error)
	}
	return orders, nil
}

func (or *orderRepository) Delete(ctx *context.Context, id uint) *errs.XError {
	order := &entities.Order{Model: &entities.Model{ID: id, IsActive: false}}
	err := or.GormDAL.Delete(ctx, order)
	if err != nil {
		return err
	}
	return nil
}
