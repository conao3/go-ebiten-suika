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

type Rect struct {
	x float32
	y float32
	w float32
	h float32
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
var binRect = Rect{
	x: (width/2)-(bin_w/2),
	y: 100,
	w: bin_w,
	h: bin_h,
}

func PickRandomRadius() float32 {
	inx := rand.Intn(len(lstR))
	return lstR[inx]
}

func isCollisionBall(ball1, ball2 Circle) bool {
	distX := float64(ball1.x - ball2.x)
	distY := float64(ball1.y - ball2.y)
	distance := math.Sqrt(math.Pow(distX, 2) + math.Pow(distY, 2))
	if float32(distance) < ball1.r + ball2.r {
		return true
	}
	return false
}

func isCollisionBallBin(ball Circle, binRect Rect) bool {
	if ball.y + ball.r > binRect.y + binRect.h {
		return true
	}
	return false
}

func isCollisionFieldBall(ball Circle, inx int, fieldBalls []Circle) bool {
	for i := range fieldBalls {
		if i != inx {
			if isCollisionBall(ball, fieldBalls[i]) {
				return true
			}
		}
	}
	return false
}

func Colision(ball Circle, binRect Rect, fieldBalls []Circle) bool {
	if isCollisionBallBin(ball, binRect) {
		return true
	}
	for _, elm := range fieldBalls {
		if isCollisionBall(ball, elm) {
			return true
		}
	}
	return false
}

func UpdateFieldBalls(binRect Rect, fieldBalls *[]Circle) {
	// drop ball if no collision
	for i := range *fieldBalls {
		attemptX := (*fieldBalls)[i].x
		attemptY := (*fieldBalls)[i].y + 1
		if !(
			isCollisionFieldBall(Circle{attemptX, attemptY, (*fieldBalls)[i].r}, i, *fieldBalls) ||
			isCollisionBallBin(Circle{attemptX, attemptY, (*fieldBalls)[i].r}, binRect)) {
			(*fieldBalls)[i].y = attemptY
		}
	}

	// ensure no collision between field balls
	for i := range *fieldBalls {
		for j := range *fieldBalls {
			if i != j {
				distX := float64((*fieldBalls)[j].x - (*fieldBalls)[i].x)
				distY := float64((*fieldBalls)[j].y - (*fieldBalls)[i].y)
				distance := math.Sqrt(math.Pow(distX, 2) + math.Pow(distY, 2))
				expectedDistance := float64((*fieldBalls)[i].r + (*fieldBalls)[j].r)
				if float32(distance) < float32(expectedDistance) {

					(*fieldBalls)[i].x += float32(distX) / float32(distance) * float32(expectedDistance - distance)
					(*fieldBalls)[i].y += float32(distY) / float32(distance) * float32(expectedDistance - distance)
				}
			}
		}
	}

	// ensure no collision between field balls and bin
	for i := range *fieldBalls {
		left := (*fieldBalls)[i].x - (*fieldBalls)[i].r
		right := (*fieldBalls)[i].x + (*fieldBalls)[i].r
		bottom := (*fieldBalls)[i].y + (*fieldBalls)[i].r
		if left < binRect.x {
			(*fieldBalls)[i].x += binRect.x - left
		}
		if right > binRect.x + binRect.w {
			(*fieldBalls)[i].x -= right - (binRect.x + binRect.w)
		}
		if bottom > binRect.y + binRect.h {
			(*fieldBalls)[i].y -= bottom - (binRect.y + binRect.h)
		}
	}
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

		if Colision(g.ball, binRect, g.fieldBalls) {
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
	UpdateFieldBalls(binRect, &g.fieldBalls)
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
	vector.StrokeRect(screen, binRect.x, binRect.y, binRect.w, binRect.h, 1, color.RGBA{0, 255, 0, 255}, false)

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
