package jwt

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
)

const APPLICATION_TYPE = 2

var (
	ErrParseToken      = errors.New("failed to parse authorization token")
	ErrClaimsAssertion = errors.New("unable to extract token claims")
)

// TokenClaims represents the role claims present in the JWT token. This structure extends
// the jwt.StandardClaims to facilitate the extraction of claims when parsing a JWT token.
// `Kind` represents the token type in the Mainflux platform and the other properties are
// part of the RFC specification: https://datatracker.ietf.org/doc/html/rfc7519.
type TokenClaims struct {
	Iss  string `json:"iss"`
	Sub  string `json:"sub"`
	Kind int    `json:"type"`
	jwt.StandardClaims
}

// GetEmail returns the email of the responsible for creating the JWT token, which can be of
// 'app' or 'user' type. Depending on this, the e-mail can be extract from the `Iss` role claim
// or the `Sub` role claim.
func GetEmail(token string) (string, error) {
	parser := new(jwt.Parser)
	// second return value can be ignore since it's only the individual parts of the token
	t, _, err := parser.ParseUnverified(token, &TokenClaims{})
	if err != nil {
		return "", ErrParseToken
	}

	claims, ok := t.Claims.(*TokenClaims)
	if !ok {
		return "", ErrClaimsAssertion
	}

	if claims.Kind == APPLICATION_TYPE {
		return claims.Iss, nil
	}

	return claims.Sub, nil
}
