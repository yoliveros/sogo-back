package main

import (
	"net/http"

	"github.com/yoliveros/sogo-back/handler"
)

func loadRouters(roter *http.ServeMux) {
	auth := &handler.Handler{}

	roter.HandleFunc("GET /auth/test", auth.Test)
}
