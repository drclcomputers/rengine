// input.go - Input handling
package engine

import (
	"math"

	"github.com/go-gl/glfw/v3.3/glfw"
)

// Player movement and rotation
func (e *Engine) HandleInput() {
	if e.Player == nil {
		return
	}

	if e.KeyState[glfw.KeyW] {
		e.movePlayer(e.Player.DirX, e.Player.DirY)
	}

	if e.KeyState[glfw.KeyS] {
		e.movePlayer(-e.Player.DirX, -e.Player.DirY)
	}

	if e.KeyState[glfw.KeyD] {
		e.movePlayer(e.Player.DirY, -e.Player.DirX)
	}

	if e.KeyState[glfw.KeyA] {
		e.movePlayer(-e.Player.DirY, e.Player.DirX)
	}

	if e.KeyState[glfw.KeyQ] {
		e.rotatePlayer(e.Player.RotationSpeed)
	}

	if e.KeyState[glfw.KeyE] {
		e.rotatePlayer(-e.Player.RotationSpeed)
	}
}

// Handles collision detection and movement
func (e *Engine) movePlayer(dirX, dirY float64) {
	newPosX := e.Player.PosX + dirX*e.Player.MoveSpeed
	newPosY := e.Player.PosY + dirY*e.Player.MoveSpeed

	if e.isWalkable(int(newPosX), int(e.Player.PosY)) {
		e.Player.PosX = newPosX
	}

	if e.isWalkable(int(e.Player.PosX), int(newPosY)) {
		e.Player.PosY = newPosY
	}
}

// Rotates the player's direction and camera plane
func (e *Engine) rotatePlayer(angle float64) {
	oldDirX := e.Player.DirX
	e.Player.DirX = e.Player.DirX*math.Cos(angle) - e.Player.DirY*math.Sin(angle)
	e.Player.DirY = oldDirX*math.Sin(angle) + e.Player.DirY*math.Cos(angle)

	oldPlaneX := e.Player.PlaneX
	e.Player.PlaneX = e.Player.PlaneX*math.Cos(angle) - e.Player.PlaneY*math.Sin(angle)
	e.Player.PlaneY = oldPlaneX*math.Sin(angle) + e.Player.PlaneY*math.Cos(angle)
}

// Checks if a position is walkable
func (e *Engine) isWalkable(x, y int) bool {
	if x < 0 || x >= e.MapWidth || y < 0 || y >= e.MapHeight {
		return false
	}
	return e.WorldMap[x][y] == 0
}

// Checks if a key is pressed
func (e *Engine) IsKeyPressed(key glfw.Key) bool {
	return e.KeyState[key]
}
