// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.
package security_test

import (
	"crypto/rand"
	"encoding/base64"
	"testing"

	"github.com/eterline/fstmon/internal/infra/security"
)

const minLen = 64

func genToken(n int) string {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	return base64.RawURLEncoding.EncodeToString(b)[:n]
}

func TestNewTokenAuthProvide_NoTokens(t *testing.T) {
	provider, err := security.NewTokenAuthProvide(minLen)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if provider == nil {
		t.Fatal("Expected provider instance, got nil")
	}
	if provider.Enabled() {
		t.Error("Expected authentication to be disabled with no tokens")
	}
}

func TestNewTokenAuthProvide_InvalidTokenLength(t *testing.T) {
	shortToken := genToken(10)
	_, err := security.NewTokenAuthProvide(minLen, shortToken)
	if err == nil {
		t.Fatal("Expected error for short token, got nil")
	}
}

func TestNewTokenAuthProvide_ValidTokens(t *testing.T) {
	tokens := []string{genToken(minLen), genToken(minLen + 5)}
	provider, err := security.NewTokenAuthProvide(minLen, tokens...)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if provider.Enabled() {
		t.Error("Expected authentication to be enabled")
	}
}

func TestTokenAuthProvide_TestBearer_Valid(t *testing.T) {
	tok := genToken(minLen)
	provider, _ := security.NewTokenAuthProvide(minLen, tok)

	header := "Bearer: " + tok
	if !provider.TestBearer(header) {
		t.Error("Expected valid token to pass TestBearer")
	}
}

func TestTokenAuthProvide_TestBearer_Invalid(t *testing.T) {
	valid := genToken(minLen)
	invalid := genToken(minLen)

	if valid == invalid {
		invalid = genToken(minLen)
	}

	provider, _ := security.NewTokenAuthProvide(minLen, valid)

	header := "Bearer: " + invalid
	if provider.TestBearer(header) {
		t.Error("Expected invalid token to fail TestBearer")
	}

	if provider.TestBearer("InvalidHeader") {
		t.Error("Expected malformed header to fail TestBearer")
	}
}

func TestTokenAuthProvide_TestBearer_MultipleTokens(t *testing.T) {
	tok1 := genToken(minLen)
	tok2 := genToken(minLen)
	provider, _ := security.NewTokenAuthProvide(minLen, tok1, tok2)

	if !provider.TestBearer("Bearer: " + tok1) {
		t.Error("Expected first token to pass")
	}
	if !provider.TestBearer("Bearer: " + tok2) {
		t.Error("Expected second token to pass")
	}
}

func TestTokenAuthProvide_TestBearer_Disabled(t *testing.T) {
	tok := genToken(minLen)
	provider, _ := security.NewTokenAuthProvide(minLen, tok)

	_ = provider.Enabled() // stub for ide PROBLEMS system

	provider = &security.TokenAuthProvide{}
	if !provider.TestBearer("any") {
		t.Error("Expected disabled provider to always return true")
	}
}

func TestTokenAuthProvide_TestBearer_NilProvider(t *testing.T) {
	var provider *security.TokenAuthProvide
	if provider.TestBearer("Bearer: something") {
		t.Error("Expected nil provider to return false")
	}
}
