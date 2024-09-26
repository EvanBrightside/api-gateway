package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// Структура ответа
type Response struct {
	Message string `json:"message"`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Создаем JSON-ответ
		response := Response{
			Message: "Settings Service",
		}

		// Устанавливаем заголовок Content-Type
		w.Header().Set("Content-Type", "application/json")

		// Кодируем структуру в JSON и отправляем в ответ
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
			return
		}
	})

	log.Fatal(http.ListenAndServe(":8081", nil))
}
