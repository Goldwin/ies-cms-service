package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/Goldwin/ies-pik-cms/pkg/auth/entities"
	"github.com/Goldwin/ies-pik-cms/pkg/auth/repositories"
	"github.com/redis/go-redis/v9"
	"github.com/vmihailenco/msgpack/v5"
)

type otpRepositoryImpl struct {
	ctx        context.Context
	txPipeline redis.Pipeliner
	client     redis.UniversalClient
}

// AddOtp implements repositories.OtpRepository.
func (o *otpRepositoryImpl) AddOtp(otp entities.Otp) error {
	bytes, err := msgpack.Marshal(otp)
	if err != nil {
		return err
	}
	key := getOtpKey(otp.EmailAddress)
	ttl := otp.ExpiredTime.Sub(time.Now())
	res, err := o.txPipeline.Set(o.ctx, key, string(bytes), ttl).Result()
	fmt.Printf("key: %s, ttl: %s, res: %s\n", key, ttl, res)
	return err
}

// GetOtp implements repositories.OtpRepository.
func (o *otpRepositoryImpl) GetOtp(email entities.EmailAddress) (*entities.Otp, error) {
	val, err := o.client.Get(o.ctx, getOtpKey(email)).Bytes()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	otp := entities.Otp{}
	if len(val) == 0 {
		return nil, nil
	}
	err = msgpack.Unmarshal(val, &otp)
	if err != nil {
		return nil, err
	}
	return &otp, nil
}

// RemoveOtp implements repositories.OtpRepository.
func (o *otpRepositoryImpl) RemoveOtp(otp entities.Otp) error {
	return o.txPipeline.Del(o.ctx, getOtpKey(otp.EmailAddress)).Err()
}

func NewOtpRepository(ctx context.Context, client redis.UniversalClient, txPipeline redis.Pipeliner) repositories.OtpRepository {
	return &otpRepositoryImpl{
		ctx:        ctx,
		client:     client,
		txPipeline: txPipeline,
	}
}

func getOtpKey(email entities.EmailAddress) string {
	return fmt.Sprintf("auth:otp:email#%s", email)
}
