package main

import (
	"bitbucket.org/incazteca/bookworm/bookworm"
	"encoding/json"
	"net/http"
)

const _10MB_LIMIT int64 = (1000 * 10000) + 1

type SuccessRes struct {
	FileText       string         `json:"file_text"`
	TotalWordCount int            `json:"total_word_count"`
	WordCounts     map[string]int `json:"word_counts"`
}

type FailRes struct {
	ErrorMsg string `json:"error"`
}

// TODO: Clean this up, It's a disgrace...

func fileUploadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "POST" {
		if r.ContentLength > _10MB_LIMIT {
			w.WriteHeader(http.StatusRequestEntityTooLarge)
			response := FailRes{"File Upload is limited to 10MB. Please submit a smaller file"}
			json.NewEncoder(w).Encode(response)
		} else {
			r.ParseMultipartForm(_10MB_LIMIT)
			file, _, err := r.FormFile("file")
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				failure := FailRes{err.Error()}
				json.NewEncoder(w).Encode(failure)
				return
			}
			defer file.Close()

			output, err := bookworm.Parse(file, r.FormValue("filter"))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				failure := FailRes{err.Error()}
				json.NewEncoder(w).Encode(failure)
				return
			}
			w.WriteHeader(http.StatusOK)
			response := SuccessRes{output.Body, output.TotalWordCount, output.WordCounts}
			json.NewEncoder(w).Encode(response)
			return
		}

	} else {
		failure := FailRes{"Page not found"}
		json.NewEncoder(w).Encode(failure)
	}
}

func main() {
	http.HandleFunc("/file/upload", fileUploadHandler)
	http.ListenAndServe(":8080", nil)
}
