package otpmodule

import "paytm-project/internal/models"

type Otp struct {
	Code     string
	Exp      int64
	Attempts int16
}

type GenerateOtpRequest struct {
	Receiver string `json:"receiver"`
}

type GenerateOtpResponse struct {
	models.BaseSuccessResponse
}

type VerifyOtpRequest struct {
	Otp      string `json:"otp"`
	Receiver string `json:"receiver"`
}

type VerifyOtpResponse struct {
	models.BaseSuccessResponse
}
