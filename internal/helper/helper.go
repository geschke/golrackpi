// Copyright 2022 Ralf Geschke. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package helper

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	b64 "encoding/base64"
	mrand "math/rand"
	"time"

	"golang.org/x/crypto/pbkdf2"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func init() {
	mrand.Seed(time.Now().UnixNano())
}

// RandSeq returns a random sequence of letters as string value
func RandSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[mrand.Intn(len(letters))]
	}
	return string(b)
}

// GetPBKDF2Hash returns a key from the password, salt and iteration count, returning a []byte of length 32
func GetPBKDF2Hash(password string, salt string, rounds int) []byte {
	// key length 32 for SHA256
	tmp := pbkdf2.Key([]byte(password), []byte(salt), rounds, 32, sha256.New)
	return tmp
}

// GetHMACSHA256 returns a HMAC hash using the SHA256 algorithm as []byte
// Credits to https://golangcode.com/generate-sha256-hmac/
func GetHMACSHA256(secret []byte, valueToEncrypt string) []byte {

	// Create a new HMAC by defining the hash type and the key (as byte array)
	h := hmac.New(sha256.New, []byte(secret))

	// Write Data to it
	h.Write([]byte(valueToEncrypt))
	mac := h.Sum(nil)
	return mac
}

// GetSHA256Hash returns a SHA256 checksum of the data as []byte
func GetSHA256Hash(valueToHash []byte) []byte {
	hash := sha256.Sum256(valueToHash)
	return hash[:]
}

// CreateClientProof returns the client proof computed by client and server signature as base64 encoded string
func CreateClientProof(clientSignature []byte, serverSignature []byte) string {
	clientSignatureLength := len(clientSignature)
	//var result [clientSignatureLength]byte
	result := make([]byte, clientSignatureLength)
	//fmt.Println("clientSignature:", clientSignature)
	//fmt.Println("serverSignature:", serverSignature)
	//fmt.Println("Length clientSignature", len(clientSignature))
	//result = new byte[clientSignature.length];
	for i := 0; i < len(clientSignature); i++ {
		result[i] = (byte(0xff & (clientSignature[i] ^ serverSignature[i])))
	}
	//fmt.Println("result:", result)

	return b64.StdEncoding.EncodeToString(result)

}

// GenerateRandomBytes returns random []byte with length n
func GenerateRandomBytes(n int) ([]byte, error) {

	b := make([]byte, n)
	_, err := rand.Read(b)

	if err != nil {
		return nil, err
	}

	return b, nil
}
