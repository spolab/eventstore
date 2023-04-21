package event

import (
	"fmt"

	"github.com/spf13/cobra"
)

const GRPC_ADDR = "grpc-addr"

var EventCommand *cobra.Command

func init() {
	EventGetCommand := &cobra.Command{
		Use:   "get [event id]",
		Short: "Displays the contents of a specific event",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("event get")
		},
	}
	EventGetCommand.ArgAliases = []string{"id"}

	EventCommand = &cobra.Command{
		Use:   "event",
		Short: "Event-related commands",
	}
	EventCommand.PersistentFlags().String(GRPC_ADDR, "localhost:9000", "Address of the GRPC listener")
	EventCommand.AddCommand(EventGetCommand)
}
