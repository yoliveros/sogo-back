package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/yoliveros/sogo-back/middleware"

	_ "github.com/go-sql-driver/mysql"
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
