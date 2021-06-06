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
)

var (
	localHost        = "configs/host"
	localExclusion   = "configs/exclusion"
	foundDomains     []string
	exclusionDomains []string
	err              error
)

func init() {
	// Remove the file
	if fileExists(localHost) {
		err = os.Remove(localHost)
		handleErrors(err)
	}
	// Read Exclusion
	if fileExists(localExclusion) {
		file, err := os.Open(localExclusion)
		handleErrors(err)
		defer file.Close()
		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			exclusionDomains = append(exclusionDomains, scanner.Text())
		}
	}
}

func main() {
	startScraping()
}

func startScraping() {
	// Replace the URLs in this section to create your own list or add new lists.
	urls := []string{
		"https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts",
		"https://raw.githubusercontent.com/lightswitch05/hosts/master/docs/lists/ads-and-tracking-extended.txt",
		"https://raw.githubusercontent.com/notracking/hosts-blocklists/master/unbound/unbound.blacklist.conf",
		"https://raw.githubusercontent.com/lightswitch05/hosts/lists/tracking-aggressive-extended.txt",
		"https://raw.githubusercontent.com/lightswitch05/hosts/lists/facebook-extended.txt",
		"https://raw.githubusercontent.com/lightswitch05/hosts/lists/hate-and-junk-extended.txt",
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
		"https://raw.githubusercontent.com/Kees1958/W3C_annual_most_used_survey_blocklist/master/TOP_EU_US_Ads_Trackers_HOST",
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
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	handleErrors(err)
	// locate all domains
	regex := regexp.MustCompile(`(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z0-9][a-z0-9-]{0,61}[a-z0-9]`)
	foundDomains = regex.FindAllString(string(body), -1)
	uniqueDomains()
}

func uniqueDomains() {
	// Make each domain one-of-a-kind.
	uniqueDomains := makeUnique(foundDomains)
	// Remove all the exclusions domains from the list.
	for a := 0; a < len(exclusionDomains); a++ {
		uniqueDomains = removeStringFromSlice(uniqueDomains, exclusionDomains[a])
	}
	// Make everything one more time unique and then save it.
	for i := 0; i < len(uniqueDomains); i++ {
		// Validate all the domains
		if validateDomainViaNameServer(uniqueDomains[i]) {
			// Keep a list of all the valid domains.
			filePath, err := os.OpenFile(localHost, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
	for i := 0; i < len(randomStrings); i++ {
		if !flag[randomStrings[i]] {
			flag[randomStrings[i]] = true
			uniqueString = append(uniqueString, randomStrings[i])
		}
	}
	return uniqueString
}

// Validate a website's domain using the name server
func validateDomainViaNameServer(domain string) bool {
	ns, _ := net.LookupNS(domain)
	return len(ns) >= 1
}

// Validate the URI
func validURL(uri string) bool {
	validUri, err := url.ParseRequestURI(uri)
	if err != nil {
		return false
	}
	_ = validUri
	return true
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
