package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/guzzlerio/rizo"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/reaandrew/code-named-talula"
)

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

func GetJSON(client *http.Client, url string) map[string]interface{} {
	//Craft a GET request to the proxy for /people
	getRequest, err := http.NewRequest("GET", url, nil)
	Expect(err).To(BeNil())
	getResponse, err := client.Do(getRequest)
	Expect(err).To(BeNil())
	defer getResponse.Body.Close()
	getResponseContent, err := ioutil.ReadAll(getResponse.Body)

	var jsonResponse map[string]interface{}
	jsonResponseError := json.Unmarshal(getResponseContent, &jsonResponse)
	Expect(jsonResponseError).To(BeNil())
	Expect(getResponse.StatusCode).To(Equal(http.StatusOK))
	return jsonResponse
}

var _ = Describe("CodeNamedTalula", func() {

	PIt("Creating a response transform without a body returns a BadRequest", func() {})
	PIt("Creating a response transform without an ID returns a BadRequest", func() {})
	PIt("Creating a response transform with a malformed ID returns a BadRequest", func() {})
	PIt("Creating a response transform with an ID that does not exist returns a NotFound", func() {})

	It("Transforms a JSON response", func() {

		factory := rizo.HTTPResponseFactory(func(w http.ResponseWriter) {
			io.WriteString(w, `{
         "firstname" : "John",
         "lastname" : "Doe",
         "age" : 33
      }`)
		})

		TestServer.Use(factory).For(rizo.RequestWithPath("/people"))

		client := &http.Client{}

		endpointAPIResponse := CreateEndpoint(client, fmt.Sprintf(`{
      "destination" : "http://%s:%d",
      "path" : "/people"
    }`, Host, TestServerPort))

		CreateResponseTransform(client, `function transform(body){
				      return {
				        name : body.firstname + " " + body.lastname,
				        age : body.age
				      }
				    }`, endpointAPIResponse.Links)

		jsonResponse := GetJSON(client, TransformURL("/people"))

		Expect(jsonResponse["name"]).ToNot(BeNil())
		Expect(jsonResponse["age"]).ToNot(BeNil())
		Expect(string(jsonResponse["name"].(string))).To(Equal("John Doe"))
		Expect(jsonResponse["age"].(float64)).To(Equal(float64(33)))

	})
})
