package orkestrator

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	auth "github.com/LootNex/HTTP-Caculator_V2/internal/auth"
	calculator "github.com/LootNex/HTTP-Caculator_V2/pkg"
	config "github.com/LootNex/HTTP-Caculator_V2/pkg/config"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Expression struct {
	Id     string  `json:"id"`
	Status string  `json:"status"`
	Result float64 `json:"result"`
}

type Request struct {
	Jwt_token          string `json:"jwt_token"`
	Expression_request string `json:"expression_request"`
}

var AllExpressions []Expression

func NewExpression(w http.ResponseWriter, r *http.Request) {

	newexpression := new(Expression)
	request := new(Request)
	json.NewDecoder(r.Body).Decode(&request)

	new_uuid := uuid.New().String()
	newexpression.Id = new_uuid
	newexpression.Status = "in process"

	json.NewEncoder(w).Encode(map[string]string{"Id": newexpression.Id})
	log.Println(request.Expression_request)
	result, err := calculator.Calc(request.Expression_request)
	log.Println("result:", result)

	if err != nil {
		newexpression.Status = "expression has problems"
		fmt.Print(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		newexpression.Status = "calculation is completed"
		newexpression.Result = result
	}

	AllExpressions = append(AllExpressions, *newexpression)

}

func GetAllExpressions(w http.ResponseWriter, r *http.Request) {

	var result struct {
		Expressions []Expression `json:"expressions"`
	}

	result.Expressions = AllExpressions

	json.NewEncoder(w).Encode(&result)

}

func GetExpression(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Path[len("/api/v1/expressions/"):]
	for i := range AllExpressions {
		if AllExpressions[i].Id == id {
			var result struct {
				Expression Expression `json:"expression"`
			}
			result.Expression = AllExpressions[i]
			json.NewEncoder(w).Encode(&result)
		}
	}

}

func AuthMiddleware(sqlDB *auth.App, next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var jwtKey = []byte("super_secret_signature")

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Missing or invalid Authrication header", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}
		var UserId string
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			UserId = claims["user_id"].(string)

		} else {
			http.Error(w, "cannot parse claims", http.StatusInternalServerError)
		}
		fmt.Println(UserId)
		exist, err := sqlDB.Compare(UserId)
		if err != nil {
			http.Error(w, fmt.Sprintf("cannot compare %v", err), http.StatusInternalServerError)
		}

		if !exist {
			http.Error(w, "you should authorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})

}

func OrkestratorRun(db *sql.DB) {

	sqlDB := auth.NewApp(db)

	http.HandleFunc("/api/v1/calculate", AuthMiddleware(sqlDB, http.HandlerFunc(NewExpression)))
	http.HandleFunc("/api/v1/expressions", GetAllExpressions)
	http.HandleFunc("/api/v1/expressions/", GetExpression)
	http.HandleFunc("/api/v1/register", sqlDB.Register)
	http.HandleFunc("/api/v1/login", sqlDB.SingIn)

	port := config.New()

	http.ListenAndServe(":"+port.Port, nil)
}
