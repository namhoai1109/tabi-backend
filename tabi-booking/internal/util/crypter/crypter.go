package crypter

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"tabi-booking/internal/model"
	"time"

	"github.com/uniplaces/carbon"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes the password using bcrypt
func HashPassword(password string) (string, error) {
	hashedPW, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPW), nil
}

// Encrypt encrypt a string
func Encrypt(key []byte, text string) (*string, error) {
	plaintext := []byte(text)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the cipherText.
	cipherText := make([]byte, aes.BlockSize+len(plaintext))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plaintext)

	// convert to base64
	result := base64.URLEncoding.EncodeToString(cipherText)
	return &result, nil
}

// Decrypt encrypt a string
func Decrypt(key []byte, cryptoText string) (*string, error) {
	cipherText, _ := base64.URLEncoding.DecodeString(cryptoText)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the cipherText.
	if len(cipherText) < aes.BlockSize {
		return nil, errors.New("cipherText too short")
	}
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(cipherText, cipherText)

	result := string(cipherText)

	return &result, nil
}

// CompareHashAndPassword matches hash with password. Returns true if hash and password match.
func CompareHashAndPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// GenRefreshToken generate refresh token
func GenRefreshToken(secret string) string {
	now, _ := carbon.NowInLocation("Asia/Ho_Chi_Minh")
	endOfDay := now.AddYears(3)
	endOfDay.SetHour(0)
	endOfDay.SetMinute(0)
	endOfDay.SetSecond(0)

	data := model.RefreshToken{
		ExpiredAt: endOfDay.Time.Unix(),
	}
	dataStr, _ := json.Marshal(data)
	encryptedData, _ := Encrypt([]byte(secret), string(dataStr))

	return *encryptedData
}

// ValidateRefreshToken validate refresh token
func ValidateRefreshToken(token, secret string) bool {
	decrypted, err := Decrypt([]byte(secret), token)
	if err != nil {
		fmt.Println("===== err Decrypt", err.Error())
		return false
	}

	result := new(model.RefreshToken)
	if err := json.Unmarshal([]byte(*decrypted), result); err != nil {
		fmt.Println("===== err Unmarshal", err.Error())
		return false
	}

	t := time.Unix(result.ExpiredAt, 0)
	return t.After(time.Now())
}
