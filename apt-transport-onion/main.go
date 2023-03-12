package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"time"

	"github.com/eyedeekay/apttransport"
	"github.com/eyedeekay/onramp"
)

func main() {
	transport := &apttransport.AptMethod{}
	transport.Main = transport.DefaultMain
	onion, err := onramp.NewOnion("apt-transport-onion")
	if err == nil {
		samHTTPClient := GenerateSAMHTTPClient(onion)
		transport.Client = samHTTPClient
		transport.Main()
	} else {
		log.Fatal(err)
	}

}

func GenerateSAMHTTPClient(onion *onramp.Onion) apttransport.AptClient {
	aptClient := &http.Client{
		Timeout: time.Duration(6) * time.Minute,
		Transport: &http.Transport{
			MaxIdleConns:          0,
			MaxIdleConnsPerHost:   2,
			DisableKeepAlives:     false,
			ResponseHeaderTimeout: time.Duration(10) * time.Minute,
			ExpectContinueTimeout: time.Duration(10) * time.Minute,
			IdleConnTimeout:       time.Duration(10) * time.Minute,
			TLSNextProto:          make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
			Dial:                  onion.Dial,
		},
		Jar:           nil,
		CheckRedirect: nil,
	}
	return aptClient
}
