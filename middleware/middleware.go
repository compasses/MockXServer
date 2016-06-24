package middleware

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
	"time"

	"github.com/Compasses/MockXServer/offline"
	"github.com/Compasses/MockXServer/online"
	"github.com/compasses/MockXServer/db"
	"github.com/compasses/MockXServer/utils"
)

type config struct {
	RunMode      string
	TLS          string
	RemoteServer string
	ListenOn     string
	LogFile      string
	GrabIF       string
}

func GetConfiguration() (conf *config, err error) {
	//get configuration
	file, err := os.Open("./config.json")
	if err != nil {
		log.Println("read file failed...", err)
		log.Println("Just run in offline mode")
		return nil, err
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println("read file failed...", err)
		log.Println("Just run in offline mode")
		return nil, err
	} else {
		json.Unmarshal(data, &conf)
		log.Println("get configuration:", string(data))
	}
	return conf, nil
}

type middleWare struct {
	conf     *config
	handler  http.Handler
	replaydb *db.ReplayDB
	runmode  int
}

func NewMiddleware() *middleWare {
	conf, err := GetConfiguration()
	if err != nil {
		return nil
	}
	// set log format
	if len(conf.LogFile) > 0 {
		f, err := os.OpenFile(conf.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Println("error opening file: %v", err)
		}

		log.SetOutput(f)
		go func(f *os.File) {
			for {
				f.Sync()
				time.Sleep(time.Second)
			}
		}(f)
	} else {
		log.Println("Not assign log file, just print to command window.")
	}
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)

	middleware := new(middleWare)
	middleware.conf = conf
	dbreplay, err := db.NewReplayDB("./ReplayDB")
	dbreplay.ReadDir("./input")
	if err != nil {
		log.Println("Open replayDB error ", err)
		return nil
	}
	middleware.replaydb = dbreplay

	if conf.RunMode == "offline" {
		log.Println("MockXServer Run in offline mode...")
		offlineDB, err := db.NewReplayDB("./OfflineDB")
		if err != nil {
			log.Println("Open OfflineDB error ", err)
			return nil
		}

		middleware.handler = offline.NewServerRouter(offlineDB.GetDB())
		middleware.runmode = 0
	} else {
		log.Println("MockXServer Run in online mode...")
		middleware.handler = online.NewProxyHandler(conf.RemoteServer, conf.GrabIF, middleware.replaydb)
		middleware.runmode = 1
	}

	return middleware
}

func (middle *middleWare) Run() {
	log.Println("Listen ON: ", middle.conf.ListenOn)
	if middle.conf.TLS == "on" {
		log.Fatal(http.ListenAndServeTLS(middle.conf.ListenOn, "cert.pem", "key.pem", middle))
	} else {
		log.Fatal(http.ListenAndServe(middle.conf.ListenOn, middle))
	}
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

func (middleware *middleWare) ServeHTTP(w http.ResponseWriter, req *http.Request) {
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
		middleware.handler.ServeHTTP(w, newRq)
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
