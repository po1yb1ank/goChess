package database

type userData struct {
	login string
	password string
	logStatus bool
}
var user userData
func SetUser(l string, p string) {
	user.login = l
	user.password = p
	user.logStatus = true
	AddDB(l, p)
}
func SetCurrentUser(l string, p string)  {
	user.login = l
	user.password = p
	user.logStatus = true
}

func ClearUser()  {
	user.login = ""
	user.password = ""
	user.logStatus = false
}