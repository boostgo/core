package authx

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"strings"

	"github.com/boostgo/core/convert"
	"github.com/golang-jwt/jwt/v4"
)

type ClaimsParser[T jwt.Claims] func(claims jwt.MapClaims) (T, error)

type JwtParser[T jwt.Claims] struct {
	secret       []byte
	algorithm    string
	publicKey    *rsa.PublicKey
	parser       *jwt.Parser
	claimsParser ClaimsParser[T]
}

func NewJwtParser[T jwt.Claims](
	secret, algorithm, publicKey string,
	claimsParser ClaimsParser[T],
) (*JwtParser[T], error) {
	publicKey = strings.ReplaceAll(publicKey, "\\n", "\n")

	parsedPublicKey, err := parseRSAPublicKey(publicKey)
	if err != nil {
		return nil, err
	}

	return &JwtParser[T]{
		secret:       convert.BytesFromString(secret),
		algorithm:    algorithm,
		publicKey:    parsedPublicKey,
		parser:       jwt.NewParser(jwt.WithoutClaimsValidation()),
		claimsParser: claimsParser,
	}, nil
}

func MustParser[T jwt.Claims](
	secret, algorithm, publicKey string,
	claimsParser ClaimsParser[T],
) *JwtParser[T] {
	parser, err := NewJwtParser[T](secret, algorithm, publicKey, claimsParser)
	if err != nil {
		panic(err)
	}

	return parser
}

func (p JwtParser[T]) Parse(tokenString string) (T, error) {
	var empty T

	token, err := p.parser.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		switch p.algorithm {
		case "HS256":
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, NewUnexpectedSigningMethodError(JwtSignMethodHS256, convert.String(token.Header["alg"]))
			}

			return p.secret, nil

		case "RS256":
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, NewUnexpectedSigningMethodError(JwtSignMethodRS256, convert.String(token.Header["alg"]))
			}

			if p.publicKey != nil {
				return p.publicKey, nil
			}

			return nil, ErrNoPublicKeyRS256

		default:
			return nil, ErrUnsupportedSigningMethod
		}
	})
	if err != nil {
		return empty, NewParseTokenError(err, tokenString)
	}

	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return empty, ErrInvalidClaims
	}

	return p.claimsParser(mapClaims)
}

func parseRSAPublicKey(publicKeyPEM string) (*rsa.PublicKey, error) {
	if publicKeyPEM == "" {
		return nil, ErrNoPublicKeyRS256.AddParam("stage", "validation")
	}

	wrappedPublicKey := "-----BEGIN PUBLIC KEY-----\n" + publicKeyPEM + "\n-----END PUBLIC KEY-----\n"

	block, _ := pem.Decode(convert.BytesFromString(wrappedPublicKey))
	if block == nil {
		return nil, ErrFailedToDecodePemBlock
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, ErrParsePkixPublicKey.SetError(err)
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, ErrNoPublicKeyRS256.AddParam("stage", "returning")
	}

	return rsaPub, nil
}
