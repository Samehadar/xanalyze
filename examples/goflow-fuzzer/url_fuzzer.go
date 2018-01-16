package main

import (
	"log"
	"os"

	"github.com/sniperkit/xanalyze/examples/goflow-fuzzer/configuration"
	"github.com/sniperkit/xanalyze/examples/goflow-fuzzer/flow"
)

type urlFuzzer struct{}

func newURLFuzzer() *urlFuzzer {
	return new(urlFuzzer)
}

func (u *urlFuzzer) run() {
	log.SetOutput(os.Stdout)

	configuration := u.readConfiguration()
	fuzz := flow.NewFuzz(configuration)
	fuzz.Start()
}

func (u *urlFuzzer) readConfiguration() *configuration.Configuration {
	configurationFactory := configuration.NewFactory()
	return configurationFactory.FromCommandLine()
}
