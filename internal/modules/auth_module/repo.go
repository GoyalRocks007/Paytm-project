package authmodule

import (
	"log"
	"paytm-project/internal/models"

	"gorm.io/gorm"
)

type IAuthRepo interface {
	CreateUser(user *User) error
	GetUserByEmail(email string) (*User, error)
	GetUserByEmailTransactional(tx *gorm.DB, email string) (*User, error)
	UpdateWalletBalance(tx *gorm.DB, userId string, delta float64) error
	UpdateUserRole(user *User) error
}

type AuthRepo struct {
	models.BaseRepo
}

func (a *AuthRepo) CreateUser(user *User) error {
	err := a.Db.Create(user).Error
	if err != nil {
		log.Println("failed to create user:", err)
		return err
	}
	log.Println("User created with ID:", user.Id)
	return nil
}

func (a *AuthRepo) GetUserByEmail(email string) (*User, error) {
	var user User
	err := a.Db.Preload("Wallet").Where("email = ?", email).First(&user).Error
	if err != nil {
		log.Println("couldn't find user with given email:", err.Error())
		return nil, err
	}
	return &user, nil
}

func (a *AuthRepo) GetUserByEmailTransactional(tx *gorm.DB, email string) (*User, error) {
	var user User
	err := tx.Preload("Wallet").Where("email = ?", email).First(&user).Error
	if err != nil {
		log.Println("couldn't find user with given email:", err.Error())
		return nil, err
	}
	return &user, nil
}

func (a *AuthRepo) UpdateWalletBalance(tx *gorm.DB, userId string, delta float64) error {
	return tx.Model(&Wallet{}).
		Where("user_id = ?", userId).
		Update("balance", gorm.Expr("balance + ?", delta)).Error
}

func (a *AuthRepo) UpdateUserRole(user *User) error {
	return a.Db.
		Model(&User{}).
		Where("email = ?", user.Email).
		Update("role", user.Role).Error
}
