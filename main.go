package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/yoliveros/sogo-back/middleware"
)

func main() {
	router := http.NewServeMux()

	loadRouters(router)

	stack := middleware.CreateStack(
		middleware.Logging,
		middleware.AllowCors,
		// middleware.IsAuthenticated,
	)

	server := http.Server{
		Addr:    ":8080",
		Handler: stack(router),
	}

	fmt.Println("Server listening on port 8080")
	log.Fatal(server.ListenAndServe())
}

func rootHandler(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("Hello, World!"))
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		w.Write([]byte("Hello, API!"))
		return
	}

	w.Write([]byte("Hello, " + id + "!"))
}
