# VendorHub API Routes Documentation

## Base URL

```
http://localhost:8080
```

---

## Route Structure Overview

```
├── /health                          # Health check
├── /auth                            # Authentication routes (public)
├── /products                        # Product routes (mixed public/protected)
├── /vendors                         # Vendor routes (public)
├── /me                              # User profile (protected)
└── /admin                           # Admin routes (protected + admin only)
```

---

## 1. HEALTH CHECK

### GET /health

**Authentication:** Not Required

**Description:** Check if server is running

**Response:** 200 OK

```
OK
```

**cURL:**

```bash
curl http://localhost:8080/health
```

---

## 2. AUTHENTICATION ROUTES (Public)

### POST /auth/signup

**Authentication:** Not Required

**Description:** Register a new user (customer, vendor, or admin)

**Request Body:**

```json
{
  "email": "vendor@example.com",
  "password": "password123",
  "first_name": "John",
  "last_name": "Doe",
  "phone": "+1234567890",
  "role": "vendor"
}
```

**Response:** 201 Created

```json
{
  "data": {
    "id": "uuid",
    "email": "vendor@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "role": "vendor",
    "is_approved": false
  },
  "message": "User registered successfully"
}
```

**cURL:**

```bash
curl -X POST http://localhost:8080/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "email": "vendor@example.com",
    "password": "password123",
    "first_name": "John",
    "last_name": "Doe",
    "phone": "+1234567890",
    "role": "vendor"
  }'
```

---

### POST /auth/login

**Authentication:** Not Required

**Description:** Login and get JWT token

**Request Body:**

```json
{
  "email": "vendor@example.com",
  "password": "password123"
}
```

**Response:** 200 OK

```json
{
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": "uuid",
      "email": "vendor@example.com",
      "role": "vendor"
    }
  },
  "message": "Login successful"
}
```

**cURL:**

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "vendor@example.com",
    "password": "password123"
  }'
```

---

## 3. PRODUCT ROUTES

### Public Product Endpoints

#### GET /products/active

**Authentication:** Not Required

**Description:** Get all active products (marketplace view)

**Query Parameters:**

- `page` (optional): Page number (default: 1)
- `page_size` (optional): Items per page (default: 20, max: 100)

**Response:** 200 OK

```json
[
  {
    "id": "uuid",
    "user_id": "vendor-uuid",
    "name": "Laptop",
    "description": "High-performance laptop",
    "price": 999.99,
    "is_active": true,
    "created_at": "2025-01-02T10:00:00Z",
    "updated_at": "2025-01-02T10:00:00Z"
  }
]
```

**cURL:**

```bash
curl http://localhost:8080/products/active
curl http://localhost:8080/products/active?page=1&page_size=20
```

---

#### GET /products/search

**Authentication:** Not Required

**Description:** Search products by name or description

**Query Parameters:**

- `q` (required): Search term

**Response:** 200 OK

```json
[
  {
    "id": "uuid",
    "user_id": "vendor-uuid",
    "name": "Laptop Computer",
    "description": "High-performance laptop for professionals",
    "price": 999.99,
    "is_active": true,
    "created_at": "2025-01-02T10:00:00Z",
    "updated_at": "2025-01-02T10:00:00Z"
  }
]
```

**cURL:**

```bash
curl "http://localhost:8080/products/search?q=laptop"
```

---

#### GET /products/price

**Authentication:** Not Required

**Description:** Get products within a price range

**Query Parameters:**

- `min` (required): Minimum price
- `max` (required): Maximum price

**Response:** 200 OK

```json
[
  {
    "id": "uuid",
    "user_id": "vendor-uuid",
    "name": "Mouse",
    "description": "Wireless mouse",
    "price": 29.99,
    "is_active": true,
    "created_at": "2025-01-02T10:00:00Z",
    "updated_at": "2025-01-02T10:00:00Z"
  }
]
```

**cURL:**

```bash
curl "http://localhost:8080/products/price?min=20&max=50"
```

---

#### GET /products?id={productId}

**Authentication:** Not Required

**Description:** Get a single product by ID

**Response:** 200 OK

```json
{
  "id": "uuid",
  "user_id": "vendor-uuid",
  "name": "Laptop",
  "description": "High-performance laptop",
  "price": 999.99,
  "is_active": true,
  "created_at": "2025-01-02T10:00:00Z",
  "updated_at": "2025-01-02T10:00:00Z"
}
```

**cURL:**

```bash
curl "http://localhost:8080/products?id=product-uuid"
```

---

### Protected Product Endpoints (Vendor Only)

#### POST /products

**Authentication:** Required (JWT Token)

**Description:** Create a new product

**Request Body:**

```json
{
  "name": "Laptop",
  "description": "High-performance laptop",
  "price": 999.99
}
```

**Response:** 201 Created

```json
{
  "id": "uuid",
  "user_id": "vendor-uuid",
  "name": "Laptop",
  "description": "High-performance laptop",
  "price": 999.99,
  "is_active": true,
  "created_at": "2025-01-02T10:00:00Z",
  "updated_at": "2025-01-02T10:00:00Z"
}
```

**cURL:**

```bash
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "Laptop",
    "description": "High-performance laptop",
    "price": 999.99
  }'
