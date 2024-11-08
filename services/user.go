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

func (us *UserService) Create(data *tenant.User, schema string) error {
	table := us.GetSchemaTable(schema, data)
	return us.db.Table(table).Create(&data).Error
}

func (us *UserService) ReadByEmail(data *tenant.User, email string, schema string) error {
	table := us.GetSchemaTable(schema, data)
	return us.db.Table(table).Where("email = ?", email).First(data).Error
}

func (us *UserService) GetSchemaTable(schema, data interface{}) string {
	stmt := &gorm.Statement{DB: us.db}
	stmt.Parse(data)
	return fmt.Sprintf(`"%s".%s`, schema, stmt.Table)
}
