package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"os"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/golang/freetype/truetype"
)

var (
	BG_COLOR     = color.RGBA{0, 43, 54, 255}
	PADDLE_COLOR = color.RGBA{203, 75, 22, 255}
	BALL_COLOR   = color.RGBA{211, 1, 2, 255}
	LINE_COLOR   = color.RGBA{7, 54, 66, 255}
	SCORE_COLOR  = color.RGBA{42, 161, 152, 255}

	UP    Vector2 = Vector2{0, 1}
	DOWN  Vector2 = Vector2{0, -1}
	LEFT  Vector2 = Vector2{-1, 0}
	RIGHT Vector2 = Vector2{1, 0}

	PLAYER_SCORE   = 0
	OPPONENT_SCORE = 0

	GAME_OVER = false
	WINNER    = 1
)

const (
	WIDTH  float64 = 1366
	HEIGHT float64 = 768

	SPEED_MULTIPLIER = 0.1
)

type Paddle struct {
	Position      Vector2
	MoveDirection Vector2
	Speed         float64
	Color         color.Color
	Shape         *imdraw.IMDraw
}

type Ball struct {
	Position      Vector2
	MoveDirection Vector2
	Speed         float64
	Color         color.Color
	Shape         *imdraw.IMDraw
}

func NewBall() Ball {
	ball := Ball{
		Position:      Vector2{WIDTH / 2, HEIGHT / 2},
		MoveDirection: LEFT,
		Speed:         7,
		Color:         BALL_COLOR,
	}

	ball.AddShape()
	return ball
}

func Point(self bool, player, opponent *Paddle, ball *Ball) {
	ball.Position = Vector2{WIDTH / 2, HEIGHT / 2}
	ball.Speed = 7

	player.Position = Vector2{25, HEIGHT / 2}
	player.MoveDirection = Vector2{0, 0}
	player.Speed = 15

	opponent.Position = Vector2{WIDTH - 25, HEIGHT / 2}
	opponent.MoveDirection = Vector2{0, 0}
	opponent.Speed = 0

	if self {
		PLAYER_SCORE++
		ball.MoveDirection = LEFT
	} else {
		OPPONENT_SCORE++
		ball.MoveDirection = RIGHT
	}

	if PLAYER_SCORE >= 10 || OPPONENT_SCORE >= 10 {
		GAME_OVER = true
		if self {
			WINNER = 1
		} else {
			WINNER = 2
		}
	}
}

func (b *Ball) AddShape() {
	b.Shape = MakeRectangle(b.Position, b.Color, 25, 25)
}

func (b *Ball) Move(player, opponent Paddle) {
	if VectorSum(b.Position, ScalarProduct(b.MoveDirection, b.Speed)).Y+(25/2) > HEIGHT {
		b.MoveDirection.FlipY()
	} else if VectorSum(b.Position, ScalarProduct(b.MoveDirection, b.Speed)).Y-(25/2) < 0 {
		b.MoveDirection.FlipY()
	} else if b.WillCollideWithPlayer(player) {
		b.MoveDirection.FlipX()
		b.MoveDirection.Add(Vector2{0, b.CollisionPoint(player)})
		b.Speed *= b.CollisionPoint(player)
		// writeToStdout(b.CollisionPoint(player))
	} else if b.WillCollideWithOpponent(opponent) {
		b.MoveDirection.FlipX()
		b.MoveDirection.Add(Vector2{0, b.CollisionPoint(opponent)})
		b.Speed *= b.CollisionPoint(opponent)
		// writeToStdout(b.CollisionPoint(opponent))
	} else if b.WillCollideWithPlayerWall() {
		Point(false, &player, &opponent, b)
	} else if b.WillCollideWithOpponentWall() {
		Point(true, &player, &opponent, b)
	}
	b.Position.Add(ScalarProduct(b.MoveDirection, b.Speed))
	b.Shape = MakeRectangle(b.Position, b.Color, 25, 25)
}

// func writeToStdout(value any) {
// 	text := fmt.Sprintf("%v", value)
// 	w := bufio.NewWriter(os.Stdout)
// 	w.Write([]byte(text))
// }

func (b *Ball) WillCollideWithPlayer(player Paddle) bool {
	target := VectorSum(b.Position, ScalarProduct(b.MoveDirection, b.Speed))

	return target.X-(25/2) < player.Position.X+(25/2) &&
		target.Y-(25/2) < player.Position.Y+100 &&
		target.Y+(25/2) > player.Position.Y-100
}

func (b *Ball) WillCollideWithPlayerWall() bool {
	target := VectorSum(b.Position, ScalarProduct(b.MoveDirection, b.Speed))
	return target.X-(25/2) <= 0
}

func (b *Ball) WillCollideWithOpponentWall() bool {
	target := VectorSum(b.Position, ScalarProduct(b.MoveDirection, b.Speed))
	return target.X+(25/2) >= WIDTH
}

func (b *Ball) CollisionPoint(paddle Paddle) float64 {
	p := (paddle.Position.Y - b.Position.Y)

	return -(((100 * p) / paddle.Position.Y) / 100)
}

func (b *Ball) WillCollideWithOpponent(opponent Paddle) bool {
	target := VectorSum(b.Position, ScalarProduct(b.MoveDirection, b.Speed))

	return target.X+(25/2) > opponent.Position.X-(25/2) &&
		target.Y-(25/2) < opponent.Position.Y+100 &&
		target.Y+(25/2) > opponent.Position.Y-100
}

