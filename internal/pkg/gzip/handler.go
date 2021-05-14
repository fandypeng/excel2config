package gzip

import (
	"compress/gzip"
	"fmt"
	bm "github.com/go-kratos/kratos/pkg/net/http/blademaster"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
)

type gzipHandler struct {
	*Options
	gzPool sync.Pool
}

func newGzipHandler(level int, options ...Option) *gzipHandler {
	var gzPool sync.Pool
	gzPool.New = func() interface{} {
		gz, err := gzip.NewWriterLevel(ioutil.Discard, level)
		if err != nil {
			panic(err)
		}
		return gz
	}
	handler := &gzipHandler{
		Options: DefaultOptions,
		gzPool:  gzPool,
	}
	for _, setter := range options {
		setter(handler.Options)
	}
	return handler
}

func (g *gzipHandler) Handle(c *bm.Context) {
	if fn := g.DecompressFn; fn != nil && c.Request.Header.Get("Content-Encoding") == "gzip" {
		fn(c)
	}

	if !g.shouldCompress(c.Request) {
		return
	}

	gz := g.gzPool.Get().(*gzip.Writer)
	defer g.gzPool.Put(gz)
	defer gz.Reset(ioutil.Discard)
	gz.Reset(c.Writer)

	c.Writer.Header().Set("Content-Encoding", "gzip")
	c.Writer.Header().Set("Vary", "Accept-Encoding")
	c.Writer = &gzipWriter{c.Writer, gz, 0}
	defer func() {
		gz.Close()
		c.Writer.Header().Set("Content-Length", fmt.Sprint(c.Writer.(*gzipWriter).Size()))
	}()
	c.Next()
}

func (g *gzipHandler) shouldCompress(req *http.Request) bool {
	if !strings.Contains(req.Header.Get("Accept-Encoding"), "gzip") ||
		strings.Contains(req.Header.Get("Connection"), "Upgrade") ||
		strings.Contains(req.Header.Get("Content-Type"), "text/event-stream") {

		return false
	}

	extension := filepath.Ext(req.URL.Path)
	if g.ExcludedExtensions.Contains(extension) {
		return false
	}

	if g.ExcludedPaths.Contains(req.URL.Path) {
		return false
	}
	if g.ExcludedPathesRegexs.Contains(req.URL.Path) {
		return false
	}

	return true
}
