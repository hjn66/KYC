package controllers

import (
	"encoding/base64"
	"encoding/json"
	"html/template"
	"net/http"

	qrcode "github.com/skip2/go-qrcode"
)

func (app *Application) HomeHandler(w http.ResponseWriter, r *http.Request) {
	userAsBytes, err := app.Fabric.Query("igkm1a0CBwXX5RV22BSxAoL0iTv9ewOQJpDsSPPnj+U=")
	if err != nil {
		http.Error(w, "Unable to query the blockchain", 500)
	}

	user := User{}
	json.Unmarshal(userAsBytes, &user)

	var png []byte
	png, err = qrcode.Encode("LastName: "+user.LastName+" FirstName: "+user.FirstName, qrcode.Medium, 256)
	encodedString := "<img src=\"data:image/png;base64," + base64.StdEncoding.EncodeToString(png) + "\"/>"
	data := &struct {
		FirstName string
		LastName  string
		BirthDate string
		PublicKey string
		QRpng     template.HTML
	}{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		BirthDate: user.BirthDate,
		PublicKey: user.PublicKey,
		QRpng:     template.HTML(encodedString),
	}

	renderTemplate(w, r, "home.html", data)
}
