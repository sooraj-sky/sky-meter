package skydns

import (
	"context"
	"net"
	"time"
)

// Create a custom Resolver that uses a specific DNS server IP
func customResolver(dnsServer string) *net.Resolver {
	return &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			dialer := net.Dialer{
				Timeout: 10 * time.Second,
				Resolver: &net.Resolver{
					PreferGo: true,
					Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
						return (&net.Dialer{
							Timeout:   10 * time.Second,
							DualStack: true,
						}).DialContext(ctx, network, dnsServer+":53")
					},
				},
			}
			return dialer.DialContext(ctx, network, address)
		},
	}
}
