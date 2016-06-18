package offline

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/boltdb/bolt"
	. "github.com/compasses/MockXServer/utils"
)

var GlobalDB *bolt.DB

//tables define
const (
	ProductTable                 = "PRODUCTS"
	SKUTable                     = "PRODUCT-SKU"
	DefaultATS           TableId = 10
	CustomerTable                = "CUSTOMER"
	AddressTable                 = "ADDRESS"
	OrderTable                   = "ORDERS"
	RecommandProductsNum         = 15
)

func RepoCheckoutShoppingCart(checkoutCart CheckoutCartPlayLoad) CheckoutShoppingCartRsp {
	result := CheckoutShoppingCartRsp{checkoutCart, nil, true}

	var cartTotal float64
	for i, val := range result.ShoppingCart.CartItems {
		var quantity float64
		switch v := val.Quantity.(type) {
		case string:
			quantity = ToFloat64FromString(v)
		default:
			quantity = float64(v.(float64))
		}

		price := ToFloat64FromString(val.UnitPrice.(string))
		result.ShoppingCart.CartItems[i].LineTotal = (float64)(quantity) * price
		cartTotal += result.ShoppingCart.CartItems[i].LineTotal.(float64)
	}

	result.ShoppingCart.CartTotal = cartTotal
	result.OrderTotal = cartTotal
	result.TaxTotal = cartTotal

	result.ShippingCosts = []map[string]interface{}{
		{
			"default":         true,
			"identifier":      "FIX-100-100",
			"type":            "FIX",
			"id":              -100,
			"carrierId":       -100,
			"name":            "EShopDefaultRate",
			"cost":            "0.0",
			"minDeliveryDays": 2,
			"maxDeliveryDays": 5,
		},
	}
	return result
}

func RepoGetSalesOrders(channelAccountId string) (result interface{}) {

	GlobalDB.Update(func(tx *bolt.Tx) error {
		orderBucket, err := tx.CreateBucketIfNotExists([]byte(OrderTable))
		if err != nil {
			HandleError(err)
			return err
		}

		cusBuk := orderBucket.Bucket([]byte(channelAccountId))
		if cusBuk == nil {
			log.Println("not found this user " + channelAccountId)
			return nil
		}
		c := cusBuk.Cursor()
		var number int
		var orders []interface{}

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var temp interface{}
			json.Unmarshal(v, &temp)
			orders = append(orders, temp)
			number++
		}
		//reverse the orders
		for i := len(orders)/2 - 1; i >= 0; i-- {
			opp := len(orders) - 1 - i
			orders[i], orders[opp] = orders[opp], orders[i]
		}

		result = map[string]interface{}{
			"value":       orders,
			"odata.count": number,
		}
		return nil
	})
	return
}

func RepoGetSalesOrder(orderId TableId, channelAccountId string) (result interface{}) {

	GlobalDB.Update(func(tx *bolt.Tx) error {
		orderBucket, err := tx.CreateBucketIfNotExists([]byte(OrderTable))
		if err != nil {
			HandleError(err)
			return err
		}

		cusBuk := orderBucket.Bucket([]byte(channelAccountId))
		if cusBuk == nil {
			result = "not found this user " + channelAccountId
			return nil
		}

		orderBytes := cusBuk.Get(orderId.ToBytes())

		json.Unmarshal(orderBytes, &result)

		return nil
	})
	return
}

