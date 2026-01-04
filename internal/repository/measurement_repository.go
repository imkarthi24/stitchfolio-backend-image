package repository

import (
	"context"

	"github.com/imkarthi24/sf-backend/internal/entities"
	"github.com/imkarthi24/sf-backend/internal/repository/common"
	"github.com/imkarthi24/sf-backend/internal/repository/scopes"
	"github.com/imkarthi24/sf-backend/pkg/db"
	"github.com/imkarthi24/sf-backend/pkg/errs"
)

type MeasurementRepository interface {
	Create(*context.Context, *entities.Measurement) *errs.XError
	Update(*context.Context, *entities.Measurement) *errs.XError
	Get(*context.Context, uint) (*entities.Measurement, *errs.XError)
	GetAll(*context.Context, string) ([]entities.Measurement, *errs.XError)
	Delete(*context.Context, uint) *errs.XError
}

type measurementRepository struct {
	txn      db.DBTransactionManager
	customDB common.CustomGormDB
}

func ProvideMeasurementRepository(txn db.DBTransactionManager, customDB common.CustomGormDB) MeasurementRepository {
	return &measurementRepository{txn: txn, customDB: customDB}
}

func (mr *measurementRepository) Create(ctx *context.Context, measurement *entities.Measurement) *errs.XError {
	res := mr.txn.Txn(ctx).Create(&measurement)
	if res.Error != nil {
		return errs.NewXError(errs.DATABASE, "Unable to save measurement", res.Error)
	}
	return nil
}

func (mr *measurementRepository) Update(ctx *context.Context, measurement *entities.Measurement) *errs.XError {
	return mr.customDB.Update(ctx, *measurement)
}

func (mr *measurementRepository) Get(ctx *context.Context, id uint) (*entities.Measurement, *errs.XError) {
	measurement := entities.Measurement{}
	res := mr.txn.Txn(ctx).Preload("Customer").Find(&measurement, id)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find measurement", res.Error)
	}
	return &measurement, nil
}

func (mr *measurementRepository) GetAll(ctx *context.Context, search string) ([]entities.Measurement, *errs.XError) {
	var measurements []entities.Measurement
	res := mr.txn.Txn(ctx).
		Scopes(scopes.Channel(), scopes.IsActive()).
		Scopes(scopes.ILike(search, "dress_type", "measurement_by")).
		Scopes(db.Paginate(ctx)).
		Preload("Customer").
		Find(&measurements)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find measurements", res.Error)
	}
	return measurements, nil
}

func (mr *measurementRepository) Delete(ctx *context.Context, id uint) *errs.XError {
	measurement := &entities.Measurement{Model: &entities.Model{ID: id, IsActive: false}}
	err := mr.customDB.Delete(ctx, measurement)
	if err != nil {
		return err
	}
	return nil
}
