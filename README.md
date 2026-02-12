# Task Manager REST API (Golang)

A production-style Task Management Service built in Go using Gin, GORM, PostgreSQL, JWT authentication, Role-Based Access Control, background workers with goroutines & channels, pagination, filtering, Swagger docs, Docker support, and clean architecture.

---

## Features

* JWT Authentication (Register / Login)
* Role Based Access Control (User / Admin)
* Task CRUD APIs
* Background Worker: auto-completes tasks after X minutes
* Pagination & filtering
* Clean architecture (handler → service → repository)
* Context usage in all layers
* UUID for IDs
* Unit tests with mocks
* Swagger (OpenAPI) documentation
* Graceful shutdown
* Dockerfile

---

## Tech Stack

* Go
* Gin
* GORM
* PostgreSQL
* JWT
* Goroutines & Channels
* Swagger (swaggo)
* Docker

---

## Project Structure

```
cmd/
config/
models/
repository/
service/
handler/
middleware/
worker/
utils/
docs/
```

---

## Environment Variables (.env)

```
DB_URL=postgres://dipak:123456@localhost:5432/taskdb?sslmode=disable
JWT_SECRET=supersecret
AUTO_COMPLETE_MINUTES=2
```

---

## Run Locally

```bash
go mod tidy
go run cmd/main.go
```

Swagger UI:

```
http://localhost:8080/swagger/index.html
```

---

## PostgreSQL Setup

```sql
CREATE USER dipak WITH PASSWORD '123456';
CREATE DATABASE taskdb OWNER dipak;
```

---

## Authentication APIs

| Method | Endpoint  | Description   |
| ------ | --------- | ------------- |
| POST   | /register | Register user |
| POST   | /login    | Get JWT token |

---

## Task APIs (JWT Required)

| Method | Endpoint    | Description                     |
| ------ | ----------- | ------------------------------- |
| POST   | /tasks      | Create task                     |
| GET    | /tasks      | List tasks (pagination, filter) |
| GET    | /tasks/{id} | Get task by ID                  |
| DELETE | /tasks/{id} | Delete task                     |

Query params:

```
/tasks?page=1&limit=10&status=pending
```

---

## Background Worker

When a task is created, its ID is pushed into a channel.
A goroutine waits for **AUTO_COMPLETE_MINUTES** and automatically marks the task as `completed` if it is still `pending` or `in_progress`.

This demonstrates Go concurrency using goroutines and channels without blocking API requests.

---

## Unit Tests

Service layer is tested using repository interfaces and mocks.

```bash
go test ./... -v
```

---

## Docker

```bash
docker build -t task-manager .
docker run -p 8080:8080 task-manager
```

---

## Swagger Docs

Interactive API documentation available at:

```
/swagger/index.html
```

---

## Architecture Highlights

* Repository pattern for DB abstraction
* Service layer for business logic
* Middleware for JWT
* Interfaces for testability
* Context passed across layers
* Clean separation of concerns

---

## Author

Dipak Bharade
