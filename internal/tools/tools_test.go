package tools_test

import (
	"testing"

	"github.com/xantinium/gophermart/internal/tools"
)

func TestCheckLuhn(t *testing.T) {
	positiveCases := []string{
		"1104942",
		"4561261212345467",
	}
	negativeCases := []string{
		"",
		"2t4",
		"11045",
		"4561261212345464",
	}

	for _, tt := range positiveCases {
		if got := tools.CheckLuhn(tt); !got {
			t.Errorf("CheckLuhn() = %v, want %v", got, true)
		}
	}
	for _, tt := range negativeCases {
		if got := tools.CheckLuhn(tt); got {
			t.Errorf("CheckLuhn() = %v, want %v", got, false)
		}
	}
}

func TestHashPassword(t *testing.T) {
	password := "some_password_123"

	hashedPassword, err := tools.HashPassword(password)
	if err != nil {
		t.Errorf("failed to hash password: %v", err)
	}

	if !tools.CheckPassword(password, hashedPassword) {
		t.Error("hash does not match")
	}
}
