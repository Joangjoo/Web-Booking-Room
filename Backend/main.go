package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginPayload struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type Room struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Capacity    int    `json:"capacity"`
}

type BookingInput struct {
	RoomID    int       `json:"room_id" binding:"required"`
	StartTime time.Time `json:"start_time" binding:"required"`
	EndTime   time.Time `json:"end_time" binding:"required"`
}

var jwtSecret = []byte("JoSecretKey")
var db *sql.DB

func main() {
	connStr := "postgres://app_user:app_password@localhost:5432/booking_db?sslmode=disable"
	var err error
	db, err = sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("Gagal membuka koneksi database: %v", err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}
	fmt.Println("🎉 Berhasil terhubung ke database booking_db!")

	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173"} 
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	router.Use(cors.New(config))

	api := router.Group("/api")
	{
		api.POST("/register", registerUser)
		api.POST("/login", loginUser)
		api.GET("/rooms", getAllRooms)
		api.GET("/rooms/:id", getRoomByID)
	}

	protected := router.Group("/api/protected")
	protected.Use(authMiddleware())
	{
		protected.GET("/test", func(c *gin.Context) {
			userID, exists := c.Get("userID")
			if !exists {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendapatkan user ID dari context"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Halo user %v! Anda berhasil mengakses rute terproteksi.", userID)})
		})
		protected.POST("/bookings", createBooking)

		// protected.POST("/bookings", createBooking)
		// protected.GET("/rooms", getAllRooms)
	}

	fmt.Println("🚀 Server berjalan di http://localhost:8080")
	router.Run(":8080")
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Header otorisasi tidak ditemukan"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Format token tidak valid. Gunakan format 'Bearer <token>'"})
			return
		}
		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Metode signing tidak terduga: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid: " + err.Error()})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID := claims["user_id"]
			c.Set("userID", userID)
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid atau claims tidak ditemukan"})
		}
	}
}

func registerUser(c *gin.Context) {
	var newUser User

	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengenkripsi password"})
		return
	}

	insertSQL := `INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3) RETURNING id`
	var userID int
	err = db.QueryRow(insertSQL, newUser.Name, newUser.Email, string(hashedPassword)).Scan(&userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat pengguna baru, email mungkin sudah terdaftar."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Pengguna berhasil terdaftar!",
		"user_id": userID,
	})
}

func loginUser(c *gin.Context) {
	var payload LoginPayload
	var userName string
	var hashedPassword string
	var userID int

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	querySQL := `SELECT id, name, password_hash FROM users WHERE email = $1`
	err := db.QueryRow(querySQL, payload.Email).Scan(&userID, &userName, &hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau password salah"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Terjadi kesalahan server"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(payload.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau password salah"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"name":    userName,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil!",
		"token":   tokenString,
	})
}

func getAllRooms(c *gin.Context) {
	querySQL := `SELECT id, name, description, capacity FROM rooms`
	rows, err := db.Query(querySQL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data ruangan"})
		return
	}
	defer rows.Close()

	var rooms []Room

	for rows.Next() {
		var room Room
		err := rows.Scan(&room.ID, &room.Name, &room.Description, &room.Capacity)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error setelah iterasi data ruangan"})
		return
	}

	c.JSON(http.StatusOK, rooms)
}

func getRoomByID(c *gin.Context) {
	// 1. Ambil ID dari parameter URL
	roomID := c.Param("id")

	var room Room
	querySQL := `SELECT id, name, description, capacity FROM rooms WHERE id = $1`

	err := db.QueryRow(querySQL, roomID).Scan(&room.ID, &room.Name, &room.Description, &room.Capacity)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Ruangan tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data ruangan"})
		return
	}

	c.JSON(http.StatusOK, room)
}

func createBooking(c *gin.Context) {
	var input BookingInput

	// 1. Bind input JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid: " + err.Error()})
		return
	}

	// 2. Validasi waktu dasar (EndTime harus setelah StartTime)
	if input.EndTime.Before(input.StartTime) || input.EndTime.Equal(input.StartTime) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Waktu selesai harus setelah waktu mulai"})
		return
	}

	// 3. Cek konflik jadwal (Logika Inti)
	var conflictCount int
	conflictQuery := `
        SELECT COUNT(*) FROM bookings
        WHERE room_id = $1 AND status = 'confirmed' AND (
            (start_time < $2 AND end_time > $3)
        )`
	// $1 = room_id, $2 = input.EndTime, $3 = input.StartTime
	err := db.QueryRow(conflictQuery, input.RoomID, input.EndTime, input.StartTime).Scan(&conflictCount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memeriksa ketersediaan ruangan"})
		return
	}

	if conflictCount > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Ruangan tidak tersedia pada jadwal yang dipilih (konflik)"})
		return
	}

	// 4. Ambil User ID dari middleware context
	userIDClaim, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendapatkan ID pengguna dari token"})
		return
	}
	// Konversi userID dari interface{} ke tipe data yang sesuai (misalnya float64 lalu int)
	userID := int(userIDClaim.(float64))

	// 5. Masukkan data booking baru ke database
	insertSQL := `INSERT INTO bookings (user_id, room_id, start_time, end_time, status) VALUES ($1, $2, $3, $4, 'confirmed') RETURNING id`
	var bookingID int
	err = db.QueryRow(insertSQL, userID, input.RoomID, input.StartTime, input.EndTime).Scan(&bookingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data booking"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Booking berhasil dibuat!",
		"booking_id": bookingID,
	})
}
