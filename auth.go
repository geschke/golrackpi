// Copyright 2022 Ralf Geschke. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package golrackpi

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"errors"
	"io"

	"crypto/sha256"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/geschke/golrackpi/internal/helper"

	"net/http"
)

const (
	endpointAuthStart         string = "/api/v1/auth/start"
	endpointAuthFinish        string = "/api/v1/auth/finish"
	endpointAuthCreateSession string = "/api/v1/auth/create_session"
)

// AuthStartRequestType defines the JSON structure of the first step in the authentication process
type AuthStartRequestType struct {
	Nonce    string `json:"nonce"`
	Username string `json:"username"`
}

// AuthFinishRequestType defines the JSON structure of the second step in the authentication process
type AuthFinishRequestType struct {
	TransactionId string `json:"transactionId"`
	Proof         string `json:"proof"`
}

// AuthCreateSessionType defines the JSON structure of the last step in the authentication process
type AuthCreateSessionType struct {
	TransactionId string `json:"transactionId"`
	Iv            string `json:"iv"`
	Tag           string `json:"tag"`
	Payload       string `json:"payload"`
}

// AuthClient is the library's instance, it contains the configuration settings with SessionId after successful authentication
type AuthClient struct {
	Scheme    string
	Server    string
	Password  string
	SessionId string
}

// New returns a blank AuthClient instance with default http scheme
func New() *AuthClient {
	client := AuthClient{
		Scheme:    "http",
		Server:    "",
		Password:  "",
		SessionId: "",
	}
	return &client
}

// NewWithParameter returns an AuthClient instance.
// It takes an AuthClient structure as parameter, so it's possible to submit all connection settings in one step.
func NewWithParameter(param AuthClient) *AuthClient {
	if param.Scheme == "https" {
		param.Scheme = "https"
	} else {
		param.Scheme = "http"
	}
	client := AuthClient{
		Scheme:   param.Scheme,
		Server:   param.Server,
		Password: param.Password,
	}
	return &client
}

// SetServer sets the IP address or FQDN of the Kostal inverter
func (c *AuthClient) SetServer(server string) {
	c.Server = server
}

// SetServer sets the password for user access of the Kostal inverter
func (c *AuthClient) SetPassword(password string) {
	c.Password = password
}

// SetServer sets the scheme (http or https) of the Kostal inverter
func (c *AuthClient) SetScheme(scheme string) {
	if scheme == "https" {
		scheme = "https"
	} else {
		scheme = "http"
	}
	c.Scheme = scheme
}

// getUrl is a helper function which creates the API URL
func (c *AuthClient) getUrl(request string) string {
	return c.Scheme + "://" + c.Server + request
}

