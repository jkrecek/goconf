package goconf

import (
	"golang.org/x/crypto/openpgp"
	"io"
	"os"
)

type EncryptionStructure struct {
	PublicKeyPath string
}

func (en *EncryptionStructure) GetEncryptionEntity() (openpgp.EntityList, error) {
	publicKeyFile := getValidFullPath(en.PublicKeyPath, getExecutePath())
	keyringFileBuffer, err := os.Open(publicKeyFile)
	if err != nil {
		return openpgp.EntityList{}, err
	}

	defer keyringFileBuffer.Close()
	entityList, err := openpgp.ReadArmoredKeyRing(keyringFileBuffer)
	if err != nil {
		return openpgp.EntityList{}, err
	}

	return entityList, nil
}

func (en *EncryptionStructure) Encrypt(ciphertext io.Writer) (io.WriteCloser, error) {
	entity, err := en.GetEncryptionEntity()
	if err != nil {
		return nil, err
	}

	return openpgp.Encrypt(ciphertext, entity, nil, nil, nil)
}

func (en *EncryptionStructure) RuntimeTest() (err error, fatal bool) {
	_, err = en.GetEncryptionEntity()
	fatal = true
	return
}