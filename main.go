package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten"
)

const (
	screenWidth, screenHeight = 640, 360
	boidCount                 = 300
)

var (
	green = color.RGBA{10, 255, 50, 255}
	boids [boidCount]*Boid
)

func update(screen *ebiten.Image) error {
	if !ebiten.IsDrawingSkipped() {
		for _, boid := range boids {
			screen.Set(int(boid.position.x+1), int(boid.position.y), green)
			screen.Set(int(boid.position.x-1), int(boid.position.y), green)
			screen.Set(int(boid.position.x), int(boid.position.y+1), green)
			screen.Set(int(boid.position.x), int(boid.position.y-1), green)
		}
	}
	return nil
}

func main() {
	for i := 0; i < boidCount; i++ {
		createBoid(i)
	}
	if err := ebiten.Run(update, screenWidth, screenHeight, 2, "Boids in a box"); err != nil {
		log.Fatal(err)
	}
}
