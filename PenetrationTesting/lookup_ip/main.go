package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
)

// Function to check if a string represents a valid IP address
func isValidIP(ip string) bool {
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil
}

// Function to perform IP lookup for a given domain and save the unique results to a file
func lookupIP(domain string, ipOutputFile, aliasOutputFile string) {
	// Open IP output file for appending
	outFile, err := os.OpenFile(ipOutputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	// Execute 'dig' command
	cmd := exec.Command("dig", "+short", domain)
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Failed to lookup IP for %s: %v", domain, err)
		return
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")

	// Process each line
	for _, line := range lines {
		result := strings.TrimSpace(line)
		elements := strings.Split(result, ":")

		// Process each element
		for _, element := range elements {
			trimmedElement := strings.TrimSpace(element)
			if isValidIP(trimmedElement) {
				// Write to IP output file
				_, err := outFile.WriteString(trimmedElement + "\n")
				if err != nil {
					log.Fatal(err)
				}
			} else {
				// Write to alias output file
				aliasFile, err := os.OpenFile(aliasOutputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					log.Fatal(err)
				}
				defer aliasFile.Close()

				_, err = aliasFile.WriteString(trimmedElement + "\n")
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}

func main() {
	// File containing the list of domains
	domainFile := "domains.txt"

	// Output file to store unique IP addresses
	ipOutputFile := "ips.txt"
	aliasOutputFile := "alias.txt"

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

	err = os.WriteFile(ipOutputFile, nil, 0644)
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(aliasOutputFile, nil, 0644)
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
		lookupIP(domain, ipOutputFile, aliasOutputFile)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("IP lookup completed. Unique IP results saved to %s and %s\n", ipOutputFile, aliasOutputFile)
}
