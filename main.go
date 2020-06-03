package main

import (
	"fmt"
	"github.com/alexedwards/scs"
	"net/http"
	"time"
	"uploadServer/controllers"
	//"uploadServer/database"
)
var sessionManager *scs.SessionManager
func SetupRoutes()	{
	sessionManager = scs.New()
	sessionManager.Lifetime = 24 * time.Hour
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/",fs))
	http.HandleFunc("/", controllers.Init)
	http.HandleFunc("/upload", controllers.UploadFile)//upload file
	http.HandleFunc("/login", controllers.Login)
	http.HandleFunc("/logout", controllers.LogOut)
	http.HandleFunc("/register", controllers.Register)
	http.HandleFunc("/main", controllers.MainPage)
	http.HandleFunc("/redirect", controllers.Redirect)
	//http.HandleFunc("/upload", uploadFile)
	http.ListenAndServeTLS(":8080","cert.pem","key.pem", nil)
	//http.ListenAndServe(":8080", nil)
}
func main() {
	fmt.Println("File upload example")
	SetupRoutes()
}


