package services

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
	"time"

	redisv9 "github.com/redis/go-redis/v9"
)

type OTPService struct {
	redis  *redisv9.Client
	prefix string
	ttl    time.Duration
}

func NewOTPService(client *redisv9.Client, ttlSeconds int) *OTPService {
	return &OTPService{redis: client, prefix: "otp:", ttl: time.Duration(ttlSeconds) * time.Second}
}

func (s *OTPService) key(phone string) string {
	return s.prefix + strings.TrimSpace(phone)
}

func (s *OTPService) Generate(ctx context.Context, phone string) (string, error) {
	// 6-digit numeric OTP
	n := int64(100000)
	m := int64(900000)
	r, err := rand.Int(rand.Reader, big.NewInt(m))
	if err != nil {
		return "", err
	}
	code := fmt.Sprintf("%06d", r.Int64()+n)
	if err := s.redis.Set(ctx, s.key(phone), code, s.ttl).Err(); err != nil {
		return "", err
	}
	return code, nil
}

func (s *OTPService) Verify(ctx context.Context, phone, code string) (bool, error) {
	stored, err := s.redis.Get(ctx, s.key(phone)).Result()
	if err == redisv9.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if stored == strings.TrimSpace(code) {
		// consume OTP
		_ = s.redis.Del(ctx, s.key(phone)).Err()
		return true, nil
	}
	return false, nil
}
