package commands_test

import (
	"fmt"
	"testing"

	"github.com/Goldwin/ies-pik-cms/pkg/auth/commands"
	"github.com/Goldwin/ies-pik-cms/pkg/auth/entities"
	"github.com/Goldwin/ies-pik-cms/pkg/auth/repositories/mocks"
	common "github.com/Goldwin/ies-pik-cms/pkg/common/commands"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type CompleteRegistrationCommandTest struct {
	ctx               *mocks.CommandContext
	accountRepository *mocks.AccountRepository
	suite.Suite
}

func (t *CompleteRegistrationCommandTest) SetupTest() {
	t.ctx = mocks.NewCommandContext(t.T())
	t.accountRepository = mocks.NewAccountRepository(t.T())
	t.ctx.EXPECT().AccountRepository().Maybe().Return(t.accountRepository)
}

func (t *CompleteRegistrationCommandTest) TestExecute_CompleteRegistrationCompletedAlready_Failed() {
	t.accountRepository.EXPECT().GetAccount(mock.AnythingOfType("entities.EmailAddress")).Return(&entities.Account{
		Email: entities.EmailAddress("p6bqK@example.com"),
		Person: entities.Person{
			FirstName:  "firstName",
			MiddleName: "middleName",
			LastName:   "lastName",
			ID:         "123",
		},
		Roles: []entities.Role{{ID: 1, Name: "Church Member", Scopes: []entities.Scope{
			entities.EventCheckIn,
			entities.EventView,
		}}},
	}, nil)
	result := commands.CompleteRegistrationCommand{
		FirstName:  "firstName",
		MiddleName: "middleName",
		LastName:   "lastName",
		Email:      "p6bqK@example.com",
	}.Execute(t.ctx)
	assert.Equal(t.T(), common.ExecutionStatusFailed, result.Status)
	assert.Equal(t.T(), commands.CompleteRegistrationErrorAlreadyCompleted, result.Error.Code)
}

func (t *CompleteRegistrationCommandTest) TestExecute_CompleteRegistrationAccountReadFailure_Failed() {
	t.accountRepository.EXPECT().GetAccount(mock.AnythingOfType("entities.EmailAddress")).Return(nil, fmt.Errorf("error"))
	result := commands.CompleteRegistrationCommand{
		FirstName:  "firstName",
		MiddleName: "middleName",
		LastName:   "lastName",
		Email:      "p6bqK@example.com",
	}.Execute(t.ctx)
	assert.Equal(t.T(), common.ExecutionStatusFailed, result.Status)
	assert.Equal(t.T(), commands.CompleteRegistrationErrorFailedToGetAccount, result.Error.Code)
}

func (t *CompleteRegistrationCommandTest) TestExecute_CompleteRegistrationAccountIsNotRegistered_Failed() {
	t.accountRepository.EXPECT().GetAccount(mock.AnythingOfType("entities.EmailAddress")).Return(nil, nil)
	result := commands.CompleteRegistrationCommand{
		FirstName:  "firstName",
		MiddleName: "middleName",
		LastName:   "lastName",
		Email:      "p6bqK@example.com",
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
		FirstName:  "firstName",
		MiddleName: "middleName",
		LastName:   "lastName",
		Email:      "p6bqK@example.com",
	}.Execute(t.ctx)
	assert.Equal(t.T(), common.ExecutionStatusFailed, result.Status)
	assert.Equal(t.T(), commands.CompleteRegistrationErrorUpdateFailure, result.Error.Code)
}

func (t *CompleteRegistrationCommandTest) TestExecute_CompleteRegistration_Success() {
	personId := "123"
	t.accountRepository.EXPECT().GetAccount(mock.AnythingOfType("entities.EmailAddress")).Return(&entities.Account{
		Email: entities.EmailAddress("p6bqK@example.com"),
	}, nil)

	t.accountRepository.EXPECT().UpdateAccount(mock.AnythingOfType("entities.Account")).RunAndReturn(
		func(a entities.Account) (*entities.Account, error) {
			a.Person.ID = personId
			return &a, nil
		},
	)
	result := commands.CompleteRegistrationCommand{
		FirstName:  "firstName",
		MiddleName: "middleName",
		LastName:   "lastName",
		Email:      "p6bqK@example.com",
	}.Execute(t.ctx)

	assert.Equal(t.T(), common.ExecutionStatusSuccess, result.Status)
	assert.Equal(t.T(), personId, result.Result.ID)
}

func TestCompleteRegistration(t *testing.T) {
	suite.Run(t, new(CompleteRegistrationCommandTest))
}
