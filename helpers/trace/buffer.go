package trace

import (
	"bufio"
	"errors"
	"fmt"
	"hash"
	"hash/crc32"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"sync"

	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"

	"gitlab.com/gitlab-org/gitlab-runner/helpers"
)

const defaultBytesLimit = 4 * 1024 * 1024 // 4MB

var errLogLimitExceeded = errors.New("log limit exceeded")

type Buffer struct {
	lock sync.RWMutex
	lw   *limitWriter
	w    io.WriteCloser

	logFile  *os.File
	bufw     *bufio.Writer
	checksum hash.Hash32

	// failedFlush indicates that a read which subsequentialy attempted to
	// flush data to the underlying writer failed. In this scenario, calls to
	// Write() will immediately attempt to flush and return any error on a
	// failure.
	failedFlush bool
}

type lengthSort []string

func (s lengthSort) Len() int {
	return len(s)
}

func (s lengthSort) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s lengthSort) Less(i, j int) bool {
	return len(s[i]) > len(s[j])
}

func (b *Buffer) SetMasked(values []string) {
	b.lock.Lock()
	defer b.lock.Unlock()

	// close existing writer to flush data
	if b.w != nil {
		b.w.Close()
	}

	defaultTransformers := []transform.Transformer{
		newSensitiveURLParamTransform(),
		encoding.Replacement.NewEncoder(),
	}

	transformers := make([]transform.Transformer, 0, len(values)+len(defaultTransformers))

	sort.Sort(lengthSort(values))
	for _, value := range values {
		transformers = append(transformers, newPhraseTransform(value))
	}

	transformers = append(transformers, defaultTransformers...)

	b.w = transform.NewWriter(b.lw, transform.Chain(transformers...))
}

func (b *Buffer) SetLimit(size int) {
	b.lock.Lock()
	defer b.lock.Unlock()

	b.lw.limit = int64(size)
}

func (b *Buffer) Size() int {
	b.lock.RLock()
	defer b.lock.RUnlock()

	if b.lw == nil {
		return 0
	}
	return int(b.lw.written)
}

func (b *Buffer) Reader(offset, n int) (io.Reader, error) {
	b.lock.Lock()
	defer b.lock.Unlock()

	// For simplicity, we read only from the file, rather than also the bufio.Writer.
	// To ensure the underlying file has the data requested, we always flush the
	// buffer.
	//
	// If a failure occurs on flushing the data, we store that an error occurred so
	// buffer.Write() can retry and additionally return any error on the write side.
	if err := b.bufw.Flush(); err != nil {
		b.failedFlush = true
		return nil, fmt.Errorf("flushing log buffer: %w", err)
	}

	return io.NewSectionReader(b.logFile, int64(offset), int64(n)), nil
}

func (b *Buffer) Write(p []byte) (int, error) {
	b.lock.Lock()
	defer b.lock.Unlock()

	// if we previously failed to flush to the underlying writer, try again
	// and return any failure immediately.
	if b.failedFlush {
		b.failedFlush = false
		if err := b.bufw.Flush(); err != nil {
			return 0, err
		}
	}

	n, err := b.w.Write(p)
	// if we get a log limit exceeded error, we've written the log limit
	// notice out to the log and will now silently not write any additional
	// data: we return len(p), nil so the caller continues as normal.
	if err == errLogLimitExceeded {
		return len(p), nil
	}
	return n, err
}

func (b *Buffer) Finish() {
	b.lock.RLock()
	defer b.lock.RUnlock()

	if b.w != nil {
		_ = b.w.Close()
	}
}

func (b *Buffer) Close() {
	_ = b.logFile.Close()
	_ = os.Remove(b.logFile.Name())
}

func (b *Buffer) Checksum() string {
	b.lock.RLock()
	defer b.lock.RUnlock()

	return fmt.Sprintf("crc32:%08x", b.checksum.Sum32())
}

type limitWriter struct {
	w       io.Writer
	written int64
	limit   int64
}

func (w *limitWriter) Write(p []byte) (int, error) {
	capacity := w.limit - w.written

	if capacity <= 0 {
		return 0, errLogLimitExceeded
	}

	if int64(len(p)) >= capacity {
		p = p[:capacity]
		n, err := w.w.Write(p)
		if err == nil {
			err = errLogLimitExceeded
		}
		if n < 0 {
			n = 0
		}
		w.written += int64(n)
		w.writeLimitExceededMessage()

		return n, err
	}

	n, err := w.w.Write(p)
	if n < 0 {
		n = 0
	}
	w.written += int64(n)
	return n, err
}

func (w *limitWriter) writeLimitExceededMessage() {
	n, _ := fmt.Fprintf(
		w.w,
		"\n%sJob's log exceeded limit of %v bytes.\n"+
			"Job execution will continue but no more output will be collected.%s\n",
		helpers.ANSI_BOLD_YELLOW,
		w.limit,
		helpers.ANSI_RESET,
	)
	w.written += int64(n)
}

func New() (*Buffer, error) {
	logFile, err := ioutil.TempFile("", "trace")
	if err != nil {
		return nil, err
	}

	buffer := &Buffer{
		logFile:  logFile,
		bufw:     bufio.NewWriter(logFile),
		checksum: crc32.NewIEEE(),
	}

	buffer.lw = &limitWriter{
		w:       io.MultiWriter(buffer.bufw, buffer.checksum),
		written: 0,
		limit:   defaultBytesLimit,
	}

	buffer.SetMasked(nil)

	return buffer, nil
}
