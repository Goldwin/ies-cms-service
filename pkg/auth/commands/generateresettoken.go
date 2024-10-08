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

const (
	ResetPasswordErrorFailedToGenOtp    CommandErrorCode = 20401
	ResetPasswordErrorAccountNotFound   CommandErrorCode = 20402
	ResetPasswordErrorCodeAlreadyExists CommandErrorCode = 20402
)

func (cmd GenerateResetTokenCommand) maybeDeleteToken(ctx CommandContext) *entities.PasswordResetCode {
	originalToken, err := ctx.PasswordResetCodeRepository().Get(cmd.Email)

	if err != nil || originalToken == nil {
		return nil
	}

	if time.Now().After(originalToken.ExpiresAt) {
		ctx.PasswordResetCodeRepository().Delete(originalToken)
		return nil
	}

	return originalToken
}

func (cmd GenerateResetTokenCommand) Execute(ctx CommandContext) CommandExecutionResult[dto.PasswordResetCodeResult] {
	originalToken := cmd.maybeDeleteToken(ctx)

	if originalToken != nil {
		return CommandExecutionResult[dto.PasswordResetCodeResult]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    ResetPasswordErrorCodeAlreadyExists,
				Message: fmt.Sprintf("Password Reset Code has been sent. Please try again after %.0f seconds.", originalToken.ExpiresAt.Sub(time.Now()).Seconds()),
			},
		}
	}

	//30 seconds minimum
	ttlMillis := max(cmd.TTLMillis, 300000)
	token, err := rand.Int(rand.Reader, big.NewInt(999999))
	if err != nil {
		return CommandExecutionResult[dto.PasswordResetCodeResult]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    GenerateOtpErrorFailedToGenOtp,
				Message: fmt.Sprintf("Failed to Generate Token: %s", err.Error()),
			},
		}
	}

	password, err := ctx.PasswordRepository().Get(cmd.Email)

	if err != nil || password == nil {
		return CommandExecutionResult[dto.PasswordResetCodeResult]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    ResetPasswordErrorAccountNotFound,
				Message: fmt.Sprintf("Failed to Get Account Detail: %s", cmd.Email),
			},
		}
	}

	strToken := strconv.Itoa(int(token.Int64()))
	_, err = ctx.PasswordResetCodeRepository().Save(&entities.PasswordResetCode{
		Email:     cmd.Email,
		Code:      strToken,
		ExpiresAt: time.Now().Add(time.Duration(ttlMillis) * time.Millisecond),
	})

	if err != nil {
		return CommandExecutionResult[dto.PasswordResetCodeResult]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    ResetPasswordErrorFailedToGenOtp,
				Message: fmt.Sprintf("Failed to Generate OTP: %s", err.Error()),
			},
		}
	}

	return CommandExecutionResult[dto.PasswordResetCodeResult]{Status: ExecutionStatusSuccess, Result: dto.PasswordResetCodeResult{
		Email: cmd.Email, Code: strToken,
	}}
}
