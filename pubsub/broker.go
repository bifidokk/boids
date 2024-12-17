package main

import "sync"

type Subscribers map[string]*Subscriber

type Broker struct {
	subscribers Subscribers
	topics      map[string]Subscribers
	mutex       sync.RWMutex
}

func NewBroker() *Broker {
	return &Broker{
		subscribers: Subscribers{},
		topics:      map[string]Subscribers{},
	}
}

func (b *Broker) AddSubscriber() *Subscriber {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	id, s := NewSubscriber()
	b.subscribers[id] = s

	return s
}

func (b *Broker) Broadcast(topics []string, body string) {
	for _, topic := range topics {
		for _, s := range b.topics[topic] {
			message := NewMessage(topic, body)

			go (func(s *Subscriber) {
				s.Signal(message)
			})(s)
		}
	}
}

func (b *Broker) Subscribe(s *Subscriber, topic string) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if b.topics[topic] == nil {
		b.topics[topic] = Subscribers{}
	}

	s.AddTopic(topic)
	b.topics[topic][s.id] = s
}

func (b *Broker) Publish(topic string, body string) {
	b.mutex.RLock()
	topics := b.topics[topic]
	b.mutex.RUnlock()

	for _, s := range topics {
		if s.active == false {
			return
		}

		message := NewMessage(topic, body)

		go (func(s *Subscriber) {
			s.Signal(message)
		})(s)
	}
}

func (b *Broker) Unsubscribe(s *Subscriber, topic string) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	delete(b.topics[topic], s.id)
	s.RemoveTopic(topic)
}

func (b *Broker) RemoveSubscriber(s *Subscriber) {
	for topic := range s.topics {
		b.Unsubscribe(s, topic)
	}

	b.mutex.Lock()
	delete(b.subscribers, s.id)
	b.mutex.Unlock()
	s.Destruct()
}
