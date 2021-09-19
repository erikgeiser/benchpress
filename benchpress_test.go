package benchpress

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/klauspost/compress/zstd"
	"github.com/ulikunitz/xz"
	"github.com/ulikunitz/xz/lzma"
)

const mb = 1024 * 1024

func BenchmarkGzip(b *testing.B) {
	gzipLevel := func(l int) string {
		switch l {
		case gzip.NoCompression:
			return "no compression"
		case gzip.BestSpeed:
			return "best speed"
		case gzip.BestCompression:
			return "best compression"
		case gzip.DefaultCompression:
			return "default"
		case gzip.HuffmanOnly:
			return "huffman only"
		default:
			return "invalid"
		}
	}

	levels := []int{
		gzip.BestSpeed,
		gzip.BestCompression,
		gzip.DefaultCompression,
	}

	var benchmarks []compressionBenchmark

	for _, level := range levels {
		level := level

		benchmarks = append(benchmarks, compressionBenchmark{
			name: fmt.Sprintf("gzip level %q", gzipLevel(level)),
			compressor: func(w io.Writer) (io.Writer, error) {
				return gzip.NewWriterLevel(w, level)
			},
			decompressor: func(r io.Reader) (io.Reader, error) {
				return gzip.NewReader(r)
			},
		})
	}

	runCompressionBenchmarks(b, benchmarks...)
}

func BenchmarkZstd(b *testing.B) {
	zstdLevel := func(l zstd.EncoderLevel) string {
		switch l {
		case zstd.SpeedFastest:
			return "fastest"
		case zstd.SpeedDefault:
			return "default"
		case zstd.SpeedBetterCompression:
			return "better compression"
		case zstd.SpeedBestCompression:
			return "best compression"
		default:
			return "invalid"
		}
	}

	levels := []zstd.EncoderLevel{
		zstd.SpeedFastest,
		zstd.SpeedDefault,
		zstd.SpeedBetterCompression,
		zstd.SpeedBestCompression,
	}

	var benchmarks []compressionBenchmark

	for _, level := range levels {
		level := level

		benchmarks = append(benchmarks, compressionBenchmark{
			name: fmt.Sprintf("zstd level %q", zstdLevel(level)),
			compressor: func(w io.Writer) (io.Writer, error) {
				return zstd.NewWriter(w, zstd.WithEncoderLevel(level))
			},
			decompressor: func(r io.Reader) (io.Reader, error) {
				return zstd.NewReader(r)
			},
		})
	}

	runCompressionBenchmarks(b, benchmarks...)
}

func BenchmarkXZ(b *testing.B) {
	runCompressionBenchmarks(b, []compressionBenchmark{{
		name: "xz",
		compressor: func(w io.Writer) (io.Writer, error) {
			return xz.NewWriter(w)
		},
		decompressor: func(r io.Reader) (io.Reader, error) {
			return xz.NewReader(r)
		},
	}}...)
}

func BenchmarkLZMA(b *testing.B) {
	runCompressionBenchmarks(b, []compressionBenchmark{{
		name: "lzma",
		compressor: func(w io.Writer) (io.Writer, error) {
			return lzma.NewWriter(w)
		},
		decompressor: func(r io.Reader) (io.Reader, error) {
			return lzma.NewReader(r)
		},
	}}...)
}

func BenchmarkLZMA2(b *testing.B) {
	runCompressionBenchmarks(b, []compressionBenchmark{{
		name: "lzma2",
		compressor: func(w io.Writer) (io.Writer, error) {
			return lzma.NewWriter2(w)
		},
		decompressor: func(r io.Reader) (io.Reader, error) {
			return lzma.NewReader2(r)
		},
	}}...)
}

type (
	compressor   func(io.Writer) (io.Writer, error)
	decompressor func(io.Reader) (io.Reader, error)
)

type compressionBenchmark struct {
	name         string
	compressor   compressor
	decompressor decompressor
}

func compress(b *testing.B, compressor compressor, w io.Writer, r io.Reader) {
	b.Helper()

	c, err := compressor(w)
	if err != nil {
		b.Fatalf("init: %v", err)
	}

	_, err = io.Copy(c, r)
	if err != nil {
		b.Fatalf("compress: %v", err)
	}

	if closer, ok := c.(io.Closer); ok {
		err = closer.Close()
		if err != nil {
			b.Fatalf("close: %v", err)
		}
	}
}

func runCompressionBenchmarks(b *testing.B, bench ...compressionBenchmark) {
	contents, err := os.ReadDir("testdata")
	if err != nil {
		b.Fatalf("listing testdata: %v", err)
	}

	if len(contents) == 0 {
		b.Fatalf("no testdata")
	}

	for _, entry := range contents {
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		data, err := os.ReadFile(filepath.Join("testdata", entry.Name()))
		if err != nil {
			b.Fatalf("reading test file: %v", err)
		}

		b.Logf("%s: initial size: %d MB", entry.Name(), len(data)/mb)

		for _, test := range bench {
			b.Run(entry.Name()+"_compress_"+test.name, func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					var compressed bytes.Buffer

					compress(b, test.compressor, &compressed, bytes.NewReader(data))
				}
			})

			var buf bytes.Buffer
			compress(b, test.compressor, &buf, bytes.NewReader(data))
			compressedData := buf.Bytes()

			b.Logf("%s: %s: %.2f MB (%.2f%%)", entry.Name(), test.name,
				float64(len(compressedData))/mb,
				100.*float64(len(compressedData))/float64(len(data)))

			b.Run(entry.Name()+"_decompress_"+test.name, func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					var dump bytes.Buffer

					d, err := test.decompressor(bytes.NewReader(compressedData))
					if err != nil {
						b.Fatalf("init: %v", err)
					}

					_, err = io.Copy(&dump, d)
					if err != nil {
						b.Fatalf("decompress: %v", err)
					}
				}
			})
		}
	}
}
