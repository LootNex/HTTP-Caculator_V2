package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Account struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (db App) Register(w http.ResponseWriter, r *http.Request) {

	account := new(Account)
	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		http.Error(w, "login or password is incorect", http.StatusBadRequest)
	}

	var exist bool
	fmt.Println(*account)
	err = db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE login = ?)", account.Login).Scan(&exist)
	if err != nil {
		http.Error(w, fmt.Sprintf("cannot check login and password %v", err), http.StatusInternalServerError)
		return
	}
	log.Println(exist)
	if exist {
		http.Error(w, "this user already exist", http.StatusBadRequest)
		return
	}

	uuid := uuid.New().String()

	_, err = db.DB.Exec("INSERT INTO users(login, password, user_id) VALUES (?, ?, ?)", account.Login, account.Password, uuid)
	if err != nil {
		http.Error(w, "problems with Database", http.StatusInternalServerError)
	}

}

func (db App) SingIn(w http.ResponseWriter, r *http.Request) {

	account := new(Account)
	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		http.Error(w, "your login or password is incorrect", http.StatusBadRequest)
	}
	var uuid string
	err = db.DB.QueryRow("SELECT user_id FROM users WHERE login = ? AND password = ?", account.Login, account.Password).Scan(&uuid)
	if err != nil {
		http.Error(w, fmt.Sprintf("cannot check login and password %v", err), http.StatusInternalServerError)
		return
	}
	if uuid == "" {
		http.Error(w, "you should register", http.StatusInternalServerError)
		return
	}

	now := time.Now()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": uuid,
		"nbf":     now.Unix(),
		"exp":     now.Add(10 * time.Minute).Unix(),
	})

	tokenString, err := token.SignedString([]byte("super_secret_signature"))
	if err != nil {
		http.Error(w, "problems with jwt"+err.Error(), http.StatusInternalServerError)
	}

	var response = struct {
		Token   string `json:"token"`
		Message string `json:"message"`
		Status  string `json:"status"`
	}{
		Token:   tokenString,
		Message: "Authetication succesful",
		Status:  "OK",
	}
	jsonData, err := json.MarshalIndent(&response, "", "  ")
	if err != nil {
		http.Error(w, "problems to send jwt"+err.Error(), http.StatusInternalServerError)
	}

	_, err = w.Write(jsonData)
	if err != nil {
		http.Error(w, "problems to send jwt"+err.Error(), http.StatusInternalServerError)
	}

	_, err = db.DB.Exec("UPDATE users SET user_id = ? WHERE login = ? AND password = ?", uuid, account.Login, account.Password)
	if err != nil {
		http.Error(w, fmt.Sprintf("cannot update users %v", err), http.StatusInternalServerError)
	}

}

func (db App) Compare(userId string) (bool, error) {
	var exist bool
	err := db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE user_id = ?)", userId).Scan(&exist)

	if err != nil {
		return false, err
	}

	if !exist {
		return false, nil
	}

	return true, nil
}
