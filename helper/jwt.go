package helper

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

type JWTHelper interface {
	CreateAndSign(customClaimBytes []byte, expiredAt int64) (string, error)
	ParseAndVerify(signed string) ([]byte, error)
}

type JwtConfig struct {
	Issuer string `json:"issuer"`
	Key    string `json:"key"`
}

type jwtHelperImpl struct {
	key    string
	issuer string
	method jwt.SigningMethod
}

func NewJWTHelper(config JwtConfig, method jwt.SigningMethod) *jwtHelperImpl {
	return &jwtHelperImpl{
		key:    config.Key,
		issuer: config.Issuer,
		method: method,
	}
}

func (h *jwtHelperImpl) CreateAndSign(customClaimBytes []byte, expiredAt int64) (string, error) {
	token := jwt.NewWithClaims(h.method, jwt.MapClaims{
		"iss":  h.issuer,
		"exp":  expiredAt,
		"data": string(customClaimBytes),
	})

	signed, err := token.SignedString([]byte(h.key))
	if err != nil {
		return "", err
	}

	return signed, nil
}

func (h *jwtHelperImpl) ParseAndVerify(signed string) ([]byte, error) {
	token, err := jwt.Parse(signed, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.key), nil
	}, jwt.WithValidMethods([]string{h.method.Alg()}),
		jwt.WithIssuer(h.issuer),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenInvalidIssuer) {
			return nil, nil
		}

		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, err
	}

	data := claims["data"].(string)

	return []byte(data), nil
}
