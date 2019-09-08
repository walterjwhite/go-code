package email

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	"github.com/walterjwhite/go-application/libraries/logging"
)

var UserTlsConfig tls.Config
var rootCertPool = x509.NewCertPool()

type EmailSenderAccount struct {
	Username           string
	Password           string
	Domain             string
	EmailAddress       string
	Server             EmailServer
	Certificates       []string
	InsecureSkipVerify bool
}

type EmailServer struct {
	Host string
	Port int
}

type EmailMessage struct {
	To  []string
	Cc  []string
	Bcc []string

	Subject string
	Body    string
}

func (e *EmailSenderAccount) Initialize() {
	e.addCerts()
	UserTlsConfig = tls.Config{InsecureSkipVerify: e.InsecureSkipVerify, ServerName: e.Server.Host, RootCAs: rootCertPool}
}

func (e *EmailSenderAccount) addCerts() {
	for _, certificate := range e.Certificates {
		addCert(certificate)
	}
}

type UnableToAddCertificateError struct {
	CertificatePath string
}

func (e *UnableToAddCertificateError) Error() string {
	return fmt.Sprintf("Failed to append root CA cert @ %v\n", e.CertificatePath)
}

func addCert(certificatePath string) {
	pem, err := ioutil.ReadFile(certificatePath)
	logging.Panic(err)

	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		logging.Panic(&UnableToAddCertificateError{CertificatePath: certificatePath})
	}
}
