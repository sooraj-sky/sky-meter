package httpreponser

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptrace"
	"time"

	models "sky-meter/models"
	skydns "sky-meter/packages/dns"
	skyenv "sky-meter/packages/env"
)

func GetHttpdata(url string, timeout time.Duration, SkipSsl bool) (httpdata []byte, httpstatuscode int, errs error) {
	// Create trace struct.
	trace, debug := trace()

	// Prepare request with trace attached to it.
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatalln("request error", err)
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))

	// MAke a request.
	res, err := client(timeout, SkipSsl).Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer res.Body.Close()

	data, err := json.MarshalIndent(debug, "", "    ")
	return data, res.StatusCode, err
}

// client returns an instance of *http.Client with the given timeout and SSL skip settings.
func client(timeout time.Duration, SkipSsl bool) *http.Client {
	// The transport function is called to create the Transport object for the client.
	return &http.Client{
		Transport: transport(SkipSsl),                        // If SkipSsl is true, the Transport object is created with InsecureSkipVerify set to true.
		Timeout:   time.Duration(timeout * time.Millisecond), // The Timeout field of the client is set to the given timeout value.
	}
	// The returned client should be used to make HTTP requests to remote servers.
}

// transport creates a new http.Transport instance with the provided SkipSsl flag
// If SkipSsl is true, DisableKeepAlives will be set to true to disable keep-alive connections
// TLSClientConfig will be set using the tlsConfig() function
func transport(SkipSsl bool) *http.Transport {
	return &http.Transport{
		DisableKeepAlives: SkipSsl,
		TLSClientConfig:   tlsConfig(SkipSsl),
	}
}

func tlsConfig(SkipSsl bool) *tls.Config {
	// Create a new tls.Config struct and initialize its fields
	return &tls.Config{
		// SSL verification
		InsecureSkipVerify: SkipSsl,
	}
}

func trace() (*httptrace.ClientTrace, *models.Debug) {

	// Set up custom DNS resolver using DNS server from environment variables.
	allEnv := skyenv.GetEnv()
	dnsServer := allEnv.DnsServer
	skydns.CustomResolver(dnsServer)


	// Create a new Debug object.
	d := &models.Debug{}

	// Create a new ClientTrace object with callback functions for different stages of the HTTP request/response cycle.
	t := &httptrace.ClientTrace{

		// Callback function for DNS start.
		DNSStart: func(info httptrace.DNSStartInfo) {
			t := time.Now().UTC().String()
			d.DNS.Start = t
			d.DNS.Host = info.Host
		},
		// Callback function for DNS end.
		DNSDone: func(info httptrace.DNSDoneInfo) {
			t := time.Now().UTC().String()
			d.DNS.End = t
			d.DNS.Address = info.Addrs
			d.DNS.Error = info.Err
		},
		// Callback function for dial start.
		ConnectStart: func(network, addr string) {
			t := time.Now().UTC().String()
			d.Dial.Start = t
		},
		// Callback function for dial end.
		ConnectDone: func(network, addr string, err error) {
			t := time.Now().UTC().String()
			d.Dial.End = t
		},
		// Callback function for connection time.
		GotConn: func(connInfo httptrace.GotConnInfo) {
			t := time.Now().UTC().String()
			d.Connection.Time = t
		},
		// Callback function for writing all request headers.
		WroteHeaders: func() {
			t := time.Now().UTC().String()
			d.WroteAllRequestHeaders.Time = t
		},
		// Callback function for writing all request.
		WroteRequest: func(wr httptrace.WroteRequestInfo) {
			t := time.Now().UTC().String()
			d.WroteAllRequest.Time = t
		},
		// Callback function for first received response byte.
		GotFirstResponseByte: func() {
			t := time.Now().UTC().String()
			d.FirstReceivedResponseByte.Time = t
		},
	}

	// Return the ClientTrace object and the Debug object.
	return t, d
}

// CallEndpoint calls an HTTP endpoint with the given URL and timeout, and returns the response body, status code, and any errors encountered.
func CallEndpoint(endpoint interface{}, timeout int, SkipSsl bool) ([]byte, int, error) {

	// Convert the endpoint parameter to a string
	url, _ := endpoint.(string)

	// Call GetHttpdata to retrieve the HTTP response data, status code, and error
	httpresdata, statusCode, err := GetHttpdata(url, time.Duration(timeout), SkipSsl)

	// Return the HTTP response data, status code, and error
	return httpresdata, statusCode, err
}
