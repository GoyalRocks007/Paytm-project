package otpmodule

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"paytm-project/internal/clients"
	"paytm-project/internal/models"
	authmodule "paytm-project/internal/modules/auth_module"
	"time"
)

var (
	NewOtpCore = func(clients *clients.ClientRegistry) IOtpCore {
		return &OtpCore{
			clients: clients,
		}
	}
)

type IOtpCore interface {
	GenerateOtp(ctx context.Context, generateOtpRequest *GenerateOtpRequest) (*GenerateOtpResponse, error)
	SendOtp(ctx context.Context, otp string, receiver string) error
	VerifyOtp(ctx context.Context, verifyOtpRequest *VerifyOtpRequest) (*VerifyOtpResponse, string, map[string]interface{}, error)
}

type OtpCore struct {
	clients *clients.ClientRegistry
}

func (oc *OtpCore) GenerateOtp(ctx context.Context, generateOtpRequest *GenerateOtpRequest) (*GenerateOtpResponse, error) {
	otp, err := OtpGenerator()
	if err != nil {
		return nil, errors.New("error generating otp")
	}
	oc.clients.RedisClient.Set(ctx,
		generateOtpRequest.Receiver,
		&Otp{
			Code:     otp,
			Attempts: 0,
			Claims:   generateOtpRequest.Claims,
		},
		time.Minute*5)
	err = oc.SendOtp(ctx, otp, generateOtpRequest.Receiver)
	if err != nil {
		log.Println("error in sending otp email for", generateOtpRequest.Receiver, err.Error())
	}
	return &GenerateOtpResponse{
		BaseSuccessResponse: models.BaseSuccessResponse{Success: true},
	}, nil

}

func (oc *OtpCore) SendOtp(ctx context.Context, otp string, receiver string) error {
	sender := os.Getenv("GOOGLE_EMAIL_ADDRESS")
	subject := os.Getenv("OTP_EMAIL_SUBJECT")
	body := fmt.Sprintf(os.Getenv("OTP_EMAIL_BODY"), otp)
	return oc.clients.EmailClient.SendEmail(
		ctx,
		sender,
		receiver,
		subject,
		body,
	)
}

func (oc *OtpCore) VerifyOtp(ctx context.Context, verifyOtpRequest *VerifyOtpRequest) (*VerifyOtpResponse, string, map[string]interface{}, error) {
	receiver := verifyOtpRequest.Receiver
	var otp Otp
	data, exist, err := oc.clients.RedisClient.Get(ctx, receiver)
	if !exist || err != nil {
		return nil, "", nil, errors.New("otp expired please try again")
	}
	jsonData, err1 := json.Marshal(data)
	if err1 != nil {
		return nil, "", nil, err
	}
	if err2 := json.Unmarshal(jsonData, &otp); err != nil {
		return nil, "", nil, err2
	}

	if otp.Code == verifyOtpRequest.Otp {
		err := oc.clients.RedisClient.Del(ctx, receiver)
		if err != nil {
			log.Println("error in deleting key", receiver, err.Error())
			return nil, "", nil, err
		}
		token, err := authmodule.GenerateJwt(map[string]interface{}{
			"receiver": verifyOtpRequest.Receiver,
		})
		if err != nil {
			return nil, "", nil, errors.New("error generating token")
		}
		return &VerifyOtpResponse{
			BaseSuccessResponse: models.BaseSuccessResponse{Success: true},
		}, token, otp.Claims, nil
	} else {
		return nil, "", nil, errors.New("invalid otp")
	}

}
