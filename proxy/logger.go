package proxy

import (
	stdlog "log"
	"os"
)

var log = stdlog.New(os.Stdout, "[PROXY] ", stdlog.Lmsgprefix|stdlog.LstdFlags)