// Login handles the complete authenciation and login process.
// In case of success it returns the session id.
func (c *AuthClient) Login() (string, error) {

	// prepare step 1 of authentication
	randomString := helper.RandSeq(12)
	base64String := b64.StdEncoding.EncodeToString([]byte(randomString))

	userName := "user" // default user name of plant owner

	// create JSON authentication request
	startRequest := AuthStartRequestType{
		Nonce:    base64String,
		Username: userName,
	}

	body, _ := json.Marshal(startRequest)

	// send step 1 authentication request
	resp, err := http.Post(c.getUrl(endpointAuthStart), "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", errors.New("could not initiate authentication")
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", errors.New("request returned with http error " + resp.Status)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("could not read authentication response")
	}

	responseReader := bytes.NewReader(responseBody)

	var result map[string]interface{}
	json.NewDecoder(responseReader).Decode(&result)

	// expect a result with the following map entries, check existence
	serverNonceResp, serverNonceOk := result["nonce"]
	roundsResp, roundsOk := result["rounds"]
	serverSaltResp, serverSaltOk := result["salt"]
	transactionIdResp, transactionIdOk := result["transactionId"]

	if !serverNonceOk || !roundsOk || !serverSaltOk || !transactionIdOk {
		return "", errors.New("authentication response has wrong format")
	}
	serverNonce := serverNonceResp.(string)
	rounds := int64(roundsResp.(float64))
	serverSalt := serverSaltResp.(string)
	transactionId := transactionIdResp.(string)

	// do some magic crypto stuff
	var saltedPassword, clientKey, serverKey, storedKey, clientSignature, serverSignature []byte

	serverSaltDecoded, _ := b64.StdEncoding.DecodeString(serverSalt)

	saltedPassword = helper.GetPBKDF2Hash(c.Password, string(serverSaltDecoded), int(rounds))
	clientKey = helper.GetHMACSHA256(saltedPassword, "Client Key")
	serverKey = helper.GetHMACSHA256(saltedPassword, "Server Key")
	storedKey = helper.GetSHA256Hash(clientKey)

	authMessage := fmt.Sprintf("n=%s,r=%s,r=%s,s=%s,i=%d,c=biws,r=%s", userName, startRequest.Nonce, string(serverNonce), string(serverSalt), rounds, string(serverNonce))

	clientSignature = helper.GetHMACSHA256(storedKey, authMessage)
	serverSignature = helper.GetHMACSHA256(serverKey, authMessage)
	clientProof := helper.CreateClientProof(clientSignature, clientKey)

	// Perform step 2 of the authentication process

	finishRequest := AuthFinishRequestType{
		TransactionId: transactionId,
		Proof:         clientProof,
	}

	finishRequestBody, _ := json.Marshal(finishRequest)

	respFinish, err := http.Post(c.getUrl(endpointAuthFinish), "application/json", bytes.NewBuffer(finishRequestBody))

	if err != nil {
		return "", errors.New("could not initiate authentication finish request")
	}
	if respFinish.StatusCode != 200 {
		return "", errors.New("request returned with http error " + respFinish.Status)
	}

	defer respFinish.Body.Close()

	responseFinishBody, err := io.ReadAll(respFinish.Body)
	if err != nil {
		//Failed to read response.
		return "", errors.New("could not read from authentication finish request")
	}

	responseFinishReader := bytes.NewReader(responseFinishBody)

	var resultFinish map[string]interface{}
	json.NewDecoder(responseFinishReader).Decode(&resultFinish)

	// signature and token only set when login was successful
	_, authOkSignature := resultFinish["signature"]
	_, authOkToken := resultFinish["token"]
	if !authOkSignature || !authOkToken {
		return "", errors.New("authentication failed")
	}

	signatureStr := resultFinish["signature"].(string)
	signature, _ := b64.StdEncoding.DecodeString(signatureStr)
	token := resultFinish["token"].(string)

	cmpBytes := bytes.Compare(signature, serverSignature)

	if cmpBytes != 0 {
		return "", errors.New("signature check error")
	}

	h := hmac.New(sha256.New, []byte(storedKey))
	// Write Data to it
	h.Write([]byte("Session Key"))
	h.Write([]byte(authMessage))
	h.Write([]byte(clientKey))

	protocolKey := h.Sum(nil)

	ivNonce, _ := helper.GenerateRandomBytes(16)

	block, err := aes.NewCipher(protocolKey)
	if err != nil {
		return "", errors.New("cipher creation error " + err.Error())

	}

	// default tag size in Go is 16
	aesgcm, err := cipher.NewGCMWithNonceSize(block, 16)
	if err != nil {
		return "", errors.New("cipher error " + err.Error())
	}

	var tag []byte
	ciphertext := aesgcm.Seal(nil, ivNonce, []byte(token), nil)

	// golang appends tag at the end of ciphertext, so we have to extract it
	// see https://stackoverflow.com/questions/68350301/extract-tag-from-cipher-aes-256-gcm-golang
	ciphertext, tag = ciphertext[:len(ciphertext)-16], ciphertext[len(ciphertext)-16:]

	// perform step 3 of authentication request to get session id

	createSessionRequest := AuthCreateSessionType{
		TransactionId: transactionId,
		Iv:            b64.StdEncoding.EncodeToString(ivNonce),
		Tag:           b64.StdEncoding.EncodeToString(tag),
		Payload:       b64.StdEncoding.EncodeToString(ciphertext),
	}

	createSessionRequestBody, _ := json.Marshal(createSessionRequest)

	respCreateSession, err := http.Post(c.getUrl(endpointAuthCreateSession), "application/json", bytes.NewBuffer(createSessionRequestBody))

	if err != nil {
		return "", errors.New("could not create session")

	}
	if respCreateSession.StatusCode != 200 {
		return "", errors.New("request returned with http error " + respCreateSession.Status)
	}
	defer respCreateSession.Body.Close()

	responseCreateSessionBody, err := io.ReadAll(respCreateSession.Body)
	if err != nil {
		return "", errors.New("could not read from create session request")
	}

	responseCreateSessionReader := bytes.NewReader(responseCreateSessionBody)

	var resultCreateSession map[string]interface{}
	json.NewDecoder(responseCreateSessionReader).Decode(&resultCreateSession)

	sessionId, sessionOk := resultCreateSession["sessionId"]
	if !sessionOk {
		return "", errors.New("session id not available")
	}

	c.SessionId = sessionId.(string)
	return c.SessionId, nil

}

// Logout deletes the current session
func (c *AuthClient) Logout() (bool, error) {

	client := http.Client{}

	request, err := http.NewRequest("POST", c.getUrl("/api/v1/auth/logout"), nil)
	if err != nil {
		return false, err
	}

	request.Header.Add("authorization", "Session "+c.SessionId)

	response, err := client.Do(request)
	if err != nil || response.StatusCode != 200 {
		return false, errors.New("logout error")
	}
	defer response.Body.Close()

	c.SessionId = ""

	return true, nil

}

// Me returns information about the current user
func (c *AuthClient) Me() (map[string]interface{}, error) {
	result := make(map[string]interface{})
	client := http.Client{}

	request, err := http.NewRequest("GET", c.getUrl("/api/v1/auth/me"), nil)
	if err != nil {
		return result, err
	}

	request.Header.Add("authorization", "Session "+c.SessionId)

	response, err := client.Do(request)
	if err != nil {
		return result, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
	}
	var jsonResult interface{}
	err = json.Unmarshal(body, &jsonResult)
	if err != nil {
		return result, err
	}

	m, mOk := jsonResult.(map[string]interface{})

	if !mOk {
		return result, errors.New("could not read response")
	}
	return m, nil

}
