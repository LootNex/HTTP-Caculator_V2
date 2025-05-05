package orkestrator_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"Calculator_V2/internal/orkestrator"
)

type Response struct {
	Id string `json:"id"`
}

func TestNewExpression(t *testing.T) {
	reqBody, _ := json.Marshal(map[string]string{
	  "expression_request": "2+3*5-(2+1)/3",
	})
  
	req := httptest.NewRequest("POST", "/api/v1/calculate", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
  
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(orkestrator.NewExpression)
  
	handler.ServeHTTP(rr, req)
  
	log.Println("Ответ сервера (строка):", rr.Body.String())
	log.Printf("Сырые байты ответа: %q\n", rr.Body.Bytes())
  
	if rr.Code != http.StatusOK {
	  t.Errorf("Ожидался статус 200, получен %d", rr.Code)
	}
  
	var response Response
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
	  t.Fatalf("Ошибка декодирования JSON: %v", err)
	}
  
	if response.Id == "" {
	  t.Errorf("Ожидался ID, но он пустой")
	}
  }

func TestGetAllExpressions(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/v1/expressions", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(orkestrator.GetAllExpressions)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Ожидался статус 200, получен %d", rr.Code)
	}

	var response struct {
		Expressions []orkestrator.Expression `json:"expressions"`
	}
	json.Unmarshal(rr.Body.Bytes(), &response)

	if len(response.Expressions) == 0 {
		t.Errorf("Ожидался список выражений, но он пуст")
	}
}

func TestGetExpression(t *testing.T) {
	reqBody := `{"expression_request": "2+3*5"}`
	req := httptest.NewRequest("POST", "/api/v1/calculate", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(orkestrator.NewExpression)
	handler.ServeHTTP(rr, req)

	var createResp Response
	json.Unmarshal(rr.Body.Bytes(), &createResp)

	getReq := httptest.NewRequest("GET", "/api/v1/expressions/"+createResp.Id, nil)
	getRR := httptest.NewRecorder()
	getHandler := http.HandlerFunc(orkestrator.GetExpression)

	getHandler.ServeHTTP(getRR, getReq)

	if getRR.Code != http.StatusOK {
		t.Errorf("Ожидался статус 200, получен %d", getRR.Code)
	}

	var getResp struct {
		Expression orkestrator.Expression `json:"expression"`
	}
	json.Unmarshal(getRR.Body.Bytes(), &getResp)

	if getResp.Expression.Id != createResp.Id {
		t.Errorf("Ожидался ID %s, получен %s", createResp.Id, getResp.Expression.Id)
	}
}
