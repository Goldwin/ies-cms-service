package dto

type AuthData struct {
	ID         string   `json:"id"`
	FirstName  string   `json:"first_name"`
	MiddleName string   `json:"middle_name"`
	LastName   string   `json:"last_name"`
	Email      string   `json:"email"`
	Scopes     []string `json:"scopes"`
}

type AuthInput struct {
	Token     string
	SecretKey []byte
}
