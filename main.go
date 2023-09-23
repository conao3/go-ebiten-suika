package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Vector struct {
	x float32
	y float32
}

type Ball struct {
	p Vector  // position
	v Vector  // velocity
	r float32  // radius
}

type Rect struct {
	p Vector  // position
	w float32  // width
	h float32  // height
}

type Game struct{
	ball Ball
	fieldBalls []Ball
	BallHorizontalV float32
}

const (
	width float32 = 400.0
	height float32 = 300.0

	bin_w float32 = 150.0
	bin_h float32 = 200.0
)

var lstR = []float32{
	5.0, 10.0, 15.0, 20.0, 25.0, 30.0,
}

var gravity = Vector{
	x: 0,
	y: 0.5,
}

var binRect = Rect{
	p: Vector {
		x: (width/2)-(bin_w/2),
		y: 80,
	},
	w: bin_w,
	h: bin_h,
}

func (v *Vector) Minus() {
	v.x *= -1
	v.y *= -1
}

func (v *Vector) Length() float32 {
	if v.x == 0 {
		return float32(math.Abs(float64(v.y)))
	} else if v.y == 0 {
		return float32(math.Abs(float64(v.x)))
	}

	return float32(math.Sqrt(float64(v.x*v.x + v.y*v.y)))
}

func (v *Vector) Angle() float32 {
	return float32(math.Atan2(float64(v.y), float64(v.x)))
}

func (v *Vector) Add(v2 Vector) {
	v.x += v2.x
	v.y += v2.y
}

func (v *Vector) Sub(v2 Vector) {
	v.Add(VectorMinus(v2))
}

func (v *Vector) Dot(v2 Vector) float32 {
	return v.x*v2.x + v.y*v2.y
}

func (v *Vector) Cross(v2 Vector) float32 {
	return v.x*v2.y - v.y*v2.x
}

func (v *Vector) AddScalar(v2 float32) {
	v.x += v2
	v.y += v2
}

func (v *Vector) SubScalar(v2 float32) {
	v.AddScalar(-v2)
}

func (v *Vector) MulScalar(v2 float32) {
	v.x *= v2
	v.y *= v2
}

func (v *Vector) DivScalar(v2 float32) {
	v.MulScalar(1/v2)
}

func (v *Vector) Normalize() {
	v.DivScalar(v.Length())
}

func VectorMinus(v Vector) Vector {
	return Vector{
		x: -v.x,
		y: -v.y,
	}
}

func VectorLength(v Vector) float32 {
	if v.x == 0 {
		return float32(math.Abs(float64(v.y)))
	} else if v.y == 0 {
		return float32(math.Abs(float64(v.x)))
	}

	return float32(math.Sqrt(float64(v.x*v.x + v.y*v.y)))
}

func VectorAngle(v Vector) float32 {
	return float32(math.Atan2(float64(v.y), float64(v.x)))
}

func VectorAdd(v1 Vector, v2 Vector) Vector {
	return Vector{
		x: v1.x + v2.x,
		y: v1.y + v2.y,
	}
}

func VectorSub(v1 Vector, v2 Vector) Vector {
	return VectorAdd(v1, VectorMinus(v2))
}

func VectorDot(v1 Vector, v2 Vector) float32 {
	return v1.x*v2.x + v1.y*v2.y
}

func VectorCross(v1 Vector, v2 Vector) float32 {
	return v1.x*v2.y - v1.y*v2.x
}

func VectorAddScalar(v Vector, v2 float32) Vector {
	return Vector{
		x: v.x + v2,
		y: v.y + v2,
	}
}

func VectorSubScalar(v Vector, v2 float32) Vector {
	return VectorAddScalar(v, -v2)
}

func VectorMulScalar(v Vector, v2 float32) Vector {
	return Vector{
		x: v.x * v2,
		y: v.y * v2,
	}
}

func VectorDivScalar(v Vector, v2 float32) Vector {
	return VectorMulScalar(v, 1/v2)
}

func VectorNormalize(v Vector) Vector {
	return VectorDivScalar(v, VectorLength(v))
}

