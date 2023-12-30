package commands_test

import (
	"errors"
	"testing"
	"time"

	"github.com/Goldwin/ies-pik-cms/pkg/auth/commands"
	"github.com/Goldwin/ies-pik-cms/pkg/auth/entities"
	"github.com/Goldwin/ies-pik-cms/pkg/auth/repositories/mocks"
	common "github.com/Goldwin/ies-pik-cms/pkg/common/commands"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type AuthCommandTest struct {
	ctx               *mocks.CommandContext
	accountRepository *mocks.AccountRepository
	suite.Suite
}

func (t *AuthCommandTest) SetupTest() {
	t.accountRepository = mocks.NewAccountRepository(t.T())
	t.ctx = mocks.NewCommandContext(t.T())
	t.ctx.EXPECT().AccountRepository().Maybe().Return(t.accountRepository)
}

func (t *AuthCommandTest) TestExecute_InvalidToken_Failed() {
	secretKey := []byte("secret-key")
	result := commands.AuthCommand{
		Token:     "token",
		SecretKey: secretKey,
	}.Execute(t.ctx)
	assert.Equal(t.T(), common.ExecutionStatusFailed, result.Status)
	assert.Equal(t.T(), commands.AuthErrorInvalidToken, result.Error.Code)
}

func (t *AuthCommandTest) TestExecute_MalformedToken_Failed() {
	secretKey := []byte("secret-key")
	token := createInvalidToken(secretKey)
	result := commands.AuthCommand{
		Token:     token,
		SecretKey: secretKey,
	}.Execute(t.ctx)
	assert.Equal(t.T(), common.ExecutionStatusFailed, result.Status)
	assert.Equal(t.T(), commands.AuthErrorInvalidToken, result.Error.Code)
}

func (t *AuthCommandTest) TestExecute_InvalidEmail_Failed() {
	secretKey := []byte("secret-key")
	token := createValidToken("invalid@@email.com", secretKey)
	result := commands.AuthCommand{
		Token:     token,
		SecretKey: secretKey,
	}.Execute(t.ctx)
	assert.Equal(t.T(), common.ExecutionStatusFailed, result.Status)
	assert.Equal(t.T(), commands.AuthErrorInvalidToken, result.Error.Code)
}

func (t *AuthCommandTest) TestExecute_WrongKey_Failed() {
	secretKey := []byte("secret-key")
	temperedKey := []byte("tempered-key")
	token := createValidToken("example@email.com", temperedKey)
	result := commands.AuthCommand{
		Token:     token,
		SecretKey: secretKey,
	}.Execute(t.ctx)
	assert.Equal(t.T(), common.ExecutionStatusFailed, result.Status)
	assert.Equal(t.T(), commands.AuthErrorInvalidToken, result.Error.Code)
}

func (t *AuthCommandTest) TestExecute_FailedToRetrieveAccountInfo_Failed() {
	secretKey := []byte("secret-key")
	token := createValidToken("example@email.com", secretKey)
	t.accountRepository.EXPECT().GetAccount(mock.AnythingOfType("entities.EmailAddress")).Return(nil, errors.New("Failed to Fetch Account"))
	result := commands.AuthCommand{
		Token:     token,
		SecretKey: secretKey,
	}.Execute(t.ctx)
	assert.Equal(t.T(), common.ExecutionStatusFailed, result.Status)
	assert.Equal(t.T(), commands.AuthErrorFailedToRetrieveAccount, result.Error.Code)
}

func (t *AuthCommandTest) TestExecute_AccountNotExists_Failed() {
	secretKey := []byte("secret-key")
	token := createValidToken("example@email.com", secretKey)
	t.accountRepository.EXPECT().GetAccount(mock.AnythingOfType("entities.EmailAddress")).Return(nil, nil)
	result := commands.AuthCommand{
		Token:     token,
		SecretKey: secretKey,
	}.Execute(t.ctx)
	assert.Equal(t.T(), common.ExecutionStatusFailed, result.Status)
	assert.Equal(t.T(), commands.AuthErrorAccountDoesNotExist, result.Error.Code)
}

func (t *AuthCommandTest) TestExecute_AccountExists_Success() {
	secretKey := []byte("secret-key")
	personId := "account-id"
	token := createValidToken("example@email.com", secretKey)
	t.accountRepository.EXPECT().GetAccount(mock.AnythingOfType("entities.EmailAddress")).Return(
		&entities.Account{
			Email: "example@email.com",
			Roles: []entities.Role{
				{
					Name: "Member",
					Scopes: []entities.Scope{
						entities.EventCheckIn,
					},
				},
			},
			Person: entities.Person{ID: personId},
		}, nil)
	result := commands.AuthCommand{
		Token:     token,
		SecretKey: secretKey,
	}.Execute(t.ctx)
	assert.Equal(t.T(), common.ExecutionStatusSuccess, result.Status)
	assert.Equal(t.T(), personId, result.Result.ID)
}

func TestAuth(t *testing.T) {
	suite.Run(t, new(AuthCommandTest))
}

func createValidToken(email string, secretKey []byte) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"email": email,
			"exp":   time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		panic(err)
	}

	return tokenString
}

func createInvalidToken(secretKey []byte) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"exp": time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		panic(err)
	}

	return tokenString
}
