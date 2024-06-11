package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

type Handler struct{}

type User struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	results, err := db.Query("CALL sp_get_users()")
	if err != nil {
		log.Fatal(err)
	}

	users := []User{}

	for results.Next() {
		var user User
		if err := results.Scan(
			&user.ID,
			&user.Name,
			&user.LastName,
			&user.Username,
			&user.Password,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			log.Fatal(err)
		}

		users = append(users, user)
	}

	json.NewEncoder(w).Encode(users)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	result, err := db.Query("CALL sp_get_user_by_username(?)", username)
	if err != nil {
		log.Fatal(err)
	}

	if !result.Next() {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("User not found"))
		return
	}

	user := User{}

	result.Scan(&user.Username, user.Password)

	if user.Password != password {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Wrong password"))
		return
	}

	w.Write([]byte("You are logged in!"))
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	log.Println("Hello, World!")
	w.Write([]byte("Hello, World!"))
}