func GenerateOrder(order OrderCreate, orderId uint64) interface{} {
	return map[string]interface{}{
		"id":            orderId,
		"docNumber":     orderId,
		"billingAddr":   order.EShopOrder.BillingAddress,
		"shippingAddr":  order.EShopOrder.ShippingAddress,
		"shippingCost":  "0",
		"subTotal":      "500",
		"grossDocTotal": "500",
		"taxTotal":      "0",
		"customer": map[string]interface{}{
			"id": order.EShopOrder.CustomerId,
		},
		"process": map[string]interface{}{
			"id":          3,
			"processName": "快递",
		},
		"salesOrderLines": []map[string]interface{}{
			{
				"amazonOrderItemCode":       nil,
				"baseDocId":                 nil,
				"baseDocLineId":             nil,
				"baseDocLineNumber":         nil,
				"baseDocNumber":             nil,
				"baseDocType":               nil,
				"canelReason":               nil,
				"costTotal":                 "0",
				"discountPercentage":        "0",
				"docCurrencyId":             1,
				"exceptionFlag":             false,
				"exceptionReasonDesc":       nil,
				"giftMessage":               nil,
				"grossLineTotal":            "500.00",
				"grossLineTotalAfterDisc":   "500.00",
				"grossLineTotalAfterDiscLC": "500.00",
				"grossLineTotalLC":          "500.00",
				"grossProfitAmount":         "0",
				"grossProfitRate":           "0",
				"grossUnitPrice":            "500.00",
				"id":                        119,
				"inventoryUomName":          "Unit",
				"inventoryUomQuantity":      "1.00",
				"invoiceStatus":             "tNotInvoiced",
				"isNonLogistical":           false,
				"isPreparingStock":          false,
				"isPromotionAppliable":      true,
				"isService":                 false,
				"lineAction":                "C,A,I,S",
				"lineCalcBase":              "byTotal",
				"lineComments":              nil,
				"lineNumber":                1,
				"lineType":                  "tProductLine",
				"logisticsStatus":           "tAllocated",
				"merchantFulfillmentItemID": nil,
				"netLineTotal":              "500.00",
				"netLineTotalAfterDisc":     "500.00",
				"netLineTotalAfterDiscLC":   "500.00",
				"netLineTotalLC":            "500.00",
				"netUnitPrice":              "500.00",
				"originLine":                nil,
				"originSkuPrice":            "500",
				"planningWhsId":             nil,
				"priceSource":               nil,
				"promotionDescription":      nil,
				"promotionHintBenefit":      nil,
				"promotionHintDes":          nil,
				"propertyDynamicMeta":       "grossLineTotalAfterDiscLC:T,grossLineTotalLC:T,salesUomName:T,grossLineTotalAfterDisc:T,promotionItem:T,grossProfitRate:T,netLineTotalLC:T,taxAmountLC:T,grossProfitAmount:T,salesUom:T,lineType:T,exceptionReasonDesc:T,inventoryUomQuantity:T,exceptionReason:T,originSkuPrice:T,netLineTotalAfterDisc:T,netLineTotalAfterDiscLC:T,inventoryUom:T,totalLC:T,exceptionFlag:T,totalAfterDiscountLC:T,costTotal:T,invoiceStatus:T,lineNumber:T,inventoryUomName:T,taxAmount:T,promotion:T",
				"purchaseOrderId":           nil,
				"purchaseOrderLineId":       nil,
				"quantity":                  "1",
				"remark":                    nil,
				"salesUomName":              "Unit",
				"shippingId":                nil,
				"shippingType":              nil,
				"skuCode":                   "ProductCode13",
				"skuMainLogoURL":            "http://internal-ci.s3.amazonaws.com/T1/2015/11/09/b09fb2bc-10f8-48fd-a4e1-e190194992a4",
				"skuName":                   "Red Silk Dress",
				"targetDocId":               nil,
				"taxAmount":                 "0",
				"taxAmountLC":               "0",
				"total":                     "500.00",
				"totalAfterDiscount":        "500.00",
				"totalAfterDiscountLC":      "500.00",
				"totalLC":                   "500.00",
				"unitCost":                  "0",
				"unitPrice":                 "500.00",
				"uomConversionRate":         "1",
				"variantValues":             "",
				"warehouseCode":             "主仓库",
				"warehouseName":             "主仓库",
				"docCurrency":               nil,
				"exceptionReason":           nil,
				"inventoryUom": map[string]interface{}{
					"description": "无",
					"id":          1,
					"name":        "无",
				},
				"promotion":     nil,
				"promotionItem": nil,
				"salesUom": map[string]interface{}{
					"description": "无",
					"id":          1,
					"name":        "无",
				},
				"sku": map[string]interface{}{
					"code":    "ProductCode13",
					"id":      57,
					"name":    "Red Silk Dress",
					"product": nil,
				},
				"warehouse": map[string]interface{}{
					"id":                 1,
					"whsName":            "主仓库",
					"virtualNodes":       nil,
					"warehouseOperators": nil,
					"defaulterUser":      nil,
				},
			},
		}}
}

