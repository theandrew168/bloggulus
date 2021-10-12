package cache

import (
	"bytes"
	"log"
	"net/http"
	"sync"
	"time"
)

type response struct {
	statusCode int
	header     http.Header
	body       bytes.Buffer
}

func (r *response) Header() http.Header {
	return r.header
}

func (r *response) Write(body []byte) (int, error) {
	return r.body.Write(body)
}

func (r *response) WriteHeader(statusCode int) {
	r.statusCode = statusCode
}

func copyHeaders(dst, src http.ResponseWriter) {
	for key, _ := range src.Header() {
		for _, value := range src.Header().Values(key) {
			dst.Header().Add(key, value)
		}
	}
}

type memoryCache struct {
	sync.RWMutex

	data map[string]response
	ttl  time.Duration
	next http.Handler
}

func (c *memoryCache) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url := r.URL.String()

	// check if URL is cached
	c.RLock()
	if resp, ok := c.data[url]; ok {
		log.Printf("cache hit: %s\n", url)

		// send response from cache
		copyHeaders(w, &resp)
		w.WriteHeader(resp.statusCode)
		w.Write(resp.body.Bytes())
		c.RUnlock()
		return
	}
	c.RUnlock()

	log.Printf("cache miss: %s\n", url)

	// call and capture the next handler
	resp := response{
		statusCode: 200,
		header:     make(http.Header),
		body:       bytes.Buffer{},
	}
	c.next.ServeHTTP(&resp, r)

	// check for errors on downstream handler
	if resp.statusCode >= 400 {
		copyHeaders(w, &resp)
		w.WriteHeader(resp.statusCode)
		w.Write(resp.body.Bytes())
		return
	}

	// update cache
	c.Lock()
	c.data[url] = resp
	c.Unlock()

	// send response
	copyHeaders(w, &resp)
	w.WriteHeader(resp.statusCode)
	w.Write(resp.body.Bytes())
}

func NewMemory(ttl time.Duration) func(http.Handler) http.Handler {
	fn := func(next http.Handler) http.Handler {
		return &memoryCache{
			data: make(map[string]response),
			ttl:  ttl,
			next: next,
		}
	}
	return fn
}
