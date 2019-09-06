package main

import (
	"fmt"
	"os"

	"github.com/mfojtik/deptrack/pkg/release"
	"github.com/mfojtik/deptrack/pkg/report"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintf(os.Stderr, "Usage: %s <config> <gitDir> <reportDir>\n", os.Args[0])
		os.Exit(255)
	}
	config, err := release.ReadConfig(os.Args[1])
	if err != nil {
		panic(err)
	}

	for _, r := range config.Releases {
		fmt.Printf("Analyzing release %q (%q)...\n", r.Name, os.Args[2])

		status, err := report.NewReport(os.Args[2], r)
		if err != nil {
			panic(err)
		}

		if err := report.WriteReport(os.Args[3], *status); err != nil {
			panic(err)
		}
	}
}
