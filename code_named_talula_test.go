package main_test

import (
	"fmt"
	"io"
	"net/http"

	"github.com/guzzlerio/rizo"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/reaandrew/code-named-talula"
)

var _ = Describe("JSON Transform", func() {

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

		jsonResponse, err := GetJSON(client, TransformURL("/people"))
		Expect(err).To(BeNil())
		Expect(jsonResponse.StatusCode).To(Equal(http.StatusOK))

		Expect(jsonResponse.JSON["name"]).ToNot(BeNil())
		Expect(jsonResponse.JSON["age"]).ToNot(BeNil())
		Expect(string(jsonResponse.JSON["name"].(string))).To(Equal("John Doe"))
		Expect(jsonResponse.JSON["age"].(float64)).To(Equal(float64(33)))

	})
})
