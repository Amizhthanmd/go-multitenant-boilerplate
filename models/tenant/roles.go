package tenant

type (
	Role struct {
		ID          string       `json:"id,omitempty" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
		Name        string       `json:"name"`
		Permissions []Permission `json:"permissions" gorm:"many2many:role_permissions;"`
		Users       []User       `json:"users" gorm:"foreignKey:RoleID"`
	}
	AddRoles struct {
		Name        string   `json:"name" binding:"required"`
		Permissions []string `json:"permissions" binding:"required"`
	}
)
