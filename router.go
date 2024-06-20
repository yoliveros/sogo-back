package main

import (
	"net/http"

	"sogo-back/handler"
)

func loadRouters(roter *http.ServeMux) {
	file_handler := &handler.File{}

	roter.HandleFunc("GET /files", file_handler.GetFiles)
	roter.HandleFunc("POST /upload", file_handler.Upload)
	roter.HandleFunc("POST /folder", file_handler.CreateFolder)
}
