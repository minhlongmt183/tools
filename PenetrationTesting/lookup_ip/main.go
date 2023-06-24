package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
)

// Function to perform IP lookup for a given domain and save the unique results to a file
func lookupIP(domain string) string {
	ip, err := net.LookupIP(domain)
	if err != nil {
		return fmt.Sprintf("Failed to lookup IP for %s", domain)
	}

	if len(ip) > 0 {
		return fmt.Sprintf("%s: %s", domain, ip[0])
	}

	return ""
}

func main() {
	// File containing the list of domains
	domainFile := "domains.txt"

	// Output file to store unique IP addresses
	outputFile := "ip.txt"

	// Check if the domain file argument is provided
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run lookup [domain_file]")
		fmt.Println("Please provide the path to the file containing the list of domains.")
		os.Exit(1)
	}

	// Assign the domain file argument to the variable
	domainFile = os.Args[1]

	// Check if the domain file exists
	_, err := os.Stat(domainFile)
	if os.IsNotExist(err) {
		fmt.Printf("The specified domain file does not exist: %s\n", domainFile)
		os.Exit(1)
	}

	// Clear the output file
	err = ioutil.WriteFile(outputFile, nil, 0644)
	if err != nil {
		log.Fatal(err)
	}

	// Read each line from the domain file
	file, err := os.Open(domainFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		domain := scanner.Text()
		result := lookupIP(domain)
		if !strings.Contains(result, "Failed to lookup IP") {
			ipFile, err := os.OpenFile(outputFile, os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				log.Fatal(err)
			}

			ipScanner := bufio.NewScanner(ipFile)
			exists := false
			for ipScanner.Scan() {
				if ipScanner.Text() == result {
					exists = true
					break
				}
			}

			ipFile.Close()

			if !exists {
				err := ioutil.WriteFile(outputFile, []byte(result+"\n"), 0644)
				if err != nil {
					log.Fatal(err)
				}
			}
		} else {
			fmt.Println(result)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("IP lookup completed. Unique IP results saved to %s\n", outputFile)
}
