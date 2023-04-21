package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"toremo.com/petclinic/eventstore/cmd/event"
	"toremo.com/petclinic/eventstore/cmd/server"
)

const (
	ENV_BIND   = "EVENSTORE_BIND"
	ENV_DRIVER = "EVENSTORE_DRIVER"
	ENV_DB_URL = "EVENSTORE_DB_URL"
)

func main() {
	var (
		bind   string
		driver string
		dburl  string
	)

	root := &cobra.Command{
		Use: "eventstore",
	}
	root.AddCommand(server.ServerCommand)
	root.AddCommand(event.EventCommand)

	if err := root.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(0)

	// Load the configuration parameters
	flag.StringVar(&bind, "bind", "0.0.0.0:9000", "The address to bind the service to")
	flag.StringVar(&driver, "driver", "", "The database driver to use (mongodb or postgresql)")
	flag.StringVar(&dburl, "db-url", "", "The connection URL to connect to the database")
	flag.Parse()

	// Validate the parameters
	if bind == "" {
		log.Fatal().Msg("bind address must be specified")
	}
}
