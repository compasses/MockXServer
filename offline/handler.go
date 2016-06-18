package offline

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/boltdb/bolt"
	. "github.com/compasses/MockXServer/utils"

	"github.com/julienschmidt/httprouter"
)

//MockServerError for some error handler
func MockServerError(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(500)
	result := map[string]interface{}{
		"odata.error": map[string]interface{}{
			"error-code":  "P129S00003",
			"path":        nil,
			"targetLabel": nil,
			"message": map[string]interface{}{
				"lang":  "zh-CN",
				"value": "系统出错。我们已经收到错误通知，正在处理。",
			},
		},
	}
	log.Printf("MockServerError Rsp: %+v", result)

	if err := json.NewEncoder(w).Encode(result); err != nil {
		panic(err)
	}
}

//NotFoundHandler not found error
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	r.ParseForm() //解析参数，默认是不会解析的

	log.Println(net.ParseIP(strings.Split(r.RemoteAddr, ":")[0]))
	dec := json.NewDecoder(r.Body)
	var result interface{}
	dec.Decode(&result)
	log.Println("Req:", result)
	w.Write([]byte("<h1>Welcome to compasses/MockXServer</h1>"))
}

//BackupDB backupdb
func BackupDB(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	err := GlobalDB.View(func(tx *bolt.Tx) error {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", `attachment; filename="EshopOfflineServerDB"`)
		w.Header().Set("Content-Length", strconv.Itoa(int(tx.Size())))
		_, err := tx.WriteTo(w)
		return err
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//PlaceOrder go order
func PlaceOrder(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()
	dec := json.NewDecoder(r.Body)
	var result OrderCreate
	dec.Decode(&result)
	log.Printf("got place order request: %+v\n", result)
	newOrder := RepoCreateOrder(result)

	log.Println("PlaceOrder Rsp: ", newOrder)

	if err := json.NewEncoder(w).Encode(newOrder); err != nil {
		panic(err)
	}
}

//GetSalesOrder return sales order
func GetSalesOrder(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()

	dec := json.NewDecoder(r.Body)
	var req interface{}
	dec.Decode(&req)
	Req := req.(map[string]interface{})
	Id := TableId(ToInt64FromString(Req["orderId"].(string)))

	log.Println("GetSalesOrder", req)
	salesOrder := RepoGetSalesOrder(Id, Req["channelAccountId"].(string))

	if err := json.NewEncoder(w).Encode(salesOrder); err != nil {
		panic(err)
	}
}

func GetSalesOrders(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()

	dec := json.NewDecoder(r.Body)
	var req interface{}
	dec.Decode(&req)
	Req := req.(map[string]interface{})

	log.Println("GetSalesOrders", req)
	salesOrders := RepoGetSalesOrders(Req["channelAccountId"].(string))

	if err := json.NewEncoder(w).Encode(salesOrders); err != nil {
		panic(err)
	}
	log.Println("GetSalesOrders Rsp", req)

}

func Checkout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()
	dec := json.NewDecoder(r.Body)
	var result CheckoutShoppingCart
	err := dec.Decode(&result)
	if err != nil {
		HandleError(err)
	}
	log.Println("CheckoutShoppingCart req: ", result)

	resp := RepoCheckoutShoppingCart(result.ShoppingCart)
	log.Println("CheckoutShoppingCart resp: ", resp)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		panic(err)
	}

}

//ATS find product ATS
func ATS(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()
	dec := json.NewDecoder(r.Body)
	var checkInfo ATSReq
	err := dec.Decode(&checkInfo)

	var rsp interface{}
	log.Println("ATS request: ", checkInfo)

	if err != nil {
		HandleError(err)
		//w.WriteHeader(http.StatusBadRequest)
		log.Println("ATS just returen the default")
		rsp = map[string]interface{}{
			"allowBackOrder": true,
		}
	} else {
		log.Printf("ATS Req %+v\n", checkInfo)
		rsp = RepoCreateATSRsp(&checkInfo)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(rsp); err != nil {
		panic(err)
	}
}

func RecommandationProducts(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()
	dec := json.NewDecoder(r.Body)
	var id RecommandInfo
	err := dec.Decode(&id)

	if err != nil {
		HandleError(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Printf("RecommandProducts Req%+v\n", id)

	RecommandIds := RepoCreateRecommandationProducts(id.ProductId)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(RecommandIds); err != nil {
		panic(err)
	}
	log.Printf("RecommandProducts Rsp%+v\n", RecommandIds)
}

func GetCustomer(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()
	var channelAccountId interface{}
	dec := json.NewDecoder(r.Body)

	err := dec.Decode(&channelAccountId)
	id := channelAccountId.(map[string]interface{})

	if err != nil {
		HandleError(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Printf("Get customer req %+v\n", channelAccountId)

	Account := RepoGetCustomer(GetIdFromStr(id["channelAccountId"].(string)))

	if err := json.NewEncoder(w).Encode(Account); err != nil {
		panic(err)
	}
	log.Printf("Return Customer Rsp %+v\n", Account)
}

func UpdateCustomer(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	if err := json.NewEncoder(w).Encode(nil); err != nil {
		panic(err)
	}

	log.Printf("customer exist")
}
func CheckEmailExistence(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	if err := json.NewEncoder(w).Encode(nil); err != nil {
		panic(err)
	}
	log.Printf("customer exist")
}

func CreateCustomer(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()
	var customer CustomerCreate
	dec := json.NewDecoder(r.Body)

	err := dec.Decode(&customer)

	if err != nil {
		HandleError(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("Customer Create Req%+v\n", customer)

	Account := RepoCreateAccount(customer)

	if err := json.NewEncoder(w).Encode(Account); err != nil {
		panic(err)
	}
	log.Printf("Customer Create Rsp%+v\n", Account)
}

func CustomerAddressNew(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()
	dec := json.NewDecoder(r.Body)

	var addInfo CustomerAddress
	err := dec.Decode(&addInfo)

	if err != nil {
		HandleError(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	Rs := RepoCreateAddress(&addInfo)

	if err = json.NewEncoder(w).Encode(Rs); err != nil {
		panic(err)
	}

	log.Printf("Create address info %+v\n", addInfo)
	log.Printf("Result %+v\n", Rs)

}

func CustomerAddressUpdate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()
	addressId := GetIdFromStr(ps.ByName("id"))

	dec := json.NewDecoder(r.Body)

	var addInfo CustomerAddress
	err := dec.Decode(&addInfo)

	if err != nil {
		HandleError(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	Rs := RepoUpdateAddress(addressId, &addInfo)
	log.Printf("Update address info %+v\n", Rs)

	if err = json.NewEncoder(w).Encode(Rs); err != nil {
		panic(err)
	}
}

func GetCustomerAddress(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()

	customerId := GetIdFromStr(r.Form["$filter"][0])
	Rs := RepoGetCustomerAddress(customerId)

	if err := json.NewEncoder(w).Encode(Rs); err != nil {
		panic(err)
	}
}
func MiscCheck(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()
	dec := json.NewDecoder(r.Body)

	var checkParam map[string]interface{}
	err := dec.Decode(&checkParam)

	if err != nil {
		HandleError(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Println("misc check parames ", checkParam)
	Rs := RetrieveByMapLevel(checkParam, []string{"miscParam", "lines"})
	lines := Rs.([]interface{})

	resp := make(map[string][]interface{})
	for _, val := range lines {
		valm := val.(map[string]interface{})
		resp["lineResult"] = append(resp["lineResult"], map[string]interface{}{
			"onChannel":      "true",
			"ats":            10,
			"allowBackOrder": "true",
			"skuId":          valm["skuId"],
			"valid":          "true",
		})
	}

	log.Println("resp ", resp)

	if err = json.NewEncoder(w).Encode(resp); err != nil {
		panic(err)
	}
}
