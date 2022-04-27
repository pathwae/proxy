package proxy

import (
	"errors"
	"reflect"
	"sync"

	"gopkg.in/yaml.v2"
)

type Change struct {
	Name    string  `json:"name"`
	Backend Backend `json:"backend"`
}

var (
	backendList map[string]*Backend
	changeLock  sync.Mutex
)

// Backend is a proxy configured service from conf.
type Backend struct {
	// To is the url where to send the request.
	To string `yaml:"to" json:"to"`

	// ForceSSL forces the connection to be ssl. If true, any HTTP connection will be redirected to HTTPS.
	ForceSSL bool `yaml:"force_ssl" json:"force_ssl"`

	// Enabled is the flag to enable or disable the service.
	Enabled *bool `yaml:"enabled,omitempty" json:"enabled,omitempty" default:"true"`
}

// LoadServers loads the servers from the conf file.
func LoadYAMLConfig(content string) map[string]*Backend {
	var servers map[string]*Backend

	err := yaml.Unmarshal([]byte(content), &servers)
	if err != nil {
		log.Fatal(err)
	}

	// give information about the servers
	for from, to := range servers {
		// default is to enable the backend
		if to.Enabled == nil {
			to.Enabled = new(bool)
			// get the default value from the tags
			tag, _ := reflect.TypeOf(to).Elem().FieldByName("Enabled")
			if tag.Tag.Get("default") != "" {
				*to.Enabled = tag.Tag.Get("default") == "true"
			} else {
				*to.Enabled = true
			}
		}

		log.Printf("Configured %s -> %s", from, to.To)
	}

	backendList = servers
	return servers
}

// SetBackend sets the backend for the given host.
func SetBackend(name string, b Backend) error {
	log.Println("Change backend:", name, b)
	if b.To == "" {
		return errors.New("Backend url is empty")
	}
	if _, ok := backendList[name]; !ok {
		return errors.New("Backend url doesn't exists in the config")
	}
	backendList[name] = &b

	go func() {
		changeLock.Lock()
		defer changeLock.Unlock()
		for _, listener := range changesListeners {
			listener <- &Change{
				Name:    name,
				Backend: b,
			}
		}
	}()
	return nil
}
