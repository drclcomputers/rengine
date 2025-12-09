// engine.go - Core engine structure and initialization
package engine

import (
	"fmt"
	"log"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type Engine struct {
	ScreenWidth   int
	ScreenHeight  int
	MapWidth      int
	MapHeight     int
	FPS           int
	WorldMap      [][]int
	Player        *Player
	FrameBuffer   []byte
	KeyState      map[glfw.Key]bool
	Textures      []*Texture
	FloorTexture  *Texture
	CeilingTexture *Texture
	Window        *glfw.Window
	ScreenTexture uint32
	Running       bool
	SprintMultiplier float64
	DebugInfo        *DebugInfo
}

type Player struct {
	PosX          float64
	PosY          float64
	DirX          float64
	DirY          float64
	PlaneX        float64
	PlaneY        float64
	MoveSpeed     float64
	RotationSpeed float64
}

func NewEngine(width, height, mapWidth, mapHeight, fps int) *Engine {
	return &Engine{
		ScreenWidth:  width,
		ScreenHeight: height,
		MapWidth:     mapWidth,
		MapHeight:    mapHeight,
		FPS:          fps,
		KeyState:     make(map[glfw.Key]bool),
		Textures:     make([]*Texture, 0),
		Running:      false,
		SprintMultiplier: 2.0,
	}
}

func (e *Engine) SetWorldMap(worldMap [][]int) {
	e.WorldMap = worldMap
}

func (e *Engine) SetPlayer(posX, posY, dirX, dirY, planeX, planeY, moveSpeed, rotSpeed float64) {
	e.Player = &Player{
		PosX:          posX,
		PosY:          posY,
		DirX:          dirX,
		DirY:          dirY,
		PlaneX:        planeX,
		PlaneY:        planeY,
		MoveSpeed:     moveSpeed,
		RotationSpeed: rotSpeed,
	}
}

func (e *Engine) Initialize(fullscreen bool) error {
	if err := glfw.Init(); err != nil {
		return fmt.Errorf("failed to initialize glfw: %v", err)
	}

	var window *glfw.Window
	var err error

	if fullscreen {
		monitor := glfw.GetPrimaryMonitor()
		mode := monitor.GetVideoMode()
		window, err = glfw.CreateWindow(mode.Width, mode.Height, "Raycaster Engine", monitor, nil)
		fmt.Printf("Monitor Resolution: %d x %d\n", mode.Width, mode.Height)
	} else {
		window, err = glfw.CreateWindow(e.ScreenWidth, e.ScreenHeight, "Raycaster Engine", nil, nil)
	}

	if err != nil {
		return fmt.Errorf("failed to create window: %v", err)
	}

	e.Window = window
	window.MakeContextCurrent()
	window.SetKeyCallback(e.keyCallback)

	if err := gl.Init(); err != nil {
		return fmt.Errorf("failed to initialize OpenGL: %v", err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version:", version)

	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(0, float64(e.ScreenWidth), float64(e.ScreenHeight), 0, -1, 1)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
	gl.Disable(gl.DEPTH_TEST)

	e.FrameBuffer = make([]byte, e.ScreenWidth*e.ScreenHeight*4)

	e.initScreenTexture()

	e.Running = true
	return nil
}

func (e *Engine) initScreenTexture() {
	gl.GenTextures(1, &e.ScreenTexture)
	gl.BindTexture(gl.TEXTURE_2D, e.ScreenTexture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(e.ScreenWidth),
		int32(e.ScreenHeight),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(e.FrameBuffer),
	)
}

func (e *Engine) keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	switch action {
	case glfw.Press:
		e.KeyState[key] = true

		if key == glfw.KeyF3 {
			e.ToggleDebug()
		}
	case glfw.Release:
		e.KeyState[key] = false
	}

	if (key == glfw.KeyEscape && action == glfw.Press) || (key == glfw.KeyC && action == glfw.Press && (mods&glfw.ModControl) != 0) {
		e.Running = false
		w.SetShouldClose(true)
	}
}

func (e *Engine) Cleanup() {
	if e.Window != nil {
		glfw.Terminate()
	}
}

func (e *Engine) IsRunning() bool {
	return e.Running && !e.Window.ShouldClose()
}

func (e *Engine) GetTexture(index int) *Texture {
	if index >= 0 && index < len(e.Textures) {
		return e.Textures[index]
	}
	log.Printf("Warning: texture index %d out of range", index)
	return nil
}
