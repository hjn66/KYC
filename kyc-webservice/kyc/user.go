package kyc

import (
	"KYC/web/controllers"
	"crypto/rsa"
	"fmt"
	"sync"
	"time"
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
	nonce      string
	GUID       int
}
type BlockChainRecord struct {
	NationalID string
	FirstName  string
	LastName   string
	BirthDate  string
	Photo      string
	PublicKey  string
}
type BlockChainDate struct {
	GUID   string
	Record BlockChainRecord
}
type QRTicket struct {
	Expiration string
	Nonce      string
}

type Login struct {
	Nonce          string
	LoginDate      time.Time
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
type Register struct {
	Nonce        string
	GUID         int
	RegisterDate time.Time
	Status       string
	User         User
}

type RegisterTable struct {
	RegisterList []*Register
	mutex        *sync.Mutex
}

type Conf struct {
	loginTable    *loginTable
	RegisterTable *RegisterTable
	app           *controllers.Application
	privateKey    *rsa.PrivateKey
}

func New(application *controllers.Application) *Conf {
	return &Conf{
		&loginTable{
			make([]*Login, 0),
			new(sync.Mutex),
		}, &RegisterTable{
			make([]*Register, 0),
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
func (logins *loginTable) GetLogin(nonce string) (*Login, error) {

	for _, login := range logins.loginList {
		if login.Nonce == nonce {
			return login, nil
		}
	}

	return nil, fmt.Errorf("Not Found")
}

// addLogin adds a new Login with the provided data.
func (registers *RegisterTable) addRegister(newRegister Register) {
	// Acquire our lock and make sure it will be released.
	registers.mutex.Lock()
	defer registers.mutex.Unlock()
	newRegister.Status = "Pending"
	newRegister.GUID = -1
	// Add entry to the RegisterTable
	registers.RegisterList = append(registers.RegisterList, &newRegister)
	return
}

// GetAllEntries returns all non-nil entries in the RegisterTable.
func (registers *RegisterTable) GetAllEntries() []*Register {
	// Placeholder for the entries we will be returning.
	entries := make([]*Register, 0)

	// Iterate through all existig entries.
	for _, entry := range registers.RegisterList {
		if entry != nil {
			// Entry is not nil, so we want to return it.
			entries = append(entries, entry)
		}
	}

	return entries
}

// GetLogin returns entry in the RegisterTable.
func (registers *RegisterTable) GetRegister(nonce string) (*Register, error) {

	for _, register := range registers.RegisterList {
		if register.Nonce == nonce {
			return register, nil
		}
	}

	return nil, fmt.Errorf("Not Found")
}

// changeRegisterStatus change status of register
func (registers *RegisterTable) changeRegisterStatus(nonce string, status string) error {
	registers.mutex.Lock()
	defer registers.mutex.Unlock()
	for index, register := range registers.RegisterList {
		if register.Nonce == nonce {
			registers.RegisterList[index].Status = status
			return nil
		}
	}

	return fmt.Errorf("Not Found")
}

// GetLogin returns entry in the RegisterTable.
func (registers *RegisterTable) setRegisterGUID(nonce string, guid int) error {
	registers.mutex.Lock()
	defer registers.mutex.Unlock()
	for index, register := range registers.RegisterList {
		if register.Nonce == nonce {
			registers.RegisterList[index].GUID = guid
			return nil
		}
	}

	return fmt.Errorf("Not Found")
}
