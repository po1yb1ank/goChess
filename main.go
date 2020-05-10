package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path"
)
func uploadFile(w http.ResponseWriter, r *http.Request){
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
	http.HandleFunc("/", serveTemplate)//upload file
	http.HandleFunc("/upload", uploadFile)
	http.ListenAndServe(":8080", nil)
}
func serveTemplate(w http.ResponseWriter, r *http.Request){

	/*
	lp := filepath.Join("templates","layout.html")
	fp := filepath.Join("templates", filepath.Clean(r.URL.Path))
	tmpl, _ := template.ParseFiles(lp, fp)
	 */
	tmpl := template.Must(template.ParseFiles(path.Join("templates", "layout.html"), path.Join("templates", "index.html")))
	if err := tmpl.ExecuteTemplate(w, "layout", nil); err != nil{
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}

}
func main() {
	fmt.Println("File upload example")
	SetupRoutes()
}


