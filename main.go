package main

import (
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	world      [][]bool
	width      int
	height     int
	lastUpdate time.Time
}

type Chunk struct {
	y    int
	data []bool
}

func NewGame(width int, height int) *Game {
	rawData := make([]bool, width*height)
	for i := range rawData {
		rawData[i] = rand.Intn(100) >= 30
	}

	world := make([][]bool, height)
	for i := range world {
		world[i], rawData = rawData[:width], rawData[width:]
	}

	return &Game{world, width, height, time.Now()}
}

func mod(a, b int) int {
	return (a%b + b) % b
}

func (g *Game) ComputeLine(y int) []bool {
	line := make([]bool, g.width)
	for x, alive := range g.world[y] {
		livingNeighbors := 0
		for sY := -1; sY <= 1; sY++ {
			for sX := -1; sX <= 1; sX++ {
				if sX == 0 && sY == 0 {
					continue
				}
				if g.world[mod(y+sY, g.height)][mod(x+sX, g.width)] {
					livingNeighbors++
				}
			}
		}

		if alive && (livingNeighbors == 2 || livingNeighbors == 3) {
			line[x] = true
		} else if !alive && livingNeighbors == 3 {
			line[x] = true
		} else if alive {
			line[x] = false
		}
	}
	return line
}

func (g *Game) ComputeNextWorld() [][]bool {
	nextWorld := make([][]bool, g.height)

	chunkChan := make(chan Chunk)

	computeChunk := func(y int) {
		chunk := Chunk{
			y,
			g.ComputeLine(y),
		}
		chunkChan <- chunk
	}

	for y := range g.world {
		go computeChunk(y)
	}

	for range nextWorld {
		chunk := <-chunkChan
		nextWorld[chunk.y] = chunk.data
	}

	return nextWorld
}

func (g *Game) Update() error {
	if time.Since(g.lastUpdate).Milliseconds() < 1000/60 {
		return nil
	}
	g.lastUpdate = time.Now()

	g.world = g.ComputeNextWorld()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for y, line := range g.world {
		for x, alive := range line {
			c := color.Black
			if alive {
				c = color.White
			}
			screen.Set(x, y, c)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.width, g.height
}

func main() {
	const SIZE = 3

	ebiten.SetWindowSize(800, 800)
	ebiten.SetWindowTitle("Conway's game of life")
	if err := ebiten.RunGame(NewGame(120*SIZE, 120*SIZE)); err != nil {
		log.Fatal(err)
	}
}
