package webservice

import (
	"net/http"

	"github.com/go-martini/martini"
)

// WebService is the interface that should be implemented by types that want to
// provide web services.
type WebService interface {
	RegisterUserPost(params martini.Params, req *http.Request) (int, string)

	GetTicketPost(params martini.Params, req *http.Request) (int, string)
	GetTicketQR(params martini.Params) (int, string)
	LoginPost(params martini.Params, req *http.Request) (int, string)
	CheckFieldPost(params martini.Params, req *http.Request) (int, string)

	HomeGet(params martini.Params) (int, string)
	WebGetUsers(params martini.Params) (int, string)
	WebGetLogins(params martini.Params, req *http.Request) (int, string)
}

// RegisterWebService adds martini routes to the relevant webservice methods

func RegisterWebService(webService WebService,
	classicMartini *martini.ClassicMartini) {

	classicMartini.Get("/", webService.HomeGet)
	classicMartini.Get("/users", webService.WebGetUsers)
	classicMartini.Get("/logins", webService.WebGetLogins)
	classicMartini.Post("/register", webService.RegisterUserPost)
	classicMartini.Post("/getticket", webService.GetTicketPost)
	classicMartini.Get("/getticketQR", webService.GetTicketQR)
	classicMartini.Post("/login", webService.LoginPost)
	classicMartini.Post("/checkfield", webService.CheckFieldPost)

}
