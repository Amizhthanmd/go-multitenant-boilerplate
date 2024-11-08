package services

import (
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

func (us *UserService) Create(data *tenant.User) error {
	return us.db.Create(&data).Error
}

func (us *UserService) ReadByEmail(data *tenant.User, email string) error {
	return us.db.First(data, "email").Error
}
