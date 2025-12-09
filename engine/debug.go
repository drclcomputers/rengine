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

func (e *Engine) DisableDebug() {
	if e.DebugInfo != nil {
		e.DebugInfo.Enabled = false
	}
}

func (e *Engine) ToggleDebug() {
	if e.DebugInfo == nil {
		e.EnableDebug()
	} else {
		e.DebugInfo.Enabled = !e.DebugInfo.Enabled
	}
}

func (e *Engine) IsDebugEnabled() bool {
	return e.DebugInfo != nil && e.DebugInfo.Enabled
}

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

func (e *Engine) drawDebugBackground(x, y, width, height int) {
	for dy := 0; dy < height; dy++ {
		for dx := 0; dx < width; dx++ {
			px := x + dx
			py := y + dy
			if px >= 0 && px < e.ScreenWidth && py >= 0 && py < e.ScreenHeight {
				idx := (py*e.ScreenWidth + px) * 4
				e.FrameBuffer[idx] = 0
				e.FrameBuffer[idx+1] = 0
				e.FrameBuffer[idx+2] = 0
				e.FrameBuffer[idx+3] = 180
			}
		}
	}
}

func (e *Engine) drawDebugText(lines []string, startX, startY int) {
	for i, line := range lines {
		y := startY + i*14
		x := startX
		for _, char := range line {
			e.drawDebugChar(x, y, char)
			x += 8 
		}
	}
}

