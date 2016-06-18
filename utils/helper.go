package utils

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
)

type TableId int64
type IDSeqs []TableId
type StrSeqs []string

//used to avoid recurision
type idseqs IDSeqs

var SuccNum int = 0
var FailNum int = 0

func (ids *IDSeqs) UnmarshalJSON(b []byte) (err error) {
	log.Println("IDSeqs got bytes: ", string(b))
	var tids *idseqs
	if err = json.Unmarshal(b, &tids); err == nil {
		*ids = IDSeqs(*tids)
		return
	}

	strid := new(StrSeqs)
	if err = json.Unmarshal(b, strid); err == nil {
		for _, val := range *strid {
			if v, err := strconv.ParseInt(val, 10, 64); err == nil {
				*ids = append(*ids, TableId(v))
			}
		}
		return
	} else {
		log.Println("got error: ", err)
	}
	return
}

func (id *TableId) UnmarshalJSON(b []byte) (err error) {
	log.Println("TableId got bytes: ", string(b))
	var tabd int64
	if err = json.Unmarshal(b, &tabd); err == nil {
		*id = TableId(tabd)
		return
	}

	s := ""
	if err = json.Unmarshal(b, &s); err == nil {
		v, err := strconv.ParseInt(s, 10, 64)
		if err == nil {
			*id = TableId(v)
		}
	}

	return
}

func GetAddressObj(addr interface{}) map[string]interface{} {
	old := addr.(map[string]interface{})

	var result = map[string]interface{}{
		"recipientName": old["userName"],
		"cityName":      old["customCity"],
		"stateId":       old["state"],
		"countryId":     old["country"],
		"street1":       old["address1"],
		"street2":       old["address2"],
		"zipCode":       old["zipCode"],
		"mobile":        old["phone"],
		"state":         old["customState"],
	}

	return result
}

func ToInt64FromString(input string) int64 {
	re, _ := strconv.ParseInt(input, 10, 8)
	return re
}

func ToFloat64FromString(input string) float64 {
	re, _ := strconv.ParseFloat(input, 8)
	return re
}

func GetSliceIntFromBytes(input []byte) []TableId {
	sizeofInt := 8
	data := make([]TableId, len(input)/sizeofInt)
	buf := bytes.NewBuffer(input)
	for i := range data {
		var re int64
		binary.Read(buf, binary.LittleEndian, &re)
		data[i] = TableId(re)
	}

	return data
}

func GetSliceBytesFromInts(input []TableId) []byte {
	buf := new(bytes.Buffer)

	for i := range input {
		binary.Write(buf, binary.LittleEndian, int64(input[i]))
	}
	return buf.Bytes()
}

func ContainsIntSlice(s []TableId, e TableId) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func (tId TableId) ToBytes() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, int64(tId))
	return buf.Bytes()
}

func (tId TableId) ToString() (str string) {
	str = strconv.FormatInt(int64(tId), 10)
	return
}

func (tId TableId) ToInt() int64 {
	return int64(tId)
}

func ToInt64FromBytes(st []byte) int64 {
	buf := bytes.NewReader(st)
	var result int64
	binary.Read(buf, binary.LittleEndian, &result)
	return result
}

//proc string like "createcustomernew(1)", and return 1
func GetIdFromStr(input string) TableId {
	valId := regexp.MustCompile(`(\d+)`)
	val, _ := strconv.Atoi(valId.FindString(input))
	return TableId(val)
}

func RetrieveByMapLevel(f interface{}, levels []string) interface{} {
	length := len(levels)

	if length <= 0 {
		return f
	}

	result := f.(map[string]interface{})

	for k, v := range result {
		if k == levels[0] {
			return RetrieveByMapLevel(v.(interface{}), levels[1:])
		}
	}

	return result
}

func HandleDecodeData(f map[string]interface{}) {
	for k, v := range f {
		log.Println("key", k)
		switch vv := v.(type) {
		//		case string:
		//			log.Println(k, "is string", vv)
		//		case int:
		//			log.Println(k, "is int", vv)
		//		case []interface{}:
		//			log.Println(k, "is an array:")
		//			for i, u := range vv {
		//				log.Println(i, u)
		//			}
		//		case interface{}:
		//			log.Println(k, "is an interface")
		//			HandleDecodeData(f[k].(map[string]interface{}))
		default:
			log.Println(vv)
			log.Println(k, "is of a type I don't know how to handle")
		}
	}
}

func HandleError(err error) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}
	log.Println(err, file, "line:", line)
}

func LogOutPut(NeedLog bool, v ...interface{}) {
	if NeedLog {
		log.Println(v)
	}
}

func RequstFormat(NeedLog bool, req *http.Request, newbody string) {
	if !NeedLog {
		return
	}
	result := "URL: " + req.URL.String() + "\r\n"
	result += "Method: " + req.Method + "\r\n"
	result += "Body: " + newbody + "\r\n"
	result += "Header: "
	for key, _ := range req.Header {
		var vals string = ""
		for _, allV := range req.Header[key] {
			vals += allV
		}
		result += " Key: " + key + " -> " + vals + "\r\n"
	}
	LogOutPut(NeedLog, result)
}

func ResponseFormat(NeedLog bool, resp *http.Response, body string) {
	if !NeedLog {
		return
	}

	result := "Status: " + resp.Status + "\r\n"
	result += "Body: " + body + "\r\n"
	result += "Header: "
	for key, _ := range resp.Header {
		var vals string = ""
		for _, allV := range resp.Header[key] {
			vals += allV
		}
		result += " Key: " + key + " -> " + vals + "\r\n"
	}
	LogOutPut(NeedLog, result)
}

func ReflectStruct(req *http.Request) {
	s := reflect.ValueOf(req).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		log.Printf("%d: %s %s = %v\n", i,
			typeOfT.Field(i).Name, f.Type(), f.Interface())
	}
}
