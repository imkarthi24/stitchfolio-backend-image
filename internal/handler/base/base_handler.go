package base

import "github.com/imkarthi24/sf-backend/internal/handler"

type BaseHandler struct {
	HealthHandler             Health
	UserHandler               *handler.UserHandler
	ChannelHandler            *handler.ChannelHandler
	MasterConfigHandler       *handler.MasterConfigHandler
	AdminHandler              *handler.AdminHandler
	CustomerHandler           *handler.CustomerHandler
	EnquiryHandler            *handler.EnquiryHandler
	OrderHandler              *handler.OrderHandler
	OrderItemHandler          *handler.OrderItemHandler
	MeasurementHandler        *handler.MeasurementHandler
	PersonHandler             *handler.PersonHandler
	DressTypeHandler          *handler.DressTypeHandler
	OrderHistoryHandler       *handler.OrderHistoryHandler
	MeasurementHistoryHandler *handler.MeasurementHistoryHandler
	EnquiryHistoryHandler     *handler.EnquiryHistoryHandler
}

func ProvideBaseHandler(health Health,
	user *handler.UserHandler,
	channelHandler *handler.ChannelHandler,
	masterConfigHandler *handler.MasterConfigHandler,
	adminHandler *handler.AdminHandler,
	customerHandler *handler.CustomerHandler,
	enquiryHandler *handler.EnquiryHandler,
	orderHandler *handler.OrderHandler,
	orderItemHandler *handler.OrderItemHandler,
	measurementHandler *handler.MeasurementHandler,
	personHandler *handler.PersonHandler,
	dressTypeHandler *handler.DressTypeHandler,
	orderHistoryHandler *handler.OrderHistoryHandler,
	measurementHistoryHandler *handler.MeasurementHistoryHandler,
	enquiryHistoryHandler *handler.EnquiryHistoryHandler,
) BaseHandler {
	return BaseHandler{
		HealthHandler:             health,
		UserHandler:               user,
		ChannelHandler:            channelHandler,
		MasterConfigHandler:       masterConfigHandler,
		AdminHandler:              adminHandler,
		CustomerHandler:           customerHandler,
		EnquiryHandler:            enquiryHandler,
		OrderHandler:              orderHandler,
		OrderItemHandler:          orderItemHandler,
		MeasurementHandler:        measurementHandler,
		PersonHandler:             personHandler,
		DressTypeHandler:          dressTypeHandler,
		OrderHistoryHandler:       orderHistoryHandler,
		MeasurementHistoryHandler: measurementHistoryHandler,
		EnquiryHistoryHandler:     enquiryHistoryHandler,
	}
}
