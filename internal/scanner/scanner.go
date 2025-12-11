package scanner

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/VampXDH/ipfinder/internal/common"
	"github.com/VampXDH/ipfinder/internal/logger"
	"github.com/VampXDH/ipfinder/internal/source"
)

type Scanner struct {
	ipList     []string
	outputFile string
	threads    int
	verbose    bool
	silent     bool
	noColor    bool

	client    *http.Client
	writer    *OutputWriter
	log       *logger.Logger
	sources   []source.Source
	startTime time.Time
	ctx       context.Context
}

func NewScanner(ctx context.Context, ipList []string, outputFile string, threads int, verbose, silent, noColor bool) *Scanner {
	return &Scanner{
		ipList:     ipList,
		outputFile: outputFile,
		threads:    threads,
		verbose:    verbose,
		silent:     silent,
		noColor:    noColor,
		client: &http.Client{
			Timeout: 120 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
		},
		log: &logger.Logger{
			Silent:  silent,
			Verbose: verbose,
			NoColor: noColor,
		},
		sources: []source.Source{
			source.RapidDNS{},
			source.TNTcode{},
			source.WebScan{},
			source.NetworksDB{},
			source.Chaxunle{},
			source.THCOrg{},
		},
		ctx: ctx,
	}
}

func (s *Scanner) Run() error {
	s.startTime = time.Now()

	s.log.Info("Loaded %d IPs", len(s.ipList))
	s.log.Info("Using %d sources", len(s.sources))
	s.log.Info("Threads: %d", s.threads)
	s.log.Info("Output: %s", s.outputFile)
	s.log.Info("Sources: rapiddns, tntcode, webscan, networksdb, chaxunle, thc-org")
	s.log.Line()

	writer, err := NewOutputWriter(s.outputFile)
	if err != nil {
		return err
	}
	s.writer = writer
	defer s.writer.Close()

	sem := make(chan struct{}, s.threads)
	var wg sync.WaitGroup

	for _, ip := range s.ipList {
		select {
		case <-s.ctx.Done():
			s.log.Info("Cancelled by user")
			return s.ctx.Err()
		default:
		}

		ip := ip
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			s.processIP(ip)
			<-sem
		}()
	}

	wg.Wait()
	s.printSummary()
	return nil
}

func (s *Scanner) processIP(ip string) {
	s.log.Info("Processing: %s", ip)

	ipDomains := 0

	for _, src := range s.sources {
		select {
		case <-s.ctx.Done():
			s.log.Info("Cancelled by user")
			return
		default:
		}

		common.RandomSleep(1000, 3000)

		s.log.Verbosef("Querying %s with %s", ip, src.Name())

		domains, err := src.Query(ip, s.client)
		if err != nil {
			s.log.Warning("%s error for %s: %v", src.Name(), ip, err)
			continue
		}

		if len(domains) > 0 {
			count := len(domains)
			ipDomains += count

			for _, d := range domains {
				err := s.writer.Write(d)
				if err != nil {
					s.log.Warning("Failed to write domain %s: %v", d, err)
				}
			}

			s.log.Success(src.Name(), ip, count)
		} else {
			s.log.Verbosef("%s: %s - 0 domains", src.Name(), ip)
		}
	}

	if ipDomains > 0 {
		s.log.Info("%s: Total %d domains", ip, ipDomains)
	} else {
		s.log.Info("%s: No domains found", ip)
	}
}

func (s *Scanner) printSummary() {
	elapsed := time.Since(s.startTime)
	unique := s.writer.Count()

	if s.silent {
		// Silent mode cuma print jumlah domain unik
		println(unique)
		return
	}

	s.log.Line()
	s.log.Stat("=", 50)
	s.log.Info("Completed in %v", elapsed.Round(time.Second))
	s.log.Info("IPs processed: %d", len(s.ipList))
	s.log.Info("Sources used: %d", len(s.sources))
	s.log.Info("Unique domains: %d", unique)
	s.log.Info("Saved to: %s", s.outputFile)
	s.log.Stat("=", 50)
}
