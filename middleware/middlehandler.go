package middleware

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/Compasses/MockXServer/utils"
)

func (middleware *middleWare) SaveNotFound(path, method, reqBody, respBody string, statusCode int) (err error) {
	finalReq, finalResp, err := utils.JsonNormalize(reqBody, respBody, statusCode)
	if err != nil {
		log.Println("JSON Normalize error ", err)
	}

	filename := strings.Replace(path+"Incomplete.json", "/", "_", -1)
	filename = "./input/" + filename

	result := map[string]interface{}{
		path: map[string][]interface{} {
    {
    method: map[string]interface{} {
    "request": finalReq,
    "response":finalResp,
  }
  },
  },
	}

	jsonStr, err := json.MarshalIndent(result, "", "    ")

	err = ioutil.WriteFile(filename, jsonStr, 0666)
	if err != nil {
		log.Println(err)
	}
	return nil
}

func (middleware *middleWare) returnFile(w http.ResponseWriter, filename, attachName string) {
	f, err := os.Open(filename)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	fileInfo, err := f.Stat()
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", `attachment; filename=`+attachName)
	w.Header().Set("Content-Length", strconv.Itoa(int(fileInfo.Size())))
	io.Copy(w, f)
}

func (middleware *middleWare) TestTruncate() {
	dbFile := middleware.replaydb.GetDBFilePath()
	middleware.replaydb.Close()
	fmt.Println("DB file ", dbFile)
	db, err := os.OpenFile(dbFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Println("OPen error ", err)
	}
	err = db.Truncate(0)
	if err != nil {
		fmt.Println("Truncate error ", err)
	}
}

func (middleware *middleWare) Truncate(w http.ResponseWriter) {
	dbFile := middleware.replaydb.GetDBFilePath()
	middleware.replaydb.Close()
	db, err := os.OpenFile(dbFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		log.Println("OPen error ", err)
	}
	err = db.Truncate(0)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		log.Println("Truncate error ", err)
	}
	w.WriteHeader(200)
	w.Write([]byte("<h1>Truncate successful</h1>"))
}

func (middleware *middleWare) TestPactGen() {
	middleware.replaydb.ReadDir("./input")
	middleware.GenPactWithProvider()
}

func (middleware *middleWare) GenerateJSON(w http.ResponseWriter) {
	middleware.replaydb.ReadDir("./input")
	filename := middleware.replaydb.SerilizeToFile()
	middleware.returnFile(w, filename, "JSONFile.json")
}

func (middleware *middleWare) GeneratePACT(w http.ResponseWriter) {
	middleware.replaydb.ReadDir("./input")
	middleware.GenPactWithProvider()
	pactfile := middleware.GetPactFile()
	if len(pactfile) > 0 {
		middleware.returnFile(w, pactfile, "ConsumerContracts.json")
	} else {
		log.Println("No Pact File Generate ")
		w.Write([]byte("No Pact File Generate "))
	}
}
