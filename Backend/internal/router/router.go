package router

import (
	"database/sql"
	"example.com/Backend/internal/handlers" 

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(db *sql.DB) *gin.Engine {
	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173"} 
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	router.Use(cors.New(config))

	h := &handlers.Handler{DB: db}

	api := router.Group("/api")
	{
		api.POST("/register", h.RegisterUser)
		api.POST("/login", h.LoginUser)
		api.GET("/rooms", h.GetAllRooms)
		api.GET("/rooms/:id", h.GetRoomByID)
	}

	// Grup untuk rute terproteksi (perlu token JWT)
	protected := router.Group("/api/protected")
	protected.Use(h.AuthMiddleware()) // Middleware juga dipanggil sebagai method
	{
		protected.POST("/bookings", h.CreateBooking)
		// Tambahkan rute terproteksi lainnya di sini jika ada
	}

	return router
}