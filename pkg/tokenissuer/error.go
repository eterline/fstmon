package tokenissuer

import "fmt"

type IssuerError struct {
	embed string
	err   error
}

func (ie *IssuerError) Error() string {
	if ie.err == nil {
		return fmt.Sprintf("issuer error: %s", ie.embed)
	}
	return fmt.Sprintf("issuer error: %s: %v", ie.embed, ie.err)
}

func (ie *IssuerError) Wrap(err error) *IssuerError {
	ie.err = err
	return ie
}

func (ie *IssuerError) Unwrap() error {
	return ie.err
}

func newIssuerError(embed string) *IssuerError {
	return &IssuerError{embed: embed}
}

var (
	ErrIssuerNil           = newIssuerError("nil issuer")
	ErrIssuerBadTokenSize  = newIssuerError("invalid token size")
	ErrIssuerNotDiv4       = newIssuerError("token size must be multiple of 4")
	ErrIssuerTooSmall      = newIssuerError("tokenSize too small")
	ErrIssuerRandomFail    = newIssuerError("failed to read random bytes")
	ErrIssuerDecodeFail    = newIssuerError("base64 decode error")
	ErrIssuerWrongTokenLen = newIssuerError("invalid token length")
	ErrSecretEnvNotExists  = newIssuerError("secret env not exists")
)
