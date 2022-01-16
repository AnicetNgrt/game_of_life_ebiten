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

func NewGame(width int, height int, alivePerc int) *Game {
	world := make([][]bool, height)
	for y := range world {
		world[y] = make([]bool, width)
		for x := range world[y] {
			world[y][x] = rand.Intn(100) < alivePerc
		}
	}
	return &Game{world, width, height, time.Now()}
}

func (g *Game) LivingNeighbors(x int, y int) int {
	mod := func(a, b int) int {
		return (a%b + b) % b
	}

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
	return livingNeighbors
}

func (g *Game) ComputeLine(y int) []bool {
	line := make([]bool, g.width)
	for x, alive := range g.world[y] {
		livingNeighbors := g.LivingNeighbors(x, y)
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

	nextWorld := make([][]bool, g.height)
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
	const ALIVE_PERC = 70

	ebiten.SetWindowSize(640, 640)
	ebiten.SetWindowTitle("Conway's game of life")

	if err := ebiten.RunGame(NewGame(120*SIZE, 120*SIZE, ALIVE_PERC)); err != nil {
		log.Fatal(err)
	}
}
