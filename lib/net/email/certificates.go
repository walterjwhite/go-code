package email

import (
	"fmt"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"io/ioutil"
)

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
