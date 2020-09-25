package cli

import (
	"flag"
	"time"
)

type config struct {
	concurrency int
	directory   bool
	format      string
	urls        []string
	timeout     time.Duration
}

func makeConfigFromFlags() (cfg config) {
	cfg = config{}

	flag.IntVar(&cfg.concurrency, "c", 0, "maximum number of concurrent checks")
	flag.BoolVar(&cfg.directory, "d", false, "check all instances from the SecureDrop directory")
	flag.StringVar(&cfg.format, "f", "jsonl", "output format of the status checks")
	flag.DurationVar(&cfg.timeout, "t", 0, "general timeout, e.g. 20s")

	flag.Parse()

	cfg.urls = flag.Args()

	return
}
