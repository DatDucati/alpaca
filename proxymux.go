package main

import (
	"net/http"
)

type proxyMux struct {
	http.ServeMux
	proxy http.Handler
}

func (m *proxyMux) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Pass CONNECT requests and absolute-form URIs to the proxy handler.
	// If the request URL has a scheme, it is an absolute-form URI
	// (RFC 7230 Section 5.3.2).
	if req.Method == http.MethodConnect || req.URL.Scheme != "" {
		m.proxy.ServeHTTP(w, req)
		return
	}
	// The request URI is an origin-form or asterisk-form target which we
	// handle as an origin server (RFC 7230 5.3). authority-form URIs
	// are only for CONNECT, which has already been dispatched to the
	// proxy handler.
	m.ServeMux.ServeHTTP(w, req)
}
