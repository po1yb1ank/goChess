package database
var allRooms = make(map[string]*room)
type room struct {
	player1 string
	player2 string
	col1 string
	col2 string
	roomId string
	available bool
}