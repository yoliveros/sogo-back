package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"sogo-back/db"
)

type File struct {
	ID       string `json:"id"`
	ParentID string `json:"parent_id"`
	Name     string `json:"name"`
}

func (h *File) GetFiles(w http.ResponseWriter, r *http.Request) {

	rows, err := db.DB.Query("CALL sp_get_files()")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	files := []File{}

	for rows.Next() {
		file := File{}
		err = rows.Scan(&file.ID, &file.ParentID, &file.Name)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		files = append(files, file)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(files)
}

func (h *File) Upload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	var algo interface{}
	algo = "hola"
	log.Println(algo)

	var parent_id interface{}

	file, handle, err := r.FormFile("file")
	parent_id = r.FormValue("parent_id")

	log.Println(parent_id)

	if err != nil {
		log.Println("Error Retrieving the File: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	_, err = os.Stat("temp-files")
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

	if parent_id == "" {
		parent_id = &parent_id
	}

	_, err = db.DB.Exec("CALL sp_insert_files(?, ?)", parent_id, handle.Filename)
	if err != nil {
		log.Println("File writing error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File uploaded successfully"))
}
