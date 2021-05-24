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

func main() {
	response, err := http.Get("https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts")
	handleErrors(err)
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	handleErrors(err)
	regex := regexp.MustCompile(`(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z0-9][a-z0-9-]{0,61}[a-z0-9]`)
	domains := regex.FindAllString(string(body), -1)
	for _, domain := range domains {
		if validateDomain(domain) {
			filePath, err := os.OpenFile(adwareConfig, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			handleErrors(err)
			defer filePath.Close()
			fileContent := fmt.Sprint(domain, "\n")
			_, err = filePath.WriteString(fileContent)
			handleErrors(err)
			log.Println(comment)
		}
	}
}

func validateDomain(domain string) bool {
	ns, _ := net.LookupNS(domain)
	return len(ns) >= 1
}

func handleErrors(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
