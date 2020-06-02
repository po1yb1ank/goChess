package database

type userData struct {
	login string
	password string
	logStatus bool
}
var user userData
func SetUser() {

}
func IfLogged() bool {
	return user.logStatus
}