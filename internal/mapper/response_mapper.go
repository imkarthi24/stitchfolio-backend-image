package mapper

import (
	"encoding/json"
	"time"

	"github.com/imkarthi24/sf-backend/internal/entities"
	responseModel "github.com/imkarthi24/sf-backend/internal/model/response"
	"github.com/loop-kar/pixie/util"
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
	Person(e *entities.Person) (*responseModel.Person, error)
	Persons(items []entities.Person) ([]responseModel.Person, error)
	DressType(e *entities.DressType) (*responseModel.DressType, error)
	DressTypes(items []entities.DressType) ([]responseModel.DressType, error)
	Measurement(e *entities.Measurement) (*responseModel.Measurement, error)
	Measurements(items []entities.Measurement) ([]responseModel.Measurement, error)
	Order(e *entities.Order) (*responseModel.Order, error)
	Orders(items []entities.Order) ([]responseModel.Order, error)
	OrderItem(e *entities.OrderItem) (*responseModel.OrderItem, error)
	OrderItems(items []entities.OrderItem) ([]responseModel.OrderItem, error)
	OrderHistory(e *entities.OrderHistory) (*responseModel.OrderHistory, error)
	OrderHistories(items []entities.OrderHistory) ([]responseModel.OrderHistory, error)
	MeasurementHistory(e *entities.MeasurementHistory) (*responseModel.MeasurementHistory, error)
	MeasurementHistories(items []entities.MeasurementHistory) ([]responseModel.MeasurementHistory, error)
	ExpenseTracker(e *entities.ExpenseTracker) (*responseModel.ExpenseTracker, error)
	ExpenseTrackers(items []entities.ExpenseTracker) ([]responseModel.ExpenseTracker, error)
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

	var statusStr *string
	if e.Status != nil {
		str := string(*e.Status)
		statusStr = &str
	}

	var performedBy *responseModel.User
	if e.PerformedBy != nil {
		user, err := m.User(e.PerformedBy)
		if err != nil {
			return nil, err
		}
		performedBy = user
	}

	return &responseModel.EnquiryHistory{
		ID:              e.ID,
		IsActive:        e.IsActive,
		Status:          statusStr,
		EmployeeComment: e.EmployeeComment,
		CustomerComment: e.CustomerComment,
		VisitingDate:    visitingDateStr,
		CallBackDate:    callBackDateStr,
		EnquiryDate:     util.DateTimeToStringOrDefault(e.EnquiryDate, time.DateOnly),
		ResponseStatus:  string(e.ResponseStatus),
		EnquiryId:       e.EnquiryId,
		EmployeeId:      e.EmployeeId,
		Employee:        employee,
		PerformedAt:     e.PerformedAt,
		PerformedById:   e.PerformedById,
		PerformedBy:     performedBy,
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

	persons, err := m.Persons(e.Persons)
	if err != nil {
		return nil, err
	}

	enquiries, err := m.Enquiries(e.Enquiries)
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
		Persons:        persons,
		Enquiries:      enquiries,
		Orders:         orders,
	}, nil
}

func (m *responseMapper) Person(e *entities.Person) (*responseModel.Person, error) {
	if e == nil {
		return nil, nil
	}

	customer, err := m.Customer(e.Customer)
	if err != nil {
		return nil, err
	}

	measurements, err := m.Measurements(e.Measurements)
	if err != nil {
		return nil, err
	}

	return &responseModel.Person{
		ID:           e.ID,
		IsActive:     e.IsActive,
		FirstName:    e.FirstName,
		LastName:     e.LastName,
		Gender:       string(e.Gender),
		Age:          e.Age,
		CustomerId:   &e.CustomerId,
		Customer:     customer,
		Measurements: measurements,
	}, nil
}

func (m *responseMapper) Persons(items []entities.Person) ([]responseModel.Person, error) {
	result := make([]responseModel.Person, 0)
	for _, item := range items {
		mappedItem, err := m.Person(&item)
		if err != nil {
			return nil, err
		}
		result = append(result, *mappedItem)
	}
	return result, nil
}

