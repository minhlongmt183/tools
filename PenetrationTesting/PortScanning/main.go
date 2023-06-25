package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/panjf2000/ants"
)

type T struct {
	IP    string
	Ports string
}

func main() {
	var ipPortListFileName string
	var nmapOutputDirectory string
	var parallelNmapAmount int
	flag.StringVar(&ipPortListFileName, "i", "", "IP Address:Port list")
	flag.StringVar(&nmapOutputDirectory, "o", "nmap_reports", "Nmap directory to store output")
	flag.IntVar(&parallelNmapAmount, "t", 5, "Amount of Nmap running at the time")
	flag.Parse()

	if _, err := os.Stat(nmapOutputDirectory); os.IsNotExist(err) {
		err := os.Mkdir(nmapOutputDirectory, os.ModePerm)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to create output directory")
			os.Exit(1)
		}
	}

	var wg sync.WaitGroup
	pool, _ := ants.NewPoolWithFunc(parallelNmapAmount, func(i interface{}) {
		defer wg.Done()
		input := i.(T)
		nmapCmd := fmt.Sprintf("timeout -k 1m 10m nmap -sV -sC -v -T4 --open -p %s %s -oN %s/%s.txt -oX %s/%s.xml", input.Ports, input.IP,
			nmapOutputDirectory, input.IP,
			nmapOutputDirectory, input.IP)
		cmd := exec.Command("bash", "-c", nmapCmd)
		err := cmd.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error %s: %s\n", cmd, err)
		}
	})
	defer pool.Release()

	ipMap := treemap.NewWithStringComparator()
	loadIpAddressInfo(ipPortListFileName, ipMap)

	fmt.Fprintf(os.Stdout, "Total: %d\n", ipMap.Size())
	current := 0
	it := ipMap.Iterator()
	for it.Next() {
		key, value := it.Key(), it.Value()
		ports := strings.Join(value.([]string), ",")
		task := T{
			IP:    key.(string),
			Ports: ports,
		}
		current++
		fmt.Fprintf(os.Stdout, "Current progress: %d/%d\n", current, ipMap.Size())
		wg.Add(1)
		pool.Invoke(task)
	}
	wg.Wait()
}

func loadIpAddressInfo(ipPortFile string, ipInfoMap *treemap.Map) {
	f, err := os.Open(ipPortFile)
	if err != nil {
		log.Fatalf("Failed to open ipAddressList: %s", err)
	}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "http") {
			if u, err := url.Parse(line); err == nil {
				line = u.Host
			}
		}

		lineArgs := strings.SplitN(line, ":", 2)
		if len(lineArgs) != 2 {
			continue
		}
		if portList, found := ipInfoMap.Get(lineArgs[0]); found {
			portListSlice := portList.([]string)
			portListSlice = append(portListSlice, lineArgs[1])
			ipInfoMap.Put(lineArgs[0], portListSlice)
		} else {
			ipInfoMap.Put(lineArgs[0], []string{lineArgs[1]})
		}
	}
}
