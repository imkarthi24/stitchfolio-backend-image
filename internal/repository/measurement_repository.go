package repository

import (
	"context"
	"errors"

	"github.com/imkarthi24/sf-backend/internal/entities"
	"github.com/imkarthi24/sf-backend/internal/repository/scopes"
	"github.com/loop-kar/pixie/db"
	"github.com/loop-kar/pixie/errs"
	"github.com/loop-kar/pixie/util"
  "github.com/loop-kar/pixie/constants"
	
	"gorm.io/gorm"
)

type MeasurementRepository interface {
	Create(*context.Context, *entities.Measurement) *errs.XError
	BatchCreate(*context.Context, []*entities.Measurement) *errs.XError
	Update(*context.Context, *entities.Measurement) *errs.XError
	BatchUpdate(*context.Context, []*entities.Measurement) *errs.XError
	Get(*context.Context, uint) (*entities.Measurement, *errs.XError)
	GetByPersonIdAndDressTypeId(*context.Context, uint, uint) (*entities.Measurement, *errs.XError)
	GetAll(*context.Context, string) ([]entities.Measurement, *errs.XError)
	Delete(*context.Context, uint) *errs.XError
}

type measurementRepository struct {
	GormDAL
}

func ProvideMeasurementRepository(customDB GormDAL) MeasurementRepository {
	return &measurementRepository{GormDAL: customDB}
}

func (mr *measurementRepository) Create(ctx *context.Context, measurement *entities.Measurement) *errs.XError {
	res := mr.WithDB(ctx).Create(&measurement)
	if res.Error != nil {
		return errs.NewXError(errs.DATABASE, "Unable to save measurement", res.Error)
	}
	return nil
}

func (mr *measurementRepository) BatchCreate(ctx *context.Context, measurements []*entities.Measurement) *errs.XError {
	if len(measurements) == 0 {
		return nil
	}

	res := mr.txn.Txn(ctx).CreateInBatches(measurements, 100)
	if res.Error != nil {
		return errs.NewXError(errs.DATABASE, "Unable to batch save measurements", res.Error)
	}
	return nil
}

func (mr *measurementRepository) Update(ctx *context.Context, measurement *entities.Measurement) *errs.XError {
	return mr.GormDAL.Update(ctx, *measurement)
}

func (mr *measurementRepository) BatchUpdate(ctx *context.Context, measurements []*entities.Measurement) *errs.XError {
	if len(measurements) == 0 {
		return nil
	}

	for _, measurement := range measurements {
		if measurement.ID == 0 {
			continue
		}
		err := mr.customDB.Update(ctx, *measurement)
		if err != nil {
			return errs.NewXError(errs.DATABASE, "Unable to batch update measurements", err)
		}
	}

	return nil
}

func (mr *measurementRepository) Get(ctx *context.Context, id uint) (*entities.Measurement, *errs.XError) {
	measurement := entities.Measurement{}
	res := mr.WithDB(ctx).
		Preload("Person").
		Preload("DressType").
		Preload("TakenBy", scopes.SelectFields("first_name", "last_name")).
		Find(&measurement, id)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find measurement", res.Error)
	}
	return &measurement, nil
}

func (mr *measurementRepository) GetByPersonIdAndDressTypeId(ctx *context.Context, personId uint, dressTypeId uint) (*entities.Measurement, *errs.XError) {
	measurement := entities.Measurement{}
	res := mr.txn.Txn(ctx).
		Where("person_id = ? AND dress_type_id = ?", personId, dressTypeId).
		First(&measurement)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errs.NewXError(errs.DATABASE, "Unable to find measurement", res.Error)
	}
	return &measurement, nil
}

func (mr *measurementRepository) GetAll(ctx *context.Context, search string) ([]entities.Measurement, *errs.XError) {
	var measurements []entities.Measurement

	filterValue := util.ReadValueFromContext(ctx, constants.FILTER_KEY)
	var filter string
	if filterValue != nil {
		filter = filterValue.(string)
	}

	res := mr.txn.Txn(ctx).
		Scopes(scopes.Channel(), scopes.IsActive()).
		Scopes(scopes.GetMeasurements_Search(search)).
		Scopes(scopes.GetMeasurements_Filter(filter)).
		Scopes(db.Paginate(ctx)).
		Preload("Person").
		Preload("DressType").
		Preload("TakenBy", scopes.SelectFields("first_name", "last_name")).
		Find(&measurements)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find measurements", res.Error)
	}
	return measurements, nil
}

func (mr *measurementRepository) Delete(ctx *context.Context, id uint) *errs.XError {
	measurement := &entities.Measurement{Model: &entities.Model{ID: id, IsActive: false}}
	err := mr.GormDAL.Delete(ctx, measurement)
	if err != nil {
		return err
	}
	return nil
}
