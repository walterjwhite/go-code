package jwthelper

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/mitchellh/go-homedir"
	"log"
	"time"
)

func Build(subjectKey string, expirationKey string, issuerKey string, principal string, issuer string, certificateFile string, certificateFilePassword string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		subjectKey:    principal,
		expirationKey: time.Now().Add(30 * time.Minute).Unix(),
		issuerKey:     issuer,
	})

	certificateFilePath, err := homedir.Expand(certificateFile)
	if err != nil {
		log.Fatalf("Error getting private key: %v / %v", certificateFile, err)
	}

	_, privateKey := Decode(certificateFilePath, certificateFilePassword)
	tokenString, err := token.SignedString(privateKey)

	if err != nil {
		log.Fatalf("Error getting token: %v\n", err)
	}

	return tokenString
}
