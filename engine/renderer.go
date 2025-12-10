// renderer.go - Raycasting renderer with proper floor/ceiling
package engine

import (
	"fmt"
	"math"
	"time"

	"github.com/go-gl/gl/v2.1/gl"
)

func (e *Engine) Render() {
	if e.FloorTexture == nil || e.CeilingTexture == nil {
		return
	}

	e.renderFloorCeiling()

	for i := range e.ZBuffer {
		e.ZBuffer[i] = 1e30
	}

	for x := 0; x < e.ScreenWidth; x++ {
		e.castRay(x)
	}

	e.RenderSprites() 

	e.RenderHUD()

	e.RenderShootingEffects()

	e.RenderDebugOverlay()

	time.Sleep(time.Duration(1000/e.FPS))
}

func (e *Engine) RenderHUD() {
	if e.Combat == nil {
		return
	}

	healthPercent := e.Combat.PlayerHealth / e.Combat.PlayerMaxHealth
	healthBarWidth := 200
	healthBarHeight := 20
	healthX := 20
	healthY := e.ScreenHeight - 40

	for y := range healthBarHeight {
		for x := range healthBarWidth {
			px := healthX + x
			py := healthY + y
			if px >= 0 && px < e.ScreenWidth && py >= 0 && py < e.ScreenHeight {
				idx := (py*e.ScreenWidth + px) * 4
				e.FrameBuffer[idx] = 100
				e.FrameBuffer[idx+1] = 0
				e.FrameBuffer[idx+2] = 0
				e.FrameBuffer[idx+3] = 180
			}
		}
	}

	fillWidth := int(float64(healthBarWidth) * healthPercent)
	for y := range healthBarHeight {
		for x := range fillWidth {
			px := healthX + x
			py := healthY + y
			if px >= 0 && px < e.ScreenWidth && py >= 0 && py < e.ScreenHeight {
				idx := (py*e.ScreenWidth + px) * 4
				if healthPercent > 0.5 {
					e.FrameBuffer[idx] = 0
					e.FrameBuffer[idx+1] = 255
					e.FrameBuffer[idx+2] = 0
				} else if healthPercent > 0.25 {
					e.FrameBuffer[idx] = 255
					e.FrameBuffer[idx+1] = 255
					e.FrameBuffer[idx+2] = 0
				} else {
					e.FrameBuffer[idx] = 255
					e.FrameBuffer[idx+1] = 0
					e.FrameBuffer[idx+2] = 0
				}
				e.FrameBuffer[idx+3] = 255
			}
		}
	}

	ammoText := fmt.Sprintf("Ammo: %d", e.Combat.Ammo)
	e.drawDebugText([]string{ammoText}, healthX, healthY - 25)
}

func (e *Engine) RenderShootingEffects() {
	// Muzzle flash effect
	if e.IsMuzzleFlashActive() {
		e.renderMuzzleFlash()
	}

	// Hit marker (crosshair confirmation)
	if e.IsHitMarkerActive() {
		e.renderHitMarker()
	}

	// Always render crosshair
	e.renderCrosshair()
}

