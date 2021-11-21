package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProxyMux(t *testing.T) {
	var mux proxyMux
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("I'm a web server"))
	}))
	mux.proxy = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("I'm a proxy server"))
	})

	server := httptest.NewServer(&mux)
	defer server.Close()

	for _, test := range []struct{
		name string
		method string
		target string
		response string
	}{
		{"OriginForm", http.MethodGet, "/path", "I'm a web server"},
		{"AbsoluteForm", http.MethodGet, "http://www.test/path", "I'm a proxy server"},
		{"AuthorityForm", http.MethodConnect, "www.test:443", "I'm a proxy server"},
	} {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(test.method, test.target, nil)
			mux.ServeHTTP(w, req)
			resp := w.Result()
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			assert.Equal(t, test.response, string(body))
		})
	}
}
