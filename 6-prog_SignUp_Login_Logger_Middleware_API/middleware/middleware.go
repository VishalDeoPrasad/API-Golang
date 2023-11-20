package middleware

import (
	"errors"
	"golang/auth"
)

type Mid struct {
	a *auth.Auth
}

func NewMid(a *auth.Auth) (*Mid, error) {
    if a == nil {
        return nil, errors.New("auth can't be nil")
    }
    return &Mid{a: a}, nil
}
