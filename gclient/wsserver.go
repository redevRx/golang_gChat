package client

import "log"

type WsServer struct {
	clients map[*Client]bool
	Rooms map[*Room]bool
	register chan *Client
	unRegister chan *Client
	sendMessage chan []byte

}

func OnNewWsServer() *WsServer{
	return &WsServer{
		clients: make(map[*Client]bool),
		Rooms: make(map[*Room]bool),
		register: make(chan *Client),
		unRegister: make(chan *Client),
		sendMessage: make(chan []byte),
	}
}

//register and unregister
func (ws *WsServer) onRegister(client *Client){
	ws.clients[client] = true
}
func (ws *WsServer) onUnRegister(client *Client){
	if _, ok := ws.clients[client]; ok{
		log.Println("close connected :",client.conn.RemoteAddr())
		delete(ws.clients , client)
		client.conn.Close()
		log.Println("ws room :",len(ws.Rooms))
		log.Println("ws client :",len(ws.clients))
	}
}

func (ws *WsServer) OnRun(){
	for  {
		select {
		case client := <- ws.register:
			ws.onRegister(client)
			//log.Fatalln("new client")
			break
		case client := <- ws.unRegister:
			ws.onUnRegister(client)
			//log.Fatalln("client close")
			break
		case message := <- ws.sendMessage:
			for client := range ws.clients{
				select {
				case client.message <- message:
				default:
					close(client.message)
					delete(ws.clients , client)
				}
			}
		}
	}
}

//room
//find room name
func (ws *WsServer) onFindRoomName(name string) *Room{
	var myroom *Room
	for room := range ws.Rooms{
		if room.Name == name{
			myroom = room
			break
		}
	}
	return myroom
}

//create new room
func (ws *WsServer) onCreateNewRoom(name string) *Room{
	room := onNewRoom(name)
	go room.onRun()
	ws.Rooms[room] = 0 == 0

	return room
}

