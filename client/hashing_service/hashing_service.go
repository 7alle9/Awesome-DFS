package hashing_service

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"golang.org/x/crypto/pbkdf2"
)

var hasher *HashingService

type HashingService struct {
	secret string
}

func addPadding(data []byte, blockSize int) []byte {
	padding := blockSize - (len(data) % blockSize)
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

func removePadding(data []byte) ([]byte, error) {
	paddingLen := int(data[len(data)-1])
	if paddingLen <= 0 || paddingLen > len(data) {
		return nil, errors.New("invalid padding")
	}

	for i := len(data) - paddingLen; i < len(data); i++ {
		if data[i] != byte(paddingLen) {
			return nil, errors.New("invalid padding")
		}
	}

	return data[:len(data)-paddingLen], nil
}

func (s *HashingService) deriveKey(salt []byte) []byte {
	return pbkdf2.Key([]byte(s.secret), salt, 100000, 32, sha256.New)
}

func (s *HashingService) EncryptByteArray(plaintext []byte) ([]byte, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("failed to generate salt: %v", err)
	}

	key := s.deriveKey(salt)

	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return nil, fmt.Errorf("failed to generate IV: %v", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher block: %v", err)
	}

	paddedPlaintext := addPadding(plaintext, aes.BlockSize)

	ciphertext := make([]byte, len(paddedPlaintext))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, paddedPlaintext)

	result := append(salt, iv...)
	result = append(result, ciphertext...)

	return result, nil
}

func (s *HashingService) DecryptByteArray(encryptedData []byte) ([]byte, error) {
	if len(encryptedData) < 48 { // Minimum length: 16(salt) + 16(iv) + 16(ciphertext block size)
		return nil, errors.New("encrypted data is too short")
	}

	salt := encryptedData[:16]
	iv := encryptedData[16:32]
	ciphertext := encryptedData[32:]

	key := s.deriveKey(salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher block: %v", err)
	}

	if len(iv) != block.BlockSize() {
		return nil, fmt.Errorf("invalid IV length: expected %d, got %d", block.BlockSize(), len(iv))
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	decrypted := make([]byte, len(ciphertext))
	mode.CryptBlocks(decrypted, ciphertext)

	unpadded, err := removePadding(decrypted)
	if err != nil {
		return nil, fmt.Errorf("failed to remove padding: %v", err)
	}

	return unpadded, nil
}

func InitHashingService(secret string) {
	hasher = &HashingService{secret: secret}
}

func GetHasher() (*HashingService, error) {
	if hasher == nil {
		return nil, errors.New("hashing service not initialized")
	}
	return hasher, nil
}
