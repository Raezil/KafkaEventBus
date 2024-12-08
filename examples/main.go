package main

import (
	"log"

	. "github.com/Raezil/KafkaEventBus"
)

func Commit(eventstore *EventStore, event *Event) []byte {
	jsonData, err := eventstore.Commit(event)
	if err != nil {
		panic(err)
	}
	return jsonData
}

func main() {
	producer, err := NewProducer()
	if err != nil {
		log.Fatalf("Error creating producer: %v", err)
	}
	defer producer.Close()
	eventStore := NewEventStore()
	eventStore.Publish(Commit(eventStore, NewEvent(
		"HouseWasSold",
		map[string]any{
			"price": 100.00,
		},
	)))
}
