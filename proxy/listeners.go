package proxy

import (
	"crypto/tls"
	"net/http"
)

var (
	certListeners    = make(map[string][]chan *tls.Certificate)
	StatListeners    = make(map[string][]chan *Stat)
	changesListeners = make([]chan *Change, 0)
)

// RegisterCertListener registers a channel to receive a certificate for the given host. This is useful for the SSE API endpoint.
func RegisterCertListener(host string) chan *tls.Certificate {
	certLock.Lock()
	defer certLock.Unlock()
	if _, ok := certListeners[host]; !ok {
		certListeners[host] = make([]chan *tls.Certificate, 0)
	}
	listener := make(chan *tls.Certificate)
	certListeners[host] = append(certListeners[host], listener)
	return listener
}

// UnregisterCertListener unregisters a channel from receiving a certificate for the given host.
func UnregisterCertListener(cl chan *tls.Certificate) {
	certLock.Lock()
	defer certLock.Unlock()
	close(cl)
	for host, listeners := range certListeners {
		for i, listener := range listeners {
			if listener == cl {
				certListeners[host] = append(certListeners[host][:i], certListeners[host][i+1:]...)
				return
			}
		}
	}

}

// RegisterChangeListener registers a channel to receive a change for the given host. This is useful for the SSE API endpoint.
func RegisterChangesListener() chan *Change {
	changeLock.Lock()
	defer changeLock.Unlock()
	listener := make(chan *Change)
	changesListeners = append(changesListeners, listener)
	return listener
}

// UnregisterChangesListener unregisters a channel from receiving changes.
func UnregisterChangesListener(listener chan *Change) {
	changeLock.Lock()
	defer changeLock.Unlock()
	close(listener)
	for i, l := range changesListeners {
		if l == listener {
			changesListeners = append(changesListeners[:i], changesListeners[i+1:]...)
			return
		}
	}
}

// RegisterStatsListener registers a listener for stats.
func RegisterStatsListener(hostname string) chan *Stat {
	lock.Lock()
	defer lock.Unlock()
	ch := make(chan *Stat, 100)
	StatListeners[hostname] = append(StatListeners[hostname], ch)
	return ch
}

// UnregisterStatsListener unregisters a listener for stats.
func UnregisterStatsListener(ch chan *Stat) {
	lock.Lock()
	defer lock.Unlock()
	close(ch)
	for hostname, listeners := range StatListeners {
		for i, listener := range listeners {
			if listener == ch {
				listeners = append(listeners[:i], listeners[i+1:]...)
				StatListeners[hostname] = listeners
				return
			}
		}
	}
}

// doStats sends the stats to all registered listeners.
func doStat(req http.Request) {
	lock.Lock()
	defer lock.Unlock()
	hostname := req.Host
	path := req.URL.Path
	method := req.Method
	// Send the stat to all listeners for this hostname
	go func() {
		for _, listener := range StatListeners[hostname] {
			listener <- &Stat{path, method}
		}
	}()
	// Add the stat to the hostname's list of stats
	Stats[hostname] = append(Stats[hostname], &Stat{path, method})
}
