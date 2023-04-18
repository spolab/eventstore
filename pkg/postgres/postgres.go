package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
)

func AppendEvent(db *sql.DB, streamID string, expectedVersion int, eventType string, eventEncoding string, eventSource string, eventData []byte) error {
	_, err := db.Exec("CALL append_event($1, $2, $3, $4, $5, $6)", streamID, expectedVersion, eventType, eventEncoding, eventSource, eventData)
	if err != nil {
		return err
	}
	return nil
}

// Event represents a single event in the database
type Event struct {
	EventID        int64     `json:"event_id"`
	StreamID       string    `json:"stream_id"`
	StreamVersion  int       `json:"stream_version"`
	Type           string    `json:"event_type"`
	Encoding       string    `json:"event_encoding"`
	Source         string    `json:"event_source"`
	Data           string    `json:"event_data"`
	EventTimestamp time.Time `json:"event_ts"`
}

// GetEventsByStream retrieves all events belonging to a specific stream_id
func GetEventsByStream(db *sql.DB, streamID string) ([]Event, error) {
	logger := log.With().Str("function", "GetEventsByStream").Logger() // Create a logger with the function name

	logger.Info().Str("stream_id", streamID).Msg("Retrieving events by stream_id") // Log an info message

	rows, err := db.Query("SELECT event_id, stream_id, stream_version, event_type, event_data, event_ts FROM events WHERE stream_id=$1", streamID)
	if err != nil {
		logger.Error().Err(err).Str("stream_id", streamID).Msg("Failed to get events by stream_id") // Log an error message
		return nil, fmt.Errorf("failed to get events by stream_id: %v", err)
	}
	defer rows.Close()

	events := make([]Event, 0)
	for rows.Next() {
		event := Event{}
		err := rows.Scan(&event.EventID, &event.StreamID, &event.StreamVersion, &event.Type, &event.Data, &event.EventTimestamp)
		if err != nil {
			logger.Error().Err(err).Str("stream_id", streamID).Msg("Failed to scan events") // Log an error message
			return nil, fmt.Errorf("failed to scan events: %v", err)
		}
		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		logger.Error().Err(err).Str("stream_id", streamID).Msg("Error iterating events") // Log an error message
		return nil, fmt.Errorf("error iterating events: %v", err)
	}

	logger.Info().Int("num_events", len(events)).Msg("Retrieved events by stream_id") // Log an info message

	return events, nil
}
