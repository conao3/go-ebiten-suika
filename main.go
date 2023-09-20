package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct{
	ballY float32
	ballSpeed float32
	dropping bool
}

func NewGame() *Game {
	return &Game{
		ballY: 50,
		ballSpeed: 0,
		dropping: false,
	}
}

func (g *Game) Update() error {
	if g.dropping {
		g.ballY += g.ballSpeed
		g.ballSpeed += 0.2

		if g.ballY > 400 {
			g.ballY = 50
			g.ballSpeed = 0
			g.dropping = false
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	yText := fmt.Sprintf("BallY: %.2f, Dropping: %t", g.ballY, g.dropping)
	ebitenutil.DebugPrint(screen, yText)

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.dropping = true
	}

	var width float32 = 600.0
	// var height float32 = 600.0

	var bin_w float32 = 250.0
	var bin_h float32 = 300.0

	// stage
	vector.StrokeRect(screen, (width/2)-(bin_w/2), 100, bin_w, bin_h, 1, color.RGBA{0, 255, 0, 255}, false)

	// ball
	vector.StrokeCircle(screen, (width/2), g.ballY, 5, 1, color.RGBA{255, 0, 0, 255}, false)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 600, 450
}

func main() {
	ebiten.SetWindowSize(1000, 750)
	ebiten.SetWindowTitle("Hello, World!")
	game := NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
