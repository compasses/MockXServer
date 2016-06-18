package offline

import (
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
		"/sbo/service/EShopService@getATS",
		ATS,
	},
	Route{
		"Recommendation Products",
		"POST",
		"/sbo/service/ProductService@getRecommendationProductIds",
		RecommandationProducts,
	},
	Route{
		"Create Customer",
		"POST",
		"/sbo/service/CustomerService@createEShopCustomer",
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
		"/sbo/service/EShopService@miscCheck",
		MiscCheck,
	},
	Route{
		"Checkout",
		"POST",
		"/sbo/service/EShopService@checkoutShoppingCart",
		Checkout,
	},
	Route{
		"PlaceOrder",
		"POST",
		"/sbo/service/EShopService@placeOrder",
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
