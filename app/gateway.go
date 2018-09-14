package app

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

func CallService(c *gin.Context) {

	serviceHost := c.MustGet("serviceURL").(string)

	//Create client
	client := &http.Client{}

	//Check for request body
	var bodyBytes []byte
	var err error
	if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut {
		if c.Request.Body != nil {
			bodyBytes, err = ioutil.ReadAll(c.Request.Body)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong when reading request body", "detail": err.Error()})
				return
			}
		}
	}

	//Create request
	proxyReq, err := http.NewRequest(c.Request.Method, serviceHost+c.Request.URL.Path, bytes.NewReader(bodyBytes))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong while creating the request to the service", "detail": err.Error()})
		return
	}

	//Copy headers
	proxyReq.Header = make(http.Header)
	for h, val := range c.Request.Header {
		proxyReq.Header[h] = val
	}

	//Add query parameters
	q := proxyReq.URL.Query()
	for key, value := range c.Request.URL.Query() {
		q.Add(key, value[0])
	}
	proxyReq.URL.RawQuery = q.Encode()

	//Fetch Request
	response, err := client.Do(proxyReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong calling the service", "detail": err.Error()})
		return
	}
	defer response.Body.Close()

	//Parse response
	var result map[string]interface{}
	json.NewDecoder(response.Body).Decode(&result)

	c.JSON(response.StatusCode, result)
	return

}
