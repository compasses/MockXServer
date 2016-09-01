package middleware

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/compasses/MockXServer/db"
	"github.com/compasses/MockXServer/offline"
	"github.com/compasses/MockXServer/online"
	"github.com/compasses/MockXServer/utils"
)

type config struct {
	RunMode      string
	TLS          string
	RemoteServer string
	ListenOn     string
	LogFile      string
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
		middleware.handler = online.NewProxyHandler(conf.RemoteServer, middleware.replaydb)
		middleware.runmode = 1
	}

	return middleware
}

func (middleware *middleWare) PreHandle(w http.ResponseWriter, req *http.Request) bool {
	path := strings.Split(req.RequestURI, "?")
	log.Println("prehandle path: ", req.RequestURI)

	if path[0] == "/json" {
		middleware.GenerateJSON(w)
		return true
	} else if path[0] == "/pact" {
		middleware.GeneratePACT(w)
		return true
	} else if path[0] == "/truncate" {
		middleware.Truncate(w)
		return true
	}

	return false
}

func (middle *middleWare) Run() {
	log.Println("Listen ON: ", middle.conf.ListenOn)
	if middle.conf.TLS == "on" {
		log.Fatal(http.ListenAndServeTLS(middle.conf.ListenOn, "cert.pem", "key.pem", middle))
	} else {
		log.Fatal(http.ListenAndServe(middle.conf.ListenOn, middle))
	}
}

func (middleware *middleWare) HandleOffline(w http.ResponseWriter, req *http.Request) {
	newbody := make([]byte, req.ContentLength)
	req.Body.Read(newbody)
	path := strings.Split(req.RequestURI, "?")
	log.Println("offline handle , try to get ", path[0], req.Method, string(newbody))

	res, err := middleware.replaydb.GetResponse(path[0], req.Method, string(newbody))
	if err != nil || res == nil {
		log.Println("Cannot get response from replaydb on offline mode, need hanle in offline/online handler ", err)
		middleware.SaveNotFound(path[0], req.Method, string(newbody), "...xxx...", 200)

		newRq, err := http.NewRequest(req.Method, req.RequestURI, ioutil.NopCloser(bytes.NewReader(newbody)))
		if err != nil {
			log.Println("new http request failed ", err)
		}
		utils.RequstFormat(true, newRq, string(newbody))
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

func (middleware *middleWare) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	preHandled := middleware.PreHandle(w, req)
	if preHandled {
		return
	}

	if middleware.conf.RunMode == "online" {
		middleware.handler.ServeHTTP(w, req)
	} else {
		// firstly try to get from replaydb
		middleware.HandleOffline(w, req)
	}
}
