// debug.go - Debug overlay and performance monitoring
package rengine

import (
	"fmt"
	"runtime"
	"time"
)

type DebugInfo struct {
	Enabled       bool
	FPS           float64
	FrameTime     time.Duration
	RAMUsageMB    float64
	RAMAllocMB    float64
	NumGoroutines int
	PlayerPosX    float64
	PlayerPosY    float64
	PlayerDirX    float64
	PlayerDirY    float64
	
	lastFrameTime time.Time
	frameCount    int
	fpsUpdateTime time.Time
}

// Enables the debug overlay
func (e *Engine) EnableDebug() {
	if e.DebugInfo == nil {
		e.DebugInfo = &DebugInfo{
			Enabled:       true,
			lastFrameTime: time.Now(),
			fpsUpdateTime: time.Now(),
		}
	} else {
		e.DebugInfo.Enabled = true
	}
}

// Disables the debug overlay
func (e *Engine) DisableDebug() {
	if e.DebugInfo != nil {
		e.DebugInfo.Enabled = false
	}
}

// Toggles the debug overlay on/off
func (e *Engine) ToggleDebug() {
	if e.DebugInfo == nil {
		e.EnableDebug()
	} else {
		e.DebugInfo.Enabled = !e.DebugInfo.Enabled
	}
}

// Returns whether debug mode is enabled
func (e *Engine) IsDebugEnabled() bool {
	return e.DebugInfo != nil && e.DebugInfo.Enabled
}

// Updates all debug information
func (e *Engine) UpdateDebugInfo() {
	if e.DebugInfo == nil || !e.DebugInfo.Enabled {
		return
	}

	now := time.Now()
	
	e.DebugInfo.FrameTime = now.Sub(e.DebugInfo.lastFrameTime)
	e.DebugInfo.lastFrameTime = now
	
	if now.Sub(e.DebugInfo.fpsUpdateTime) >= time.Second {
		e.DebugInfo.FPS = float64(e.DebugInfo.frameCount) / now.Sub(e.DebugInfo.fpsUpdateTime).Seconds()
		e.DebugInfo.frameCount = 0
		e.DebugInfo.fpsUpdateTime = now
	}
	
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	e.DebugInfo.RAMUsageMB = float64(m.Alloc) / 1024 / 1024
	e.DebugInfo.RAMAllocMB = float64(m.TotalAlloc) / 1024 / 1024
	
	e.DebugInfo.NumGoroutines = runtime.NumGoroutine()
	
	if e.Player != nil {
		e.DebugInfo.PlayerPosX = e.Player.PosX
		e.DebugInfo.PlayerPosY = e.Player.PosY
		e.DebugInfo.PlayerDirX = e.Player.DirX
		e.DebugInfo.PlayerDirY = e.Player.DirY
	}
}

// Draws the debug information on screen
func (e *Engine) RenderDebugOverlay() {
	if e.DebugInfo == nil || !e.DebugInfo.Enabled {
		return
	}

	debugLines := []string{
		fmt.Sprintf("FPS: %.1f", e.DebugInfo.FPS),
		fmt.Sprintf("Frame Time: %.2fms", e.DebugInfo.FrameTime.Seconds()*1000),
		fmt.Sprintf("RAM Usage: %.2f MB", e.DebugInfo.RAMUsageMB),
		fmt.Sprintf("RAM Total Alloc: %.2f MB", e.DebugInfo.RAMAllocMB),
		fmt.Sprintf("Goroutines: %d", e.DebugInfo.NumGoroutines),
		"",
		fmt.Sprintf("Player Pos: (%.2f, %.2f)", e.DebugInfo.PlayerPosX, e.DebugInfo.PlayerPosY),
		fmt.Sprintf("Player Dir: (%.2f, %.2f)", e.DebugInfo.PlayerDirX, e.DebugInfo.PlayerDirY),
		"",
		fmt.Sprintf("Map Size: %dx%d", e.MapWidth, e.MapHeight),
		fmt.Sprintf("Screen: %dx%d", e.ScreenWidth, e.ScreenHeight),
		fmt.Sprintf("Textures Loaded: %d", len(e.Textures)),
	}

	e.drawDebugBackground(10, 10, 280, len(debugLines)*16+10)
	
	e.drawDebugText(debugLines, 15, 20)
}

