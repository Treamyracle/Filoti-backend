package handler

import (
	"log"
	"net/http" // Required for http.ResponseWriter, http.Request, and http.ListenAndServe for local testing
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"filoti-backend/config"
	"filoti-backend/routes"
)

var r *gin.Engine

func init() {
	// Load .env (only for local development)
	if os.Getenv("VERCEL_ENV") != "production" { // Check if not in Vercel production env
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found, using environment variables")
		}
	}

	// Initialize DB (will stay alive as long as Lambda container is warm)
	config.InitDB()

	r = routes.SetupRouter()
}

// Handler function for Vercel
// Vercel will call this function directly for incoming requests
func Handler(w http.ResponseWriter, req *http.Request) {
	// r is the global Gin engine initialized in init()
	r.ServeHTTP(w, req)
}

// main function for local development (optional, but good practice)
// This will be ignored by Vercel.
func main() {
	// This part is for local testing using `go run main.go`
	// It uses the same 'r' instance initialized in init()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on port %s (for local development)", port)
	// Use http.ListenAndServe with the global Gin router
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Failed to run server locally: %v", err)
	}
}
