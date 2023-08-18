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
	Left            = iota
	Up
	Right
	Down
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
	direction      = Right
	loopCount      int
	speed          = 25
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

func update() error {
	// Handle arrow key events to change the snake's direction
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) && direction != Right {
		direction = Left
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) && direction != Left {
		direction = Right
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) && direction != Down {
		direction = Up
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) && direction != Up {
		direction = Down
	}

	// Increment the loop counter
	loopCount++

	// If the loop count reaches the desired speed, update the snake's position
	if loopCount >= speed {
		loopCount = 0 // Reset the loop count

		// Move the body
		if len(snake.Body) > 0 {
			// Save the old position of the head
			oldX, oldY := snakeX, snakeY

			// Update the position of the head based on the direction
			switch direction {
			case Left:
				snakeX -= itemSize
			case Right:
				snakeX += itemSize
			case Up:
				snakeY -= itemSize
			case Down:
				snakeY += itemSize
			}

			// Move the rest of the body
			for i := 0; i < len(snake.Body); i++ {
				// Save the current position of the body part
				currentX, currentY := snake.Body[i].X, snake.Body[i].Y

				// Move the body part to the previous position of the part in front of it
				snake.Body[i].X, snake.Body[i].Y = oldX, oldY

				// Update the old position to the current position for the next iteration
				oldX, oldY = currentX, currentY
			}
		}
		handleFoodCollisions()
		handleBodyCollisions()
		handleEdgeCollisions()
	}

	return nil
}

func handleFoodCollisions() {
	foodsMutex.Lock()
	defer foodsMutex.Unlock()

	for i := len(foods) - 1; i >= 0; i-- {
		food := foods[i]
		if (snakeX < food.X+itemSize) && (snakeX+itemSize > food.X) && (snakeY < food.Y+itemSize) && (snakeY+itemSize > food.Y) {
			foods = append(foods[:i], foods[i+1:]...) // Remove the eaten food
			addBodyPart()
		}
	}
}

func handleBodyCollisions() {
	for i := len(snake.Body) - 1; i >= 0; i-- {
		part := snake.Body[i]
		if (snakeX < part.X+itemSize) && (snakeX+itemSize > part.X) && (snakeY < part.Y+itemSize) && (snakeY+itemSize > part.Y) {
			//TODO: GAME OVER
			fmt.Println("GAME OVER")
			break
		}
	}
}

func handleEdgeCollisions() {
	if (snakeX < 0) || (snakeX+itemSize > screenWidth) || (snakeY < 0) || (snakeY+itemSize > screenHeight) {
		//TODO: GAME OVER
		fmt.Println("GAME OVER")
	}
}

func addBodyPart() {
	snake.BodySize++

	c := 200 - uint8(len(snake.Body))*20
	newBodyPart := SnakeBodyPart{
		X:     snake.Body[len(snake.Body)-1].X,
		Y:     snake.Body[len(snake.Body)-1].Y,
		Color: color.RGBA{c, c, c, 255},
	}

	snake.Body = append(snake.Body, newBodyPart)
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
	update()
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
