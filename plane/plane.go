package plane

import (
	// "fmt"
	"math"
)

type Plane struct {
	// Size
	X         float64 // X posiiton of player
	Y         float64 // Y position of player
	Direction float64 // Direction of player
	r         float64 // Distance between center and lefttop of img
	thetaP    float64 // Angle of vector between center and lefttop of img
	Life      int
	// Mobility
	speed     float64 // speed of player
	accel     float64 // accel magnitude of player
	brake     float64 // brake magnitude of player
	rot_speed float64 // rotation speed of player
	// Weapon
	Gun_interval int
}

//Initialise player parameters
func (p *Plane) NewPlayer(px, py, sx, sy int,
	speed, accel, brake, rot_speed float64, interval int) {
	p.X = float64(sx)/2 - float64(px)/2
	p.Y = float64(sy)/2 - float64(py)/2
	p.r = math.Sqrt(float64(px)*float64(px)+float64(py)*float64(py)) / 2.0
	p.thetaP = math.Atan(float64(py) / float64(px))
	p.speed = speed
	p.accel = accel
	p.brake = brake
	p.rot_speed = rot_speed
	p.Gun_interval = interval
}

func (p *Plane) MovePlayer(keys [4]bool, sx, sy int) {
	p.rotateCenter(keys)
	p.translate(keys, sx, sy)
}

func (p *Plane) MoveEnemy(sx, sy int) {
  p.Y = p.Y - p.speed*p.brake
  p.checkScreenEdge(sx, sy)
}

func (p *Plane) translate(keys [4]bool, sx, sy int) {
	var m float64
	if keys[0] == true {
		m = p.accel
	} else if keys[2] == true {
		m = p.brake
	} else {
		m = 1.0
	}
	// switch dir {
	// case 1: //KeyW
	// 	m = p.accel
	// case 3: //KeyS
	// 	m = p.brake
	// default:
	// 	m = 1.0
	// }
	p.Y = p.Y - m*p.speed*math.Cos(p.Direction)
	p.X = p.X + m*p.speed*math.Sin(p.Direction)
	p.checkScreenEdge(sx, sy)
}

// Update angle for rotation
func (p *Plane) rotate(keys [4]bool) {
	if keys[1] { //KeyA
		p.Direction = p.Direction - p.rot_speed
	}
	if keys[3] { //KeyD
		p.Direction = p.Direction + p.rot_speed
	}
}

// Rotate plater around the center of img
func (p *Plane) rotateCenter(keys [4]bool) {
	xc1, yc1 := p.getCenter()
	p.rotate(keys)
	xc2, yc2 := p.getCenter()
	p.X = p.X + (xc1 - xc2)
	p.Y = p.Y + (yc1 - yc2)
}

// Detece edge of screen and move to the opposite side
func (p *Plane) checkScreenEdge(sx, sy int) {
	if p.X < 0 {
		p.X = p.X + float64(sx)
	}
	if p.X > float64(sx) {
		p.X = p.X - float64(sx)
	}
	if p.Y < 0 {
		p.Y = p.Y + float64(sy)
	}
	if p.Y > float64(sy) {
		p.Y = p.Y - float64(sy)
	}
}

//Get center position of player img
func (p *Plane) getCenter() (xc, yc float64) {
	xc = p.X + p.r*math.Cos(p.Direction+p.thetaP)
	yc = p.X + p.r*math.Sin(p.Direction+p.thetaP)
	return
}

func (p *Plane) CountdownInterval() {
	p.Gun_interval = p.Gun_interval - 1
}