// Render muzzle flash as screen edge flash
func (e *Engine) renderMuzzleFlash() {
	flashIntensity := 100
	borderSize := 5

	// Top border
	for y := 0; y < borderSize; y++ {
		for x := 0; x < e.ScreenWidth; x++ {
			idx := (y*e.ScreenWidth + x) * 4
			e.FrameBuffer[idx] = byte(min(255, int(e.FrameBuffer[idx])+flashIntensity))
			e.FrameBuffer[idx+1] = byte(min(255, int(e.FrameBuffer[idx+1])+flashIntensity))
			e.FrameBuffer[idx+2] = 0
			e.FrameBuffer[idx+3] = 255
		}
	}

	// Bottom border
	for y := e.ScreenHeight - borderSize; y < e.ScreenHeight; y++ {
		for x := 0; x < e.ScreenWidth; x++ {
			idx := (y*e.ScreenWidth + x) * 4
			e.FrameBuffer[idx] = byte(min(255, int(e.FrameBuffer[idx])+flashIntensity))
			e.FrameBuffer[idx+1] = byte(min(255, int(e.FrameBuffer[idx+1])+flashIntensity))
			e.FrameBuffer[idx+2] = 0
			e.FrameBuffer[idx+3] = 255
		}
	}

	// Left border
	for y := 0; y < e.ScreenHeight; y++ {
		for x := 0; x < borderSize; x++ {
			idx := (y*e.ScreenWidth + x) * 4
			e.FrameBuffer[idx] = byte(min(255, int(e.FrameBuffer[idx])+flashIntensity))
			e.FrameBuffer[idx+1] = byte(min(255, int(e.FrameBuffer[idx+1])+flashIntensity))
			e.FrameBuffer[idx+2] = 0
			e.FrameBuffer[idx+3] = 255
		}
	}

	// Right border
	for y := 0; y < e.ScreenHeight; y++ {
		for x := e.ScreenWidth - borderSize; x < e.ScreenWidth; x++ {
			idx := (y*e.ScreenWidth + x) * 4
			e.FrameBuffer[idx] = byte(min(255, int(e.FrameBuffer[idx])+flashIntensity))
			e.FrameBuffer[idx+1] = byte(min(255, int(e.FrameBuffer[idx+1])+flashIntensity))
			e.FrameBuffer[idx+2] = 0
			e.FrameBuffer[idx+3] = 255
		}
	}
}

func (e *Engine) renderCrosshair() {
	centerX := e.ScreenWidth / 2
	centerY := e.ScreenHeight / 2
	size := 10
	thickness := 2
	gap := 3

	// Horizontal line (left)
	for x := centerX - size - gap; x < centerX - gap; x++ {
		for y := centerY - thickness/2; y < centerY + thickness/2; y++ {
			if x >= 0 && x < e.ScreenWidth && y >= 0 && y < e.ScreenHeight {
				idx := (y*e.ScreenWidth + x) * 4
				e.FrameBuffer[idx] = 255
				e.FrameBuffer[idx+1] = 255
				e.FrameBuffer[idx+2] = 255
				e.FrameBuffer[idx+3] = 255
			}
		}
	}

	// Horizontal line (right)
	for x := centerX + gap; x < centerX + size + gap; x++ {
		for y := centerY - thickness/2; y < centerY + thickness/2; y++ {
			if x >= 0 && x < e.ScreenWidth && y >= 0 && y < e.ScreenHeight {
				idx := (y*e.ScreenWidth + x) * 4
				e.FrameBuffer[idx] = 255
				e.FrameBuffer[idx+1] = 255
				e.FrameBuffer[idx+2] = 255
				e.FrameBuffer[idx+3] = 255
			}
		}
	}

	// Vertical line (top)
	for y := centerY - size - gap; y < centerY - gap; y++ {
		for x := centerX - thickness/2; x < centerX + thickness/2; x++ {
			if x >= 0 && x < e.ScreenWidth && y >= 0 && y < e.ScreenHeight {
				idx := (y*e.ScreenWidth + x) * 4
				e.FrameBuffer[idx] = 255
				e.FrameBuffer[idx+1] = 255
				e.FrameBuffer[idx+2] = 255
				e.FrameBuffer[idx+3] = 255
			}
		}
	}

	// Vertical line (bottom)
	for y := centerY + gap; y < centerY + size + gap; y++ {
		for x := centerX - thickness/2; x < centerX + thickness/2; x++ {
			if x >= 0 && x < e.ScreenWidth && y >= 0 && y < e.ScreenHeight {
				idx := (y*e.ScreenWidth + x) * 4
				e.FrameBuffer[idx] = 255
				e.FrameBuffer[idx+1] = 255
				e.FrameBuffer[idx+2] = 255
				e.FrameBuffer[idx+3] = 255
			}
		}
	}
}

