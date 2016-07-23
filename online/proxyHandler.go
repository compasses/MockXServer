package online

import (
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/compasses/MockXServer/db"
	. "github.com/compasses/MockXServer/utils"
)

type ProxyRoute struct {
	client *http.Client
	url    string
	db     *db.ReplayDB
}

func NewProxyHandler(newurl string, db *db.ReplayDB) *ProxyRoute {
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	return &ProxyRoute{
		client: &http.Client{Transport: tr},
		url:    newurl,
		db:     db}
}

func (proxy *ProxyRoute) doReq(NeedLog bool, path, method, requestBody string, newRq *http.Request) (resp *http.Response, res []byte) {
	now := time.Now()
	resp, err := proxy.client.Do(newRq)
	if resp != nil {
		defer resp.Body.Close()
	}

	LogOutPut(NeedLog, "Time used: ", time.Since(now))
	if err != nil {
		log.Println("get error ", err)
	} else {
		if resp.Header.Get("Content-Encoding") == "gzip" {
			resp.Body, err = gzip.NewReader(resp.Body)
			if err != nil {
				panic(err)
			}
		}

		res, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("ioutil read err ", err)
		}

		if NeedLog {
			if resp.StatusCode == 500 || resp.StatusCode == 404 {
				FailNum++
			} else {
				SuccNum++
				err = proxy.db.StoreRequest(path, method, requestBody, string(res), resp.StatusCode)
			}
		}

		LogOutPut(NeedLog, "Get response : ")
		ResponseFormat(NeedLog, resp, string(res))

		if err != nil {
			log.Println("Store data failed ", err)
		}
	}
	return
}

func (proxy *ProxyRoute) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	newbody, err := ioutil.ReadAll(req.Body)
	//make([]byte, req.ContentLength)
	//req.Body.Read(newbody)
	if err != nil {
		log.Println("Read request failed..", err)
		return
	}

	newRq, err := http.NewRequest(req.Method, proxy.url+req.RequestURI, ioutil.NopCloser(bytes.NewReader(newbody)))
	if err != nil {
		log.Println("new request error ", err)
	}

	newRq.Header = req.Header
	path := strings.Split(req.RequestURI, "?")

	LogOutPut(true, "online handle, New Request: ")
	RequstFormat(true, newRq, string(newbody))
	resphttp, res := proxy.doReq(true, path[0], req.Method, string(newbody), newRq)
	for key, _ := range resphttp.Header {
		w.Header().Set(key, strings.Join(resphttp.Header[key], ";"))
	}

	w.Write(res)

}
