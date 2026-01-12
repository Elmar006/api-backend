REST API на Go

REST API для управления задачами (TODO) с аутентификацией пользователей.

Ключевые особенности

- Аутентификация и авторизация через JWT токены
- Создание, чтение и удаление задач через защищенные эндпоинты REST API
- Регистрация и вход пользователей
- Чистая архитектура с разделением на handlers, models, database, auth
- Интеграция с базой данных PostgreSQL с использованием GORM для ORM и автоматических миграций
- Маршрутизация через легковесный роутер chi с middleware для логирования и восстановления
- Обработка ошибок и ответы в формате JSON
- Конкурентобезопасная работа с базой данных
- Unit-тесты с проверкой краевых случаев
- Полная контейнеризация приложения и базы (Docker + Docker Compose)
- Health check эндпоинт для мониторинга
- Автоматизация тестирования и деплоя через GitHub Actions

Стек технологий

- Язык: Go 1.25
- База данных: PostgreSQL (с GORM для ORM)
- Аутентификация: JWT (JSON Web Tokens)
- Роутер: github.com/go-chi/chi/v5
- Middleware: Логирование, Recovery, RequestID
- Конфигурация: Переменные окружения через github.com/joho/godotenv
- Тестирование: Встроенный пакет тестирования Go
- Контейнеризация: Docker, Docker Compose
- CI/CD: GitHub Actions

Структура проекта

```
├── cmd/
│   └── main.go            # Точка входа приложения
├── internal/
│   ├── auth/              # Логика аутентификации (JWT)
│   ├── database/          # Подключение к БД и миграции
│   ├── handlers/          # HTTP обработчики
│   └── models/            # Модели данных (GORM)
├── compose.yaml           # Конфигурация Docker Compose
├── Dockerfile             # Docker образ приложения
├── go.mod                 # Зависимости Go
└── README.md              # Документация
```

Эндпоинты API

Публичные эндпоинты (не требуют аутентификации)

- GET /health - Проверка работоспособности сервера
  ```
  Ответ: 200 OK
  Тело: "OK"
  ```

- POST /register - Регистрация нового пользователя
  ```
  Тело запроса:
  {
    "email": "user@example.com",
    "password": "securepassword"
  }
  
  Ответ: 201 Created
  ```

- POST /login - Вход пользователя и получение токена
  ```
  Тело запроса:
  {
    "email": "user@example.com",
    "password": "securepassword"
  }
  
  Ответ: 200 OK
  Тело: { "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." }
  ```

Защищенные эндпоинты (требуют JWT токен в заголовке Authorization)

- GET /me - Получить информацию о текущем пользователе
  ```
  Заголовок: Authorization: Bearer <JWT_TOKEN>
  Ответ: 200 OK
  Тело: { "id": 1, "email": "user@example.com", "role": "user" }
  ```

- GET /tasks - Получить все задачи текущего пользователя
  ```
  Заголовок: Authorization: Bearer <JWT_TOKEN>
  Ответ: 200 OK
  Тело: [{"id":1,"description":"Задача 1","note":"Детали","user_id":1}]
  ```

- POST /tasks - Создать новую задачу
  ```
  Заголовок: Authorization: Bearer <JWT_TOKEN>
  Тело запроса:
  {
    "description": "Новая задача",
    "note": "Опциональная заметка"
  }
  Ответ: 201 Created
  ```

- GET /task/{id} - Получить задачу по ID
  ```
  Заголовок: Authorization: Bearer <JWT_TOKEN>
  Ответ: 200 OK
  Тело: {"id":1,"description":"Задача 1","note":"Детали","user_id":1}
  ```

- DELETE /task/{id} - Удалить задачу по ID
  ```
  Заголовок: Authorization: Bearer <JWT_TOKEN>
  Ответ: 200 OK
  ```

Быстрый старт

1. Клонирование репозитория
```bash
git clone <repository-url>
cd api-backend
```

2. Запуск с Docker Compose
```bash
docker compose up -d
```

3. Проверка работоспособности
```bash
curl http://localhost:3000/health
# Ответ: OK
```

4. Регистрация пользователя
```bash
curl -X POST http://localhost:3000/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

5. Вход и получение токена
```bash
curl -X POST http://localhost:3000/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
# Сохраните токен из ответа
```

6. Использование API с токеном
```bash
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# Получить информацию о пользователе
curl -H "Authorization: Bearer $TOKEN" http://localhost:3000/me

# Создать задачу
curl -X POST http://localhost:3000/tasks \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"description":"Моя первая задача","note":"Важная заметка"}'

# Получить все задачи
curl -H "Authorization: Bearer $TOKEN" http://localhost:3000/tasks
```


Docker Compose

Сервисы:
- api: Go приложение (порт 3000)
- db: PostgreSQL база данных (порт 5432)

Тестирование

```bash
# Запуск unit-тестов
cd internal/handlers
go test ./...

# Запуск тестов с покрытием
go test ./... -cover

# Запуск конкретного теста
go test ./internal/handlers -v
```

CI/CD

Проект использует GitHub Actions для автоматизации:
- При пуше в main ветку: Запускаются тесты (go test) и проверка кода (go vet)
- При создании тега: Собирается и публикуется Docker образ в Docker Hub

Контакты:
Email: birembekove@gmail.com
Telegram: @viy1ix
github: github.com/Elmar006
