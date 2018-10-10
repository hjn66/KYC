package controllers

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func (app *Application) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	data := &struct {
		TransactionId string
		Success       bool
		Response      bool
	}{
		TransactionId: "",
		Success:       false,
		Response:      false,
	}
	if r.FormValue("submitted") == "true" {
		user := User{}

		h := sha256.New()
		h.Write([]byte(r.FormValue("nationalID")))
		userkey := base64.StdEncoding.EncodeToString(h.Sum(nil))

		h.Reset()
		h.Write([]byte(r.FormValue("bithdate")))
		user.BirthDate = base64.StdEncoding.EncodeToString(h.Sum(nil))

		h.Reset()
		h.Write([]byte(r.FormValue("firstname")))
		user.FirstName = base64.StdEncoding.EncodeToString(h.Sum(nil))

		h.Reset()
		h.Write([]byte(r.FormValue("lastName")))
		user.LastName = base64.StdEncoding.EncodeToString(h.Sum(nil))

		file, _, errfile := r.FormFile("publicKey")
		if errfile != nil {
			fmt.Println("publicKey")
			fmt.Println(errfile)
			return
		}
		pbyte := make([]byte, 1000)
		nbytes, _ := file.Read(pbyte)
		user.PublicKey = string(pbyte[:nbytes])

		userAsBytes, _ := json.Marshal(user)

		txid, err := app.Fabric.RegisterUser(userkey, userAsBytes)

		if err != nil {
			http.Error(w, "Unable to register user in the blockchain", 500)
		}
		data.TransactionId = txid
		data.Success = true
		data.Response = true

		expiration := time.Now().Add(10 * time.Minute)
		cookie := http.Cookie{Name: "nationalID", Value: r.FormValue("nationalID"), Expires: expiration}
		http.SetCookie(w, &cookie)
	}
	renderTemplate(w, r, "register.html", data)
}
