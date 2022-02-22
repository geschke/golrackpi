package helper

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	b64 "encoding/base64"
	"fmt"
	"golang.org/x/crypto/pbkdf2"
	mrand "math/rand"

	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func init() {
	mrand.Seed(time.Now().UnixNano())
}

func RandSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[mrand.Intn(len(letters))]
	}
	return string(b)
}

func GetPBKDF2Hash(password string, salt string, rounds int) []byte {
	fmt.Println("in getPBKDF2Hash")
	// key length 32 for SHA256
	tmp := pbkdf2.Key([]byte(password), []byte(salt), rounds, 32, sha256.New)
	fmt.Println("hash ist:", tmp)
	return tmp
}

func GetHMACSHA256(secret []byte, valueToEncrypt string) []byte {
	//fmt.Printf("Secret: %s Data: %s\n", secret, valueToEncrypt)

	// Create a new HMAC by defining the hash type and the key (as byte array)
	h := hmac.New(sha256.New, []byte(secret))

	// Write Data to it
	h.Write([]byte(valueToEncrypt))
	mac := h.Sum(nil)
	fmt.Println("MAC:", mac)
	return mac
}

func GetSHA256Hash(valueToHash []byte) []byte {
	hash := sha256.Sum256(valueToHash)
	return hash[:]
}

func CreateClientProof(clientSignature []byte, serverSignature []byte) string {
	clientSignatureLength := len(clientSignature)
	//var result [clientSignatureLength]byte
	result := make([]byte, clientSignatureLength)
	fmt.Println("clientSignature:", clientSignature)
	fmt.Println("serverSignature:", serverSignature)
	fmt.Println("Length clientSignature", len(clientSignature))
	//result = new byte[clientSignature.length];
	for i := 0; i < len(clientSignature); i++ {
		result[i] = (byte(0xff & (clientSignature[i] ^ serverSignature[i])))
	}
	fmt.Println("result:", result)

	return b64.StdEncoding.EncodeToString(result)

}

func GenerateRandomBytes(n int) ([]byte, error) {

	b := make([]byte, n)
	_, err := rand.Read(b)

	if err != nil {
		return nil, err
	}

	return b, nil
}
