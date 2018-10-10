// Copyright 2013 Bruno Albuquerque (bga@bug-br.org.br).
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not
// use this file except in compliance with the License. You may obtain a copy of
// the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations under
// the License.

// WebService related methods.

package kyc

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	random "math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-martini/martini"
	qrcode "github.com/skip2/go-qrcode"
)

// WebGet implements webservice.WebGet.
func (conf *Conf) WebGet(params martini.Params) (int, string) {
	if len(params) == 0 {
		// No params. Return entire collection encoded as JSON.
		encodedEntries, err := conf.app.Fabric.QueryRange("1000", "2000")

		if err != nil {
			// Failed encoding collection.
			return http.StatusInternalServerError, "internal error"
		}

		// Return encoded entries.
		return http.StatusOK, string(encodedEntries)
	}
	return 0, ""
}

// WebGetLogins implements webservice.WebGetLogins.
func (conf *Conf) WebGetLogins(params martini.Params, req *http.Request) (int, string) {
	if len(req.URL.Query()) == 0 {
		// No Query. Return entire collection encoded as JSON.
		encodedEntries, err := json.Marshal(conf.loginTable.GetAllEntries())
		if err != nil {
			// Failed encoding collection.
			return http.StatusInternalServerError, "internal error"
		}

		// Return encoded entries.
		return http.StatusOK, string(encodedEntries)
	} else {
		nounce := req.URL.Query().Get("nounce")
		if nounce != "" {
			resLogin, err := conf.loginTable.GetLogin(nounce)
			if err != nil {
				// Nonce not Found
				return http.StatusOK, "{}"
			}
			encodedLogin, err := json.Marshal(resLogin)

			if err != nil {
				// Failed encoding resLogin
				return http.StatusInternalServerError, "internal error"
			}
			return http.StatusOK, string(encodedLogin)
		}
	}
	return http.StatusOK, "{}"
}

// WebGet implements webservice.WebGet.
func (conf *Conf) GetTicketQR(params martini.Params) (int, string) {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	nonce := make([]rune, 10)
	for i := range nonce {
		nonce[i] = letterRunes[random.Intn(len(letterRunes))]
	}

	expiration := time.Now().Add(100 * time.Minute)
	var qrticket QRTicket
	qrticket.Expiration = expiration.Format(time.RFC3339)
	qrticket.Nounce = string(nonce)

	encodedTicket, _ := json.Marshal(qrticket)
	publicKey := conf.privateKey.PublicKey
	sha256 := sha256.New()
	encrypted, err := rsa.EncryptOAEP(sha256, rand.Reader, &publicKey, encodedTicket, nil)
	if err != nil {
		fmt.Printf("EncryptOAEP: %s\n", err)
	}
	data := struct {
		T string
		N string
		O string
		F string
	}{
		T: base64.StdEncoding.EncodeToString(encrypted),
		N: string(nonce),
		O: "Melli Bank",
		F: "FLI",
	}
	encodedData, err := json.Marshal(data)

	var png []byte
	png, err = qrcode.Encode(string(encodedData), qrcode.High, 512)
	template, err := ioutil.ReadFile("getticket.html")
	// png, err = ioutil.ReadFile("/home/hojjat/Desktop/111.png")
	html := string(template)
	html = strings.Replace(html, "{QRIMAGE}", base64.StdEncoding.EncodeToString(png), 1)
	html = strings.Replace(html, "{TICKET}", data.T, 1)
	html = strings.Replace(html, "{NOUNCE}", data.N, -1)

	return http.StatusOK, html
}

func (conf *Conf) GetTicketPost(params martini.Params,
	req *http.Request) (int, string) {

	// Make sure Body is closed when we are done.
	defer req.Body.Close()
	// req.
	// Read request body.
	requestBody, err := ioutil.ReadAll(req.Body)
	fmt.Println("requestBody: " + string(requestBody))
	if err != nil {
		return http.StatusInternalServerError, "internal error"
	}

	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	nonce := make([]rune, 10)
	for i := range nonce {
		nonce[i] = letterRunes[random.Intn(len(letterRunes))]
	}

	expiration := time.Now().Add(10 * time.Minute)
	guid := struct {
		GUID int
	}{
		GUID: -1,
	}
	err = json.Unmarshal(requestBody, &guid)
	if err != nil || guid.GUID == -1 {
		return http.StatusInternalServerError, "Bad Format Parameter: Need GUID"
	}
	var ticket Ticket
	ticket.Expiration = expiration.Format(time.RFC3339)
	ticket.Nounce = string(nonce)
	ticket.GUID = guid.GUID

	encodedTicket, _ := json.Marshal(ticket)
	publicKey := conf.privateKey.PublicKey
	sha256 := sha256.New()
	encrypted, err := rsa.EncryptOAEP(sha256, rand.Reader, &publicKey, encodedTicket, nil)
	if err != nil {
		fmt.Printf("EncryptOAEP: %s\n", err)
	}
	data := struct {
		Ticket string
		Nounce string
	}{
		Ticket: base64.StdEncoding.EncodeToString(encrypted),
		Nounce: string(nonce),
	}
	encodedData, err := json.Marshal(data)
	return http.StatusOK, string(encodedData)
}