func RepoCreateOrder(order OrderCreate) interface{} {
	var newOrder interface{}

	GlobalDB.Update(func(tx *bolt.Tx) error {
		orderBucket, err := tx.CreateBucketIfNotExists([]byte(OrderTable))
		if err != nil {
			HandleError(err)
			return err
		}
		var cusOrderBuck *bolt.Bucket
		if order.EShopOrder.ChannelAccountId == nil {
			cusOrderBuck, err = orderBucket.CreateBucketIfNotExists([]byte("GUESTUSER"))
		} else {
			cusOrderBuck, err = orderBucket.CreateBucketIfNotExists([]byte(order.EShopOrder.ChannelAccountId.(string)))
		}
		if err != nil {
			HandleError(err)
			return err
		}
		newId, _ := orderBucket.NextSequence()
		newOrder = GenerateOrder(order, newId)

		orderBytes, _ := json.Marshal(newOrder)
		cusOrderBuck.Put(TableId(newId).ToBytes(), orderBytes)

		return nil
	})

	return newOrder
}

func GetProductATS(ProductId TableId) int64 {
	var atsQua int64
	log.Println("Try to process ats ", GlobalDB)

	GlobalDB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(SKUTable))

		if err != nil {
			HandleError(err)
			return err
		}

		bb, errbb := b.CreateBucketIfNotExists(ProductId.ToBytes())

		if errbb != nil {
			HandleError(errbb)
			return err
		}

		ats := bb.Get([]byte("ats"))
		if ats == nil {
			//new product id, need initialize the ats info
			bb.Put([]byte("ats"), DefaultATS.ToBytes())
			atsQua = int64(DefaultATS)
		} else {
			result := ToInt64FromBytes(ats)
			if err != nil {
				HandleError(err)
				return err
			}
			atsQua = result
		}
		return nil
	})
	return atsQua
}

func RepoCreateATSRsp(req *ATSReq) []ATSRsp {
	var rest []ATSRsp

	for _, atsR := range req.SkuIds {

		atsQua := GetProductATS(atsR)
		rsp := ATSRsp{
			SkuId:          atsR,
			Ats:            atsQua,
			AllowBackOrder: true,
		}
		rest = append(rest, rsp)
	}
	log.Printf("ATS Rsp %+v\n", rest)
	return rest
}

func RepoCreateRecommandationProducts(Id TableId) []TableId {
	var ProductId []TableId

	GlobalDB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(ProductTable))
		if err != nil {
			HandleError(err)
			return err
		}

		c := b.Cursor()
		num := 0

		for k, _ := c.First(); k != nil && (num < RecommandProductsNum); k, _ = c.Next() {
			ProductId = append(ProductId, TableId(ToInt64FromBytes(k)))
			num++
		}
		_, errbb := b.CreateBucketIfNotExists(Id.ToBytes())

		if errbb != nil {
			HandleError(errbb)
			return err
		}

		return nil
	})

	return ProductId
}

