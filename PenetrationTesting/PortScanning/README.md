# enmap
---
This tool is written to support oscp exam.

This will scan all port TCP and UDP parallel.

# p-scan

p-scan is a Go-based tool for parallelized Nmap scanning. It allows you to scan multiple IP addresses and ports simultaneously using Nmap.

## Usage


Run p-scan with the desired options.

```bash
./p-scan -i <ipPortListFileName> -o <nmapOutputDirectory> -t <parallelNmapAmount>
```
Replace the following parameters:  

- <ipPortListFileName>: The file containing the list of IP addresses and ports to scan.
- <nmapOutputDirectory>: The directory to store the Nmap output files.
- <parallelNmapAmount>: The number of Nmap instances to run in parallel.

## Example
Here's an example usage:

```bash
./p-scan -i ip-port-list.txt -o nmap_reports -t 5
```
In this example, the tool reads the IP addresses and ports from the ip-port-list.txt file and runs Nmap scans on them. The output files will be stored in the nmap_reports directory. The tool will run 5 Nmap instances in parallel.

## Note
- Make sure to have RustScan installed on your system for optimal scanning results.

- Consider increasing the file limit if you encounter the warning "File limit is lower than the default batch size. Consider upping with --ulimit. May cause harm to sensitive servers."

- This tool is inspired by theblackturtle.

## Recommendations
For optimal scanning workflow, you can follow these recommendations:

Use RustScan to identify open ports:

```bash
rustscan --scripts None --range 1-65535 -b 20000 -a {{.upAddr}} --accessible -t 2500 > raw-open-ports.txt
```
This command scans all ports (1-65535) for the specified IP address and outputs the raw results to the raw-open-ports.txt file.

Extract the open ports from the raw results:

```bash
cat raw-open-ports.txt | grep 'Open ' | awk '{print $2}' > open-ports.txt
```
This command filters the raw results to extract the open ports and stores them in the open-ports.txt file.

Use p-scan to perform Nmap scanning on the extracted open ports:

```bash
p-scan -t 4 -i open-ports.txt -o nmap_output_folder
```
This command uses p-scan to perform parallelized Nmap scanning on the open ports. Specify the number of parallel Nmap instances with -t, provide the open-ports.txt file as input with -i, and specify the output directory with -o.

Please note that this README.md file assumes the presence of RustScan and its usage before running p-scan. Adjust the recommendations section based on your requirements and the availability of tools on your system.

Hope this helps! Let me know if you have any more questions.