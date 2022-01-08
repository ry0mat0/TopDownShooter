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
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"TopDownShooter/bullet"
	// "TopDownShooter/controll"
	"TopDownShooter/plane"
)

const (
	pi       = math.Pi
	screenX  = 800
	screenY  = 600
	fontSize = 10

	// game modes
	modeTitle    = 0
	modeGame     = 1
	modeGameover = 2

	//player
	speed       = 2.5
	accel       = 2.0
	brake       = 0.5
	rotSpeed    = 0.03
	playerSizeX = 50
	playerSizeY = 50

	//bullets
	maxBulletCount = 20
	bulletSizeX    = 2
	bulletSizeY    = 10
	bulletSpeed    = 8.0
	interval       = 10
)

var (
	playerImg  *ebiten.Image
	enemyImg   *ebiten.Image
	bulletImg  *ebiten.Image
	arcadeFont font.Face
)

func init() {
	var err error
	playerImg, _, err = ebitenutil.NewImageFromFile("image/plane.png")
	if err != nil {
		log.Fatal(err)
	}
	bulletImg, _, err = ebitenutil.NewImageFromFile("image/bullet.png")
	if err != nil {
		log.Fatal(err)
	}
	enemyImg, _, err = ebitenutil.NewImageFromFile("image/enemy.png")
	if err != nil {
		log.Fatal(err)
	}
	tt, err := opentype.Parse(fonts.PressStart2P_ttf)
	if err != nil {
		log.Fatal(err)
	}
	const dpi = 72
	arcadeFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
}

type Game struct {
	mode    int
	plane  *plane.Plane
	bullets [maxBulletCount]*bullet.Bullet
	enemy   *plane.Plane
}

// NewGame method
func NewGame() *Game {
	g := &Game{}
	g.init()
	return g
}

func (g *Game) init() {
	g.plane = &plane.Plane{}
	g.plane.NewPlayer(playerSizeX, playerSizeY, screenX, screenY,
		speed, accel, brake, rotSpeed, interval)
	g.enemy = &plane.Plane{}
	g.enemy.NewPlayer(playerSizeX, playerSizeY, screenX, screenY,
		speed, accel, brake, rotSpeed, interval)
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
    g.plane.MovePlayer(g.moveKey(), screenX, screenY)
    g.enemy.MoveEnemy(screenX, screenY)
		for i := 0; i < maxBulletCount; i++ {
			g.bullets[i].Move(screenX, screenY)
		}
		if g.isKeyPressed() {
			if g.plane.Gun_interval == 0 {
				i := checkEmptyBullet(g.bullets)
				if i >= 0 {
					g.bullets[i].NewBullet(g.plane.X, g.plane.Y, g.plane.Direction, bulletSpeed)
					g.plane.Gun_interval = interval
				}
			} else {
				g.plane.CountdownInterval()
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
		text.Draw(screen, "PRESS SPACE KEY", arcadeFont, 245, 240, color.Black)
	case modeGame:
		g.DrawUI(screen)
		g.DrawPlayer(screen)
    g.DrawBullets(screen)
    g.DrawEnemy(screen)
	case modeGameover:
	}
}

func (g *Game) DrawPlayer(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Rotate(g.plane.Direction)
	op.GeoM.Translate(g.plane.X, g.plane.Y)
	s := strconv.FormatFloat(g.plane.X, 'f', 1, 64)
	t := strconv.FormatFloat(g.plane.Y, 'f', 1, 64)
	u := strconv.FormatFloat(g.plane.Direction*180.0/pi, 'f', 1, 64)
	ebitenutil.DebugPrintAt(screen, strings.Join([]string{"X is ", s}, ""), 0, 15)
	ebitenutil.DebugPrintAt(screen, strings.Join([]string{"Y is ", t}, ""), 0, 30)
	ebitenutil.DebugPrintAt(screen, strings.Join([]string{"D is ", u}, ""), 0, 45)
	screen.DrawImage(playerImg, op)
}

func (g *Game) DrawBullets(screen *ebiten.Image) {
	for i := 0; i < maxBulletCount; i++ {
		if g.bullets[i].Visible {
			str := strings.Join([]string{"bullet", strconv.FormatInt(int64(i), 10), "is visible"}, "")
			ebitenutil.DebugPrintAt(screen, str, 0, 60+i*15)
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Rotate(g.bullets[i].Direction)
			op.GeoM.Translate(g.bullets[i].X, g.bullets[i].Y)
			screen.DrawImage(bulletImg, op)
		}
	}
}

func (g *Game) DrawEnemy(screen *ebiten.Image) {
  op :=&ebiten.DrawImageOptions{}
  op.GeoM.Translate(g.enemy.X, g.enemy.Y)
  screen.DrawImage(enemyImg, op)
}
func (g *Game) DrawUI(screen *ebiten.Image) {
	text.Draw(screen, "LIFE", arcadeFont, 10, screenY-10, color.White)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return screenX, screenY
}

func (g *Game) isKeyPressed() bool {
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		return true
	}
	return false
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
