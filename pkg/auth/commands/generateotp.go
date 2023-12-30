package commands

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
	"time"

	"github.com/Goldwin/ies-pik-cms/pkg/auth/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/auth/entities"
	"github.com/Goldwin/ies-pik-cms/pkg/auth/repositories"
	. "github.com/Goldwin/ies-pik-cms/pkg/common/commands"
)

type GenerateOtpCommand struct {
	Email     string
	TTLMillis int64
}

const (
	GenerateOtpErrorInvalidEmail     AppErrorCode = 20001
	GenerateOtpErrorFailedToGenOtp   AppErrorCode = 20002
	GenerateOtpErrorFailedToStoreOtp AppErrorCode = 20003
	GenerateOtpErrorOtpExists        AppErrorCode = 20004
)

func (cmd GenerateOtpCommand) Execute(ctx repositories.CommandContext) AppExecutionResult[dto.OtpResult] {
	otp, _ := ctx.OtpRepository().GetOtp(entities.EmailAddress(cmd.Email))
	if otp != nil {
		expireSecond := otp.ExpiredTime.Sub(time.Now()).Seconds()
		if expireSecond > 0 {
			return AppExecutionResult[dto.OtpResult]{
				Status: ExecutionStatusFailed,
				Error: AppErrorDetail{
					Code:    GenerateOtpErrorOtpExists,
					Message: fmt.Sprintf("OTP Already Exists. Please wait for %.0f seconds and try again.", expireSecond),
				},
			}
		}
	}
	//30 seconds minimum
	ttlMillis := max(cmd.TTLMillis, 30000)
	password, err := rand.Int(rand.Reader, big.NewInt(999999))
	if err != nil {
		return AppExecutionResult[dto.OtpResult]{
			Status: ExecutionStatusFailed,
			Error: AppErrorDetail{
				Code:    GenerateOtpErrorFailedToGenOtp,
				Message: fmt.Sprintf("Failed to Generate OTP: %s", err.Error()),
			},
		}
	}

	salt, err := rand.Int(rand.Reader, big.NewInt(999999))
	if err != nil {
		return AppExecutionResult[dto.OtpResult]{
			Status: ExecutionStatusFailed,
			Error: AppErrorDetail{
				Code:    GenerateOtpErrorFailedToGenOtp,
				Message: fmt.Sprintf("Failed to Generate OTP: %s", err.Error()),
			},
		}
	}

	n := password.Int64()
	passwordBytes := []byte(fmt.Sprintf("%06v", n))
	passwordAndSalt := append(passwordBytes, salt.Bytes()...)
	passwordHash := sha256.Sum256(passwordAndSalt)

	result := entities.Otp{
		EmailAddress: entities.EmailAddress(cmd.Email),
		PasswordHash: passwordHash[:],
		Salt:         salt.Bytes(),
		ExpiredTime:  time.Now().Add(time.Duration(ttlMillis) * time.Millisecond),
	}

	if !result.EmailAddress.IsValid() {
		return AppExecutionResult[dto.OtpResult]{
			Status: ExecutionStatusFailed,
			Error: AppErrorDetail{
				Code:    GenerateOtpErrorInvalidEmail,
				Message: "Invalid Email Address",
			},
		}
	}

	err = ctx.OtpRepository().AddOtp(result)
	if err != nil {
		return AppExecutionResult[dto.OtpResult]{
			Status: ExecutionStatusFailed,
			Error: AppErrorDetail{
				Code:    GenerateOtpErrorFailedToStoreOtp,
				Message: fmt.Sprintf("Failed to Generate OTP: %s", err.Error()),
			},
		}
	}
	return AppExecutionResult[dto.OtpResult]{Status: ExecutionStatusSuccess, Result: dto.OtpResult{
		Email: cmd.Email,
		OTP:   passwordBytes,
	}}
}
