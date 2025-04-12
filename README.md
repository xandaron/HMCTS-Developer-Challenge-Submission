# HMCTS Developer Technical Test Submission

## Overview

This repository contains my submission for the **DTS Developer Technical Test**. The project focuses on building a task management system for HMCTS caseworkers. The API is fully functional and built in **Go**, with **Docker** used for containerized deployment and a **MySQL** database.

> ✅ API functionality is complete  
> ❌ Unit tests not yet implemented  
> ❌ Frontend UI is minimal and currently under development

---

## 🚀 Getting Started

To run the project locally using Docker:

```bash
git clone https://github.com/xandaron/HMCTS-Developer-Challenge-Submission.git
cd HMCTS-Developer-Challenge-Submission
docker-compose up --build
```

The backend will be available at: https://localhost:443

### Demo Login

After starting the application, you can log in with:
- Username: demo
- Password: demo123

### Local Development (without Docker)

To run the project locally without Docker:

1. Install Go (version 1.24+)
2. Install MySQL
3. Create the database:
   ```bash
   mysql -u root -p
   CREATE DATABASE mydb;
   exit;
   ```
4. Initialize the database schema:
   ```bash
   mysql -u root -p mydb < database/seed/init.sql
   ```
5. Set up the environment variables:
   ```
   DB_HOST=localhost
   DB_PORT=3306
   DB_USER=user
   DB_PASSWORD=password
   DB_NAME=mydb
   ```
6. Run the application:
   ```bash
   go run main.go
   ```

## 📦 Tech Stack

- Backend: Golang (Go)
- Database: MySQL
- Containerization: Docker + Docker Compose
- Templating / Frontend: HTML templates (WIP)
- Styles: TailwindCSS

## 🏗️ Application Architecture

The application follows a simple MVC architecture:

- **API Layer**: Handles HTTP requests and responses (`/api` directory)
- **Database Layer**: Manages database connections and queries (`/database` directory)
- **Session Management**: Handles user authentication and sessions (`/session` directory)
- **Error Handling**: Centralized error handling (`/errors` directory)
- **Templates**: Frontend HTML templates (`/templates` directory)

### Authentication Flow

The application uses session-based authentication:

1. User registers or logs in through the `/api/signup` or `/api/login` endpoints
2. A session cookie is created and stored on the client
3. Protected routes check for a valid session before allowing access
4. Sessions expire after 5 minutes of inactivity

<!-- ## 📱 Screenshots -->

## 📚 API Documentation

Base URL
https://localhost:443

### Endpoints

#### Login

<details>
<summary><code>POST</code> <code><b>/api/login</b></code></summary>

##### Checks login credentials and returns a session cookie if successful

##### Parameters

> | name | type     | data type   | description                                           |
> | ---- | -------- | ----------- | ----------------------------------------------------- |
> | None | required | object JSON | `json {"username":<username>, "password":<password>}` |

##### Responses

> | http code | content-type                | response                                   |
> | --------- | --------------------------- | ------------------------------------------ |
> | `200`     | `text/plain; charset=UTF-8` |                                            |
> | `400`     | `text/plain; charset=UTF-8` | `Invalid JSON`                             |
> | `400`     | `application/json`          | `{"message":"user not found"}`             |
> | `400`     | `application/json`          | `{"message":"incorrect password"}`         |
> | `400`     | `application/json`          | `{"message":"empty username or password"}` |
> | `500`     | `text/plain; charset=UTF-8` | `Internal Server Error`                    |

##### Example cURL

```bash
curl -X POST https://localhost:443/api/login -H "content-Type: application/json" -d "{ \"username\": \"test\", \"password\": \"12345\" }" -c cookies.txt -k
```

</details>

#### Sign Up

<details>
<summary><code>POST</code> <code><b>/api/signup</b></code></summary>

##### Creates a user account and returns a session cookie if successful

##### Parameters

> | name | type     | data type   | description                                           |
> | ---- | -------- | ----------- | ----------------------------------------------------- |
> | None | required | object JSON | `json {"username":<username>, "password":<password>}` |

##### Responses

> | http code | content-type                | response                                   |
> | --------- | --------------------------- | ------------------------------------------ |
> | `200`     | `text/plain; charset=UTF-8` |                                            |
> | `400`     | `text/plain; charset=UTF-8` | `Invalid JSON`                             |
> | `400`     | `application/json`          | `{"message":"user already exists"}`        |
> | `400`     | `application/json`          | `{"message":"empty username or password"}` |
> | `500`     | `text/plain; charset=UTF-8` | `Internal Server Error`                    |

##### Example cURL

