package proxy

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	temporaryCertStore = "/tmpcerts"
	Organization       = "Pathwae"
	keySize            = 2048
)

var (
	certLock = &sync.Mutex{}
)

// SetCertCache add a certificate to the cache. It also notifies all listeners for the given host.
func SetCertCache(servername string, cert *tls.Certificate) {
	certLock.Lock()
	log.Println("Setting certificate for " + servername)
	certCache[servername] = cert
	certLock.Unlock()
	go func() {
		for _, listener := range certListeners[servername] {
			listener <- cert
		}
	}()
}

// MakeCert generates a self-signed certificate for the given host.
func MakeCert(host string) *tls.Certificate {

	// make the temporaryCertDir if it doesn't exist
	tempCertDir := temporaryCertStore
	if os.Getenv("TESTMODE") == "1" {
		tempCertDir = "./tmpcerts"
	}

	if _, err := os.Stat(tempCertDir); os.IsNotExist(err) {
		os.Mkdir(tempCertDir, 0755)
	}

	log.Println("Generating temporary certificate for " + host)

	// generate a key pair
	priv, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		log.Fatalf("failed to generate private key: %s", err)
	}

	// start building the template
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Issuer: pkix.Name{
			Organization: []string{Organization},
		},
		Subject: pkix.Name{
			CommonName:   host,
			Organization: []string{Organization},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(4 * 365 * 24 * time.Hour),
		// for given host
		DNSNames: []string{host},
	}

	// sign the template with the private key
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		log.Fatalf("Failed to create certificate: %s", err)
	}

	// names of certs:
	certfilename := filepath.Join(tempCertDir, host+".crt")
	keyfilename := filepath.Join(tempCertDir, host+".key")

	// write cert and key files
	certOut, err := os.Create(certfilename)
	if err != nil {
		log.Fatalf("failed to open "+certfilename+" for writing: %s", err)
	}
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()

	keyOut, err := os.OpenFile(keyfilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("failed to open "+keyfilename+" for writing: %s", err)
	}

	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
	keyOut.Close()

	// return the cert
	cert, err := tls.LoadX509KeyPair(certfilename, keyfilename)
	if err != nil {
		log.Fatalf("failed to load certificate: %s", err)
	}
	return &cert

}
