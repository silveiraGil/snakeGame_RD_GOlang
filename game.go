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
	itemSize        = 15
)

type game struct{}

type Food struct {
	X, Y  int
	Color color.Color
}

type SnakeBodyPart struct {
	X, Y  int
	Color color.Color
}

type Snake struct {
	Head     SnakeBodyPart
	Body     []SnakeBodyPart
	BodySize int
}

var (
	foods          []*Food
	foodsMutex     sync.Mutex // Mutex to protect the foods slice
	snake          Snake
	snakeX, snakeY int
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

	x := rand.Intn(screenWidth - itemSize)
	y := rand.Intn(screenHeight - itemSize)

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

func initializeSnake() {
	// Set snake's head color
	snake.Head.Color = color.RGBA{77, 77, 77, 255}
	snake.BodySize = 2

	// Set snake's body colors
	snake.Body = make([]SnakeBodyPart, snake.BodySize) // Initialize the slice with the correct size
	for i := range snake.Body {
		// Set the color for each body part
		// The body parts closer to the head are lighter in color
		c := 200 - uint8(i)*20
		snake.Body[i].Color = color.RGBA{c, c, c, 255}
	}

	// Set the initial position for the snake's head (center of the screen)
	snakeX = screenWidth / 2
	snakeY = screenHeight / 2

	snake.Head.X = snakeX
	snake.Head.Y = snakeY

	// Set the initial positions for the snake's body parts
	for i := range snake.Body {
		snake.Body[i].X = snake.Head.X - (i+1)*itemSize
		snake.Body[i].Y = snake.Head.Y
	}
}

func drawFood(screen *ebiten.Image) {
	foodsMutex.Lock()
	defer foodsMutex.Unlock()

	for _, food := range foods {
		vector.DrawFilledRect(screen, float32(food.X), float32(food.Y), itemSize, itemSize, food.Color, true)
	}
}

func drawSnake(screen *ebiten.Image) {
	// Draw snake's head
	vector.DrawFilledRect(screen, float32(snakeX), float32(snakeY), float32(itemSize), float32(itemSize), snake.Head.Color, true)

	// Draw snake's body
	for _, part := range snake.Body {
		bodyColor := color.RGBA{127, 127, 127, 100}
		vector.DrawFilledRect(screen, float32(part.X), float32(part.Y), float32(itemSize), float32(itemSize), bodyColor, true)
	}
}

func (g *game) Draw(screen *ebiten.Image) {
	setScreenColor(screen)
	drawFood(screen)
	drawSnake(screen)
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

	// Initialize the snake
	initializeSnake()

	// Start the game
	if err := ebiten.RunGame(&game{}); err != nil {
		fmt.Println("Error:", err)
	}
}
