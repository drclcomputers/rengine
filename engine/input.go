// input.go - Input handling
package engine

import (
	"math"

	"github.com/go-gl/glfw/v3.3/glfw"
)

func (e *Engine) HandleInput() {
    if e.Player == nil {
        return
    }
    
    // Sprint
    if e.KeyState[glfw.KeyLeftShift] {
        e.Player.SprintMultiplier = 2.0
    } else {
        e.Player.SprintMultiplier = 1.0
    }
    
    // Movement
    if e.KeyState[glfw.KeyW] {
        e.Player.Move(e.Player.DirX, e.Player.DirY, 0.016, e.WorldMap, e.MapWidth, e.MapHeight)
    }
    if e.KeyState[glfw.KeyS] {
        e.Player.Move(-e.Player.DirX, -e.Player.DirY, 0.016, e.WorldMap, e.MapWidth, e.MapHeight)
    }
    if e.KeyState[glfw.KeyD] {
        e.Player.Move(e.Player.DirY, -e.Player.DirX, 0.016, e.WorldMap, e.MapWidth, e.MapHeight)
    }
    if e.KeyState[glfw.KeyA] {
        e.Player.Move(-e.Player.DirY, e.Player.DirX, 0.016, e.WorldMap, e.MapWidth, e.MapHeight)
    }
    
    // Rotation
    if e.KeyState[glfw.KeyQ] {
        e.Player.Rotate(e.Player.RotationSpeed)
    }
    if e.KeyState[glfw.KeyE] {
        e.Player.Rotate(-e.Player.RotationSpeed)
    }
    
    // Check pickups
    e.CheckItemPickups()
}

func (e *Engine) keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
    if action==glfw.Press{
        e.KeyState[key] = true
        
        if key==glfw.KeySpace { e.Player.Shoot() }
        if key==glfw.KeyR { e.Player.Reload() }
        if key==glfw.KeyF { e.HandleUseKey() }
        if key==glfw.KeyF3 { e.ToggleDebug() } 
        //if key==glfw.KeyM { e.ToggleMinimap() }
        if key==glfw.KeyN { e.MinimapCorner = (e.MinimapCorner + 1) % 4 }
        if key==glfw.Key1 || key==glfw.Key2 || key==glfw.Key3{ e.Player.SwitchWeapon(int(key - glfw.Key1)) }
        
	} 
	if action == glfw.Release {
        e.KeyState[key] = false
    }
    
    if key == glfw.KeyEscape && action == glfw.Press {
        e.Running = false
        w.SetShouldClose(true)
    }
}

func (e *Engine) movePlayer(dirX, dirY float64) {
	newPosX := e.Player.PosX + dirX*e.Player.MoveSpeed*e.Player.SprintMultiplier
	newPosY := e.Player.PosY + dirY*e.Player.MoveSpeed*e.Player.SprintMultiplier

	if IsWalkable(int(newPosX), int(e.Player.PosY), e.MapWidth, e.MapHeight, e.WorldMap) {
		e.Player.PosX = newPosX
	}

	if IsWalkable(int(e.Player.PosX), int(newPosY), e.MapWidth, e.MapHeight, e.WorldMap) {
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

func IsWalkable(x, y, MapWidth, MapHeight int, WorldMap [][]int) bool {
	if x < 0 || x >= MapWidth || y < 0 || y >= MapHeight {
		return false
	}
	return WorldMap[x][y] == 0
}

func (e *Engine) IsKeyPressed(key glfw.Key) bool {
	return e.KeyState[key]
}
