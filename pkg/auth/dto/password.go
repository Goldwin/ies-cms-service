package dto

type PasswordResetInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Code     string `json:"code"`
}

type PasswordResult struct {
	Email string `json:"email"`
}

type PasswordResetCodeResult struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}
