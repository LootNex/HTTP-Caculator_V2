package auth

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Account struct {
	login    string
	password string
}

type App struct{
	DB *sql.DB
}

func NewApp(db *sql.DB) *App{
	return &App{
		DB: db,
	}
}

func (a App) Register(w http.ResponseWriter, r *http.Request) {

	account := new(Account)
	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil{
		http.Error(w, "login or password is incorect", http.StatusBadRequest)
	}

	_, err = a.DB.Exec("INSERT INTO users(login, password) VALUES ($1, $2)", account.login, account.password)
	if err != nil{
		http.Error(w, "problems with Database", http.StatusInternalServerError)
	}

}

func (db App) SingIn(w http.ResponseWriter, r *http.Request) {

	account := new(Account)
	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil{
		fmt.Errorf("your login or password is incorrect")
	}

	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name": account.login,
		"nbf":  now.Unix(),
		"exp":  now.Add(5 * time.Minute).Unix(),
	})


	tokenString, err := token.SignedString("super_secret_signature")
	if err != nil{
		http.Error(w, "problems with jwt", http.StatusInternalServerError)
	}

	_, err = w.Write([]byte(tokenString))
	if err != nil{
		http.Error(w, "problems with jwt", http.StatusInternalServerError)
	}

}
