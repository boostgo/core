package httpx

import (
	"context"
	"io"
	"net/http"
	"time"
)

type CacheSetter interface {
	Set(ctx context.Context, request *http.Request, responseBody []byte, ttl time.Duration) error
}

type CacheGetter interface {
	Get(ctx context.Context, request *http.Request) ([]byte, bool, error)
}

type CacheDistributor interface {
	CacheSetter
	CacheGetter
}

type cacheResponseWriter struct {
	inner  http.ResponseWriter
	writer io.Writer
}

func NewCacheResponseWriter(inner http.ResponseWriter, writer io.Writer) http.ResponseWriter {
	return &cacheResponseWriter{
		inner:  inner,
		writer: writer,
	}
}

func (writer *cacheResponseWriter) Header() http.Header {
	return writer.inner.Header()
}

func (writer *cacheResponseWriter) Write(b []byte) (int, error) {
	return writer.writer.Write(b)
}

func (writer *cacheResponseWriter) WriteHeader(statusCode int) {
	writer.inner.WriteHeader(statusCode)
}
