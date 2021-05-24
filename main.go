package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
)

var (
	adwareConfig    = "configs/adware"
	malwareConfig   = "configs/malware"
	privacyConfig   = "configs/privacy"
	sexualConfig    = "configs/sexual"
	socialConfig    = "configs/social"
	exclusionConfig = "configs/exclusion"
)

func init() {
	//
}

func main() {
	validateAndSave("https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts", adwareConfig)
}

func validateAndSave(url, path string) {
	// Send a request to acquire all the information you need.
	response, err := http.Get(url)
	handleErrors(err)
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	handleErrors(err)
	// locate all domains
	regex := regexp.MustCompile(`(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z0-9][a-z0-9-]{0,61}[a-z0-9]`)
	domains := regex.FindAllString(string(body), -1)
	// Make each domain one-of-a-kind.
	uniqueDomains := makeUnique(domains)
	for i := 0; i < len(uniqueDomains); i++ {
		if validateDomain(uniqueDomains[i]) {
			// a file including all of the domains
			filePath, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			handleErrors(err)
			defer filePath.Close()
			fileContent := fmt.Sprint(uniqueDomains[i], "\n")
			_, err = filePath.WriteString(fileContent)
			handleErrors(err)
		}
	}
}

// Take a list of domains and make them one-of-a-kind
func makeUnique(randomStrings []string) []string {
	flag := make(map[string]bool)
	var uniqueString []string
	for _, randomString := range randomStrings {
		if !flag[randomString] {
			flag[randomString] = true
			uniqueString = append(uniqueString, randomString)
		}
	}
	return uniqueString
}

// Validate a website's domain
func validateDomain(domain string) bool {
	ns, _ := net.LookupNS(domain)
	return len(ns) >= 1
}

// Make a decision about how to handle errors.
func handleErrors(err error) {
	if err != nil {
		log.Println(err)
	}
}
