package server

import (
	"context"
	"net"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	eventstore "toremo.com/petclinic/eventstore/gen"
	"toremo.com/petclinic/eventstore/pkg/mongodb"
)

const (
	FLAG_DRIVER    = "driver"
	FLAG_DB_URL    = "db-url"
	FLAG_GRPC_ADDR = "grpc-addr"
)

var ServerCommand *cobra.Command

func init() {
	ServerCommand = &cobra.Command{
		Use:   "server",
		Short: "Start the EventStore server",
		Args:  cobra.ArbitraryArgs,
		RunE:  serverCommandImpl,
	}
	ServerCommand.Flags().String(FLAG_DRIVER, "", "Database driver to use (postgresql, mongodb)")
	ServerCommand.Flags().String(FLAG_DB_URL, "", "Connection string for the database")
	ServerCommand.Flags().String(FLAG_GRPC_ADDR, "0.0.0.0:9000", "Address to bind the gRPC service")
	ServerCommand.MarkFlagRequired(FLAG_DRIVER)
	ServerCommand.MarkFlagRequired(FLAG_DB_URL)
}

func serverCommandImpl(cmd *cobra.Command, args []string) error {
	bind := cmd.Flags().Lookup(FLAG_GRPC_ADDR).Value.String()
	dburl := cmd.Flags().Lookup(FLAG_DB_URL).Value.String()
	// Allocate a listener
	listener, err := net.Listen("tcp", bind)
	if err != nil {
		return err
	}

	// Establishes a connection to the mongodb database
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dburl))
	if err != nil {
		return err
	}
	defer func() {
		_ = client.Disconnect(ctx)
	}()

	// Creates the server and registers the service
	server := grpc.NewServer()
	impl, err := mongodb.NewMongoDriver(ctx, client, "eventstore")
	if err != nil {
		return err
	}
	eventstore.RegisterEventStoreServer(server, impl)

	// Start the serveer
	log.Info().Str("address", bind).Msg("server started")
	err = server.Serve(listener)
	if err != nil {
		return err
	}
	return nil
}
