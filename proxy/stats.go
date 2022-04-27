package proxy

import (
	"net/http"
	"runtime"
	"sync"
)

var (
	Stats    map[string][]*Stat
	StatChan chan http.Request
	lock     sync.Mutex
)

type Stat struct {
	Path   string `json:"path"`
	Method string `json:"method"`
}

func init() {
	lock = sync.Mutex{}
	StatChan = make(chan http.Request, 100)
	Stats = make(map[string][]*Stat)
	go func() {
		for {
			doStat(<-StatChan) // send the request to stats handler
		}
	}()
}

// GetStats returns the stats for a hostname.
func GetStat(hostname string) []*Stat {
	lock.Lock()
	defer lock.Unlock()
	return Stats[hostname]
}

// GetMemory returns the memory usage of the current process.
func GetMemory() int {
	// get the runtime memory used by the process
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return int(m.Alloc)
}
