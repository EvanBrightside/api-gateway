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

- **401 Unauthorized**: Authentication error if the token is invalid or missing.

**Example response**:

```json
{
  "error": "Invalid token"
}
```

- **502 Bad Gateway**: Error proxying the request to the target service.

**Example response**:

```json
{
  "error": "Bad Gateway"
}
```

## Monitoring

The API Gateway exposes metrics for Prometheus and can be visualized in Grafana.

### GET /metrics

**Description**: Returns Prometheus metrics.

- **URL**: `/metrics/`
- **Method**: `GET`

### Example Request using curl

```bash
curl http://localhost:8080/metrics
```

**Response**: Returns metrics in Prometheus format.

## Prometheus and Grafana Setup

### Prometheus

Prometheus is set up to scrape the `/metrics` endpoint from the API Gateway.

#### Access Prometheus:
Once the containers are running, you can access the Prometheus dashboard by visiting:

```text
http://localhost:9090
```

#### Verify Scraping

You can verify that Prometheus is correctly scraping the metrics from the API Gateway by searching for the metric `api_gateway_requests_total` on the Prometheus UI.

### Grafana

Grafana can be used to visualize the metrics scraped by Prometheus. Here's how to set it up.

#### Access Grafana

Once the containers are running, you can access the Grafana dashboard by visiting:

```text
http://localhost:3000
```

#### Login to Grafana:
Default credentials for Grafana:

- **Username**: `admin`
- **Password**: `admin`

You'll be prompted to change the password on the first login.

### Add Prometheus as a Data Source

1. Go to the **Grafana Dashboard** → **Configuration** → **Data Sources**.
2. Click **Add data source**.
3. Select **Prometheus**.
4. Set the **URL** to `http://prometheus:9090` (or `http://localhost:9090` if you're running Prometheus locally).
5. Click **Save & Test** to verify the connection.

### Create a Dashboard

1. After adding Prometheus as a data source, go to **Create** → **Dashboard**.
2. Add a new **Graph** panel.
3. In the **Metrics** tab, add your Prometheus metric, for example:

    ```text
    api_gateway_requests_total
    ```

4. You can now visualize the total number of requests processed by the API Gateway.

## Possible Errors

- **502 Bad Gateway**: Occurs if the target service is unavailable or there is an error proxying the request.

- **401 Unauthorized**: Authentication error, occurs when the JWT token is missing or invalid.

- **400 Bad Request**: Request error, such as when the request body is malformed.

## Logs

All logs for the API Gateway are available in real time through Docker:

```bash
docker-compose logs -f
```

## Notes

- Ensure that services like settings-service and callback-router are running correctly, as the API Gateway relies on them.

- You can configure environment variables in the docker-compose.yml file if you need to specify custom parameters for each service.
