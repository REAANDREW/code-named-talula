package main

import (
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
)

var (
	//AdminServer ...
	AdminServer = gin.Default()
)

//StartAdminServer ...
func StartAdminServer() {
	AdminServer.POST("/endpoints", func(c *gin.Context) {
		var endpoint EndpointDTO
		c.BindJSON(&endpoint)
		newEndpoint := Endpoint{
			ID:          uuid.NewV4(),
			Destination: endpoint.Destination,
			Path:        endpoint.Path,
		}
		endpoints[newEndpoint.ID] = newEndpoint
		response := APIResponse{
			Links: []Link{
				Link{
					Href:   AdminURL("/endpoints/%s/transforms/response", newEndpoint.ID),
					Method: "PUT",
					Rel:    "set_response_transform",
				},
			},
		}
		fmt.Printf("Response Built %v \n", response)
		c.JSON(201, &response)
	})

	AdminServer.PUT("/endpoints/:id/transforms/response", func(c *gin.Context) {
		contentType := c.ContentType()
		payload, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			panic(err)
		}

		uuidValue, err := uuid.FromString(c.Param("id"))
		if err != nil {
			panic(err)
		}

		if val, ok := endpoints[uuidValue]; ok {
			//do something here
			val.ResponseTransform = Transform{
				ContentType: contentType,
				Transform:   string(payload),
			}
		} else {
			panic("Cannot find the endpoint with that UUID")
		}
		response := APIResponse{
			Links: []Link{},
		}
		c.JSON(201, &response)
	})

	AdminServer.Run(fmt.Sprintf("%s:%d", Host, AdminPort))
}
