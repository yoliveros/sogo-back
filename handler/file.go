package handler

import (
	"io"
	"log"
	"net/http"
	"os"
)

type File struct{}

func (h *File) Upload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	file, handle, err := r.FormFile("file")

	if err != nil {
		log.Println("Error Retrieving the File: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	stat, err := os.Stat("temp-images")
	log.Println("stat: ", stat, err)
	if err != nil {
		os.Mkdir("temp-files", 0755)
	}

	// Create a temporary file within our temp-images directory
	tempFile, err := os.CreateTemp("temp-files", "upload-.*"+handle.Filename)
	if err != nil {
		log.Println("Temp file creation error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer tempFile.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Println("File reading error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = tempFile.Write(fileBytes)
	if err != nil {
		log.Println("File writing error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File uploaded successfully"))
}
