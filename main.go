package main // <-- PENTING: package main

import (
	"log"
	"net/http"
	"os"

	"filoti-backend/config" // Import dari root module
	"filoti-backend/routes" // Import dari root module
)

func main() {
	// 1. Inisialisasi DB
	//    config.InitDB() akan otomatis memuat .env berkat modifikasi kita
	config.InitDB()

	// 2. Setup router
	r := routes.SetupRouter()

	// 3. Jalankan server lokal
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on http://localhost:%s (for local development)", port)

	// Kita gunakan server http standar Go yang di-wrap oleh Gin
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Failed to run server locally: %v", err)
	}
}
