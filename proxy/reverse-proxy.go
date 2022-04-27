package proxy

import (
	"crypto/tls"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

const certDir = "/certs"

type ProxyError struct {
	URL string
	Err string
	At  time.Time
}

// ReverseProxy is an HTTP Handler that takes an incoming request and sends it to another server, proxying the response back to the client.
type ReverseProxy struct {
	// http.Server Composing.
	*http.Server

	// list of configured backends to proxy to.
	backends map[string]*Backend
}

// NewReversProxy returns a new ReverseProxy to handle HTTP requests.
func NewReversProxy(conf, addr string) *ReverseProxy {
	s := &ReverseProxy{
		Server: &http.Server{
			Addr: addr,
		},
		backends: LoadYAMLConfig(conf),
	}

	s.Server.Handler = s // force to use the ServeHTTP method
	return s
}

// ServeHTTP implements the http.Handler interface. This method is called by the http.Server for each request. Whatever it is made
// in HTTP or HTTPS. It will proxy the request to the right server with httputil.ReverseProxy.
func (rp *ReverseProxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// Get the target host and make a new request to it.
	target, ok := rp.backends[req.Host]
	if !ok {
		http.Error(rw, "no such host", http.StatusBadRequest)
		return
	}

	// if backend is disabled, return a 503
	if !*target.Enabled {
		http.Error(rw, "backend is disabled", http.StatusServiceUnavailable)
		return
	}

	// if req.TLS is nil and server.ForceSSL, redirect
	if target.ForceSSL && req.TLS == nil {
		http.Redirect(
			rw, req,
			"https://"+req.Host+req.URL.String(),
			http.StatusMovedPermanently,
		)
		return
	}

	// fix the url if no scheme is provided
	if !strings.HasPrefix(target.To, "http") {
		target.To = "http://" + target.To
	}

	// we must proxy the "targe" host to "to" host
	to, err := url.Parse(target.To)
	if err != nil {
		http.Error(rw, "url parse: "+err.Error(), http.StatusBadRequest)
		return
	}

	// make seom stats
	go func(r *http.Request) {
		// send the content to stats, not the pointer
		StatChan <- *r
	}(req)

	// create a ReverseProxy
	proxy := &httputil.ReverseProxy{
		Director: func(proxied *http.Request) {
			proxied.URL.Scheme = to.Scheme
			proxied.URL.Host = to.Host
			host := req.Host
			if strings.Contains(host, ":") {
				host = strings.Split(host, ":")[0]
			}
			proxied.Header.Set("Host", host)
			proxied.Header.Set("X-Forwarded-For", req.RemoteAddr)
			proxied.Header.Set("X-Forwarded-Host", req.Host)
			if req.TLS != nil {
				proxied.Header.Set("X-Forwarded-Proto", "https")
			} else {
				proxied.Header.Set("X-Forwarded-Proto", "http")
			}
		},
	}
	proxy.ServeHTTP(rw, req)
}

// GetCerts is a callback function for tls.Config.GetCertificate - it search the right certificate to use.
// It will manage widlcard certificates. If the certificate is not found, so a temporary one will be created.
func (rp *ReverseProxy) GetCerts(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
	// standard name
	if cert, ok := certCache[info.ServerName]; ok {
		return cert, nil
	}

	// find if there is a certificate for the host with wildcard
	for name, cert := range certCache {
		if strings.HasPrefix(name, "*") {
			domain := strings.TrimPrefix(name, "*")
			if strings.HasSuffix(info.ServerName, domain) {
				cert.Leaf.DNSNames = append(cert.Leaf.DNSNames, info.ServerName)
				certCache[info.ServerName] = cert // cache
				return cert, nil
			}
		}
	}

	// no provided certificate... OK, let's create one
	log.Println("No certificate found for", info.ServerName, "creating one temporary...")
	cert := MakeCert(info.ServerName)
	SetCertCache(info.ServerName, cert)
	return cert, nil
}
