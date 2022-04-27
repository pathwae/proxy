package api

import (
	"crypto/x509"
	"encoding/json"
	"net/http"
	"pathwae/proxy"
	"time"
)

type ServerStatus struct {
	Stat bool
}

// BackendSSE is a server side event handler for a given Backend.
func BackendSSE(w http.ResponseWriter, r *http.Request) {
	// prepare sse
	enableCORS(w)
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	//server name is set at the end of the url
	serverName := r.URL.Path[len("/api/v1/sse/status/"):]
	server := proxy.GetBackend(serverName)
	if server == nil {
		log.Printf("Server %s not found", serverName)
		return
	}

	// Register a listener of stats
	statListener := proxy.RegisterStatsListener(serverName)
	defer proxy.UnregisterStatsListener(statListener)

	// create a cert listener
	certLitener := proxy.RegisterCertListener(serverName)
	defer proxy.UnregisterCertListener(certLitener)

	// client close listener
	notify := w.(http.CloseNotifier).CloseNotify()

	for {
		select {
		case <-notify:
			log.Println("SSE: Client disconnected")
			return

		case cert := <-certLitener:
			info, err := x509.ParseCertificate(cert.Certificate[0])
			if err != nil {
				log.Println("SSE: Error parsing certificate")
				continue
			}
			cert.Leaf = info

			certInfo := CertInfo{
				Issuer:     cert.Leaf.Issuer.String(),
				Subject:    cert.Leaf.Subject.String(),
				NotBefore:  cert.Leaf.NotBefore.String(),
				NotAfter:   cert.Leaf.NotAfter.String(),
				CommonName: cert.Leaf.Subject.CommonName,
				DNSNames:   cert.Leaf.DNSNames,
			}

			certJson, err := json.Marshal(certInfo)
			if err != nil {
				log.Println("SSE: Error marshalling cert info")
				continue
			}
			// send sse
			w.Write([]byte("event: cert\n"))
			w.Write([]byte("data: " + string(certJson) + "\n\n"))
			w.(http.Flusher).Flush()

		case stat := <-statListener:
			// Send the stat to the client
			statJSON, _ := json.Marshal(stat)
			tosend := "event: stat\n" +
				"data: " + string(statJSON) + "\n\n"
			w.Write([]byte(tosend))
			w.(http.Flusher).Flush()

		case <-time.Tick(1 * time.Second):
			// send the server status
			server = proxy.GetBackend(serverName)
			serverStatus := proxy.IsBackendUp(*server)
			status := ServerStatus{Stat: serverStatus}
			message, _ := json.Marshal(&status)

			// the message type is "stats"
			tosend := "event: status\n" +
				"data: " + string(message) + "\n\n"
			w.Write([]byte(tosend))
			// flush
			w.(http.Flusher).Flush()
		}
	}

}

// GlobalSSE is a server side event handler for the reverse proxy itself (gives memory info, and others things).
func GlobalSSE(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// client close listener
	notify := w.(http.CloseNotifier).CloseNotify()

	changesListener := proxy.RegisterChangesListener()
	defer proxy.UnregisterChangesListener(changesListener)

	for {
		select {
		case <-notify:
			log.Println("Global SSE: Client disconnected")
			return

		case changes := <-changesListener:
			changesJSON, _ := json.Marshal(changes)
			tosend := "event: changes\n" +
				"data: " + string(changesJSON) + "\n\n"
			w.Write([]byte(tosend))
			w.(http.Flusher).Flush()

		case <-time.Tick(1 * time.Second):
			// send the server status
			memory := proxy.GetMemory()
			message, _ := json.Marshal(memory)

			// the message type is "stats"
			tosend := "event: memory\n" +
				"data: " + string(message) + "\n\n"
			w.Write([]byte(tosend))

			// flush
			w.(http.Flusher).Flush()
		}
	}
}
