package main

import(
	"net/http"
	"fmt"
)
func uploadFile(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Uploading file")
	//parse input
	r.ParseMultipartForm(10 << 20) //20Mb
	//retrieve file from posted form-adta
	//write temp file onto server
	//retun code
	file, handler, err := r.FormFile("MyFile")
	if err != nil	{
		fmt.Println("Error retrieving file from form-data")
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Printf("Upload File: %+v\n",handler.Filename)
	fmt.Printf("File size: %+v\n",handler.Size)
	fmt.Printf("MIME header: %+v\n",handler.Header)
}
func SetupRoutes()	{
	http.HandleFunc("/upload", uploadFile)
	http.ListenAndServe(":8080", nil)
}
func main() {
	fmt.Println("File upload example")
	SetupRoutes()
}


