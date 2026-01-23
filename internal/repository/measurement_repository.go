package repository

import (
	"context"
	"fmt"

	"github.com/imkarthi24/sf-backend/internal/entities"
	"github.com/imkarthi24/sf-backend/internal/repository/common"
	"github.com/imkarthi24/sf-backend/internal/repository/scopes"
	"github.com/imkarthi24/sf-backend/pkg/db"
	"github.com/imkarthi24/sf-backend/pkg/errs"
	"github.com/imkarthi24/sf-backend/pkg/util"
)

type MeasurementRepository interface {
	Create(*context.Context, *entities.Measurement) *errs.XError
	BatchCreate(*context.Context, []*entities.Measurement) *errs.XError
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

func (mr *measurementRepository) BatchCreate(ctx *context.Context, measurements []*entities.Measurement) *errs.XError {
	if len(measurements) == 0 {
		return nil
	}

	// Convert []*entities.Measurement to []interface{} for batch create
	measurementsInterface := make([]interface{}, len(measurements))
	for i, m := range measurements {
		measurementsInterface[i] = m
	}

	res := mr.txn.Txn(ctx).CreateInBatches(measurementsInterface, 100)
	if res.Error != nil {
		return errs.NewXError(errs.DATABASE, "Unable to batch save measurements", res.Error)
	}
	return nil
}

func (mr *measurementRepository) Update(ctx *context.Context, measurement *entities.Measurement) *errs.XError {
	return mr.customDB.Update(ctx, *measurement)
}

func (mr *measurementRepository) Get(ctx *context.Context, id uint) (*entities.Measurement, *errs.XError) {
	measurement := entities.Measurement{}
	res := mr.txn.Txn(ctx).
		Preload("Person").
		Preload("DressType").
		Preload("TakenBy", scopes.SelectFields("first_name", "last_name")).
		Find(&measurement, id)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find measurement", res.Error)
	}
	return &measurement, nil
}

func (mr *measurementRepository) GetAll(ctx *context.Context, search string) ([]entities.Measurement, *errs.XError) {
	var measurements []entities.Measurement
	query := mr.txn.Txn(ctx).
		Scopes(scopes.Channel(), scopes.IsActive())

	if !util.IsNilOrEmptyString(&search) {
		formattedSearch := util.EncloseWithPercentageOperator(search)
		whereClause := fmt.Sprintf(
			"(dress_type ILIKE %s OR measurement_by ILIKE %s OR EXISTS (SELECT 1 FROM \"stitch\".\"Customers\" WHERE \"Customers\".id = \"stitch\".\"Measurements\".customer_id AND (\"Customers\".phone_number ILIKE %s OR \"Customers\".first_name ILIKE %s OR \"Customers\".last_name ILIKE %s)))",
			formattedSearch, formattedSearch, formattedSearch, formattedSearch, formattedSearch,
		)
		query = query.Where(whereClause)
	}

	res := query.
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
	err := mr.customDB.Delete(ctx, measurement)
	if err != nil {
		return err
	}
	return nil
}
