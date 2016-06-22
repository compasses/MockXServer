package middleware

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/SEEK-Jobs/pact-go"
	"github.com/SEEK-Jobs/pact-go/provider"
	"github.com/compasses/MockXServer/utils"
)

const PactsDir = "./pacts"

type ProviderAPIClient struct {
	baseURL string
}

func (c *ProviderAPIClient) ClientRun(method, path string, reqBody interface{}) error {
	url := fmt.Sprintf("%s%s", c.baseURL, path)
	reqb := utils.JsonInterfaceToByte(reqBody)

	log.Println("going to verify: ", url+" "+method)

	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqb))
	if err != nil {
		log.Println(err)
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Println(err)
		return err
	}
	defer resp.Body.Close()

	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("ioutil read err ", err, " body", string(res))
		return err
	}
	//log.Println("Got response: ", string(res))

	return nil
}

func (middleware *middleWare) GetPactFile() string {
	files, err := ioutil.ReadDir(PactsDir)
	if err != nil {
		log.Println(err)
		return ""
	}
	//now just run the first file
	for _, file := range files {
		log.Println("upload pact file name ", file.Name())
		return PactsDir + "/" + file.Name()
	}
	return ""
}

func (middleware *middleWare) buildPact(consumerName, providerName string) pact.Builder {
	return pact.
		NewConsumerPactBuilder(&pact.BuilderConfig{PactPath: PactsDir}).
		ServiceConsumer(consumerName).
		HasPactWith(providerName)
}

func (middleware *middleWare) RunPact(builder pact.Builder, path, method string, reqBody, respBody interface{}, statusCode int,
	consumerName, providerName string) {
	ms, msUrl := builder.GetMockProviderService()

	request := provider.NewJSONRequest(method, path, "", nil)
	if reqBody != nil {
		err := request.SetBody(reqBody)
		if err != nil {
			log.Println("Set request error ", err, " reqBody ", respBody)
		}
	}

	header := make(http.Header)
	header.Add("content-type", "application/json")
	response := provider.NewJSONResponse(statusCode, header)
	if respBody != nil {
		err := response.SetBody(respBody)
		if err != nil {
			log.Println("Set Response error ", err, " respBody ", respBody)
		}
	}

	//Register interaction for this test scope
	if err := ms.Given(consumerName).
		UponReceiving(providerName).
		With(*request).
		WillRespondWith(*response); err != nil {
		log.Println(err)
	}

	//log.Println("Register: ", " Request ", string(req), " response ", respBody)

	//test
	client := &ProviderAPIClient{baseURL: msUrl}
	if err := client.ClientRun(method, path, reqBody); err != nil {
		log.Println(err)
	}

	//Verify registered interaction
	if err := ms.VerifyInteractions(); err != nil {
		log.Println(err)
	}

	//Clear interaction for this test scope, if you need to register and verify another interaction for another test scope
	ms.ClearInteractions()
}

func (middleware *middleWare) GenPactWithProvider() {
	builder := middleware.buildPact("EShop Online Store", "EShop Adaptor")
	//map[string]map[string][]interface{}
	//"Path", "Method", "[req..., rsp...,]"
	interactMap := middleware.replaydb.GetJSONMap()
	for path, value := range interactMap {
		for method, interacts := range value {
			var count int = 0
			for _, detailMapel := range interacts {
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
				count++
				for k, v := range responseMap {
					status, _ := strconv.Atoi(k)
					//fmt.Println("\r\nstore:", request, "response", v)
					consumName := "mock server for " + path + " method " + method + " " + strconv.Itoa(count)
					provideName := "pact contract for " + path + " method " + method + " " + strconv.Itoa(count)
					middleware.RunPact(builder, path, method, request, v, status, consumName,
						provideName)
					count++
					break
				}
			}
		}
	}

	//Finally, build to produce the pact json file
	if err := builder.Build(); err != nil {
		log.Println(err)
	}
}
