package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"secret-im/service/signalserver/cmd/api/textsecure"
	"time"
	"github.com/golang/protobuf/proto"
)

func main()  {

	u := url.URL{Scheme: "ws", Host: "127.0.0.1:38888", Path: "/v1/websocket"}
	log.Println("connecting to ", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()
	done := make(chan struct{})


	wreq:=&textsecure.WebSocketRequestMessage{}
	wreq.Body=[]byte(`{
						"destination": "B",
						"messages": [
							{
								"content": "fsdf.....",
								"destinationDeviceId": 1,
								"destinationRegistrationId": 1577,
								"type": 6
							}
						],
						"online": false,
						"timestamp": 1623151469837
					}`)
	wreq.Id=1
	wreq.Path="/api/v1/message"
	wreq.Verb="PUT"
	w:=&textsecure.WebSocketMessage{}
	w.Type=textsecure.WebSocketMessage_REQUEST
	w.Request=wreq
	log.Println("webSocketMessage is ",w.String())
	msg,err:=proto.Marshal(w)
	if err!=nil{
		log.Println("proto.Marshal(w) error:",err)
	}else{
		c.WriteMessage(websocket.BinaryMessage,msg)
	}




	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	time.Sleep(time.Minute)

}
