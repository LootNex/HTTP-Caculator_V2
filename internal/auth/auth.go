package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Account struct {
	login    string
	password string
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
		http.Error(w, "your login or password is incorrect", http.StatusBadRequest)
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
