package token

import (
	"fmt"

	"crypto/ecdsa"

	"github.com/golang-jwt/jwt"
)

func ValidateJWTToken(publicKey *ecdsa.PublicKey, tokenString string) (jwt.MapClaims, error) {
	jwtToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil {
		return nil, err
	}
	if jwtToken == nil {
		return nil, ErrInvalidToken
	}

	if claims, ok := jwtToken.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}

	return nil, ErrInvalidToken
}
