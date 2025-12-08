// input.go - Input handling
package rengine

import (
	"math"

	"github.com/go-gl/glfw/v3.3/glfw"
)

func (e *Engine) HandleInput() {
	if e.Player == nil {
		return
	}

	speedMultiplier := 1.0
	if e.IsSprinting() {
		speedMultiplier = e.SprintMultiplier
	}

	if e.KeyState[glfw.KeySpace] && e.HeightMap != nil {
		e.PlayerJump()
	}

	if e.KeyState[glfw.KeyW] {
		if e.HeightMap != nil {
			e.MovePlayerWithHeight(e.Player.DirX, e.Player.DirY, speedMultiplier)
		} else {
			e.movePlayer(e.Player.DirX, e.Player.DirY, speedMultiplier)
		}
	}

	if e.KeyState[glfw.KeyS] {
		if e.HeightMap != nil {
			e.MovePlayerWithHeight(-e.Player.DirX, -e.Player.DirY, speedMultiplier)
		} else {
			e.movePlayer(-e.Player.DirX, -e.Player.DirY, speedMultiplier)
		}
	}

	if e.KeyState[glfw.KeyD] {
		if e.HeightMap != nil {
			e.MovePlayerWithHeight(e.Player.DirY, -e.Player.DirX, speedMultiplier)
		} else {
			e.movePlayer(e.Player.DirY, -e.Player.DirX, speedMultiplier)
		}
	}

	if e.KeyState[glfw.KeyA] {
		if e.HeightMap != nil {
			e.MovePlayerWithHeight(-e.Player.DirY, e.Player.DirX, speedMultiplier)
		} else {
			e.movePlayer(-e.Player.DirY, e.Player.DirX, speedMultiplier)
		}
	}

	if e.KeyState[glfw.KeyQ] {
		e.rotatePlayer(e.Player.RotationSpeed * speedMultiplier)
	}

	if e.KeyState[glfw.KeyE] {
		e.rotatePlayer(-e.Player.RotationSpeed * speedMultiplier)
	}

	if e.HeightMap != nil {
		e.UpdatePlayerVerticalPosition()
	}
}

func (e *Engine) movePlayer(dirX, dirY, speedMultiplier float64) {
	newPosX := e.Player.PosX + dirX*e.Player.MoveSpeed*speedMultiplier
	newPosY := e.Player.PosY + dirY*e.Player.MoveSpeed*speedMultiplier

	if e.isWalkable(int(newPosX), int(e.Player.PosY)) {
		e.Player.PosX = newPosX
	}

	if e.isWalkable(int(e.Player.PosX), int(newPosY)) {
		e.Player.PosY = newPosY
	}
}

func (e *Engine) rotatePlayer(angle float64) {
	oldDirX := e.Player.DirX
	e.Player.DirX = e.Player.DirX*math.Cos(angle) - e.Player.DirY*math.Sin(angle)
	e.Player.DirY = oldDirX*math.Sin(angle) + e.Player.DirY*math.Cos(angle)

	oldPlaneX := e.Player.PlaneX
	e.Player.PlaneX = e.Player.PlaneX*math.Cos(angle) - e.Player.PlaneY*math.Sin(angle)
	e.Player.PlaneY = oldPlaneX*math.Sin(angle) + e.Player.PlaneY*math.Cos(angle)
}

func (e *Engine) isWalkable(x, y int) bool {
	if x < 0 || x >= e.MapWidth || y < 0 || y >= e.MapHeight {
		return false
	}
	return e.WorldMap[x][y] == 0
}

func (e *Engine) IsKeyPressed(key glfw.Key) bool {
	return e.KeyState[key]
}
