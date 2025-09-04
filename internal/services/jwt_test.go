package services

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWTService_GenerateParse(t *testing.T) {
	svc := NewJWTService("secret", 1)
	uid := uuid.New()
	tok, err := svc.Generate(uid)
	if err != nil || tok == "" {
		t.Fatalf("generate: %v, tok=%q", err, tok)
	}
	claims, err := svc.Parse(tok)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if claims.UserID != uid.String() {
		t.Fatalf("expected uid %s, got %s", uid, claims.UserID)
	}
}

func TestJWTService_Expired(t *testing.T) {
	svc := NewJWTService("secret", 0)
	uid := uuid.New()
	tok, err := svc.Generate(uid)
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	time.Sleep(2 * time.Second)
	_, err = svc.Parse(tok)
	if err == nil {
		t.Fatalf("expected error for expired token")
	}
} 