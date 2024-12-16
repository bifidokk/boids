package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"log"
	"math"
	"sync"
)

const (
	screenWidth, screenHeight = 320, 320
	trainLength               = 70
)

var (
	trains        [4]*Train
	intersections [4]*Intersection
	colours       = [4]color.RGBA{
		{233, 33, 40, 255},
		{78, 151, 210, 255},
		{251, 170, 26, 255},
		{11, 132, 54, 255},
	}

	white = color.RGBA{R: 185, G: 185, B: 185, A: 255}
)

type Game struct{}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Layout(_, _ int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) Draw(screen *ebiten.Image) {
	DrawTrains(screen)
	DrawTracks(screen)
	DrawIntersections(screen)
}

func main() {
	for i := 0; i < 4; i++ {
		trains[i] = &Train{
			Id:          i,
			TrainLength: trainLength,
			Front:       0,
		}
	}

	for i := 0; i < 4; i++ {
		intersections[i] = &Intersection{
			Id:       i,
			Mutex:    sync.Mutex{},
			LockedBy: -1,
		}
	}

	go MoveTrain2(trains[0], 300, []*Crossing{{Position: 125, Intersection: intersections[0]},
		{Position: 175, Intersection: intersections[1]},
	})

	go MoveTrain2(trains[1], 300, []*Crossing{{Position: 125, Intersection: intersections[1]},
		{Position: 175, Intersection: intersections[2]},
	})

	go MoveTrain2(trains[2], 300, []*Crossing{{Position: 125, Intersection: intersections[2]},
		{Position: 175, Intersection: intersections[3]},
	})

	go MoveTrain2(trains[3], 300, []*Crossing{{Position: 125, Intersection: intersections[3]},
		{Position: 175, Intersection: intersections[0]},
	})

	ebiten.SetWindowSize(320*3, 320*3)
	ebiten.SetWindowTitle("Trains in a box")

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}

func DrawIntersections(screen *ebiten.Image) {
	drawIntersection(screen, intersections[0], 145, 145)
	drawIntersection(screen, intersections[1], 175, 145)
	drawIntersection(screen, intersections[2], 175, 175)
	drawIntersection(screen, intersections[3], 145, 175)
}

func DrawTracks(screen *ebiten.Image) {
	for i := 0; i < 300; i++ {
		screen.Set(10+i, 135, white)
		screen.Set(185, 10+i, white)
		screen.Set(310-i, 185, white)
		screen.Set(135, 310-i, white)
	}
}

func DrawTrains(screen *ebiten.Image) {
	drawXTrain(screen, 0, 1, 10, 135)
	drawYTrain(screen, 1, 1, 10, 185)
	drawXTrain(screen, 2, -1, 310, 185)
	drawYTrain(screen, 3, -1, 310, 135)
}

func drawIntersection(screen *ebiten.Image, intersection *Intersection, x int, y int) {
	c := white
	if intersection.LockedBy >= 0 {
		c = colours[intersection.LockedBy]
	}
	screen.Set(x-1, y, c)
	screen.Set(x, y-1, c)
	screen.Set(x, y, c)
	screen.Set(x+1, y, c)
	screen.Set(x, y+1, c)
}

func drawXTrain(screen *ebiten.Image, id int, dir int, start int, yPos int) {
	s := start + (dir * (trains[id].Front - trains[id].TrainLength))
	e := start + (dir * trains[id].Front)
	for i := math.Min(float64(s), float64(e)); i <= math.Max(float64(s), float64(e)); i++ {
		screen.Set(int(i)-dir, yPos-1, colours[id])
		screen.Set(int(i), yPos, colours[id])
		screen.Set(int(i)-dir, yPos+1, colours[id])
	}
}

func drawYTrain(screen *ebiten.Image, id int, dir int, start int, xPos int) {
	s := start + (dir * (trains[id].Front - trains[id].TrainLength))
	e := start + (dir * trains[id].Front)
	for i := math.Min(float64(s), float64(e)); i <= math.Max(float64(s), float64(e)); i++ {
		screen.Set(xPos-1, int(i)-dir, colours[id])
		screen.Set(xPos, int(i), colours[id])
		screen.Set(xPos+1, int(i)-dir, colours[id])
	}
}
