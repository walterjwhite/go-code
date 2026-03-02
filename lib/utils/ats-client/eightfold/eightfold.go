package eightfold

import (
	"fmt"

	atsclient "github.com/walterjwhite/go-code/lib/utils/ats-client"
)

type Account = atsclient.Account

type EightfoldATS struct{}

func (e *EightfoldATS) GetName() string {
	return "eightfold"
}

func (e *EightfoldATS) RegisterAccount(executor *atsclient.Executor, account *Account) error {
	fmt.Println("Registering account on Eightfold...")
	return nil
}

func (e *EightfoldATS) LoginAccount(executor *atsclient.Executor, email, password string) error {
	fmt.Println("Logging into Eightfold...")
	return nil
}

func (e *EightfoldATS) ApplyForJob(executor *atsclient.Executor, resumePath string, qaMap map[string]string, aiEnabled bool) error {
	fmt.Println("Applying for job on Eightfold...")
	return nil
}
