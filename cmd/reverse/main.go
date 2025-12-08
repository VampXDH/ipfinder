package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/VampXDH/ipfinder/internal/scanner"
)

func main() {
	var (
		singleIP   string
		ipFile     string
		outputFile string
		threads    int
		verbose    bool
		silent     bool
		noColor    bool
	)

	// Seed random (untuk User-Agent & delay)
	rand.Seed(time.Now().UnixNano())

	// Parse flags
	flag.StringVar(&singleIP, "d", "", "Single IP address to scan")
	flag.StringVar(&ipFile, "l", "", "File containing list of IPs")
	flag.StringVar(&outputFile, "o", "results/domains.txt", "Output file")
	flag.IntVar(&threads, "t", 30, "Number of concurrent threads")
	flag.BoolVar(&verbose, "v", false, "Verbose output")
	flag.BoolVar(&silent, "silent", false, "Silent mode (only shows count)")
	flag.BoolVar(&noColor, "no-color", false, "Disable color output")

	flag.Usage = func() {
		fmt.Println(`
   _      ____         __       
  (_)__  / _(_)__  ___/ /__ ____
 / / _ \/ _/ / _ \/ _  / -_) __/
/_/ .__/_//_/_//_/\_,_/\__/_/   
  /_/                           

Reverse IP Finder v3.0 - All Sources

Usage: reverseip [options]

Options:
  -d string      Single IP address to scan
  -l string      File containing list of IPs
  -o string      Output file (default: results/domains.txt)
  -t int         Number of concurrent threads (default: 30)
  -v             Verbose output
  -silent        Silent mode (only shows count)
  -no-color      Disable color output
  -h, -help      Show this help message

Examples:
  reverseip -d 8.8.8.8
  reverseip -l ips.txt -t 100 -o results.txt
  reverseip -d 1.1.1.1 -v
  reverseip -l ips.txt -silent`)
	}

	flag.Parse()

	// Show help kalau gak ada flag sama sekali
	if flag.NFlag() == 0 {
		flag.Usage()
		return
	}

	// Validasi input
	if singleIP == "" && ipFile == "" {
		if noColor {
			fmt.Println("[ERROR] Either -d (single IP) or -l (IP list file) must be specified")
		} else {
			fmt.Println("\033[91m[ERROR] Either -d (single IP) or -l (IP list file) must be specified\033[0m")
		}
		os.Exit(1)
	}

	// Load IPs
	var ipList []string
	if singleIP != "" {
		ipList = []string{strings.TrimSpace(singleIP)}
	} else {
		f, err := os.Open(ipFile)
		if err != nil {
			if noColor {
				fmt.Printf("[ERROR] File error: %v\n", err)
			} else {
				fmt.Printf("\033[91m[ERROR] File error: %v\033[0m\n", err)
			}
			os.Exit(1)
		}
		defer f.Close()

		sc := bufio.NewScanner(f)
		for sc.Scan() {
			ip := strings.TrimSpace(sc.Text())
			if ip == "" {
				continue
			}
			if strings.HasPrefix(ip, "#") || strings.HasPrefix(ip, "//") {
				continue
			}
			ipList = append(ipList, ip)
		}

		if len(ipList) == 0 {
			if noColor {
				fmt.Println("[ERROR] No valid IPs found in file")
			} else {
				fmt.Println("\033[91m[ERROR] No valid IPs found in file\033[0m")
			}
			os.Exit(1)
		}
	}

	// Banner
	if !silent {
		fmt.Println(`
   _      ____         __       
  (_)__  / _(_)__  ___/ /__ ____
 / / _ \/ _/ / _ \/ _  / -_) __/
/_/ .__/_//_/_//_/\_,_/\__/_/   
  /_/                           `)
		if noColor {
			fmt.Println("\nReverse IP Finder v3.0 - All Sources Enabled\n")
		} else {
			fmt.Println("\n\033[92mReverse IP Finder v3.0 - All Sources Enabled\033[0m\n")
		}
	}

	// Run scanner
	s := scanner.NewScanner(ipList, outputFile, threads, verbose, silent, noColor)
	if err := s.Run(); err != nil {
		if noColor {
			fmt.Printf("[ERROR] Scanner error: %v\n", err)
		} else {
			fmt.Printf("\033[91m[ERROR] Scanner error: %v\033[0m\n", err)
		}
		os.Exit(1)
	}
}
