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

package kyc

import (
	"crypto/rsa"
	"fmt"
	"sync"

	"github.com/chainHero/heroes-service/web/controllers"
)

// GuestBookEntry represents a single entry in a Guest Book. It contains the
// usual fields.
type User struct {
	NationalID string
	FirstName  string
	LastName   string
	BirthDate  string
	Photo      string
	PublicKey  string
}

type Ticket struct {
	Expiration string
	Nounce     string
	GUID       int
}

type QRTicket struct {
	Expiration string
	Nounce     string
}

type Login struct {
	Nounce         string
	Expiration     string
	GUID           int
	FirstName      string
	LastName       string
	Image          string
	CheckFirstName bool
	CheckLastName  bool
	CheckImage     bool
}

type loginTable struct {
	loginList []*Login
	mutex     *sync.Mutex
}

type Conf struct {
	loginTable *loginTable
	app        *controllers.Application
	privateKey *rsa.PrivateKey
}

func NewLoginTable() *loginTable {
	return &loginTable{
		make([]*Login, 0),
		new(sync.Mutex),
	}
}

func New(application *controllers.Application) *Conf {
	return &Conf{
		&loginTable{
			make([]*Login, 0),
			new(sync.Mutex),
		},
		application,
		GetPrivateKey(),
	}
}

// AddEntry adds a new GuestBookEntry with the provided data.
func (logins *loginTable) addLogin(newLogin Login) {
	// Acquire our lock and make sure it will be released.
	logins.mutex.Lock()
	defer logins.mutex.Unlock()

	// Add entry to the LoginTable
	logins.loginList = append(logins.loginList, &newLogin)
	// Return the Id for the new entry.
	return
}

// GetAllEntries returns all non-nil entries in the Guest Book.
func (logins *loginTable) GetAllEntries() []*Login {
	// Placeholder for the entries we will be returning.
	entries := make([]*Login, 0)

	// Iterate through all existig entries.
	for _, entry := range logins.loginList {
		if entry != nil {
			// Entry is not nil, so we want to return it.
			entries = append(entries, entry)
		}
	}

	return entries
}

// GetAllEntries returns all non-nil entries in the Guest Book.
func (logins *loginTable) GetLogin(nounce string) (*Login, error) {

	for _, login := range logins.loginList {
		if login.Nounce == nounce {
			return login, nil
		}
	}

	return nil, fmt.Errorf("Not Found")
}
