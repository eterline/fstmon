package tokenissuer

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

// TokenLen represents predefined token sizes in bytes.
type TokenLen int

const (
	TokenShort  TokenLen = 32
	TokenMedium TokenLen = 48
	TokenLong   TokenLen = 64
	TokenMax    TokenLen = 128
)

// TokenIssuer represents a token generator and validator using HMAC-based authentication.
type TokenIssuer struct {
	secret     []byte
	tokenSize  int
	macLen     int
	payloadLen int
}

/*
NewTokenIssuerEnv creates a new TokenIssuer using a secret retrieved from an environment variable.

	Parameters:
	- envName: the name of the environment variable containing the secret.
	- tokenSize: predefined token length from tokenLen constants.

	Returns:
	- *TokenIssuer instance
	- error if environment variable is missing or token size is invalid
*/
func NewTokenIssuerEnv(envName string, tokenSize TokenLen) (*TokenIssuer, error) {
	rawSecret := os.Getenv(envName)
	if rawSecret == "" {
		return nil, ErrSecretEnvNotExists
	}
	return NewTokenIssuer([]byte(rawSecret), tokenSize)
}

/*
NewTokenIssuer creates a new TokenIssuer instance using a given secret and token size.

	Parameters:
	- secret: byte slice containing the secret key.
	- tokenSize: predefined token length from tokenLen constants.

	Returns:
	- *TokenIssuer instance
	- error if token size is invalid
*/
func NewTokenIssuer(secret []byte, tokenSize TokenLen) (*TokenIssuer, error) {
	if tokenSize <= 0 {
		return nil, ErrIssuerTooSmall
	}

	if tokenSize%4 != 0 {
		return nil, ErrIssuerNotDiv4
	}

	combinedRaw := (tokenSize / 4) * 3
	const defaultMacLen = 16

	if combinedRaw <= defaultMacLen {
		return nil, ErrIssuerTooSmall.Wrap(fmt.Errorf("need > %d raw bytes", defaultMacLen))
	}

	ti := &TokenIssuer{
		// internal safety token store... It prevents source secret slice modify
		secret:     bytes.Clone(secret),
		tokenSize:  int(tokenSize),
		macLen:     defaultMacLen,
		payloadLen: int(combinedRaw) - defaultMacLen,
	}

	return ti, nil
}

/*
Token generates a random token as a byte slice with HMAC appended.

	Returns:
	- raw token bytes (payload + truncated HMAC)
	- error if token generation fails
*/
func (t *TokenIssuer) Token() (token []byte, err error) {
	if t == nil {
		return nil, ErrIssuerNil
	}

	payload := make([]byte, t.payloadLen)
	if _, err := io.ReadFull(rand.Reader, payload); err != nil {
		return nil, ErrIssuerRandomFail.Wrap(err)
	}

	m := hmac.New(sha512.New, t.secret)
	m.Write(payload)
	full := m.Sum(nil)
	mac := full[:t.macLen]

	raw := make([]byte, 0, t.payloadLen+t.macLen)
	raw = append(raw, payload...)
	raw = append(raw, mac...)

	return raw, nil
}

/*
TokenString generates a base64 string token.

	Returns:
	- token string
	- error if generation fails or length does not match expected tokenSize
*/
func (t *TokenIssuer) TokenString() (tokenBase64 string, err error) {
	raw, err := t.Token()
	if err != nil {
		return "", err
	}

	str := base64.RawURLEncoding.EncodeToString(raw)
	if len(str) != t.tokenSize {
		return "", ErrIssuerWrongTokenLen
	}

	return str, nil
}

/*
TestToken validates a raw token byte slice by comparing its HMAC.

	Parameters:
	- raw: token bytes (payload + truncated HMAC)

	Returns:
	- bool: true if token is valid
	- error: if TokenIssuer is nil or other internal error occurs
*/
func (t *TokenIssuer) TestToken(rawToken []byte) (ok bool, err error) {
	if t == nil {
		return false, ErrIssuerNil
	}

	if len(rawToken) != t.payloadLen+t.macLen {
		return false, nil
	}

	payload := rawToken[:t.payloadLen]
	macGiven := rawToken[t.payloadLen:]

	h := hmac.New(sha512.New, t.secret)
	h.Write(payload)
	expected := h.Sum(nil)[:t.macLen]

	if subtle.ConstantTimeCompare(macGiven, expected) == 1 {
		return true, nil
	}

	return false, nil
}

/*
TestTokenString validates a token string by decoding it from base64 and checking HMAC.

	Parameters:
	- token: base64 string

	Returns:
	- bool: true if token is valid
	- error: if decoding fails or internal validation fails
*/
func (t *TokenIssuer) TestTokenString(tokenBase64 string) (bool, error) {
	raw, err := base64.RawURLEncoding.DecodeString(tokenBase64)
	if err != nil {
		return false, ErrIssuerDecodeFail.Wrap(err)
	}
	return t.TestToken(raw)
}
