# ipfinder
`ipfinder` is an IP information lookup tool capable of performing reverse IP searches, allowing users to identify domains hosted on the same IP address. The tool leverages passive online data sources to ensure accurate and reliable results.

`ipfinder` features a simple modular architecture optimized for speed. It is designed to perform one primary function—reverse IP lookup, and it excels at doing exactly that.

We built it in compliance with all applicable passive-source licenses and usage limitations. Its passive model ensures speed, efficiency, and confidentiality, making it a valuable asset for penetration testers and bug bounty hunters.


# Features
<h1 align="left">
  <img src="static/ipfinder-run.png" alt="ipfinder" width="700px"></a>
  <br>
</h1>


# Usage

```sh
ipfinder -h
```

This will display help for the tool. Here are all the switches it supports.

```yaml
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
  ipfinder -d 8.8.8.8
  ipfinder -l ips.txt -t 100 -o results.txt
  ipfinder -d 1.1.1.1 -v
  ipfinder -l ips.txt -silent
```

# Installation
```sh
go install -v github.com/VampXDH/ipfinder/cmd/ipfinder@latest
```
