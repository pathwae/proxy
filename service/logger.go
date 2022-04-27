package service

import (
	stdlog "log"
	"os"
)

var log = stdlog.New(os.Stdout, "[SERVICE] ", stdlog.Lmsgprefix|stdlog.LstdFlags)
