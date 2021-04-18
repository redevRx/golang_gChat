package client

import (
	"encoding/json"
	"log"
)

//message type
const TypeMessage = "gchat-message"
const TypeRoom = "gchat-room"
const TypeLeaveRoom = "gchat-leave-room"
const TypeJoinChannel = "gchat-join-channel"
const TypeLeaveChannel = "gchat-leave-channel"
const TypeCallOffer = "gchat-offer"
const TypeCallAnswer = "gchat-answer"


type Message struct {
	MessageType string `json:"messageType"`
	Message string `json:"message"`
	RoomName string `json:"room"`
	//Sender chan *websocket.Conn
}

func (message *Message) OnEncodeMessage() []byte {
	json,err := json.Marshal(&message)
	if err != nil {
		log.Fatal(err)
	}
	return json
}
