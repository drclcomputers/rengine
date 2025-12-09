// heightmap.go - Height map system for stairs, slopes, and vertical movement
package rengine

import "math"

type HeightMap struct {
	Data          [][]float64 
	MaxClimbAngle float64    
	StepHeight    float64   
	Gravity       float64
	JumpForce     float64 
}

func (e *Engine) SetupHeightMap() {
	e.HeightMap = &HeightMap{
		Data:          make([][]float64, e.MapWidth),
		MaxClimbAngle: 0.785, // 45 degrees
		StepHeight:    0.5,   
		Gravity:       0.01,
		JumpForce:     0.2,
	}

	for i := 0; i < e.MapWidth; i++ {
		e.HeightMap.Data[i] = make([]float64, e.MapHeight)
	}
}

func (e *Engine) SetCellHeight(x, y int, height float64) {
	if e.HeightMap == nil {
		e.SetupHeightMap()
	}
	if x >= 0 && x < e.MapWidth && y >= 0 && y < e.MapHeight {
		e.HeightMap.Data[x][y] = height
	}
}

func (e *Engine) GetGroundHeight(x, y float64) float64 {
	if e.HeightMap == nil {
		return 0.0
	}

	mapX := int(x)
	mapY := int(y)

	if mapX < 0 || mapX >= e.MapWidth || mapY < 0 || mapY >= e.MapHeight {
		return 0.0
	}

	fx := x - float64(mapX)
	fy := y - float64(mapY)

	h00 := e.HeightMap.Data[mapX][mapY]
	
	h10 := h00
	h01 := h00
	h11 := h00
	
	if mapX+1 < e.MapWidth {
		h10 = e.HeightMap.Data[mapX+1][mapY]
		if mapY+1 < e.MapHeight {
			h11 = e.HeightMap.Data[mapX+1][mapY+1]
		}
	}
	if mapY+1 < e.MapHeight {
		h01 = e.HeightMap.Data[mapX][mapY+1]
	}

	h0 := h00*(1-fx) + h10*fx
	h1 := h01*(1-fx) + h11*fx
	return h0*(1-fy) + h1*fy
}

func (e *Engine) SetPlayerHeight(height float64) {
	if e.Player != nil {
		e.Player.PosZ = height
	}
}

func (e *Engine) GetPlayerHeight() float64 {
	if e.Player == nil {
		return 0.0
	}
	return e.Player.PosZ
}

func (e *Engine) UpdatePlayerVerticalPosition() {
	if e.Player == nil || e.HeightMap == nil {
		return
	}

	groundHeight := e.GetGroundHeight(e.Player.PosX, e.Player.PosY)

	if e.Player.PosZ > groundHeight {
		e.Player.VerticalVelocity -= e.HeightMap.Gravity
		e.Player.PosZ += e.Player.VerticalVelocity

		if e.Player.PosZ <= groundHeight {
			e.Player.PosZ = groundHeight
			e.Player.VerticalVelocity = 0
			e.Player.IsGrounded = true
		}
	} else {
		e.Player.PosZ = groundHeight
		e.Player.IsGrounded = true
		e.Player.VerticalVelocity = 0
	}
}

func (e *Engine) PlayerJump() {
	if e.Player == nil || e.HeightMap == nil {
		return
	}

	if e.Player.IsGrounded {
		e.Player.VerticalVelocity = e.HeightMap.JumpForce
		e.Player.IsGrounded = false
	}
}

func (e *Engine) CanClimbSlope(currentHeight, targetHeight, distance float64) bool {
	if e.HeightMap == nil {
		return true
	}

	heightDiff := targetHeight - currentHeight
	angle := math.Atan2(heightDiff, distance)

	return math.Abs(angle) <= e.HeightMap.MaxClimbAngle || math.Abs(heightDiff) <= e.HeightMap.StepHeight
}

func (e *Engine) MovePlayerWithHeight(dirX, dirY, speedMultiplier float64) {
	if e.Player == nil {
		return
	}

	newPosX := e.Player.PosX + dirX*e.Player.MoveSpeed*speedMultiplier
	newPosY := e.Player.PosY + dirY*e.Player.MoveSpeed*speedMultiplier

	if !e.isWalkable(int(newPosX), int(e.Player.PosY)) ||
	   !e.isWalkable(int(e.Player.PosX), int(newPosY)) {
		return
	}

	if e.HeightMap != nil {
		currentHeight := e.GetGroundHeight(e.Player.PosX, e.Player.PosY)
		targetHeight := e.GetGroundHeight(newPosX, newPosY)
		distance := math.Sqrt(dirX*dirX + dirY*dirY) * e.Player.MoveSpeed * speedMultiplier

		if !e.CanClimbSlope(currentHeight, targetHeight, distance) {
			return 
		}
	}

	if e.isWalkable(int(newPosX), int(e.Player.PosY)) {
		e.Player.PosX = newPosX
	}

	if e.isWalkable(int(e.Player.PosX), int(newPosY)) {
		e.Player.PosY = newPosY
	}

	if e.HeightMap != nil {
		e.UpdatePlayerVerticalPosition()
	}
}

func (e *Engine) CreateStairs(x1, y1, x2, y2 int, startHeight, endHeight float64) {
	if e.HeightMap == nil {
		e.SetupHeightMap()
	}

	dx := x2 - x1
	dy := y2 - y1
	steps := max(abs(dx), abs(dy))

	if steps == 0 {
		return
	}

	heightStep := (endHeight - startHeight) / float64(steps)

	for i := 0; i <= steps; i++ {
		t := float64(i) / float64(steps)
		x := x1 + int(float64(dx)*t)
		y := y1 + int(float64(dy)*t)
		height := startHeight + heightStep*float64(i)
		
		e.SetCellHeight(x, y, height)
	}
}

func (e *Engine) CreateRamp(x1, y1, x2, y2 int, startHeight, endHeight float64) {
	if e.HeightMap == nil {
		e.SetupHeightMap()
	}

	for x := min(x1, x2); x <= max(x1, x2); x++ {
		for y := min(y1, y2); y <= max(y1, y2); y++ {
			tx := float64(x-x1) / float64(x2-x1)
			ty := float64(y-y1) / float64(y2-y1)
			t := (tx + ty) / 2.0
			
			if t < 0 {
				t = 0
			}
			if t > 1 {
				t = 1
			}
			
			height := startHeight + (endHeight-startHeight)*t
			e.SetCellHeight(x, y, height)
		}
	}
}

func abs(x int) int {
	if x < 0 { return -x }
	return x
}
