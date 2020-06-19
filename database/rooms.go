package database

import (
	"fmt"
	"github.com/gorilla/websocket"
)

type roomStr struct {
	player1 *websocket.Conn
	player2 *websocket.Conn
	col1 string
	col2 string
	roomId string
	available bool
}
var allRooms = make(map[string]*roomStr)
var cRoom string

func InitRooms()  {
}
func NewRoom(room string){
	fmt.Println("entered room creation")
	cRoom = room
	fmt.Println("croom set")
	if _, ok := allRooms[room]; ok {
		fmt.Println("room already exists")
		return
	}
	fmt.Println("check pass")
	var thisRoom roomStr
	thisRoom.player1 = nil
	thisRoom.player2 = nil
	thisRoom.col1 = "w"
	thisRoom.col2 = "b"
	thisRoom.roomId = room
	thisRoom.available = true
	allRooms[room] = &thisRoom
	fmt.Println("created room")

}
func AddRoomPlayer(ws *websocket.Conn)  {
	fmt.Println("for",&ws)
	if thisRoom, ok := allRooms[cRoom];ok{
		if thisRoom.player1 == nil{
			thisRoom.player1 = ws
			allRooms[cRoom] = thisRoom
			return
		}
		if thisRoom.player1 != nil && thisRoom.player2 == nil{
			thisRoom.player2 = ws
			thisRoom.available = false
			allRooms[cRoom] = thisRoom
			return
		}
	}
}
func IfRoomAvailable() bool  {
	fmt.Println(cRoom)
	if thisRoom, ok := allRooms[cRoom]; ok{
		fmt.Println("status", thisRoom.available)
		return thisRoom.available
	}
	fmt.Println("no check")
	return false
}