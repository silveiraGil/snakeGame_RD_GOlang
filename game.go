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
	maxFoods        = 4
	itemSize        = 15
	Left            = iota
	Up
	Right
	Down
)

type game struct{}

type Food struct {
	X, Y     int
	Color    color.Color
	name     string
	Creation time.Time
	Duration time.Duration
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
	foods                                           []*Food
	foodsMutex                                      sync.Mutex // Mutex to protect the foods slice
	snake                                           Snake
	direction                                       = Right
	snakeX, snakeY, loopCount, totalScore           int
	speed                                           = 25
	redFruits, greenFruits, blueFruits, whiteFruits int
	isGameOver                                      bool
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
	var name string
	maxSeconds := 15
	minSeconds := 5
	randomSeconds := rand.Intn(maxSeconds-minSeconds+1) + minSeconds
	lifetime := time.Duration(randomSeconds) * time.Second

	switch rand.Intn(4) {
	case 0:
		foodColor = color.RGBA{255, 0, 0, 255}
		name = "Red"
	case 1:
		foodColor = color.RGBA{0, 0, 255, 255}
		name = "Blue"
	case 2:
		foodColor = color.RGBA{255, 255, 255, 255}
		name = "White"
	case 3:
		foodColor = color.RGBA{0, 255, 0, 255}
		name = "Green"
	}

	food := &Food{
		X:        x,
		Y:        y,
		Color:    foodColor,
		name:     name,
		Creation: time.Now(),
		Duration: lifetime,
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
		bodyColor := part.Color
		vector.DrawFilledRect(screen, float32(part.X), float32(part.Y), float32(itemSize), float32(itemSize), bodyColor, true)
	}
}

func update() error {
	if isGameOver {
		return nil
	}

	expireFood()
	handleArrowKeyEvents()

	// Increment the loop counter
	loopCount++
	// If the loop count reaches the desired speed, update the snake's position
	if loopCount >= speed {
		loopCount = 0
		moveSnake()
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
			setScore(food.name)
			foods = append(foods[:i], foods[i+1:]...) // Remove the eaten food
			addBodyPart(food)
		}
	}
}

func handleBodyCollisions() {
	for i := len(snake.Body) - 1; i >= 0; i-- {
		part := snake.Body[i]
		if (snakeX < part.X+itemSize) && (snakeX+itemSize > part.X) && (snakeY < part.Y+itemSize) && (snakeY+itemSize > part.Y) {
			isGameOver = true
			break
		}
	}
}

func handleEdgeCollisions() {
	if (snakeX < 0) || (snakeX+itemSize > screenWidth) || (snakeY < 0) || (snakeY+itemSize > screenHeight) {
		isGameOver = true
	}
}

func addBodyPart(food *Food) {
	snake.BodySize++

	newBodyPart := SnakeBodyPart{
		X:     snake.Body[len(snake.Body)-1].X,
		Y:     snake.Body[len(snake.Body)-1].Y,
		Color: food.Color,
	}

	snake.Body = append(snake.Body, newBodyPart)
}

func updateScore() {
	// Update the game title with the fruit counts and total score
	title := fmt.Sprintf("Snake Game      |      Red: %d, Green: %d, Blue: %d, White: %d, Score: %d", redFruits, greenFruits, blueFruits, whiteFruits, totalScore)
	ebiten.SetWindowTitle(title)
}

func setScore(foodName string) {
	var score int
	switch foodName {
	case "Red":
		redFruits++
		score = rand.Intn(5)
	case "Blue":
		blueFruits++
		score = rand.Intn(4)
	case "White":
		whiteFruits++
		score = rand.Intn(3)
	case "Green":
		greenFruits++
		score = rand.Intn(2)
	}

	totalScore = totalScore + score

	if speed <= 6 {
		speed = speed + score
	} else {
		speed = speed - score
	}

	updateScore()
}

func expireFood() {
	foodsMutex.Lock()
	defer foodsMutex.Unlock()
	if len(foods) > 0 {
		for i := range foods {
			if time.Since(foods[i].Creation) > foods[i].Duration {
				foods = append(foods[:i], foods[i+1:]...)
				break
			}
		}
	}
}

func handleArrowKeyEvents() {
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
}

func moveSnake() {
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
}

func (g *game) Draw(screen *ebiten.Image) {
	if isGameOver {
		title := fmt.Sprintf("GAME OVER      |      Red: %d, Green: %d, Blue: %d, White: %d, Score: %d", redFruits, greenFruits, blueFruits, whiteFruits, totalScore)
		ebiten.SetWindowTitle(title)
	} else {
		setScreenColor(screen)
		drawFood(screen)
		drawSnake(screen)
	}
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
	updateScore()

	// Start the goroutine to create food
	go createFoodRoutine()

	// Initialize the snake
	initializeSnake()

	// Start the game
	if err := ebiten.RunGame(&game{}); err != nil {
		fmt.Println("Error:", err)
	}
}
