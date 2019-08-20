package certs

import (
	"crypto/rsa"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"io/ioutil"
)

func ParseRsaPrivateKeyFromPemFile(privFile string) (*rsa.PrivateKey, error) {
	fileContents, err := ioutil.ReadFile(privFile)
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
