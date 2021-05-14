package gzip

import (
	"compress/gzip"
	"net/http"

	bm "github.com/go-kratos/kratos/pkg/net/http/blademaster"
)

const (
	BestCompression    = gzip.BestCompression
	BestSpeed          = gzip.BestSpeed
	DefaultCompression = gzip.DefaultCompression
	NoCompression      = gzip.NoCompression
)

func Gzip(level int, options ...Option) bm.HandlerFunc {
	return newGzipHandler(level, options...).Handle
}

type gzipWriter struct {
	http.ResponseWriter
	writer *gzip.Writer
	size   int
}

func (g *gzipWriter) WriteString(s string) (int, error) {
	g.Header().Del("Content-Length")
	n, err := g.writer.Write([]byte(s))
	g.size += n
	return n, err
}

func (g *gzipWriter) Write(data []byte) (int, error) {
	g.Header().Del("Content-Length")
	n, err := g.writer.Write(data)
	g.size += n
	return n, err
}

func (g *gzipWriter) WriteHeader(code int) {
	g.Header().Del("Content-Length")
	g.size = 0
	g.ResponseWriter.WriteHeader(code)
}

func (g *gzipWriter) Size() int {
	return g.size
}
