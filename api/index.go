package handler // <-- PENTING: package handler

import (
	"net/http"
	"sync" // <-- Impor sync

	"filoti-backend/config" // Import dari root module
	"filoti-backend/routes" // Import dari root module

	"github.com/gin-gonic/gin"
)

// Gunakan sync.Once untuk inisialisasi yang aman untuk serverless
var (
	once   sync.Once
	router *gin.Engine
)

// initRouter akan menginisialisasi DB dan Router HANYA SEKALI
func initRouter() {
	// 1. Inisialisasi DB (InitDB akan memuat env var dari Vercel)
	config.InitDB()

	// 2. Setup router
	router = routes.SetupRouter()
}

// Handler adalah entry point Vercel.
// Ini adalah satu-satunya func yang dibutuhkan Vercel.
func Handler(w http.ResponseWriter, r *http.Request) {
	// 1. Panggil initRouter() menggunakan sync.Once
	//    Ini memastikan DB & Router hanya di-init satu kali,
	//    bahkan jika Vercel "membekukan" dan "mencairkan" fungsi ini.
	once.Do(initRouter)

	// 2. Sajikan permintaan menggunakan router Gin yang sudah diinisialisasi
	router.ServeHTTP(w, r)
}
