package source

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/VampXDH/ipfinder/internal/common"
)

// ===================== INTERFACE =====================

type Source interface {
	Name() string
	Query(ctx context.Context, ip string, client *http.Client) ([]string, error)
}

// ===================== RapidDNS =====================

type RapidDNS struct{}

func (r RapidDNS) Name() string { return "rapiddns" }

func (r RapidDNS) Query(ctx context.Context, ip string, client *http.Client) ([]string, error) {
	url := fmt.Sprintf("https://rapiddns.io/sameip/%s?full=1#result", ip)
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	req.Header.Set("User-Agent", common.GetRandomUserAgent())
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Referer", "https://rapiddns.io/")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	html := string(body)

	re := regexp.MustCompile(`<td[^>]*>\s*([a-zA-Z0-9.-]+\.[a-zA-Z]{2,})\s*</td>`)
	matches := re.FindAllStringSubmatch(html, -1)

	var domains []string
	for _, m := range matches {
		if len(m) > 1 {
			domain := common.NormalizeDomain(m[1])
			if domain != "" && !strings.Contains(domain, "rapiddns") {
				domains = append(domains, domain)
			}
		}
	}
	return common.UniqueStrings(domains), nil
}

// ===================== WebScan =====================

type WebScan struct{}

func (w WebScan) Name() string { return "webscan" }

func (w WebScan) Query(ctx context.Context, ip string, client *http.Client) ([]string, error) {
	url := fmt.Sprintf("https://api.webscan.cc/?action=query&ip=%s", ip)
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	req.Header.Set("User-Agent", common.GetRandomUserAgent())
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	if err := json.Unmarshal(body, &results); err != nil {
		return nil, err
	}

	var domains []string
	for _, item := range results {
		if domain, ok := item["domain"].(string); ok {
			domain = common.NormalizeDomain(domain)
			if domain != "" {
				domains = append(domains, domain)
			}
		}
	}
	return common.UniqueStrings(domains), nil
}

// ===================== TNTcode =====================

type TNTcode struct{}

func (t TNTcode) Name() string { return "tntcode" }

func (t TNTcode) Query(ctx context.Context, ip string, client *http.Client) ([]string, error) {
	url := fmt.Sprintf("https://domains.tntcode.com/ip/%s", ip)
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	req.Header.Set("User-Agent", common.GetRandomUserAgent())
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	html := string(body)

	re := regexp.MustCompile(`<textarea[^>]*>([\s\S]*?)</textarea>`)
	matches := re.FindAllStringSubmatch(html, -1)

	var domains []string
	for _, m := range matches {
		if len(m) > 1 {
			lines := strings.Split(m[1], "\n")
			for _, line := range lines {
				domain := common.NormalizeDomain(strings.TrimSpace(line))
				if domain != "" && !strings.Contains(domain, "tntcode") {
					domains = append(domains, domain)
				}
			}
		}
	}
	return common.UniqueStrings(domains), nil
}

// ===================== NetworksDB =====================

type NetworksDB struct{}

func (n NetworksDB) Name() string { return "networksdb" }

func (n NetworksDB) Query(ctx context.Context, ip string, client *http.Client) ([]string, error) {
	url := fmt.Sprintf("https://networksdb.io/domains-on-ip/%s", ip)
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	req.Header.Set("User-Agent", common.GetRandomUserAgent())
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	html := string(body)

	re := regexp.MustCompile(`<pre[^>]*class="[^"]*threecols[^"]*"[^>]*>([\s\S]*?)</pre>`)
	matches := re.FindAllStringSubmatch(html, -1)

	var domains []string
	for _, m := range matches {
		if len(m) > 1 {
			lines := strings.Split(m[1], "\n")
			for _, line := range lines {
				domain := common.NormalizeDomain(strings.TrimSpace(line))
				if domain != "" && !strings.Contains(domain, "networksdb") {
					domains = append(domains, domain)
				}
			}
		}
	}
	return common.UniqueStrings(domains), nil
}

// ===================== Chaxunle =====================

type Chaxunle struct{}

func (c Chaxunle) Name() string { return "chaxunle" }

func (c Chaxunle) Query(ctx context.Context, ip string, client *http.Client) ([]string, error) {
	url := fmt.Sprintf("https://www.chaxunle.cn/ip/%s.html", ip)
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	req.Header.Set("User-Agent", common.GetRandomUserAgent())
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	html := string(body)

	re := regexp.MustCompile(`([a-zA-Z0-9][a-zA-Z0-9.-]+\.[a-zA-Z]{2,})`)
	matches := re.FindAllStringSubmatch(html, -1)

	var domains []string
	for _, m := range matches {
		if len(m) > 1 {
			domain := common.NormalizeDomain(m[1])
			if domain != "" &&
				!strings.Contains(domain, "chaxunle") &&
				!strings.Contains(domain, "baidu") &&
				!strings.Contains(domain, "qq.com") {
				domains = append(domains, domain)
			}
		}
	}
	return common.UniqueStrings(domains), nil
}

// ===================== THC.org =====================

type THCOrg struct{}

func (t THCOrg) Name() string { return "thc-org" }

func (t THCOrg) Query(ctx context.Context, ip string, client *http.Client) ([]string, error) {
	baseURL := fmt.Sprintf("https://ip.thc.org/%s", ip)

	var allDomains []string
	nextPage := baseURL
	pageCount := 0

	for nextPage != "" && pageCount < 20 {
		select {
		case <-ctx.Done():
			return common.UniqueStrings(allDomains), ctx.Err()
		default:
		}

		req, err := http.NewRequestWithContext(ctx, "GET", nextPage, nil)
		if err != nil {
			return allDomains, err
		}
		req.Header.Set("User-Agent", common.GetRandomUserAgent())
		req.Header.Set("Accept", "text/plain")

		resp, err := client.Do(req)
		if err != nil {
			return allDomains, err
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return allDomains, err
		}

		domains, newNext := parseTHCResponse(string(body))
		allDomains = append(allDomains, domains...)
		nextPage = newNext
		pageCount++

		time.Sleep(100 * time.Millisecond)
	}

	return common.UniqueStrings(allDomains), nil
}

func parseTHCResponse(body string) ([]string, string) {
	var domains []string
	var nextPage string

	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	cleaned := ansiRegex.ReplaceAllString(body, "")

	scanner := bufio.NewScanner(strings.NewReader(cleaned))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.Contains(line, "Next Page:") {
			parts := strings.SplitN(line, "Next Page:", 2)
			if len(parts) > 1 {
				nextPage = strings.TrimSpace(parts[1])
				continue
			}
		}

		if strings.HasPrefix(line, ";") || strings.HasPrefix(line, ";;") || line == "" {
			continue
		}

		if strings.Contains(line, ".") && !strings.Contains(line, " ") {
			domain := common.NormalizeDomain(line)
			if domain != "" {
				domains = append(domains, domain)
			}
		}
	}

	return domains, nextPage
}

// ===================== NOTE UNTUK KONTRIBUTOR =====================
//
// Untuk menambah source baru:
//
// 1. Buat struct baru, contoh:
/*
   type MyAPI struct{}

   func (m MyAPI) Name() string { return "myapi" }

   func (m MyAPI) Query(ip string, client *http.Client) ([]string, error) {
       // implement logic call API
   }
*/
//
// 2. Tambahkan ke slice `sources` di `internal/scanner/scanner.go`:
//    sources: []source.Source{
//        source.RapidDNS{},
//        ...
//        source.MyAPI{},
//    }
//
