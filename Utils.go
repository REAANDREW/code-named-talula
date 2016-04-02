package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/satori/go.uuid"
)

func SafeUUID(uuidValue uuid.UUID) string {
	return strings.Replace(uuidValue.String(), "-", "_", -1)
}

//AdminURL ...
func AdminURL(path string, args ...interface{}) string {
	return fmt.Sprintf("http://%s:%d%s", Host, AdminPort, fmt.Sprintf(path, args...))
}

//TransformURL ...
func TransformURL(path string, args ...interface{}) string {
	return fmt.Sprintf("http://%s:%d%s", Host, ProxyPort, fmt.Sprintf(path, args...))
}

//JSONResponse ...
type JSONResponse struct {
	JSON       map[string]interface{}
	StatusCode int
}

//GetJSON ...
func GetJSON(client *http.Client, url string) (JSONResponse, error) {
	//Craft a GET request to the proxy for /people
	getRequest, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return JSONResponse{}, err
	}
	getResponse, err := client.Do(getRequest)
	if err != nil {
		return JSONResponse{}, err
	}
	defer getResponse.Body.Close()
	getResponseContent, err := ioutil.ReadAll(getResponse.Body)
	if err != nil {
		return JSONResponse{}, err
	}

	var jsonResponse map[string]interface{}
	jsonResponseError := json.Unmarshal(getResponseContent, &jsonResponse)
	if jsonResponseError != nil {
		return JSONResponse{}, err
	}
	return JSONResponse{
		JSON:       jsonResponse,
		StatusCode: getResponse.StatusCode,
	}, nil
}
