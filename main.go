/* Update lists with a simple application. */
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
	adwareConfig  = "configs/adware.conf"
	malwareConfig = "configs/malware.conf"
	privacyConfig = "configs/privacy.conf"
	sexualConfig  = "configs/sexual.conf"
	socialConfig  = "configs/social.conf"
)

func init() {
	//
}

func main() {
	validateAndSave("https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts", adwareConfig)
}

func validateAndSave(url, path string) {
	response, err := http.Get(url)
	handleErrors(err)
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	handleErrors(err)
	regex := regexp.MustCompile(`(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z0-9][a-z0-9-]{0,61}[a-z0-9]`)
	domains := regex.FindAllString(string(body), -1)
	// unique
	uniqueDomains := makeUnique(domains)
	for i := 0; i < len(uniqueDomains); i++ {
		if validateDomain(uniqueDomains[i]) {
			if fileExists(path) {
				os.Remove(path)
			}
			if !fileExists(path) {
				os.Create(path)
			}
			filePath, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			handleErrors(err)
			defer filePath.Close()
			fileContent := fmt.Sprint(uniqueDomains[i], "\n")
			_, err = filePath.WriteString(fileContent)
			handleErrors(err)
		}
	}
}

// Take in a list of domain and make them uniquie
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

// Validate a domain
func validateDomain(domain string) bool {
	ns, _ := net.LookupNS(domain)
	return len(ns) >= 1
}

// Decide what to do with errors
func handleErrors(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Check if a file exists
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