// Render hit marker (X shape when you hit an enemy)
func (e *Engine) renderHitMarker() {
	centerX := e.ScreenWidth / 2
	centerY := e.ScreenHeight / 2
	size := 8
	thickness := 2

	// Draw X shape
	for i := -size; i <= size; i++ {
		// Diagonal line 1 (\)
		for t := -thickness/2; t <= thickness/2; t++ {
			x := centerX + i + t
			y := centerY + i + t
			if x >= 0 && x < e.ScreenWidth && y >= 0 && y < e.ScreenHeight {
				idx := (y*e.ScreenWidth + x) * 4
				e.FrameBuffer[idx] = 255
				e.FrameBuffer[idx+1] = 0
				e.FrameBuffer[idx+2] = 0
				e.FrameBuffer[idx+3] = 255
			}
		}

		// Diagonal line 2 (/)
		for t := -thickness/2; t <= thickness/2; t++ {
			x := centerX + i + t
			y := centerY - i + t
			if x >= 0 && x < e.ScreenWidth && y >= 0 && y < e.ScreenHeight {
				idx := (y*e.ScreenWidth + x) * 4
				e.FrameBuffer[idx] = 255
				e.FrameBuffer[idx+1] = 0
				e.FrameBuffer[idx+2] = 0
				e.FrameBuffer[idx+3] = 255
			}
		}
	}
}

func (e *Engine) renderFloorCeiling() {
	for y := e.ScreenHeight/2 + 1; y < e.ScreenHeight; y++ {
		rayDirX0 := e.Player.DirX - e.Player.PlaneX
		rayDirY0 := e.Player.DirY - e.Player.PlaneY
		rayDirX1 := e.Player.DirX + e.Player.PlaneX
		rayDirY1 := e.Player.DirY + e.Player.PlaneY

		p := y - e.ScreenHeight/2

		posZ := 0.5 * float64(e.ScreenHeight)

		rowDistance := posZ / float64(p)

		floorStepX := rowDistance * (rayDirX1 - rayDirX0) / float64(e.ScreenWidth)
		floorStepY := rowDistance * (rayDirY1 - rayDirY0) / float64(e.ScreenWidth)

		floorX := e.Player.PosX + rowDistance*rayDirX0
		floorY := e.Player.PosY + rowDistance*rayDirY0

		for x := 0; x < e.ScreenWidth; x++ {
			tx := int(float64(e.FloorTexture.Width) * (floorX - math.Floor(floorX))) & (e.FloorTexture.Width - 1)
			ty := int(float64(e.FloorTexture.Height) * (floorY - math.Floor(floorY))) & (e.FloorTexture.Height - 1)

			floorX += floorStepX
			floorY += floorStepY

			r, g, b := e.FloorTexture.GetPixel(tx, ty)
			factor := 1.0 / (1.0 + rowDistance*0.1)
			r = byte(float64(r) * factor)
			g = byte(float64(g) * factor)
			b = byte(float64(b) * factor)

			idx := (y*e.ScreenWidth + x) * 4
			e.FrameBuffer[idx] = r
			e.FrameBuffer[idx+1] = g
			e.FrameBuffer[idx+2] = b
			e.FrameBuffer[idx+3] = 255

			r2, g2, b2 := e.CeilingTexture.GetPixel(tx, ty)
			r2 = byte(float64(r2) * factor)
			g2 = byte(float64(g2) * factor)
			b2 = byte(float64(b2) * factor)

			ceilingY := e.ScreenHeight - y - 1
			idx2 := (ceilingY*e.ScreenWidth + x) * 4
			e.FrameBuffer[idx2] = r2
			e.FrameBuffer[idx2+1] = g2
			e.FrameBuffer[idx2+2] = b2
			e.FrameBuffer[idx2+3] = 255
		}
	}
}

