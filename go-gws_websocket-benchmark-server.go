package main

import (
	"github.com/lxzan/gws"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

var connectionsCounter int = 0

//Struct for events
type event_data struct {
	C  int32 `json:"c"`
	Ts int64 `json:"ts"`
}

type Handler struct {
	gws.BuiltinEventHandler
}

// Starts a http listening service on the given port
// and passes any requests to the entrypoint to the "serve" callback function
func main() {
	log.SetFlags(0)

	//settings for the http service
	var addr = flag.String("addr", "0.0.0.0:8080", "http service address")

	//define entrypoint as well as the callback function which will handle requests
	http.HandleFunc("/", serve)

	//start listening for incoming connections
	log.Fatal(http.ListenAndServe(*addr, nil))
}

// Called once per incoming connection
// Handles events like when a new client connects
// and when the server receives a message from the client
func serve(w http.ResponseWriter, r *http.Request) {

	handler := new(Handler)

	upgrader := gws.NewUpgrader(handler, &gws.ServerOption{
		// CompressEnabled: true,
	})

	//upgrade the connection from a HTTP connection to a websocket connection
	c, err := upgrader.Upgrade(w, r)
	if err != nil {
		log.Println(err.Error())
		return
	}

	connectionsCounter++
	log.Printf("New connection (%d)", connectionsCounter)

	// send newly connected client initial timestamp
	err = notify(c, 0)
	if errorCheck(err) {
		return
	}
	go func() {
		c.ReadLoop()
	}()

	//keep connection open
	// defer c.Close()

	//continuously listen for incoming messages
	// for {

	// 	// read in incoming messages
	// 	mt, message, err := c.ReadMessage()
	// 	_ = mt
	// 	if errorCheck(err) {
	// 		return
	// 	}

	// 	log.Printf("recv: %s", message)

	// 	// decode incoming message into a struct
	// 	var json_data event_data
	// 	err = json.Unmarshal(message, &json_data)
	// 	if errorCheck(err) {
	// 		return
	// 	}

	// 	// notify client with event for message with count "c"
	// 	notify(c, json_data.C)
	// }
}

// Checks for errors after an operation
// Return - bool - true if an error occurred
func errorCheck(err error) bool {
	if err != nil {
		// if error is not nil print error
		fmt.Println(err)
		return true
	}
	return false
}

// Gets the current unix timestamp of the server
// Return - int64 -The current unix timestamp
func getTimestamp() int64 {
	return time.Now().Unix()
}

// Creates a JSON string containing the message count and the current timestamp
// Param - c - int32 - The message count
// Return - []byte - A JSON string (byte array) containing the message count and the current timestamp
func getEvent(c int32) []byte {
	// create an event struct for the time that message "c" is received by the server
	var event event_data
	event.C = c
	event.Ts = getTimestamp()

	// convert json struct into a byte array
	event_string, err := json.Marshal(event)
	if err != nil {
		fmt.Println(err)
		return []byte{}
	}

	return event_string
}

// Send a connected client an event JSON string
// Param - ws - *gws.Conn The client connection the outgoing message is for
// Param - c - int32 - The message count
// Return - error - Error object containing a possible error that occurred
func notify(ws *gws.Conn, c int32) error {

	// send the given connection the event timestamp for message "c"
	return ws.WriteMessage(1, getEvent(c))
}

func (c *Handler) OnMessage(socket *gws.Conn, message *gws.Message) {
	defer message.Close()

	// log.Printf("recv: %s", message)

	// decode incoming message into a struct
	var json_data event_data
	err := json.Unmarshal(message.Data.Bytes(), &json_data)
	if errorCheck(err) {
		return
	}

	// notify client with event for message with count "c"
	notify(socket, json_data.C)

	// socket.WriteMessage(message.Opcode, message.Bytes())
	// _ = message.Close()
}