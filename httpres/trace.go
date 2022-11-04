package httpreponser

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptrace"
	models "sky-meter/models"
	"time"
)

func GetHttpdata(url string) (httpdata []byte, httpsatauscode int) {
	// Create trace struct.
	trace, debug := trace()

	// Prepare request with trace attached to it.
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatalln("request error", err)
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))

	// MAke a request.
	res, err := client().Do(req)
	if err != nil {
		log.Fatalln("client error", err)
	}
	defer res.Body.Close()

	data, err := json.MarshalIndent(debug, "", "    ")
	return data, res.StatusCode
}

func client() *http.Client {
	return &http.Client{
		Transport: transport(),
	}
}

func transport() *http.Transport {
	return &http.Transport{
		DisableKeepAlives: true,
		TLSClientConfig:   tlsConfig(),
	}
}

func tlsConfig() *tls.Config {
	return &tls.Config{
		InsecureSkipVerify: true,
	}
}

func trace() (*httptrace.ClientTrace, *models.Debug) {
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
