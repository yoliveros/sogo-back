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

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	results, err := db.Query("CALL sp_get_users()")
	if err != nil {
		log.Fatal(err)
	}
	defer results.Close()

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
	login := Login{}

	json.NewDecoder(r.Body).Decode(&login)
	defer r.Body.Close()

	if login.Username == "" || login.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing username or password"))
		return
	}

	result, err := db.Query("CALL sp_get_user_by_username(?)", login.Username)
	if err != nil {
		log.Fatal(err)
	}
	defer result.Close()

	user := User{}

	if result.Next() {
		result.Scan(&user.Username, &user.Password)
	}

	if user.Username == "" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("User not found"))
		return
	}

	log.Println(user)
	log.Println(login)

	if user.Username != login.Username || user.Password != login.Password {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid credentials"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login successful"))
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	log.Println("Hello, World!")
	w.Write([]byte("Hello, World!"))
}
