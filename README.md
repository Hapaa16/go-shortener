# 🚀 Go URL Shortener (Gin + MySQL)

A simple and lightweight **URL shortener service** built using **Go (Gin framework)** and **MySQL**.  
It provides RESTful APIs to shorten URLs, update or delete existing ones, and retrieve access statistics.

https://roadmap.sh/projects/url-shortening-service
---

## 🐳 Running the Project with Docker Compose

### 1. Clone the Repository

```bash
git clone https://github.com/Hapaa16/go-shortener.git
cd go-shortener
```

### 2. Set Up Environment Variables

Create a `.env` file in the root directory (or export variables manually):

```bash
MYSQL_DATABASE=appdb
MYSQL_USER=shortener-user
MYSQL_PASSWORD=pass123
MYSQL_ROOT_PASSWORD=pass123
```

### 3. Build and Start Containers

```bash
docker-compose up -d
```

The app should now be running on **http://localhost:8080** 🎉

---

## 📦 Example Usage

### 🔗 Shorten a URL

```bash
curl -X POST -H "Content-Type: application/json" \
  -d '{"url": "https://google.com"}' \
  http://localhost:8080/api/shorten
```

### ✏️ Update an Existing Shortened URL

```bash
curl -X PUT -H "Content-Type: application/json" \
  -d '{"url": "https://example.com"}' \
  http://localhost:8080/api/shorten/abc
```

### ❌ Delete a Shortened URL

```bash
curl -X DELETE http://localhost:8080/api/shorten/abc
```

### 🔍 Retrieve a Shortened URL

```bash
curl -X GET http://localhost:8080/api/shorten/abc
```

### 📊 Get URL Access Statistics

```bash
curl -X GET http://localhost:8080/api/shorten/abc/stats
```

---

## 🔗 API Endpoints Overview

| Method | Endpoint                        | Description                      |
|--------:|----------------------------------|----------------------------------|
| **POST**   | `/api/shorten`                  | Create a shortened URL           |
| **PUT**    | `/api/shorten/:code`            | Update an existing short URL     |
| **DELETE** | `/api/shorten/:code`            | Delete a short URL               |
| **GET**    | `/api/shorten/:code`            | Redirect to original URL            |
| **GET**    | `/api/shorten/:code/stats`      | Get access statistics            |

---

## 🧱 Tech Stack

- **Go** — Web server using the Gin framework  
- **MySQL** — Persistent data storage  



