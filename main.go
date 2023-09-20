package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Circle struct {
	x float32
	y float32
	r float32
}

type Game struct{
	ball Circle
	ballSpeed float32
	dropping bool
	fieldBalls []Circle
}

const (
	width float32 = 600.0
	// var height float32 = 600.0

	bin_w float32 = 250.0
	bin_h float32 = 300.0
)

var lstR = []float32{
	5.0, 10.0, 15.0, 20.0, 25.0, 30.0,
}

func PickRandomRadius() float32 {
	inx := rand.Intn(len(lstR))
	return lstR[inx]
}

func Colision(ball Circle, bottom float32, fieldBalls []Circle) bool {
	if ball.y + ball.r > bottom {
		return true
	}
	for _, elm := range fieldBalls {
		distX := float64(ball.x - elm.x)
		distY := float64(ball.y - elm.y)
		distance := math.Sqrt(math.Pow(distX, 2) + math.Pow(distY, 2))
		if float32(distance) < ball.r + elm.r {
			return true
		}
	}
	return false
}

func NewCircle() Circle {
	return Circle{
		x: 300,
		y: 50,
		r: PickRandomRadius(),
	}
}

func NewGame() *Game {
	return &Game{
		ball: NewCircle(),
		ballSpeed: 0,
		dropping: false,
	}
}

func (g *Game) Update() error {
	if g.dropping {
		g.ball.y += g.ballSpeed
		g.ballSpeed += 0.2

		if Colision(g.ball, 100 + bin_h, g.fieldBalls) {
			g.fieldBalls = append(g.fieldBalls, Circle{
				x: g.ball.x,
				y: g.ball.y,
				r: g.ball.r,
			})
			g.ball = NewCircle()
			g.ballSpeed = 0
			g.dropping = false
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	yText := fmt.Sprintf("Ball.y: %.2f, Dropping: %t", g.ball.y, g.dropping)
	ebitenutil.DebugPrint(screen, yText)

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.dropping = true
	}

	if !g.dropping {
		if ebiten.IsKeyPressed(ebiten.KeyLeft) {
			g.ball.x -= 3
		} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
			g.ball.x += 3
		}
	}

	// stage
	vector.StrokeRect(screen, (width/2)-(bin_w/2), 100, bin_w, bin_h, 1, color.RGBA{0, 255, 0, 255}, false)

	// ball
	vector.StrokeCircle(screen, g.ball.x, g.ball.y, g.ball.r, 1, color.RGBA{255, 0, 0, 255}, false)

	// field ball
	for _, ball := range g.fieldBalls {
		vector.StrokeCircle(screen, ball.x, ball.y, ball.r, 1, color.RGBA{0, 0, 255, 255}, false)
	}
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
