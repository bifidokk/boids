package main

import (
	"math"
	"math/rand"
	"time"
)

type Boid struct {
	position Vector2D
	velocity Vector2D
	id       int
}

func (b *Boid) start() {
	for {
		b.moveOne()
		time.Sleep(5 * time.Millisecond)
	}
}

func (b *Boid) moveOne() {
	b.velocity = b.velocity.Add(b.calculateAcceleration()).Limit(-1, 1)

	boidMap[int(b.position.x)][int(b.position.y)] = -1
	b.position = b.position.Add(b.velocity)
	boidMap[int(b.position.x)][int(b.position.y)] = b.id

	next := b.position.Add(b.velocity)

	if next.x >= screenWidth || next.x < 0 {
		b.velocity = Vector2D{-b.velocity.x, b.velocity.y}
	}

	if next.y >= screenHeight || next.y < 0 {
		b.velocity = Vector2D{b.velocity.x, -b.velocity.y}
	}
}

func (b *Boid) calculateAcceleration() Vector2D {
	upper, lower := b.position.AddV(viewRadius), b.position.AddV(-viewRadius)
	averageVelocity := Vector2D{0, 0}
	count := 0

	for i := math.Max(lower.x, 0); i <= math.Min(upper.x, screenWidth); i++ {
		for j := math.Max(lower.y, 0); j <= math.Min(upper.y, screenHeight); j++ {
			if boidId := boidMap[int(i)][int(j)]; boidId != -1 && boidId != b.id {
				if distance := boids[boidId].position.Distance(b.position); distance < viewRadius {
					count++
					averageVelocity = averageVelocity.Add(boids[boidId].velocity)
				}
			}
		}
	}

	acceleration := Vector2D{0, 0}

	if count > 0 {
		averageVelocity = averageVelocity.DivideV(float64(count))
		acceleration = averageVelocity.Subtract(b.velocity).MultiplyV(adjustmentRate)
	}

	return acceleration
}

func createBoid(bid int) {
	b := Boid{
		position: Vector2D{rand.Float64() * screenWidth, rand.Float64() * screenHeight},
		velocity: Vector2D{(rand.Float64() * 2) - 1.0, (rand.Float64() * 2) - 1.0},
		id:       bid,
	}

	boids[bid] = &b
	boidMap[int(b.position.x)][int(b.position.y)] = b.id

	go b.start()
}
