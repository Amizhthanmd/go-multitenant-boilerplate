package services

import (
	"fmt"
	"go-multitenant-boilerplate/models/tenant"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserService struct {
	db     *gorm.DB
	logger *zap.Logger
}

func InitializeUserService(db *gorm.DB, log *zap.Logger) *UserService {
	return &UserService{db: db, logger: log}
}

// Users
func (us *UserService) Create(data *tenant.User, schema string) error {
	table := us.GetSchemaTable(schema, data)
	return us.db.Table(table).Create(&data).Error
}

func (us *UserService) ReadByEmail(data *tenant.User, email string) error {
	return us.db.Where("email = ?", email).First(data).Error
}

func (us *UserService) GetUserById(data *tenant.User, id string, schema string) error {
	table := us.GetSchemaTable(schema, data)
	return us.db.Table(table).Where("id = ?", id).First(data).Error
}

func (us *UserService) ListUsers(data *[]tenant.User, limit, offset int, schema string) error {
	table := us.GetSchemaTable(schema, data)
	return us.db.Table(table).Order("created_at ASC").Limit(limit).Offset(offset).Find(data).Error
}

func (us *UserService) UpdateUser(data *tenant.User, schema string) error {
	table := us.GetSchemaTable(schema, data)
	var existingUser tenant.User
	err := us.db.Table(table).Where("id = ?", data.ID).First(&existingUser).Error
	if err != nil {
		return fmt.Errorf("user not found: %v", err)
	}
	return us.db.Table(table).Where("id = ?", data.ID).Updates(data).Error
}

func (us *UserService) DeleteUser(id string, schema string) error {
	table := us.GetSchemaTable(schema, tenant.User{})
	return us.db.Table(table).Where("id = ?", id).Delete(&tenant.User{}).Error
}

// Permissions
func (us *UserService) GetPermissions(data *[]tenant.Permission, schema string) error {
	table := us.GetSchemaTable(schema, data)
	return us.db.Table(table).Find(&data).Error
}

func (us *UserService) ListPermissions(data *[]tenant.Permission, schema string) error {
	table := us.GetSchemaTable(schema, data)
	return us.db.Table(table).Find(data).Error
}

func (us *UserService) GetPermissionsByIds(data *[]tenant.Permission, Ids []string) error {
	return us.db.Where("id IN ?", Ids).Find(&data).Error
}

// Roles
func (us *UserService) CreateRoles(data *tenant.Role) error {
	return us.db.Create(&data).Error
}

func (us *UserService) GetRolesById(data *tenant.Role, id string) error {
	return us.db.Preload("Permissions").Where("id = ?", id).First(&data).Error
}

func (us *UserService) GetSchemaTable(schema, data interface{}) string {
	stmt := &gorm.Statement{DB: us.db}
	stmt.Parse(data)
	return fmt.Sprintf(`"%s".%s`, schema, stmt.Table)
}