func PickRandomRadius() float32 {
	inx := rand.Intn(len(lstR))
	return lstR[inx]
}

func NewBall() Ball {
	return Ball{
		p: Vector{
			x: width/2,
			y: 50,
		},
		v: Vector{
			x: 0,
			y: 0,
		},
		r: PickRandomRadius(),
	}
}

func NewGame() *Game {
	return &Game{
		ball: NewBall(),
	}
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustReleased(ebiten.KeySpace) {
		g.fieldBalls = append(g.fieldBalls, g.ball)
		g.ball = NewBall()
		g.ball.p.x = g.fieldBalls[len(g.fieldBalls)-1].p.x + (rand.Float32()-0.5)
	}

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		if g.BallHorizontalV > 0 {
			g.BallHorizontalV = 0
		}
		g.BallHorizontalV -= 0.3
		g.ball.p.x += g.BallHorizontalV
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		if g.BallHorizontalV < 0 {
			g.BallHorizontalV = 0
		}
		g.BallHorizontalV += 0.3
		g.ball.p.x += g.BallHorizontalV
	} else {
		g.BallHorizontalV = 0
	}

	for i := range g.fieldBalls {
		g.fieldBalls[i].v.Add(gravity)
	}
	for i := range g.fieldBalls {
		g.fieldBalls[i].p.Add(g.fieldBalls[i].v)
	}
	for i := range g.fieldBalls {
		for j := range g.fieldBalls {
			if i != j {
				d := VectorSub(g.fieldBalls[i].p, g.fieldBalls[j].p)
				dl := VectorLength(d)
				r := g.fieldBalls[i].r + g.fieldBalls[j].r
				if dl < r {
					diff := r - dl
					dn := VectorNormalize(d)
					g.fieldBalls[i].p.Add(VectorMulScalar(dn, diff))
					g.fieldBalls[j].p.Sub(VectorMulScalar(dn, diff))

					vDiff := VectorSub(g.fieldBalls[i].v, g.fieldBalls[j].v)
					A := VectorMulScalar(dn, VectorDot(vDiff, dn) * 0.25)
					g.fieldBalls[i].v.Sub(A)
					g.fieldBalls[j].v.Add(A)
				}
			}
		}
	}
	for i := range g.fieldBalls {
		bin_bottom := binRect.p.y + binRect.h
		if g.fieldBalls[i].p.y + g.fieldBalls[i].r > bin_bottom {
			g.fieldBalls[i].p.y = bin_bottom - g.fieldBalls[i].r
			g.fieldBalls[i].v.y *= -1 * 0.3
		}
	}
	for i := range g.fieldBalls {
		if g.fieldBalls[i].v.Length() < 0.1 {
			g.fieldBalls[i].v.x = 0
			g.fieldBalls[i].v.y = 0
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// if (len(g.fieldBalls) >= 1) {
	// 	yText := fmt.Sprintf("y: %f, length: %f", g.fieldBalls[0].p.y, g.fieldBalls[0].v.Length())
	// 	ebitenutil.DebugPrint(screen, yText)
	// }
	if (len(g.fieldBalls) >= 2) {
		yText := fmt.Sprintf("length: %f, r1+r2: %f", VectorLength(VectorSub(g.fieldBalls[0].p, g.fieldBalls[1].p)), g.fieldBalls[0].r + g.fieldBalls[1].r)
		ebitenutil.DebugPrint(screen, yText)
	}

	// stage
	vector.StrokeRect(screen, binRect.p.x, binRect.p.y, bin_w, bin_h, 1, color.RGBA{0, 255, 0, 255}, false)

	// ball
	vector.StrokeCircle(screen, g.ball.p.x, g.ball.p.y, g.ball.r, 1, color.RGBA{255, 0, 0, 255}, false)

	// field ball
	for _, ball := range g.fieldBalls {
		vector.StrokeCircle(screen, ball.p.x, ball.p.y, ball.r, 1, color.RGBA{0, 0, 255, 255}, false)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return int(width), int(height)
}

func main() {
	ebiten.SetWindowSize(1000, 750)
	ebiten.SetWindowTitle("Hello, World!")
	game := NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
