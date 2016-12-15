package utils

import (
	"bytes"
	"crypto/tls"
	"database/sql"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/franela/goreq"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func TestBytesToInt(t *testing.T) {
	test1 := int64(256)
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, test1)
	//binary.PutVarint(buf.Bytes(), int64(test1))
	fmt.Println("byte is ", buf.Bytes())
	bufr := bytes.NewBuffer(buf.Bytes())
	var rest int64
	err := binary.Read(bufr, binary.LittleEndian, &rest)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("rest value ", rest)

	num, _ := binary.Varint(buf.Bytes())
	fmt.Println("num is ", num)

	var testS []TableId
	testS = append(testS, 34, 3443, 344242)
	tids := GetSliceBytesFromInts(testS)
	fmt.Println("length is ", len(tids), "val is ", tids)

	tidr := GetSliceIntFromBytes(tids)
	fmt.Println("int value is ", tidr)

	t4 := ToInt64FromBytes(tids[0:8])
	fmt.Println("t4 val is ", t4)

}

func TestBytesToInt2(t *testing.T) {
	var intVar = 123
	fmt.Println("intVar is : ", intVar)

	intByte := []byte(strconv.Itoa(intVar))

	fmt.Println("intByte is : ", intByte)
}

func TestGetIntFromStr(t *testing.T) {
	input := "CreateCustomerNew(23)"
	valId := regexp.MustCompile(`(\d+)`)
	val, _ := strconv.Atoi(valId.FindString(input))
	fmt.Println("id :", val)
	var ni []interface{}
	ni = append(ni, "ok")
	ni = append(ni, 2323)

	fmt.Println("interfaces values ", ni)
}

func TestHttpGet(t *testing.T) {
	res, err := goreq.Request{
		Method:    "GET",
		Uri:       "https://cnpvgvb1ep052.pvgl.sap.corp:29900",
		ShowDebug: true,
	}.Do()

	//	dec := json.NewDecoder(res.Response.Body)
	//	var result interface{}
	//	dec.Decode(&result)

	//	//fmt.Println("Result ", res.Response)
	//	fmt.Println("resp result ", result)
	fmt.Println("error is ", err)
	if err == nil {
		nres1, _ := ioutil.ReadAll(res.Response.Body)
		fmt.Println("raw body is ", string(nres1))
	}

	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{},
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get("https://www.baidu.com")
	fmt.Println("err is ", err)
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	var result2 interface{}
	err = dec.Decode(&result2)
	fmt.Println("decode error ", err)
	fmt.Println("body is ", result2)
	nres, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("raw body is ", string(nres))

}

func TestDBCase(t *testing.T) {
	db, err := sql.Open("mysql", "root:12345@tcp(cnpvgvb1ep140.pvgl.sap.corp:3306)/")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	// Select the DataBase
	db.Exec("use ESHOPDB116")
	// Execute the query
	//    rows, err := db.Query("show databases")
	//    if err != nil {
	//        panic(err.Error()) // proper error handling instead of panic in your app
	//    }
	rows, err := db.Query("select option_value from wp_options where option_name = 'eshopSetting'")
	// Get column names
	//    columns, err := rows.Columns()
	//	fmt.Println("columns: ", columns)

	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	//	// Make a slice for the values
	//	values := make([]sql.RawBytes, 1)

	//	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	//	// references into such a slice
	//	// See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
	//	scanArgs := make([]interface{}, len(values))
	//	for i := range values {
	//		scanArgs[i] = &values[i]
	//	}
	//	fmt.Printf("scanArgs: %v\n", values)
	var res []byte
	// Fetch rows
	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(&res)

		var dataS interface{}
		json.Unmarshal(res, &dataS)
		m := dataS.(map[string]interface{})
		fmt.Println(": ", m["shopName"])
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}

		// Now do something with the data.
		//        // Here we just print each column as a string.
		//        var value string
		//        for _, col := range values {
		//            // Here we can check if the value is nil (NULL value)
		//            if col == nil {
		//                value = "NULL"
		//            } else {
		//                value = string(col)
		//            }

		//        }
		//        fmt.Println("-----------------------------------")
	}
	if err = rows.Err(); err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
}

func load(direction string, enc encoding.Encoding) (func() transform.Transformer, error) {

	newTransformer := enc.NewEncoder
	if direction == "Decode" {
		newTransformer = enc.NewDecoder
	}

	return newTransformer, nil
}

func TestConversion(t *testing.T) {
	f, err := os.OpenFile("testlogfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		t.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println("This is a test log entry")

	sr := strings.NewReader("你好，世界")

	newTransformer, _ := load("Decode", simplifiedchinese.GBK)

	rInUTF8 := transform.NewReader(sr, newTransformer())
	res := make([]byte, 100)
	rInUTF8.Read(res)

	log.Println("t is ", "真的吗", string(res))
}
