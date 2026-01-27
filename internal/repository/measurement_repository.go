package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/imkarthi24/sf-backend/internal/entities"
	"github.com/imkarthi24/sf-backend/internal/repository/scopes"
	"github.com/loop-kar/pixie/db"
	"github.com/loop-kar/pixie/errs"
	"github.com/loop-kar/pixie/util"

	"gorm.io/gorm"
)

type MeasurementRepository interface {
	Create(*context.Context, *entities.Measurement) *errs.XError
	BatchCreate(*context.Context, []*entities.Measurement) *errs.XError
	Update(*context.Context, *entities.Measurement) *errs.XError
	BatchUpdate(*context.Context, []*entities.Measurement) *errs.XError
	Get(*context.Context, uint) (*entities.Measurement, *errs.XError)
	GetByPersonIdAndDressTypeId(*context.Context, uint, uint) (*entities.Measurement, *errs.XError)
	GetAll(*context.Context, string) ([]entities.GroupedMeasurement, *errs.XError)
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

	res := mr.WithDB(ctx).CreateInBatches(measurements, 100)
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
		err := mr.GormDAL.Update(ctx, *measurement)
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
	res := mr.WithDB(ctx).
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

func (mr *measurementRepository) GetAll(ctx *context.Context, search string) ([]entities.GroupedMeasurement, *errs.XError) {
	var groupedMeasurements []entities.GroupedMeasurement

	// filterValue := util.ReadValueFromContext(ctx, constants.FILTER_KEY)
	// var filter string
	// if filterValue != nil {
	// 	filter = filterValue.(string)
	// }

	query := mr.WithDB(ctx).Table(entities.Measurement{}.TableNameForQuery()).
		Scopes(scopes.IsActive("E"), scopes.Channel("E")).
		Scopes(db.Paginate(ctx)).
		Select(`E.person_id, 
			STRING_AGG(DISTINCT DT.name, ',' ORDER BY DT.name) as dress_types`).
		Joins(`INNER JOIN "stich"."DressTypes" DT ON DT.id = E.dress_type_id`).
		Joins(`INNER JOIN "stich"."Persons" P ON P.id = E.person_id`).
		Group("person_id")

	if !util.IsNilOrEmptyString(&search) {
		formattedSearch := util.EncloseWithPercentageOperator(search)
		whereClause := fmt.Sprintf(
			`(C.first_name ILIKE %s OR C.last_name ILIKE %s OR CONCAT(C.first_name, ' ', C.last_name) ILIKE %s)`,
			formattedSearch, formattedSearch, formattedSearch,
		)
		query = query.Joins(`INNER JOIN "stich"."Customers" C ON C.id = P.customer_id`).
			Where(whereClause)
	}

	res := query.Scan(&groupedMeasurements)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "UNABLE_TO_FIND_MEASUREMENTS", res.Error)
	}
	return groupedMeasurements, nil
}

func (mr *measurementRepository) Delete(ctx *context.Context, id uint) *errs.XError {
	measurement := &entities.Measurement{Model: &entities.Model{ID: id, IsActive: false}}
	err := mr.GormDAL.Delete(ctx, measurement)
	if err != nil {
		return err
	}
	return nil
}
