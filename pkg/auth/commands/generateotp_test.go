package commands_test

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/Goldwin/ies-pik-cms/pkg/auth/commands"
	commandMock "github.com/Goldwin/ies-pik-cms/pkg/auth/commands/mocks"
	"github.com/Goldwin/ies-pik-cms/pkg/auth/entities"
	"github.com/Goldwin/ies-pik-cms/pkg/auth/repositories/mocks"
	common "github.com/Goldwin/ies-pik-cms/pkg/common/commands"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ToggleRandReader struct {
	willFail bool
	reader   io.Reader
}

func (f *ToggleRandReader) Read(p []byte) (n int, err error) {
	if f.willFail {
		return 0, fmt.Errorf("failed to read")
	}
	f.willFail = true
	return f.reader.Read(p)
}

type GenerateOtpCommandTest struct {
	ctx           *commandMock.CommandContext
	otpRepository *mocks.OtpRepository
	suite.Suite
}

func (t *GenerateOtpCommandTest) SetupTest() {
	t.otpRepository = mocks.NewOtpRepository(t.T())
	t.ctx = commandMock.NewCommandContext(t.T())

	t.ctx.EXPECT().OtpRepository().Maybe().Return(t.otpRepository)
}

func (t *GenerateOtpCommandTest) TestExecuteGenerateOtpSuccess() {
	var otp *entities.Otp
	t.otpRepository.EXPECT().Get(mock.Anything).Return(nil, nil)
	t.otpRepository.EXPECT().Save(mock.Anything).RunAndReturn(func(o *entities.Otp) (*entities.Otp, error) {
		otp = o
		return o, nil
	})
	result := commands.GenerateOtpCommand{
		Email:     "p6bq1@example.com",
		TTLMillis: 0,
	}.Execute(t.ctx)

	assert.Equal(t.T(), common.ExecutionStatusSuccess, result.Status)
	resultHash := sha256.Sum256(append(result.Result.OTP, otp.Salt[:]...))
	assert.Equal(t.T(), otp.PasswordHash[:], resultHash[:])
}
func (t *GenerateOtpCommandTest) TestExecuteGenerateOtpFailedToStoreFailed() {
	t.otpRepository.EXPECT().Get(mock.Anything).Return(nil, nil)
	t.otpRepository.EXPECT().Save(mock.Anything).Return(nil, fmt.Errorf("failed to add otp"))
	result := commands.GenerateOtpCommand{
		Email:     "p6bq2@example.com",
		TTLMillis: 0,
	}.Execute(t.ctx)

	assert.Equal(t.T(), common.ExecutionStatusFailed, result.Status)
	assert.Equal(t.T(), commands.GenerateOtpErrorFailedToStoreOtp, result.Error.Code)
}

func (t *GenerateOtpCommandTest) TestExecuteGenerateOtpInvalidEmailFailed() {
	t.otpRepository.EXPECT().Get(mock.Anything).Return(nil, nil)
	result := commands.GenerateOtpCommand{
		Email:     "p6bq3@@example.com",
		TTLMillis: 0,
	}.Execute(t.ctx)

	assert.Equal(t.T(), common.ExecutionStatusFailed, result.Status)
	assert.Equal(t.T(), commands.GenerateOtpErrorInvalidEmail, result.Error.Code)
}

func (t *GenerateOtpCommandTest) TestExecuteGenerateOtpFailedToGenOtpFailed() {
	t.otpRepository.EXPECT().Get(mock.Anything).Return(nil, nil)
	oldReader := rand.Reader
	rand.Reader = &ToggleRandReader{willFail: true, reader: rand.Reader}
	result := commands.GenerateOtpCommand{
		Email:     "p6bq4@example.com",
		TTLMillis: -1,
	}.Execute(t.ctx)

	rand.Reader = oldReader
	assert.Equal(t.T(), common.ExecutionStatusFailed, result.Status)
	assert.Equal(t.T(), commands.GenerateOtpErrorFailedToGenOtp, result.Error.Code)
}

func (t *GenerateOtpCommandTest) TestExecuteGenerateOtpFailedToGenSaltFailed() {
	t.otpRepository.EXPECT().Get(mock.Anything).Return(nil, nil)
	oldReader := rand.Reader
	newReader := &ToggleRandReader{willFail: false, reader: rand.Reader}
	rand.Reader = newReader
	result := commands.GenerateOtpCommand{
		Email:     "p6bq5@example.com",
		TTLMillis: -1,
	}.Execute(t.ctx)

	rand.Reader = oldReader
	assert.Equal(t.T(), common.ExecutionStatusFailed, result.Status)
	assert.Equal(t.T(), commands.GenerateOtpErrorFailedToGenOtp, result.Error.Code)
}

func (t *GenerateOtpCommandTest) TestExecuteGenerateOtpFailedOtpExistsFailed() {
	t.otpRepository.EXPECT().Get(mock.Anything).Return(&entities.Otp{
		EmailAddress: "",
		PasswordHash: []byte{},
		Salt:         []byte{},
		ExpiresAt:    time.Now().Add(time.Minute),
	}, nil)
	result := commands.GenerateOtpCommand{
		Email:     "p6bqK@example.com",
		TTLMillis: 0,
	}.Execute(t.ctx)

	assert.Equal(t.T(), common.ExecutionStatusFailed, result.Status)
	assert.Equal(t.T(), commands.GenerateOtpErrorOtpExists, result.Error.Code)
}

func TestGenerateOtp(t *testing.T) {
	suite.Run(t, new(GenerateOtpCommandTest))
}
