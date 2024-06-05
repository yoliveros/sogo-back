package handler

import (
	"log"
	"net/http"
)

type Handler struct{}

func (h *Handler) Test(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	log.Println("recived")
	w.Write([]byte("Test reached"))
}
