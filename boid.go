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

func (b *Boid) moveOne() {
	acc := b.accelerate()
	rwLock.Lock()
	b.velocity = b.velocity.Add(acc).Limit(-1, 1)
	boidMap[int(b.position.x)][int(b.position.y)] = -1
	b.position = b.position.Add(b.velocity)
	boidMap[int(b.position.x)][int(b.position.y)] = b.id
	rwLock.Unlock()

	next := b.position.Add(b.velocity)
	if next.x >= screenWidth || next.x < 0 {
		b.velocity = Vector2D{-b.velocity.x, b.velocity.y}
	}
	if next.y >= screenHeight || next.y < 0 {
		b.velocity = Vector2D{b.velocity.x, -b.velocity.y}
	}
}

func (b *Boid) accelerate() Vector2D {
	upper, lower := b.position.AddV(viewRadius), b.position.AddV(-viewRadius)
	avgVelocity := Vector2D{0, 0}
	avgPosition := Vector2D{0, 0}
	separation := Vector2D{0, 0}
	count := 0.0

	rwLock.RLock() // NOTE: The mutexes in GO are not re-entrant.
	for i := math.Max(lower.x, 0); i < math.Min(upper.x, screenWidth); i++ {
		for j := math.Max(lower.y, 0); j < math.Min(upper.y, screenHeight); j++ {
			if otherBoid := boidMap[int(i)][int(j)]; otherBoid != -1 && otherBoid != b.id {
				if dist := boids[otherBoid].position.Distance(b.position); dist < viewRadius {
					count++
					avgVelocity = avgVelocity.Add(boids[otherBoid].velocity)
					avgPosition = avgPosition.Add(boids[otherBoid].position)
					separation = separation.Add(b.position.Subtract(boids[otherBoid].position).DivisionV(dist))
				}
			}
		}
	}
	rwLock.RUnlock()

	acc := Vector2D{b.borderBounce(b.position.x, screenWidth), b.borderBounce(b.position.y, screenHeight)}
	if count > 0 {
		avgPosition, avgVelocity = avgPosition.DivisionV(count), avgVelocity.DivisionV(count)
		accAlignment := avgVelocity.Subtract(b.velocity).MultiplyV(adjustionRate)
		accCohesion := avgPosition.Subtract(b.position).MultiplyV(adjustionRate)
		accSeparation := separation.MultiplyV(adjustionRate)
		acc = accAlignment.Add(accCohesion).Add(accSeparation)
	}

	return acc
}

func (b *Boid) borderBounce(position, max float64) float64 {
	if position < viewRadius {
		return 1 / position
	} else if position > max-viewRadius {
		return 1 / (position - max)
	}
	return 0
}

func (b *Boid) start() {
	for {
		b.moveOne()
		time.Sleep(5 * time.Millisecond)
	}
}

func createBoid(bid int) {
	b := Boid{
		position: Vector2D{x: rand.Float64() * screenWidth, y: rand.Float64() * screenHeight},
		velocity: Vector2D{x: rand.Float64()*2 - 1.0, y: rand.Float64()*2 - 1.0},
		id:       bid,
	}
	boids[bid] = &b
	boidMap[int(b.position.x)][int(b.position.y)] = bid
	go b.start()
}
