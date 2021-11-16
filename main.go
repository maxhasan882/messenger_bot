package main

import (
	"MessengerBot/soc"
	"encoding/json"
	"flag"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

type Payload struct {
	Url       string `json:"url"`
	StickerId int64  `json:"sticker_id"`
}

type Attachments struct {
	Type    string  `json:"type"`
	Payload Payload `json:"payload"`
}

type Sender struct {
	ID string `json:"id"`
}

type Message struct {
	Mid         string        `json:"mid"`
	Text        string        `json:"text"`
	Attachments []Attachments `json:"attachments"`
}

type Messaging struct {
	Sender  Sender  `json:"sender"`
	Message Message `json:"message"`
}

type Entry struct {
	ID        string      `json:"id"`
	Time      int64       `json:"time"`
	Messaging []Messaging `json:"messaging"`
}

type MessageObject struct {
	Object string  `json:"object"`
	Entry  []Entry `json:"entry"`
}

func main() {
	dataChan := make(chan soc.MessageData)
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
		var requestBody MessageObject
		//buf := new(bytes.Buffer)
		//buf.ReadFrom(c.Request.Body)
		//str := buf.String()
		//log.Println(str)
		err := c.Bind(&requestBody)
		if err != nil {
			log.Println("Bing Error: ", err.Error())
		}
		for _, entry := range requestBody.Entry {
			for _, data := range entry.Messaging {
				if data.Message.Text != "" {
					rawData, _ := json.Marshal(soc.RawData{Data: data.Message.Text, Type: "text"})
					messageData := soc.MessageData{RoomId: data.Sender.ID, Data: rawData}
					dataChan <- messageData
				} else {
					var attachments []string
					for _, attachment := range data.Message.Attachments {
						attachments = append(attachments, attachment.Payload.Url)
					}
					rawData, _ := json.Marshal(soc.RawData{Attachments: attachments, Type: "attachment"})
					messageData := soc.MessageData{RoomId: data.Sender.ID, Data: rawData}
					dataChan <- messageData
				}
			}
		}
		log.Println("Message: ", requestBody.Entry[0].Messaging[0].Message.Text, "Client ID: ",
			requestBody.Entry[0].Messaging[0].Sender.ID)
		c.JSON(200, "")
	})
	r.GET("/", func(c *gin.Context) {
		log.Println("Inside Home")
		c.JSON(200, gin.H{"Message": "This is Home"})
	})
	flag.Parse()
	hub := soc.NewHub()
	go hub.Run()
	r.LoadHTMLFiles("index.html")
	r.GET("/room/:roomId", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	r.GET("/ws/:roomId", func(c *gin.Context) {
		roomId := c.Param("roomId")
		soc.ServeWs(hub, c.Writer, c.Request, roomId, dataChan)
	})
	err := r.Run()
	if err != nil {
		return
	}
}
