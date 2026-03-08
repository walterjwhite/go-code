package encryption

type Encryptor interface {
	Encrypt(data []byte) ([]byte, error)

	Decrypt(ciphertext []byte) ([]byte, error)
}