```

---

#### PUT /products?id={productId}

**Authentication:** Required (JWT Token)

**Description:** Update a product

**Request Body:**

```json
{
  "name": "Updated Laptop",
  "description": "Updated description",
  "price": 1099.99,
  "is_active": true
}
```

**Response:** 200 OK

```json
{
  "id": "uuid",
  "user_id": "vendor-uuid",
  "name": "Updated Laptop",
  "description": "Updated description",
  "price": 1099.99,
  "is_active": true,
  "created_at": "2025-01-02T10:00:00Z",
  "updated_at": "2025-01-02T11:00:00Z"
}
```

**cURL:**

```bash
curl -X PUT "http://localhost:8080/products?id=product-uuid" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "Updated Laptop",
    "price": 1099.99
  }'
```

---

#### DELETE /products?id={productId}

**Authentication:** Required (JWT Token)

**Description:** Delete a product

**Response:** 200 OK

```json
{
  "message": "product deleted successfully"
}
```

**cURL:**

```bash
curl -X DELETE "http://localhost:8080/products?id=product-uuid" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

#### PUT /products/status?id={productId}

**Authentication:** Required (JWT Token)

**Description:** Activate/Deactivate a product

**Request Body:**

```json
{
  "is_active": false
}
```

**Response:** 200 OK

```json
{
  "id": "uuid",
  "user_id": "vendor-uuid",
  "name": "Laptop",
  "description": "High-performance laptop",
  "price": 999.99,
  "is_active": false,
  "created_at": "2025-01-02T10:00:00Z",
  "updated_at": "2025-01-02T12:00:00Z"
}
```

**cURL:**

```bash
curl -X PUT "http://localhost:8080/products/status?id=product-uuid" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "is_active": false
  }'
```

---

#### GET /products/my

**Authentication:** Required (JWT Token)

**Description:** Get all products for authenticated vendor

**Response:** 200 OK

```json
[
  {
    "id": "uuid",
    "user_id": "vendor-uuid",
    "name": "Laptop",
    "description": "High-performance laptop",
    "price": 999.99,
    "is_active": true,
    "created_at": "2025-01-02T10:00:00Z",
    "updated_at": "2025-01-02T10:00:00Z"
  }
]
```

**cURL:**

```bash
curl http://localhost:8080/products/my \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## 4. VENDOR ROUTES (Public)

#### GET /vendors/{id}/products

**Authentication:** Not Required

**Description:** Get all products from a vendor

**Path Parameters:**

- `id` (required): Vendor/User UUID

**Response:** 200 OK

```json
[
  {
    "id": "uuid",
    "user_id": "vendor-uuid",
    "name": "Laptop",
    "description": "High-performance laptop",
    "price": 999.99,
    "is_active": true,
    "created_at": "2025-01-02T10:00:00Z",
    "updated_at": "2025-01-02T10:00:00Z"
  }
]
```

**cURL:**

```bash
curl http://localhost:8080/vendors/vendor-uuid/products
```

---

#### GET /vendors/{id}/products/active

**Authentication:** Not Required

**Description:** Get active products from a vendor

**Response:** 200 OK

```json
[
  {
    "id": "uuid",
    "user_id": "vendor-uuid",
    "name": "Laptop",
    "description": "High-performance laptop",
    "price": 999.99,
    "is_active": true,
    "created_at": "2025-01-02T10:00:00Z",
    "updated_at": "2025-01-02T10:00:00Z"
  }
]
```

**cURL:**

```bash
curl http://localhost:8080/vendors/vendor-uuid/products/active
```

---

## 5. USER ROUTES (Protected)

#### GET /me

**Authentication:** Required (JWT Token)

**Description:** Get authenticated user profile

**Response:** 200 OK

```json
{
  "id": "uuid",
  "email": "user@example.com",
  "first_name": "John",
  "last_name": "Doe",
  "role": "vendor",
  "is_approved": true,
  "created_at": "2025-01-01T10:00:00Z",
  "updated_at": "2025-01-02T10:00:00Z"
}
```

**cURL:**

```bash
curl http://localhost:8080/me \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## 6. ADMIN ROUTES (Protected + Admin Only)

