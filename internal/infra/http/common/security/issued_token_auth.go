package security

import (
	"regexp"

	"github.com/eterline/fstmon/pkg/tokenissuer"
)

type IssuedTokenAuthProvide struct {
	enable bool
	issuer *tokenissuer.TokenIssuer
	re     *regexp.Regexp
}

func NewIssuedTokenAuthProvide(enable bool) (*IssuedTokenAuthProvide, error) {
	var (
		err error
		i   *tokenissuer.TokenIssuer
	)

	if enable {
		i, err = tokenissuer.NewTokenIssuerEnv("FSTMON_SECRET", tokenissuer.TokenMedium)
		if err != nil {
			return nil, err
		}
	}

	itap := &IssuedTokenAuthProvide{
		issuer: i,
		enable: enable,
		re:     regexp.MustCompile(`^Bearer:\s*(.+)`),
	}

	return itap, nil
}

/*
TestBearer – validates a Bearer token extracted from the Authorization header.

	Comparison is performed in constant time using subtle.ConstantTimeCompare to prevent timing attacks.
	If the provider is nil, authentication is considered disabled and always returns true.
*/
func (tap *IssuedTokenAuthProvide) TestBearer(bearer string) bool {
	if tap == nil {
		return false
	}

	if !tap.Enabled() {
		return true
	}

	mch := tap.re.FindStringSubmatch(bearer)
	if len(mch) == 0 {
		return false
	}

	ok, _ := tap.issuer.TestTokenString(mch[1])
	return ok
}

// Enabled – returns whether token authentication is enabled.
// Authentication is enabled when the provider is non-nil.
func (tap *IssuedTokenAuthProvide) Enabled() bool {
	return tap != nil && tap.enable
}

func (tap *IssuedTokenAuthProvide) Issue() (string, error) {
	return tap.issuer.TokenString()
}
