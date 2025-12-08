package scanner

import (
	"os"
	"sync"
)

type OutputWriter struct {
	file     *os.File
	mu       sync.Mutex
	seen     map[string]bool
	filename string
	count    int
}

func NewOutputWriter(filename string) *OutputWriter {
	// Create directory if needed
	_ = os.MkdirAll("results", 0755)

	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}

	return &OutputWriter{
		file:     file,
		seen:     make(map[string]bool),
		filename: filename,
		count:    0,
	}
}

func (w *OutputWriter) Write(domain string) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.seen[domain] {
		return
	}

	w.seen[domain] = true
	w.count++
	w.file.WriteString(domain + "\n")
}

func (w *OutputWriter) Close() {
	w.file.Close()
}

func (w *OutputWriter) Count() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.count
}
