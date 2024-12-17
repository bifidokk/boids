package main

import (
	"fmt"
	"math/rand"
	"sync"
)

type Subscriber struct {
	id       string
	messages chan *Message
	topics   map[string]bool
	active   bool
	mutex    sync.RWMutex
}

func NewSubscriber() (string, *Subscriber) {
	id := string(rand.Intn(1000000000))
	return id, &Subscriber{
		id:       id,
		messages: make(chan *Message),
		topics:   map[string]bool{},
		active:   true,
	}
}

func (s *Subscriber) AddTopic(topic string) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	s.topics[topic] = true
}

func (s *Subscriber) RemoveTopic(topic string) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	delete(s.topics, topic)
}

func (s *Subscriber) GetTopics() []string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var topics []string
	for topic := range s.topics {
		topics = append(topics, topic)
	}

	return topics
}

func (s *Subscriber) Signal(msg *Message) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if s.active {
		s.messages <- msg
	}
}

func (s *Subscriber) Destruct() {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	s.active = false
	close(s.messages)
}

func (s *Subscriber) Listen() {
	for {
		if msg, ok := <-s.messages; ok {
			fmt.Printf("Subscriber %s, received: %s from topic: %s\n", s.id, msg.GetBody(), msg.GetTopic())
		}
	}
}
