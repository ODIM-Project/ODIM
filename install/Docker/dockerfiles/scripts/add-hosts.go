//(C) Copyright [2020] Hewlett Packard Enterprise Development LP
//
//Licensed under the Apache License, Version 2.0 (the "License"); you may
//not use this file except in compliance with the License. You may obtain
//a copy of the License at
//
//    http:#www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//License for the specific language governing permissions and limitations
// under the License.

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

const (
	// NULL is a constant for empty string
	NULL          = ""
	hostsFilePath = "/etc/hosts"
	contentHeader = "# --- User configured entries --- BEGIN"
	contentFooter = "# --- User configured entries --- END"
)

func main() {

	var inputFile string

	flag.StringVar(&inputFile, "file", "", "Path of the file which contains hosts info")
	flag.Parse()

	if inputFile == NULL {
		flag.PrintDefaults()
		os.Exit(1)
	}

	data, err := ioutil.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("Failed to read %s with error %v", inputFile, err)
	}

	if len(data) < 3 {
		log.Println("User configuration is empty, exiting")
		os.Exit(0)
	}

	fd, err := os.OpenFile(hostsFilePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open %s with error %v", hostsFilePath, err)
	}
	defer fd.Close()

	hostsData := fmt.Sprintf("\n%s\n%s\n%s\n", contentHeader, string(data), contentFooter)
	if _, err := fd.Write([]byte(hostsData)); err != nil {
		log.Fatalf("Failed to write to %s with error %v", hostsFilePath, err)
	}

	os.Exit(0)
}
