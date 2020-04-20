package main

import(
	"io/ioutil"
	"net/http"
	"fmt"
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
	http.HandleFunc("/upload", uploadFile)
	http.ListenAndServe(":8080", nil)
}
func main() {
	fmt.Println("File upload example")
	SetupRoutes()
}