// RegisterUserPost implements webservice.RegisterUserPost.
func (conf *Conf) RegisterUserPost(params martini.Params,
	req *http.Request) (int, string) {

	// Make sure Body is closed when we are done.
	defer req.Body.Close()

	// Read request body.
	requestBody, err := ioutil.ReadAll(req.Body)
	// fmt.Println("requestBody: " + string(requestBody))

	if err != nil {
		return http.StatusInternalServerError, "Internal error"
	}

	if len(params) != 0 {
		// No keys in params. This is not supported.
		return http.StatusMethodNotAllowed, "Method not allowed"
	}

	lastuser, err := conf.app.Fabric.Query("LastUser")
	if err != nil {
		return http.StatusInternalServerError, "Unable to query LastUser GUID"
	}
	userkey := string(lastuser)

	// Unmarshal entry sent by the user.
	var user, invokeUser User
	err = json.Unmarshal(requestBody, &user)
	fmt.Println("------------------Register-------------------")
	fmt.Println("firstName:" + user.FirstName + " lastName:" + user.LastName + " Image:" + user.Photo + " NationalID:" + user.NationalID)

	h := sha256.New()
	h.Write([]byte(user.NationalID))
	invokeUser.NationalID = base64.StdEncoding.EncodeToString(h.Sum(nil))

	h.Reset()
	h.Write([]byte(user.BirthDate))
	invokeUser.BirthDate = base64.StdEncoding.EncodeToString(h.Sum(nil))

	h.Reset()
	h.Write([]byte(user.FirstName))
	invokeUser.FirstName = base64.StdEncoding.EncodeToString(h.Sum(nil))

	h.Reset()
	h.Write([]byte(user.LastName))
	invokeUser.LastName = base64.StdEncoding.EncodeToString(h.Sum(nil))

	h.Reset()
	h.Write([]byte(user.Photo))
	invokeUser.Photo = base64.StdEncoding.EncodeToString(h.Sum(nil))

	invokeUser.PublicKey = user.PublicKey

	userAsBytes, _ := json.Marshal(invokeUser)

	txid, err := conf.app.Fabric.RegisterUser(userkey, userAsBytes)

	if err != nil {
		return http.StatusInternalServerError, "Unable to register user in the blockchain"
	}
	userID, _ := strconv.Atoi(userkey)
	data := &struct {
		GUID int
		TXID string
	}{
		GUID: userID,
		TXID: txid,
	}
	encodedData, err := json.Marshal(data)
	return http.StatusOK, string(encodedData)
}

