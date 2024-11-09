package tenant

type Permission struct {
	ID    string `json:"id,omitempty" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name  string `json:"name"`
	Roles []*Role `json:"roles" gorm:"many2many:role_permissions;"`
}
