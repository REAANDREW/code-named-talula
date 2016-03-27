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

func findLinkByRel(rel string, links []Link) LinkResult {
	var result = LinkResult{}
	for _, link := range links {
		if link.Rel == rel {
			result.Found = true
			result.Result = link
			break
		}
	}
	return result
}

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
		apiResponse := CreateEndpoint(client, fmt.Sprintf(`{
      "destination" : "http://%s:%d",
      "path" : "/people"
    }`, Host, TestServerPort))

		fmt.Printf("Response %v \n", apiResponse)
		findResult := findLinkByRel("set_response_transform", apiResponse.Links)
		Expect(findResult.Found).To(Equal(true))
		createScriptLink := findResult.Result

		//Add the transformation script
		transformString := []byte(`function transform(body){
				      return {
				        name : body.firstname + " " + body.lastname,
				        age : body.age
				      }
				    }`)
		transformRequest, err := http.NewRequest(createScriptLink.Method, createScriptLink.Href, bytes.NewBuffer(transformString))
		Expect(err).To(BeNil())
		transformResponse, err := client.Do(transformRequest)
		Expect(err).To(BeNil())
		Expect(transformResponse.StatusCode).To(Equal(http.StatusCreated))

		//Craft a GET request to the proxy for /people
		getPeopleRequest, err := http.NewRequest("GET", TransformURL("/people"), nil)
		Expect(err).To(BeNil())
		getPeopleResponse, err := client.Do(getPeopleRequest)
		Expect(err).To(BeNil())
		defer getPeopleResponse.Body.Close()
		getPeopleResponseContent, err := ioutil.ReadAll(getPeopleResponse.Body)

		var jsonResponse map[string]interface{}
		jsonResponseError := json.Unmarshal(getPeopleResponseContent, &jsonResponse)

		Expect(jsonResponseError).To(BeNil())
		Expect(getPeopleResponse.StatusCode).To(Equal(http.StatusOK))
		Expect(jsonResponse["name"]).ToNot(BeNil())
		Expect(jsonResponse["age"]).ToNot(BeNil())
		Expect(string(jsonResponse["name"].(string))).To(Equal("John Doe"))
		Expect(jsonResponse["age"].(float64)).To(Equal(float64(33)))

	})
})
