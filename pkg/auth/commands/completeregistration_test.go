package commands_test

import (
	"fmt"
	"testing"

	"github.com/Goldwin/ies-pik-cms/pkg/auth/commands"
	"github.com/Goldwin/ies-pik-cms/pkg/auth/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/auth/entities"
	"github.com/Goldwin/ies-pik-cms/pkg/auth/repositories/mocks"
	common "github.com/Goldwin/ies-pik-cms/pkg/common/commands"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type CompleteRegistrationCommandTest struct {
	ctx                *mocks.CommandContext
	accountRepository  *mocks.AccountRepository
	passwordRepository *mocks.PasswordRepository
	suite.Suite
}

const (
	emailAddressType = "entities.EmailAddress"
	emailAddress     = "p6bqK@example.com"
	password         = "p@ssw0rd"
)

func (t *CompleteRegistrationCommandTest) SetupTest() {
	t.ctx = mocks.NewCommandContext(t.T())
	t.accountRepository = mocks.NewAccountRepository(t.T())
	t.passwordRepository = mocks.NewPasswordRepository(t.T())
	t.ctx.EXPECT().AccountRepository().Maybe().Return(t.accountRepository)
	t.ctx.EXPECT().PasswordRepository().Maybe().Return(t.passwordRepository)
}

func (t *CompleteRegistrationCommandTest) TestExecuteCompleteRegistrationCompletedAlreadyFailed() {
	t.accountRepository.EXPECT().GetAccount(mock.AnythingOfType(emailAddressType)).Return(&entities.Account{
		Email: entities.EmailAddress(emailAddress),
		Roles: []entities.Role{{ID: 1, Name: "Church Member", Scopes: []entities.Scope{
			entities.EventCheckIn,
			entities.EventView,
		}}},
	}, nil)
	result := commands.CompleteRegistrationCommand{
		Input: dto.CompleteRegistrationInput{
			FirstName:       "firstName",
			MiddleName:      "middleName",
			LastName:        "lastName",
			Email:           emailAddress,
			Password:        []byte(password),
			ConfirmPassword: []byte(password),
		},
	}.Execute(t.ctx)
	assert.Equal(t.T(), common.ExecutionStatusFailed, result.Status)
	assert.Equal(t.T(), commands.CompleteRegistrationErrorAlreadyCompleted, result.Error.Code)
}

func (t *CompleteRegistrationCommandTest) TestExecuteCompleteRegistrationAccountReadFailureFailed() {
	t.accountRepository.EXPECT().GetAccount(mock.AnythingOfType(emailAddressType)).Return(nil, fmt.Errorf("error"))
	result := commands.CompleteRegistrationCommand{
		Input: dto.CompleteRegistrationInput{
			FirstName:       "firstName",
			MiddleName:      "middleName",
			LastName:        "lastName",
			Email:           emailAddress,
			Password:        []byte(password),
			ConfirmPassword: []byte(password),
		},
	}.Execute(t.ctx)
	assert.Equal(t.T(), common.ExecutionStatusFailed, result.Status)
	assert.Equal(t.T(), commands.CompleteRegistrationErrorFailedToGetAccount, result.Error.Code)
}

func (t *CompleteRegistrationCommandTest) TestExecuteCompleteRegistrationAccountIsNotRegisteredFailed() {
	t.accountRepository.EXPECT().GetAccount(mock.AnythingOfType(emailAddressType)).Return(nil, nil)
	result := commands.CompleteRegistrationCommand{
		Input: dto.CompleteRegistrationInput{
			FirstName:       "firstName",
			MiddleName:      "middleName",
			LastName:        "lastName",
			Email:           emailAddress,
			Password:        []byte(password),
			ConfirmPassword: []byte(password),
		},
	}.Execute(t.ctx)
	assert.Equal(t.T(), common.ExecutionStatusFailed, result.Status)
	assert.Equal(t.T(), commands.CompleteRegistrationErrorAccountIsNotRegistered, result.Error.Code)
}

func (t *CompleteRegistrationCommandTest) TestExecuteCompleteRegistrationFailedToUpdateFailed() {
	t.accountRepository.EXPECT().GetAccount(mock.AnythingOfType(emailAddressType)).Return(&entities.Account{
		Email: entities.EmailAddress(emailAddress),
	}, nil)

	t.accountRepository.EXPECT().UpdateAccount(mock.AnythingOfType("entities.Account")).Return(nil, fmt.Errorf("failed to update"))
	result := commands.CompleteRegistrationCommand{
		Input: dto.CompleteRegistrationInput{
			FirstName:       "firstName",
			MiddleName:      "middleName",
			LastName:        "lastName",
			Email:           emailAddress,
			Password:        []byte(password),
			ConfirmPassword: []byte(password),
		},
	}.Execute(t.ctx)
	assert.Equal(t.T(), common.ExecutionStatusFailed, result.Status)
	assert.Equal(t.T(), commands.CompleteRegistrationErrorUpdateFailure, result.Error.Code)
}

func (t *CompleteRegistrationCommandTest) TestExecuteCompleteRegistrationSuccess() {
	t.accountRepository.EXPECT().GetAccount(mock.AnythingOfType(emailAddressType)).Return(&entities.Account{
		Email: entities.EmailAddress(emailAddress),
	}, nil)

	t.accountRepository.EXPECT().UpdateAccount(mock.AnythingOfType("entities.Account")).RunAndReturn(
		func(a entities.Account) (*entities.Account, error) {
			return &a, nil
		},
	)
	result := commands.CompleteRegistrationCommand{
		Input: dto.CompleteRegistrationInput{
			FirstName:       "firstName",
			MiddleName:      "middleName",
			LastName:        "lastName",
			Email:           emailAddress,
			Password:        []byte(password),
			ConfirmPassword: []byte(password),
		},
	}.Execute(t.ctx)

	assert.Equal(t.T(), common.ExecutionStatusSuccess, result.Status)
}

func TestCompleteRegistration(t *testing.T) {
	suite.Run(t, new(CompleteRegistrationCommandTest))
}
