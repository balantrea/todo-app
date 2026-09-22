# Todo App

A REST API for managing users, todo lists, and tasks.

The project is built with **Go**, **Gin**, **PostgreSQL**, **SQLX**, **JWT**, and **Docker Compose**.

---

## 🚀 Features

* User registration
* User authentication with JWT
* Create todo lists
* Get all user's todo lists
* Get a todo list by ID
* Update a todo list
* Delete a todo list
* Create todo items inside a list
* Get all items from a list
* Get an item by ID
* Update todo items
* Mark items as completed
* Delete todo items
* User ownership validation
* PostgreSQL database migrations
* Docker Compose setup
* Graceful HTTP server shutdown
* 5-second graceful shutdown timeout

---

## 🛠 Tech Stack

| Technology     | Purpose                                |
| -------------- | -------------------------------------- |
| Go             | Main programming language              |
| Gin            | HTTP web framework                     |
| PostgreSQL     | Database                               |
| SQLX           | Database access                        |
| JWT            | Authentication                         |
| Viper          | Configuration                          |
| Zerolog        | Logging                                |
| Docker         | Containerization                       |
| Docker Compose | Application and database orchestration |
| golang-migrate | Database migrations                    |

---

## 📁 Project Structure

```text
todo-app/
├── cmd/
│   └── main.go
│
├── configs/
│   └── config.yml
│
├── internal/
│   └── pkg/
│       ├── handler/
│       │   ├── auth.go
│       │   ├── handler.go
│       │   ├── items.go
│       │   ├── lists.go
│       │   ├── middleware.go
│       │   └── response.go
│       │
│       ├── repository/
│       │   ├── auth_postgres.go
│       │   ├── list_item_postgres.go
│       │   ├── postgres.go
│       │   ├── repository.go
│       │   └── todo_list_postgres.go
│       │
│       └── service/
│           ├── auth.go
│           ├── list_item.go
│           ├── service.go
│           └── todo_list.go
│
├── schema/
│   ├── 000001_init.up.sql
│   └── 000001_init.down.sql
│
├── Dockerfile
├── docker-compose.yml
├── Makefile
├── go.mod
├── go.sum
├── .env
└── todo.go
```

---

## 🏗 Architecture

The application follows a layered architecture:

```text
                  HTTP Request
                       │
                       ▼
                ┌────────────┐
                │  Handler   │
                └─────┬──────┘
                      │
                      ▼
                ┌────────────┐
                │  Service   │
                └─────┬──────┘
                      │
                      ▼
               ┌──────────────┐
               │ Repository   │
               └──────┬───────┘
                      │
                      ▼
                ┌────────────┐
                │ PostgreSQL │
                └────────────┘
```

### Handler

Responsible for HTTP-related operations:

* Parsing request parameters
* Reading JSON request bodies
* Validating input
* Extracting the user ID from JWT
* Returning HTTP responses

### Service

Contains the application's business logic.

For example, before creating an item, the service verifies that the requested list belongs to the authenticated user.

### Repository

Responsible for communication with PostgreSQL:

* `SELECT`
* `INSERT`
* `UPDATE`
* `DELETE`
* Transactions
* SQL queries and joins

---

# 🔐 Authentication

The API uses **JWT-based authentication**.

## Sign Up

```http
POST /auth/sign-up
```

Request body:

```json
{
  "name": "John",
  "username": "john",
  "password": "password"
}
```

Response:

```json
{
  "id": 1
}
```

---

## Sign In

```http
POST /auth/sign-in
```

Request body:

```json
{
  "username": "john",
  "password": "password"
}
```

Response:

```json
{
  "token": "YOUR_JWT_TOKEN"
}
```

The JWT token must be included in protected requests:

```http
Authorization: Bearer YOUR_JWT_TOKEN
```

---

# 📋 Lists API

All `/api/*` endpoints require authentication.

## Create a List

