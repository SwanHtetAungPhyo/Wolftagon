
# Swan Htet Aung Phyo
 Solution 
---

# 🛡️ Auth Service

This service offers secure **user registration**, **login**, **JWT-based authentication**, **role-based access control (RBAC)**, **token refresh**, and **logout** functionality using **Go** and **PostgreSQL**, with Redis for email verification and token blacklisting.

---
# Base Url - https://localhost:8081
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

---

# 📝 Combined API Documentation


---


### **1.1 User Registration**

#### **Endpoint**

`POST /auth/register`

This endpoint is used to register a new user in the system. Upon successful registration, a verification email is sent to the user's email address.

---

#### **Request**

```http
POST {{baseUrl}}/auth/register
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

#### **Response**

```json
{
  "status_code": 201,
  "message": "Registration successful",
  "data": {
    "message": "Registration successful. Verification email sent."
  }
}
```

---

#### **Response Status Codes**

* `201 Created`: The user was successfully registered, and a verification email was sent.

---

### **1.2 Verify Email**

#### **Endpoint**

`POST /auth/verify`

After registration, the user will receive an email with a verification code. This endpoint is used to verify the user's email address.

---

#### **Request**

```http
POST {{baseUrl}}/auth/verify
Content-Type: application/json
```

```json
{
  "code": "046416",
  "email": "swanhtetaungp@gmail.com"
}
```

---

#### **Response**

```json
{
  "status_code": 200,
  "message": "Email verified successfully"
}
```

---

#### **Response Status Codes**

* `200 OK`: Email successfully verified.

---


## 🔐 2. Login

### **Endpoint**

`POST /auth/login`

Authenticates a user and returns:

* An **access token** in the response
* A **refresh token** via `Set-Cookie`

---

### **Request**

```http
POST {{baseUrl}}/auth/login
Content-Type: application/json
```

```json
{
  "email": "swanhtetaungp@gmail.com",
  "password": "securePassword123"
}
```

---

### **Response**

```json
{
  "status_code": 200,
  "message": "Login successful",
  "data": {
    "message": "Login successful",
    "token": "<ACCESS_TOKEN>",
    "user_metadata": {
      "user_id": "3dac6eda-e7cf-407e-8281-06cf4c64b552",
      "email": "swanhtetaungp@gmail.com",
      "first_name": "Swan Htet",
      "role_name": "user"
    }
  }
}
```

---

## 👤 3. Get User Info

### **Endpoint**

`GET /user`

Returns user-specific info if authenticated with an access token.

---

### **Request**

```http
GET {{baseUrl}}/user
Authorization: Bearer <ACCESS_TOKEN>
Content-Type: application/json
```

---

### **Response**

```json
{
  "message": "Welcome  Wolftagon User",
  "success": true
}
```

---

## 🔒 4. Access Admin Route with User Token

### **Endpoint**

`GET /admin`

Fails if the token does not belong to an admin.

---

### **Request**

```http
GET {{baseUrl}}/admin
Authorization: Bearer <USER_ACCESS_TOKEN>
Content-Type: application/json
```

---

### **Response**

```json
{
  "error": "Access denied"
}
```

**HTTP Status**: `403 Forbidden`

---

## ♻️ 5. Refresh Token

### **Endpoint**

`GET /refresh`

Refreshes the access and refresh tokens using the existing refresh token from a secure HTTP-only cookie.

---

### **Request**

```http
GET {{baseUrl}}/refresh
Authorization: Bearer <OLD_ACCESS_TOKEN>
```

(Cookies with `refresh_token` must be included, usually automatically sent by browser or client.)

---

### **Response**

```json
{
  "status_code": 200,
  "message": "Token refreshed",
  "data": {
    "access_token": "<NEW_ACCESS_TOKEN>",
    "refresh_token": "<NEW_REFRESH_TOKEN>"
  }
}
```

Here’s the documentation for the **Logout** process:

---

## 🚪 User Logout

### **Endpoint**

`POST /logout`

This endpoint is used to log out a user from the system. The access token is invalidated, and the refresh token is cleared.

---

#### **Request**

```http
POST {{baseUrl}}/logout
Content-Type: application/json
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDY3OTI0MjcsImlhdCI6MTc0Njc5MDYyNywicm9sZSI6InVzZXIiLCJzdWIiOiJlY2FlYTRmZC1mNmViLTQwZTYtYjMzZC0yZjhmNzUzMzYzMGQiLCJ0eXBlIjowfQ.Ppbgprw4TEwrH17frHTxH5JISyzB-JuQW3Q_0zljrdY
```

---

#### **Response**

```json
{
  "status_code": 200,
  "message": "Logout successful"
}
```

---

#### **Response Status Codes**

* `200 OK`: The user has been successfully logged out.

---


---

## 🔁 Token Lifecycle Summary

| Step                    | Token Used             | Purpose                                             |
|-------------------------|------------------------|-----------------------------------------------------|
| Login                   | —                      | Get access + refresh tokens                         |
| Access User/Admin Route | Access Token           | Validate identity and role                          |
| Refresh Token           | Refresh Token (in cookie) | Get new access + refresh tokens for token rotation  |
| Logout                  | Auth header            | Revoke all the access token and refresh are revoked |

---



