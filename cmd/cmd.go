package cmd

import (
	"flag"
	"log"
	"os"
)

var logger = log.New(os.Stdout, "", 0)
var verbose bool

func addVerboseFlag(subCommand *flag.FlagSet) {
	subCommand.BoolFunc("verbose", "enable more detailed logging", func(string) error {
		verbose = true
		return nil
	})
}

func verboseLog(message string) {
	if verbose {
		logger.Print(message)
	}
}