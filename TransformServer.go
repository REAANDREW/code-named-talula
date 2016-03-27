package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

var (
	//TransformServer ...
	TransformServer = gin.Default()
)

//StartTransformServer ...
func StartTransformServer() {
	TransformServer.GET("*path", func(c *gin.Context) {
		c.Status(200)
		urlParsed, err := url.ParseRequestURI(c.Request.RequestURI)
		if err != nil {
			panic(err)
		}
		path := urlParsed.Path

		for _, endpoint := range endpoints {
			if endpoint.Path == path {
				destinationURL, err := url.Parse(endpoint.Destination)
				if err != nil {
					panic(err)
				}
				destinationURL.RawQuery = urlParsed.RawQuery
				destinationURL.Path = urlParsed.Path
				client := &http.Client{}
				proxyRequest, err := http.NewRequest(c.Request.Method, destinationURL.String(), c.Request.Body)
				if err != nil {
					panic(err)
				}

				proxyResponse, err := client.Do(proxyRequest)
				if err != nil {
					panic(err)
				}
				defer proxyResponse.Body.Close()
				proxyResponseBody, err := ioutil.ReadAll(proxyResponse.Body)
				if V8Worker == nil {
					panic("Something gone wrong")
				}
				result := V8Worker.SendSync(string(proxyResponseBody))
				fmt.Println(result)
				c.Data(proxyResponse.StatusCode, proxyResponse.Header.Get("Content-Type"), []byte(result))
			}
		}
	})
	TransformServer.Run(fmt.Sprintf("%s:%d", Host, ProxyPort))
}
