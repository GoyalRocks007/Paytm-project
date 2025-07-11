package authmodule

import (
	"paytm-project/internal/models"

	"gorm.io/gorm"
)

type Resource string
type Role string

const (
	AdminRoleName Role = "ADMIN"
	UserRoleName  Role = "USER"
)

type User struct {
	models.BaseModel
	Name     string
	Contact  string
	Email    string `gorm:"unique"`
	Username string `gorm:"unique"`
	Wallet   Wallet `gorm:"foreignKey:UserId" json:"wallet"`
	Password string
	Role     Role
}

type Wallet struct {
	models.BaseModel
	UserId  string
	Balance float64
}

type SignupRequestDto struct {
	Name     string `json:"name" validate:"required"`
	Contact  string `json:"contact" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type SignupResponseDto struct {
	models.BaseSuccessResponse
}

type LoginRequestDto struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponseDto struct {
	models.BaseSuccessResponse
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if err := u.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}
	u.Role = UserRoleName
	u.Wallet = Wallet{
		Balance: 10000,
	}
	return
}

func (r *Role) ValidateRole() bool {
	if *r == UserRoleName || *r == AdminRoleName {
		return true
	}
	return false
}
