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
			fmt.Println("ws1")
			thisRoom.player1 = ws
			allRooms[cRoom] = thisRoom
			thisRoom.player1.WriteMessage(1,[]byte(thisRoom.col1))
			return
		}
		if thisRoom.player1 != nil && thisRoom.player2 == nil{
			fmt.Println("ws2")
			thisRoom.player2 = ws
			thisRoom.available = false
			allRooms[cRoom] = thisRoom
			thisRoom.player2.WriteMessage(1,[]byte(thisRoom.col2))
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
func PosChange(mt int, pos []byte)  {
	if thisRoom, ok := allRooms[cRoom]; ok{
		thisRoom.player1.WriteMessage(mt, pos)
		thisRoom.player2.WriteMessage(mt, pos)
	}
}
func SendSides(mt int)  {
	if thisRoom, ok := allRooms[cRoom]; ok{
		fmt.Println([]byte(thisRoom.col1))
		thisRoom.player1.WriteMessage(mt, []byte(thisRoom.col1))
		thisRoom.player2.WriteMessage(mt, []byte(thisRoom.col2))
	}
}