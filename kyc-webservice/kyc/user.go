package kyc

import (
	"crypto/rsa"
	"fmt"
	"sync"

	"KYC/web/controllers"
)

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

// addLogin adds a new Login with the provided data.
func (logins *loginTable) addLogin(newLogin Login) {
	// Acquire our lock and make sure it will be released.
	logins.mutex.Lock()
	defer logins.mutex.Unlock()

	// Add entry to the LoginTable
	logins.loginList = append(logins.loginList, &newLogin)
	// Return the Id for the new entry.
	return
}

// GetAllEntries returns all non-nil entries in the loginTable.
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

// GetLogin returns entry in the loginTable.
func (logins *loginTable) GetLogin(nounce string) (*Login, error) {

	for _, login := range logins.loginList {
		if login.Nounce == nounce {
			return login, nil
		}
	}

	return nil, fmt.Errorf("Not Found")
}
