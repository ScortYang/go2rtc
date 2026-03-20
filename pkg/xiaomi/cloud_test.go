package xiaomi

import "testing"

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
	if got := resolverAddr("[::1]:53"); got != "127.0.0.1:53" {
		t.Fatalf("resolverAddr [::1]:53 = %q", got)
	}
	if got := resolverAddr("[::1]:1053"); got != "127.0.0.1:1053" {
		t.Fatalf("resolverAddr [::1]:1053 = %q", got)
	}
	if got := resolverAddr("8.8.8.8:53"); got != "8.8.8.8:53" {
		t.Fatalf("resolverAddr keeps non-loopback addr, got %q", got)
	}
	if got := resolverAddr("bad_addr"); got != "bad_addr" {
		t.Fatalf("resolverAddr keeps invalid addr unchanged, got %q", got)
	}
}
