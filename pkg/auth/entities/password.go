package entities

type PasswordDetail struct {
	EmailAddress EmailAddress
	Salt         []byte
	PasswordHash []byte
}
