package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	fcert := "/go/src/app/cert.pem"
	fkey := "/go/src/app/key.pem"

	// Load our CA certificate
	clientCACert, err := ioutil.ReadFile(fcert)
	if err != nil {
		log.Fatal("Unable to open cert", err)
	}

	clientCertPool := x509.NewCertPool()
	clientCertPool.AppendCertsFromPEM(clientCACert)

	tlsConfig := &tls.Config{
		// Reject any TLS certificate that cannot be validated
		ClientAuth: tls.RequireAndVerifyClientCert,
		// Ensure that we only use our "CA" to validate certificates
		ClientCAs: clientCertPool,
		// PFS because we can but this will reject client with RSA certificates
		CipherSuites: []uint16{tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384, tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256},
		// Force it server side
		PreferServerCipherSuites: true,
		// TLS 1.2 because we can
		MinVersion: tls.VersionTLS12,
	}

	tlsConfig.BuildNameToCertificate()

	proxyURL, err := url.Parse("http://caddy")
	if err != nil {
		log.Fatalln("Unable to parse proxy URL")
	}
	proxy := httputil.NewSingleHostReverseProxy(proxyURL)
	httpServer := &http.Server{
		Handler:   proxy,
		Addr:      ":1025",
		TLSConfig: tlsConfig,
	}

	log.Fatalln(httpServer.ListenAndServeTLS(fcert, fkey))
}
