package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

type game struct{}

func setScreenColor(screen *ebiten.Image) {
	screen.Fill(color.Black)
}

func (g *game) Draw(screen *ebiten.Image) {
	setScreenColor(screen)
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

	// Start the game
	if err := ebiten.RunGame(&game{}); err != nil {
		fmt.Println("Error:", err)
	}
}
