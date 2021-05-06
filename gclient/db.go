package client

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

//structure connection to database
type DbConnect struct {
	//Db is instance for connection to database
	DB *sql.DB
}

//structure for keep chat in chat table
//
//type ChatData struct {
//	Message string
//	Image string
//	Sticker string
//	Video string
//	//UserId => senderId
//	UserId string
//	itemType string
//	//room name
//	Room string
//	//time is time that client send message to server
//	Timestamp time.Time
//	//message ID
//	MessageId string
//	UserName string
//
//}

//structure user list
type UserModel struct {
	Room string
	Name string
}

// set parameter that use for connection to database server host -> ..., port -> ...
func onNewConnect(db *sql.DB) *DbConnect {
	return &DbConnect{
		DB: db,
	}
}

//create table if table room is null and table chat is null
//will new create but alread not create
func (db *DbConnect) onCreateTable(){
	if _ , err := db.DB.Exec("create table if not exists room(id int not null auto_increment primary key , room varchar(60) , userName varchar(20))"); err != nil{
		log.Fatal("create room table error :",err)
	}
	//
	if _,err := db.DB.Exec("create table if not exists chat(id int null null auto_increment primary key , userId varchar(255) , room varchar(60) , itemType varchar(20) ,message varchar(255) , image varchar(255) , sticker varchar(255) , video varchar(255), timestamp time , messageId varchar(50))"); err != nil{
		log.Fatal("create table chat error :",err)
	}
}

//if there are clients register room will put room name in database
//in table name room
//
//name is room name that put
func (db *DbConnect) onKeepRoom(name string , userName string) error{
	//defind error type
	var error error

	//check room and name already in table
	rows , error := db.DB.Query("select userName , room from room where userName =?",userName)
	if error != nil{
		log.Println(error)
		error = error
	}

	defer rows.Close()

	var checkAlreadyRoomName Message
	//
	for rows.Next() {
		rows.Scan(&checkAlreadyRoomName.UserName , &checkAlreadyRoomName.RoomName)
	}

	if checkAlreadyRoomName.UserName == userName {
		//
		log.Println("already this user name in database :",checkAlreadyRoomName.RoomName)
	}else {
		//"insert into room (room , userName) values(?,?)",name,userName
		_ , err := db.DB.Exec("insert into room (room , userName) values(?,?)",name,userName)
		if err != nil{
			log.Println("insert room in table toom error :",err)
			error = err
		}
		//
	}

	return error
}

//if user change name update name in table room
//name is room name
//and user name
func (db *DbConnect) onUpdateRoom(name string,userName string) error{
	_,err := db.DB.Exec("update room set userName = ? where room = ?",userName,name)
	if err != nil{
		log.Println("update user name in room table error :",err)
	}
	//
	return err
}

//if clients remove chat room will remove room name from database
//table name room
//
//name is room name that will remove out
func (db *DbConnect) onUnRoom(name string) error{
	_,err := db.DB.Exec("delete from room where room =?",name)
	if err != nil{
		log.Println(err)
	}
	//
	return err
}

//get user that live in room
//room name -> get listData
func (db *DbConnect) getUserInRoom(name string) ([]UserModel, error) {
	result , err := db.DB.Query("select room , userName from room where room =?",name)
	//
	defer result.Close()
	//
	if err != nil{
		log.Println("get user that live in room error :",err)
	}
	orderUserList := make([]UserModel,50)
	for result.Next() {
		var user UserModel
		result.Scan(&user.Room , &user.Name)

		orderUserList = append(orderUserList , user)
	}
	//
	return orderUserList , nil
}

//if clients send message to room server will keep history chat between client with client
// such as : message image or video
//
//chat is chatData structure for save
//
//item type is type of message -> text video
// and message is -> data chat that client send
func (db *DbConnect) onKeepMessage(chat *Message , url string) error{
	//create sql string
	var sqlString = "insert into chat"
	var item = ""
	//check item type it is ?
	switch chat.ItemType {
	case ItemTypeMessage:
		sqlString += "(userId ,room , itemType , message , timestamp , messageId) values(?,?,?,?,?,?)"
		item += url
		break
	case ItemTypeImage:
		sqlString += "(userId ,room , itemType , image , timestamp , messageId) values(?,?,?,?,?,?)"
		item += url
		break
	case ItemTypeSticker:
		sqlString += "(userId ,room , itemType , sticker , timestamp , messageId) values(?,?,?,?,?,?)"
		item += url
		break
	case ItemTypeVideo:
		sqlString += "(userId ,room , itemType , video , timestamp , messageId) values(?,?,?,?,?,?)"
		item += url
		break
	default:
		break
	}
	_ , err := db.DB.Exec(sqlString , chat.UserId,chat.RoomName , chat.ItemType , item ,time.Now().Format("2006-01-02 15:04:05"),chat.MessageId)
	if err != nil{
		log.Fatal(err)
	}
	//
	return err
}

//get all message that client send between client
//name -> room name use for search message
func (db *DbConnect) getChat(name string) ([]*Message, error){
	result , err := db.DB.Query("select userId ,room ,itemType ,message,image ,sticker ,video ,timestamp ,messageId from chat where room =?",name)
	if err != nil{
		log.Println("get list message error :",err)
	}
	defer result.Close()
	//
	var orderChat []*Message
	//
	for result.Next() {
		var message Message
		result.Scan(&message.UserId,&message.RoomName,&message.ItemType,&message.Message , &message.Image , &message.Sticker , &message.Video ,&message.Timestamp ,&message.MessageId)
		//log.Println(message)
		orderChat = append(orderChat,&message)
	}
	//for i := range orderChat{
	//	log.Println(&orderChat[i])
	//}
	//
	return orderChat, err
}

//if client remove message server will remove message in database
//that table name is chat
//by use itemType and message
func (db *DbConnect) onUnMessage(messageId string) error{
	//create sql string
	var sqlString = "delete from chat where messageId =?"
	//check item type it is ?
	//switch itemType {
	//case ItemTypeMessage:
	//	sqlString += "where message = $1"
	//	break
	//case ItemTypeImage:
	//	sqlString += "where image = $1"
	//	break
	//case ItemTypeSticker:
	//	sqlString += "where sticker = $1"
	//	break
	//case ItemTypeVideo:
	//	sqlString += "where video = $1"
	//	break
	//default:
	//	break
	//}
	_ , err := db.DB.Exec(sqlString, messageId)
	if err != nil{
		log.Println(err)
	}
	//
	return err
}

//create connection to database mysql and return insatnce sql.DB
//for use manament database
func OnConnectionDatabase() (*DbConnect,error){
	//stringConn := mConn.User+":"+mConn.Password+"@tcp"+"("+mConn.Host+""+mConn.Port+")"+"/"+mConn.Dbname
	db , err := sql.Open("mysql","redev:redev@tcp(localhost:3306)/exvention_face_db")
	if err != nil{
		//log.Println("connection to database error :",err)
	}
	//keep connection parameter
	mConn := onNewConnect(db)
	//create table
	mConn.onCreateTable()
	//
	return mConn ,err
}