package main

import (
	"database/sql"
	"log"
	"os" // <-- Tambahkan import os

	"example.com/Backend/internal/router"
	"github.com/joho/godotenv" // <-- Tambahkan import godotenv
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: Gagal memuat file .env. Menggunakan environment variable sistem.")
	}

	connStr := os.Getenv("DB_SOURCE")
	if connStr == "" {
		log.Fatal("DB_SOURCE environment variable tidak ditemukan")
	}

	db, err := sql.Open("pgx", connStr)
	log.Println("🎉 Berhasil terhubung ke database!")

	port := os.Getenv("API_PORT")
	if port == "" {
		port = ":8080" 
	}
	
	r := router.SetupRouter(db)

	log.Printf("🚀 Server berjalan di http://localhost%s\n", port)
	if err := r.Run(port); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}