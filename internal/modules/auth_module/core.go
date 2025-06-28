package authmodule

import (
	"errors"
	"log"
	"paytm-project/internal/models"
)

type IAuthCore interface {
	UserSignup(signupRequest *SignupRequestDto) (*SignupResponseDto, error)
	UserLogin(loginRequest *LoginRequestDto) (string, error)
	GetRepo() IAuthRepo
}

type AuthCore struct {
	AuthRepo IAuthRepo
}

func (ac *AuthCore) GetRepo() IAuthRepo {
	return ac.AuthRepo
}

func (a *AuthCore) UserSignup(signupRequest *SignupRequestDto) (*SignupResponseDto, error) {
	hashedPassword, perr := HashPassword(signupRequest.Password)
	if perr != nil {
		log.Println("error while hashing the password", perr.Error())
		return nil, perr
	}
	signupRequest.Password = hashedPassword
	user := SignupRequestDtoToUserMapper(signupRequest)
	err := a.AuthRepo.CreateUser(user)

	if err != nil {
		return nil, err
	}
	return &SignupResponseDto{BaseSuccessResponse: models.BaseSuccessResponse{
		Success: true,
	}}, nil
}

func (a *AuthCore) UserLogin(loginRequest *LoginRequestDto) (string, error) {
	user, err := a.AuthRepo.GetUserByEmail(loginRequest.Email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}
	if !CheckPassword(user.Password, loginRequest.Password) {
		return "", errors.New("invalid email or password")
	}

	claims := map[string]interface{}{
		"email":   user.Email,
		"wallet":  user.Wallet.Id,
		"balance": user.Wallet.Balance,
		"role":    user.Role,
	}

	token, terr := GenerateJwt(claims)
	if terr != nil {
		log.Println("error generating token")
		return "", terr
	}

	return token, nil
}
