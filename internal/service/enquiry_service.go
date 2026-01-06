package service

import (
	"context"

	"github.com/imkarthi24/sf-backend/internal/mapper"
	requestModel "github.com/imkarthi24/sf-backend/internal/model/request"
	responseModel "github.com/imkarthi24/sf-backend/internal/model/response"
	"github.com/imkarthi24/sf-backend/internal/repository"
	"github.com/imkarthi24/sf-backend/pkg/errs"
)

type EnquiryService interface {
	SaveEnquiry(*context.Context, requestModel.Enquiry) *errs.XError
	UpdateEnquiry(*context.Context, requestModel.Enquiry, uint) *errs.XError
	UpdateEnquiryCustomer(*context.Context, uint, uint) *errs.XError
	Get(*context.Context, uint) (*responseModel.Enquiry, *errs.XError)
	GetAll(*context.Context, string) ([]responseModel.Enquiry, *errs.XError)
	Delete(*context.Context, uint) *errs.XError
}

type enquiryService struct {
	enquiryRepo  repository.EnquiryRepository
	customerRepo repository.CustomerRepository
	mapper       mapper.Mapper
	respMapper   mapper.ResponseMapper
}

func ProvideEnquiryService(repo repository.EnquiryRepository, customerRepo repository.CustomerRepository, mapper mapper.Mapper, respMapper mapper.ResponseMapper) EnquiryService {
	return enquiryService{
		enquiryRepo:  repo,
		customerRepo: customerRepo,
		mapper:       mapper,
		respMapper:   respMapper,
	}
}

func (svc enquiryService) SaveEnquiry(ctx *context.Context, enquiry requestModel.Enquiry) *errs.XError {
	var customerId *uint

	if enquiry.PhoneNumber != "" {
		existingCustomer, err := svc.customerRepo.GetByPhoneNumber(ctx, enquiry.PhoneNumber)
		if err != nil {
			return err
		}

		if existingCustomer != nil {
			customerId = &existingCustomer.ID
		} else {
			customerRequest := requestModel.Customer{
				Name:           enquiry.Name,
				Email:          enquiry.Email,
				PhoneNumber:    enquiry.PhoneNumber,
				WhatsappNumber: enquiry.WhatsappNumber,
				Address:        enquiry.Address,
				IsActive:       false,
			}
			dbCustomer, mapErr := svc.mapper.Customer(customerRequest)
			if mapErr != nil {
				return errs.NewXError(errs.INVALID_REQUEST, "Unable to map customer", mapErr)
			}

			createErr := svc.customerRepo.Create(ctx, dbCustomer)
			if createErr != nil {
				return createErr
			}
			customerId = &dbCustomer.ID

			dbCustomer.IsActive = false
			updateErr := svc.customerRepo.Update(ctx, dbCustomer)
			if updateErr != nil {
				return updateErr
			}
		}
	} else if enquiry.CustomerId != nil {
		customerId = enquiry.CustomerId
	}

	dbEnquiry, err := svc.mapper.Enquiry(enquiry)
	if err != nil {
		return errs.NewXError(errs.INVALID_REQUEST, "Unable to save enquiry", err)
	}
	dbEnquiry.CustomerId = customerId

	errr := svc.enquiryRepo.Create(ctx, dbEnquiry)
	if errr != nil {
		return errr
	}

	return nil
}

func (svc enquiryService) UpdateEnquiry(ctx *context.Context, enquiry requestModel.Enquiry, id uint) *errs.XError {
	dbEnquiry, err := svc.mapper.Enquiry(enquiry)
	if err != nil {
		return errs.NewXError(errs.INVALID_REQUEST, "Unable to update enquiry", err)
	}

	dbEnquiry.ID = id
	errr := svc.enquiryRepo.Update(ctx, dbEnquiry)
	if errr != nil {
		return errr
	}
	return nil
}

func (svc enquiryService) UpdateEnquiryCustomer(ctx *context.Context, enquiryId uint, customerId uint) *errs.XError {
	enquiry, err := svc.enquiryRepo.Get(ctx, enquiryId)
	if err != nil {
		return err
	}

	_, customerErr := svc.customerRepo.Get(ctx, customerId)
	if customerErr != nil {
		return customerErr
	}

	customerIdPtr := &customerId
	enquiry.CustomerId = customerIdPtr
	errr := svc.enquiryRepo.Update(ctx, enquiry)
	if errr != nil {
		return errr
	}
	return nil
}

func (svc enquiryService) Get(ctx *context.Context, id uint) (*responseModel.Enquiry, *errs.XError) {
	enquiry, err := svc.enquiryRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	mappedEnquiry, mapErr := svc.respMapper.Enquiry(enquiry)
	if mapErr != nil {
		return nil, errs.NewXError(errs.MAPPING_ERROR, "Failed to map Enquiry data", mapErr)
	}

	return mappedEnquiry, nil
}

func (svc enquiryService) GetAll(ctx *context.Context, search string) ([]responseModel.Enquiry, *errs.XError) {
	enquiries, err := svc.enquiryRepo.GetAll(ctx, search)
	if err != nil {
		return nil, err
	}

	mappedEnquiries, mapErr := svc.respMapper.Enquiries(enquiries)
	if mapErr != nil {
		return nil, errs.NewXError(errs.MAPPING_ERROR, "Failed to map Enquiry data", mapErr)
	}

	return mappedEnquiries, nil
}

func (svc enquiryService) Delete(ctx *context.Context, id uint) *errs.XError {
	err := svc.enquiryRepo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
