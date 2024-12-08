Example

```
func main() {
	producer, err := NewProducer()
	if err != nil {
		log.Fatalf("Error creating producer: %v", err)
	}
	defer producer.Close()

	// Example with HouseWasSold projection
	event := NewEvent(
		"HouseWasSold",
		map[string]any{
			"price": 100.00,
		},
	)

	// Commit the event to Kafka and get the JSON data
	jsonData, err := event.Commit(&producer)
	if err != nil {
		panic(err)
	}

	// Now we can unmarshal the JSON back into the Event struct and check Projection
	dispatcher := NewDispatcher()
	dispatcher.Publish(jsonData)
}
```
