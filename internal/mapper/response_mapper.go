package mapper

import (
	"encoding/json"
	"time"

	"github.com/imkarthi24/sf-backend/internal/entities"
	responseModel "github.com/imkarthi24/sf-backend/internal/model/response"
	"github.com/imkarthi24/sf-backend/pkg/util"
)

type responseMapper struct{}

type ResponseMapper interface {
	UserBrowse([]entities.User) []responseModel.User
	User(*entities.User) (*responseModel.User, error)

	Channels([]entities.Channel) []responseModel.Channel
	Channel(*entities.Channel) *responseModel.Channel

	Enquiry(e *entities.Enquiry) (*responseModel.Enquiry, error)
	Enquiries(enquiries []entities.Enquiry) ([]responseModel.Enquiry, error)

	EnquiryHistory(e *entities.EnquiryHistory) (*responseModel.EnquiryHistory, error)
	EnquiryHistories(enquiryHistories []entities.EnquiryHistory) ([]responseModel.EnquiryHistory, error)

	MasterConfig(e *entities.MasterConfig) (*responseModel.MasterConfig, error)
	MasterConfigs(items []entities.MasterConfig) ([]responseModel.MasterConfig, error)

	Customer(e *entities.Customer) (*responseModel.Customer, error)
	Customers(items []entities.Customer) ([]responseModel.Customer, error)
	Measurement(e *entities.Measurement) (*responseModel.Measurement, error)
	Measurements(items []entities.Measurement) ([]responseModel.Measurement, error)
	Order(e *entities.Order) (*responseModel.Order, error)
	Orders(items []entities.Order) ([]responseModel.Order, error)
	OrderItem(e *entities.OrderItem) (*responseModel.OrderItem, error)
	OrderItems(items []entities.OrderItem) ([]responseModel.OrderItem, error)
}

func ProvideResponseMapper() ResponseMapper {
	return &responseMapper{}
}

func (*responseMapper) Channel(channel *entities.Channel) *responseModel.Channel {

	return &responseModel.Channel{
		Channel:               channel,
		ChannelOwnerFirstName: channel.OwnerUser.FirstName,
		ChannelOwnerLastName:  channel.OwnerUser.LastName,
		PhoneNumber:           channel.OwnerUser.PhoneNumber,
		Email:                 channel.OwnerUser.Email,
	}
}

func (m *responseMapper) Channels(channels []entities.Channel) []responseModel.Channel {
	res := make([]responseModel.Channel, 0)
	for _, chnl := range channels {
		res = append(res, *m.Channel(&chnl))
	}

	return res
}

func (m *responseMapper) UserBrowse(users []entities.User) []responseModel.User {

	res := make([]responseModel.User, 0)
	for _, user := range users {
		mappedUser, _ := m.User(&user)
		res = append(res, *mappedUser)
	}

	return res
}

func (m *responseMapper) User(usr *entities.User) (*responseModel.User, error) {
	if usr == nil {
		return nil, nil
	}

	return &responseModel.User{
		ID:                  usr.ID,
		IsActive:            usr.IsActive,
		FirstName:           usr.FirstName,
		LastName:            usr.LastName,
		Extension:           usr.Extension,
		PhoneNumber:         usr.PhoneNumber,
		Email:               usr.Email,
		Role:                string(usr.Role),
		IsLoginDisabled:     usr.IsLoginDisabled,
		IsLoggedIn:          usr.IsLoggedIn,
		LastLoginTime:       usr.LastLoginTime,
		LoginFailureCounter: usr.LoginFailureCounter,
		Experience:          usr.Experience,
		Department:          usr.Department,
	}, nil
}

func (m *responseMapper) EnquiryHistories(enquiryHistories []entities.EnquiryHistory) ([]responseModel.EnquiryHistory, error) {
	if len(enquiryHistories) == 0 || enquiryHistories == nil {
		return nil, nil
	}

	histories := make([]responseModel.EnquiryHistory, 0)

	for _, history := range enquiryHistories {
		mappedHistory, err := m.EnquiryHistory(&history)
		if err != nil {
			return nil, err
		}
		histories = append(histories, *mappedHistory)
	}
	return histories, nil
}

