package commands_test

import (
	"crypto/sha256"
	"errors"
	"testing"
	"time"

	"github.com/Goldwin/ies-pik-cms/pkg/auth/commands"
	"github.com/Goldwin/ies-pik-cms/pkg/auth/entities"
	"github.com/Goldwin/ies-pik-cms/pkg/auth/repositories/mocks"
	common "github.com/Goldwin/ies-pik-cms/pkg/common/commands"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type SignInCommandTest struct {
	ctx                *mocks.CommandContext
	otpRepository      *mocks.OtpRepository
	accountRepository  *mocks.AccountRepository
	passwordRepository *mocks.PasswordRepository
	suite.Suite
}

func (t *SignInCommandTest) SetupTest() {
	t.otpRepository = mocks.NewOtpRepository(t.T())
	t.accountRepository = mocks.NewAccountRepository(t.T())
	t.passwordRepository = mocks.NewPasswordRepository(t.T())
	t.ctx = mocks.NewCommandContext(t.T())

	t.ctx.EXPECT().OtpRepository().Maybe().Return(t.otpRepository)
	t.ctx.EXPECT().AccountRepository().Maybe().Return(t.accountRepository)
	t.ctx.EXPECT().PasswordRepository().Maybe().Return(t.passwordRepository)
}

func (t *SignInCommandTest) TestExecute_OtpSignIn_Success() {
	email := "p6bqK@example.com"
	password := []byte("password")
	salt := []byte("salt")
	passwordHash := sha256.Sum256(append(password, salt...))

	t.otpRepository.EXPECT().GetOtp(mock.AnythingOfType("entities.EmailAddress")).Return(&entities.Otp{
		EmailAddress: entities.EmailAddress(email),
		PasswordHash: passwordHash[:],
		Salt:         salt,
		ExpiredTime:  time.Now().Add(time.Minute),
	}, nil)

	t.otpRepository.EXPECT().RemoveOtp(mock.AnythingOfType("entities.Otp")).Maybe().Return(nil)
	t.accountRepository.EXPECT().GetAccount(mock.AnythingOfType("entities.EmailAddress")).Return(&entities.Account{
		Email: entities.EmailAddress(email),
	}, nil)

	result := commands.SigninCommand{
		Email:     email,
		Password:  password,
		Method:    commands.SigninMethodOTP,
		SecretKey: []byte("secret-key"),
	}.Execute(t.ctx)

	assert.Equal(t.T(), common.ExecutionStatusSuccess, result.Status)
	assert.NotEmpty(t.T(), result.Result)
}

func (t *SignInCommandTest) TestExecute_OtpSignInAccountNotExists_Success() {
	personId := "person-id-1"
	email := "p6bqK@example.com"
	password := []byte("password")
	salt := []byte("salt")
	passwordHash := sha256.Sum256(append(password, salt...))

	t.otpRepository.EXPECT().GetOtp(mock.AnythingOfType("entities.EmailAddress")).Return(&entities.Otp{
		EmailAddress: entities.EmailAddress(email),
		PasswordHash: passwordHash[:],
		Salt:         salt,
		ExpiredTime:  time.Now().Add(time.Minute),
	}, nil)

	t.otpRepository.EXPECT().RemoveOtp(mock.AnythingOfType("entities.Otp")).Maybe().Return(nil)
	t.accountRepository.EXPECT().GetAccount(mock.AnythingOfType("entities.EmailAddress")).Return(nil, nil)
	t.accountRepository.EXPECT().AddAccount(mock.AnythingOfType("entities.Account")).RunAndReturn(
		func(a entities.Account) (*entities.Account, error) {
			a.Person.ID = personId
			return &a, nil
		},
	)

	result := commands.SigninCommand{
		Email:     email,
		Password:  password,
		Method:    commands.SigninMethodOTP,
		SecretKey: []byte("secret-key"),
	}.Execute(t.ctx)

	assert.Equal(t.T(), common.ExecutionStatusSuccess, result.Status)
	assert.Equal(t.T(), personId, result.Result.AuthData.ID)
}

func (t *SignInCommandTest) TestExecute_OtpSignInAccountNotExistsFailedToCreate_Failed() {
	email := "p6bqK@example.com"
	password := []byte("password")
	salt := []byte("salt")
	passwordHash := sha256.Sum256(append(password, salt...))

	t.otpRepository.EXPECT().GetOtp(mock.AnythingOfType("entities.EmailAddress")).Return(&entities.Otp{
		EmailAddress: entities.EmailAddress(email),
		PasswordHash: passwordHash[:],
		Salt:         salt,
		ExpiredTime:  time.Now().Add(time.Minute),
	}, nil)

	t.otpRepository.EXPECT().RemoveOtp(mock.AnythingOfType("entities.Otp")).Maybe().Return(nil)
	t.accountRepository.EXPECT().GetAccount(mock.AnythingOfType("entities.EmailAddress")).Return(nil, nil)
	t.accountRepository.EXPECT().AddAccount(mock.AnythingOfType("entities.Account")).Return(nil, errors.New("failed to create account"))

	result := commands.SigninCommand{
		Email:     email,
		Password:  password,
		Method:    commands.SigninMethodOTP,
		SecretKey: []byte("secret-key"),
	}.Execute(t.ctx)

	assert.Equal(t.T(), common.ExecutionStatusFailed, result.Status)
	assert.Equal(t.T(), commands.SignInErrorFailedToCreateNewAccount, result.Error.Code)
}

func (t *SignInCommandTest) TestExecute_OtpSignInFailedToRetrieveAccount_Failed() {
	email := "p6bqK@example.com"
	password := []byte("password")
	salt := []byte("salt")
	passwordHash := sha256.Sum256(append(password, salt...))

	t.otpRepository.EXPECT().GetOtp(mock.AnythingOfType("entities.EmailAddress")).Return(&entities.Otp{
		EmailAddress: entities.EmailAddress(email),
		PasswordHash: passwordHash[:],
		Salt:         salt,
		ExpiredTime:  time.Now().Add(time.Minute),
	}, nil)

	t.otpRepository.EXPECT().RemoveOtp(mock.AnythingOfType("entities.Otp")).Maybe().Return(nil)
	t.accountRepository.EXPECT().GetAccount(mock.AnythingOfType("entities.EmailAddress")).Return(nil, errors.New("Failed to get account"))

	result := commands.SigninCommand{
		Email:     email,
		Password:  password,
		Method:    commands.SigninMethodOTP,
		SecretKey: []byte("secret-key"),
	}.Execute(t.ctx)

	assert.Equal(t.T(), common.ExecutionStatusFailed, result.Status)
}

func (t *SignInCommandTest) TestExecute_OtpSignInFailedToInvalidate_Failed() {
	email := "p6bqK@example.com"
	password := []byte("password")
	salt := []byte("salt")
	passwordHash := sha256.Sum256(append(password, salt...))

	t.otpRepository.EXPECT().GetOtp(mock.AnythingOfType("entities.EmailAddress")).Return(&entities.Otp{
		EmailAddress: entities.EmailAddress(email),
		PasswordHash: passwordHash[:],
		Salt:         salt,
		ExpiredTime:  time.Now().Add(time.Minute),
	}, nil)

	t.otpRepository.EXPECT().RemoveOtp(mock.AnythingOfType("entities.Otp")).Maybe().Return(errors.New("Failed to invalidate otp"))

	result := commands.SigninCommand{
		Email:     email,
		Password:  password,
		Method:    commands.SigninMethodOTP,
		SecretKey: []byte("secret-key"),
	}.Execute(t.ctx)

	assert.Equal(t.T(), common.ExecutionStatusFailed, result.Status)
	assert.Equal(t.T(), commands.SigninErrorOTPFailedToInvalidateOTP, result.Error.Code)
}

func (t *SignInCommandTest) TestExecute_OtpSignInWrongOtp_Failed() {
	email := "p6bqK@example.com"
	password := []byte("password")
	salt := []byte("salt")
	passwordHash := sha256.Sum256(append(password, salt...))

	t.otpRepository.EXPECT().RemoveOtp(mock.AnythingOfType("entities.Otp")).Maybe().Return(nil)

	t.otpRepository.EXPECT().GetOtp(mock.AnythingOfType("entities.EmailAddress")).Return(&entities.Otp{
		EmailAddress: entities.EmailAddress(email),
		PasswordHash: passwordHash[:],
		Salt:         salt,
		ExpiredTime:  time.Now().Add(time.Minute),
	}, nil)

	result := commands.SigninCommand{
		Email:     email,
		Password:  []byte("wrong-password"),
		Method:    commands.SigninMethodOTP,
		SecretKey: []byte("secret-key"),
	}.Execute(t.ctx)

	assert.Equal(t.T(), common.ExecutionStatusFailed, result.Status)
	assert.Equal(t.T(), commands.SigninErrorWrongOtp, result.Error.Code)
}

func (t *SignInCommandTest) TestExecute_OtpSignInExpired_Failed() {
	email := "p6bqK@example.com"
	password := []byte("password")
	salt := []byte("salt")
	passwordHash := sha256.Sum256(append(password, salt...))

	t.otpRepository.EXPECT().RemoveOtp(mock.AnythingOfType("entities.Otp")).Maybe().Return(nil)

	t.otpRepository.EXPECT().GetOtp(mock.AnythingOfType("entities.EmailAddress")).Return(&entities.Otp{
		EmailAddress: entities.EmailAddress(email),
		PasswordHash: passwordHash[:],
		Salt:         salt,
		ExpiredTime:  time.Now().Add(-time.Minute),
	}, nil)

	result := commands.SigninCommand{
		Email:     email,
		Password:  password,
		Method:    commands.SigninMethodOTP,
		SecretKey: []byte("secret-key"),
	}.Execute(t.ctx)

	assert.Equal(t.T(), common.ExecutionStatusFailed, result.Status)
	assert.Equal(t.T(), commands.SigninErrorOtpExired, result.Error.Code)
}

func (t *SignInCommandTest) TestExecute_OtpNotFound_Failed() {
	email := "p6bqK@example.com"
	password := []byte("password")

	t.otpRepository.EXPECT().RemoveOtp(mock.AnythingOfType("entities.Otp")).Maybe().Return(nil)
	t.otpRepository.EXPECT().GetOtp(mock.AnythingOfType("entities.EmailAddress")).Return(nil, nil)

	result := commands.SigninCommand{
		Email:     email,
		Password:  password,
		Method:    commands.SigninMethodOTP,
		SecretKey: []byte("secret-key"),
	}.Execute(t.ctx)

	assert.Equal(t.T(), common.ExecutionStatusFailed, result.Status)
}

func (t *SignInCommandTest) TestExecute_OtpFailedToRetrieve_Failed() {
	email := "p6bqK@example.com"
	password := []byte("password")

	t.otpRepository.EXPECT().RemoveOtp(mock.AnythingOfType("entities.Otp")).Maybe().Return(nil)
	t.otpRepository.EXPECT().GetOtp(mock.AnythingOfType("entities.EmailAddress")).Return(nil, errors.New("Failed to Retrieve stored OTP"))

	result := commands.SigninCommand{
		Email:     email,
		Password:  password,
		Method:    commands.SigninMethodOTP,
		SecretKey: []byte("secret-key"),
	}.Execute(t.ctx)

	assert.Equal(t.T(), common.ExecutionStatusFailed, result.Status)
	assert.Equal(t.T(), commands.SigninErrorFailedToGetOtp, result.Error.Code)
}

func (t *SignInCommandTest) TestExecute_PasswordLogin_Failed() {
	email := "p6bqK@example.com"
	password := []byte("password")

	t.passwordRepository.EXPECT().Get(mock.AnythingOfType("entities.EmailAddress")).Return(
		nil,
		errors.New("Failed to Retrieve stored password"),
	)

	result := commands.SigninCommand{
		Email:     email,
		Password:  password,
		Method:    commands.SigninMethodPassword,
		SecretKey: []byte("secret-key"),
	}.Execute(t.ctx)

	assert.Equal(t.T(), common.ExecutionStatusFailed, result.Status)
	assert.Equal(t.T(), commands.SignInErrorPasswordFailedToGetAccountDetail, result.Error.Code)
}

func TestSignIn(t *testing.T) {
	suite.Run(t, new(SignInCommandTest))
}
