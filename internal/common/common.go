package common

import (
	"math/rand"
	"net"
	"regexp"
	"strings"
	"time"
)

var (
	reProto = regexp.MustCompile(`^https?://`)
	reBad   = regexp.MustCompile(`[^a-z0-9.-]`)
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var browsers = []string{
	"Chrome",
	"Firefox",
	"Safari",
	"Edge",
	"Opera",
}

var osList = []string{
	"Windows NT 10.0",
	"Windows NT 6.3",
	"Windows NT 6.2",
	"Windows NT 6.1",
	"Macintosh; Intel Mac OS X 10_15",
	"Macintosh; Intel Mac OS X 10_14",
	"Macintosh; Intel Mac OS X 10_13",
	"X11; Linux x86_64",
	"X11; Ubuntu; Linux x86_64",
}

var chromeVersions = []string{
	"120.0.0.0", "119.0.0.0", "118.0.0.0", "117.0.0.0",
	"116.0.0.0", "115.0.0.0", "114.0.0.0", "113.0.0.0",
	"112.0.0.0", "111.0.0.0", "110.0.0.0", "109.0.0.0",
	"108.0.0.0", "107.0.0.0", "106.0.0.0", "105.0.0.0",
	"104.0.0.0", "103.0.0.0", "102.0.0.0", "101.0.0.0",
}

var firefoxVersions = []string{
	"120.0", "119.0", "118.0", "117.0", "116.0", "115.0",
	"114.0", "113.0", "112.0", "111.0", "110.0", "109.0",
	"108.0", "107.0", "106.0", "105.0", "104.0", "103.0",
	"102.0", "101.0",
}

func GetRandomUserAgent() string {
	browser := browsers[rand.Intn(len(browsers))]
	os := osList[rand.Intn(len(osList))]

	switch browser {
	case "Chrome":
		version := chromeVersions[rand.Intn(len(chromeVersions))]
		if strings.Contains(os, "Windows") {
			return "Mozilla/5.0 (" + os + "; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/" + version + " Safari/537.36"
		}
		return "Mozilla/5.0 (" + os + ") AppleWebKit/537.36 (KHTML, like Gecko) Chrome/" + version + " Safari/537.36"

	case "Firefox":
		version := firefoxVersions[rand.Intn(len(firefoxVersions))]
		if strings.Contains(os, "Windows") {
			return "Mozilla/5.0 (" + os + "; Win64; x64; rv:" + version + ") Gecko/20100101 Firefox/" + version
		}
		return "Mozilla/5.0 (" + os + "; rv:" + version + ") Gecko/20100101 Firefox/" + version

	default:
		version := chromeVersions[rand.Intn(len(chromeVersions))]
		return "Mozilla/5.0 (" + os + "; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/" + version + " Safari/537.36"
	}
}

func RandomSleep(min, max int) {
	sleepTime := min + rand.Intn(max-min+1)
	time.Sleep(time.Duration(sleepTime) * time.Millisecond)
}

func NormalizeDomain(raw string) string {
	if raw == "" {
		return ""
	}

	raw = reProto.ReplaceAllString(raw, "")
	parts := strings.Split(raw, "/")
	raw = parts[0]
	parts = strings.Split(raw, ":")
	raw = parts[0]
	raw = strings.TrimPrefix(raw, "www.")
	raw = strings.ToLower(strings.TrimSpace(raw))

	if !strings.Contains(raw, ".") {
		return ""
	}

	raw = reBad.ReplaceAllString(raw, "")
	return raw
}

func UniqueStrings(slice []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, s := range slice {
		if !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}
	return result
}

func IsValidIP(ip string) bool {
	if ip == "" {
		return false
	}
	parsed := net.ParseIP(ip)
	return parsed != nil
}
