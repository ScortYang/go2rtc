package xiaomi

import (
	"net/http"
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
