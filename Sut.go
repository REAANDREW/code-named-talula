package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/gomega"
)

//CreateEndpoint ...
func CreateEndpoint(client *http.Client, data string) APIResponse {
	buffer := bytes.NewBuffer([]byte(data))
	endpointRequest, err := http.NewRequest("POST", AdminURL("/endpoints"), buffer)
	Expect(err).To(BeNil())
	scriptResponse, err := client.Do(endpointRequest)
	Expect(err).To(BeNil())
	Expect(scriptResponse.StatusCode).To(Equal(http.StatusCreated))
	defer scriptResponse.Body.Close()
	endpointResponseBody, err := ioutil.ReadAll(scriptResponse.Body)
	Expect(err).To(BeNil())
	var apiResponse APIResponse
	err = json.Unmarshal(endpointResponseBody, &apiResponse)
	Expect(err).To(BeNil())
	return apiResponse
}

//CreateResponseTransform ...
func CreateResponseTransform(client *http.Client, data string, links []Link) APIResponse {
	findResult := FindLinkByRel("set_response_transform", links)
	Expect(findResult.Found).To(Equal(true))
	createScriptLink := findResult.Result

	//Add the transformation script
	transformString := []byte(data)
	transformRequest, err := http.NewRequest(createScriptLink.Method, createScriptLink.Href, bytes.NewBuffer(transformString))
	Expect(err).To(BeNil())
	transformResponse, err := client.Do(transformRequest)
	Expect(err).To(BeNil())
	Expect(transformResponse.StatusCode).To(Equal(http.StatusCreated))
	defer transformResponse.Body.Close()
	transformResponseBody, err := ioutil.ReadAll(transformResponse.Body)
	var apiResponse APIResponse
	err = json.Unmarshal(transformResponseBody, &apiResponse)
	Expect(err).To(BeNil())
	return apiResponse
}
