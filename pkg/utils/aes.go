package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

var RandAESKey []byte

func init() {
	RandAESKey = make([]byte, 16)
	_, err := rand.Read(RandAESKey[:])
	if err != nil {
		panic(err)
	}
}

// EncryptAES_CBC encrypts the plaintext using AES in CBC mode.
func EncryptAES_CBC(plainText, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Padding the plaintext before encryption.
	padding := aes.BlockSize - len(plainText)%aes.BlockSize
	padText := append(plainText, bytes.Repeat([]byte{byte(padding)}, padding)...)
	cipherText := make([]byte, len(padText))

	// Generate an initialization vector (IV) for CBC mode, same size as the block size.
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return "", err
	}

	// Initialize the CBC encryption mode with the key and IV.
	mode := cipher.NewCBCEncrypter(block, iv)
	// Perform the encryption.
	mode.CryptBlocks(cipherText, padText)

	// Combine the IV and ciphertext, usually needed for decryption.
	ivCipherText := append(iv, cipherText...)

	// Encode the result in Base64 for easier transmission or storage.
	return base64.StdEncoding.EncodeToString(ivCipherText), nil
}

func DecryptAES_CBC(cipherTextBase64 string, key []byte) (string, error) {
	cipherText, err := base64.StdEncoding.DecodeString(cipherTextBase64)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Separate the IV from the ciphertext.
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]
	plainText := make([]byte, len(cipherText))

	// Decrypt the data.
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plainText, cipherText)

	// Remove the padding and convert to string.
	padding := plainText[len(plainText)-1]
	if int(padding) > len(plainText) || int(padding) > aes.BlockSize {
		return "", fmt.Errorf("invalid padding")
	}
	plainText = plainText[:len(plainText)-int(padding)]
	return string(plainText), nil
}
