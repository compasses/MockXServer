package offline

import . "github.com/compasses/MockXServer/utils"

type ATSReq struct {
	SkuIds   IDSeqs
	ChanelId TableId
}

type ATSRsp struct {
	SkuId          TableId `json:"skuId"`
	Ats            int64   `json:"ats"`
	AllowBackOrder bool    `json:"allowBackOrder"`
}

type RecommandInfo struct {
	ChannelId  int
	ProductId  TableId `json:",string"`
	CurrencyId int
}

type Customer struct {
	Id           TableId `json:"id,string"`
	CustomerType string  `json:"customerType"`
	Email        string  `json:"email"`
}

type CustomerCreateRsp struct {
	CustomerCode     string  `json:"customerCode"`
	CustomerID       TableId `json:"customerID"`
	ChannelAccountID TableId `json:"channelAccountID"`
	//FailType		 string `json:"failType"`
}

type CustomerAddress struct {
	Id            TableId     `json:"id"`
	CustomerInfo  Customer    `json:"customer"`
	AddressInfo   interface{} `json:"address"`
	DefaultBillTo bool        `json:"defaultBillTo"`
	DefaultShipTo bool        `json:"defaultShipTo"`
}

type CustomerCreate struct {
	ChannelId    TableId
	Account      string
	Customer     Customer
	CustomerType string
	AccountInfo  CustomerCreateRsp
	Addresses    []CustomerAddress
}

type CartItems struct {
	SkuId              interface{} `json:"skuId"`
	UnitPrice          interface{} `json:"unitPrice"`
	Quantity           interface{} `json:"quantity"`
	TaxAmount          interface{} `json:"taxAmount"`
	DiscountPercentage interface{} `json:"discountPercentage"`
	LineTotal          interface{} `json:"lineTotal"`
	LineTotalAfterDisc interface{} `json:"lineTotalAfterDisc"`
	StandardPrice      interface{} `json:"standardPrice"`
	Remark             interface{} `json:"remark"`
}

type ShoppingCart struct {
	CartTotal          interface{} `json:"cartTotal"`
	DiscountPercentage interface{} `json:"discountPercentage"`
	DiscountSum        interface{} `json:"discountSum"`
	PriceMethod        interface{} `json:"priceMethod"`
	CartItems          []CartItems `json:"cartItems"`
}

type CheckoutCartPlayLoad struct {
	ShippingAddress    interface{}  `json:"shippingAddress"`
	BillingAddress     interface{}  `json:"billingAddress"`
	CustomerId         interface{}  `json:"customerId"`
	ChannelAccountId   interface{}  `json:"channelAccountId"`
	ChannelId          interface{}  `json:"channelId"`
	ShoppingCart       ShoppingCart `json:"shoppingCart"`
	ShippingMethod     interface{}  `json:"shippingMethod"`
	Promotion          interface{}  `json:"promotion"`
	TaxTotal           interface{}  `json:"taxTotal"`
	OrderTotal         interface{}  `json:"orderTotal"`
	DiscountPercentage interface{}  `json:"discountPercentage"`
	DiscountSum        interface{}  `json:"discountSum"`
}

type CheckoutShoppingCart struct {
	ShoppingCart CheckoutCartPlayLoad `json:"shoppingCart"`
}

type CheckoutShoppingCartRsp struct {
	CheckoutCartPlayLoad
	ShippingCosts         interface{} `json:"shippingCosts"`
	EnableExpressDelivery bool        `json:"enableExpressDelivery"`
}

type OrderCreate struct {
	EShopOrder CheckoutCartPlayLoad `json:"eShopOrder"`
}
