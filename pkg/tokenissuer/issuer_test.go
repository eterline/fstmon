package tokenissuer_test

import (
	"encoding/base64"
	"os"
	"strings"
	"testing"

	"github.com/eterline/dumbsec/tokenissuer"
)

func Test_NewTokenIssuer_Valid(t *testing.T) {
	secret := []byte("supersecret")
	sizes := []tokenissuer.TokenLen{
		tokenissuer.TokenShort,
		tokenissuer.TokenMedium,
		tokenissuer.TokenLong,
		tokenissuer.TokenMax,
	}

	for _, size := range sizes {
		ti, err := tokenissuer.NewTokenIssuer(secret, size)
		if err != nil {
			t.Fatalf("expected no error for size %d, got %v", size, err)
		}
		if ti == nil {
			t.Fatalf("expected non-nil TokenIssuer for size %d", size)
		}
	}
}

func Test_NewTokenIssuer_InvalidSizes(t *testing.T) {
	secret := []byte("secret")

	if _, err := tokenissuer.NewTokenIssuer(secret, 0); err == nil {
		t.Fatal("expected error for token size 0")
	}

	if _, err := tokenissuer.NewTokenIssuer(secret, 7); err == nil {
		t.Fatal("expected error for token size not divisible by 4")
	}

	if _, err := tokenissuer.NewTokenIssuer(secret, 4); err == nil {
		t.Fatal("expected error for too small token size")
	}
}

func Test_NewTokenIssuerEnv(t *testing.T) {
	envName := "TEST_SECRET"
	os.Setenv(envName, "envsecret")
	defer os.Unsetenv(envName)

	ti, err := tokenissuer.NewTokenIssuerEnv(envName, tokenissuer.TokenShort)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if ti == nil {
		t.Fatal("expected non-nil TokenIssuer")
	}

	_, err = tokenissuer.NewTokenIssuerEnv("NON_EXISTENT_ENV", tokenissuer.TokenShort)
	if err == nil {
		t.Fatal("expected error for missing environment variable")
	}
}

func Test_TokenAndTestToken(t *testing.T) {
	secret := []byte("supersecret")
	ti, _ := tokenissuer.NewTokenIssuer(secret, tokenissuer.TokenMedium)

	tokens := make(map[string]struct{})

	for i := 0; i < 100; i++ {
		raw, err := ti.Token()
		if err != nil {
			t.Fatalf("unexpected error generating token: %v", err)
		}
		ok, err := ti.TestToken(raw)
		if err != nil {
			t.Fatalf("unexpected error testing token: %v", err)
		}
		if !ok {
			t.Errorf("valid token did not validate")
		}

		str := base64.RawURLEncoding.EncodeToString(raw)
		tokens[str] = struct{}{}
	}

	if len(tokens) != 100 {
		t.Errorf("expected 100 unique tokens, got %d", len(tokens))
	}
}

func Test_TokenStringAndTestTokenString(t *testing.T) {
	secret := []byte("anothersecret")
	ti, _ := tokenissuer.NewTokenIssuer(secret, tokenissuer.TokenLong)

	for i := 0; i < 50; i++ {
		str, err := ti.TokenString()
		if err != nil {
			t.Fatalf("unexpected error generating token string: %v", err)
		}

		ok, err := ti.TestTokenString(str)
		if err != nil {
			t.Fatalf("unexpected error testing token string: %v", err)
		}
		if !ok {
			t.Errorf("valid token string did not validate")
		}
	}
}

