package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"

	"golang.org/x/crypto/bcrypt"
)

const cost = 14

func HashPassword(password string) (string, error) {
	if len(password) == 0 {
		return "", nil
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(bytes), err
}

func CheckHashPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}


func padKey(key []byte, blockSize int) []byte {
	paddedKey := make([]byte, blockSize)
	copy(paddedKey, key)
	return paddedKey
}

func Encrypt(key []byte, password string) (string, error) {
	plaintext := []byte(password)

	block, err := aes.NewCipher(padKey(key, aes.BlockSize))
	if err != nil {
		return "", err
	}

	// Генерация случайного вектора инициализации (IV)
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	// Шифрование данных
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// Возвращаем зашифрованные данные в формате base64
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(key []byte, ciphertext string) (string, error) {
	ciphertextBytes, err := base64.URLEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(padKey(key, aes.BlockSize))
	if err != nil {
		return "", err
	}

	// Получение IV из зашифрованных данных
	iv := ciphertextBytes[:aes.BlockSize]
	ciphertextBytes = ciphertextBytes[aes.BlockSize:]

	// Расшифровка данных
	stream := cipher.NewCFBDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertextBytes))
	stream.XORKeyStream(plaintext, ciphertextBytes)

	return string(plaintext), nil
}
