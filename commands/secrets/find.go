package main

import (
	"fmt"
	"github.com/walterjwhite/go-application/libraries/security/secrets"
	"strings"
)

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
