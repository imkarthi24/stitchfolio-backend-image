package repository

import (
	"context"

	"github.com/imkarthi24/sf-backend/internal/entities"
	"github.com/imkarthi24/sf-backend/internal/repository/common"
	"github.com/imkarthi24/sf-backend/internal/repository/scopes"
	"github.com/imkarthi24/sf-backend/pkg/db"
	"github.com/imkarthi24/sf-backend/pkg/errs"
)

type EnquiryRepository interface {
	Create(*context.Context, *entities.Enquiry) *errs.XError
	Update(*context.Context, *entities.Enquiry) *errs.XError
	Get(*context.Context, uint) (*entities.Enquiry, *errs.XError)
	GetAll(*context.Context, string) ([]entities.Enquiry, *errs.XError)
	Delete(*context.Context, uint) *errs.XError
}

type enquiryRepository struct {
	txn      db.DBTransactionManager
	customDB common.CustomGormDB
}

func ProvideEnquiryRepository(txn db.DBTransactionManager, customDB common.CustomGormDB) EnquiryRepository {
	return &enquiryRepository{txn: txn, customDB: customDB}
}

func (er *enquiryRepository) Create(ctx *context.Context, enquiry *entities.Enquiry) *errs.XError {
	res := er.txn.Txn(ctx).Create(&enquiry)
	if res.Error != nil {
		return errs.NewXError(errs.DATABASE, "Unable to save enquiry", res.Error)
	}
	return nil
}

func (er *enquiryRepository) Update(ctx *context.Context, enquiry *entities.Enquiry) *errs.XError {
	return er.customDB.Update(ctx, *enquiry)
}

func (er *enquiryRepository) Get(ctx *context.Context, id uint) (*entities.Enquiry, *errs.XError) {
	enquiry := entities.Enquiry{}
	res := er.txn.Txn(ctx).Preload("Customer").Find(&enquiry, id)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find enquiry", res.Error)
	}
	return &enquiry, nil
}

func (er *enquiryRepository) GetAll(ctx *context.Context, search string) ([]entities.Enquiry, *errs.XError) {
	var enquiries []entities.Enquiry
	res := er.txn.Txn(ctx).
		Scopes(scopes.Channel(), scopes.IsActive()).
		Scopes(scopes.ILike(search, "subject", "notes", "status")).
		Scopes(db.Paginate(ctx)).
		Preload("Customer").
		Find(&enquiries)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find enquiries", res.Error)
	}
	return enquiries, nil
}

func (er *enquiryRepository) Delete(ctx *context.Context, id uint) *errs.XError {
	enquiry := &entities.Enquiry{Model: &entities.Model{ID: id, IsActive: false}}
	err := er.customDB.Delete(ctx, enquiry)
	if err != nil {
		return err
	}
	return nil
}
