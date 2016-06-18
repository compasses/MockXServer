package utils

import (
	"encoding/json"
	"log"
	"strconv"
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

func JsonNormalizeSingle(body string) (stream []byte, err error) {
	streamMap, err := TOJsonInterface([]byte(body))
	stream, err = json.Marshal(streamMap)
	if err != nil {
		log.Println("Marshal failed ", err, streamMap)
	}
	return
}

func JsonNormalize(reqBody, respBody string, statusCode int) (reqStram, respStram []byte, err error) {
	reqMap, err := TOJsonInterface([]byte(reqBody))
	respMap, err := TOJsonInterface([]byte(respBody))

	rspStream := map[string]interface{}{
		strconv.Itoa(statusCode): respMap,
	}

	respStram, err = json.Marshal(rspStream)
	if err != nil {
		log.Println("Marshal failed ", err, rspStream)
	}

	reqStram, err = json.Marshal(reqMap)
	if err != nil {
		log.Println("Marshal failed ", err, reqMap)
	}

	return
}
