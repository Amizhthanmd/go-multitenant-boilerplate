package tenant

import "time"

type User struct {
	ID           string    `json:"id,omitempty" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	FirstName    string    `json:"first_name" binding:"required"`
	LastName     string    `json:"last_name" binding:"required"`
	Organization string    `json:"organization" binding:"required"`
	Role         string    `json:"role" binding:"required"`
	Email        string    `json:"email" binding:"required"`
	Password     string    `json:"password" binding:"required"`
	CreatedAt    time.Time `json:"created_at,omitempty" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at,omitempty" gorm:"autoUpdateTime"`
}

type UserLogin struct {
	Organization string `json:"organization" binding:"required"`
	Email        string `json:"email" binding:"required"`
	Password     string `json:"password" binding:"required"`
}
