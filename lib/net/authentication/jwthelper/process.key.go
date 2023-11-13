package jwthelper

import (
	"crypto/rsa"
	"crypto/x509"
	"errors"
	"golang.org/x/crypto/pkcs12"
	"log"
	"os"
)

func Decode(filename string, password string) (*x509.Certificate, *rsa.PrivateKey) {
	b, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error reading file: %v / %v\n", filename, err)
	}

	certificate, privateKey, err := decode(b, password)
	if err != nil {
		log.Fatalf("Error decoding: %v\n", err)
	}

	return certificate, privateKey
}

func decode(p12 []byte, password string) (*x509.Certificate, *rsa.PrivateKey, error) {
	privateKey, cert, err := pkcs12.Decode(p12, password)
	if err != nil {
		return nil, nil, err
	}

	if err := verify(cert); err != nil {
		return nil, nil, err
	}

	priv, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, nil, errors.New("Expected RSA Private Key Type")
	}

	return cert, priv, nil
}

func verify(cert *x509.Certificate) error {
	_, err := cert.Verify(x509.VerifyOptions{})
	if err != nil {
		return nil
	}

	switch e := err.(type) {
	case x509.CertificateInvalidError:
		switch e.Reason {
		case x509.Expired:
			return ErrExpired
		default:
			return err
		}
	case x509.UnknownAuthorityError:
		return nil
	default:
		return err
	}
}

var (
	ErrExpired = errors.New("Certificate has expired or is not yet valid.")
)
