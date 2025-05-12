package orkestrator

import (
	"context"
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

type Expressions struct {
	Id     string  `json:"id"`
	Expression string `json:"expression"`
	Result float64 `json:"result"`
	Status string `json:"status"`
}

type Request struct {
	Jwt_token          string `json:"jwt_token"`
	Expression_request string `json:"expression_request"`
}

type HandlerContext struct{
	DB *auth.App
	user_id string
}

func NewExpression(w http.ResponseWriter, r *http.Request) {

	newexpression := new(Expressions)
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

	contex, ok := r.Context().Value(DB).(*HandlerContext)
	if !ok {
		http.Error(w, "cannot get user_id and sqlDB", http.StatusInternalServerError)
		return
	}

	_, err = contex.DB.DB.Exec("INSERT INTO expressions(expression, result, user_id) VALUES(?,?,?)", request.Expression_request, result, contex.user_id)
	if err != nil{
		http.Error(w, "cannot add expression", http.StatusInternalServerError)
	}

}

func GetAllExpressions(w http.ResponseWriter, r *http.Request) {

	contex, ok := r.Context().Value(DB).(*HandlerContext)
	if !ok {
		http.Error(w, "cannot get user_id and sqlDB", http.StatusInternalServerError)
		return
	}
	

	rows, err := contex.DB.DB.Query("SELECT id, expression, result FROM expressions WHERE user_id = ?", contex.user_id)
	if err != nil{
		http.Error(w, "cannot get expressions", http.StatusInternalServerError)
		return
	}
	log.Println("данные получили")
	defer rows.Close()

	var expressions []Expressions

	for rows.Next(){
		log.Println("данные перебираем")
		var exp Expressions

		if err := rows.Scan(&exp.Id, &exp.Expression, &exp.Result); err != nil{
			http.Error(w, "cannot scan expression", http.StatusInternalServerError)
		}
		expressions = append(expressions, exp)
		
		if err := rows.Err(); err != nil{
			http.Error(w, "cannot read rows", http.StatusInternalServerError)
		}

	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.Encode(&expressions)


}

func GetExpression(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Path[len("/api/v1/expressions/"):]
	log.Println(id)
	contex, ok := r.Context().Value(DB).(*HandlerContext)
	if !ok {
		http.Error(w, "cannot get user_id and sqlDB", http.StatusInternalServerError)
		return
	}

	var GetExpressionResponse struct{
		Expression string `json:"expression"`
		Result string `json:"result"`
	}

	err := contex.DB.DB.QueryRow("SELECT expression, result FROM expressions WHERE id = ?", id).Scan(&GetExpressionResponse.Expression, &GetExpressionResponse.Result)
	if err != nil{
		http.Error(w, "cannot get expression", http.StatusInternalServerError)
		return
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.Encode(&GetExpressionResponse)


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
		log.Println("!", UserId)
		exist, err := sqlDB.Compare(UserId)
		if err != nil {
			http.Error(w, fmt.Sprintf("cannot compare %v", err), http.StatusInternalServerError)
		}

		if !exist {
			http.Error(w, "you should authorized", http.StatusUnauthorized)
			return
		}

		contex := HandlerContext{
			DB: sqlDB,
			user_id: UserId,
		}
		ctx := context.WithValue(r.Context(), DB, &contex)
		next.ServeHTTP(w, r.WithContext(ctx))
	})

}

func OrkestratorRun(db *sql.DB) {

	sqlDB := auth.NewApp(db)

	http.HandleFunc("/api/v1/calculate", AuthMiddleware(sqlDB, http.HandlerFunc(NewExpression)))
	http.HandleFunc("/api/v1/expressions", AuthMiddleware(sqlDB, http.HandlerFunc(GetAllExpressions)))
	http.HandleFunc("/api/v1/expressions/", AuthMiddleware(sqlDB, http.HandlerFunc(GetExpression)))
	http.HandleFunc("/api/v1/register", sqlDB.Register)
	http.HandleFunc("/api/v1/login", sqlDB.SingIn)

	port := config.New()

	http.ListenAndServe(":"+port.Port, nil)
}
