
# Swan Htet Aung Phyo
 Solution 
---

# 🛡️ Auth Service

This service offers secure **user registration**, **login**, **JWT-based authentication**, **role-based access control (RBAC)**, **token refresh**, and **logout** functionality using **Go** and **PostgreSQL**, with Redis for email verification and token blacklisting.

---

## ⚙️ Tech Stack

* **Go (Fiber)** – Fast and minimalist web framework
* **PostgreSQL** – Relational database for persistent storage
* **Redis** – In-memory store used for:

    * 🔐 Storing email verification codes
    * ❌ Storing revoked JWTs

---

## 🔐 Features

* ✅ **Register** users with roles (e.g., `user`, `admin`)
* ✉️ **Email verification** using Redis (disabled for assessment)
* 🔑 **Login** with hashed password validation
* 🪪 **JWT Authentication** with access and refresh tokens
* 🧠 **RBAC** for endpoint-level permission control
* 🔁 **Token refresh** via refresh token
* 🚪 **Logout** invalidates access token via Redis
* 🔒 **Protected routes**:

    * `/user` → Requires valid user token
    * `/admin` → Requires admin privileges

---

## 🚀 Running the Application

```bash
# Start PostgreSQL and Redis containers
make docker

# Run the application
make run
```

> ⚠️ **Note for Assessment:**
> Email verification is **disabled**. Only pre-registered emails (like yours) are supported by Resend Mail.

---

## 📦 API Endpoints

| Method | Endpoint         | Description                       |
| ------ | ---------------- | --------------------------------- |
| POST   | `/auth/register` | Register a new user               |
| POST   | `/auth/verify`   | Verify email with code (disabled) |
| POST   | `/auth/login`    | Login and receive tokens          |
| GET    | `/refresh`       | Refresh the JWT tokens            |
| POST   | `/logout`        | Logout and revoke token           |
| GET    | `/user`          | Protected route for users         |
| GET    | `/admin`         | Protected route for admins        |

---

## 🗃️ PostgreSQL Schema

```sql
CREATE TABLE roles (
    role_id UUID PRIMARY KEY,
    role_name VARCHAR(100) NOT NULL
);

CREATE TABLE users (
    user_id UUID PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    age INT NOT NULL,
    password TEXT NOT NULL,
    verified BOOLEAN NOT NULL,
    email VARCHAR(150) NOT NULL UNIQUE,
    role_id UUID NOT NULL,
    CONSTRAINT fk_role
        FOREIGN KEY (role_id)
        REFERENCES roles(role_id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT
);
```

---

## 🔐 Auth Flow Overview

1. **Register:** Submit user info and desired role.
2. **(Disabled)** **Verify Email:** Receive and verify email code.
3. **Login:** Receive access and refresh tokens.
4. **Access:** Use access token for protected routes.
5. **Refresh:** Use refresh token to obtain a new access token.
6. **Logout:** Invalidate current token (blacklisted in Redis).

---

## 🌐 Base URL

```txt
https://localhost:8081
```

---

## 📋 Request Examples

### ✅ Register

```http
POST /auth/register
Content-Type: application/json
```

```json
{
  "first_name": "Swan Htet",
  "last_name": "Aung Phyo",
  "email": "swanhtetaungp@gmail.com",
  "password": "securePassword123",
  "role_name": "user",
  "age": 25
}
```

---

### 🔑 Login

```http
POST /auth/login
Content-Type: application/json
```

```json
{
  "email": "swanhtetaungp@gmail.com",
  "password": "securePassword123"
}
```

---

### 🔄 Refresh Token

```http
GET /refresh
Authorization: Bearer <refresh_token>
```

---

### 🚪 Logout

```http
POST /logout
Content-Type: application/json
Authorization: Bearer <access_token>
```

---

### 👤 Get User Data

```http
GET /user
Authorization: Bearer <access_token>
```

---

### 🔒 Get Admin Data

```http
GET /admin
Authorization: Bearer <access_token>
```

---

## 💻 CURL Examples

<details>
<summary>Click to expand</summary>

### ✅ Register

```bash
curl -X POST https://localhost:8081/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Swan Htet",
    "last_name": "Aung Phyo",
    "email": "swanhtetaungp@gmail.com",
    "password": "securePassword123",
    "role_name": "user",
    "age": 25
  }' -k
```

### 🔑 Login

```bash
curl -X POST https://localhost:8081/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "swanhtetaungp@gmail.com",
    "password": "securePassword123"
  }' -k
```

### 🔄 Refresh Token

```bash
curl -X GET https://localhost:8081/refresh \
  -H "Authorization: Bearer <your_refresh_token>" -k
```

### 🚪 Logout

```bash
curl -X POST https://localhost:8081/logout \
  -H "Authorization: Bearer <your_access_token>" \
  -H "Content-Type: application/json" -k
```

### 👤 User Endpoint

```bash
curl -X GET https://localhost:8081/user \
  -H "Authorization: Bearer <your_user_token>" -k
```

### 🔒 Admin Endpoint

```bash
curl -X GET https://localhost:8081/admin \
  -H "Authorization: Bearer <your_admin_token>" -k
```

</details>

---

