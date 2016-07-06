package db

import (
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/compasses/MockXServer/utils"
)

func (replay *ReplayDB) ReadJsonFiles(filePath string) {
	now := time.Now()
	log.Println("going to read file", filePath)
	stream, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Println("going to read file error happened ", err, stream)
		return
	}

	res, err := utils.TOJsonInterface(stream)
	paths := res.(map[string]interface{})["paths"]
	pathsMap := paths.(map[string]interface{})
	for path, val := range pathsMap {
		valMap := val.(map[string]interface{})
		for method, detail := range valMap {
			detailMap := detail.([]interface{})
			for _, detailMapel := range detailMap {
				detailMapItem := detailMapel.(map[string]interface{})
				request, ok := detailMapItem["request"]
				if !ok {
					log.Println("missing request, continue ", detailMapItem)
					continue
				}
				respose, ok := detailMapItem["response"]
				if !ok {
					log.Println("missing response, continue ", detailMapItem)
					continue
				}
				responseMap := respose.(map[string]interface{})
				for k, v := range responseMap {
					status, _ := strconv.Atoi(k)
					//fmt.Println("\r\nstore:", request, "response", v)
					replay.StoreRequestFromJson(path, method, request, v, status)
					break
				}
			}
		}
	}
	log.Println("Time Used:", time.Since(now))
}

func (replay *ReplayDB) ReadDir(dir string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Println(err)
		return
	}

	for _, file := range files {
		if strings.Contains(file.Name(), "Incomplete") {
			continue
		}
		replay.ReadJsonFiles(dir + "/" + file.Name())
	}
}
