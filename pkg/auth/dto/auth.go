package dto

type AuthData struct {
	Email  string   `json:"email"`
	Scopes []string `json:"scopes"`
}

type AuthInput struct {
	Token     string
	SecretKey []byte
}