func (conf *Conf) LoginPost(params martini.Params,
	req *http.Request) (int, string) {
	// Make sure Body is closed when we are done.
	defer req.Body.Close()

	// Read request body.
	requestBody, err := ioutil.ReadAll(req.Body)
	fmt.Println("requestBody: " + string(requestBody))

	if err != nil {
		return http.StatusInternalServerError, "Internal error"
	}

	if len(params) != 0 {
		// No keys in params. This is not supported.
		return http.StatusMethodNotAllowed, "Method not allowed"
	}

	loginResponse := &struct {
		LoginSuccessful bool
		Message         string
	}{
		LoginSuccessful: false,
		Message:         "",
	}

	loginData := struct {
		Ticket       string
		NationalId   string
		SignedNounce string
	}{
		Ticket:       "",
		NationalId:   "",
		SignedNounce: "",
	}
	err = json.Unmarshal(requestBody, &loginData)
	if err != nil {
		loginResponse.Message = "Bad Request Format, Need Ticket token by getticket, NationalID and encrypted Nounce"
		encodedLoginResponse, _ := json.Marshal(loginResponse)
		return http.StatusOK, string(encodedLoginResponse)
	}

	var ticket Ticket
	sha256 := sha256.New()
	decodedTikcet, err := base64.StdEncoding.DecodeString(loginData.Ticket)
	if err != nil {
		loginResponse.Message = "Invalid Ticket - Base64 decode Error"
		encodedLoginResponse, _ := json.Marshal(loginResponse)
		return http.StatusOK, string(encodedLoginResponse)
	}
	decryptedTicket, err := rsa.DecryptOAEP(sha256, rand.Reader, conf.privateKey, decodedTikcet, nil)
	if err != nil {
		loginResponse.Message = "Invalid Ticket - Can not decrypte"
		encodedLoginResponse, _ := json.Marshal(loginResponse)
		return http.StatusOK, string(encodedLoginResponse)
	}
	err = json.Unmarshal(decryptedTicket, &ticket)
	if err != nil {
		loginResponse.Message = "Invalid Ticket - Bad format ticket"
		encodedLoginResponse, _ := json.Marshal(loginResponse)
		return http.StatusOK, string(encodedLoginResponse)
	}
	ticketTime, err := time.Parse(time.RFC3339, ticket.Expiration)
	if err != nil {
		loginResponse.Message = "Invalid Ticket - Time Parse Problem"
		encodedLoginResponse, _ := json.Marshal(loginResponse)
		return http.StatusOK, string(encodedLoginResponse)
	}
	if time.Now().After(ticketTime) {
		loginResponse.Message = "Expired Ticket"
		encodedLoginResponse, _ := json.Marshal(loginResponse)
		return http.StatusOK, string(encodedLoginResponse)
	}
	retrivedUserByte, err := conf.app.Fabric.Query(strconv.Itoa(ticket.GUID))
	fmt.Println("----------GUID in Login--------------------:" + strconv.Itoa(ticket.GUID))
	if err != nil {
		loginResponse.Message = "BlockChain Error - Can not query blockchain"
		encodedLoginResponse, _ := json.Marshal(loginResponse)
		return http.StatusOK, string(encodedLoginResponse)
	}
	if retrivedUserByte == nil {
		loginResponse.Message = "User Not Found"
		encodedLoginResponse, _ := json.Marshal(loginResponse)
		return http.StatusOK, string(encodedLoginResponse)
	}
	var retrivedUser User
	err = json.Unmarshal(retrivedUserByte, &retrivedUser)

	sha256.Reset()
	sha256.Write([]byte(loginData.NationalId))
	encodedNationalID := base64.StdEncoding.EncodeToString(sha256.Sum(nil))

	if encodedNationalID != retrivedUser.NationalID {
		loginResponse.Message = "Login Faild - NationalID not mached"
		encodedLoginResponse, _ := json.Marshal(loginResponse)
		return http.StatusOK, string(encodedLoginResponse)
	}

	pembyte := []byte(retrivedUser.PublicKey)
	block, _ := pem.Decode(pembyte)
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	publicKey := pub.(*rsa.PublicKey)

	sha256.Reset()
	sha256.Write([]byte(ticket.Nounce))
	hashedNounce := sha256.Sum(nil)
	signedNounce, _ := base64.StdEncoding.DecodeString(loginData.SignedNounce)
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashedNounce, signedNounce)
	if err != nil {
		loginResponse.Message = "Login Failed - Signature verification Error!"
		encodedLoginResponse, _ := json.Marshal(loginResponse)
		return http.StatusOK, string(encodedLoginResponse)
	}

	loginResponse.LoginSuccessful = true
	loginResponse.Message = "Login Successful"
	encodedLoginResponse, _ := json.Marshal(loginResponse)
	return http.StatusOK, string(encodedLoginResponse)
}

