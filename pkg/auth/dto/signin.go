package dto

type SignInResult struct {
	AccessToken string   `json:"access_token"`
	AuthData    AuthData `json:"client_info"`
}

type SignInInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Method   string
}

type CompleteRegistrationInput struct {
	FirstName       string `json:"first_name"`
	MiddleName      string `json:"middle_name"`
	LastName        string `json:"last_name"`
	Email           string
	Password        []byte `json:"password"`
	ConfirmPassword []byte `json:"confirm_password"`
}
