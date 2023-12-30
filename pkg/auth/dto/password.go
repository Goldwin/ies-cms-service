package dto

type PasswordInput struct {
	Email           string `json:"email"`
	Password        []byte `json:"password"`
	ConfirmPassword []byte `json:"confirm_password"`
}

type PasswordResult struct {
	Email string `json:"email"`
}
