package icims

import (
	"fmt"

	atsclient "github.com/walterjwhite/go-code/lib/utils/ats-client"
)

type Account = atsclient.Account

type IcimsATS struct{}

func (i *IcimsATS) GetName() string {
	return "icims"
}

func (i *IcimsATS) RegisterAccount(executor *atsclient.Executor, account *Account) error {
	fmt.Println("Registering account on iCIMS...")
	return nil
}

func (i *IcimsATS) LoginAccount(executor *atsclient.Executor, email, password string) error {
	fmt.Println("Logging into iCIMS...")
	return nil
}

func (i *IcimsATS) ApplyForJob(executor *atsclient.Executor, resumePath string, qaMap map[string]string, aiEnabled bool) error {
	fmt.Println("Applying for job on iCIMS...")
	return nil
}
