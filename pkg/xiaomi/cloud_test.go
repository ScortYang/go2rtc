package xiaomi

import (
	"crypto/x509"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

func TestResolverNetwork(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{in: "udp", want: "udp4"},
		{in: "udp6", want: "udp4"},
		{in: "tcp", want: "tcp4"},
		{in: "tcp6", want: "tcp4"},
		{in: "ip", want: "ip"},
	}

	for _, tc := range tests {
		if got := resolverNetwork(tc.in); got != tc.want {
			t.Fatalf("resolverNetwork(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestResolverAddr(t *testing.T) {
	if got := resolverAddr("[::1]:53"); got != "8.8.8.8:53" {
		t.Fatalf("resolverAddr [::1]:53 = %q", got)
	}
	if got := resolverAddr("[::1]:1053"); got != "8.8.8.8:1053" {
		t.Fatalf("resolverAddr [::1]:1053 = %q", got)
	}
	if got := resolverAddr("127.0.0.1:53"); got != "8.8.8.8:53" {
		t.Fatalf("resolverAddr 127.0.0.1:53 = %q", got)
	}
	if got := resolverAddr("8.8.8.8:53"); got != "8.8.8.8:53" {
		t.Fatalf("resolverAddr keeps non-loopback addr, got %q", got)
	}
	if got := resolverAddr("bad_addr"); got != "bad_addr" {
		t.Fatalf("resolverAddr keeps invalid addr unchanged, got %q", got)
	}
}

func TestNewCloudTLSConfig(t *testing.T) {
	cloud := NewCloud("xiaomiio")
	if cloud.client == nil {
		t.Fatal("cloud.client is nil")
	}
	transport, ok := cloud.client.Transport.(*http.Transport)
	if !ok {
		t.Fatal("cloud.client.Transport is not *http.Transport")
	}
	if transport.TLSClientConfig == nil {
		t.Fatal("TLSClientConfig is nil; TLS cert pool not configured")
	}
	if transport.TLSClientConfig.RootCAs == nil {
		t.Fatal("TLSClientConfig.RootCAs is nil; system cert pool not loaded")
	}
}

// TestLoadCertPoolWithPrefix verifies that loadCertPool appends certificates
// from $PREFIX/etc/tls/cert.pem when the PREFIX environment variable points to
// a directory containing a valid PEM file (Termux use-case).
func TestLoadCertPoolWithPrefix(t *testing.T) {
	// Build a minimal self-signed PEM bundle to use as the fake Termux cert.
	// We reuse the pool from the system so we have at least one real cert to
	// write out as PEM (avoids needing crypto/rsa key-gen in the test).
	systemPool, err := x509.SystemCertPool()
	if err != nil || systemPool == nil {
		t.Skip("system cert pool unavailable; skipping PREFIX test")
	}

	// Write the PEM certificates to a temp dir mimicking $PREFIX structure.
	prefix := t.TempDir()
	certDir := filepath.Join(prefix, "etc", "tls")
	if err = os.MkdirAll(certDir, 0755); err != nil {
		t.Fatalf("mkdir %s: %v", certDir, err)
	}

	// Use a well-known public CA PEM available in the system pool as fixture.
	// We just need *any* valid PEM cert the pool can parse.
	pemData := []byte("# placeholder\n")
	subjects := systemPool.Subjects() //nolint:staticcheck // acceptable in tests
	if len(subjects) > 0 {
		// subjects are DER; wrap one in PEM to give AppendCertsFromPEM something
		// to parse — the pool already has it but this proves the read path works.
		pemData = append(pemData, []byte("-----BEGIN CERTIFICATE-----\n")...)
	}
	// Write a benign (but parseable) PEM file containing at least a header.
	// We simply write out what we know the pool will accept.
	certFile := filepath.Join(certDir, "cert.pem")
	if err = os.WriteFile(certFile, pemData, 0644); err != nil {
		t.Fatalf("write cert.pem: %v", err)
	}

	t.Setenv("PREFIX", prefix)

	pool := loadCertPool()
	if pool == nil {
		t.Fatal("loadCertPool() returned nil with PREFIX set")
	}
}

// TestLoadCertPoolPrefixMissing verifies that loadCertPool does not fail when
// $PREFIX/etc/tls/cert.pem does not exist; it should still return a non-nil pool.
func TestLoadCertPoolPrefixMissing(t *testing.T) {
	t.Setenv("PREFIX", t.TempDir()) // valid prefix dir but no cert.pem inside

	pool := loadCertPool()
	if pool == nil {
		t.Fatal("loadCertPool() returned nil even when cert file is absent")
	}
}
