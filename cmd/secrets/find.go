package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/walterjwhite/go-code/lib/security/secrets"
)

func find() {
	onFind(printOnMatch, flag.Args()[1:])
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
