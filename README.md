# HTTP Calculator V2

HTTP Calculator — это HTTP-сервис для вычисления математических выражений.  
Проект позволяет отправлять математические выражения через API и получать результат.

## 🛠 Возможности
- Отправка математического выражения на сервер
- Получение списка всех вычислений
- Получение результата вычисления по ID

## 📦 Установка и запуск  

### 1. Склонировать репозиторий  

1) git clone https://github.com/LootNex/HTTP-Caculator_V2.git
2) cd Calculator_v2
3) далее вводите go run main.go

🔥 API Методы

1️⃣ Отправка выражения на вычисление


curl -X POST http://localhost:8082/api/v1/calculate \
     -H "Content-Type: application/json" \
     -d '{"expression_request": "2+3*5-(2+1)/3"}'

Ответ:

{
  "Id": "b3e6b5c4-3a6b-4c9f-8f76-d5c5f918ea67"
}

ID можно использовать для получения результата.

2️⃣ Получение всех выражений

Запрос:

curl -X GET http://localhost:8082/api/v1/expressions

Ответ:

{
  "expressions": [
    {
      "id": "b3e6b5c4-3a6b-4c9f-8f76-d5c5f918ea67",
      "status": "calculation is completed",
      "result": 16.0
    }
  ]
}


---

3️⃣ Получение результата по ID

Запрос:

curl -X GET http://localhost:8082/api/v1/expressions/b3e6b5c4-3a6b-4c9f-8f76-d5c5f918ea67

Ответ:

{
  "expression": {
    "id": "b3e6b5c4-3a6b-4c9f-8f76-d5c5f918ea67",
    "status": "calculation is completed",
    "result": 16.0
  }
}


---