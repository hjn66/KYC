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
