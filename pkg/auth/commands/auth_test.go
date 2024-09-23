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

func (t *AuthCommandTest) TestExecuteInvalidTokenFailed() {
	secretKey := []byte("secret-key3")
	result := commands.AuthCommand{
		Token:     "token",
		SecretKey: secretKey,
	}.Execute(t.ctx)
	assert.Equal(t.T(), common.ExecutionStatusFailed, result.Status)
	assert.Equal(t.T(), commands.AuthErrorInvalidToken, result.Error.Code)
}

func (t *AuthCommandTest) TestExecuteMalformedTokenFailed() {
	secretKey := []byte("secret-key4")
	token := createInvalidToken(secretKey)
	result := commands.AuthCommand{
		Token:     token,
		SecretKey: secretKey,
	}.Execute(t.ctx)
	assert.Equal(t.T(), common.ExecutionStatusFailed, result.Status)
	assert.Equal(t.T(), commands.AuthErrorInvalidToken, result.Error.Code)
}

func (t *AuthCommandTest) TestExecuteInvalidEmailFailed() {
	secretKey := []byte("secret-key1")
	token := createValidToken("invalid@@email.com", secretKey)
	result := commands.AuthCommand{
		Token:     token,
		SecretKey: secretKey,
	}.Execute(t.ctx)
	assert.Equal(t.T(), common.ExecutionStatusFailed, result.Status)
	assert.Equal(t.T(), commands.AuthErrorInvalidToken, result.Error.Code)
}

func (t *AuthCommandTest) TestExecuteWrongKeyFailed() {
	secretKey := []byte("secret-key2")
	temperedKey := []byte("tempered-key")
	token := createValidToken("example1@email.com", temperedKey)
	result := commands.AuthCommand{
		Token:     token,
		SecretKey: secretKey,
	}.Execute(t.ctx)
	assert.Equal(t.T(), common.ExecutionStatusFailed, result.Status)
	assert.Equal(t.T(), commands.AuthErrorInvalidToken, result.Error.Code)
}

func (t *AuthCommandTest) TestExecuteFailedToRetrieveAccountInfoFailed() {
	secretKey := []byte("secret-key3")
	token := createValidToken("example2@email.com", secretKey)
	t.accountRepository.EXPECT().Get(mock.Anything).Return(nil, errors.New("Failed to Fetch Account"))
	result := commands.AuthCommand{
		Token:     token,
		SecretKey: secretKey,
	}.Execute(t.ctx)
	assert.Equal(t.T(), common.ExecutionStatusFailed, result.Status)
	assert.Equal(t.T(), commands.AuthErrorFailedToRetrieveAccount, result.Error.Code)
}

func (t *AuthCommandTest) TestExecuteAccountNotExistsFailed() {
	secretKey := []byte("secret-key")
	token := createValidToken("example5@email.com", secretKey)
	t.accountRepository.EXPECT().Get(mock.Anything).Return(nil, nil)
	result := commands.AuthCommand{
		Token:     token,
		SecretKey: secretKey,
	}.Execute(t.ctx)
	assert.Equal(t.T(), common.ExecutionStatusFailed, result.Status)
	assert.Equal(t.T(), commands.AuthErrorAccountDoesNotExist, result.Error.Code)
}

func (t *AuthCommandTest) TestExecuteAccountExistsSuccess() {
	secretKey := []byte("secret-key")
	token := createValidToken("example@email.com", secretKey)
	t.accountRepository.EXPECT().Get(mock.Anything).Return(
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
		}, nil)
	result := commands.AuthCommand{
		Token:     token,
		SecretKey: secretKey,
	}.Execute(t.ctx)
	assert.Equal(t.T(), common.ExecutionStatusSuccess, result.Status)
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
