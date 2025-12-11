package engine

import (
	"fmt"
	"log"
	"time"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type Engine struct {
	// Display
	ScreenWidth   int
	ScreenHeight  int
	FrameBuffer   []byte
	ScreenTexture uint32
	Window        *glfw.Window
	
	// World
	MapWidth      int
	MapHeight     int
	WorldMap      [][]int
	
	// Player
	Player        *Player
	
	// Rendering
	Textures      []*Texture
	FloorTexture  *Texture
	CeilingTexture *Texture
	ZBuffer       []float64
	
	// Entities
	Sprites       []*Sprite
	Doors         []*Door
	
	// Input
	KeyState      map[glfw.Key]bool
	
	// Game state
	Running       bool
	FPS           int
	LastUpdateTime time.Time
	
	// Debug
	DebugInfo     *DebugInfo
	
	// HUD
	HitMarkerTime     time.Time
	HitMarkerDuration time.Duration
	MinimapEnabled    bool
	MinimapCorner     int // 0=TL, 1=TR, 2=BR, 3=BL
}

func NewEngine(width, height, mapWidth, mapHeight, fps int) *Engine {
	return &Engine{
		ScreenWidth:       width,
		ScreenHeight:      height,
		MapWidth:          mapWidth,
		MapHeight:         mapHeight,
		FPS:               fps,
		KeyState:          make(map[glfw.Key]bool),
		Textures:          make([]*Texture, 0),
		Sprites:           make([]*Sprite, 0),
		Doors:             make([]*Door, 0),
		ZBuffer:           make([]float64, width),
		Running:           false,
		LastUpdateTime:    time.Now(),
		HitMarkerDuration: time.Millisecond * 200,
		MinimapEnabled:    true,
		MinimapCorner:     1, // Top-right by default
	}
}

func (e *Engine) SetWorldMap(worldMap [][]int) {
	e.WorldMap = worldMap
}

func (e *Engine) SetPlayer(player *Player) {
	e.Player = player
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
