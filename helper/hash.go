package helper

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type HashHelper interface {
	Hash(str string) (string, error)
	Check(str string, hash []byte) (bool, error)
}

type HashConfig struct {
	HashCost int `json:"hash_cost"`
}

type hashHelperImpl struct {
	hashCost int
}

func NewHashHelper(config HashConfig) *hashHelperImpl {
	return &hashHelperImpl{
		hashCost: config.HashCost,
	}
}

func (h *hashHelperImpl) Hash(str string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(str), h.hashCost)
	if err != nil {
		return "", fmt.Errorf("[helper][hash][Hash][bcrypt.GenerateFromPassword] Error: %w", err)
	}

	return string(hash), nil
}

func (h *hashHelperImpl) Check(str string, hash []byte) (bool, error) {
	err := bcrypt.CompareHashAndPassword(hash, []byte(str))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil
		}
		return false, fmt.Errorf("[helper][hash][Check][bcrypt.CompareHashAndPassword] Error: %w", err)
	}

	return true, nil
}

func HashSHA512(str string) string {
	hash := sha512.Sum512([]byte(str))

	return hex.EncodeToString(hash[:])
}
