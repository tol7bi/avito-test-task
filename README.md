# 📦 PVZ Backend Service

Проект реализует backend-сервис для сотрудников Пунктов Выдачи Заказов (ПВЗ), позволяющий вести учёт поставок товаров, проверку заказов и фиксацию приёмок. Разработано как тестовое задание для стажёра backend-направления (весенняя волна 2025, Avito).

---

## 🚀 Стек технологий

- **Язык**: Go 1.21+
- **Фреймворк**: Fiber
- **База данных**: PostgreSQL (без ORM, `pgx`)
- **gRPC**: список ПВЗ
- **Prometheus**: метрики
- **JWT**: авторизация (в том числе dummy)

---

## ⚙️ Установка и запуск

```bash
git clone https://github.com/tol7bi/pvz-backend.git
cd pvz-backend
docker-compose up --build
```

## Порты

Сервис | Порт
HTTP API | 8080
gRPC API | 3000
Prometheus | 9000
PostgreSQL | 5432


## 🔐 Авторизация
Dummy
```http
POST /dummyLogin
Body: { "role": "employee" | "moderator" }
Response: { "token": "JWT_TOKEN" }
```
Пользовательская

```http
POST /register
POST /login
```

---

## 📚 Endpoints

| Метод | Endpoint                          | Описание                                | Роль      |
|-------|-----------------------------------|------------------------------------------|-----------|
| POST  | /register                         | Регистрация                              | -         |
| POST  | /login                            | Вход                                     | -         |
| POST  | /dummyLogin                       | Получение тестового JWT-токена           | -         |
| POST  | /pvz                              | Создание ПВЗ                             | moderator |
| GET   | /pvz                              | Список ПВЗ с приёмками и товарами        | employee  |
| POST  | /receptions                       | Инициация новой приёмки                  | employee  |
| POST  | /products                         | Добавление товара в текущую приёмку      | employee  |
| POST  | /pvz/:id/delete_last_product      | Удаление последнего добавленного товара  | employee  |
| POST  | /pvz/:id/close_last_reception     | Закрытие текущей приёмки                 | employee  |

## 📂 Структура проекта

```bash
├── cmd/                 # main.go
├── internal/
│   ├── app/             # маршруты
│   ├── http/            # HTTP-хендлеры
│   ├── grpc/            # gRPC-сервер
│   ├── metrics/         # Prometheus
│   ├── middleware/      # JWT, роли, JSON-логирование
│   ├── models/          # структуры
│   └── repository/      # SQL-логика
├── tests/               # unit-тесты
├── Dockerfile
├── docker-compose.yml
└── README.md
``` 
## 🛰 gRPC

Endpoint: pvz.v1.PVZService/GetPVZList

Описание: возвращает все созданные ПВЗ

Пример запроса:

```bash
grpcurl -plaintext localhost:3000 pvz.v1.PVZService/GetPVZList
```

---
 
## 📈 Метрики Prometheus

Метрики доступны на http://localhost:9000/metrics

Технические
- go_goroutines

- http_requests_total

- go_gc_duration_seconds

Бизнесовые
- pvz_created_total

- receptions_created_total

- products_created_total

---
## 👤 Автор

Tolebi Raptayev @tol7bi