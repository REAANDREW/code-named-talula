package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os/exec"
	"path/filepath"
	"time"

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

var _ = Describe("CodeNamedTalula", func() {

	PIt("Creating a response transform without a body returns a BadRequest", func() {})
	PIt("Creating a response transform without an ID returns a BadRequest", func() {})
	PIt("Creating a response transform with a malformed ID returns a BadRequest", func() {})
	PIt("Creating a response transform with an ID that does not exist returns a NotFound", func() {})

	It("Transforms a JSON response", func() {
		/*
		   Setup a test http server on port 4000 which will return a JSON response of
		   {
		        "first name" : "John",
		        "last name" : "Doe",
		        "age" : 33
		    }
		*/

		testPort := 4000
		testServer := rizo.CreateRequestRecordingServer(testPort)

		factory := rizo.HTTPResponseFactory(func(w http.ResponseWriter) {
			io.WriteString(w, `{
         "firstname" : "John",
         "lastname" : "Doe",
         "age" : 33
      }`)
		})

		testServer.Use(factory).For(rizo.RequestWithPath("/people"))
		testServer.Start()
		defer func() {
			testServer.Stop()
		}()

		/* Start the application listening */

		exePath, err := filepath.Abs("./code-named-talula")
		if err != nil {
			panic(err)
		}

		cmd := exec.Command(exePath)
		stdout, err := cmd.StdoutPipe()
		Expect(err).To(BeNil())
		stderr, err := cmd.StderrPipe()
		Expect(err).To(BeNil())
		defer func() {
			cmd.Process.Kill()

			stdoutOutput, err := ioutil.ReadAll(stdout)
			if err != nil {
				panic(err)
			}
			fmt.Println("-----------STDOUT")
			fmt.Println(string(stdoutOutput))

			stderrOutput, err := ioutil.ReadAll(stderr)
			if err != nil {
				panic(err)
			}
			fmt.Println("-----------STDERR")
			fmt.Println(string(stderrOutput))
		}()

		if err != nil {
			panic(err)
		}

		cmd.Start()
		time.Sleep(time.Second * 1)

		/*
		   Configure the application to:
		   - proxy traffic for path http://localhost:3000/people to http://localhost:4000/people
		   - transform the response so that it returns
		     {
		       "name" : "John Doe",
		       "age" : 33
		     }
		*/

		client := &http.Client{}
		bodyString := []byte(fmt.Sprintf(`{
      "destination" : "http://%s:%d",
      "path" : "/people"
    }`, Host, testPort))
		endpointRequest, err := http.NewRequest("POST", AdminURL("/endpoints"), bytes.NewBuffer(bodyString))
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
