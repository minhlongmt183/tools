# IP Lookup

IP Lookup is a Go program that performs IP lookup for a given list of domains and saves the unique IP addresses and aliases to separate files.

## Prerequisites

- Go (version 1.13 or higher)
- dig command-line tool

## Usage

1. Clone the repository or download the source code.

2. Navigate to the project directory.

3. Ensure that the `domains.txt` file is present in the project directory. This file should contain a list of domains, with each domain on a new line.

4. Open a terminal and run the following command to build the Go program:

   ```shell
   go build
   ```
   Run the program with the following command:

   ```shell
   ./ip-lookup [domain_file]
   ```
   Replace [domain_file] with the path to the file containing the list of domains. For example:
   
   ```shell
   ./ip-lookup domains.txt
   ```
   The program will perform IP lookup for each domain in the file and save the unique IP addresses to the ips.txt file, and aliases (non-IP results) to the alias.txt file.
   
   After the program finishes execution, you will see the following message:
   
   ```shell
   IP lookup completed. Unique IP results saved to ips.txt and alias.txt
   Open the ips.txt and alias.txt files to view the results.  
   ```
## Author 
This code was written by Edisc.