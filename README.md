HTTP Calculator V2

Описание

HTTP-сервис для вычисления математических выражений, работающий с очередями задач и gRPC-агентами. Поддерживает регистрацию, аутентификацию и авторизованный доступ ко всем эндпоинтам API.


---

Запуск

1. Установи зависимости и собери проект

 - git clone https://github.com/LootNex/HTTP-Caculator_V2
 - go mod tidy
 - go run cmd/main.go


---

Аутентификация

Перед обращением к защищённым маршрутам API необходимо пройти регистрацию и авторизацию, чтобы получить JWT токен.
Токен находится в валидном состоянии только 10 минут! Когда время истечёт, придется снова логиниться.

Регистрация пользователя

curl -X POST http://localhost:8080/api/v1/register -H "Content-Type: application/json" -d '{"login": "...", "password": "..."}'

Авторизация

curl -X POST http://localhost:8080/api/v1/login -H "Content-Type: application/json" -d '{"login": "...", "password": "..."}'

Ответ:

{
  "token": "<JWT-TOKEN>"
  "message": "..."
  "status": "..."
}


---

Заголовки авторизации

К каждому защищённому эндпоинту необходимо добавлять заголовок:

Authorization: Bearer <JWT-TOKEN>


---

Эндпоинты

POST /api/v1/calculate

Отправить выражение на вычисление.

  curl -X POST http://localhost:8080/api/v1/calculate  
  -H "Content-Type: application/json"  
  -H "Authorization: Bearer Ваш токен"
  -d '{"expression_request": "2+3*5-(2+1)/3"}'

Ответ:

{
  "Id": "uuid-выражения"
}


---

GET /api/v1/expressions

Получить все выражения пользователя.

curl -X GET http://localhost:8080/api/v1/expressions
  -H "Authorization: Bearer <JWT-TOKEN>"

Ответ:

[
  {
    "id": "...",
    "expression": "...",
    "result": "..."
  },
  ...
]


---

GET /api/v1/expressions/{id}

Получить одно выражение по ID.

curl -X GET http://localhost:8080/api/v1/expressions/<EXPRESSION_ID>
  -H "Authorization: Bearer <JWT-TOKEN>"

Ответ:

{
  "expression": "...",
  "result": "..."
}


---

Примечания

Все ответы отдаются в формате JSON.

При ошибках возвращается 500 Internal Server Error с сообщением об ошибке.

Заголовок Authorization обязателен для всех /api/v1/* маршрутов.