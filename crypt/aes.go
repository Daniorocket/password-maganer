package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"log"
)

func AesEncrypt(key []byte, message string, iv []byte) (encmess string, err error) {
	plainText := []byte(message)
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Println(err)
		return "", err
	}
	cipherText := make([]byte, aes.BlockSize+len(plainText))
	IV := cipherText[:aes.BlockSize]
	for i, _ := range iv {
		IV[i] = iv[i]
	}
	stream := cipher.NewCFBEncrypter(block, IV)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	encmess = base64.URLEncoding.EncodeToString(cipherText)
	return encmess, nil
}
func AesDecrypt(key []byte, securemess string, iv []byte) (decodedmess string, err error) {
	cipherText, err := base64.URLEncoding.DecodeString(securemess)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(cipherText) < aes.BlockSize {
		err = errors.New("Ciphertext block size is too short!")
		return "", err
	}

	cipherText = cipherText[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)
	decodedmess = string(cipherText)
	return
}
