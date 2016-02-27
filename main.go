package main

import (
	"fmt"
	"net/http"
)

// Look into using multipart/form-data

type File struct {
	body string
}

func fileHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Wazzup")
}

func main() {
	http.HandleFunc("/file/upload", fileHandler)
	http.ListenAndServe(":8080", nil)
}