func Test_InvalidTokens(t *testing.T) {
	secret := []byte("secretkey")
	ti, _ := tokenissuer.NewTokenIssuer(secret, tokenissuer.TokenShort)

	// wrong length byte slice
	ok, err := ti.TestToken([]byte{0x01, 0x02})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Errorf("invalid token unexpectedly validated")
	}

	// modify valid token
	validRaw, _ := ti.Token()
	validRaw[0] ^= 0xFF
	ok, err = ti.TestToken(validRaw)
	if err != nil {
		t.Fatalf("unexpected error testing modified token: %v", err)
	}
	if ok {
		t.Errorf("tampered token unexpectedly validated")
	}

	// invalid base64 string
	ok, err = ti.TestTokenString("invalid@@@")
	if err == nil {
		t.Errorf("expected error decoding invalid base64 string")
	}
	if ok {
		t.Errorf("invalid string unexpectedly validated")
	}

	// tampered token string
	validStr, _ := ti.TokenString()
	tampered := "A" + validStr[1:]
	ok, err = ti.TestTokenString(tampered)
	if err != nil {
		t.Fatalf("unexpected error testing tampered token string: %v", err)
	}
	if ok {
		t.Errorf("tampered token string unexpectedly validated")
	}
}

func Test_NilTokenIssuer(t *testing.T) {
	var ti *tokenissuer.TokenIssuer

	if _, err := ti.Token(); err != tokenissuer.ErrIssuerNil {
		t.Errorf("expected ErrIssuerNil, got %v", err)
	}

	if _, err := ti.TokenString(); err != tokenissuer.ErrIssuerNil {
		t.Errorf("expected ErrIssuerNil, got %v", err)
	}

	ok, err := ti.TestToken([]byte{0x01})
	if err != tokenissuer.ErrIssuerNil || ok {
		t.Errorf("expected ErrIssuerNil and false, got %v, %v", err, ok)
	}

	ok, err = ti.TestTokenString("abcd")
	if err != tokenissuer.ErrIssuerNil || ok {
		t.Errorf("expected ErrIssuerNil and false, got %v, %v", err, ok)
	}
}

func Test_Base64RoundTrip(t *testing.T) {
	secret := []byte("roundtripsecret")
	ti, _ := tokenissuer.NewTokenIssuer(secret, tokenissuer.TokenMedium)

	raw, _ := ti.Token()
	str := base64.RawURLEncoding.EncodeToString(raw)
	decoded, err := base64.RawURLEncoding.DecodeString(str)
	if err != nil {
		t.Fatalf("unexpected error decoding base64: %v", err)
	}
	if string(decoded) != string(raw) {
		t.Errorf("decoded bytes do not match original")
	}
}

func Test_MultipleTokenLengths(t *testing.T) {
	secret := []byte("multiplesecret")
	lengths := []tokenissuer.TokenLen{
		tokenissuer.TokenShort,
		tokenissuer.TokenMedium,
		tokenissuer.TokenLong,
		tokenissuer.TokenMax,
	}
	for _, l := range lengths {
		ti, err := tokenissuer.NewTokenIssuer(secret, l)
		if err != nil {
			t.Fatalf("unexpected error creating TokenIssuer: %v", err)
		}
		str, err := ti.TokenString()
		if err != nil {
			t.Fatalf("unexpected error generating token string: %v", err)
		}
		ok, err := ti.TestTokenString(str)
		if err != nil {
			t.Fatalf("unexpected error testing token string: %v", err)
		}
		if !ok {
			t.Errorf("valid token did not validate")
		}
	}
}

func Test_PayloadVariation(t *testing.T) {
	secret := []byte("variationsecret")
	ti, _ := tokenissuer.NewTokenIssuer(secret, tokenissuer.TokenMedium)

	prev := make(map[string]struct{})
	for i := 0; i < 100; i++ {
		raw, _ := ti.Token()
		str := base64.RawURLEncoding.EncodeToString(raw)
		if _, exists := prev[str]; exists {
			t.Errorf("duplicate token detected: %s", str)
		}
		prev[str] = struct{}{}
	}
}

func Test_TokenLengthConsistency(t *testing.T) {
	secret := []byte("lengthsecret")
	ti, _ := tokenissuer.NewTokenIssuer(secret, tokenissuer.TokenLong)

	for i := 0; i < 50; i++ {
		str, _ := ti.TokenString()
		if strings.TrimSpace(str) == "" {
			t.Errorf("token string is empty or whitespace")
		}
	}
}
