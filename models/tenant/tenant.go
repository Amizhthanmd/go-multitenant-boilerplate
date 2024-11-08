package tenant

import "time"

type Tenant struct {
	ID           string    `json:"id,omitempty" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	FirstName    string    `json:"first_name" required:"true"`
	LastName     string    `json:"last_name" required:"true"`
	Organization string    `json:"organization" required:"true"`
	Dsn          string    `json:"dsn"`
	Email        string    `json:"email" required:"true"`
	Password     string    `json:"password" required:"true"`
	CreatedAt    time.Time `json:"created_at,omitempty" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at,omitempty" gorm:"autoUpdateTime"`
}
