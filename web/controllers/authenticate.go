package controllers

import (
	"encoding/base64"
	"html/template"
	"math/rand"
	"net/http"
	"time"

	qrcode "github.com/skip2/go-qrcode"
)

func (app *Application) AuthenticateHandler(w http.ResponseWriter, r *http.Request) {

	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	nonce := make([]rune, 10)
	for i := range nonce {
		nonce[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	expiration := time.Now().Add(10 * time.Minute)
	cookie := http.Cookie{Name: "nonce", Value: string(nonce), Expires: expiration}
	http.SetCookie(w, &cookie)
	var png []byte
	png, _ = qrcode.Encode("NationalCode,FirstName,LastName\n", qrcode.Medium, 256)
	encodedString := "<img src=\"data:image/png;base64," + base64.StdEncoding.EncodeToString(png) + "\"/>"

	data := &struct {
		QRpng template.HTML
	}{
		QRpng: template.HTML(encodedString),
	}

	renderTemplate(w, r, "authenticate.html", data)
}
