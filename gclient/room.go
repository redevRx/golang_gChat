package client

type Room struct {
	Name string
	Clients map[*Client]bool
	Register chan *Client
	UnRegister chan *Client
	OnMessage chan *Message
}

func onNewRoom(name string) *Room{
	return &Room{
		Name: name,
		Clients: make(map[*Client]bool),
		Register: make(chan *Client),
		UnRegister: make(chan *Client),
		OnMessage: make(chan *Message),
	}
}

//register room
func (r *Room) onRoomRegister(client *Client){
	r.Clients[client] = 0 == 0
}
//unRegister room
func (r *Room) onRoomUnRegister(client *Client){
	if _,ok := r.Clients[client]; ok{
		//there is this room
		delete(r.Clients , client)
		client.conn.Close()
	}
}


//room loop check event for client
func (r *Room) onRun(){
	for {
		select {
		case room := <- r.Register:
			//register
			r.onRoomRegister(room)
			break
		case room := <- r.UnRegister:
			//un register room
			r.onRoomUnRegister(room)
			break
		case message := <- r.OnMessage:
			for client := range r.Clients{
				select {
				case client.message <- message.OnEncodeMessage():
				default:
					close(client.message)
					delete(r.Clients,client)
					break
				}
			}
		}
	}
}