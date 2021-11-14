package main

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"strconv"
)

func main() {
	r := gin.Default()
	r.GET("/test", func(c *gin.Context) {
		log.Println("Inside Test")
		c.JSON(200, gin.H{"Message": "This is test"})
	})
	r.GET("/webhook", func(c *gin.Context) {
		query := c.Query("hub.challenge")
		log.Println("Inside Get, Token: ", query)
		intVar, err := strconv.ParseInt(query, 0, 64)
		if err != nil {
			log.Println("Error: ", err)
		}
		c.JSON(200, intVar)
	})
	r.POST("/webhook", func(c *gin.Context) {
		jsonData, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			log.Println("Error: ", err)
		}
		log.Println("Response: ", jsonData)
		c.JSON(200, "")
	})
	r.GET("/", func(c *gin.Context) {
		log.Println("Inside Home")
		c.JSON(200, gin.H{"Message": "This is Home"})
	})
	err := r.Run()
	if err != nil {
		return
	}
}
