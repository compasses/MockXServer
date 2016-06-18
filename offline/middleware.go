package offline

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/compasses/MockXServer/db"
	"github.com/compasses/MockXServer/utils"
	"github.com/boltdb/bolt"
	"github.com/julienschmidt/httprouter"
)

type offlinemiddleware struct {
	router   *httprouter.Router
	replaydb *db.ReplayDB
}

func NewMiddleware() *offlinemiddleware {
	router := httprouter.New()
	db, err := db.NewReplayDB()
	db.ReadDir("./input")
	if err != nil {
		log.Println("Open replayDB error ", err)
	}

	for _, route := range routes {
		httpHandle := Logger(route.HandleFunc, route.Name)

		router.Handle(
			route.Method,
			route.Pattern,
			httpHandle,
		)
	}

	router.NotFound = LoggerNotFound(NotFoundHandler)
	GlobalDB, err = bolt.Open("./OfflineDB", 0666, nil)
	if err != nil {
		log.Println("Open OfflineDB error ", err)
	}
	return &offlinemiddleware{
		router:   router,
		replaydb: db,
	}
}

func (middleware *offlinemiddleware) returnFile(w http.ResponseWriter, filename, attachName string) {
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

func (middleware *offlinemiddleware) TestTruncate() {
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

func (middleware *offlinemiddleware) Truncate(w http.ResponseWriter) {
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

func (middleware *offlinemiddleware) TestPactGen() {
	middleware.replaydb.ReadDir("./input")
	middleware.GenPactWithProvider()
}

func (middleware *offlinemiddleware) GenerateJSON(w http.ResponseWriter) {
	middleware.replaydb.ReadDir("./input")
	filename := middleware.replaydb.SerilizeToFile()
	middleware.returnFile(w, filename, "JSONFile.json")
}

func (middleware *offlinemiddleware) GeneratePACT(w http.ResponseWriter) {
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

func (middleware *offlinemiddleware) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	newbody := make([]byte, req.ContentLength)
	req.Body.Read(newbody)
	path := strings.Split(req.RequestURI, "?")

	log.Println("try to get ", path[0], req.Method, string(newbody))
	if path[0] == "/json" {
		middleware.GenerateJSON(w)
		return
	} else if path[0] == "/pact" {
		middleware.GeneratePACT(w)
		return
	} else if path[0] == "/truncate" {
		middleware.Truncate(w)
		return
	}

	res, err := middleware.replaydb.GetResponse(path[0], req.Method, string(newbody))
	if err != nil || res == nil {
		log.Println("Cannot get response from replaydb on offline mode, need hanle in offline handler ", err)
		newRq, err := http.NewRequest(req.Method, req.RequestURI, ioutil.NopCloser(bytes.NewReader(newbody)))
		if err != nil {
			log.Println("new http request failed ", err)
		}
		middleware.router.ServeHTTP(w, newRq)
	} else {
		result, _ := utils.TOJsonInterface(res)
		log.Println("Get response from replaydb on offline mode ", (result))
		resultmap := result.(map[string]interface{})
		for key, value := range resultmap {
			status, _ := strconv.Atoi(key)
			w.WriteHeader(status)
			stream := []byte("")
			if value != nil {
				stream, err = json.Marshal(value)
				if err != nil {
					log.Println("Marshal failed ", err, stream)
				}
			}
			_, err = w.Write(stream)
			if err != nil {
				log.Println("Get response from replaydb  but write error ", err)
			}
			break
		}
	}
}

func ServerRouter() *httprouter.Router {
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

	return router
}
