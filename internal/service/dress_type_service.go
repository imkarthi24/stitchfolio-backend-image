package service

import (
	"context"

	"github.com/imkarthi24/sf-backend/internal/mapper"
	requestModel "github.com/imkarthi24/sf-backend/internal/model/request"
	responseModel "github.com/imkarthi24/sf-backend/internal/model/response"
	"github.com/imkarthi24/sf-backend/internal/repository"
	"github.com/loop-kar/pixie/errs"
)

type DressTypeService interface {
	SaveDressType(*context.Context, requestModel.DressType) *errs.XError
	UpdateDressType(*context.Context, requestModel.DressType, uint) *errs.XError
	Get(*context.Context, uint) (*responseModel.DressType, *errs.XError)
	GetAll(*context.Context, string) ([]responseModel.DressType, *errs.XError)
	Delete(*context.Context, uint) *errs.XError
}

type dressTypeService struct {
	dressTypeRepo repository.DressTypeRepository
	mapper        mapper.Mapper
	respMapper    mapper.ResponseMapper
}

func ProvideDressTypeService(repo repository.DressTypeRepository, mapper mapper.Mapper, respMapper mapper.ResponseMapper) DressTypeService {
	return dressTypeService{
		dressTypeRepo: repo,
		mapper:        mapper,
		respMapper:    respMapper,
	}
}

func (svc dressTypeService) SaveDressType(ctx *context.Context, dressType requestModel.DressType) *errs.XError {
	dbDressType, err := svc.mapper.DressType(dressType)
	if err != nil {
		return errs.NewXError(errs.INVALID_REQUEST, "Unable to save dress type", err)
	}

	errr := svc.dressTypeRepo.Create(ctx, dbDressType)
	if errr != nil {
		return errr
	}

	return nil
}

func (svc dressTypeService) UpdateDressType(ctx *context.Context, dressType requestModel.DressType, id uint) *errs.XError {
	dbDressType, err := svc.mapper.DressType(dressType)
	if err != nil {
		return errs.NewXError(errs.INVALID_REQUEST, "Unable to update dress type", err)
	}

	dbDressType.ID = id
	errr := svc.dressTypeRepo.Update(ctx, dbDressType)
	if errr != nil {
		return errr
	}
	return nil
}

func (svc dressTypeService) Get(ctx *context.Context, id uint) (*responseModel.DressType, *errs.XError) {
	dressType, err := svc.dressTypeRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	mappedDressType, mapErr := svc.respMapper.DressType(dressType)
	if mapErr != nil {
		return nil, errs.NewXError(errs.MAPPING_ERROR, "Failed to map DressType data", mapErr)
	}

	return mappedDressType, nil
}

func (svc dressTypeService) GetAll(ctx *context.Context, search string) ([]responseModel.DressType, *errs.XError) {
	dressTypes, err := svc.dressTypeRepo.GetAll(ctx, search)
	if err != nil {
		return nil, err
	}

	mappedDressTypes, mapErr := svc.respMapper.DressTypes(dressTypes)
	if mapErr != nil {
		return nil, errs.NewXError(errs.MAPPING_ERROR, "Failed to map DressType data", mapErr)
	}

	return mappedDressTypes, nil
}

func (svc dressTypeService) Delete(ctx *context.Context, id uint) *errs.XError {
	err := svc.dressTypeRepo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
