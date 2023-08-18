package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth     = 640
	screenHeight    = 480
	foodSpawnPeriod = 3 * time.Second
	maxFoods        = 3
	foodSize        = 15
)

type game struct{}

type Food struct {
	X, Y  int
	Color color.Color
}

var (
	foods      []*Food
	foodsMutex sync.Mutex // Mutex to protect the foods slice
)

func setScreenColor(screen *ebiten.Image) {
	screen.Fill(color.Black)
}

func createFoodRoutine() {
	for {
		if len(foods) < maxFoods {
			createFood()
		}
		time.Sleep(foodSpawnPeriod)
	}
}

func createFood() {
	foodsMutex.Lock()
	defer foodsMutex.Unlock()

	x := rand.Intn(screenWidth - foodSize)
	y := rand.Intn(screenHeight - foodSize)

	var foodColor color.Color

	switch rand.Intn(3) {
	case 0:
		foodColor = color.RGBA{255, 0, 0, 255}
	case 1:
		foodColor = color.RGBA{0, 0, 255, 255}
	case 2:
		foodColor = color.RGBA{255, 255, 255, 255}
	}

	food := &Food{
		X:     x,
		Y:     y,
		Color: foodColor,
	}

	foods = append(foods, food)
}

func drawFood(screen *ebiten.Image) {
	foodsMutex.Lock()
	defer foodsMutex.Unlock()

	for _, food := range foods {
		vector.DrawFilledRect(screen, float32(food.X), float32(food.Y), foodSize, foodSize, food.Color, true)
	}
}

func (g *game) Draw(screen *ebiten.Image) {
	setScreenColor(screen)
	drawFood(screen)
}

func (g *game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *game) Update() error {
	return nil
}

func main() {

	ebiten.SetWindowSize(screenWidth, screenHeight)
	title := "Snake Game"
	ebiten.SetWindowTitle(title)

	// Start the goroutine to create food
	go createFoodRoutine()

	// Start the game
	if err := ebiten.RunGame(&game{}); err != nil {
		fmt.Println("Error:", err)
	}
}
