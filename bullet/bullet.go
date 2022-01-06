package bullet

import (
  "math"
)

type Bullet struct {
	X         float64 // X posiiton of bullet
	Y         float64 // Y position of bullet
	Direction float64 // Direction of player
	r         float64 // Distance between center and lefttop of img
	thetaP    float64 // Angle of vector between center and lefttop of img
	speed     float64 // speed of bullet
	Visible   bool
}

func (b *Bullet) NewBullet(px, py, pd, bs float64) {
	b.X = px
	b.Y = py
	b.Direction = pd
	b.speed = bs
  b.Visible = true
}

func (b *Bullet) Move(){
  b.X = b.X + b.speed*math.Sin(b.Direction)
  b.Y = b.Y - b.speed*math.Cos(b.Direction)
}

