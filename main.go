package main

import(
	"net/http"
	"fmt"
)
func uploadFile(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Uploading file")
}
func SetupRoutes()	{
	http.HandleFunc("/upload", uploadFile)
	http.ListenAndServe(":8080", nil)
}
func main() {
	fmt.Println("File upload example")
}


