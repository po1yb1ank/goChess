package controllers

import (
	"fmt"
	"strings"

	//"github.com/patrickmn/go-cache"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	//"time"
	"uploadServer/database"
	//"github.com/alexedwards/scs"
)
func SetDataBase (){

}
func Init (w http.ResponseWriter, r* http.Request){
	/* init page, there we will setup database
	 */
	if database.DataBaseStatus() == true{
		if database.IfLogged() == true{
			MainPage(w, r)
		} else {
			ServeTemplate(w, r, "login")
		}
	}
	//problem: if user already logged
	if database.DataBaseStatus() == false{
		database.SetDataBase()
		ServeTemplate(w, r, "login")
	}
}
func ServeTemplate(w http.ResponseWriter, r *http.Request, typeOf string){
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
func Register(w http.ResponseWriter, r *http.Request){
	ServeTemplate(w,r,"register")
	//if user already logged in, response reject
	if database.IfLogged() == true{
		ServeTemplate(w, r, "main")
	} else {
		//parse forms
		r.ParseForm()
		l := strings.Join(r.Form["login"], "")
		p := strings.Join(r.Form["password"], "")

		if l != "" && p != ""{
			database.SetUser(l, p)
			if database.IfLogged() == true {
				http.Redirect(w, r, "/main", 307)
			}
		}
	}

}
func UploadFile(w http.ResponseWriter, r *http.Request){
	ServeTemplate(w, r, "main")
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
func MainPage (w http.ResponseWriter, r *http.Request){
	ServeTemplate(w, r, "main")
	if database.IfLogged(){
		fmt.Println(database.SeekDB("poly"))
	}
}
func LogOut (w http.ResponseWriter, r *http.Request){
	//this is logout func
	//logout execution
	//redirect to login
	ServeTemplate(w, r, "login")
}