func (e *Engine) castRay(x int) {
	cameraX := 2*float64(x)/float64(e.ScreenWidth) - 1
	rayDirX := e.Player.DirX + e.Player.PlaneX*cameraX
	rayDirY := e.Player.DirY + e.Player.PlaneY*cameraX

	mapX := int(e.Player.PosX)
	mapY := int(e.Player.PosY)

	deltaDistX := math.Abs(1 / rayDirX)
	deltaDistY := math.Abs(1 / rayDirY)

	var sideDistX, sideDistY float64
	var stepX, stepY int
	var side int

	if rayDirX < 0 {
		stepX = -1
		sideDistX = (e.Player.PosX - float64(mapX)) * deltaDistX
	} else {
		stepX = 1
		sideDistX = (float64(mapX) + 1.0 - e.Player.PosX) * deltaDistX
	}

	if rayDirY < 0 {
		stepY = -1
		sideDistY = (e.Player.PosY - float64(mapY)) * deltaDistY
	} else {
		stepY = 1
		sideDistY = (float64(mapY) + 1.0 - e.Player.PosY) * deltaDistY
	}

	hit := false
	for !hit {
		if sideDistX < sideDistY {
			sideDistX += deltaDistX
			mapX += stepX
			side = 0
		} else {
			sideDistY += deltaDistY
			mapY += stepY
			side = 1
		}

		if mapX < 0 || mapX >= e.MapWidth || mapY < 0 || mapY >= e.MapHeight {
			return
		}
		if e.WorldMap[mapX][mapY] > 0 {
			hit = true
		}
	}

	var perpWallDist float64
	if side == 0 {
		perpWallDist = (float64(mapX) - e.Player.PosX + (1-float64(stepX))/2) / rayDirX
	} else {
		perpWallDist = (float64(mapY) - e.Player.PosY + (1-float64(stepY))/2) / rayDirY
	}

	e.ZBuffer[x] = perpWallDist

	lineHeight := int(float64(e.ScreenHeight) / perpWallDist)
	drawStart := -lineHeight/2 + e.ScreenHeight/2
	drawStart = max(0, drawStart)
	drawEnd := lineHeight/2 + e.ScreenHeight/2
	if drawEnd >= e.ScreenHeight {
		drawEnd = e.ScreenHeight - 1
	}

	wallType := e.WorldMap[mapX][mapY] - 1
	if wallType < 0 || wallType >= len(e.Textures) {
		return
	}
	texture := e.Textures[wallType]

	var wallX float64
	if side == 0 {
		wallX = e.Player.PosY + perpWallDist*rayDirY
	} else {
		wallX = e.Player.PosX + perpWallDist*rayDirX
	}
	wallX -= math.Floor(wallX)

	texX := int(wallX * float64(texture.Width))
	if side == 0 && rayDirX > 0 {
		texX = texture.Width - texX - 1
	}
	if side == 1 && rayDirY < 0 {
		texX = texture.Width - texX - 1
	}

	e.drawTexturedStripe(x, drawStart, drawEnd, lineHeight, texX, texture, perpWallDist)
}

func (e *Engine) drawTexturedStripe(x, drawStart, drawEnd, lineHeight, texX int, texture *Texture, perpWallDist float64) {
	for y := drawStart; y < drawEnd; y++ {
		d := y*256 - e.ScreenHeight*128 + lineHeight*128
		texY := ((d * texture.Height) / lineHeight) / 256
		texY = max(0, texY)
		if texY >= texture.Height {
			texY = texture.Height - 1
		}

		r, g, b := texture.GetPixel(texX, texY)

		factor := 1.0 / (1.0 + perpWallDist*0.1)
		r = byte(float64(r) * factor)
		g = byte(float64(g) * factor)
		b = byte(float64(b) * factor)

		fbIdx := (y*e.ScreenWidth + x) * 4
		e.FrameBuffer[fbIdx] = r
		e.FrameBuffer[fbIdx+1] = g
		e.FrameBuffer[fbIdx+2] = b
		e.FrameBuffer[fbIdx+3] = 255
	}
}

func (e *Engine) DrawFrameBuffer() {
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.Enable(gl.TEXTURE_2D)
	gl.BindTexture(gl.TEXTURE_2D, e.ScreenTexture)
	gl.TexSubImage2D(
		gl.TEXTURE_2D,
		0,
		0, 0,
		int32(e.ScreenWidth),
		int32(e.ScreenHeight),
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(e.FrameBuffer),
	)

	gl.Begin(gl.QUADS)
	gl.TexCoord2f(0, 0)
	gl.Vertex2i(0, 0)
	gl.TexCoord2f(1, 0)
	gl.Vertex2i(int32(e.ScreenWidth), 0)
	gl.TexCoord2f(1, 1)
	gl.Vertex2i(int32(e.ScreenWidth), int32(e.ScreenHeight))
	gl.TexCoord2f(0, 1)
	gl.Vertex2i(0, int32(e.ScreenHeight))
	gl.End()
	gl.Disable(gl.TEXTURE_2D)
}
