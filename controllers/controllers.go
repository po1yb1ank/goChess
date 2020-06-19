package controllers

import (
	"fmt"
	"github.com/gorilla/websocket"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strings"
	"sync"
	"uploadServer/database"
	"github.com/gorilla/sessions"
)
//pics
var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)
var doOnce sync.Once
var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
}
func Init(w http.ResponseWriter, r *http.Request) {
	/* init page, there we will setup database
	 */
	session, _ := store.Get(r, "cookie-name")
	if database.DataBaseStatus() == true {
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth  {
			http.Redirect(w,r,"/login",302)
		} else {
			http.Redirect(w,r,"/main",302)
		}
	}
	//problem: if user already logged
	if database.DataBaseStatus() == false {
		database.InitRooms()
		database.SetDataBase()
		ServeTemplate(w, r, "login")
	}
}
func ServeTemplate(w http.ResponseWriter, r *http.Request, typeOf string) {
	var tmpl *template.Template
	switch typeOf {
	case "login":
		tmpl = template.Must(template.ParseFiles(path.Join("templates", "index.html"), path.Join("templates", "login.html")))
	case "main":
		tmpl = template.Must(template.ParseFiles(path.Join("templates", "index.html"), path.Join("templates", "main.html")))
	case "room":
		tmpl = template.Must(template.ParseFiles(path.Join("templates", "index.html"), path.Join("templates", "room.html")))
	case "register":
		tmpl = template.Must(template.ParseFiles(path.Join("templates", "index.html"), path.Join("templates", "register.html")))
	}
	if err := tmpl.ExecuteTemplate(w, "main", nil); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}
func Register(w http.ResponseWriter, r *http.Request) {
	//if user already logged in, response reject
	var l, p string
	session, _ := store.Get(r, "cookie-name")
	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		if r.Method == "GET" {
			ServeTemplate(w, r, "register")
		}
		if r.Method == "POST" {
			r.ParseForm()
			l = strings.Join(r.Form["login"], "")
			p = strings.Join(r.Form["password"], "")
			fmt.Println("login: ",l)
			fmt.Println("pass: ",p)
			if l != "" && p != "" {
				database.SetUser(l, p)
				/*if database.IfLogged() == true {
					http.Redirect(w, r, "https://127.0.0.1:8080/redirect", 301)
				}*/
				session.Values["authenticated"] = true
				session.Save(r, w)
				http.Redirect(w, r, "/main", 303)
			}
		}
	}else{
		http.Redirect(w, r, "/main", 303)
	}

}
func Redirect(w http.ResponseWriter, r *http.Request) {

	http.Redirect(w, r, "/main", 303)
}
func UploadFile(w http.ResponseWriter, r *http.Request) {
	ServeTemplate(w, r, "main")
	fmt.Fprintf(w, "Uploading file\n")
	//parse input
	r.ParseMultipartForm(10 << 20) //20Mb
	//retrieve file from posted form-data
	//retun code
	file, handler, err := r.FormFile("File")
	if err != nil {
		fmt.Println("Error retrieving file from form-data")
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Printf("Upload File: %+v\n", handler.Filename)
	fmt.Printf("File size: %+v\n", handler.Size)
	fmt.Printf("MIME header: %+v\n", handler.Header)

	//write temp file onto server
	tempFile, err := ioutil.TempFile("temp-images", "upload-*.file")
	if err != nil {
		fmt.Println("Uploading error")
		fmt.Println(err)
		return
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	tempFile.Write(fileBytes)

	fmt.Fprintf(w, "Successfully uploaded file\n")

}
func Reader(ws *websocket.Conn)  {
	for  {
		messageType, p, err := ws.ReadMessage()
		if err != nil{
			//fmt.Println("error at reader", err)
			return
		}
		fmt.Println(string(p))
		if err := ws.WriteMessage(messageType, p); err != nil{
			fmt.Println(err)
			return
		}
	}
}
func MainPage(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/login", 303)
	} else {
		if r.Method == "GET" {
			ServeTemplate(w, r, "main")
		}
		if r.Method == "POST" {
			r.ParseForm()
			room := strings.Join(r.Form["room"], "")
			fmt.Println(room)
			database.NewRoom(room)
			if database.IfRoomAvailable() == true{
				http.Redirect(w, r, "/room", 303)
			}else{
				fmt.Println("room is full")
				http.Redirect(w, r, "/main", 303)
			}
		}
	}
}
func Login(w http.ResponseWriter, r *http.Request)  {
	session, _ := store.Get(r, "cookie-name")
	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		if r.Method == "GET" {
			ServeTemplate(w, r, "login")
		}
		if r.Method == "POST" {
			r.ParseForm()
			l := strings.Join(r.Form["login"], "")
			p := strings.Join(r.Form["password"], "")
			fmt.Println("login: ",l)
			fmt.Println("pass: ",p)
			if l != "" && p != "" {
				if k,v := database.SeekDB(l); k!= "" && v == p && l == k{
					fmt.Println("found!")
					session.Values["authenticated"] = true
					session.Save(r, w)
					http.Redirect(w,r,"/main", 303)
				}else{
					http.Redirect(w,r,"/login", 303)
					//ServeTemplate(w, r, "login")
				}
			}
		}
	}else {
		http.Redirect(w,r,"/main", 301)
	}
}
func LogOut(w http.ResponseWriter, r *http.Request) {
	//this is logout func
	//logout execution
	//redirect to login
	database.ClearUser()
	session, _ := store.Get(r, "cookie-name")
	session.Values["authenticated"] = false
	session.Save(r, w)
	http.Redirect(w, r, "/login",302)
	ServeTemplate(w, r, "login")
}
func WS (w http.ResponseWriter, r *http.Request){
	upgrader.CheckOrigin = func(r *http.Request) bool {return true}
	ws, err := upgrader.Upgrade(w, r,nil )
	if err != nil{
		fmt.Println("error:", err)
		return
	}
	fmt.Println("Connected")
	database.AddRoomPlayer(ws)
	Reader(ws)
}
func Room (w http.ResponseWriter, r *http.Request){
	ServeTemplate(w, r, "room")
}