package models

import "time"
import "github.com/lib/pq"

type User struct {
	ID       int    `json:"id"`
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
	IsAvailable bool            `json:"is_available"` 
    Features    pq.StringArray  `json:"features"`
}

type BookingInput struct {
	RoomID    int       `json:"room_id" binding:"required"`
	StartTime time.Time `json:"start_time" binding:"required"`
	EndTime   time.Time `json:"end_time" binding:"required"`
	Purpose   string    `json:"purpose"`
}

type BookingDetails struct {
	ID        int       `json:"id"`
	RoomName  string    `json:"room_name"`
	UserName  string    `json:"user_name"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Purpose   string    `json:"purpose"`
	Status    string    `json:"status"`
}