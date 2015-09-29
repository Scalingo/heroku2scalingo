package signals

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Scalingo/heroku2scalingo/io"
)

var (
	CatchQuitSignals = true
)

func Handle() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	for sig := range signals {
		if CatchQuitSignals {
			io.Warnf("%v catched, aborting…\n", sig)
			os.Exit(-127)
		}
	}
}
