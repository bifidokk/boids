package main

import (
	"fmt"
	"math/rand"
	"time"
)

var availableTopics = map[string]string{
	"A": "AAA",
	"B": "BBB",
	"C": "CCC",
	"D": "DDD",
}

func main() {
	broker := NewBroker()
	s1 := broker.AddSubscriber()
	broker.Subscribe(s1, availableTopics["A"])
	broker.Subscribe(s1, availableTopics["B"])

	s2 := broker.AddSubscriber()
	broker.Subscribe(s2, availableTopics["C"])
	broker.Subscribe(s2, availableTopics["D"])

	go publish(broker)

	go s1.Listen()
	go s2.Listen()

	fmt.Scanln()
	fmt.Println("Done")
}

func publish(broker *Broker) {
	topicKeys := make([]string, 0, len(availableTopics))
	topicValues := make([]string, 0, len(availableTopics))

	for k, v := range availableTopics {
		topicKeys = append(topicKeys, k)
		topicValues = append(topicValues, v)
	}

	for {
		randValue := topicValues[rand.Intn(len(topicValues))]
		msg := fmt.Sprintf("%f", rand.Float64())
		fmt.Printf("Publishing %s to %s topic\n", msg, randValue)
		go broker.Publish(randValue, msg)

		r := rand.Intn(4)
		time.Sleep(time.Duration(r) * time.Second)
	}
}
