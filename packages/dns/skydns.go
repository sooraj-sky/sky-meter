package skydns

import (
	"context"
	"net"
	"time"
)

// CustomResolver returns a custom resolver that resolves DNS queries using the specified DNS server.
func CustomResolver(dnsServer string) *net.Resolver {
	return &net.Resolver{

		// create a new net.Resolver with PreferGo set to true
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {

			// create a new net.Dialer with a timeout of 10 seconds and a custom resolver that uses the specified DNS server
			dialer := net.Dialer{
				Timeout: 10 * time.Second,
				Resolver: &net.Resolver{
					PreferGo: true,
					Dial: func(ctx context.Context, network, address string) (net.Conn, error) {

						// create a new net.Dialer with a timeout of 10 seconds and DualStack set to true, and use it to dial the specified DNS server
						return (&net.Dialer{
							Timeout:   10 * time.Second,
							DualStack: true,
						}).DialContext(ctx, network, dnsServer+":53")
					},
				},
			}
			// use the dialer to dial the specified network and address
			return dialer.DialContext(ctx, network, address)
		},
	}
}
