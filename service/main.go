package service

import (
	"crypto/tls"
	"pathwae/api"
	"pathwae/proxy"
)

func Start(conf string) {
	log.Println("Starting services...")
	// create tls config for the tlsServer
	tlsServer := proxy.NewReversProxy(conf, ":https") // :443
	config := &tls.Config{
		GetCertificate: tlsServer.GetCerts,
	}
	tlsServer.TLSConfig = config
	log.Println("Starting https reverse proxy...")
	go tlsServer.ListenAndServeTLS("", "")

	// create a http server
	httpServer := proxy.NewReversProxy(conf, ":http") // :80
	log.Println("Starting http reverse proxy...")
	go httpServer.ListenAndServe()

	// start the API
	log.Println("Starting API on port :8080...")
	log.Fatal(api.Start(":8080"))
}
