package services

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	redisv9 "github.com/redis/go-redis/v9"
)

func TestOTPService_GenerateAndVerify(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("miniredis: %v", err)
	}
	defer mr.Close()

	rdb := redisv9.NewClient(&redisv9.Options{Addr: mr.Addr()})
	svc := NewOTPService(rdb, 1, 3, 60) // for line 53
	ctx := context.Background()

	code, err := svc.Generate(ctx, "+15551234567")
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	if code == "" || len(code) != 6 {
		t.Fatalf("expected 6-digit code, got %q", code)
	}

	ok, err := svc.Verify(ctx, "+15551234567", code)
	if err != nil || !ok {
		t.Fatalf("verify expected true, got %v, err=%v", ok, err)
	}
	// second verify should fail (consumed)
	ok, err = svc.Verify(ctx, "+15551234567", code)
	if err != nil {
		t.Fatalf("verify second err: %v", err)
	}
	if ok {
		t.Fatalf("expected false after consume")
	}
}

func TestOTPService_TTL(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("miniredis: %v", err)
	}
	defer mr.Close()

	rdb := redisv9.NewClient(&redisv9.Options{Addr: mr.Addr()})
	svc := NewOTPService(rdb, 1, 3, 60) // for line 53
	ctx := context.Background()

	code, err := svc.Generate(ctx, "+15557654321")
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	// advance miniredis time to trigger TTL expiry
	mr.FastForward(2 * time.Second)
	ok, err := svc.Verify(ctx, "+15557654321", code)
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	if ok {
		t.Fatalf("expected false after expiry")
	}
}
