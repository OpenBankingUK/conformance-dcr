package certs

import (
	"crypto/rsa"
	"os"

	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
)

func ParseRsaPrivateKeyFromPemFile(privFile string) (*rsa.PrivateKey, error) {
	fileContents, err := os.ReadFile(privFile)
	if err != nil {
		return nil, errors.Wrap(err, "parsing rsa private key from file")
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(fileContents)
	if err != nil {
		return nil, errors.Wrap(err, "parsing rsa private key from file")
	}

	err = privateKey.Validate()
	if err != nil {
		return nil, errors.Wrap(err, "parsing rsa private key from file")
	}

	return privateKey, nil
}
