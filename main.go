/* Update lists with a simple application. */
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
)

func main() {
	response, err := http.Get("https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts")
	handleErrors(err)
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	handleErrors(err)
	regex := regexp.MustCompile(`(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z0-9][a-z0-9-]{0,61}[a-z0-9]`)
	comments := regex.FindAllString(string(body), -1)
	for _, comment := range comments {
		filePath, err := os.OpenFile("file-name", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		handleErrors(err)
		defer filePath.Close()
		fileContent := fmt.Sprint(comment, "\n")
		_, err = filePath.WriteString(fileContent)
		handleErrors(err)
		log.Println(comment)
	}
}

func handleErrors(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
