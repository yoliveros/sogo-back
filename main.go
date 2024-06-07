package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/yoliveros/sogo-back/middleware"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	router := http.NewServeMux()
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/sogo")
	if err != nil {
		log.Fatal(err)
	}

	users, err := db.Query("CALL sp_get_users()")
	if err != nil {
		log.Fatal(err)
	}

	for users.Next() {
		var id string
		var name string
		var last_name string
		var email string
		var password string
		var created_at string
		var updated_at string
		if err := users.Scan(&id, &name, &last_name, &email, &password, &created_at, &updated_at); err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, name, email, password, created_at, updated_at)
	}

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