#### GET /admin/vendors/pending

**Authentication:** Required (JWT Token)
**Role:** Admin Only

**Description:** Get all pending vendors

**Response:** 200 OK

```json
[
  {
    "id": "uuid",
    "email": "vendor@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "role": "vendor",
    "is_approved": false,
    "created_at": "2025-01-01T10:00:00Z"
  }
]
```

**cURL:**

```bash
curl http://localhost:8080/admin/vendors/pending \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

---

#### GET /admin/vendors/approved

**Authentication:** Required (JWT Token)
**Role:** Admin Only

**Description:** Get all approved vendors

**Response:** 200 OK

```json
[
  {
    "id": "uuid",
    "email": "vendor@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "role": "vendor",
    "is_approved": true,
    "created_at": "2025-01-01T10:00:00Z",
    "updated_at": "2025-01-02T10:00:00Z"
  }
]
```

**cURL:**

```bash
curl http://localhost:8080/admin/vendors/approved \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

---

#### POST /admin/vendors/{id}/approve

**Authentication:** Required (JWT Token)
**Role:** Admin Only

**Description:** Approve a vendor

**Path Parameters:**

- `id` (required): Vendor UUID

**Response:** 200 OK

```json
{
  "data": {
    "id": "uuid",
    "email": "vendor@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "role": "vendor",
    "is_approved": true
  },
  "message": "Vendor approved successfully"
}
```

**cURL:**

```bash
curl -X POST http://localhost:8080/admin/vendors/vendor-uuid/approve \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

---

## Route Summary Table

| Method | Endpoint                        | Auth | Role   | Description                |
| ------ | ------------------------------- | ---- | ------ | -------------------------- |
| GET    | `/health`                       | ✗    | -      | Health check               |
| POST   | `/auth/signup`                  | ✗    | -      | Register user              |
| POST   | `/auth/login`                   | ✗    | -      | Login user                 |
| GET    | `/products/active`              | ✗    | -      | Get all active products    |
| GET    | `/products/search`              | ✗    | -      | Search products            |
| GET    | `/products/price`               | ✗    | -      | Filter by price            |
| GET    | `/products?id={id}`             | ✗    | -      | Get single product         |
| POST   | `/products`                     | ✓    | vendor | Create product             |
| PUT    | `/products?id={id}`             | ✓    | vendor | Update product             |
| DELETE | `/products?id={id}`             | ✓    | vendor | Delete product             |
| PUT    | `/products/status?id={id}`      | ✓    | vendor | Toggle status              |
| GET    | `/products/my`                  | ✓    | vendor | Get my products            |
| GET    | `/vendors/{id}/products`        | ✗    | -      | Get vendor products        |
| GET    | `/vendors/{id}/products/active` | ✗    | -      | Get vendor active products |
| GET    | `/me`                           | ✓    | -      | Get profile                |
| GET    | `/admin/vendors/pending`        | ✓    | admin  | List pending vendors       |
| GET    | `/admin/vendors/approved`       | ✓    | admin  | List approved vendors      |
| POST   | `/admin/vendors/{id}/approve`   | ✓    | admin  | Approve vendor             |

---

## Error Responses

### 400 Bad Request

```json
{
  "error": "Invalid request body"
}
```

### 401 Unauthorized

```json
{
  "error": "Unauthorized"
}
```

### 403 Forbidden

```json
{
  "error": "Only vendors can create products"
}
```

### 404 Not Found

```json
{
  "error": "Product not found"
}
```

### 409 Conflict

```json
{
  "error": "Email already exists"
}
```

### 500 Internal Server Error

```json
{
  "error": "Internal server error"
}
```

---

## Middleware

All routes use the following global middleware:

- `RequestID`: Adds unique request ID
- `RealIP`: Extracts real client IP
- `Logger`: Logs all requests
- `Recoverer`: Recovers from panics
- `Timeout`: 15-second timeout for all requests

Protected routes additionally use:

- `JWTAuth`: Validates JWT token

Admin routes additionally use:

- `AdminOnly`: Checks admin role
