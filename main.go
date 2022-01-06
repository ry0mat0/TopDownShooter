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

	"TopDownShooter/bullet"
	// "TopDownShooter/controll"
	"TopDownShooter/player"
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
	speed       = 2.0
	accel       = 3.0
	brake       = 0.5
	rotSpeed    = 0.02
	playerSizeX = 50
	playerSizeY = 50

	//bullets
	maxBulletCount = 10
	bulletSizeX    = 2
	bulletSizeY    = 10
	bulletSpeed    = 5.0
)

var playerImg *ebiten.Image
var bulletImg *ebiten.Image

func init() {
	var err error
	playerImg, _, err = ebitenutil.NewImageFromFile("image/plane.png")
	bulletImg, _, err = ebitenutil.NewImageFromFile("image/bullet.png")
	if err != nil {
		log.Fatal(err)
	}
}

type Game struct {
	mode    int
	player  *player.Player
	bullets [maxBulletCount]*bullet.Bullet
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
	for i := 0; i < maxBulletCount; i++ {
		g.bullets[i] = &bullet.Bullet{}
	}
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
		for i := 0; i < maxBulletCount; i++ {
			g.bullets[i].Move(screenX, screenY)
		}
		if g.isKeyJustPressed() {
			i := checkEmptyBullet(g.bullets)
			if i >= 0 {
				g.bullets[i].NewBullet(g.player.X, g.player.Y, g.player.Direction, bulletSpeed)
			}
		}
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
		g.DrawBullets(screen)
	case modeGameover:
		ebitenutil.DebugPrint(screen, "GameOver")
	}
}

func (g *Game) DrawPlayer(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Rotate(g.player.Direction)
	op.GeoM.Translate(g.player.X, g.player.Y)
	s := strconv.FormatFloat(g.player.X, 'f', 1, 64)
	t := strconv.FormatFloat(g.player.Y, 'f', 1, 64)
	u := strconv.FormatFloat(g.player.Direction*180.0/pi, 'f', 1, 64)
	ebitenutil.DebugPrintAt(screen, strings.Join([]string{"X is ", s}, ""), 0, 15)
	ebitenutil.DebugPrintAt(screen, strings.Join([]string{"Y is ", t}, ""), 0, 30)
	ebitenutil.DebugPrintAt(screen, strings.Join([]string{"D is ", u}, ""), 0, 45)
	screen.DrawImage(playerImg, op)
}

func (g *Game) DrawBullets(screen *ebiten.Image) {
	for i := 0; i < maxBulletCount; i++ {
		if g.bullets[i].Visible {
			ebitenutil.DebugPrintAt(screen, "bullet is visible", 0, 60+i*15)
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Rotate(g.bullets[i].Direction)
			op.GeoM.Translate(g.bullets[i].X, g.bullets[i].Y)
			screen.DrawImage(bulletImg, op)
		}
	}
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

func checkEmptyBullet(bullets [maxBulletCount]*bullet.Bullet) int {
	for i := 0; i < maxBulletCount; i++ {
		if !bullets[i].Visible {
			return i
		}
	}
	return -1
}

func main() {
	ebiten.SetWindowSize(screenX, screenY)
	ebiten.SetWindowTitle("EbitenGame")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
