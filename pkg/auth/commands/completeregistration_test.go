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

func (t *CompleteRegistrationCommandTest) SetupTest() {
	t.ctx = mocks.NewCommandContext(t.T())
	t.accountRepository = mocks.NewAccountRepository(t.T())
	t.passwordRepository = mocks.NewPasswordRepository(t.T())
	t.ctx.EXPECT().AccountRepository().Maybe().Return(t.accountRepository)
	t.ctx.EXPECT().PasswordRepository().Maybe().Return(t.passwordRepository)
}

func (t *CompleteRegistrationCommandTest) TestExecute_CompleteRegistrationCompletedAlready_Failed() {
	t.accountRepository.EXPECT().GetAccount(mock.AnythingOfType("entities.EmailAddress")).Return(&entities.Account{
		Email: entities.EmailAddress("p6bqK@example.com"),
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
			Email:           "p6bqK@example.com",
			Password:        []byte("p@ssw0rd"),
			ConfirmPassword: []byte("p@ssw0rd"),
		},
	}.Execute(t.ctx)
	assert.Equal(t.T(), common.ExecutionStatusFailed, result.Status)
	assert.Equal(t.T(), commands.CompleteRegistrationErrorAlreadyCompleted, result.Error.Code)
}

func (t *CompleteRegistrationCommandTest) TestExecute_CompleteRegistrationAccountReadFailure_Failed() {
	t.accountRepository.EXPECT().GetAccount(mock.AnythingOfType("entities.EmailAddress")).Return(nil, fmt.Errorf("error"))
	result := commands.CompleteRegistrationCommand{
		Input: dto.CompleteRegistrationInput{
			FirstName:       "firstName",
			MiddleName:      "middleName",
			LastName:        "lastName",
			Email:           "p6bqK@example.com",
			Password:        []byte("p@ssw0rd"),
			ConfirmPassword: []byte("p@ssw0rd"),
		},
	}.Execute(t.ctx)
	assert.Equal(t.T(), common.ExecutionStatusFailed, result.Status)
	assert.Equal(t.T(), commands.CompleteRegistrationErrorFailedToGetAccount, result.Error.Code)
}

func (t *CompleteRegistrationCommandTest) TestExecute_CompleteRegistrationAccountIsNotRegistered_Failed() {
	t.accountRepository.EXPECT().GetAccount(mock.AnythingOfType("entities.EmailAddress")).Return(nil, nil)
	result := commands.CompleteRegistrationCommand{
		Input: dto.CompleteRegistrationInput{
			FirstName:       "firstName",
			MiddleName:      "middleName",
			LastName:        "lastName",
			Email:           "p6bqK@example.com",
			Password:        []byte("p@ssw0rd"),
			ConfirmPassword: []byte("p@ssw0rd"),
		},
	}.Execute(t.ctx)
	assert.Equal(t.T(), common.ExecutionStatusFailed, result.Status)
	assert.Equal(t.T(), commands.CompleteRegistrationErrorAccountIsNotRegistered, result.Error.Code)
}

func (t *CompleteRegistrationCommandTest) TestExecute_CompleteRegistrationFailedToUpdate_Failed() {
	t.accountRepository.EXPECT().GetAccount(mock.AnythingOfType("entities.EmailAddress")).Return(&entities.Account{
		Email: entities.EmailAddress("p6bqK@example.com"),
	}, nil)

	t.accountRepository.EXPECT().UpdateAccount(mock.AnythingOfType("entities.Account")).Return(nil, fmt.Errorf("failed to update"))
	result := commands.CompleteRegistrationCommand{
		Input: dto.CompleteRegistrationInput{
			FirstName:       "firstName",
			MiddleName:      "middleName",
			LastName:        "lastName",
			Email:           "p6bqK@example.com",
			Password:        []byte("p@ssw0rd"),
			ConfirmPassword: []byte("p@ssw0rd"),
		},
	}.Execute(t.ctx)
	assert.Equal(t.T(), common.ExecutionStatusFailed, result.Status)
	assert.Equal(t.T(), commands.CompleteRegistrationErrorUpdateFailure, result.Error.Code)
}

func (t *CompleteRegistrationCommandTest) TestExecute_CompleteRegistration_Success() {
	t.accountRepository.EXPECT().GetAccount(mock.AnythingOfType("entities.EmailAddress")).Return(&entities.Account{
		Email: entities.EmailAddress("p6bqK@example.com"),
	}, nil)

	t.accountRepository.EXPECT().UpdateAccount(mock.AnythingOfType("entities.Account")).RunAndReturn(
		func(a entities.Account) (*entities.Account, error) {
			return &a, nil
		},
	)
	t.passwordRepository.EXPECT().Save(mock.AnythingOfType("entities.PasswordDetail")).Return(nil)
	result := commands.CompleteRegistrationCommand{
		Input: dto.CompleteRegistrationInput{
			FirstName:       "firstName",
			MiddleName:      "middleName",
			LastName:        "lastName",
			Email:           "p6bqK@example.com",
			Password:        []byte("p@ssw0rd"),
			ConfirmPassword: []byte("p@ssw0rd"),
		},
	}.Execute(t.ctx)

	assert.Equal(t.T(), common.ExecutionStatusSuccess, result.Status)
}

func TestCompleteRegistration(t *testing.T) {
	suite.Run(t, new(CompleteRegistrationCommandTest))
}
