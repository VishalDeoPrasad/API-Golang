package auth

import (
	"crypto/rsa"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

type ctxKey int

const Key ctxKey = 1

type Auth struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func NewAuth(privKey *rsa.PrivateKey, pubKey *rsa.PublicKey) (*Auth, error) {
	if privKey == nil || pubKey == nil {
		err := errors.New("private/public key is not present")
		return nil, err
	}
	return &Auth{
		privateKey: privKey,
		publicKey:  pubKey,
	}, nil
}

func (a *Auth) GenerateToken(claims jwt.RegisteredClaims) (string, error) {
	//NewWithClaims creates a new Token with the specified signing method and claims.
	tkn := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenStr, err := tkn.SignedString(a.privateKey)
	if err != nil {
		return "", fmt.Errorf("signing token %w", err)
	}

	return tokenStr, nil
}

func (a *Auth) ValidateToken(token string) (jwt.RegisteredClaims, error) {
	var c jwt.RegisteredClaims
	// Parse the token with the registered claims.
	tkn, err := jwt.ParseWithClaims(token, &c, func(token *jwt.Token) (interface{}, error) {
		return a.publicKey, nil
	})
	if err != nil {
		return jwt.RegisteredClaims{}, fmt.Errorf("parsing token %w", err)
	}
	// Check if the parsed token is valid.
	if !tkn.Valid {
		return jwt.RegisteredClaims{}, errors.New("invalid token")
	}
	return c, nil
}
