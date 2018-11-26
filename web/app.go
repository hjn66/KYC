package web

import (
	"fmt"

	"KYC/kyc-webservice/kyc"

	"KYC/kyc-webservice/webservice"

	"KYC/web/controllers"

	"github.com/go-martini/martini"
)

func Serve(app *controllers.Application) {

	martiniClassic := martini.Classic()
	kyc := kyc.New(app)
	webservice.RegisterWebService(kyc, martiniClassic)
	fmt.Println("Listening (http://localhost:8003/) ...")
	martiniClassic.RunOnAddr(":8003")

}
