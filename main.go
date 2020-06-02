package main

import (
	"fmt"
	"net/http"
	"uploadServer/controllers"
	//"uploadServer/database"
)

func SetupRoutes()	{
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/",fs))
	http.HandleFunc("/", controllers.Init)
	http.HandleFunc("/upload", controllers.UploadFile)//upload file
	http.HandleFunc("/logout", controllers.LogOut)
	http.HandleFunc("/register", controllers.Register)
	http.HandleFunc("/main", controllers.MainPage)
	//http.HandleFunc("/upload", uploadFile)
	http.ListenAndServe(":8080", nil)
}
func main() {
	fmt.Println("File upload example")
	SetupRoutes()
}