```bash
curl -X POST https://localhost:443/api/signup -H "content-Type: application/json" -d "{ \"username\": \"newuser\", \"password\": \"12345\" }" -c cookies.txt -k
```

</details>

#### Logout

<details>
<summary><code>GET</code> <code><b>/api/logout</b></code></summary>

##### Ends the user session

##### Responses

> | http code | content-type                | response |
> | --------- | --------------------------- | -------- |
> | `303`     | `text/plain; charset=UTF-8` |          |

##### Example cURL

```bash
curl -X GET https://localhost:443/api/logout -b cookies.txt -k
```

</details>

#### Tasks

<details>
<summary><code>GET</code> <code><b>/api/tasks/</b></code></summary>

##### Get all tasks associated with the current user

##### Responses

> | http code | content-type                | response                                                                                                                                                                              |
> | --------- | --------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `200`     | `application/json`          | `{"tasks": [ {"id": <id>, "user_id": <user_id>, "name": <name>, "description": <description>, "status": <status>, "created_at": <creation date/time> "deadline": <deadline>}, ... ]}` |
> | `401`     | `text/plain; charset=UTF-8` | `Unauthorized`                                                                                                                                                                        |
> | `500`     | `text/plain; charset=UTF-8` | `Internal Server Error`                                                                                                                                                               |

##### Example cURL

```bash
curl -X GET https://localhost:443/api/tasks/ -b cookies.txt -k
```

</details>

<details>
<summary><code>GET</code> <code><b>/api/tasks/task_id</b></code></summary>

##### Get a task associated with the current user

##### Responses

> | http code | content-type                | response                                                                                                                                                          |
> | --------- | --------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `200`     | `application/json`          | `{"id": <id>, "user_id": <user_id>, "name": <name>, "description": <description>, "status": <status>, "created_at": <creation date/time> "deadline": <deadline>}` |
> | `401`     | `text/plain; charset=UTF-8` | `Unauthorized`                                                                                                                                                    |
> | `404`     | `text/plain; charset=UTF-8` | `Task Not Found`                                                                                                                                                  |
> | `500`     | `text/plain; charset=UTF-8` | `Internal Server Error`                                                                                                                                           |

##### Example cURL

```bash
curl -X GET https://localhost:443/api/tasks/<task_id> -b cookies.txt -k
```

</details>

<details>
<summary><code>POST</code> <code><b>/api/tasks/</b></code></summary>

##### Create a new task

##### Parameters

> | name | type     | data type   | description                                                                                       |
> | ---- | -------- | ----------- | ------------------------------------------------------------------------------------------------- |
> | None | required | object JSON | `json {"name": <name>, "description": <description>, "status": <status>, "deadline": <deadline>}` |

##### Responses

> | http code | content-type                | response                |
> | --------- | --------------------------- | ----------------------- |
> | `201`     | `text/plain; charset=UTF-8` |                         |
> | `400`     | `text/plain; charset=UTF-8` | `Invalid JSON`          |
> | `400`     | `text/plain; charset=UTF-8` | `Missing JSON Data`     |
> | `401`     | `text/plain; charset=UTF-8` | `Unauthorized`          |
> | `500`     | `text/plain; charset=UTF-8` | `Internal Server Error` |

##### Example cURL

```bash
curl -X POST https://localhost:443/api/tasks/ -H "content-Type: application/json" -d "{\"name\": \"test\", \"description\": \"\", \"status\": \"INCOMPLETE\", \"deadline\": \"2025-04-16 00:00:00\"}" -b cookies.txt -k
```

</details>

<details>
<summary><code>PUT</code> <code><b>/api/tasks/task_id</b></code></summary>

##### Modify a task

##### Parameters

> | name | type     | data type   | description                                                                                       |
> | ---- | -------- | ----------- | ------------------------------------------------------------------------------------------------- |
> | None | required | object JSON | `json {"name": <name>, "description": <description>, "status": <status>, "deadline": <deadline>}` |

##### Responses

> | http code | content-type                | response                |
> | --------- | --------------------------- | ----------------------- |
> | `204`     | `text/plain; charset=UTF-8` |                         |
> | `400`     | `text/plain; charset=UTF-8` | `Invalid JSON`          |
> | `400`     | `text/plain; charset=UTF-8` | `Missing JSON Data`     |
> | `400`     | `text/plain; charset=UTF-8` | `Task ID Required`      |
> | `401`     | `text/plain; charset=UTF-8` | `Unauthorized`          |
> | `404`     | `text/plain; charset=UTF-8` | `Task Not Found`        |
> | `500`     | `text/plain; charset=UTF-8` | `Internal Server Error` |

