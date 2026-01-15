package services

import (
	"github.com/Mahaveer86619/FrameSense/pkg/db"
	"github.com/Mahaveer86619/FrameSense/pkg/models"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService() *UserService {
	return &UserService{
		db: db.GetDB(),
	}
}

func (r *UserService) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserService) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserService) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserService) Exists(email, username string) bool {
	var count int64
	r.db.Model(&models.User{}).Where("email = ? OR username = ?", email, username).Count(&count)
	return count > 0
}
