package main

import (
	"fmt"

	"strings"

	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/secrets"
)

// TODO: add support for flags
// instead of specifying the key type (email, user, pass), use a flag instead (-e, -u, -p)
func main() {
	_ = application.Configure()

	secrets.Find(secrets.NewFind(), printOnMatch)
}

func printOnMatch(filePath string) {
	key := removeProject(removeValue(filePath))
	fmt.Println(key)
}

func removeValue(input string) string {
	return strings.TrimSuffix(input, "/value")
}

func removeProject(input string) string {
	return strings.TrimPrefix(input, secrets.SecretsConfigurationInstance.RepositoryPath+"/")
}
