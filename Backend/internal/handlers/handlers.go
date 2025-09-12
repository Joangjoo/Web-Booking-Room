package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"example.com/Backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	DB *sql.DB
}

func (h *Handler) AuthMiddleware() gin.HandlerFunc {
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

		jwtSecret := []byte(os.Getenv("JWT_SECRET_KEY"))
		if len(jwtSecret) == 0 {
			log.Println("Kunci JWT tidak terkonfigurasi di server")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Kesalahan konfigurasi server"})
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("metode signing tidak terduga: %v", token.Header["alg"])
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

func (h *Handler) RegisterUser(c *gin.Context) {
	var newUser models.User

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
	err = h.DB.QueryRow(insertSQL, newUser.Name, newUser.Email, string(hashedPassword)).Scan(&userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat pengguna baru, email mungkin sudah terdaftar."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Pengguna berhasil terdaftar!",
		"user_id": userID,
	})
}

func (h *Handler) LoginUser(c *gin.Context) {
	var payload models.LoginPayload
	var userName string
	var hashedPassword string
	var userID int

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	querySQL := `SELECT id, name, password_hash FROM users WHERE email = $1`
	err := h.DB.QueryRow(querySQL, payload.Email).Scan(&userID, &userName, &hashedPassword)
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

	jwtSecret := []byte(os.Getenv("JWT_SECRET_KEY"))
	if len(jwtSecret) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kunci JWT tidak terkonfigurasi di server"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"name":    userName,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token JWT"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil!",
		"token":   tokenString,
	})
}

func (h *Handler) GetAllRooms(c *gin.Context) {
	querySQL := `SELECT id, name, description, capacity, is_available, features FROM rooms`
	rows, err := h.DB.Query(querySQL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data ruangan"})
		return
	}
	defer rows.Close()

	var rooms []models.Room

	for rows.Next() {
		var room models.Room
		err := rows.Scan(&room.ID, &room.Name, &room.Description, &room.Capacity, &room.IsAvailable, &room.Features)
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

func (h *Handler) GetRoomByID(c *gin.Context) {
	roomID := c.Param("id")

	var room models.Room
	querySQL := `SELECT id, name, description, capacity, is_available, features FROM rooms WHERE id = $1`

	err := h.DB.QueryRow(querySQL, roomID).Scan(&room.ID, &room.Name, &room.Description, &room.Capacity, &room.IsAvailable, &room.Features)
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

func (h *Handler) CreateBooking(c *gin.Context) {
	var input models.BookingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid: " + err.Error()})
		return
	}

	if input.EndTime.Before(input.StartTime) || input.EndTime.Equal(input.StartTime) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Waktu selesai harus setelah waktu mulai"})
		return
	}

	var conflictCount int
	conflictQuery := `
        SELECT COUNT(*) FROM bookings
        WHERE room_id = $1 AND status = 'confirmed' AND (
            (start_time < $2 AND end_time > $3)
        )`
	err := h.DB.QueryRow(conflictQuery, input.RoomID, input.EndTime, input.StartTime).Scan(&conflictCount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memeriksa ketersediaan ruangan"})
		return
	}

	if conflictCount > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Ruangan tidak tersedia pada jadwal yang dipilih (konflik)"})
		return
	}

	userIDClaim, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendapatkan ID pengguna dari token"})
		return
	}
	userID := int(userIDClaim.(float64))

	insertSQL := `INSERT INTO bookings (user_id, room_id, start_time, end_time, status, purpose) VALUES ($1, $2, $3, $4, 'confirmed', $5) RETURNING id`
	var bookingID int
	err = h.DB.QueryRow(insertSQL, userID, input.RoomID, input.StartTime, input.EndTime, input.Purpose).Scan(&bookingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data booking"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Booking berhasil dibuat!",
		"booking_id": bookingID,
	})
}
