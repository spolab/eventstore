package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
	v1 "toremo.com/petclinic/eventstore/gen"
)

type PostgresDriver struct {
	v1.UnimplementedEventStoreServer
	db *sql.DB
}

func (pd *PostgresDriver) Append(ctx context.Context, req *v1.AppendRequest) (*v1.AppendResponse, error) {
	_, err := pd.db.Exec("CALL append_event($1, $2, $3, $4, $5, $6)", req.StreamId, req.ExpectedVersion, req.EventType, req.Encoding, req.Source, req.Data)
	if err != nil {
		return nil, err
	}
	return &v1.AppendResponse{}, nil
}

// GetEventsByStream retrieves all events belonging to a specific stream_id
func (pd *PostgresDriver) Get(ctx context.Context, req *v1.GetRequest) (*v1.GetResponse, error) {
	logger := log.With().Str("function", "Get").Str("stream_id", req.StreamId).Logger() // Create a logger with the function name

	logger.Info().Msg("Retrieving events by stream_id") // Log an info message

	rows, err := pd.db.Query("SELECT event_id, stream_id, stream_version, event_type, event_encoding, event_source, event_data, event_ts FROM events WHERE stream_id=$1 ORDER BY stream_version", req.StreamId)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get events by stream_id") // Log an error message
		return nil, fmt.Errorf("failed to get events by stream_id: %v", err)
	}
	defer rows.Close()

	events := make([]*v1.Event, 0)
	for rows.Next() {
		event := &v1.Event{}
		err := rows.Scan(&event.EventId, &event.StreamId, &event.Version, &event.EventType, &event.Encoding, &event.Source, &event.Data, &event.Timestamp)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to scan events") // Log an error message
			return nil, fmt.Errorf("failed to scan events: %v", err)
		}
		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		logger.Error().Err(err).Msg("Error iterating events") // Log an error message
		return nil, fmt.Errorf("error iterating events: %v", err)
	}

	logger.Info().Int("num_events", len(events)).Msg("Retrieved events by stream_id") // Log an info message

	return &v1.GetResponse{Events: events}, nil
}

func NewPostgresDriver(db *sql.DB) (*PostgresDriver, error) {
	if db == nil {
		return nil, errors.New("database must not be nil")
	}
	return &PostgresDriver{db: db}, nil
}
