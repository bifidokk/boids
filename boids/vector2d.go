package main

import "math"

type Vector2D struct {
	x float64
	y float64
}

func (v1 Vector2D) Add(v2 Vector2D) Vector2D {
	return Vector2D{v1.x + v2.x, v1.y + v2.y}
}

func (v1 Vector2D) Subtract(v2 Vector2D) Vector2D {
	return Vector2D{v1.x - v2.x, v1.y - v2.y}
}

func (v1 Vector2D) Multiply(v2 Vector2D) Vector2D {
	return Vector2D{v1.x * v2.x, v1.y * v2.y}
}

func (v1 Vector2D) AddV(d float64) Vector2D {
	return Vector2D{v1.x + d, v1.y + d}
}

func (v1 Vector2D) MultiplyV(d float64) Vector2D {
	return Vector2D{v1.x * d, v1.y * d}
}

func (v1 Vector2D) DivideV(d float64) Vector2D {
	return Vector2D{v1.x / d, v1.y / d}
}

func (v1 Vector2D) Limit(lower, upper float64) Vector2D {
	return Vector2D{
		math.Min(math.Max(v1.x, lower), upper),
		math.Min(math.Max(v1.y, lower), upper),
	}
}

func (v1 Vector2D) Distance(v2 Vector2D) float64 {
	return math.Pow(v1.x-v2.x, 2) + math.Pow(v1.y-v2.y, 2)
}
