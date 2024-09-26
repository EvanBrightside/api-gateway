package main

import (
	"bytes"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

// JWT ключ для подписи токенов
var jwtKey = []byte("my_secret_key")

// Структура для JWT Claims
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// Создаем Prometheus метрики
var (
	totalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_gateway_requests_total",
			Help: "Total number of requests processed",
		},
		[]string{"path", "method"},
	)
)

func init() {
	// Регистрация метрик в Prometheus
	prometheus.MustRegister(totalRequests)
}

func main() {
	router := gin.Default()

	// Middleware для логирования
	router.Use(prometheusMiddleware())

	// Маршруты
	router.GET("/metrics", gin.WrapH(promhttp.Handler())) // для Prometheus

	router.POST("/auth", auth)

	authorized := router.Group("/api")
	authorized.Use(authMiddleware())
	{
		// GET маршруты
		authorized.GET("/settings/*path", reverseProxy("http://settings:8081"))
		authorized.GET("/callback-router/*path", reverseProxy("http://callback-router:8082"))

		// POST маршруты
		authorized.POST("/settings/*path", proxyToServiceWithBody("http://settings:8081"))
		authorized.POST("/callback-router/*path", proxyToServiceWithBody("http://callback-router:8082"))
	}

	log.Fatal(router.Run(":8080"))
}

// Middleware для логирования метрик в Prometheus
func prometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		totalRequests.WithLabelValues(c.FullPath(), c.Request.Method).Inc()
		c.Next()
	}
}

// Функция обратного проксирования для GET запросов
func reverseProxy(target string) gin.HandlerFunc {
	targetURL, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	return func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

// Функция обратного проксирования для POST запросов с JSON
func proxyToServiceWithBody(target string) gin.HandlerFunc {
	targetURL, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	return func(c *gin.Context) {
		// Логируем тело и заголовки запроса
		log.Printf("Proxying request to %s with method %s", targetURL.String(), c.Request.Method)

		// Чтение тела запроса
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
			return
		}

		// Логируем тело запроса и Content-Length
		log.Printf("Request body: %s", string(body))
		log.Printf("Content-Length: %d", len(body))

		// Проксируем запрос на целевой сервис
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body)) // чтобы тело запроса было доступно для проксирования
		c.Request.ContentLength = int64(len(body))
		c.Request.Header.Set("Content-Type", "application/json")

		// Перенаправление запроса через прокси
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

// Middleware для проверки JWT
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "details": err.Error()})
			c.Abort()
			return
		}

		// Проверка срока действия токена
		if claims.ExpiresAt < time.Now().UTC().Unix() {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Функция для создания токена (пример логина)
func auth(c *gin.Context) {
	var jsonInput map[string]string
	if err := c.ShouldBindJSON(&jsonInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	username, exists := jsonInput["username"]
	if !exists || username != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	expirationTime := time.Now().UTC().Add(5 * time.Minute)
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(jwtKey)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create token"})
		return
	}

	// Преобразуем время истечения в формат UTC
	expirationFormatted := expirationTime.Format("2006-01-02 15:04:05 MST")

	c.JSON(http.StatusOK, gin.H{
		"token":           tokenStr,
		"expires_at_unix": claims.ExpiresAt,    // Время истечения токена в формате Unix
		"expires_at_utc":  expirationFormatted, // Время истечения токена в читаемом формате
	})
}
