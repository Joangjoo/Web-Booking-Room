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

	protected := router.Group("/api/protected")
	protected.Use(h.AuthMiddleware())
	{
		protected.POST("/bookings", h.CreateBooking)
		protected.POST("/rooms", h.CreateRoom)
		protected.PUT("/rooms/:id", h.UpdateRoom)
		protected.DELETE("/rooms/:id", h.DeleteRoom)
		protected.GET("/users/me", h.MeHandler)
		protected.GET("/bookings", h.GetAllBookings)
	}

	return router
}