```http
POST /api/list/
```

Request body:

```json
{
  "title": "Shopping",
  "description": "Things I need to buy"
}
```

Response:

```json
{
  "id": 1
}
```

---

## Get All Lists

```http
GET /api/list/
```

Response:

```json
{
  "data": [
    {
      "id": 1,
      "title": "Shopping",
      "description": "Things I need to buy"
    }
  ]
}
```

---

## Get a List by ID

```http
GET /api/list/:id
```

Example:

```http
GET /api/list/1
```

---

## Update a List

```http
PUT /api/list/:id
```

Request body:

```json
{
  "title": "Updated shopping list",
  "description": "Updated description"
}
```

Partial updates are supported:

```json
{
  "title": "New title"
}
```

---

## Delete a List

```http
DELETE /api/list/:id
```

Example:

```http
DELETE /api/list/1
```

---

# ✅ Items API

Todo items belong to a specific list.

Creating and retrieving all items is done through the list:

```text
/api/list/:listId/items
```

Operations on a specific item are performed using the item's ID:

```text
/api/items/:itemId
```

---

## Create an Item

```http
POST /api/list/:listId/items/
```

Example:

```http
POST /api/list/1/items/
```

Request body:

```json
{
  "title": "Buy milk",
  "description": "2 liters"
}
```

Response:

```json
{
  "id": 1
}
```

---

## Get All Items from a List

```http
GET /api/list/:listId/items/
```

Example:

```http
GET /api/list/1/items/
```

Response:

```json
{
  "data": [
    {
      "id": 1,
      "title": "Buy milk",
      "description": "2 liters",
      "done": false
    },
    {
      "id": 2,
      "title": "Buy bread",
      "description": "White bread",
      "done": true
    }
  ]
}
```

---

## Get an Item by ID

```http
GET /api/items/:itemId
```

Example:

```http
GET /api/items/1
```

Response:

```json
{
  "id": 1,
  "title": "Buy milk",
  "description": "2 liters",
  "done": false
}
```

---

## Update an Item

```http
PUT /api/items/:itemId
```

Individual fields can be updated independently.

### Update title

```json
{
  "title": "Buy oat milk"
}
```

### Update description

```json
{
  "description": "3 liters"
}
```

### Mark item as completed

```json
{
  "done": true
}
```

### Mark item as incomplete

```json
{
  "done": false
}
```

The `*bool` field used for `done` allows the application to distinguish between:

```text
nil   → field was not provided
true  → explicitly set to true
false → explicitly set to false
```

---

## Delete an Item

```http
DELETE /api/items/:itemId
```

Example:

```http
DELETE /api/items/1
```

---

# 🗄 Database

The application uses PostgreSQL.

The main tables are:

```text
users
   │
   ▼
users_lists
   │
   ▼
todo_lists
   │
   ▼
list_item
   │
   ▼
todo_items
```

### `users`

Stores application users:

```text
id
name
username
password_hash
```

### `todo_lists`

Stores todo lists:

```text
id
title
description
```

### `users_lists`

Connects users with their lists:

```text
user_id
list_id
```

### `todo_items`

Stores individual tasks:

```text
id
title
description
done
```

### `list_item`

Connects items with lists:

```text
list_id
item_id
```

This relationship structure allows the application to determine which items belong to which lists and whether a list or item belongs to the authenticated user.

---

# 🐳 Running with Docker

## Requirements

Make sure you have installed:

* Docker
* Docker Compose

---

## 1. Configure `.env`

Create a `.env` file in the project root.

Example:

```env
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=todo
SIGNING_KEY=your-secret-key
```

`SIGNING_KEY` is used to sign JWT tokens.

---

## 2. Start the Application

```bash
docker compose up -d
```

Or use the Makefile:

```bash
make up
```

Docker Compose will start:

```text
PostgreSQL
    ↓
Database migrations
    ↓
Todo API
```

The application will be available at:

