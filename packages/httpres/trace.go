package httpreponser

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptrace"
	"os"
	"time"

	models "sky-meter/models"
	skydns "sky-meter/packages/dns"
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

func client(timeout time.Duration, SkipSsl bool) *http.Client {
	return &http.Client{
		Transport: transport(SkipSsl),
		Timeout:   time.Duration(timeout * time.Millisecond),
	}
}

func transport(SkipSsl bool) *http.Transport {
	return &http.Transport{
		DisableKeepAlives: SkipSsl,
		TLSClientConfig:   tlsConfig(),
	}
}

func tlsConfig() *tls.Config {
	return &tls.Config{
		InsecureSkipVerify: true,
	}
}

func trace() (*httptrace.ClientTrace, *models.Debug) {
	//DNS settings
	dnsServer := os.Getenv("dnsserver") // Replace with your desired DNS server IP address
	resolver := skydns.CustomResolver(dnsServer)
	fmt.Sprintln(resolver)
	d := &models.Debug{}

	t := &httptrace.ClientTrace{

		DNSStart: func(info httptrace.DNSStartInfo) {
			t := time.Now().UTC().String()
			//log.Println(t, "dns start")
			d.DNS.Start = t
			d.DNS.Host = info.Host
		},
		DNSDone: func(info httptrace.DNSDoneInfo) {
			t := time.Now().UTC().String()
			//log.Println(t, "dns end")
			d.DNS.End = t
			d.DNS.Address = info.Addrs
			d.DNS.Error = info.Err
		},
		ConnectStart: func(network, addr string) {
			t := time.Now().UTC().String()
			//log.Println(t, "dial start")
			d.Dial.Start = t
		},
		ConnectDone: func(network, addr string, err error) {
			t := time.Now().UTC().String()
			//log.Println(t, "dial end")
			d.Dial.End = t
		},
		GotConn: func(connInfo httptrace.GotConnInfo) {
			t := time.Now().UTC().String()
			//log.Println(t, "conn time")
			d.Connection.Time = t
		},
		WroteHeaders: func() {
			t := time.Now().UTC().String()
			//log.Println(t, "wrote all request headers")
			d.WroteAllRequestHeaders.Time = t
		},
		WroteRequest: func(wr httptrace.WroteRequestInfo) {
			t := time.Now().UTC().String()
			//log.Println(t, "wrote all request")
			d.WroteAllRequest.Time = t
		},
		GotFirstResponseByte: func() {
			t := time.Now().UTC().String()
			//log.Println(t, "first received response byte")
			d.FirstReceivedResponseByte.Time = t
		},
	}

	return t, d
}

func CallEndpoint(endpoint interface{}, timeout int, SkipSsl bool) ([]byte, int, error) {
	url, _ := endpoint.(string)
	httpresdata, statusCode, err := GetHttpdata(url, time.Duration(timeout), SkipSsl)
	return httpresdata, statusCode, err
}
