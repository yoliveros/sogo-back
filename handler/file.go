package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"sogo-back/db"
)

const TEMP_FOLDER = "temp-files"

type File struct {
	ID       string      `json:"id"`
	ParentID interface{} `json:"parent_id"`
	Name     string      `json:"name"`
	IsFile   bool        `json:"is_file"`
}

func (h *File) GetFiles(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("CALL sp_get_files()")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	defer rows.Close()

	files := []File{}

	for rows.Next() {
		file := File{}
		err = rows.Scan(&file.ID, &file.ParentID, &file.Name, &file.IsFile)
		if file.ParentID != nil {
			file.ParentID = string(file.ParentID.([]uint8))
		}

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		files = append(files, file)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(files)
}

func (h *File) createTempFolder() {
	_, err := os.Stat(TEMP_FOLDER)
	if err != nil {
		os.Mkdir(TEMP_FOLDER, 0755)
	}
}

func (h *File) Upload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	var parent_id interface{}
	file, handle, err := r.FormFile("file")
	parent_id = r.FormValue("parent_id")

	if err != nil {
		log.Println("Error Retrieving the File: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	defer file.Close()

	h.createTempFolder()

	// Create a temporary file within our temp-images directory
	tempFile, err := os.CreateTemp(TEMP_FOLDER, "upload-.*"+handle.Filename)
	if err != nil {
		log.Println("Temp file creation error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	defer tempFile.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Println("File reading error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	_, err = tempFile.Write(fileBytes)
	if err != nil {
		log.Println("File writing error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if parent_id == "" {
		parent_id = nil
	}

	_, err = db.DB.Exec("CALL sp_insert_files(?, ?, ?)", parent_id, handle.Filename, true)
	if err != nil {
		log.Println("File writing error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File uploaded successfully"))
}

func (h *File) CreateFolder(w http.ResponseWriter, r *http.Request) {
	file := File{}
	body, err := io.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	err = json.Unmarshal(body, &file)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	h.createTempFolder()

	_, err = os.Stat(TEMP_FOLDER + "/" + file.Name)
	if err != nil {
		os.Mkdir(TEMP_FOLDER+"/"+file.Name, 0755)
	}

	_, err = db.DB.Exec("CALL sp_insert_files(?, ?, ?)", file.ParentID, file.Name, file.IsFile)
	if err != nil {
		log.Println("File writing error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Folder created successfully"))
}
