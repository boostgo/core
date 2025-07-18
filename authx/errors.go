package authx

import "github.com/boostgo/core/errorx"

var (
	ErrParseToken               = errorx.New("authx.jwt.parse_token").SetError(errorx.ErrUnauthorized)
	ErrInvalidClaims            = errorx.New("authx.jwt.invalid_claims")
	ErrUnexpectedSigningMethod  = errorx.New("authx.jwt.unexpected_signing_method")
	ErrNoPublicKeyRS256         = errorx.New("authx.jwt.no_public_key_rs256")
	ErrUnsupportedSigningMethod = errorx.New("authx.jwt.unsupported_signing_method")
	ErrFailedToDecodePemBlock   = errorx.New("authx.jwt.failed_to_decode_pem_block")
	ErrParsePkixPublicKey       = errorx.New("authx.jwt.parse_pkix_public_key")

	ErrNoToken           = errorx.New("auth.no_token").SetError(errorx.ErrUnauthorized)
	ErrNoAccess          = errorx.New("auth.no_access").SetError(errorx.ErrForbidden)
	ErrNoGroups          = errorx.New("auth.no_groups").SetError(errorx.ErrUnauthorized)
	ErrInvalidGroups     = errorx.New("auth.invalid_groups").SetError(errorx.ErrUnauthorized)
	ErrResourceForbidden = errorx.New("auth.resource_forbidden").SetError(errorx.ErrForbidden)
)

type parseTokenContext struct {
	Token string `json:"token"`
}

func NewParseTokenError(err error, token string) error {
	return ErrParseToken.
		SetError(err).
		SetData(parseTokenContext{
			Token: token,
		})
}

type unexpectedSigningMethodContext struct {
	Expected JwtSignMethod `json:"expected"`
	Actual   string        `json:"actual"`
}

func NewUnexpectedSigningMethodError(expected JwtSignMethod, actual string) error {
	return ErrUnexpectedSigningMethod.SetData(unexpectedSigningMethodContext{
		Expected: expected,
		Actual:   actual,
	})
}
