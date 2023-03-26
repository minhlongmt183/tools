package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
)

func main() {
	// Define command-line flags
	inputIP := flag.String("i", "", "target IP address")
	outputDir := flag.String("o", "nmap_reports", "output directory for Nmap reports, default: nmap_reports ")
	parallelism := flag.Int("t", 1, "number of parallel Nmap processes to run")

	// Parse command-line flags
	flag.Parse()

	// Check that required flags are set
	if *inputIP == "" {
		fmt.Println("Usage: ./portscan -i IP -o nmap_reports -t <Amount of Nmap running at time>")
		os.Exit(1)
	}

	// Create output directory if it doesn't exist
	if _, err := os.Stat(*outputDir); os.IsNotExist(err) {
		os.Mkdir(*outputDir, os.ModePerm)
	}

	// Define wait group to ensure all Nmap processes complete
	var wg sync.WaitGroup

	// Process targets in parallel
	ip := *inputIP
	output := *outputDir + "/"

	runRustscan(ip, output)

	//// Wait for available slot in semaphore
	//semaphore := make(chan bool, *parallelism)
	//semaphore <- true
	//
	//// Run Nmap process
	//wg.Add(1)
	//go func(target string) {
	//	defer func() { <-semaphore }()
	//	defer wg.Done()
	//
	//	report := fmt.Sprintf("%s/%s.xml", *outputDir, target)
	//	args := []string{"-oX", report, target}
	//	cmd := exec.Command("nmap", args...)
	//	err := cmd.Run()
	//	if err != nil {
	//		fmt.Printf("Error running Nmap on %s: %v\n", target, err)
	//	}
	//}(ip)

	// Wait for all Nmap processes to complete
	fmt.Println(parallelism)
	wg.Wait()
}

func runRustscan(ip string, outputDir string) {

	// Check if rustscan is installed
	_, err := exec.LookPath("rustscan")
	if err != nil {
		fmt.Println("rustscan is not installed. Please install it before running this program.")
		return
	}

	// Run rustscan command and capture output
	cmd := exec.Command("rustscan", "--scripts", "None", "--range", "1-65535", "-b", "20000", "-a", ip, "--accessible", "-t", "2500")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Error running rustscan command: %v\nstderr: %s", err, stderr.String())
	}

	// Write raw output to file
	err = os.WriteFile(outputDir+"raw-open-ports.txt", stdout.Bytes(), 0644)
	if err != nil {
		log.Fatalf("Error writing raw output to file: %v", err)
	}

	// Filter output to only include open ports
	var openPorts []string
	lines := strings.Split(stdout.String(), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Open ") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				openPorts = append(openPorts, fields[1])
			}
		}
	}

	// Write open ports to file
	openPortsStr := strings.Join(openPorts, "\n")
	file, err := os.Create(outputDir + "open-ports.txt")
	if err != nil {
		log.Fatalf("Error creating open-ports.txt file: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(openPortsStr)
	if err != nil {
		log.Fatalf("Error writing open ports to file: %v", err)
	}

	fmt.Println("Done")
}
