package dto

type PasswordInput struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type PasswordResult struct {
	Email string `json:"email"`
}
