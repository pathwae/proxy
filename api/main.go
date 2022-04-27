package api

import (
	"crypto/x509"
	"encoding/json"
	"net/http"
	"pathwae/proxy"
)

type CertInfo struct {
	Issuer     string   `json:"issuer"`
	Subject    string   `json:"subject"`
	NotBefore  string   `json:"notBefore"`
	NotAfter   string   `json:"notAfter"`
	CommonName string   `json:"commonName"`
	DNSNames   []string `json:"dnsNames"`
}

// allow CORS
func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Length")
}

func GetBackends(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	// get the server list from poxy and return a json response
	servers := proxy.GetBackends()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(servers)
}

func GetCerts(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Panicln(err)
			w.Write([]byte(err.(error).Error()))
		}
	}()
	enableCORS(w)

	// get the server name given in the path
	serverName := r.URL.Path[len("/api/v1/cert/"):]

	if serverName == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("No server name given"))
		return
	}

	certs := proxy.GetCertForBackend(serverName)
	if certs == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No certs found for server " + serverName))
		return
	}

	// decode the certs
	info, err := x509.ParseCertificate(certs.Certificate[0])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	certs.Leaf = info

	// generate a certinfo struct
	certInfo := CertInfo{
		Issuer:     certs.Leaf.Issuer.String(),
		Subject:    certs.Leaf.Subject.String(),
		NotBefore:  certs.Leaf.NotBefore.String(),
		NotAfter:   certs.Leaf.NotAfter.String(),
		CommonName: certs.Leaf.Subject.CommonName,
		DNSNames:   certs.Leaf.DNSNames,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(certInfo)
}

func GetBackendState(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	// get the server name given in the path
	serverName := r.URL.Path[len("/api/v1/state/"):]
	server := proxy.GetBackend(serverName)

	if server == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No server found with name " + serverName))
		return
	}

	state := proxy.IsBackendUp(*server)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(state)
}

func GetBackendStats(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	// get the server name given in the path
	serverName := r.URL.Path[len("/api/v1/stats/"):]
	server := proxy.GetBackend(serverName)

	if server == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No server found with name " + serverName))
		return
	}

	stats := proxy.GetStat(serverName)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func GetRuntimeMemAlloc(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(proxy.GetMemory())
}

func GetVersion(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(proxy.GetVersion())
}

func SetBackend(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method == "OPTIONS" {
		return
	}
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// name is given in the path
	name := r.URL.Path[len("/api/v1/backend/"):]
	log.Println("Received request to set backend: " + name)
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("No server name given"))
		return
	}

	// Get post data as JSON, this should be a proxy.Backend
	var backend proxy.Backend
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&backend)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	err = proxy.SetBackend(name, backend)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}

func Start(port string) error {
	mux := http.NewServeMux()
	// add the CORS middleware
	mux.HandleFunc("/api/v1/servers", GetBackends)
	mux.HandleFunc("/api/v1/cert/", GetCerts)
	mux.HandleFunc("/api/v1/state/", GetBackendState)
	mux.HandleFunc("/api/v1/stats/", GetBackendStats)
	mux.HandleFunc("/api/v1/runtime/mem/alloc", GetRuntimeMemAlloc)
	mux.HandleFunc("/api/v1/version", GetVersion)
	mux.HandleFunc("/api/v1/sse/status/", BackendSSE)
	mux.HandleFunc("/api/v1/sse/global", GlobalSSE)
	mux.HandleFunc("/api/v1/backend/", SetBackend)
	mux.Handle("/", NewStaticHander("./web"))

	log.Println("Starting server on address " + port)
	return http.ListenAndServe(port, mux)
}
