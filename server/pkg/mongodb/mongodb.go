package mongodb

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	v1 "toremo.com/petclinic/eventstore/gen"
	"toremo.com/petclinic/eventstore/pkg/errors"
)

const (
	StreamsCollectionName = "streams"
	EventsCollectionName  = "events"
)

type MongoJournal struct {
	v1.UnimplementedJournalServer
	Client       *mongo.Client
	DatabaseName string
}

// Stream represents a stream document in the MongoDB streams collection.
type Stream struct {
	ID            string `bson:"_id"`
	StreamVersion int64  `bson:"stream_version"`
}

// Event represents an event document in the MongoDB events collection.
type Event struct {
	ID            string `bson:"_id,omitempty"`
	StreamID      string `bson:"stream_id"`
	StreamVersion int64  `bson:"stream_version"`
	Kind          string `bson:"kind"`
	Encoding      string `bson:"encoding"`
	Source        string `bson:"source"`
	Data          []byte `bson:"data"`
	Timestamp     string `bson:"timestamp"`
}

// Append appends a new event to the specified stream.
func (md *MongoJournal) AppendEvent(ctx context.Context, req *v1.AppendEventRequest) (*v1.AppendEventResponse, error) {
	log.Info().Msg("Start Append")

	// Get the streams and events collections
	streamsCollection := md.Client.Database(md.DatabaseName).Collection(StreamsCollectionName)
	eventsCollection := md.Client.Database(md.DatabaseName).Collection(EventsCollectionName)

	// Define a filter and update document for the upsert operation
	filter := bson.M{"_id": req.StreamId, "stream_version": req.ExpectedVersion}
	update := bson.M{"$inc": bson.M{"stream_version": 1}}

	if req.ExpectedVersion == 0 {
		// Perform the upsert operation on the streams collection
		_, err := streamsCollection.InsertOne(ctx, &Stream{ID: req.StreamId, StreamVersion: 1})
		if err != nil {
			if mongo.IsDuplicateKeyError(err) {
				return nil, errors.InvalidStreamVersionError()
			}
			return nil, err
		}
	} else {
		// Perform the update operation on the streams collection
		updateResult, err := streamsCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			return nil, err
		}
		if updateResult.ModifiedCount != 1 {
			return nil, errors.InvalidStreamVersionError()
		}
	}

	// At this point, we have ensured that the stream version is correct, so we can insert the new event
	newEvent := Event{
		StreamID:      req.GetStreamId(),
		StreamVersion: req.ExpectedVersion + 1,
		Kind:          req.GetEventType(),
		Source:        req.GetSource(),
		Encoding:      req.GetEncoding(),
		Data:          req.GetData(),
		Timestamp:     time.Now().String(),
	}

	if _, err := eventsCollection.InsertOne(ctx, newEvent); err != nil {
		return nil, err
	}

	log.Info().Msg("End Append")
	return &v1.AppendEventResponse{}, nil
}

// Append appends a new event to the specified stream.
// func (md *MongoDriver) AppendEvent(ctx context.Context, req *v1.AppendEventRequest) (*v1.AppendEventResponse, error) {
// 	log.Info().Msg("Start Append")
// 	// Start a new session and defer its closure
// 	// Get the streams and events collections
// 	streamsCollection := md.Client.Database(md.DatabaseName).Collection(StreamsCollectionName)
// 	eventsCollection := md.Client.Database(md.DatabaseName).Collection(EventsCollectionName)

// 	// Find the stream with the specified streamID
// 	streamFilter := bson.M{"_id": req.StreamId}
// 	streamResult := streamsCollection.FindOne(ctx, streamFilter)

// 	var stream Stream
// 	if err := streamResult.Decode(&stream); err != nil {
// 		// Stream not found, so create a new stream with stream_version 0
// 		stream = Stream{
// 			ID:            req.StreamId,
// 			StreamVersion: 1,
// 		}
// 		if _, err := streamsCollection.InsertOne(ctx, stream); err != nil {
// 			if mongo.IsDuplicateKeyError(err) {
// 				return nil, fmt.Errorf("optimistic concurrency violation: another client created the stream at the same time")
// 			}
// 			return nil, err
// 		}
// 	} else {
// 		// Stream found, so check for optimistic concurrency
// 		if stream.StreamVersion != req.ExpectedVersion {
// 			return nil, fmt.Errorf("optimistic concurrency violation: expected version %d, actual version %d", req.GetExpectedVersion(), stream.StreamVersion)
// 		}
// 		stream.StreamVersion++
// 		updateResult, err := streamsCollection.UpdateOne(ctx, streamFilter, bson.M{"$set": bson.M{"stream_version": stream.StreamVersion}})
// 		if err != nil {
// 			return nil, err
// 		}
// 		if updateResult.ModifiedCount != 1 {
// 			return nil, fmt.Errorf("failed to update stream version")
// 		}
// 	}

// 	newEvent := Event{
// 		StreamID:      req.GetStreamId(),
// 		StreamVersion: stream.StreamVersion,
// 		Kind:          req.GetEventType(),
// 		Source:        req.GetSource(),
// 		Encoding:      req.GetEncoding(),
// 		Data:          req.GetData(),
// 		Timestamp:     time.Now().String(),
// 	}

// 	if _, err := eventsCollection.InsertOne(ctx, newEvent); err != nil {
// 		return nil, err
// 	}
// 	log.Info().Msg("End Append")
// 	return &v1.AppendEventResponse{}, nil
// }

func (d *MongoJournal) GetStreamEvents(ctx context.Context, req *v1.GetStreamEventsRequest) (*v1.GetStreamEventsResponse, error) {
	eventsCollection := d.Client.Database(d.DatabaseName).Collection(EventsCollectionName)

	filter := bson.M{"stream_id": req.StreamId}

	cursor, err := eventsCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var events []*v1.Event
	for cursor.Next(ctx) {
		event := &Event{}
		err := cursor.Decode(&event)
		if err != nil {
			return nil, err
		}
		events = append(events, &v1.Event{StreamId: event.StreamID, Version: event.StreamVersion, EventType: event.Kind, Encoding: event.Encoding, Source: event.Source, Data: event.Data, Timestamp: event.Timestamp})
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return &v1.GetStreamEventsResponse{Events: events}, nil
}

func NewMongoJournal(ctx context.Context, client *mongo.Client, dbName string) (*MongoJournal, error) {
	eventsCollection := client.Database(dbName).Collection(EventsCollectionName)
	streamIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "stream_id", Value: 1},
			{Key: "stream_version", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}
	_, err := eventsCollection.Indexes().CreateOne(ctx, streamIndex)
	if err != nil {
		return nil, err
	}
	log.Info().Msg("initialized mongo driver")
	return &MongoJournal{
		Client:       client,
		DatabaseName: dbName,
	}, nil
}
