package cli

import (
	"flag"
	"time"
)

type config struct {
	directory bool
	format    string
	timeout   time.Duration
}

func makeConfigFromFlags() (cfg config) {
	cfg = config{}

	flag.BoolVar(&cfg.directory, "d", false, "check all instances from the SecureDrop directory")
	flag.StringVar(&cfg.format, "f", "jsonl", "output format of the status checks")
	flag.DurationVar(&cfg.timeout, "t", 0, "general timeout, e.g. 20s")

	flag.Parse()

	return
}
