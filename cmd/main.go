package main

import (
	"context"
	"net"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	v1 "toremo.com/petclinic/eventstore/gen"
	"toremo.com/petclinic/eventstore/pkg/mongodb"
)

const (
	LISTEN_ADDRESS = "LISTEN_ADDRESS"
)

func main() {
	logger := log.With().Str("component", "main").Logger()

	// Read the parameters from the environment
	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		logger.Fatal().Err(err).Msg("error reading the .env file")
	}
	addr := os.Getenv(LISTEN_ADDRESS)

	// Allocate a listener
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Fatal().Str("address", addr).Err(err).Msg("error allocating the tcp listener")
	}

	// Establishes a connection to the mongodb database
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URL")))
	if err != nil {
		logger.Fatal().Err(err).Msg("error establishing a connection to mongodb")
	}
	defer func() {
		_ = client.Disconnect(ctx)
	}()

	// Creates the server and registers the service
	server := grpc.NewServer()
	impl, err := mongodb.NewMongoDriver(ctx, client, "eventstore")
	if err != nil {
		logger.Fatal().Err(err).Msg("creating the driver")
	}
	v1.RegisterEventStoreServer(server, impl)

	// Start the serveer
	logger.Info().Str("address", addr).Msg("server started")
	err = server.Serve(listener)
	if err != nil {
		logger.Fatal().Err(err).Msg("starting the server")
	}
}
