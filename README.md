🚀 Task Manager REST API (Golang)

A production-style Task Management Service built in Go using Gin, GORM, PostgreSQL, JWT authentication, Role-Based Access Control, background workers (goroutines & channels), pagination, filtering, Swagger docs, Docker, and clean architecture.

✨ Features

JWT Authentication (Register / Login)

Role Based Access Control (User / Admin)

Task CRUD APIs

Background Worker: auto-completes tasks after X minutes

Pagination & filtering support

Clean Architecture (Handler → Service → Repository)

Context propagation across layers

UUID for primary keys

Unit tests using mocks & interfaces

Swagger (OpenAPI) documentation

Graceful server shutdown

Docker support

Makefile for standardized commands

🧰 Tech Stack

Go (Golang)

Gin Web Framework

GORM ORM

PostgreSQL

JWT Authentication

Goroutines & Channels

Swagger (swaggo)

Docker

📁 Project Structure
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
Makefile
Dockerfile

⚙️ Environment Variables (.env)

Create a .env file in root:

DB_URL=postgres://dipak:123456@localhost:5432/taskdb?sslmode=disable
JWT_SECRET=supersecret
AUTO_COMPLETE_MINUTES=2
PORT=8080

🐘 PostgreSQL Setup
CREATE USER dipak WITH PASSWORD '123456';
CREATE DATABASE taskdb OWNER dipak;

▶️ Run Locally (Developer Mode)
make setup     # install swag & air (first time)
make swag      # generate swagger docs
make run       # start server


Swagger UI:

http://localhost:8080/swagger/index.html

🏗️ Build & Run Binary (Production Style)
make build
./bin/task-manager


or

make start

🧪 Run Tests
make test

🔐 Authentication APIs
Method	Endpoint	Description
POST	/register	Register new user
POST	/login	Login & get JWT
✅ Task APIs (JWT Required)
Method	Endpoint	Description
POST	/tasks	Create task
GET	/tasks	List tasks (pagination & filtering)
GET	/tasks/{id}	Get task by ID
DELETE	/tasks/{id}	Delete task

Query Parameters:

/tasks?page=1&limit=10&status=pending

⚙️ Background Worker (Goroutines & Channels)

When a task is created, its ID is pushed into a channel.

A goroutine listens to the channel.

After AUTO_COMPLETE_MINUTES, it marks the task as completed
if still pending or in_progress.

This demonstrates Go concurrency without blocking API requests.

📘 Swagger API Docs

After running:

make swag


Open:

/swagger/index.html

🐳 Docker Support

Build image:

docker build -t task-manager .


Run container:

docker run -p 8080:8080 --env-file .env task-manager

🧠 Architecture Highlights (Interview Points)

Repository pattern for DB abstraction

Service layer for business logic

Middleware for JWT authentication

Interfaces for testability and mocking

Context passed through all layers

Graceful shutdown using context and http.Server

Makefile for standardized development workflow

👨‍💻 Author

Dipak Bharade