func (m *responseMapper) Enquiries(enquiries []entities.Enquiry) ([]responseModel.Enquiry, error) {
	result := make([]responseModel.Enquiry, 0)

	for _, enquiry := range enquiries {
		mappedEnquiry, err := m.Enquiry(&enquiry)
		if err != nil {
			return nil, err
		}
		result = append(result, *mappedEnquiry)
	}

	return result, nil
}

func (m *responseMapper) Enquiry(e *entities.Enquiry) (*responseModel.Enquiry, error) {
	if e == nil {
		return nil, nil
	}

	customer, err := m.Customer(e.Customer)
	if err != nil {
		return nil, err
	}

	return &responseModel.Enquiry{
		ID:                  e.ID,
		IsActive:            e.IsActive,
		Subject:             e.Subject,
		Notes:               e.Notes,
		Status:              string(e.Status),
		CustomerId:          e.CustomerId,
		Customer:            customer,
		Source:              e.Source,
		ReferredBy:          e.ReferredBy,
		ReferrerPhoneNumber: e.ReferrerPhoneNumber,
	}, nil
}

func (m *responseMapper) EnquiryHistory(e *entities.EnquiryHistory) (*responseModel.EnquiryHistory, error) {
	if e == nil {
		return nil, nil
	}

	var employee *responseModel.User
	var err error
	if e.Employee != nil {
		employee, err = m.User(e.Employee)
		if err != nil {
			return nil, err
		}
	}

	var visitingDateStr *string
	if e.VisitingDate != nil {
		str := util.DateTimeToStringOrDefault(e.VisitingDate, time.DateOnly)
		visitingDateStr = &str
	}

	var callBackDateStr *string
	if e.CallBackDate != nil {
		str := util.DateTimeToStringOrDefault(e.CallBackDate, time.DateOnly)
		callBackDateStr = &str
	}

	return &responseModel.EnquiryHistory{
		ID:              e.ID,
		IsActive:        e.IsActive,
		EmployeeComment: e.EmployeeComment,
		CustomerComment: e.CustomerComment,
		VisitingDate:    visitingDateStr,
		CallBackDate:    callBackDateStr,
		EnquiryDate:     util.DateTimeToStringOrDefault(&e.EnquiryDate, time.DateOnly),
		ResponseStatus:  string(e.ResponseStatus),
		EnquiryId:       e.EnquiryId,
		EmployeeId:      e.EmployeeId,
		Employee:        employee,
	}, nil
}
func (m *responseMapper) MasterConfig(e *entities.MasterConfig) (*responseModel.MasterConfig, error) {
	return &responseModel.MasterConfig{
		Id:            e.ID,
		IsActive:      e.IsActive,
		Name:          e.Name,
		Type:          e.Type,
		CurrentValue:  e.CurrentValue,
		DefaultValue:  e.DefaultValue,
		UseDefault:    e.UseDefault,
		PreviousValue: e.PreviousValue,
		Description:   e.Description,
		Format:        e.Format,
	}, nil
}

func (m *responseMapper) MasterConfigs(items []entities.MasterConfig) ([]responseModel.MasterConfig, error) {
	var mappedItems []responseModel.MasterConfig
	for _, item := range items {
		mappedItem, err := m.MasterConfig(&item)
		if err != nil {
			return nil, err
		}
		mappedItems = append(mappedItems, *mappedItem)
	}

	return mappedItems, nil
}

func (m *responseMapper) Customer(e *entities.Customer) (*responseModel.Customer, error) {
	if e == nil {
		return nil, nil
	}

	enquiries, err := m.Enquiries(e.Enquiries)
	if err != nil {
		return nil, err
	}

	measurements, err := m.Measurements(e.Measurements)
	if err != nil {
		return nil, err
	}

	orders, err := m.Orders(e.Orders)
	if err != nil {
		return nil, err
	}

	return &responseModel.Customer{
		ID:             e.ID,
		IsActive:       e.IsActive,
		FirstName:      e.FirstName,
		LastName:       e.LastName,
		Email:          e.Email,
		PhoneNumber:    e.PhoneNumber,
		WhatsappNumber: e.WhatsappNumber,
		Address:        e.Address,
		Enquiries:      enquiries,
		Measurements:   measurements,
		Orders:         orders,
	}, nil
}

