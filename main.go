package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/obrafy/planning/infrastructure"
)

var (
	version     = "undefined"
	environment = flag.String("environment", "", "The specific environment to run")
	showVersion = flag.Bool("version", false, "Show version and exit")
)

func main() {
	flag.Parse()

	if showVersion != nil && *showVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	if app, err := infrastructure.NewApp(*environment); err == nil {
		app.Run()
	} else {
		fmt.Printf("Error creating application: %s", err.Error())
	}
}
