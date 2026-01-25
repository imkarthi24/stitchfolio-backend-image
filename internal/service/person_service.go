package service

import (
	"context"

	"github.com/imkarthi24/sf-backend/internal/mapper"
	requestModel "github.com/imkarthi24/sf-backend/internal/model/request"
	responseModel "github.com/imkarthi24/sf-backend/internal/model/response"
	"github.com/imkarthi24/sf-backend/internal/repository"
	"github.com/loop-kar/pixie/errs"
)

type PersonService interface {
	SavePerson(*context.Context, requestModel.Person) *errs.XError
	UpdatePerson(*context.Context, requestModel.Person, uint) *errs.XError
	Get(*context.Context, uint) (*responseModel.Person, *errs.XError)
	GetAll(*context.Context, string) ([]responseModel.Person, *errs.XError)
	Delete(*context.Context, uint) *errs.XError
	GetByCustomerId(*context.Context, uint) ([]responseModel.Person, *errs.XError)
}

type personService struct {
	personRepo repository.PersonRepository
	mapper     mapper.Mapper
	respMapper mapper.ResponseMapper
}

func ProvidePersonService(repo repository.PersonRepository, mapper mapper.Mapper, respMapper mapper.ResponseMapper) PersonService {
	return personService{
		personRepo: repo,
		mapper:     mapper,
		respMapper: respMapper,
	}
}

func (svc personService) SavePerson(ctx *context.Context, person requestModel.Person) *errs.XError {
	dbPerson, err := svc.mapper.Person(person)
	if err != nil {
		return errs.NewXError(errs.INVALID_REQUEST, "Unable to save person", err)
	}

	errr := svc.personRepo.Create(ctx, dbPerson)
	if errr != nil {
		return errr
	}

	return nil
}

func (svc personService) UpdatePerson(ctx *context.Context, person requestModel.Person, id uint) *errs.XError {
	dbPerson, err := svc.mapper.Person(person)
	if err != nil {
		return errs.NewXError(errs.INVALID_REQUEST, "Unable to update person", err)
	}

	dbPerson.ID = id
	errr := svc.personRepo.Update(ctx, dbPerson)
	if errr != nil {
		return errr
	}
	return nil
}

func (svc personService) Get(ctx *context.Context, id uint) (*responseModel.Person, *errs.XError) {
	person, err := svc.personRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	mappedPerson, mapErr := svc.respMapper.Person(person)
	if mapErr != nil {
		return nil, errs.NewXError(errs.MAPPING_ERROR, "Failed to map Person data", mapErr)
	}

	return mappedPerson, nil
}

func (svc personService) GetAll(ctx *context.Context, search string) ([]responseModel.Person, *errs.XError) {
	persons, err := svc.personRepo.GetAll(ctx, search)
	if err != nil {
		return nil, err
	}

	mappedPersons, mapErr := svc.respMapper.Persons(persons)
	if mapErr != nil {
		return nil, errs.NewXError(errs.MAPPING_ERROR, "Failed to map Person data", mapErr)
	}

	return mappedPersons, nil
}

func (svc personService) Delete(ctx *context.Context, id uint) *errs.XError {
	err := svc.personRepo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (svc personService) GetByCustomerId(ctx *context.Context, customerId uint) ([]responseModel.Person, *errs.XError) {
	persons, err := svc.personRepo.GetByCustomerId(ctx, customerId)
	if err != nil {
		return nil, err
	}

	mappedPersons, mapErr := svc.respMapper.Persons(persons)
	if mapErr != nil {
		return nil, errs.NewXError(errs.MAPPING_ERROR, "Failed to map Person data", mapErr)
	}

	return mappedPersons, nil
}
