package commands

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/Goldwin/ies-pik-cms/pkg/auth/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/auth/entities"

	. "github.com/Goldwin/ies-pik-cms/pkg/common/commands"
)

type GenerateResetTokenCommand struct {
	Email     string
	TTLMillis int64
}

func (cmd GenerateResetTokenCommand) Execute(ctx CommandContext) CommandExecutionResult[dto.PasswordResetTokenResult] {
	//30 seconds minimum
	ttlMillis := max(cmd.TTLMillis, 300000)
	token, err := rand.Int(rand.Reader, big.NewInt(999999))
	if err != nil {
		return CommandExecutionResult[dto.PasswordResetTokenResult]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    GenerateOtpErrorFailedToGenOtp,
				Message: fmt.Sprintf("Failed to Generate Token: %s", err.Error()),
			},
		}
	}

	strToken := strconv.Itoa(int(token.Int64()))
	err = ctx.PasswordRepository().SaveResetToken(entities.EmailAddress(cmd.Email), strToken, time.Duration(ttlMillis)*time.Millisecond)
	if err != nil {
		return CommandExecutionResult[dto.PasswordResetTokenResult]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    GenerateOtpErrorFailedToStoreOtp,
				Message: fmt.Sprintf("Failed to Generate OTP: %s", err.Error()),
			},
		}
	}
	return CommandExecutionResult[dto.PasswordResetTokenResult]{Status: ExecutionStatusSuccess, Result: dto.PasswordResetTokenResult{
		Email: cmd.Email, Token: strToken,
	}}
}
