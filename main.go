package main

import (
	"fmt"
	"log"
	"net/http"

	"sogo-back/db"
	"sogo-back/middleware"
)

func main() {
	router := http.NewServeMux()

	db.InitDB()
	defer db.DeinitDB()

	loadRouters(router)

	stack := middleware.CreateStack(
		middleware.Logging,
		middleware.AllowCors,
	)

	server := http.Server{
		Addr:    ":8080",
		Handler: stack(router),
	}

	fmt.Println("Server listening on port 8080")
	log.Fatal(server.ListenAndServe())
}
