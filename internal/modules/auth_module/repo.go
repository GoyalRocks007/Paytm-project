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
	GetUserIdByEmail(email string) (string, error)
	GetWalletWithLock(tx *gorm.DB, userId string) (*Wallet, error)
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

func (a *AuthRepo) GetUserIdByEmail(email string) (string, error) {
	var user User
	err := a.Db.Where("email = ?", email).Select("id").First(&user).Error
	if err != nil {
		log.Println("couldn't find user with given email:", err.Error())
		return "", err
	}
	return user.Id, nil
}

func (a *AuthRepo) GetUserByEmail(email string) (*User, error) {
	var user User
	err := a.Db.Where("email = ?", email).First(&user).Error
	if err != nil {
		log.Println("couldn't find user with given email:", err.Error())
		return nil, err
	}
	return &user, nil
}

func (a *AuthRepo) GetUserById(userId string) (*User, error) {
	var user User
	err := a.Db.Preload("Wallet").Where("id = ?", userId).First(&user).Error
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

func (a *AuthRepo) GetWalletWithLock(tx *gorm.DB, userId string) (*Wallet, error) {
	var wallet Wallet
	err := tx.Set("gorm:query_option", "FOR UPDATE").
		Where("user_id = ?", userId).
		First(&wallet).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (a *AuthRepo) UpdateUserRole(user *User) error {
	return a.Db.
		Model(&User{}).
		Where("email = ?", user.Email).
		Update("role", user.Role).Error
}
