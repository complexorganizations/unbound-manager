package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sync"

	"github.com/openrdap/rdap"
)

var (
	localHost        = "configs/host"
	localExclusion   = "configs/exclusion"
	foundDomains     []string
	exclusionDomains []string
	err              error
	wg               sync.WaitGroup
)

func init() {
	// Remove the file
	if fileExists(localHost) {
		err = os.Remove(localHost)
		handleErrors(err)
	}
	// Read Exclusion
	if fileExists(localExclusion) {
		// exclusion domain
		exclusionDomains = readAndAppend(localExclusion, exclusionDomains)
	}
}

func main() {
	startScraping()
	// once done scraping, make everything unique.
	makeEverythingUnique()
}

func startScraping() {
	// Replace the URLs in this section to create your own list or add new lists.
	urls := []string{
		"https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts",
		"https://raw.githubusercontent.com/lightswitch05/hosts/master/docs/lists/ads-and-tracking-extended.txt",
		"https://raw.githubusercontent.com/notracking/hosts-blocklists/master/unbound/unbound.blacklist.conf",
		"https://raw.githubusercontent.com/lightswitch05/hosts/master/docs/lists/tracking-aggressive-extended.txt",
		"https://raw.githubusercontent.com/lightswitch05/hosts/master/docs/lists/facebook-extended.txt",
		"https://raw.githubusercontent.com/lightswitch05/hosts/master/docs/lists/hate-and-junk-extended.txt",
		"https://raw.githubusercontent.com/AdAway/adaway.github.io/master/hosts.txt",
		"https://raw.githubusercontent.com/AdguardTeam/AdguardFilters/master/EnglishFilter/sections/adservers.txt",
		"https://raw.githubusercontent.com/tg12/pihole-phishtank-list/master/list/phish_domains.txt",
		"https://raw.githubusercontent.com/HorusTeknoloji/TR-PhishingList/master/url-lists.txt",
		"https://raw.githubusercontent.com/Spam404/lists/master/main-blacklist.txt",
		"https://raw.githubusercontent.com/FadeMind/hosts.extras/master/add.Risk/hosts",
		"https://raw.githubusercontent.com/PolishFiltersTeam/KADhosts/master/KADhosts.txt",
		"https://raw.githubusercontent.com/FadeMind/hosts.extras/master/add.Spam/hosts",
		"https://raw.githubusercontent.com/matomo-org/referrer-spam-blacklist/master/spammers.txt",
		"https://raw.githubusercontent.com/VeleSila/yhosts/master/hosts",
		"https://raw.githubusercontent.com/RooneyMcNibNug/pihole-stuff/master/SNAFU.txt",
		"https://raw.githubusercontent.com/anudeepND/blacklist/master/adservers.txt",
		"https://raw.githubusercontent.com/FadeMind/hosts.extras/master/UncheckyAds/hosts",
		"https://raw.githubusercontent.com/bigdargon/hostsVN/master/hosts",
		"https://raw.githubusercontent.com/jdlingyu/ad-wars/master/hosts",
		"https://raw.githubusercontent.com/FadeMind/hosts.extras/master/add.2o7Net/hosts",
		"https://raw.githubusercontent.com/crazy-max/WindowsSpyBlocker/master/data/hosts/spy.txt",
		"https://raw.githubusercontent.com/Kees1958/W3C_annual_most_used_survey_blocklist/master/TOP_EU_US_Ads_Trackers_ABP",
		"https://raw.githubusercontent.com/Perflyst/PiHoleBlocklist/master/android-tracking.txt",
		"https://raw.githubusercontent.com/Perflyst/PiHoleBlocklist/master/SmartTV.txt",
		"https://raw.githubusercontent.com/Perflyst/PiHoleBlocklist/master/AmazonFireTV.txt",
		"https://raw.githubusercontent.com/DandelionSprout/adfilt/master/Alternate%20versions%20Anti-Malware%20List/AntiMalwareHosts.txt",
		"https://raw.githubusercontent.com/anudeepND/blacklist/master/facebook.txt",
	}
	for i := 0; i < len(urls); i++ {
		// Validate the URI before beginning the scraping process.
		if validURL(urls[i]) {
			saveTheDomains(urls[i])
		}
	}
}

