package dto

type PasswordResetInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

type PasswordResult struct {
	Email string `json:"email"`
}

type PasswordResetTokenResult struct {
	Email string `json:"email"`
	Token string `json:"token"`
}
