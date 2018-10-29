// Encryption functions
package kyc

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func savePEMKey(fileName string, key *rsa.PrivateKey) {
	outFile, err := os.Create(fileName)
	checkError(err)
	defer outFile.Close()

	var privateKey = &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	err = pem.Encode(outFile, privateKey)
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}

func GetPrivateKey() *rsa.PrivateKey {
	var privateKey *rsa.PrivateKey
	if _, err := os.Stat("private.pem"); err == nil {
		//Pair Key exist and must load
		privateKeyFile, err := os.Open("private.pem")
		checkError(err)
		pemfileinfo, _ := privateKeyFile.Stat()
		var size int64 = pemfileinfo.Size()
		pembytes := make([]byte, size)
		buffer := bufio.NewReader(privateKeyFile)
		_, err = buffer.Read(pembytes)
		data, _ := pem.Decode([]byte(pembytes))
		privateKeyFile.Close()

		privateKey, _ = x509.ParsePKCS1PrivateKey(data.Bytes)

	} else {
		//Pair Key not exist and must generate and save
		reader := rand.Reader
		bitSize := 2048

		privateKey, err = rsa.GenerateKey(reader, bitSize)
		checkError(err)

		savePEMKey("private.pem", privateKey)
	}
	return privateKey
}
