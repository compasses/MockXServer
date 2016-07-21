package offline

import (
	"github.com/boltdb/bolt"
	"github.com/julienschmidt/httprouter"
)

type Route struct {
	Name       string
	Method     string
	Pattern    string
	HandleFunc httprouter.Handle
}

type Routes []Route

var routes = Routes{
	Route{
		"ATS Check",
		"POST",
		"/api/EshopAdapter/Product/v1/getATS",
		ATS,
	},
	Route{
		"Recommendation Products",
		"POST",
		"/api/EshopAdapter/Product/v1/getRecommendationProductIds",
		RecommandationProducts,
	},
	Route{
		"Create Customer",
		"POST",
		"/api/EshopAdapter/Customer/v1/",
		CreateCustomer,
	},
	Route{
		"check email exist",
		"POST",
		"/sbo/service/CustomerService@isEmailAccountExist",
		CheckEmailExistence,
	},
	Route{
		"Get Customer",
		"POST",
		"/sbo/service/CustomerService@getCustomer",
		GetCustomer,
	},
	Route{
		"Address New",
		"POST",
		"/sbo/service/CustomerAddressNew",
		CustomerAddressNew,
	},
	Route{
		"Update Customer",
		"POST",
		"/sbo/service/	EShopService@updateCustomer",
		UpdateCustomer,
	},
	Route{
		"Address Update",
		"PUT",
		"/sbo/service/:id/",
		CustomerAddressUpdate,
	},
	Route{
		"Address Retrieve",
		"GET",
		"/sbo/service/CustomerAddressNew/",
		GetCustomerAddress,
	},
	Route{
		"MiscCheck",
		"POST",
		"/api/EshopAdapter/Product/v1/miscCheck",
		MiscCheck,
	},
	Route{
		"Checkout",
		"POST",
		"/api/EshopAdapter/Order/v1/checkoutShoppingCart",
		Checkout,
	},
	Route{
		"PlaceOrder",
		"POST",
		"/api/EshopAdapter/Order/v1/placeOrder",
		PlaceOrder,
	},
	Route{
		"GetSalesOrder",
		"POST",
		"/sbo/service/EShopService@getSalesOrder",
		GetSalesOrder,
	},
	Route{
		"GetSalesOrders",
		"POST",
		"/sbo/service/EShopService@getSalesOrders",
		GetSalesOrders,
	},
	Route{
		"BackUpDatabase",
		"GET",
		"/backupDB",
		BackupDB,
	},
	// Route{
	// 	"GeneratePACT",
	// 	"GET",
	// 	"/pact",
	// 	GeneratePACT,
	// },
	// Route{
	// 	"GenerateJSON",
	// 	"GET",
	// 	"/json",
	// 	GenerateJSON,
	// },
}

func NewServerRouter(db *bolt.DB) *httprouter.Router {
	router := httprouter.New()

	for _, route := range routes {
		httpHandle := Logger(route.HandleFunc, route.Name)

		router.Handle(
			route.Method,
			route.Pattern,
			httpHandle,
		)
	}

	router.NotFound = LoggerNotFound(NotFoundHandler)
	GlobalDB = db

	return router
}
