package KafkaEventBus

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDispatcher(t *testing.T) {
	dispatcher := NewDispatcher()

	// Test if the dispatcher contains the "HouseWasSold" event handler
	handler, exists := (*dispatcher)["HouseWasSold"]
	assert.True(t, exists, "Expected HouseWasSold handler to exist")
	assert.NotNil(t, handler, "Handler should not be nil")

	// Test invalid event data
	eventArgs := map[string]interface{}{"price": "invalid"}
	result, err := handler(eventArgs)
	assert.Error(t, err, "Expected error due to invalid price type")
	assert.Empty(t, result.Message, "Message should be empty when error occurs")
}

func TestHouseWasSoldEvent(t *testing.T) {
	dispatcher := NewDispatcher()

	// Valid price value
	eventArgs := map[string]interface{}{"price": 300000.0}
	handler := (*dispatcher)["HouseWasSold"]
	result, err := handler(eventArgs)

	// Assert no error and correct result message
	assert.NoError(t, err, "Expected no error")
	assert.Equal(t, "House was sold for 300000.00", result.Message, "Expected correct result message")
}

func TestNewEvent(t *testing.T) {
	event := NewEvent("HouseWasSold", map[string]interface{}{"price": 200000.0})

	// Test if the event is created correctly
	assert.NotNil(t, event, "Event should not be nil")
	assert.Equal(t, "HouseWasSold", event.Projection, "Expected projection to be 'HouseWasSold'")
	assert.Equal(t, map[string]interface{}{"price": 200000.0}, event.Args, "Expected event args to match")
}

func TestCommitEvent(t *testing.T) {
	// Setting up a mock Kafka producer
	producer, err := NewProducer()
	if err != nil {
		t.Fatalf("Error setting up Kafka producer: %v", err)
	}

	// Create an event to commit
	event := NewEvent("HouseWasSold", map[string]interface{}{"price": 250000.0})

	// Commit the event (send to Kafka)
	_, err = event.Commit(&producer)
	assert.NoError(t, err, "Expected no error while committing event")
}

func TestPublishEvent(t *testing.T) {
	dispatcher := NewDispatcher()

	// Create an event in JSON format
	event := NewEvent("HouseWasSold", map[string]interface{}{"price": 100000.0})
	eventData, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("Error marshalling event: %v", err)
	}

	// Publish event
	dispatcher.Publish(eventData)

	// You can add some assertions based on the log output or any side effects
	// For testing purposes, we would need to verify logging output or mock logging
}
