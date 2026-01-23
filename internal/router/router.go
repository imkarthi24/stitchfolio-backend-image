package router

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/imkarthi24/sf-backend/internal/config"
	"github.com/imkarthi24/sf-backend/internal/constants"
	baseHandler "github.com/imkarthi24/sf-backend/internal/handler/base"
	"github.com/imkarthi24/sf-backend/internal/log/newreliclog"
	router "github.com/imkarthi24/sf-backend/internal/router/middleware"
	"github.com/imkarthi24/sf-backend/pkg/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	docs "github.com/imkarthi24/sf-backend/docs"
)

func InitRouter(handler baseHandler.BaseHandler, srvConfig config.ServerConfig) *gin.Engine {

	g := gin.Default()
	g.Use(gin.Recovery())

	// Middlewares
	g.Use(middleware.NewRelicMiddleWare(newreliclog.Get()))
	g.Use(middleware.LogMiddleware())

	g.Use(middleware.Security())
	g.Use(middleware.CORS())
	g.Use(middleware.RequestParser())
	g.Use(gzip.Gzip(gzip.DefaultCompression))

	docs.SwaggerInfo.Host = srvConfig.Host
	appRouter := g.Group(constants.API_PREFIX_V1)
	{
		appRouter.GET(constants.HEALTH, handler.HealthHandler.Health)
		appRouter.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		//**************NON JWT ENDPOINTS**************************//

		nonJwtEndpoints := appRouter.Group("user")
		{
			nonJwtEndpoints.POST("login", handler.UserHandler.Login)
			nonJwtEndpoints.POST("forgot-password", handler.UserHandler.ForgotPassword)
			nonJwtEndpoints.POST("reset-password", handler.UserHandler.ResetPassword)
		}
		nonJwtEndpoints = appRouter.Group("channel")
		{
			nonJwtEndpoints.GET("", handler.ChannelHandler.GetAllChannels)
		}

		// External Endpoints
		// Used with channel-id and channel-name header
		externalEndpoints := appRouter.Group("external", router.GenerateExternalSession())
		{
			configEndpoint := externalEndpoints.Group("masterConfig")
			{
				configEndpoint.GET("value", handler.MasterConfigHandler.GetValue)
			}
		}

		//**************JWT ENDPOINTS**************************//

		userEndpoints := appRouter.Group("user", router.VerifyJWT(srvConfig.JwtSecretKey))
		{
			userEndpoints.POST("", handler.UserHandler.SaveUser)
			userEndpoints.PUT(":id", handler.UserHandler.UpdateUser)

			userEndpoints.GET(":id", handler.UserHandler.Get)
			userEndpoints.GET("autocomplete", handler.UserHandler.GetUsersForAutoComplete)
			userEndpoints.GET("", handler.UserHandler.GetAllUsers)
			userEndpoints.GET("refresh-token", handler.UserHandler.RefreshToken)
			userEndpoints.GET("switch-channel/:id", handler.UserHandler.SwitchChannel)

			userEndpoints.DELETE(":id", handler.UserHandler.Delete)

			configEndpoints := userEndpoints.Group("config", router.VerifyJWT(srvConfig.JwtSecretKey))
			{
				configEndpoints.POST("", handler.UserHandler.SaveUserConfig)

				configEndpoints.PUT(":id", handler.UserHandler.UpdateUserConfig)

				configEndpoints.GET("", handler.UserHandler.GetUserConfig)
			}

			channelEndpoints := userEndpoints.Group("channel", router.VerifyJWT(srvConfig.JwtSecretKey))
			{
				channelEndpoints.POST("", handler.UserHandler.SaveUserChannelDetail)
				channelEndpoints.PUT(":id", handler.UserHandler.UpdateUserChannelDetail)
				channelEndpoints.DELETE(":id", handler.UserHandler.DeleteUserChannelDetail)
				channelEndpoints.GET("accessible", handler.UserHandler.GetUserAccessibleChannels)
			}
		}

		channelEndpoints := appRouter.Group("channel", router.VerifyJWT(srvConfig.JwtSecretKey))
		{
			channelEndpoints.POST("", handler.ChannelHandler.SaveChannel)

			channelEndpoints.PUT(":id", handler.ChannelHandler.UpdateChannel)

			channelEndpoints.GET(":id", handler.ChannelHandler.Get)
			// channelEndpoints.GET("", handler.ChannelHandler.GetAllChannels)
			// this has been made a non JWT endpoint
			channelEndpoints.GET("autocomplete", handler.ChannelHandler.ChannelAutoComplete)

			channelEndpoints.DELETE(":id", handler.ChannelHandler.Delete)
		}

		masterConfigEndpoints := appRouter.Group("masterConfig", router.VerifyJWT(srvConfig.JwtSecretKey))
		{
			masterConfigEndpoints.POST("", handler.MasterConfigHandler.Create)
			masterConfigEndpoints.POST("values", handler.MasterConfigHandler.GetMultipleValues)

			masterConfigEndpoints.PUT(":id", handler.MasterConfigHandler.Update)

			masterConfigEndpoints.GET("/browse", handler.MasterConfigHandler.Browse)
			masterConfigEndpoints.GET(":id", handler.MasterConfigHandler.Get)
			masterConfigEndpoints.GET("value", handler.MasterConfigHandler.GetValue)
		}

		adminEndpoints := appRouter.Group("admin", router.VerifyJWT(srvConfig.JwtSecretKey))
		{
			adminEndpoints.POST("switch-branch", handler.AdminHandler.SwitchBranch)
		}

		customerEndpoints := appRouter.Group("customer", router.VerifyJWT(srvConfig.JwtSecretKey))
		{
			customerEndpoints.POST("", handler.CustomerHandler.SaveCustomer)
			customerEndpoints.PUT(":id", handler.CustomerHandler.UpdateCustomer)
			customerEndpoints.GET("autocomplete", handler.CustomerHandler.AutocompleteCustomer)
			customerEndpoints.GET(":id", handler.CustomerHandler.Get)
			customerEndpoints.GET("", handler.CustomerHandler.GetAllCustomers)
			customerEndpoints.DELETE(":id", handler.CustomerHandler.Delete)
		}

		enquiryEndpoints := appRouter.Group("enquiry", router.VerifyJWT(srvConfig.JwtSecretKey))
		{
			enquiryEndpoints.POST("", handler.EnquiryHandler.SaveEnquiry)
			enquiryEndpoints.PUT(":id/customer", handler.EnquiryHandler.UpdateEnquiryAndCustomer)
			enquiryEndpoints.PUT(":id", handler.EnquiryHandler.UpdateEnquiry)
			enquiryEndpoints.GET(":id", handler.EnquiryHandler.Get)
			enquiryEndpoints.GET("", handler.EnquiryHandler.GetAllEnquiries)
			enquiryEndpoints.DELETE(":id", handler.EnquiryHandler.Delete)
		}

		orderEndpoints := appRouter.Group("order", router.VerifyJWT(srvConfig.JwtSecretKey))
		{
			orderEndpoints.POST("", handler.OrderHandler.SaveOrder)
			orderEndpoints.PUT(":id", handler.OrderHandler.UpdateOrder)
			orderEndpoints.GET(":id", handler.OrderHandler.Get)
			orderEndpoints.GET("", handler.OrderHandler.GetAllOrders)
			orderEndpoints.DELETE(":id", handler.OrderHandler.Delete)
		}

		orderItemEndpoints := appRouter.Group("order-item", router.VerifyJWT(srvConfig.JwtSecretKey))
		{
			orderItemEndpoints.POST("", handler.OrderItemHandler.SaveOrderItem)
			orderItemEndpoints.PUT(":id", handler.OrderItemHandler.UpdateOrderItem)
			orderItemEndpoints.GET(":id", handler.OrderItemHandler.Get)
			orderItemEndpoints.GET("", handler.OrderItemHandler.GetAllOrderItems)
			orderItemEndpoints.DELETE(":id", handler.OrderItemHandler.Delete)
		}

		measurementEndpoints := appRouter.Group("measurement", router.VerifyJWT(srvConfig.JwtSecretKey))
		{
			measurementEndpoints.POST("", handler.MeasurementHandler.SaveMeasurement)
			measurementEndpoints.POST("bulk", handler.MeasurementHandler.SaveBulkMeasurements)
			measurementEndpoints.PUT(":id", handler.MeasurementHandler.UpdateMeasurement)
			measurementEndpoints.GET(":id", handler.MeasurementHandler.Get)
			measurementEndpoints.GET("", handler.MeasurementHandler.GetAllMeasurements)
			measurementEndpoints.DELETE(":id", handler.MeasurementHandler.Delete)
		}

		personEndpoints := appRouter.Group("person", router.VerifyJWT(srvConfig.JwtSecretKey))
		{
			personEndpoints.POST("", handler.PersonHandler.SavePerson)
			personEndpoints.PUT(":id", handler.PersonHandler.UpdatePerson)
			personEndpoints.GET(":id", handler.PersonHandler.Get)
			personEndpoints.GET("", handler.PersonHandler.GetAllPersons)
			personEndpoints.GET("customer/:customerId", handler.PersonHandler.GetByCustomerId)
			personEndpoints.DELETE(":id", handler.PersonHandler.Delete)
		}

		dressTypeEndpoints := appRouter.Group("dress-type", router.VerifyJWT(srvConfig.JwtSecretKey))
		{
			dressTypeEndpoints.POST("", handler.DressTypeHandler.SaveDressType)
			dressTypeEndpoints.PUT(":id", handler.DressTypeHandler.UpdateDressType)
			dressTypeEndpoints.GET(":id", handler.DressTypeHandler.Get)
			dressTypeEndpoints.GET("", handler.DressTypeHandler.GetAllDressTypes)
			dressTypeEndpoints.DELETE(":id", handler.DressTypeHandler.Delete)
		}

		orderHistoryEndpoints := appRouter.Group("order-history", router.VerifyJWT(srvConfig.JwtSecretKey))
		{
			orderHistoryEndpoints.POST("", handler.OrderHistoryHandler.SaveOrderHistory)
			orderHistoryEndpoints.GET(":id", handler.OrderHistoryHandler.Get)
			orderHistoryEndpoints.GET("", handler.OrderHistoryHandler.GetAllOrderHistories)
			orderHistoryEndpoints.GET("order/:orderId", handler.OrderHistoryHandler.GetByOrderId)
		}

		measurementHistoryEndpoints := appRouter.Group("measurement-history", router.VerifyJWT(srvConfig.JwtSecretKey))
		{
			measurementHistoryEndpoints.POST("", handler.MeasurementHistoryHandler.SaveMeasurementHistory)
			measurementHistoryEndpoints.GET(":id", handler.MeasurementHistoryHandler.Get)
			measurementHistoryEndpoints.GET("", handler.MeasurementHistoryHandler.GetAllMeasurementHistories)
			measurementHistoryEndpoints.GET("measurement/:measurementId", handler.MeasurementHistoryHandler.GetByMeasurementId)
		}

		enquiryHistoryEndpoints := appRouter.Group("enquiry-history", router.VerifyJWT(srvConfig.JwtSecretKey))
		{
			enquiryHistoryEndpoints.POST("", handler.EnquiryHistoryHandler.SaveEnquiryHistory)
			enquiryHistoryEndpoints.GET(":id", handler.EnquiryHistoryHandler.Get)
			enquiryHistoryEndpoints.GET("", handler.EnquiryHistoryHandler.GetAllEnquiryHistories)
			enquiryHistoryEndpoints.GET("enquiry/:enquiryId", handler.EnquiryHistoryHandler.GetByEnquiryId)
		}
	}
	return g
}
