# Share-A-Ride API Design Document

## Overview

This document outlines the technical design for Share-A-Ride's REST APIs, focusing on Authentication and Driver Management.

## API Standards

### Base URL

```http
https://api.share-a-ride.com/v1
```

### Response Format

All API responses follow this standard format:

```json
{
    "success": true,
    "data": {},
    "error": null,
    "metadata": {
        "timestamp": "2024-02-20T12:00:00Z"
    }
}
```

### Error Format

```json
{
    "success": false,
    "data": null,
    "error": {
        "code": "AUTH001",
        "message": "Invalid credentials",
        "details": {}
    },
    "metadata": {
        "timestamp": "2024-02-20T12:00:00Z"
    }
}
```

### Authentication

- All authenticated endpoints require a Bearer token in the Authorization header
- Format: `Authorization: Bearer <jwt_token>`

## 1. Authentication APIs

### 1.1 Register User

```http
POST /auth/register
```

Request Body:

```json
{
    "name": "string",
    "email": "string",
    "phone": "string",
    "password": "string",
    "user_type": "rider|driver"
}
```

Response (201 Created):

```json
{
    "success": true,
    "data": {
        "user": {
            "id": "uuid",
            "name": "string",
            "email": "string",
            "phone": "string",
            "user_type": "rider|driver",
            "created_at": "timestamp"
        },
        "tokens": {
            "access_token": "string",
            "refresh_token": "string"
        }
    }
}
```

Validation Rules:

- Email must be unique and valid format
- Phone must be unique and valid format
- Password must be at least 8 characters
- user_type must be either "rider" or "driver"

### 1.2 Login

```http
POST /auth/login
```

Request Body:

```json
{
    "email": "string",
    "password": "string"
}
```

Response (200 OK):

```json
{
    "success": true,
    "data": {
        "user": {
            "id": "uuid",
            "name": "string",
            "email": "string",
            "phone": "string",
            "user_type": "rider|driver"
        },
        "tokens": {
            "access_token": "string",
            "refresh_token": "string"
        }
    }
}
```

### 1.3 Refresh Token

```http
POST /auth/refresh
```

Request Body:

```json
{
    "refresh_token": "string"
}
```

Response (200 OK):

```json
{
    "success": true,
    "data": {
        "tokens": {
            "access_token": "string",
            "refresh_token": "string"
        }
    }
}
```

## 2. Driver Management APIs

### 2.1 Submit Driver Verification

```http
POST /drivers/verify
Authorization: Bearer <token>
```

Request Body:

```json
{
    "license_number": "string",
    "vehicle": {
        "type": "car|bike",
        "model": "string",
        "plate_number": "string"
    },
    "documents": [
        {
            "type": "license|registration|insurance",
            "file_url": "string"
        }
    ]
}
```

Response (202 Accepted):

```json
{
    "success": true,
    "data": {
        "verification_id": "uuid",
        "status": "pending",
        "submitted_at": "timestamp"
    }
}
```

### 2.2 Update Driver Status

```http
PUT /drivers/status
Authorization: Bearer <token>
```

Request Body:

```json
{
    "is_available": boolean,
    "current_location": {
        "latitude": number,
        "longitude": number
    }
}
```

Response (200 OK):

```json
{
    "success": true,
    "data": {
        "driver_id": "uuid",
        "is_available": boolean,
        "current_location": {
            "latitude": number,
            "longitude": number
        },
        "updated_at": "timestamp"
    }
}
```

### 2.3 Get Driver's Ride History

```http
GET /drivers/rides
Authorization: Bearer <token>
```

Query Parameters:

- `status`: optional (completed|cancelled|ongoing)
- `from_date`: optional (ISO date)
- `to_date`: optional (ISO date)
- `page`: optional (default: 1)
- `limit`: optional (default: 10)

Response (200 OK):

```json
{
    "success": true,
    "data": {
        "rides": [
            {
                "id": "uuid",
                "rider": {
                    "id": "uuid",
                    "name": "string"
                },
                "pickup_location": {
                    "latitude": number,
                    "longitude": number
                },
                "dropoff_location": {
                    "latitude": number,
                    "longitude": number
                },
                "status": "string",
                "fare": number,
                "created_at": "timestamp",
                "completed_at": "timestamp"
            }
        ]
    },
    "metadata": {
        "total": number,
        "page": number,
        "limit": number,
        "has_more": boolean
    }
}
```

## Data Models

### User

```go
type User struct {
    ID        string    `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    Phone     string    `json:"phone"`
    Password  string    `json:"-"`
    UserType  string    `json:"user_type"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### Driver

```go
type Driver struct {
    ID              string    `json:"id"`
    UserID          string    `json:"user_id"`
    LicenseNumber   string    `json:"license_number"`
    Vehicle         Vehicle   `json:"vehicle"`
    IsVerified      bool      `json:"is_verified"`
    IsAvailable     bool      `json:"is_available"`
    CurrentLocation Location  `json:"current_location"`
    Documents       []Document `json:"documents"`
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
}
```

### Vehicle

```go
type Vehicle struct {
    Type        string `json:"type"`
    Model       string `json:"model"`
    PlateNumber string `json:"plate_number"`
}
```

### Document

```go
type Document struct {
    Type    string `json:"type"`
    FileURL string `json:"file_url"`
}
```

### Location

```go
type Location struct {
    Latitude  float64 `json:"latitude"`
    Longitude float64 `json:"longitude"`
}
```

## Error Codes

### Authentication Errors

- AUTH001: Invalid credentials
- AUTH002: Token expired
- AUTH003: Invalid token
- AUTH004: User not found
- AUTH005: Email already exists
- AUTH006: Phone already exists

### Driver Management Errors

- DRV001: Driver not found
- DRV002: Invalid vehicle type
- DRV003: Invalid document type
- DRV004: Missing required documents
- DRV005: Driver not verified
- DRV006: Invalid location coordinates

## Security Considerations

1. **Password Storage**
   - Passwords must be hashed using bcrypt
   - Minimum password length: 8 characters

2. **JWT Token**
   - Access token expiry: 1 hour
   - Refresh token expiry: 7 days
   - Token must include user ID and role

3. **Rate Limiting**
   - Login attempts: 5 per minute per IP
   - API calls: 100 per minute per user

4. **Input Validation**
   - All input must be validated and sanitized
   - File uploads limited to 5MB per file
   - Supported document formats: PDF, JPG, PNG

## Database Schema

### users

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(20) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    user_type VARCHAR(10) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
```

### drivers

```sql
CREATE TABLE drivers (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    license_number VARCHAR(50) UNIQUE NOT NULL,
    vehicle_type VARCHAR(20) NOT NULL,
    vehicle_model VARCHAR(100) NOT NULL,
    vehicle_plate VARCHAR(20) NOT NULL,
    is_verified BOOLEAN DEFAULT FALSE,
    is_available BOOLEAN DEFAULT FALSE,
    current_latitude DECIMAL(10,8),
    current_longitude DECIMAL(11,8),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
```

### driver_documents

```sql
CREATE TABLE driver_documents (
    id UUID PRIMARY KEY,
    driver_id UUID REFERENCES drivers(id),
    document_type VARCHAR(20) NOT NULL,
    file_url VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL
);
```
