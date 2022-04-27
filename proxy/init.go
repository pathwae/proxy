package proxy

import (
	"crypto/tls"
	"crypto/x509"
	"os"
	"path/filepath"
)

var certCache = make(map[string]*tls.Certificate, 0)
var errorLogs = make(map[string][]*ProxyError, 0) // unuesd at this time

func init() {

	certs, err := filepath.Glob(certDir + string(filepath.Separator) + "*")
	if os.Getenv("TESTMODE") == "1" {
		certs, err = filepath.Glob("./certs/*")
	}

	if err != nil {
		log.Fatal(err)
	}
	// classify cert and keys
	for _, cert := range certs {
		for _, key := range certs {
			if key == cert {
				continue
			}
			// try to load cert with key
			tlsCert, err := tls.LoadX509KeyPair(cert, key)
			if err != nil {
				continue
			}

			x509Cert, err := x509.ParseCertificate(tlsCert.Certificate[0])
			if err != nil {
				continue
			}
			tlsCert.Leaf = x509Cert

			// log the cert info
			log.Println("Load cert with key success:", cert, key)
			log.Println("Cert info:", x509Cert.Subject.CommonName, x509Cert.DNSNames, x509Cert.Subject)

			// if cert is decoded, record !
			for _, name := range x509Cert.DNSNames {
				SetCertCache(name, &tlsCert)
			}
		}
	}
}
