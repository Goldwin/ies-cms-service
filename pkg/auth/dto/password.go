package dto

type PasswordInput struct {
	Email    string `json:"email"`
	Password []byte `json:"password"`
	Token    string `json:"token"`
}

type PasswordResult struct {
	Email string `json:"email"`
}

type PasswordResetTokenResult struct {
	Email string `json:"email"`
	Token string `json:"token"`
}
