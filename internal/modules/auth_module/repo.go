package authmodule

import (
	"log"
	"paytm-project/internal/models"
)

type IAuthRepo interface {
	CreateUser(user *User) error
	GetUserByEmail(email string) (*User, error)
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
	err := a.Db.Where("email = ?", email).First(&user).Error
	if err != nil {
		log.Println("couldn't find user with given email:", err.Error())
		return nil, err
	}
	return &user, nil
}