func (conf *Conf) CheckFieldPost(params martini.Params,
	req *http.Request) (int, string) {
	// Make sure Body is closed when we are done.
	defer req.Body.Close()

	// Read request body.
	requestBody, err := ioutil.ReadAll(req.Body)
	fmt.Println("requestBody: " + string(requestBody))

	if err != nil {
		return http.StatusInternalServerError, "Internal error"
	}

	if len(params) != 0 {
		// No keys in params. This is not supported.
		return http.StatusMethodNotAllowed, "Method not allowed"
	}

	loginResponse := &struct {
		CheckFirstName bool
		CheckLastName  bool
		CheckImage     bool
		Message        string
	}{
		CheckFirstName: false,
		CheckLastName:  false,
		CheckImage:     false,
		Message:        "",
	}

	checkFieldData := struct {
		Ticket       string
		FirstName    string
		LastName     string
		Image        string
		GUID         int
		SignedNounce string
	}{
		Ticket:       "",
		FirstName:    "",
		LastName:     "",
		Image:        "",
		GUID:         -1,
		SignedNounce: "",
	}
	err = json.Unmarshal(requestBody, &checkFieldData)
	fmt.Println("----------------checkFieldData---------------")
	fmt.Println("firstName:" + checkFieldData.FirstName + " lastName:" + checkFieldData.LastName + " Image:" + checkFieldData.Image)
	fmt.Println(checkFieldData)
	if err != nil {
		loginResponse.Message = "Bad Request Format, Need Ticket token by getticketQR, FirstName, LastName, Image, GUID and signed Nounce"
		encodedLoginResponse, _ := json.Marshal(loginResponse)
		return http.StatusOK, string(encodedLoginResponse)
	}

	var qrticket QRTicket
	sha256 := sha256.New()
	decodedTikcet, err := base64.StdEncoding.DecodeString(checkFieldData.Ticket)
	if err != nil {
		loginResponse.Message = "Invalid Ticket - Base64 decode Error"
		encodedLoginResponse, _ := json.Marshal(loginResponse)
		return http.StatusOK, string(encodedLoginResponse)
	}
	decryptedTicket, err := rsa.DecryptOAEP(sha256, rand.Reader, conf.privateKey, decodedTikcet, nil)
	if err != nil {
		loginResponse.Message = "Invalid Ticket - Can not decrypte"
		encodedLoginResponse, _ := json.Marshal(loginResponse)
		return http.StatusOK, string(encodedLoginResponse)
	}
	err = json.Unmarshal(decryptedTicket, &qrticket)
	if err != nil {
		loginResponse.Message = "Invalid Ticket - Bad format ticket"
		encodedLoginResponse, _ := json.Marshal(loginResponse)
		return http.StatusOK, string(encodedLoginResponse)
	}
	ticketTime, err := time.Parse(time.RFC3339, qrticket.Expiration)
	if err != nil {
		loginResponse.Message = "Invalid Ticket - Time Parse Problem"
		encodedLoginResponse, _ := json.Marshal(loginResponse)
		return http.StatusOK, string(encodedLoginResponse)
	}
	if time.Now().After(ticketTime) {
		loginResponse.Message = "Expired Ticket"
		encodedLoginResponse, _ := json.Marshal(loginResponse)
		return http.StatusOK, string(encodedLoginResponse)
	}
	retrivedUserByte, err := conf.app.Fabric.Query(strconv.Itoa(checkFieldData.GUID))
	if err != nil {
		loginResponse.Message = "BlockChain Error - Can not query blockchain"
		encodedLoginResponse, _ := json.Marshal(loginResponse)
		return http.StatusOK, string(encodedLoginResponse)
	}
	if retrivedUserByte == nil {
		loginResponse.Message = "User Not Found"
		encodedLoginResponse, _ := json.Marshal(loginResponse)
		return http.StatusOK, string(encodedLoginResponse)
	}
	var retrivedUser User
	err = json.Unmarshal(retrivedUserByte, &retrivedUser)

	pembyte := []byte(retrivedUser.PublicKey)
	block, _ := pem.Decode(pembyte)
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	publicKey := pub.(*rsa.PublicKey)

	sha256.Reset()
	sha256.Write([]byte(qrticket.Nounce))
	hashedNounce := sha256.Sum(nil)
	signedNounce, _ := base64.StdEncoding.DecodeString(checkFieldData.SignedNounce)
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashedNounce, signedNounce)
	if err != nil {
		loginResponse.Message = "Login Failed - Signature verification Error!"
		encodedLoginResponse, _ := json.Marshal(loginResponse)
		return http.StatusOK, string(encodedLoginResponse)
	}

	sha256.Reset()
	sha256.Write([]byte(checkFieldData.FirstName))
	encodedFirstName := base64.StdEncoding.EncodeToString(sha256.Sum(nil))
	if encodedFirstName == retrivedUser.FirstName {
		loginResponse.Message = "FirstName mached"
		loginResponse.CheckFirstName = true
	} else {
		loginResponse.Message = "FirstName not mached"
	}

	sha256.Reset()
	sha256.Write([]byte(checkFieldData.LastName))
	encodedLastName := base64.StdEncoding.EncodeToString(sha256.Sum(nil))
	if encodedLastName == retrivedUser.LastName {
		loginResponse.Message += ", LastName mached"
		loginResponse.CheckLastName = true
	} else {
		loginResponse.Message += ", LastName not mached"
	}

	sha256.Reset()
	sha256.Write([]byte(checkFieldData.Image))
	encodedImage := base64.StdEncoding.EncodeToString(sha256.Sum(nil))
	if encodedImage == retrivedUser.Photo {
		loginResponse.Message += ", Image mached"
		loginResponse.CheckImage = true
	} else {
		loginResponse.Message += ", Image not mached"
	}
	var login Login
	login.CheckFirstName = loginResponse.CheckFirstName
	login.CheckLastName = loginResponse.CheckLastName
	login.CheckImage = loginResponse.CheckImage
	login.GUID = checkFieldData.GUID
	login.Nounce = qrticket.Nounce
	login.FirstName = checkFieldData.FirstName
	login.LastName = checkFieldData.LastName
	login.Image = checkFieldData.Image
	conf.loginTable.addLogin(login)
	encodedLoginResponse, _ := json.Marshal(loginResponse)
	return http.StatusOK, string(encodedLoginResponse)
}