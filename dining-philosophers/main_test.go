package main

import (
	"testing"
	"time"
)

func Test_dine(t *testing.T) {
	for i := 0; i < 10; i++ {
		philosophersOrderFinished = []string{}
		dine()

		if len(philosophersOrderFinished) != 5 {
			t.Errorf("philosophersOrderFinished should be 5 but got %v", len(philosophersOrderFinished))
		}
	}
}

func Test_dineWithVaryingDelays(t *testing.T) {
	var tests = []struct {
		name  string
		delay time.Duration
	}{
		{"zero delay", time.Second * 0},
		{"250ms delay", time.Millisecond * 250},
		{"500ms delay", time.Millisecond * 500},
	}

	for _, test := range tests {
		philosophersOrderFinished = []string{}

		eatTime = test.delay
		idleTime = test.delay
		sleepTime = test.delay

		dine()

		if len(philosophersOrderFinished) != 5 {
			t.Errorf("philosophersOrderFinished should be 5 but got %v", len(philosophersOrderFinished))
		}
	}
}
