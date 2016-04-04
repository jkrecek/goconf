package goconf

import (
	"golang.org/x/crypto/openpgp"
	"io"
	"os"
)

type EncryptionStructure struct {
	PublicKeyPath string
}

func (en *EncryptionStructure) GetEncryptionEntity() openpgp.EntityList {
	publicKeyFile := getValidFullPath(en.PublicKeyPath, getExecutePath())
	keyringFileBuffer, _ := os.Open(publicKeyFile)
	defer keyringFileBuffer.Close()
	entityList, err := openpgp.ReadArmoredKeyRing(keyringFileBuffer)
	if err != nil {
		panic(err)
	}

	return entityList
}

func (en *EncryptionStructure) Encrypt(ciphertext io.Writer) (io.WriteCloser, error) {
	return openpgp.Encrypt(ciphertext, en.GetEncryptionEntity(), nil, nil, nil)
}
