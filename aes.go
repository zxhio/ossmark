package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
)

var aesRandKey []byte

func init() {
	aesRandKey = make([]byte, 16)
	_, err := rand.Read(aesRandKey[:])
	if err != nil {
		panic(err)
	}
}

// PKCS7Padding adds padding to the input data according to the PKCS7 scheme.
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// PKCS7UnPadding removes the padding from the input data that was added by PKCS7Padding.
func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:length-unpadding]
}

// EncryptAES_CBC encrypts the plaintext using AES in CBC mode.
func EncryptAES_CBC(plaintext, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Padding the plaintext before encryption.
	plaintext = PKCS7Padding(plaintext, block.BlockSize())

	// Generate an initialization vector (IV) for CBC mode, same size as the block size.
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return "", err
	}

	// Initialize the CBC encryption mode with the key and IV.
	mode := cipher.NewCBCEncrypter(block, iv)
	ciphertext := make([]byte, len(plaintext))
	// Perform the encryption.
	mode.CryptBlocks(ciphertext, plaintext)

	// Combine the IV and ciphertext, usually needed for decryption.
	ivCiphertext := append(iv, ciphertext...)

	// Encode the result in Base64 for easier transmission or storage.
	return base64.StdEncoding.EncodeToString(ivCiphertext), nil
}

// DecryptAES_CBC decrypts the ciphertext previously encrypted with EncryptAES_CBC.
func DecryptAES_CBC(ciphertextBase64 string, key []byte) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextBase64)
	if err != nil {
		return "", err
	}
	// Separate the IV from the ciphertext.
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	// Decrypt the data.
	mode.CryptBlocks(plaintext, ciphertext)

	// Remove the padding and convert to string.
	plaintext = PKCS7UnPadding(plaintext)
	return string(plaintext), nil
}
