package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/url"
	"strings"
	"testing"
)

func newPacRunner(t *testing.T, expr string) *PacRunner {
	pr, err := NewPacRunner(strings.NewReader(fmt.Sprintf(
		"function FindProxyForURL(url, host) { return %s }", expr)))
	require.Nil(t, err)
	return pr
}

func checkProxy(t *testing.T, pr *PacRunner, addr, expected string) {
	u, err := url.Parse(addr)
	require.Nil(t, err)
	proxy, err := pr.FindProxyForURL(u)
	require.Nil(t, err)
	assert.Equal(t, expected, proxy)
}

func TestDirect(t *testing.T) {
	pr := newPacRunner(t, `"DIRECT"`)
	checkProxy(t, pr, "https://anz.com", "DIRECT")
}

func TestIsPlainHostName(t *testing.T) {
	pr := newPacRunner(t, `isPlainHostName(host) ? "y" : "n"`)
	checkProxy(t, pr, "https://www", "y")
	checkProxy(t, pr, "https://anz.com", "n")
}

func TestDnsDomainIs(t *testing.T) {
	pr := newPacRunner(t, `dnsDomainIs(host, ".anz.com") ? "y" : "n"`)
	checkProxy(t, pr, "https://www.anz.com", "y")
	checkProxy(t, pr, "https://www", "n")
}

func TestDnsDomainIsAnySuffix(t *testing.T) {
	// See https://bugs.chromium.org/p/chromium/issues/detail?id=299649.
	pr := newPacRunner(t, `dnsDomainIs(host, "anz.com") ? "y" : "n"`)
	checkProxy(t, pr, "https://notanz.com", "y")
}

func TestShExpMatch(t *testing.T) {
	pr := newPacRunner(t, `shExpMatch(url, "*/b/*") ? "y" : "n"`)
	checkProxy(t, pr, "http://anz.com/a/b/c.html", "y")
	checkProxy(t, pr, "http://anz.com/d/e/f.html", "n")
}
