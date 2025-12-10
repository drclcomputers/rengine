// main.go - Example game using the raycaster engine
package examples

import (
	"log"
	"runtime"
	"time"

	"rengine/engine"

	"github.com/go-gl/glfw/v3.3/glfw"
)

func init() {
	runtime.LockOSThread()
}

func Example() {
	worldMap := [][]int{
    	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	    {1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	    {1, 0, 2, 2, 2, 0, 0, 0, 3, 3, 0, 0, 0, 3, 3, 3, 0, 0, 2, 2, 0, 0, 0, 1},
	    {1, 0, 2, 0, 2, 0, 0, 0, 3, 0, 0, 0, 0, 3, 0, 3, 0, 0, 2, 0, 0, 0, 0, 1},
	    {1, 0, 2, 0, 2, 0, 0, 0, 3, 3, 0, 0, 0, 3, 3, 3, 0, 0, 2, 2, 0, 0, 0, 1},
	    {1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	    {1, 0, 0, 3, 3, 3, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 1},
	    {1, 0, 0, 3, 0, 3, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 1},
	    {1, 0, 0, 3, 3, 3, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 1},
	    {1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	    {1, 2, 2, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1},
	    {1, 2, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 1, 0, 1, 0, 0, 0, 1, 0, 1, 1},
	    {1, 2, 2, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1},
	    {1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	    {1, 3, 3, 3, 0, 0, 0, 2, 2, 2, 0, 0, 0, 3, 3, 3, 0, 0, 0, 2, 2, 2, 0, 1},
	    {1, 3, 0, 3, 0, 0, 0, 2, 0, 2, 0, 0, 0, 3, 0, 3, 0, 0, 0, 2, 0, 2, 0, 1},
	    {1, 3, 3, 3, 0, 0, 0, 2, 2, 2, 0, 0, 0, 3, 3, 3, 0, 0, 0, 2, 2, 2, 0, 1},
	    {1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	    {1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	}

	eng := engine.NewEngine(640, 480, len(worldMap), len(worldMap[1]), 24)

	eng.SetWorldMap(worldMap)

	eng.SetPlayer(
		1.5, 1.5,
		-1.0, 0.0,
		0.0, 0.66,
		0.05,
		0.05,
	)

	if err := eng.Initialize(true); err != nil {
		log.Fatalln(err)
	}
	defer eng.Cleanup()

	_, err := eng.LoadFloorTexture("examples/floor.jpg", engine.TextureTypeJPEG)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = eng.LoadCeilingTexture("examples/ceiling.jpg", engine.TextureTypeJPEG)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = eng.LoadTexture("examples/wall.jpg", engine.TextureTypeJPEG)
	if err != nil {
		log.Fatalln(err)
	}
	_, err = eng.LoadTexture("examples/wall2.jpg", engine.TextureTypeJPEG)
	if err != nil {
		log.Fatalln(err)
	}
	_, err = eng.LoadTexture("examples/wall3.jpg", engine.TextureTypeJPEG)
	if err != nil {
		log.Fatalln(err)
	}

	eng.EnableDebug()

	barrelTex, err := eng.LoadSpriteTexture("examples/pillar.jpg", engine.TextureTypeJPEG)
	if err != nil {
		log.Println("Warning: Could not load barrel sprite:", err)
	} else {
		b1 := engine.NewSprite(5.5, 5.5, barrelTex)
		b1.Scale = 0.5
		b1.VMove = -1.5
		eng.AddSprite(b1)
	}

	pillarTex, err := eng.LoadSpriteTexture("examples/pillar.jpg", engine.TextureTypeJPEG)
	if err != nil {
		log.Println("Warning: Could not load pillar sprite:", err)
	} else {
		pillar := engine.NewSprite(10.5, 10.5, pillarTex)
		pillar.Scale = 1.5
		eng.AddSprite(pillar)
	}

	for eng.IsRunning() {
		eng.UpdateDebugInfo()

		eng.HandleInput()

		eng.Render()

		eng.DrawFrameBuffer()

		eng.Window.SwapBuffers()
		glfw.PollEvents()

		time.Sleep(time.Duration(eng.FPS) * time.Millisecond)
	}
}
