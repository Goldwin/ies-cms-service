package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/Goldwin/ies-pik-cms/pkg/auth/entities"
	"github.com/Goldwin/ies-pik-cms/pkg/auth/repositories"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"github.com/vmihailenco/msgpack/v5"
)

type otpRepositoryImpl struct {
	ctx        context.Context
	txPipeline redis.Pipeliner
	client     redis.UniversalClient
}

// Delete implements repositories.OtpRepository.
func (o *otpRepositoryImpl) Delete(otp *entities.Otp) error {
	return o.txPipeline.Del(o.ctx, getOtpKey(string(otp.EmailAddress))).Err()
}

// Get implements repositories.OtpRepository.
func (o *otpRepositoryImpl) Get(email string) (*entities.Otp, error) {
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

// List implements repositories.OtpRepository.
func (o *otpRepositoryImpl) List(emails []string) ([]*entities.Otp, error) {
	var err error
	result := lo.Map(emails, func(e string, _ int) *entities.Otp {
		var otp entities.Otp
		key := getOtpKey(e)
		val, err2 := o.client.Get(o.ctx, key).Bytes()
		if err != nil {
			return nil
		}
		if err2 != nil && err2 != redis.Nil {
			err = err2
			return nil
		}
		err = msgpack.Unmarshal(val, &otp)
		return &otp
	})
	return result, err
}

// Save implements repositories.OtpRepository.
func (o *otpRepositoryImpl) Save(otp *entities.Otp) (*entities.Otp, error) {
	bytes, err := msgpack.Marshal(otp)
	if err != nil {
		return nil, err
	}
	key := getOtpKey(string(otp.EmailAddress))
	ttl := otp.ExpiresAt.Sub(time.Now())
	_, err = o.txPipeline.Set(o.ctx, key, string(bytes), ttl).Result()
	return otp, err
}

func NewOtpRepository(ctx context.Context, client redis.UniversalClient, txPipeline redis.Pipeliner) repositories.OtpRepository {
	return &otpRepositoryImpl{
		ctx:        ctx,
		client:     client,
		txPipeline: txPipeline,
	}
}

func getOtpKey(email string) string {
	return fmt.Sprintf("auth:otp:email#%s", email)
}
