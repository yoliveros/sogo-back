package main

import (
	"net/http"

	"github.com/yoliveros/sogo-back/handler"
)

func loadRouters(roter *http.ServeMux) {
	auth := &handler.Handler{}
	handler.InitDB()

	roter.HandleFunc("GET /auth/users", auth.GetUsers)
	roter.HandleFunc("POST /auth/login", auth.Login)
	roter.HandleFunc("POST /auth/register", auth.Register)
}
