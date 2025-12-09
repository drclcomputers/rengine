// main.go - Example game using the raycaster engine
package examples

import (
	"log"
	"runtime"
	"time"

	"github.com/go-gl/glfw/v3.3/glfw"
	"rengine/engine"
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

	eng := engine.NewEngine(640, 480, len(worldMap), len(worldMap[1]), 25)

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