type Vector2 struct {
	X float64
	Y float64
}

func (v *Vector2) Add(other Vector2) {
	v.X += other.X
	v.Y += other.Y
}

func (v *Vector2) FlipX() {
	v.X = -v.X
}

func (v *Vector2) FlipY() {
	v.Y = -v.Y
}

func VectorSum(one, two Vector2) Vector2 {
	return Vector2{
		X: one.X + two.X,
		Y: one.Y + two.Y,
	}
}

func VectorProduct(one, two Vector2) Vector2 {
	return Vector2{
		X: one.X * two.X,
		Y: one.Y * two.Y,
	}
}

func ScalarProduct(vector Vector2, scalar float64) Vector2 {
	return Vector2{
		X: vector.X * scalar,
		Y: vector.Y * scalar,
	}
}

type Vector4 struct {
	TL float64
	TR float64
	BR float64
	BL float64
}

func NewPaddle(player bool) Paddle {
	var paddle Paddle
	if player {
		paddle = Paddle{
			Position:      Vector2{25, HEIGHT / 2},
			MoveDirection: Vector2{0, 0},
			Speed:         15,
			Color:         PADDLE_COLOR,
		}
	} else {
		paddle = Paddle{
			Position:      Vector2{WIDTH - 25, HEIGHT / 2},
			MoveDirection: Vector2{0, 0},
			Speed:         15,
			Color:         PADDLE_COLOR,
		}
	}

	paddle.AddShape()
	return paddle
}

func (p *Paddle) AddShape() {
	p.Shape = MakeRectangle(p.Position, p.Color, 25, 200)
}

func (p *Paddle) Move(direction Vector2) {
	if VectorSum(p.Position, ScalarProduct(direction, p.Speed)).Y+100 > HEIGHT {
		return
	}
	if VectorSum(p.Position, ScalarProduct(direction, p.Speed)).Y-100 < 0 {
		return
	}
	p.Position.Add(ScalarProduct(direction, p.Speed))
	p.Shape = MakeRectangle(p.Position, p.Color, 25, 200)
}

func (p *Paddle) AutoMove(ball Ball) {

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)
	offset := r.Float64()
	y := (ball.Position.Y + (offset * 2))

	if y+100 > HEIGHT || y-100 < 0 {
		return
	}

	p.Position.Y = y
	p.Shape = MakeRectangle(p.Position, p.Color, 25, 200)
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Pong",
		Bounds: pixel.R(0, 0, WIDTH, HEIGHT),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	opp := NewPaddle(false)
	player := NewPaddle(true)
	ball := NewBall()
	line := MakeRectangle(Vector2{WIDTH / 2, HEIGHT / 2}, LINE_COLOR, 10, HEIGHT)

	ttf, err := os.ReadFile("assets/upheavtt.ttf")
	if err != nil {
		panic(err)
	}

	parsedTTF, err := truetype.Parse(ttf)
	if err != nil {
		panic(err)
	}

	gameFont := truetype.NewFace(parsedTTF, &truetype.Options{54, 0, 0, 0, 0, 0})

	atlas := text.NewAtlas(
		gameFont,
		[]rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'},
		[]rune{'G', 'A', 'M', 'E', 'O', 'V', 'R', 'P', 'L', 'Y', 'W', 'I', 'N', 'S', ' '},
	)
	playerScore := text.New(pixel.V(WIDTH/2-50-30, HEIGHT-50), atlas)
	playerScore.Color = SCORE_COLOR
	opponentScore := text.New(pixel.V(WIDTH/2+50, HEIGHT-50), atlas)
	opponentScore.Color = SCORE_COLOR

	endText := text.New(pixel.V(WIDTH/2, HEIGHT/2), atlas)
	endText.Color = SCORE_COLOR

	// fmt.Fprintf(endText, "GAME OVER\nPLAYER %d WINS", 1)

	for !win.Closed() {
		win.Clear(BG_COLOR)
		playerScore.Clear()
		opponentScore.Clear()

		fmt.Fprintf(playerScore, "%d", PLAYER_SCORE)
		fmt.Fprintf(opponentScore, "%d", OPPONENT_SCORE)

		if GAME_OVER {
			playerScore.Draw(win, pixel.IM)
			opponentScore.Draw(win, pixel.IM)

			lines := []string{
				"GAME OVER",
				fmt.Sprintf("PLAYER %d WINS", WINNER),
			}

			for _, line := range lines {
				endText.Dot.X -= endText.BoundsOf(line).W() / 2
				fmt.Fprintln(endText, line)
			}

			endText.Draw(win, pixel.IM)
			win.Update()
			for {

			}
		}

		line.Draw(win)

		opp.Shape.Draw(win)
		player.Shape.Draw(win)
		ball.Shape.Draw(win)

		playerScore.Draw(win, pixel.IM)
		opponentScore.Draw(win, pixel.IM)

		ball.Move(player, opp)

		if win.Pressed(pixelgl.KeyW) {
			player.Move(UP)
		} else if win.Pressed(pixelgl.KeyS) {
			player.Move(DOWN)
		}

		if win.Pressed(pixelgl.KeyUp) {
			opp.Move(UP)
		} else if win.Pressed(pixelgl.KeyDown) {
			opp.Move(DOWN)
		}

		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