##### Example cURL

```bash
curl -X PUT https://localhost:443/api/tasks/<task_id> -H "content-Type: application/json" -d "{\"name\": \"test\", \"description\": \"\", \"status\": \"INCOMPLETE\", \"deadline\": \"2025-04-16 00:00:00\"}" -b cookies.txt -k
```

</details>

<details>
<summary><code>DELETE</code> <code><b>/api/tasks/task_id</b></code></summary>

##### Delete a task

##### Responses

> | http code | content-type                | response                |
> | --------- | --------------------------- | ----------------------- |
> | `204`     | `text/plain; charset=UTF-8` |                         |
> | `400`     | `text/plain; charset=UTF-8` | `Task ID Required`      |
> | `401`     | `text/plain; charset=UTF-8` | `Unauthorized`          |
> | `404`     | `text/plain; charset=UTF-8` | `Task Not Found`        |
> | `500`     | `text/plain; charset=UTF-8` | `Internal Server Error` |

##### Example cURL

```bash
curl -X DELETE https://localhost:443/api/tasks/<task_id> -b cookies.txt -k
```

</details>

<details>
<summary><code>OPTIONS</code> <code><b>/api/tasks/</b></code></summary>

##### Handles CORS preflight requests

##### Responses

> | http code | content-type                | response |
> | --------- | --------------------------- | -------- |
> | `204`     | `text/plain; charset=UTF-8` |          |

##### Example cURL

```bash
curl -X OPTIONS https://localhost:443/api/tasks/ -k
```

</details>

## ⚙️ Validation & Error Handling

All endpoints implement session validation using cookies and return appropriate error codes and messages:

- 401 Unauthorized: Invalid or missing user session
- 500 Internal Server Error: Unexpected issues

All 500 errors are logged internally for debugging.

## 🔒 Security Considerations

- Passwords are hashed using Argon2id
- Session IDs are randomly generated
- All API endpoints validate user permissions
- HTTPS is implemented with self-signed certificates

**Note:** For production deployment, additional security measures would be needed:

- Production-grade SSL certificates
- Rate limiting

## 🗃️ Database Schema

### users

| Field         | Type         | Null | Key | Default | Extra          |
| ------------- | ------------ | ---- | --- | ------- | -------------- |
| id            | int unsigned | NO   | PRI | NULL    | auto_increment |
| name          | varchar(32)  | NO   |     | NULL    |                |
| password_hash | varchar(255) | NO   |     | NULL    |                |

### tasks

| Field         | Type                          | Null | Key | Default           | Extra             |
| ------------- | ----------------------------- | ---- | --- | ----------------- | ----------------- |
| id            | int unsigned                  | NO   | PRI | NULL              | auto_increment    |
| user_id       | int unsigned                  | NO   | MUL | NULL              |                   |
| name          | tinytext                      | NO   |     | NULL              |                   |
| description   | text                          | YES  |     | NULL              |                   |
| status        | enum('COMPLETE','INCOMPLETE') | NO   |     | INCOMPLETE        |                   |
| creation_time | timestamp                     | NO   |     | CURRENT_TIMESTAMP | DEFAULT_GENERATED |
| deadline      | timestamp                     | NO   |     | NULL              |                   |

## 🧪 Testing (Planned)

Unit testing setup is not yet complete. Planned testing strategy:

- **Unit tests** for core business logic using Go's testing package
- **Integration tests** for API endpoints
- **Database tests** with a test database
- **End-to-end tests** for critical user flows

## 🔮 Future Improvements

1. **Complete unit test implementation** for all components
2. Enhance the frontend with a modern JavaScript framework (React/Vue)
3. Implement more robust error handling and validation
4. Implement sorting and filtering capabilities
5. Add user profile management
6. Implement password reset functionality
7. Add task categories and priority levels

## 📌 Notes

- While the frontend is currently basic (using server-side rendered HTML), the API is fully decoupled and can be easily integrated with any modern frontend framework.
- Docker simplifies setup—no need to manually install Go or MySQL.

## 📎 Useful Commands

Rebuild containers

```bash
docker-compose up --build
```

Stop and remove containers

```bash
docker-compose down
```

## ✅ Task Checklist

- ✅ Dockerized deployment
- ✅ API - Create task
- ✅ API - Get task by ID
- ✅ API - Get all tasks
- ✅ API - Update task status
- ✅ API - Delete task
- ✅ API validation
- ✅ API error handling
- ✅ API documentation
- ✅ Database integration
- ❌ Unit tests
- ✅ Basic frontend implementation
- ✅ User authentication
