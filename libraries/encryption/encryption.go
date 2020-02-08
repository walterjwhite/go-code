package encryption

type Encryption interface {
	GetDecryptionKey() []byte
	GetEncryptionKey() []byte
}

type Encrypter interface {
	//EncryptFile(filename string, data []byte)
	Encrypt(plaintext []byte) (encrypted, salt []byte)
	EncryptFile(inFilename, outFilename, saltFile string)
}

type Decrypter interface {
	Decrypt(encrypted, salt []byte) []byte
	DecryptFile(inFilename, outFilename, saltFile string)
}
