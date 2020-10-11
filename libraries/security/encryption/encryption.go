package encryption

type Encryption interface {
	GetDecryptionKey() []byte
	GetEncryptionKey() []byte
}

type Encrypter interface {
	Encrypt(plaintext []byte) []byte
}

type Decrypter interface {
	Decrypt(encrypted []byte) []byte
}
