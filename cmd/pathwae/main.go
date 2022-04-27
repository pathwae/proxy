package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	stdlog "log"
	"os"
	"pathwae/proxy"
	"pathwae/service"
	"runtime"
)

var (
	Version = "0.0.1" // changed at build time
	log     = stdlog.New(os.Stdout, "[CMD] ", stdlog.Lmsgprefix|stdlog.LstdFlags)
)

const globalConfFile = "/global/config.yaml"

func main() {
	// set the max proc to the number of cpus
	runtime.GOMAXPROCS(runtime.NumCPU())

	showVersion := flag.Bool("version", false, "Show version")
	flag.Parse()
	if *showVersion {
		fmt.Println(Version)
		return
	}
	proxy.Version = Version

	// can be empty, not a problem
	confToLoad := os.Getenv("CONFIG_FILE")

	confFile := globalConfFile
	if os.Getenv("TESTMODE") == "1" {
		confFile = "/tmp/config.yaml"
	}

	// if CONFIG is provided, so write the content in the temporary directory
	if len(os.Getenv("CONFIG")) > 0 {
		err := ioutil.WriteFile(confFile, []byte(os.Getenv("CONFIG")), 0644)
		if err != nil {
			log.Fatalf("Failed to write configuration file: %v", err)
		}
		// force
		confToLoad = confFile
	}

	// now, if confToLoad is empty, there is a problem...
	if len(confToLoad) == 0 {
		log.Fatalf("No configuration CONFIG_FILE of CONFIG provided")
	}

	// else, we can now read the configuration file
	content, err := ioutil.ReadFile(confToLoad)
	if err != nil {
		log.Fatal(err)
	}
	conf := string(content)

	if conf == "" {
		log.Println("You didn't provide a config file in CONFIG_FILE or CONFIG env variable, this means that no backends will be loaded")
	}

	// get conf from CONF envuronment
	service.Start(conf)

}
