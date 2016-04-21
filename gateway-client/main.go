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

var proxy *httputil.ReverseProxy

func main() {
	// Load our TLS key pair to use for authentication
	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		log.Fatalln("Unable to load cert", err)
	}

	// Load our CA certificate
	clientCACert, err := ioutil.ReadFile("cert.pem")
	if err != nil {
		log.Fatal("Unable to open cert", err)
	}

	clientCertPool := x509.NewCertPool()
	clientCertPool.AppendCertsFromPEM(clientCACert)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      clientCertPool,
	}

	tlsConfig.BuildNameToCertificate()

	proxyURL, err := url.Parse("https://home.strothmann.xyz:1025")
	if err != nil {
		log.Fatalln("Unable to parse proxy URL")
	}
	proxy = httputil.NewSingleHostReverseProxy(proxyURL)
	proxy.Transport = &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	log.Fatalln(http.ListenAndServe(":8080", proxy))
}

// func proxyHandler(w http.ResponseWriter, r *http.Request) {
//   r.
// 	resp, err := client.Do(r)
// 	if err != nil {
// 		log.Println("error", err)
// 		w.WriteHeader(http.StatusBadGateway)
// 		w.Write([]byte(err.Error()))
// 		return
// 	}
// 	b, err := ioutil.ReadAll(resp.Body)
// 	defer resp.Body.Close()
// 	_, err = w.Write(b)
// 	if err != nil {
// 		log.Println(err)
// 	}
// }
