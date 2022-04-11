package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type Message struct {
	Greeting string `json:"greeting"`
}

var (
	wsUpgrader = websocket.Upgrader{
		ReadBufferSize: 1024,
		WriteBufferSize: 1024,
	}

	wsConn *websocket.Conn
)

func WsEndpoint(w http.ResponseWriter, r *http.Request)  {
	defer wsConn.Close()

	wsUpgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	wsConn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("could not upgrade: %s\n", err.Error())
		return
	}

	// even loop
	for {
		var msg Message

		err := wsConn.ReadJSON(&msg)
		if err != nil {
			fmt.Printf("Error reading JSON : %s", err.Error())
			continue
		}

		fmt.Printf("Message received : %s", msg.Greeting)
		sendMessage("Hello Client, From Server!")
	}
}

func sendMessage(message string)  {
	err := wsConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		fmt.Printf("error sending message: %s\n", err.Error())
	}
}

func main()  {

	router := mux.NewRouter()

	router.HandleFunc("/socket", WsEndpoint)

	log.Fatal(http.ListenAndServe(":9100", router))

}