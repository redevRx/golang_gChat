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
const TypeGetMessage = "gchat-get-all-message"

//
//defind item type
//item type is when client send message to server and server
//will check message it is ? such as: text image sticker or video
const ItemTypeMessage = "text"
//defind item type
//item type is when client send message to server and server
//will check message it is ? such as: text image sticker or video
const ItemTypeImage = "image"
//defind item type
//item type is when client send message to server and server
//will check message it is ? such as: text image sticker or video
const ItemTypeSticker = "sticker"
//defind item type
//item type is when client send message to server and server
//will check message it is ? such as: text image sticker or video
const ItemTypeVideo = "video"


type Message struct {
	RoomName string `json:"room"`
	//UserId => senderId
	UserId string `json:"userId"`
	UserName string `json:"userName"`
	MessageType string `json:"messageType"`
	MessageId string `json:"messageId"`
	ItemType string `json:"itemContentType"`
	Message string `json:"message"`
	//Sender chan *websocket.Conn
	//
	Image [] byte `json:"image"`
	Sticker string `json:"sticker"`
	Video [] byte  `json:"video"`
	//room name
	//time is time that client send message to server
	Timestamp string
	//message ID
}

func (message *Message) OnEncodeMessage() []byte {
	json,err := json.Marshal(&message)
	if err != nil {
		log.Fatal(err)
	}
	return json
}
