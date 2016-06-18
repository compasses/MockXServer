package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/compasses/MockXServer/offline"
	"github.com/compasses/MockXServer/online"
)

type config struct {
	RunMode      string
	TLS          string
	RemoteServer string
	ListenOn     string
	LogFile      string
	GrabIF       string
}

var GlobalServerStatus int64 = 0
var localServer string = "localhost:8080"
var GlobalConfig config

func GetConfiguration() error {
	//get configuration
	file, err := os.Open("./config.json")
	if err != nil {
		log.Println("read file failed...", err)
		log.Println("Just run in offline mode")
		return err
	}
	var conf config
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println("read file failed...", err)
		log.Println("Just run in offline mode")
		return err
	} else {
		json.Unmarshal(data, &conf)
		log.Println("get configuration:", string(data))
	}
	GlobalConfig = conf
	return nil
}

func RunDefaultServer(handler http.Handler) {
	log.Println("Listen ON: ", GlobalConfig.ListenOn)
	if GlobalConfig.TLS == "on" {
		log.Fatal(http.ListenAndServeTLS(GlobalConfig.ListenOn, "cert.pem", "key.pem", handler))
	} else {
		log.Fatal(http.ListenAndServe(GlobalConfig.ListenOn, handler))
	}
}

func StartServer() {
	err := GetConfiguration()

	if err != nil {
		log.Println("Run in default, server: ", localServer, "offline, on http")
		router := offline.ServerRouter()
		log.Fatal(http.ListenAndServe(localServer, router))
		return
	}

	if len(GlobalConfig.LogFile) > 0 {
		f, err := os.OpenFile(GlobalConfig.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
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

	log.Println("Begin API LOG------------------------")
	if GlobalConfig.RunMode == "offline" {
		log.Println("API Run in offline mode...")
		router := offline.NewMiddleware()
		RunDefaultServer(router)
	} else {
		log.Println("API Run in online mode...")
		proxy := online.NewProxyHandler(GlobalConfig.RemoteServer, GlobalConfig.GrabIF)
		RunDefaultServer(proxy)
	}
}

const banner string = `

			Mock Server

`

func main() {
	log.Println(banner)
	log.Printf("Git commit:%s\n", Version)
	log.Printf("Build time:%s\n", Compile)
	// router := offline.NewMiddleware()
	// router.TestTruncate()
	StartServer()
}