```text
http://localhost:8000
```

---

# 🛠 Makefile

The project provides several useful commands.

### Start

```bash
make up
```

### Stop

```bash
make down
```

### Restart

```bash
make restart
```

### View logs

```bash
make logs
```

### Build

```bash
make build
```

### Rebuild

```bash
make rebuild
```

### Run migrations

```bash
make migrate-up
```

### Roll back the latest migration

```bash
make migrate-down
```

### Remove Docker volumes

```bash
make clean
```

> ⚠️ `make clean` removes Docker volumes, so PostgreSQL data will be deleted.

---

# 🔄 Database Migrations

Database migrations are located in:

```text
schema/
```

Initial migration:

```text
000001_init.up.sql
```

Rollback migration:

```text
000001_init.down.sql
```

When the application is started through Docker Compose, the database migrations are applied before the application starts.

---

# 🔒 Graceful Shutdown

The application handles the following operating system signals:

```text
SIGINT
SIGTERM
```

When a shutdown signal is received:

1. The HTTP server stops accepting new requests.
2. Active requests are given time to finish.
3. A **5-second timeout** is applied.
4. The PostgreSQL connection is closed.
5. The application exits.

```text
SIGTERM / SIGINT
       │
       ▼
Graceful HTTP Shutdown
       │
       │ max 5 seconds
       ▼
Close PostgreSQL
       │
       ▼
Application Exit
```

---

# 🧪 Example

After starting the application, you can test the API with `curl`.

## 1. Register

```bash
curl -X POST http://localhost:8000/auth/sign-up \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John",
    "username": "john",
    "password": "password"
  }'
```

## 2. Sign In

```bash
curl -X POST http://localhost:8000/auth/sign-in \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john",
    "password": "password"
  }'
```

Copy the returned JWT token.

## 3. Create a List

```bash
curl -X POST http://localhost:8000/api/list/ \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Shopping",
    "description": "Things to buy"
  }'
```

## 4. Create an Item

```bash
curl -X POST http://localhost:8000/api/list/1/items/ \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Milk",
    "description": "2 liters"
  }'
```

## 5. Mark the Item as Completed

```bash
curl -X PUT http://localhost:8000/api/items/1 \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "done": true
  }'
```

---

# 📌 API Overview

| Method   | Endpoint               | Description         | Auth |
| -------- | ---------------------- | ------------------- | ---- |
| `POST`   | `/auth/sign-up`        | Register a user     | ❌    |
| `POST`   | `/auth/sign-in`        | Authenticate a user | ❌    |
| `POST`   | `/api/list/`           | Create a list       | ✅    |
| `GET`    | `/api/list/`           | Get all lists       | ✅    |
| `GET`    | `/api/list/:id`        | Get a list          | ✅    |
| `PUT`    | `/api/list/:id`        | Update a list       | ✅    |
| `DELETE` | `/api/list/:id`        | Delete a list       | ✅    |
| `POST`   | `/api/list/:id/items/` | Create an item      | ✅    |
| `GET`    | `/api/list/:id/items/` | Get list items      | ✅    |
| `GET`    | `/api/items/:id`       | Get an item         | ✅    |
| `PUT`    | `/api/items/:id`       | Update an item      | ✅    |
| `DELETE` | `/api/items/:id`       | Delete an item      | ✅    |

---

# 📚 What This Project Demonstrates

This project demonstrates practical REST API development with Go, including:

* RESTful routing
* HTTP middleware
* JWT authentication
* Handler / Service / Repository architecture
* PostgreSQL integration
* SQL JOINs
* Database transactions
* `QueryRow`, `Exec`, and `Select`
* Database migrations
* Docker and Docker Compose
* Configuration with Viper
* Structured logging with Zerolog
* Graceful server shutdown
* Partial updates using pointer fields
* User ownership and authorization checks

---

## 📄 License

This project is intended for educational and personal use.