func saveTheDomains(url string) {
	// Send a request to acquire all the information you need.
	response, err := http.Get(url)
	handleErrors(err)
	body, err := io.ReadAll(response.Body)
	handleErrors(err)
	bodyAsString := string(body)
	if bodyAsString == "404: Not Found" {
		log.Fatalln("Error: ", url)
	}
	// locate all domains
	regex := regexp.MustCompile(`(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z0-9][a-z0-9-]{0,61}[a-z0-9]`)
	foundDomains = regex.FindAllString(string(body), -1)
	defer response.Body.Close()
	// Make each domain one-of-a-kind.
	uniqueDomains := makeUnique(foundDomains)
	// Remove all the exclusions domains from the list.
	for a := 0; a < len(exclusionDomains); a++ {
		uniqueDomains = removeStringFromSlice(uniqueDomains, exclusionDomains[a])
	}
	fmt.Println("Domains:", len(uniqueDomains))
	for i := 0; i < len(uniqueDomains); i++ {
		wg.Add(1)
		go makeDomainsUnique(uniqueDomains[i])
	}
	wg.Wait()
}

func makeDomainsUnique(uniqueDomains string) {
	// Validate all the domains
	if validateDomainViaLookupNS(uniqueDomains) || validateDomainViaLookupAddr(uniqueDomains) || validateDomainViaLookupCNAME(uniqueDomains) || validateDomainViaLookupMX(uniqueDomains) || validateDomainViaLookupTXT(uniqueDomains) || domainRegistration(uniqueDomains) {
		// Keep a list of all the valid domains.
		writeToFile(localHost, uniqueDomains)
	}
	wg.Done()
}

// Take a list of domains and make them one-of-a-kind
func makeUnique(randomStrings []string) []string {
	flag := make(map[string]bool)
	var uniqueString []string
	for i := 0; i < len(randomStrings); i++ {
		if !flag[randomStrings[i]] {
			flag[randomStrings[i]] = true
			uniqueString = append(uniqueString, randomStrings[i])
		}
	}
	return uniqueString
}

// Validate a domain
func validateDomainViaLookupNS(domain string) bool {
	valid, _ := net.LookupNS(domain)
	return len(valid) >= 1
}

func validateDomainViaLookupAddr(domain string) bool {
	valid, _ := net.LookupAddr(domain)
	return len(valid) >= 1
}

func validateDomainViaLookupCNAME(domain string) bool {
	valid, _ := net.LookupCNAME(domain)
	return len(valid) >= 1
}

func validateDomainViaLookupMX(domain string) bool {
	valid, _ := net.LookupMX(domain)
	return len(valid) >= 1
}

func validateDomainViaLookupTXT(domain string) bool {
	valid, _ := net.LookupTXT(domain)
	return len(valid) >= 1
}

func domainRegistration(domain string) bool {
	client := &rdap.Client{}
	_, ok := client.QueryDomain(domain)
	return ok == nil
}

// Validate the URI
func validURL(uri string) bool {
	_, err = url.ParseRequestURI(uri)
	return err == nil
}

// Make a decision about how to handle errors.
func handleErrors(err error) {
	if err != nil {
		log.Println(err)
	}
}

// Check to see if a file already exists.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// Remove a string from a slice
func removeStringFromSlice(originalSlice []string, removeString string) []string {
	for i := 0; i < len(originalSlice); i++ {
		if originalSlice[i] == removeString {
			return append(originalSlice[:i], originalSlice[i+1:]...)
		}
	}
	return originalSlice
}

// Save to a file
func writeToFile(pathInSystem string, content string) {
	filePath, err := os.OpenFile(pathInSystem, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	handleErrors(err)
	_, err = filePath.WriteString(content + "\n")
	handleErrors(err)
	filePath.Close()
}

// Read and append to array
func readAndAppend(fileLocation string, arrayName []string) []string {
	file, err := os.Open(fileLocation)
	handleErrors(err)
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		arrayName = append(arrayName, scanner.Text())
	}
	file.Close()
	return arrayName
}

// make everything unique
func makeEverythingUnique() {
	var finalDomainList []string
	finalDomainList = readAndAppend(localHost, finalDomainList)
	uniqueDomains := makeUnique(finalDomainList)
	for i := 0; i < len(uniqueDomains); i++ {
		writeToFile(localHost, uniqueDomains[i])
	}
}
