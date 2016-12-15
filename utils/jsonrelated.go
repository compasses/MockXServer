package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
)

func TOJsonInterface(body []byte) (result interface{}, err error) {
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Println("Unmarshal failed ", err, body)
	}
	return
}

func JsonInterfaceToByte(body interface{}) (result []byte) {
	var err error
	result, err = json.Marshal(body)
	if err != nil {
		log.Println("Marshal failed ", err, body)
	}
	return
}

func JsonInterfaceToByteByNumber(body []byte) (result []byte) {
	d := json.NewDecoder(strings.NewReader(string(body)))
	d.UseNumber()
	var x interface{}
	if err := d.Decode(&x); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("decoded to %#v\n", x)
	result, err := json.Marshal(x)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func JsonNormalizeSingle(body string) (stream []byte, err error) {
	streamMap, err := TOJsonInterface([]byte(body))
	stream, err = json.Marshal(streamMap)
	if err != nil {
		log.Println("Marshal failed ", err, streamMap)
	}
	return
}

func JsonNormalize(reqBody, respBody string, statusCode int) (reqStram, respStram []byte, err error) {
	// reqMap, err := TOJsonInterface([]byte(reqBody))
	// respMap, err := TOJsonInterface([]byte(respBody))

	rspStream := map[string]string{
		strconv.Itoa(statusCode): respBody,
	}

	respStram, err = json.Marshal(rspStream)
	if err != nil {
		log.Println("Marshal failed ", err, rspStream)
	}

	reqStram = []byte(reqBody)
	//
	// reqStram, err = json.Marshal(reqMap)
	// if err != nil {
	// 	log.Println("Marshal failed ", err, reqMap)
	// }

	return
}
