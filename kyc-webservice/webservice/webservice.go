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

package webservice

import (
	"net/http"

	"github.com/go-martini/martini"
)

// WebService is the interface that should be implemented by types that want to
// provide web services.
type WebService interface {
	// WebPost wraps the POST method. Again an empty params means that the
	// request should be applied to the collection. A non-empty param will
	// have an "id" key that refers to the entry that should be processed
	// (note this specific case is usually not supported unless each entry
	// is also a collection).
	RegisterUserPost(params martini.Params, req *http.Request) (int, string)

	GetTicketPost(params martini.Params, req *http.Request) (int, string)
	GetTicketQR(params martini.Params) (int, string)
	LoginPost(params martini.Params, req *http.Request) (int, string)
	CheckFieldPost(params martini.Params, req *http.Request) (int, string)

	// WebGet is Just as above, but for the GET method. It returns all the
	// users in the ledger.
	WebGet(params martini.Params) (int, string)
	WebGetLogins(params martini.Params, req *http.Request) (int, string)
}

// RegisterWebService adds martini routes to the relevant webservice methods
// based on the path returned by GetPath. Each method is registered once for
// the collection and once for each id in the collection.
func RegisterWebService(webService WebService,
	classicMartini *martini.ClassicMartini) {

	classicMartini.Get("/users", webService.WebGet)
	classicMartini.Get("/logins", webService.WebGetLogins)
	classicMartini.Post("/register", webService.RegisterUserPost)
	classicMartini.Post("/getticket", webService.GetTicketPost)
	classicMartini.Get("/getticketQR", webService.GetTicketQR)
	classicMartini.Post("/login", webService.LoginPost)
	classicMartini.Post("/checkfield", webService.CheckFieldPost)

}
