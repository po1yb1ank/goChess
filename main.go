package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"html/template"
)
func serveTemplate(w http.ResponseWriter, r *http.Request, typeOf string){
	var tmpl *template.Template
	switch typeOf {
	case "login":
		tmpl = template.Must(template.ParseFiles(path.Join("templates", "index.html"), path.Join("templates", "login.html")))
	case "main":
		tmpl = template.Must(template.ParseFiles(path.Join("templates", "index.html"), path.Join("templates", "main.html")))
	case "register":
		tmpl = template.Must(template.ParseFiles(path.Join("templates", "index.html"), path.Join("templates", "register.html")))
	}
	if err := tmpl.ExecuteTemplate(w, "main", nil); err != nil{
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}
func register(w http.ResponseWriter, r *http.Request){
	serveTemplate(w,r,"register")

}
func uploadFile(w http.ResponseWriter, r *http.Request){
	serveTemplate(w, r, "login")
	fmt.Fprintf(w, "Uploading file\n")
	//parse input
	r.ParseMultipartForm(10 << 20) //20Mb
	//retrieve file from posted form-data
	//retun code
	file, handler, err := r.FormFile("File")
	if err != nil	{
		fmt.Println("Error retrieving file from form-data")
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Printf("Upload File: %+v\n",handler.Filename)
	fmt.Printf("File size: %+v\n",handler.Size)
	fmt.Printf("MIME header: %+v\n",handler.Header)

	//write temp file onto server
	tempFile, err := ioutil.TempFile("temp-images", "upload-*.png")
	if err != nil{
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
func mainPage (w http.ResponseWriter, r *http.Request){
	serveTemplate(w, r, "main")
}
func logOut (w http.ResponseWriter, r *http.Request){
	//this is logout func
	//logout execution
	//redirect to login
	serveTemplate(w, r, "login")
}
func SetupRoutes()	{
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/",fs))
	http.HandleFunc("/", uploadFile)//upload file
	http.HandleFunc("/logout", logOut)
	http.HandleFunc("/register", register)
	http.HandleFunc("/main", mainPage)
	//http.HandleFunc("/upload", uploadFile)
	http.ListenAndServe(":8080", nil)
}
func main() {
	fmt.Println("File upload example")
	SetupRoutes()
}


