package tenant

import "time"

type User struct {
	ID           string    `json:"id,omitempty" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Organization string    `json:"organization"`
	RoleID       string    `json:"role_id" gorm:"type:uuid"` 
	Role         Role      `json:"role" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	CreatedAt    time.Time `json:"created_at,omitempty" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at,omitempty" gorm:"autoUpdateTime"`
}

type UserLogin struct {
	Organization string `json:"organization" binding:"required"`
	Email        string `json:"email" binding:"required"`
	Password     string `json:"password" binding:"required"`
}
