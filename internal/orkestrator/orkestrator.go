package orkestrator

import (
	"Calculator_V2/internal/agent"
	calculator "Calculator_V2/pkg"
	config "Calculator_V2/pkg/config"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type Expression struct {
	Id     string  `json:"id"`
	Status string  `json:"status"`
	Result float64 `json:"result"`
}

type Request struct {
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

func TaskHandler(w http.ResponseWriter, r *http.Request) {

	switch {
	case r.Method == "GET":

		if len(calculator.Tasks) == 0{
			http.Error(w, "Нет задач", http.StatusNoContent)
			return
		}
		json.NewEncoder(w).Encode(&calculator.Tasks[0])

	case r.Method == "POST":

		solvedtask := new(agent.Solved_Task)

		json.NewDecoder(r.Body).Decode(&solvedtask)
		
		calculator.Tasks[0].Operation_time = solvedtask.Operation_time

		calculator.Task_Ready <- solvedtask.Result

	}

}

func OrkestratorRun() {

	http.HandleFunc("/api/v1/calculate", NewExpression)
	http.HandleFunc("/api/v1/expressions", GetAllExpressions)
	http.HandleFunc("/api/v1/expressions/", GetExpression)
	http.HandleFunc("/internal/task", TaskHandler)

	port := config.New()

	http.ListenAndServe(":"+port.Port, nil)
}
