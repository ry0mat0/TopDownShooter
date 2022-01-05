package main

import (
	"fmt"
	"image/color"
	_ "image/png"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	// "golang.org/x/exp/shiny/screen"
	"ebiten_hello/player"
)

const (
	pi      = math.Pi
	screenX = 640
	screenY = 480

	// game modes
	modeTitle    = 0
	modeGame     = 1
	modeGameover = 2

	//player
	speed       = 1.0
	accel       = 3.0
	brake       = 0.5
	rotSpeed    = 0.02
	playerSizeX = 100
	playerSizeY = 100
)

var gop_img *ebiten.Image

func init() {
	var err error
	gop_img, _, err = ebitenutil.NewImageFromFile("image/plane.png")
	if err != nil {
		log.Fatal(err)
	}
}

type Game struct {
	mode   int
	player *player.Player
}

// NewGame method
func NewGame() *Game {
	g := &Game{}
	g.init()
	return g
}

func (g *Game) init() {
	g.player = &player.Player{}
	g.player.NewPlayer(playerSizeX, playerSizeY, screenX, screenY, speed, accel, brake, rotSpeed)
}

func (g *Game) Update() error {
	switch g.mode {
	case modeTitle:
		if g.isKeyJustPressed() {
			g.mode = modeGame
			fmt.Println("mode Changed")
		}
	case modeGame:
		g.player.Move(g.moveKey(), screenX, screenY)
	case modeGameover:
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0xff, 0xff})
	switch g.mode {
	case modeTitle:
		ebitenutil.DebugPrint(screen, "Title")
	case modeGame:
		ebitenutil.DebugPrint(screen, "Game")
		g.DrawPlayer(screen)
	case modeGameover:
		ebitenutil.DebugPrint(screen, "GameOver")
	}
}

func (g *Game) DrawPlayer(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Rotate(g.player.Direction)
	op.GeoM.Translate(g.player.X, g.player.Y)
	op.GeoM.Scale(1.0, 1.0)
	s := strconv.FormatFloat(g.player.X, 'f', 1, 64)
	t := strconv.FormatFloat(g.player.Y, 'f', 1, 64)
	u := strconv.FormatFloat(g.player.Direction*180.0/pi, 'f', 1, 64)
	ebitenutil.DebugPrintAt(screen, strings.Join([]string{"X is ", s}, ""), 0, 15)
	ebitenutil.DebugPrintAt(screen, strings.Join([]string{"Y is ", t}, ""), 0, 30)
	ebitenutil.DebugPrintAt(screen, strings.Join([]string{"D is ", u}, ""), 0, 45)
	screen.DrawImage(gop_img, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return screenX, screenY
}

func (g *Game) isKeyJustPressed() bool {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		return true
	}
	return false
}

func (g *Game) moveKey() (keys [4]bool) {
	keys = [4]bool{
		ebiten.IsKeyPressed(ebiten.KeyW),
		ebiten.IsKeyPressed(ebiten.KeyA),
		ebiten.IsKeyPressed(ebiten.KeyS),
		ebiten.IsKeyPressed(ebiten.KeyD),
	}
	return
}

func main() {
	ebiten.SetWindowSize(screenX, screenY)
	ebiten.SetWindowTitle("EbitenGame")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
