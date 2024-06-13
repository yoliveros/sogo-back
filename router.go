package main

import (
	"net/http"

	"sogo-back/handler"
)

func loadRouters(roter *http.ServeMux) {
	file_handler := &handler.File{}

	roter.HandleFunc("POST /upload", file_handler.Upload)
}
