package trace

import (
	"math"
	"sync"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVariablesMasking(t *testing.T) {
	//nolint:lll
	input := "This is the secret message cont@ining :secret duplicateValues ffixx prefix prefix_mask suffix mask_suffix middle dd"

	maskedValues := []string{
		"is",
		"duplicateValue",
		"duplicateValue",
		":secret",
		"cont@ining",
		"fix",
		"prefix",
		"prefix_mask",
		"suffix",
		"mask_suffix",
		"dd",
		"middle",
	}

	buffer, err := New()
	require.NoError(t, err)
	defer buffer.Close()

	buffer.SetMasked(maskedValues)

	_, err = buffer.Write([]byte(input))
	require.NoError(t, err)

	buffer.Finish()

	content, err := buffer.Bytes(0, 1000)
	require.NoError(t, err)

	//nolint:lll
	assert.Equal(t, "Th[MASKED] [MASKED] the secret message [MASKED] [MASKED] [MASKED]s f[MASKED]x [MASKED] [MASKED] [MASKED] [MASKED] [MASKED] [MASKED]", string(content))
}

func TestTraceLimit(t *testing.T) {
	traceMessage := "This is the long message"

	buffer, err := New()
	require.NoError(t, err)
	defer buffer.Close()

	buffer.SetLimit(10)
	assert.Equal(t, 0, buffer.Size())

	for i := 0; i < 100; i++ {
		n, err := buffer.Write([]byte(traceMessage))
		require.NoError(t, err)
		require.Greater(t, n, 0)
	}

	buffer.Finish()

	content, err := buffer.Bytes(0, 1000)
	require.NoError(t, err)

	expectedContent := "This is th\n" +
		"\x1b[33;1mJob's log exceeded limit of 10 bytes.\n" +
		"Job execution will continue but no more output will be collected.\x1b[0;m\n"
	assert.Equal(t, len(expectedContent), buffer.Size(), "unexpected buffer size")
	assert.Equal(t, "crc32:295921ca", buffer.Checksum())
	assert.Equal(t, expectedContent, string(content))
}

func TestDelayedMask(t *testing.T) {
	buffer, err := New()
	require.NoError(t, err)
	defer buffer.Close()

	n, err := buffer.Write([]byte("data before mask\n"))
	assert.NoError(t, err)
	assert.Greater(t, n, 0)

	buffer.SetMasked([]string{"mask_me"})

	n, err = buffer.Write([]byte("data mask_me masked\n"))
	assert.NoError(t, err)
	assert.Greater(t, n, 0)

	buffer.Finish()

	content, err := buffer.Bytes(0, 1000)
	require.NoError(t, err)

	expectedContent := "data before mask\ndata [MASKED] masked\n"
	assert.Equal(t, len(expectedContent), buffer.Size(), "unexpected buffer size")
	assert.Equal(t, "crc32:690f62e1", buffer.Checksum())
	assert.Equal(t, expectedContent, string(content))
}

func TestDelayedLimit(t *testing.T) {
	buffer, err := New()
	require.NoError(t, err)
	defer buffer.Close()

	n, err := buffer.Write([]byte("data before limit\n"))
	assert.NoError(t, err)
	assert.Greater(t, n, 0)

	buffer.SetLimit(20)

	n, err = buffer.Write([]byte("data after limit\n"))
	assert.NoError(t, err)
	assert.Greater(t, n, 0)

	buffer.Finish()

	content, err := buffer.Bytes(0, 1000)
	require.NoError(t, err)

	expectedContent := "data before limit\nda\n\x1b[33;1mJob's log exceeded limit of 20 bytes.\n" +
		"Job execution will continue but no more output will be collected.\x1b[0;m\n"
	assert.Equal(t, len(expectedContent), buffer.Size(), "unexpected buffer size")
	assert.Equal(t, "crc32:559aa46f", buffer.Checksum())
	assert.Equal(t, expectedContent, string(content))
}

func TestTraceRace(t *testing.T) {
	buffer, err := New()
	require.NoError(t, err)
	defer buffer.Close()

	buffer.SetLimit(1000)

	load := []func(){
		func() { _, _ = buffer.Write([]byte("x")) },
		func() { buffer.SetMasked([]string{"x"}) },
		func() { buffer.SetLimit(1000) },
		func() { buffer.Checksum() },
		func() { buffer.Size() },
		func() { _, _ = buffer.Bytes(0, 1000) },
	}

	var wg sync.WaitGroup
	for _, fn := range load {
		wg.Add(1)
		go func(fn func()) {
			defer wg.Done()

			for i := 0; i < 100; i++ {
				fn()
			}
		}(fn)
	}

	wg.Wait()

	buffer.Finish()

	_, err = buffer.Bytes(0, 1000)
	require.NoError(t, err)
}

func TestFixupInvalidUTF8(t *testing.T) {
	buffer, err := New()
	require.NoError(t, err)
	defer buffer.Close()

	buffer.SetMasked([]string{"hello", "\xfe"})

	// \xfe and \xff are both invalid
	// \xfe we're masking though, so will be replaced with [MASKED]
	// \xff will be replaced by the "unicode replacement character" \ufffd
	// this ensures that masking happens prior to the utf8 fix
	_, err = buffer.Write([]byte("hello a\xfeb a\xffb\n"))
	require.NoError(t, err)

	content, err := buffer.Bytes(0, 1000)
	require.NoError(t, err)

	assert.True(t, utf8.ValidString(string(content)))
	assert.Equal(t, "[MASKED] a[MASKED]b a\ufffdb\n", string(content))
}

const logLineStr = "hello world, this is a lengthy log line including secrets such as 'hello', and " +
	"https://example.com/?rss_token=foo&rss_token=bar and http://example.com/?authenticity_token=deadbeef and " +
	"https://example.com/?rss_token=foobar. it's longer than most log lines, but probably a good test for " +
	"anything that's benchmarking how fast it is to write log lines."

var logLineByte = []byte(logLineStr)

func benchmarkBuffer10k(b *testing.B, line []byte) {
	for i := 0; i < b.N; i++ {
		func() {
			buffer, err := New()
			require.NoError(b, err)
			defer buffer.Close()

			buffer.SetLimit(math.MaxInt64)
			buffer.SetMasked([]string{"hello"})

			const N = 10000
			b.ReportAllocs()
			b.SetBytes(int64(len(line) * N))
			for i := 0; i < N; i++ {
				_, _ = buffer.Write(line)
			}
			buffer.Finish()
		}()
	}
}

func BenchmarkBuffer10k(b *testing.B) {
	b.Run("10k-logline", func(b *testing.B) {
		benchmarkBuffer10k(b, logLineByte)
	})

	b.Run("10k-small", func(b *testing.B) {
		benchmarkBuffer10k(b, []byte("hello"))
	})
}
