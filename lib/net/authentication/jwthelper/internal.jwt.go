package jwthelper

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/mitchellh/go-homedir"
	"time"
)

func Build(subjectKey string, expirationKey string, issuerKey string, principal string, issuer string, certificateFile string, certificateFilePassword string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		subjectKey:    principal,
		expirationKey: time.Now().Add(30 * time.Minute).Unix(),
		issuerKey:     issuer,
		"iat":         time.Now().Unix(), // Add issued at time
		"nbf":         time.Now().Unix(), // Not before time
	})

	certificateFilePath, err := homedir.Expand(certificateFile)
	if err != nil {
		return "", err
	}

	_, privateKey := Decode(certificateFilePath, certificateFilePassword)

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
