package auth

import (
	"crypto/rsa"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

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
