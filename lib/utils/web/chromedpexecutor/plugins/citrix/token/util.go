package token

func Get(input string) []int {
	token := make([]int, 6)

	for i, char := range input {
		value := int(char - '0')

		token[i] = value
	}

	return token
}
