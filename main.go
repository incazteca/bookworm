package main

import (
	"encoding/json"
	"net/http"
)

const _10MB_LIMIT int64 = 1000 * 10000

type successRes struct {
	fileContent    string         `json:"file_content"`
	totalWordCount int            `json:"total_word_count"`
	wordCounts     map[string]int `json:"word_counts"`
}

type failRes struct {
	error string `json:"error"`
}

// Look into using multipart/form-data

func fileUploadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "POST" {
		if r.ContentLength > _10MB_LIMIT {
			w.WriteHeader(http.StatusRequestEntityTooLarge)
			failure := failRes{"File Upload is limited to 10MB. Please submit a smaller file"}
			json.NewEncoder(w).Encode(failure)
		}
	} else {
		http.Error(w, "The requested page is not available", http.StatusNotFound)
	}
}

func main() {
	http.HandleFunc("/file/upload", fileUploadHandler)
	http.ListenAndServe(":8080", nil)
}