func RepoGetCustomer(channelAccountId TableId) interface{} {
	var result CustomerCreate
	key := []byte(CustomerTable)

	GlobalDB.Update(func(tx *bolt.Tx) error {
		pb, err := tx.CreateBucketIfNotExists(key)

		if err != nil {
			HandleError(err)
			return err
		}

		account := pb.Get(channelAccountId.ToBytes())
		if account == nil {
			HandleError(err)
			return err
		}

		json.Unmarshal(account, &result)
		return nil
	})

	return map[string]interface{}{
		"id":                    result.AccountInfo.CustomerID,
		"displayName":           "Jet He",
		"email":                 "update@customer.com", //result.Account,
		"checkDuplication":      true,
		"creationTime":          "2015-11-17T05:56:58.812Z",
		"creatorDisplayName":    "ERP SUITE",
		"creditLimit":           "250000",
		"creditBalance":         "100000",
		"outstandingPayment":    "300000",
		"customerCode":          "3",
		"customerName":          "Jet Jet",
		"customerType":          "CORPORATE_CUSTOMER",
		"dateOfBirth":           nil,
		"facebookAccount":       nil,
		"fax":                   "22323232",
		"firstName":             "GG",
		"gender":                nil,
		"googleAccount":         nil,
		"lastMarketingCampaign": nil,
		"lastName":              "HH",
		"linkedINAccount":       nil,
		"mainContact":           nil,
		"marketingStatus":       "Unknown",
		"membershipBalance":     0,
		"membershipId":          nil,
		"membershipSwitchOn":    false,
		"membershipTotalEarn":   0,
		"mobile":                1588888888,
		"ownerCode":             1,
		"ownerDisplayName":      "ERP SUITE",
		"phone":                 "12121212",
		"portraitId":            0,
		"position":              nil,
		"remarks":               "Good job",
		"resetPasswordLink":     nil,
		"socialImageUrl":        nil,
		"stage":                 "SUSPECT",
		"status":                "ACTIVE",
		"targetGroup":           nil,
		"taxType":               "LIABLE",
		"title":                 nil,
		"twitterAccount":        nil,
		"twitterDisplayName":    nil,
		"updateTime":            "2015-11-17T05:56:58.812Z",
		"updatorDisplayName":    "ERP SUITE",
		"vatRegistrationNumber": nil,
		"versionNum":            0,
		"webSite":               nil,
		"weiboDisplayName":      nil,
		"ext_default_UDF5":      "11",
		"priceList":             nil,
		"customerGroup":         nil,
		"defaultBillToAddress":  nil,
		"defaultShipToAddress":  nil,
		"industry":              nil,
		"language":              nil,
		"linkedCustomer": map[string]interface{}{
			"displayName":          "Jet Jet",
			"id":                   result.AccountInfo.CustomerID,
			"priceList":            nil,
			"customerGroup":        nil,
			"defaultBillToAddress": nil,
			"defaultShipToAddress": nil,
			"industry":             nil,
			"language":             nil,
			"linkedCustomer":       nil,
			"membershipLevel":      nil,
			"paymentAccount":       nil,
			"paymentTerm":          nil,
			"serviceLevelPlan":     nil,
			"twitter":              nil,
			"weibo":                nil,
		},
		"membershipLevel": map[string]interface{}{
			"id":   1,
			"name": "普通会员",
		},
		"paymentAccount": map[string]interface{}{
			"acctName": "Cash",
			"id":       1,
			"acctType": nil,
			"country":  nil,
		},
		"paymentTerm": map[string]interface{}{
			"id":   1,
			"name": "现金基础",
		},
		"serviceLevelPlan": map[string]interface{}{
			"displayName": "Gold",
			"id":          1,
		},
		"twitter": nil,
		"weibo":   nil,
	}
}

func RepoCreateAccount(customer CustomerCreate) CustomerCreateRsp {
	var result CustomerCreateRsp
	key := []byte(CustomerTable)

	GlobalDB.Update(func(tx *bolt.Tx) error {
		pb, err := tx.CreateBucketIfNotExists(key)

		if err != nil {
			HandleError(err)
			return err
		}

		customerId, _ := pb.NextSequence()

		b, err := pb.CreateBucketIfNotExists([]byte(customer.Account))

		if err != nil {
			HandleError(err)
			return err
		}

		user := b.Get([]byte("User"))
		if user == nil {
			//create new user
			customer.AccountInfo.CustomerID = TableId(customerId)
			customer.AccountInfo.CustomerCode = "offline" + strconv.FormatInt(customer.AccountInfo.CustomerID.ToInt(), 10)
			customer.AccountInfo.ChannelAccountID = TableId(customer.ChannelId * customer.AccountInfo.CustomerID)
			//test
			cusStream, _ := json.Marshal(&customer)
			b.Put([]byte("User"), cusStream)
			pb.Put(customer.AccountInfo.ChannelAccountID.ToBytes(), cusStream)
			//customer.AccountInfo.FailType = "CUSTOMERTYPEMISSMATCH"
			result = customer.AccountInfo
			//			pbaddr, err := tx.CreateBucketIfNotExists([]byte(AddressTable))
			//			if err != nil {
			//				HandleError(err)
			//				return err
			//			}
			//			pbaddr.Put(customer.AccountInfo.CustomerID.ToBytes(), cusStream)
		} else {
			json.Unmarshal(user, &customer)
			result = customer.AccountInfo
		}

		return nil
	})

	return result
}