func (m *responseMapper) DressType(e *entities.DressType) (*responseModel.DressType, error) {
	if e == nil {
		return nil, nil
	}

	return &responseModel.DressType{
		ID:           e.ID,
		IsActive:     e.IsActive,
		Name:         e.Name,
		Description:  e.Description,
		Measurements: e.Measurements,
	}, nil
}

func (m *responseMapper) DressTypes(items []entities.DressType) ([]responseModel.DressType, error) {
	result := make([]responseModel.DressType, 0)
	for _, item := range items {
		mappedItem, err := m.DressType(&item)
		if err != nil {
			return nil, err
		}
		result = append(result, *mappedItem)
	}
	return result, nil
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

	person, err := m.Person(e.Person)
	if err != nil {
		return nil, err
	}

	dressType, err := m.DressType(e.DressType)
	if err != nil {
		return nil, err
	}

	var takenBy string
	if e.TakenBy != nil {
		takenBy = e.TakenBy.FirstName + " " + e.TakenBy.LastName
	}

	var personName string
	if person != nil {
		personName = person.FirstName + " " + person.LastName
	}

	return &responseModel.Measurement{
		ID:          e.ID,
		IsActive:    e.IsActive,
		Values:      json.RawMessage(e.Value),
		PersonId:    &e.PersonId,
		Person:      person,
		PersonName:  personName,
		DressTypeId: &e.DressTypeId,
		DressType:   dressType,
		TakenById:   e.TakenById,
		TakenBy:     takenBy,
		UpdatedAt:   e.UpdatedAt,
		UpdatedById: e.UpdatedById,
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

	// customer, err := m.Customer(e.Customer)
	// if err != nil {
	// 	return nil, err
	// }

	orderItems, err := m.OrderItems(e.OrderItems)
	if err != nil {
		return nil, err
	}

	orderQuantity := e.OrderQuantity
	orderValue := e.OrderValue
	if orderQuantity == 0 && orderValue == 0 && len(e.OrderItems) > 0 {
		for _, item := range e.OrderItems {
			orderQuantity += item.Quantity
			orderValue += item.Total
		}
	}

	return &responseModel.Order{
		ID:                   e.ID,
		IsActive:             e.IsActive,
		Status:               string(e.Status),
		Notes:                e.Notes,
		AdditionalCharges:    e.AdditionalCharges,
		ExpectedDeliveryDate: e.ExpectedDeliveryDate,
		DeliveredDate:        e.DeliveredDate,
		CustomerId:           e.CustomerId,
		// Customer:       customer,
		CustomerName:   e.Customer.FirstName + " " + e.Customer.LastName,
		OrderTakenById: e.OrderTakenById,
		OrderTakenBy:   e.OrderTakenBy.FirstName + " " + e.OrderTakenBy.LastName,
		OrderQuantity:  orderQuantity,
		OrderValue:     orderValue,
		CreatedAt:      e.CreatedAt,
		UpdatedAt:      e.UpdatedAt,
		OrderItems:     orderItems,
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

	person, err := m.Person(e.Person)
	if err != nil {
		return nil, err
	}

	measurement, err := m.Measurement(e.Measurement)
	if err != nil {
		return nil, err
	}

	return &responseModel.OrderItem{
		ID:                   e.ID,
		IsActive:             e.IsActive,
		Description:          e.Description,
		Quantity:             e.Quantity,
		Price:                e.Price,
		Total:                e.Total,
		AdditionalCharges:    e.AdditionalCharges,
		ExpectedDeliveryDate: e.ExpectedDeliveryDate,
		DeliveredDate:        e.DeliveredDate,
		PersonId:             e.PersonId,
		Person:               person,
		MeasurementId:        e.MeasurementId,
		Measurement:          measurement,

		OrderId: e.OrderId,
		Order:   order,
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

func (m *responseMapper) OrderHistory(e *entities.OrderHistory) (*responseModel.OrderHistory, error) {
	if e == nil {
		return nil, nil
	}

	var status *string
	if e.Status != nil {
		s := string(*e.Status)
		status = &s
	}

	var orderItemData string
	if e.OrderItemData != nil {
		orderItemData = string(*e.OrderItemData)
	}

	order, err := m.Order(e.Order)
	if err != nil {
		return nil, err
	}

	performedBy, err := m.User(e.PerformedBy)
	if err != nil {
		return nil, err
	}

	return &responseModel.OrderHistory{
		ID:                   e.ID,
		IsActive:             e.IsActive,
		Action:               string(e.Action),
		ChangedFields:        e.ChangedFields,
		Status:               status,
		ExpectedDeliveryDate: e.ExpectedDeliveryDate,
		DeliveredDate:        e.DeliveredDate,
		OrderItemId:          e.OrderItemId,
		OrderItemData:        orderItemData,
		OrderId:              e.OrderId,
		Order:                order,
		PerformedAt:          e.PerformedAt,
		PerformedById:        e.PerformedById,
		PerformedBy:          performedBy,
	}, nil
}

func (m *responseMapper) OrderHistories(items []entities.OrderHistory) ([]responseModel.OrderHistory, error) {
	result := make([]responseModel.OrderHistory, 0)
	for _, item := range items {
		mappedItem, err := m.OrderHistory(&item)
		if err != nil {
			return nil, err
		}
		result = append(result, *mappedItem)
	}
	return result, nil
}

func (m *responseMapper) MeasurementHistory(e *entities.MeasurementHistory) (*responseModel.MeasurementHistory, error) {
	if e == nil {
		return nil, nil
	}

	measurement, err := m.Measurement(e.Measurement)
	if err != nil {
		return nil, err
	}

	performedBy, err := m.User(e.PerformedBy)
	if err != nil {
		return nil, err
	}

	var oldValues json.RawMessage
	if len(e.OldValues) > 0 {
		oldValues = json.RawMessage(e.OldValues)
	}

	return &responseModel.MeasurementHistory{
		ID:            e.ID,
		IsActive:      e.IsActive,
		Action:        string(e.Action),
		OldValues:     oldValues,
		MeasurementId: e.MeasurementId,
		Measurement:   measurement,
		PerformedAt:   e.PerformedAt,
		PerformedById: e.PerformedById,
		PerformedBy:   performedBy,
	}, nil
}

func (m *responseMapper) MeasurementHistories(items []entities.MeasurementHistory) ([]responseModel.MeasurementHistory, error) {
	result := make([]responseModel.MeasurementHistory, 0)
	for _, item := range items {
		mappedItem, err := m.MeasurementHistory(&item)
		if err != nil {
			return nil, err
		}
		result = append(result, *mappedItem)
	}
	return result, nil
}

func (m *responseMapper) ExpenseTracker(e *entities.ExpenseTracker) (*responseModel.ExpenseTracker, error) {
	if e == nil {
		return nil, nil
	}

	return &responseModel.ExpenseTracker{
		ID:           e.ID,
		IsActive:     e.IsActive,
		PurchaseDate: e.PurchaseDate,
		BillNumber:   e.BillNumber,
		CompanyName:  e.CompanyName,
		Material:     e.Material,
		Price:        e.Price,
		Location:     e.Location,
		Notes:        e.Notes,
		CreatedAt:    e.CreatedAt,
		UpdatedAt:    e.UpdatedAt,
		CreatedById:  e.CreatedById,
		UpdatedById:  e.UpdatedById,
	}, nil
}

func (m *responseMapper) ExpenseTrackers(items []entities.ExpenseTracker) ([]responseModel.ExpenseTracker, error) {
	result := make([]responseModel.ExpenseTracker, 0)
	for _, item := range items {
		mappedItem, err := m.ExpenseTracker(&item)
		if err != nil {
			return nil, err
		}
		result = append(result, *mappedItem)
	}
	return result, nil
}
