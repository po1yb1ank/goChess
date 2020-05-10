package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"html/template"
)
func serveTemplate(w http.ResponseWriter, r *http.Request){
	tmpl := template.Must(template.ParseFiles(path.Join("templates", "layout.html"), path.Join("templates", "index.html")))
	if err := tmpl.ExecuteTemplate(w, "layout", nil); err != nil{
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}
func uploadFile(w http.ResponseWriter, r *http.Request){
	serveTemplate(w, r)
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
func SetupRoutes()	{
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/",fs))
	http.HandleFunc("/", uploadFile)//upload file
	//http.HandleFunc("/upload", uploadFile)
	http.ListenAndServe(":8080", nil)
}
func main() {
	fmt.Println("File upload example")
	SetupRoutes()
}


