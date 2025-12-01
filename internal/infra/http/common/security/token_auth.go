// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.
package security

import (
	"crypto/subtle"
	"errors"
	"regexp"
)

type tokenBearerLenPolicy int

const (
	_ tokenBearerLenPolicy = 1 << (iota + 4)
	PolicyLow
	PolicyStart
	PolicyMid
	PolicyStrong
	PolicyOver
)

// TokenAuthProvide – provides constant-time Bearer token authentication
// using an internal pool of preconfigured tokens.
type TokenAuthProvide struct {
	enable    bool
	tokenPool [][]byte
	re        *regexp.Regexp
}

/*
NewTokenAuthProvide – creates a new TokenAuthProvide instance.

	Accepts one or more Bearer tokens. Each token must be at least minLen bytes,
	otherwise an error is returned.
	If no tokens are provided, the function returns (nil, nil) meaning
	authentication is disabled.
*/
func NewTokenAuthProvide(minLen tokenBearerLenPolicy, token ...string) (*TokenAuthProvide, error) {
	if len(token) == 0 {
		return &TokenAuthProvide{}, nil
	}

	p := &TokenAuthProvide{
		enable:    true,
		tokenPool: make([][]byte, 0, len(token)),
		re:        regexp.MustCompile(`^Bearer:\s*(.+)`),
	}

	for _, t := range token {
		tokenBytes := []byte(t)

		if len(tokenBytes) < int(minLen) {
			return nil, errors.New("invalid token – length must be above or eq 64 bytes")
		}

		p.tokenPool = append(p.tokenPool, tokenBytes)
	}

	return p, nil
}

/*
TestBearer – validates a Bearer token extracted from the Authorization header.

	Comparison is performed in constant time using subtle.ConstantTimeCompare to prevent timing attacks.
	If the provider is nil, authentication is considered disabled and always returns true.
*/
func (tap *TokenAuthProvide) TestBearer(bearer string) bool {
	if tap == nil {
		return false
	}

	if !tap.enable {
		return true
	}

	mch := tap.re.FindSubmatch([]byte(bearer))
	if len(mch) == 0 {
		return false
	}

	token := mch[1]

	for _, compareWith := range tap.tokenPool {
		if len(token) != len(compareWith) {
			continue
		}

		if subtle.ConstantTimeCompare(token, compareWith) == 1 {
			return true
		}
	}

	return false
}

// Enabled – returns whether token authentication is enabled.
// Authentication is enabled when the provider is non-nil.
func (tap *TokenAuthProvide) Enabled() bool {
	return tap == nil
}
