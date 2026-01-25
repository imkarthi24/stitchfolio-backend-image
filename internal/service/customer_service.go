package service

import (
	"context"

	"github.com/imkarthi24/sf-backend/internal/entities"
	"github.com/imkarthi24/sf-backend/internal/mapper"
	requestModel "github.com/imkarthi24/sf-backend/internal/model/request"
	responseModel "github.com/imkarthi24/sf-backend/internal/model/response"
	"github.com/imkarthi24/sf-backend/internal/repository"
	"github.com/loop-kar/pixie/errs"
)

type CustomerService interface {
	SaveCustomer(*context.Context, requestModel.Customer) *errs.XError
	UpdateCustomer(*context.Context, requestModel.Customer, uint) *errs.XError
	Get(*context.Context, uint) (*responseModel.Customer, *errs.XError)
	GetAll(*context.Context, string) ([]responseModel.Customer, *errs.XError)
	Delete(*context.Context, uint) *errs.XError
	AutocompleteCustomer(*context.Context, string) ([]responseModel.CustomerAutoComplete, *errs.XError)
}

type customerService struct {
	customerRepo repository.CustomerRepository
	personRepo   repository.PersonRepository
	mapper       mapper.Mapper
	respMapper   mapper.ResponseMapper
}

func ProvideCustomerService(repo repository.CustomerRepository, personRepo repository.PersonRepository, mapper mapper.Mapper, respMapper mapper.ResponseMapper) CustomerService {
	return customerService{
		customerRepo: repo,
		personRepo:   personRepo,
		mapper:       mapper,
		respMapper:   respMapper,
	}
}

func (svc customerService) SaveCustomer(ctx *context.Context, customer requestModel.Customer) *errs.XError {
	dbCustomer, err := svc.mapper.Customer(customer)
	if err != nil {
		return errs.NewXError(errs.INVALID_REQUEST, "Unable to save customer", err)
	}

	errr := svc.customerRepo.Create(ctx, dbCustomer)
	if errr != nil {
		return errr
	}

	// Create a person with the customer's name (first name + last name)
	person := &entities.Person{
		Model:      &entities.Model{IsActive: true},
		FirstName:  customer.FirstName,
		LastName:   customer.LastName,
		CustomerId: dbCustomer.ID,
		Age:        &customer.Age,
		Gender:     entities.Gender(customer.Gender),
	}

	errr = svc.personRepo.Create(ctx, person)
	if errr != nil {
		return errr
	}

	return nil
}

func (svc customerService) UpdateCustomer(ctx *context.Context, customer requestModel.Customer, id uint) *errs.XError {
	dbCustomer, err := svc.mapper.Customer(customer)
	if err != nil {
		return errs.NewXError(errs.INVALID_REQUEST, "Unable to update customer", err)
	}

	dbCustomer.ID = id
	errr := svc.customerRepo.Update(ctx, dbCustomer)
	if errr != nil {
		return errr
	}
	return nil
}

func (svc customerService) Get(ctx *context.Context, id uint) (*responseModel.Customer, *errs.XError) {
	customer, err := svc.customerRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	mappedCustomer, mapErr := svc.respMapper.Customer(customer)
	if mapErr != nil {
		return nil, errs.NewXError(errs.MAPPING_ERROR, "Failed to map Customer data", mapErr)
	}

	return mappedCustomer, nil
}

func (svc customerService) GetAll(ctx *context.Context, search string) ([]responseModel.Customer, *errs.XError) {
	customers, err := svc.customerRepo.GetAll(ctx, search)
	if err != nil {
		return nil, err
	}

	mappedCustomers, mapErr := svc.respMapper.Customers(customers)
	if mapErr != nil {
		return nil, errs.NewXError(errs.MAPPING_ERROR, "Failed to map Customer data", mapErr)
	}

	return mappedCustomers, nil
}

func (svc customerService) Delete(ctx *context.Context, id uint) *errs.XError {
	err := svc.customerRepo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (svc customerService) AutocompleteCustomer(ctx *context.Context, search string) ([]responseModel.CustomerAutoComplete, *errs.XError) {
	customers, err := svc.customerRepo.AutocompleteCustomer(ctx, search)
	if err != nil {
		return nil, err
	}

	res := make([]responseModel.CustomerAutoComplete, 0)
	for _, customer := range customers {
		res = append(res, responseModel.CustomerAutoComplete{
			CustomerID:  customer.ID,
			FirstName:   customer.FirstName,
			LastName:    customer.LastName,
			Email:       customer.Email,
			PhoneNumber: customer.PhoneNumber,
		})
	}

	return res, nil
}
