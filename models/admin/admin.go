package admin

type Admin struct {
	ID        int    `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at,omitempty" gorm:"autoCreateTime"`
	UpdatedAt string `json:"updated_at,omitempty" gorm:"autoUpdateTime"`
}
