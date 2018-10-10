package controllers

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"math/rand"
	"net/http"
	"time"
)

func (app *Application) LoginHandler(w http.ResponseWriter, r *http.Request) {
	data := &struct {
		Message  string
		Nonce    string
		Success  bool
		Response bool
	}{
		Message:  "",
		Nonce:    "",
		Success:  false,
		Response: false,
	}
	if r.FormValue("submitted") == "true" {
		user := User{}

		h := sha256.New()
		h.Write([]byte(r.FormValue("nationalID")))
		userkey := base64.StdEncoding.EncodeToString(h.Sum(nil))

		userAsBytes, err := app.Fabric.Query(userkey)

		if err != nil || userAsBytes == nil {
			http.Error(w, "Unable to query the blockchain", 500)
		} else if userAsBytes == nil {
			data.Message = "User Not Registered"
			data.Success = false
		} else {
			json.Unmarshal(userAsBytes, &user)

			h.Reset()
			h.Write([]byte(r.FormValue("firstname")))
			providedFirstName := base64.StdEncoding.EncodeToString(h.Sum(nil))

			h.Reset()
			h.Write([]byte(r.FormValue("lastName")))
			providedLastName := base64.StdEncoding.EncodeToString(h.Sum(nil))

			if user.FirstName != providedFirstName || user.LastName != providedLastName {
				data.Message = "User Data missmach!"
				data.Success = false
			} else {
				data.Message = "User Found and match!"
				data.Success = true
			}
		}

	} else {
		letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
		nonce := make([]rune, 6)
		for i := range nonce {
			nonce[i] = letterRunes[rand.Intn(len(letterRunes))]
		}
		expiration := time.Now().Add(10 * time.Minute)
		cookie := http.Cookie{Name: "nonce", Value: string(nonce), Expires: expiration}
		http.SetCookie(w, &cookie)
		data.Nonce = string(nonce)
	}
	renderTemplate(w, r, "login.html", data)
}
