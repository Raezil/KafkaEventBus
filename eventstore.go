package KafkaEventBus

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

type Result struct {
	Message string
}

type EventStore struct {
	Dispatcher *Dispatcher
	Producer   *sarama.SyncProducer
}

func NewEventStore() *EventStore {
	dispatcher := NewDispatcher()
	producer, err := NewProducer()
	if err != nil {
		panic(err)
	}
	return &EventStore{
		Dispatcher: dispatcher,
		Producer:   &producer,
	}
}

type Dispatcher map[string]func(map[string]interface{}) (Result, error)

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		"HouseWasSold": func(m map[string]interface{}) (Result, error) {
			price, ok := m["price"].(float64)
			if !ok {
				return Result{}, fmt.Errorf("price not provided or invalid")
			}

			result := fmt.Sprintf("House was sold for %.2f", price)
			return Result{Message: result}, nil
		},
	}
}

func (eventstore *EventStore) Publish(jsonData []byte) {
	var event Event
	err := json.Unmarshal(jsonData, &event)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
	handler, exists := (*eventstore.Dispatcher)[event.Projection]
	if !exists {
		panic("handler does not exist!")
	}
	result, err := handler(event.Args)
	if err != nil {
		panic(err)
	}
	log.Println(result)
}

type Event struct {
	ID         string         `json:"id"`
	Projection string         `json:"projection"`
	Args       map[string]any `json:"args"`
}

type HouseWasSold struct{}

func NewProducer() (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, config)
	if err != nil {
		log.Fatalf("Error creating producer: %v", err)
	}
	return producer, err
}

func NewEvent(projection string, args map[string]any) *Event {
	id := uuid.New()
	return &Event{
		ID:         id.String(),
		Projection: projection,
		Args:       args,
	}
}

func (eventstore *EventStore) Commit(event *Event) ([]byte, error) {
	// Marshal the struct into JSON
	jsonData, err := json.Marshal(event)
	if err != nil {
		log.Fatalf("Error marshalling struct to JSON: %v", err)
	}

	// Send the message as JSON
	message := &sarama.ProducerMessage{
		Topic: getStructName(event.Projection),
		Value: sarama.StringEncoder(jsonData),
	}

	// Send the message to Kafka
	partition, offset, err := (*eventstore.Producer).SendMessage(message)
	if err != nil {
		log.Fatalf("Error sending message: %v", err)
		return nil, err
	}

	log.Printf("Message sent to partition %d with offset %d\n", partition, offset)
	return jsonData, nil
}

// Updated getStructName to handle the `any` type properly
func getStructName(v interface{}) string {
	// Handle the case where v is a pointer (if the struct is passed as a pointer)
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem() // Dereference the pointer to get the actual struct type
	}
	// Return the name of the struct type
	return t.Name()
}