func RepoCreateAddress(customer *CustomerAddress) (result interface{}) {
	key := []byte(AddressTable)

	GlobalDB.Update(func(tx *bolt.Tx) error {
		pb, err := tx.CreateBucketIfNotExists(key)

		if err != nil {
			HandleError(err)
			return err
		}

		paddr, err := pb.CreateBucketIfNotExists(customer.CustomerInfo.Id.ToBytes())
		if err != nil {
			HandleError(err)
			return err
		}

		//create the bucket for store address info
		addressBucket, _ := paddr.CreateBucketIfNotExists([]byte("addresses"))
		addressId, _ := addressBucket.NextSequence()
		customer.Id = TableId(addressId)

		streamD, _ := json.Marshal(customer)
		addressBucket.Put(TableId(addressId).ToBytes(), streamD)
		result = customer

		return nil
	})
	return
}

func RepoUpdateAddress(addressId TableId, customer *CustomerAddress) (result interface{}) {
	key := []byte(AddressTable)

	GlobalDB.Update(func(tx *bolt.Tx) error {
		pb, err := tx.CreateBucketIfNotExists(key)

		if err != nil {
			HandleError(err)
			return err
		}

		cusBuk := pb.Bucket(customer.CustomerInfo.Id.ToBytes())
		if cusBuk == nil {
			result = "not found this user:" + customer.CustomerInfo.Id.ToString()
			return nil
		} else {
			addressBucket := cusBuk.Bucket([]byte("addresses"))
			if addressBucket == nil {
				result = "not found this account:" + customer.CustomerInfo.Email
				return nil
			}

			//create the bucket for store address info
			oldAddress := addressBucket.Get(addressId.ToBytes())
			var oldCustomerAddr CustomerAddress
			json.Unmarshal(oldAddress, &oldCustomerAddr)

			//set to new one
			oldCustomerAddr.AddressInfo = customer.AddressInfo
			oldCustomerAddr.DefaultShipTo = customer.DefaultShipTo
			oldCustomerAddr.DefaultBillTo = customer.DefaultBillTo
			streamD, _ := json.Marshal(oldCustomerAddr)

			addressBucket.Put(TableId(addressId).ToBytes(), streamD)
			result = oldCustomerAddr
		}
		return nil
	})
	return
}

func RepoGetCustomerAddress(customerId TableId) (result interface{}) {
	key := []byte(AddressTable)

	GlobalDB.Update(func(tx *bolt.Tx) error {
		pb, err := tx.CreateBucketIfNotExists(key)

		if err != nil {
			HandleError(err)
			return err
		}
		cusBuk := pb.Bucket(customerId.ToBytes())
		if cusBuk == nil {
			result = "not found this user:" + customerId.ToString()
			return nil
		} else {

			countinfo := make(map[string]interface{})
			var count int64 = 0

			//create the bucket for store address info
			addressBucket := cusBuk.Bucket([]byte("addresses"))
			if addressBucket == nil {
				countinfo["odata.count"] = 0
				result = countinfo
			} else {
				var bos []interface{}
				cur := addressBucket.Cursor()
				for k, v := cur.First(); k != nil; k, v = cur.Next() {
					count++
					var Addr CustomerAddress
					json.Unmarshal(v, &Addr)
					bos = append(bos, Addr)
				}

				countinfo["odata.count"] = 150
				countinfo["value"] = bos
				result = countinfo
			}
		}
		return nil
	})
	return
}
