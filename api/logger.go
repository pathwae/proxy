package api

import (
	stdlog "log"
	"os"
)

var log = stdlog.New(os.Stdout, "[API] ", stdlog.Lmsgprefix|stdlog.LstdFlags)