func (m *responseMapper) Customers(items []entities.Customer) ([]responseModel.Customer, error) {
	result := make([]responseModel.Customer, 0)
	for _, item := range items {
		mappedItem, err := m.Customer(&item)
		if err != nil {
			return nil, err
		}
		result = append(result, *mappedItem)
	}
	return result, nil
}

func (m *responseMapper) Measurement(e *entities.Measurement) (*responseModel.Measurement, error) {
	if e == nil {
		return nil, nil
	}

	customer, err := m.Customer(e.Customer)
	if err != nil {
		return nil, err
	}

	measurementTakenBy, err := m.User(e.MeasurementTakenBy)
	if err != nil {
		return nil, err
	}

	return &responseModel.Measurement{
		ID:              e.ID,
		IsActive:        e.IsActive,
		MeasurementDate: util.DateTimeToStringOrDefault(e.MeasurementDate, time.DateOnly),
		MeasurementBy:   e.MeasurementBy,
		DressType:       e.DressType,
		Measurements:    json.RawMessage(e.Measurements),
		CustomerId:      e.CustomerId,
		// Customer:             customer,
		MeasurementTakenById: e.MeasurementTakenById,
		// MeasurementTakenBy:   measurementTakenBy,
		UpdatedAt:          e.UpdatedAt,
		UpdatedById:        e.UpdatedById,
		CustomerName:       customer.FirstName + " " + customer.LastName,
		MeasurementTakenBy: measurementTakenBy.FirstName + " " + measurementTakenBy.LastName,
	}, nil
}

func (m *responseMapper) Measurements(items []entities.Measurement) ([]responseModel.Measurement, error) {
	result := make([]responseModel.Measurement, 0)
	for _, item := range items {
		mappedItem, err := m.Measurement(&item)
		if err != nil {
			return nil, err
		}
		result = append(result, *mappedItem)
	}
	return result, nil
}

func (m *responseMapper) Order(e *entities.Order) (*responseModel.Order, error) {
	if e == nil {
		return nil, nil
	}

	customer, err := m.Customer(e.Customer)
	if err != nil {
		return nil, err
	}

	orderItems, err := m.OrderItems(e.OrderItems)
	if err != nil {
		return nil, err
	}

	return &responseModel.Order{
		ID:         e.ID,
		IsActive:   e.IsActive,
		Status:     string(e.Status),
		CustomerId: e.CustomerId,
		Customer:   customer,
		OrderItems: orderItems,
	}, nil
}

func (m *responseMapper) Orders(items []entities.Order) ([]responseModel.Order, error) {
	result := make([]responseModel.Order, 0)
	for _, item := range items {
		mappedItem, err := m.Order(&item)
		if err != nil {
			return nil, err
		}
		result = append(result, *mappedItem)
	}
	return result, nil
}

func (m *responseMapper) OrderItem(e *entities.OrderItem) (*responseModel.OrderItem, error) {
	if e == nil {
		return nil, nil
	}

	order, err := m.Order(e.Order)
	if err != nil {
		return nil, err
	}

	measurement, err := m.Measurement(e.Measurement)
	if err != nil {
		return nil, err
	}

	return &responseModel.OrderItem{
		ID:            e.ID,
		IsActive:      e.IsActive,
		Description:   e.Description,
		Quantity:      e.Quantity,
		Price:         e.Price,
		Total:         e.Total,
		OrderId:       e.OrderId,
		Order:         order,
		MeasurementId: e.MeasurementId,
		Measurement:   measurement,
	}, nil
}

func (m *responseMapper) OrderItems(items []entities.OrderItem) ([]responseModel.OrderItem, error) {
	result := make([]responseModel.OrderItem, 0)
	for _, item := range items {
		mappedItem, err := m.OrderItem(&item)
		if err != nil {
			return nil, err
		}
		result = append(result, *mappedItem)
	}
	return result, nil
}
