package scanner

import (
	"bufio"
	"os"
	"sync"
	"sync/atomic"
)

type OutputWriter struct {
	file     *os.File
	writer   *bufio.Writer
	mu       sync.Mutex
	seen     map[string]bool
	filename string
	count    uint64
}

func NewOutputWriter(filename string) (*OutputWriter, error) {
	_ = os.MkdirAll("results", 0755)

	file, err := os.Create(filename)
	if err != nil {
		return nil, err
	}

	writer := bufio.NewWriter(file)

	return &OutputWriter{
		file:     file,
		writer:   writer,
		seen:     make(map[string]bool),
		filename: filename,
		count:    0,
	}, nil
}

func (w *OutputWriter) Write(domain string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.seen[domain] {
		return nil
	}

	_, err := w.writer.WriteString(domain + "\n")
	if err != nil {
		return err
	}

	w.seen[domain] = true
	atomic.AddUint64(&w.count, 1)
	return nil
}

func (w *OutputWriter) Close() {
	w.writer.Flush()
	w.file.Close()
}

func (w *OutputWriter) Count() int {
	return int(atomic.LoadUint64(&w.count))
}
