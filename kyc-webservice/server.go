// Create server for webservice
package main

import (
	"kyc-webservice/kyc"
	"kyc-webservice/webservice"

	"github.com/go-martini/martini"
)

func main() {
	martiniClassic := martini.Classic()
	kyc := kyc.NewUserList()
	webservice.RegisterWebService(kyc, martiniClassic)
	martiniClassic.RunOnAddr(":8003")
}
