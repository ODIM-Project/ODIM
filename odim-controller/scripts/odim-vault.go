//(C) Copyright [2020] Hewlett Packard Enterprise Development LP
//
//Licensed under the Apache License, Version 2.0 (the "License"); you may
//not use this file except in compliance with the License. You may obtain
//a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//License for the specific language governing permissions and limitations
// under the License.

package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"flag"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/pbkdf2"
)

const (
	// NULL is a constant for empty string
	NULL = ""
	// Number of iterations for generating key using PBKDF2
	pbkdf2IterationVal = 100000
	// Length of the key to be generated by PBKDF2
	aes256CryptoKeyLen = 32
)

// readFile reads the file passed and returns the
// content as byte error. Exits if read fails.
func readFile(fpath string) []byte {
	data, err := ioutil.ReadFile(fpath)
	if err != nil {
		log.Fatalf("Failed to read from %s with error %v", fpath, err)
	}
	dataLen := len(data)
	if data[dataLen-1] == '\n' {
		return data[:dataLen-1]
	}
	return data
}

// writeFile writes the passed content to the file
// at given path. Exits if write fails.
func writeFile(fpath string, data []byte) {
	if err := ioutil.WriteFile(fpath, data, 0640); err != nil {
		log.Fatalf("Failed to write to %s with error %v", fpath, err)
	}
	return
}

// generateRandomKey generates a random key containing
// aplbhabets, numerals and special characters from the
// defined set and returns the key of requested length
func generateRandomKey(length int) (key string) {
	// characters set out of which a random string will be generated
	keyCharSet := strings.Split("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_-{}<>+[]$?@:;()%,", NULL)
	keyCharSetLen := len(keyCharSet)

	// set a seed value for the random function
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	key = NULL

	// generate key for the requested length
	for i := 0; i < length; i++ {
		key += keyCharSet[rnd.Intn(keyCharSetLen)]
	}
	return key
}

// encodeFileData encodes the data present in the
// passed file, and overwrites with the base64
// encoded content.
func encodeFileData(fpath string) {
	data := readFile(fpath)
	encData := base64.StdEncoding.EncodeToString(data)
	writeFile(fpath, []byte(encData))
	return
}

// decodeFileData decodes and returns the plain data
// present in the passed file.
func decodeFileData(fpath string) []byte {
	data := readFile(fpath)
	decData, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil || len(decData) <= 0 {
		log.Fatalf("Empty file or failed to decode data with error %v", err)
	}
	return decData
}

// getCryptoKey generates a key of requested length required
// for crypting data using the base64 encoded key present
// in keyFilePath and salt value passed.
func getCryptoKey(keyFilePath, salt string, keyLength int) []byte {
	key := decodeFileData(keyFilePath)
	return pbkdf2.Key(key, []byte(salt), pbkdf2IterationVal, keyLength, sha256.New)
}

// encryptFileData encrypts the content present in the passed
// file using the key present in keyFilePath and overwrites
// fpath file content with the encrypted data.
func encryptFileData(fpath, keyFilePath string) {
	// generate salt value required for generating
	// crypto key using PBKDF2 algo.
	salt := generateRandomKey(32)

	// obtain the crypto key of lenth 32-bytes,
	// using the key passed in keyFilePath.
	key := getCryptoKey(keyFilePath, salt, aes256CryptoKeyLen)

	// obtain AES cipher block object
	aesCipher, err := aes.NewCipher(key)
	if err != nil {
		log.Fatalf("Failed to obtain cipher object with error %v", err)
	}

	// obtain AES GCM cipher object
	aesGCM, err := cipher.NewGCM(aesCipher)
	if err != nil {
		log.Fatalf("Failed to obtain GCM wrapper object with error %v", err)
	}

	// generate nonce of length defined by AES GCM object
	nonce := []byte(generateRandomKey(aesGCM.NonceSize()))

	// encrypt and store data in the file passed
	plainData := readFile(fpath)
	cipherData := aesGCM.Seal(nil, nonce, plainData, nil)

	// overwrite file with the salt, nonce and encrypted data
	fd, err := os.OpenFile(fpath, os.O_WRONLY|os.O_TRUNC, 0644)
	defer fd.Close()

	// first 32-bytes will be salt data
	if _, err := fd.Write([]byte(salt)); err != nil {
		log.Fatalf("Failed to write data1 with error %v", err)
	}

	// next 16-bytes will be nonce data
	if _, err := fd.Write(nonce); err != nil {
		log.Fatalf("Failed to write data2 with error %v", err)
	}

	// write the encrypted data at the last
	if _, err := fd.Write(cipherData); err != nil {
		log.Fatalf("Failed to write data3 with error %v", err)
	}

	return
}

// decryptFileData decrypts and returns the encrypted data
// present in fpath file using the key in keyFilePath.
func decryptFileData(fpath, keyFilePath string) {
	// read the encrypted data present in fpath
	cipherData := readFile(fpath)

	if len(cipherData) < 32 {
		log.Fatalf("file content has been altered, not proceeding with decryption")
	}

	// read first 32-bytes of the read content,
	// which will be the salt value.
	salt := string(cipherData[:32])

	// obtain the crypto key of lenth 32-bytes,
	// using the key passed in keyFilePath.
	key := getCryptoKey(keyFilePath, salt, aes256CryptoKeyLen)

	// obtain AES cipher block object
	aesCipher, err := aes.NewCipher(key)
	if err != nil {
		log.Fatalf("Failed to obtain cipher object with error %v", err)
	}

	// obtain AES GCM object
	aesGCM, err := cipher.NewGCM(aesCipher)
	if err != nil {
		log.Fatalf("Failed to obtain GCM wrapper object with error %v", err)
	}

	// read next bytes of data of length aesGCM.NonceSize()
	// from 32nd byte to get the nonce value.
	nonceEndPos := 32 + aesGCM.NonceSize()
	nonce := cipherData[32:nonceEndPos]

	// decrypt the remaining content of the file
	plaintext, err := aesGCM.Open(nil, nonce, cipherData[nonceEndPos:], nil)
	if err != nil {
		log.Fatalf("Failed to decrypt data with error %v", err)
	}

	// print the decrypted data to console
	print(string(plaintext))
}

func main() {

	var encodeFile, encryptFile, decryptFile, keyFile string

	flag.StringVar(&encodeFile, "encode", "", "Path of the file which content to be encoded")
	flag.StringVar(&encryptFile, "encrypt", "", "Path of the file which content to be encrypted, -key is must")
	flag.StringVar(&decryptFile, "decrypt", "", "Path of the file which content to be decrypted, -key is must")
	flag.StringVar(&keyFile, "key", "", "Path of the key file, must be combined with -encrypt and -decrypt")
	flag.Parse()

	if encodeFile == NULL && encryptFile == NULL && decryptFile == NULL {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if encodeFile != NULL {
		encodeFileData(encodeFile)
	}

	if encryptFile != NULL {
		if keyFile == NULL {
			flag.PrintDefaults()
			os.Exit(1)
		}
		encryptFileData(encryptFile, keyFile)
	}

	if decryptFile != NULL {
		if keyFile == NULL {
			flag.PrintDefaults()
			os.Exit(1)
		}
		decryptFileData(decryptFile, keyFile)
	}

	os.Exit(0)
}
