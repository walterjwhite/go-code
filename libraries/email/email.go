package email

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
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

func addCert(certificatePath string) {
	pem, err := ioutil.ReadFile(certificatePath)
	if err != nil {
		log.Fatal(err)
	}

	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		log.Fatalf("Failed to append root CA cert @ %v\n", certificatePath)
	}
}
