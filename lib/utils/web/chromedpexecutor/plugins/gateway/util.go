package gateway

func Get(input string) [] /*6*/ int {
	token := make([]int, 6)

	for i, char := range input {
		value := int(char - '0')

		token[i] = value
	}

	return token
}