func (e *Engine) drawDebugChar(x, y int, char rune) {
	pattern := e.getCharPattern(char)
	
	for row := 0; row < 7; row++ {
		for col := 0; col < 5; col++ {
			if pattern[row]&(1<<(4-col)) != 0 {
				px := x + col
				py := y + row
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

func (e *Engine) GetFPS() float64 {
	if e.DebugInfo == nil {
		return 0
	}
	return e.DebugInfo.FPS
}

func (e *Engine) GetRAMUsage() float64 {
	if e.DebugInfo == nil {
		return 0
	}
	return e.DebugInfo.RAMUsageMB
}

func (e *Engine) getCharPattern(char rune) [7]byte {
	switch char {
	case '0':
		return [7]byte{0x0E, 0x11, 0x13, 0x15, 0x19, 0x11, 0x0E}
	case '1':
		return [7]byte{0x04, 0x0C, 0x04, 0x04, 0x04, 0x04, 0x0E}
	case '2':
		return [7]byte{0x0E, 0x11, 0x01, 0x02, 0x04, 0x08, 0x1F}
	case '3':
		return [7]byte{0x1F, 0x02, 0x04, 0x02, 0x01, 0x11, 0x0E}
	case '4':
		return [7]byte{0x02, 0x06, 0x0A, 0x12, 0x1F, 0x02, 0x02}
	case '5':
		return [7]byte{0x1F, 0x10, 0x1E, 0x01, 0x01, 0x11, 0x0E}
	case '6':
		return [7]byte{0x06, 0x08, 0x10, 0x1E, 0x11, 0x11, 0x0E}
	case '7':
		return [7]byte{0x1F, 0x01, 0x02, 0x04, 0x08, 0x08, 0x08}
	case '8':
		return [7]byte{0x0E, 0x11, 0x11, 0x0E, 0x11, 0x11, 0x0E}
	case '9':
		return [7]byte{0x0E, 0x11, 0x11, 0x0F, 0x01, 0x02, 0x0C}
	case 'A':
		return [7]byte{0x0E, 0x11, 0x11, 0x1F, 0x11, 0x11, 0x11}
	case 'B':
		return [7]byte{0x1E, 0x11, 0x11, 0x1E, 0x11, 0x11, 0x1E}
	case 'C':
		return [7]byte{0x0E, 0x11, 0x10, 0x10, 0x10, 0x11, 0x0E}
	case 'D':
		return [7]byte{0x1C, 0x12, 0x11, 0x11, 0x11, 0x12, 0x1C}
	case 'E':
		return [7]byte{0x1F, 0x10, 0x10, 0x1E, 0x10, 0x10, 0x1F}
	case 'F':
		return [7]byte{0x1F, 0x10, 0x10, 0x1E, 0x10, 0x10, 0x10}
	case 'G':
		return [7]byte{0x0E, 0x11, 0x10, 0x17, 0x11, 0x11, 0x0F}
	case 'H':
		return [7]byte{0x11, 0x11, 0x11, 0x1F, 0x11, 0x11, 0x11}
	case 'I':
		return [7]byte{0x0E, 0x04, 0x04, 0x04, 0x04, 0x04, 0x0E}
	case 'L':
		return [7]byte{0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x1F}
	case 'M':
		return [7]byte{0x11, 0x1B, 0x15, 0x15, 0x11, 0x11, 0x11}
	case 'N':
		return [7]byte{0x11, 0x11, 0x19, 0x15, 0x13, 0x11, 0x11}
	case 'O':
		return [7]byte{0x0E, 0x11, 0x11, 0x11, 0x11, 0x11, 0x0E}
	case 'P':
		return [7]byte{0x1E, 0x11, 0x11, 0x1E, 0x10, 0x10, 0x10}
	case 'R':
		return [7]byte{0x1E, 0x11, 0x11, 0x1E, 0x14, 0x12, 0x11}
	case 'S':
		return [7]byte{0x0E, 0x11, 0x10, 0x0E, 0x01, 0x11, 0x0E}
	case 'T':
		return [7]byte{0x1F, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04}
	case 'U':
		return [7]byte{0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x0E}
	case 'W':
		return [7]byte{0x11, 0x11, 0x11, 0x15, 0x15, 0x1B, 0x11}
	case 'a':
		return [7]byte{0x00, 0x00, 0x0E, 0x01, 0x0F, 0x11, 0x0F}
	case 'b':
		return [7]byte{0x10, 0x10, 0x16, 0x19, 0x11, 0x11, 0x1E}
	case 'c':
		return [7]byte{0x00, 0x00, 0x0E, 0x10, 0x10, 0x11, 0x0E}
	case 'd':
		return [7]byte{0x01, 0x01, 0x0D, 0x13, 0x11, 0x11, 0x0F}
	case 'e':
		return [7]byte{0x00, 0x00, 0x0E, 0x11, 0x1F, 0x10, 0x0E}
	case 'g':
		return [7]byte{0x00, 0x00, 0x0F, 0x11, 0x0F, 0x01, 0x0E}
	case 'i':
		return [7]byte{0x04, 0x00, 0x0C, 0x04, 0x04, 0x04, 0x0E}
	case 'l':
		return [7]byte{0x0C, 0x04, 0x04, 0x04, 0x04, 0x04, 0x0E}
	case 'm':
		return [7]byte{0x00, 0x00, 0x1A, 0x15, 0x15, 0x11, 0x11}
	case 'n':
		return [7]byte{0x00, 0x00, 0x16, 0x19, 0x11, 0x11, 0x11}
	case 'o':
		return [7]byte{0x00, 0x00, 0x0E, 0x11, 0x11, 0x11, 0x0E}
	case 'p':
		return [7]byte{0x00, 0x00, 0x1E, 0x11, 0x1E, 0x10, 0x10}
	case 'r':
		return [7]byte{0x00, 0x00, 0x16, 0x19, 0x10, 0x10, 0x10}
	case 's':
		return [7]byte{0x00, 0x00, 0x0F, 0x10, 0x0E, 0x01, 0x1E}
	case 't':
		return [7]byte{0x08, 0x08, 0x1C, 0x08, 0x08, 0x09, 0x06}
	case 'u':
		return [7]byte{0x00, 0x00, 0x11, 0x11, 0x11, 0x13, 0x0D}
	case 'x':
		return [7]byte{0x00, 0x00, 0x11, 0x0A, 0x04, 0x0A, 0x11}
	case 'y':
		return [7]byte{0x00, 0x00, 0x11, 0x11, 0x0F, 0x01, 0x0E}
	case 'z':
		return [7]byte{0x00, 0x00, 0x1F, 0x02, 0x04, 0x08, 0x1F}
	case ' ':
		return [7]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	case '.':
		return [7]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04}
	case ':':
		return [7]byte{0x00, 0x00, 0x04, 0x00, 0x00, 0x04, 0x00}
	case '(':
		return [7]byte{0x02, 0x04, 0x08, 0x08, 0x08, 0x04, 0x02}
	case ')':
		return [7]byte{0x08, 0x04, 0x02, 0x02, 0x02, 0x04, 0x08}
	case ',':
		return [7]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x04, 0x08}
	case '-':
		return [7]byte{0x00, 0x00, 0x00, 0x1F, 0x00, 0x00, 0x00}
	case '/':
		return [7]byte{0x00, 0x01, 0x02, 0x04, 0x08, 0x10, 0x00}
	default:
		return [7]byte{0x0E, 0x11, 0x11, 0x11, 0x11, 0x11, 0x0E}
	}
}
