package proxy

import (
	"crypto/tls"
	"net"
	"net/url"
	"strings"
	"time"
)

// GetBackends returns a list of all servers.
func GetBackends() map[string]*Backend {
	return backendList
}

// GetBackend returns the server with the given name.
func GetBackend(backendName string) *Backend {
	if s, ok := backendList[backendName]; ok {
		return s
	}
	return nil
}

// GetCertForBackend returns the certificate for the given server.
func GetCertForBackend(backendName string) *tls.Certificate {
	return certCache[backendName]
}

// IsBackendUp returns true if the server is up and running.
func IsBackendUp(backend Backend) bool {

	if backend.Enabled == nil || !*backend.Enabled {
		return false
	}

	if !strings.HasPrefix(backend.To, "http") {
		backend.To = "http://" + backend.To
	}
	//parse the url
	parsed, err := url.Parse(backend.To)
	if err != nil {
		return false
	}
	//get the hostname
	hostname := parsed.Hostname()
	//get the port
	port := parsed.Port()
	if port == "" {
		if parsed.Scheme == "https" {
			port = "443"
		} else {
			port = "80"
		}
	}

	// check connection
	conn, err := net.DialTimeout("tcp", hostname+":"+port, time.Second*5)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// GetErrors returns the errors for the given server.
// Disabled at this time...
func GetErrors(backendName string) []*ProxyError {
	return errorLogs[backendName]
}
