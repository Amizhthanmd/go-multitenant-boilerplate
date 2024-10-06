package tenant

import "time"

type Tenant struct {
	ID           string    `json:"id,omitempty" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Organization string    `json:"organization"`
	Dsn          string    `json:"dsn"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	CreatedAt    time.Time `json:"created_at,omitempty" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at,omitempty" gorm:"autoUpdateTime"`
}
