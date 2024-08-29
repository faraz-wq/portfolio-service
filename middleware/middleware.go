package middleware

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func APIKeyAuth(next http.HandlerFunc) http.HandlerFunc {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	API_KEY := os.Getenv("API_KEY")

	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")
		if apiKey != API_KEY {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}
