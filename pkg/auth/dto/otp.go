package dto

type OtpInput struct {
	Email string `json:"email"`
}

type OtpResult struct {
	Email string `json:"email"`
	OTP   []byte `json:"otp"`
}
