// Example 3: Hill Climb Racing

package examples

import (
	"log"
	"runtime"
	"time"

	"github.com/go-gl/glfw/v3.3/glfw"
	rengine "rengine/engine"
)

func init() {
	runtime.LockOSThread()
}

func Hillclimb() {
	// Create map
	eng := rengine.NewEngine(640, 480, 40, 40, 33)

	// Create an open terrain map
	terrainMap := make([][]int, 40)
	for i := range terrainMap {
		terrainMap[i] = make([]int, 40)
		for j := range terrainMap[i] {
			terrainMap[i][j] = 0 // Open terrain
		}
	}

	// Add some boundary walls
	for i := 0; i < 40; i++ {
		terrainMap[i][0] = 1
		terrainMap[i][39] = 1
		terrainMap[0][i] = 1
		terrainMap[39][i] = 1
	}

	// Add some obstacle rocks
	terrainMap[15][15] = 1
	terrainMap[16][15] = 1
	terrainMap[25][20] = 1
	terrainMap[26][20] = 1
	terrainMap[20][30] = 1
	terrainMap[21][30] = 1

	eng.SetWorldMap(terrainMap)

	// Start position
	eng.SetPlayer(
		5.0, 5.0,    // Position
		1.0, 0.0,    // Direction
		0.0, 0.66,   // Camera plane
		0.1,         // Move speed
		0.03,        // Rotation speed
	)

	// Initialize engine
	if err := eng.Initialize(true); err != nil {
		log.Fatalln(err)
	}
	defer eng.Cleanup()

	// Load textures
	_, err := eng.LoadFloorTexture("examples/floor.jpg", rengine.TextureTypeJPEG)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = eng.LoadCeilingTexture("examples/ceiling.png", rengine.TextureTypePNG)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = eng.LoadTexture("examples/wall.jpg", rengine.TextureTypeJPEG)
	if err != nil {
		log.Fatalln(err)
	}

	// Set first person racing mode
	eng.SetGameMode(rengine.GameModeFirstPerson)

	// Setup vehicle
	eng.SetupVehicle(0.2, 0.008, 0.02)

	// Enable height map
	eng.SetupHeightMap()

	// Create various terrain features
	
	// Hill 1 - Gentle slope up
	eng.CreateRamp(8, 8, 15, 8, 0.0, 3.0)
	
	// Hill 2 - Coming down
	eng.CreateRamp(15, 8, 22, 8, 3.0, 0.0)
	
	// Hill 3 - Big mountain
	for x := 10; x < 20; x++ {
		for y := 15; y < 25; y++ {
			// Create a pyramid-like hill
			dx := float64(x - 15)
			dy := float64(y - 20)
			dist := dx*dx + dy*dy
			height := 5.0 - dist/10.0
			if height < 0 {
				height = 0
			}
			eng.SetCellHeight(x, y, height)
		}
	}
	
	// Stairs going up a cliff
	eng.CreateStairs(25, 10, 30, 10, 0.0, 4.0)
	
	// Steep ramp
	eng.CreateRamp(5, 25, 10, 35, 0.0, 5.0)

	// Enable debug
	eng.EnableDebug()

	log.Println("Hill Climb Racing Game")
	log.Println("Controls:")
	log.Println("  W - Accelerate")
	log.Println("  S - Brake")
	log.Println("  A - Turn Left")
	log.Println("  D - Turn Right")
	log.Println("  F3 - Toggle Debug")
	log.Println("  ESC - Quit")
	log.Println("")
	log.Println("Navigate the hilly terrain!")
	log.Println("Watch out - steep slopes will slow you down!")

	// Main game loop
	for eng.IsRunning() {
		eng.UpdateDebugInfo()

		// Vehicle controls
		accelerate := eng.IsKeyPressed(glfw.KeyW)
		brake := eng.IsKeyPressed(glfw.KeyS)
		turnLeft := eng.IsKeyPressed(glfw.KeyA)
		turnRight := eng.IsKeyPressed(glfw.KeyD)

		// Update vehicle
		eng.UpdateVehiclePhysics(accelerate, brake, turnLeft, turnRight)

		// Adjust camera based on terrain height
		groundHeight := eng.GetGroundHeight(eng.Player.PosX, eng.Player.PosY)
		eng.SetCameraHeight(groundHeight + 0.5) // Half unit above ground

		eng.Render()
		eng.DrawFrameBuffer()
		eng.Window.SwapBuffers()
		glfw.PollEvents()
		time.Sleep(time.Duration(eng.FPS) * time.Millisecond)
	}
}