// Draws a semi-transparent background for debug text
func (e *Engine) drawDebugBackground(x, y, width, height int) {
	for dy := 0; dy < height; dy++ {
		for dx := 0; dx < width; dx++ {
			px := x + dx
			py := y + dy
			if px >= 0 && px < e.ScreenWidth && py >= 0 && py < e.ScreenHeight {
				idx := (py*e.ScreenWidth + px) * 4
				e.FrameBuffer[idx] = byte(float64(e.FrameBuffer[idx]) * 0.3)
				e.FrameBuffer[idx+1] = byte(float64(e.FrameBuffer[idx+1]) * 0.3)
				e.FrameBuffer[idx+2] = byte(float64(e.FrameBuffer[idx+2]) * 0.3)
			}
		}
	}
}

// Draws text on the framebuffer (simple bitmap font)
func (e *Engine) drawDebugText(lines []string, startX, startY int) {
	for i, line := range lines {
		y := startY + i*16
		for j, char := range line {
			x := startX + j*8
			e.drawDebugChar(x, y, char)
		}
	}
}

// Draws a single character (very simple 8x16 representation)
func (e *Engine) drawDebugChar(x, y int, char rune) {
	for dy := 0; dy < 12; dy++ {
		for dx := 0; dx < 6; dx++ {
			if e.shouldDrawPixel(char, dx, dy) {
				px := x + dx
				py := y + dy
				if px >= 0 && px < e.ScreenWidth && py >= 0 && py < e.ScreenHeight {
					idx := (py*e.ScreenWidth + px) * 4
					e.FrameBuffer[idx] = 255
					e.FrameBuffer[idx+1] = 255
					e.FrameBuffer[idx+2] = 0
					e.FrameBuffer[idx+3] = 255
				}
			}
		}
	}
}

func (e *Engine) shouldDrawPixel(char rune, x, y int) bool {
	switch char {
	case ' ':
		return false
	case '.':
		return y == 10 && x >= 2 && x <= 3
	case ':':
		return (y == 3 && x >= 2 && x <= 3) || (y == 7 && x >= 2 && x <= 3)
	case '(':
		return x == 1 && y >= 2 && y <= 9
	case ')':
		return x == 4 && y >= 2 && y <= 9
	case ',':
		return (y == 10 && x >= 2 && x <= 3) || (y == 11 && x == 1)
	case '-':
		return y == 6 && x >= 1 && x <= 4
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return (x == 1 || x == 4) && y >= 2 && y <= 9 || 
		       (y == 2 || y == 9) && x >= 1 && x <= 4
	default:
		return (x == 1 || x == 4) && y >= 2 && y <= 9 || 
		       y == 2 && x >= 1 && x <= 4 ||
		       y == 6 && x >= 1 && x <= 4
	}
}

// Prints debug information to console
func (e *Engine) PrintDebugInfo() {
	if e.DebugInfo == nil || !e.DebugInfo.Enabled {
		return
	}

	fmt.Println("=== Debug Info ===")
	fmt.Printf("FPS: %.1f\n", e.DebugInfo.FPS)
	fmt.Printf("Frame Time: %.2fms\n", e.DebugInfo.FrameTime.Seconds()*1000)
	fmt.Printf("RAM Usage: %.2f MB\n", e.DebugInfo.RAMUsageMB)
	fmt.Printf("Goroutines: %d\n", e.DebugInfo.NumGoroutines)
	fmt.Printf("Player: (%.2f, %.2f) Dir: (%.2f, %.2f)\n", 
		e.DebugInfo.PlayerPosX, e.DebugInfo.PlayerPosY,
		e.DebugInfo.PlayerDirX, e.DebugInfo.PlayerDirY)
}

// Returns the current FPS
func (e *Engine) GetFPS() float64 {
	if e.DebugInfo == nil {
		return 0
	}
	return e.DebugInfo.FPS
}

// Returns current RAM usage in MB
func (e *Engine) GetRAMUsage() float64 {
	if e.DebugInfo == nil {
		return 0
	}
	return e.DebugInfo.RAMUsageMB
}
