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
