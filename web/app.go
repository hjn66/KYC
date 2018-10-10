package web

import (
	"fmt"

	"github.com/chainHero/heroes-service/kyc-webservice/kyc"
	"github.com/chainHero/heroes-service/kyc-webservice/webservice"
	"github.com/chainHero/heroes-service/web/controllers"
	"github.com/go-martini/martini"
)

func Serve(app *controllers.Application) {

	martiniClassic := martini.Classic()
	// loginTable := kyc.NewLoginTable()
	kyc := kyc.New(app)
	webservice.RegisterWebService(kyc, martiniClassic)
	fmt.Println("Listening (http://localhost:8003/) ...")
	martiniClassic.RunOnAddr(":8003")

}
