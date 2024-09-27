# API Gateway

An API Gateway for routing and proxying requests to various internal services, such as `settings-service` and `callback-router`, while handling JWT authentication and exposing metrics for monitoring.

## Requirements

To successfully deploy and run this project, you need the following:

- **Docker**: Installed Docker version 19.03 or higher
- **Docker Compose**: Installed Docker Compose version 1.27 or higher
- **Postman/Insomnia** or any other HTTP client for testing the API

## Deploying the Project using Docker Compose

1. **Clone the repository:**

    ```bash
    git clone https://github.com/EvanBrightside/api-gateway.git
    cd api-gateway
    ```

2. **Start the containers using Docker Compose:**

    Run the following command to build and start all containers:

    ```bash
    docker-compose build
    docker-compose up
    ```

3. **Verify that all services are running:**

    Make sure all containers have started successfully. You should see logs for the API Gateway and proxied services like `settings-service` and `callback-router`.

## API Documentation

### Authentication

To access protected routes, JWT authentication is required. You can obtain a token by sending a request to the `/auth` route.

#### POST /auth

**Description**: Authenticate and receive a JWT token.

**URL**: `/auth/`

**Method**: `POST`

**Request Body** (JSON):

```json
{
  "username": "admin"
}
```

**Response Body** (JSON):

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_at_unix": 1727348342,
  "expires_at_utc": "2024-09-26 10:59:02 UTC"
}
```

**Example Request using curl:**

```bash
curl -X POST http://localhost:8080/auth -d '{"username": "admin"}' -H "Content-Type: application/json"
```

### Protected Routes

#### POST /api/settings/

**Description**: Proxies a request to the settings-service.

**URL**: `/api/settings/`

**Method**: `POST`

**Authentication Required**: Yes (Bearer JWT token)

**Request Body** (JSON):

```json
{
  "setting": "value"
}
```

**Response Body** (JSON):

```json
{
  "message": "Settings Service"
}
```

**Example Request using curl**:

```bash
curl -X POST http://localhost:8080/api/settings/ -d '{"setting": "value"}' -H "Authorization: Bearer <your-token>" -H "Content-Type: application/json"
```

#### POST /api/callback-router/

**Description**: Proxies a request to the callback-router service.

**URL**: `/api/callback-router/`

**Method**: `POST`

**Authentication Required**: Yes (Bearer JWT token)

**Request Body** (JSON):

```json
{
  "callback": "value"
}
```

**Response BODY** (JSON):

```json
{
  "message": "Callback Router Service"
}
```

**Example Request using curl**:

```bash
curl -X POST http://localhost:8080/api/callback-router/ -d '{"callback": "value"}' -H "Authorization: Bearer <your-token>" -H "Content-Type: application/json"
```

## Possible Errors

### 401 Unauthorized

Authentication error if the token is invalid or missing.

**Example response**:

```json
{
  "error": "Invalid token"
}
```

### 502 Bad Gateway

Error proxying the request to the target service.

**Example response**:

```json
{
  "error": "Bad Gateway"
}
```

