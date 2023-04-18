package postgres_test

import (
	"database/sql"
	"os"
	"testing"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"toremo.com/petclinic/eventstore/pkg/postgres"
)

func TestAppendEvent(t *testing.T) {
	// Open a connection to the database
	db, err := sql.Open("postgres", os.Getenv(("POSTGRES_URL")))
	require.NoError(t, err)
	defer db.Close()

	// Create a new stream_id and event_id
	streamID := uuid.NewString()
	eventType := "event_type_1"
	eventEncoding := "application/json"
	eventSource := "TestAppendEvent"
	eventData := `{"key": "value"}`

	t.Run("FirstInsertReturns1", func(t *testing.T) {
		// Call the function with expected_version = 0 (should insert the first event)
		err := postgres.AppendEvent(db, streamID, 0, eventType, eventEncoding, eventSource, []byte(eventData))
		require.NoError(t, err)
	})

	t.Run("SecondInsertReturns2", func(t *testing.T) {
		// Call the function with expected_version = 1 (should insert the second event)
		err := postgres.AppendEvent(db, streamID, 1, "event_type_2", "event_encoding_2", "event_source_2", []byte(`{"key": "value"}`))
		require.NoError(t, err)
	})

	t.Run("StaleInsertReturnsError", func(t *testing.T) {
		// Call the function with expected_version = 0 (should fail because stream version is 2)
		err = postgres.AppendEvent(db, streamID, 1, "event_type_3", "event_encoding_3", "event_source_3", []byte(`{"key": "value"}`))
		assert.Error(t, err)
		assert.Equal(t, "pq: Expected stream_version 1 but got 2", err.Error())
	})
}

func TestGetEventsByStream(t *testing.T) {
	// Set up the database connection
	db, err := sql.Open("postgres", "postgres://postgres:password123@localhost/?sslmode=disable")
	require.NoError(t, err)
	defer db.Close()

	// Clear the events and streams tables before starting the test
	_, err = db.Exec("DELETE FROM events")
	require.NoError(t, err)
	_, err = db.Exec("DELETE FROM streams")
	require.NoError(t, err)

	// Define some test data
	streamID := uuid.NewString()
	eventData1 := `{"name": "event1"}`
	eventData2 := `{"name": "event2"}`

	// Insert some test data into the database using AppendEvent
	require.NoError(t, postgres.AppendEvent(db, streamID, 0, "event1", "encoding1", "source1", []byte(eventData1)))
	require.NoError(t, postgres.AppendEvent(db, streamID, 1, "event2", "encoding2", "source2", []byte(eventData2)))

	// Call the function being tested
	events, err := postgres.GetEventsByStream(db, streamID)

	// Verify the results
	require.NoError(t, err)
	assert.Len(t, events, 2)
	assert.Equal(t, eventData1, string(events[0].Data))
	assert.Equal(t, eventData2, string(events[1].Data))
	assert.Equal(t, streamID, events[0].StreamID)
	assert.Equal(t, streamID, events[1].StreamID)
	assert.Equal(t, int64(0), events[0].StreamVersion)
	assert.Equal(t, int64(1), events[1].StreamVersion)
	assert.Equal(t, "event1", events[0].Type)
	assert.Equal(t, "event2", events[1].Type)
	assert.Equal(t, "encoding1", events[0].Encoding)
	assert.Equal(t, "encoding2", events[1].Encoding)
	assert.Equal(t, "source1", events[0].Source)
	assert.Equal(t, "source2", events[1].Source)
}
