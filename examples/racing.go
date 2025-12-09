// Example 2: Racing Game

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

func Racing() {
	// Create larger map for racing
	eng := rengine.NewEngine(640, 480, 50, 50, 33)

	// Create a race track (0 = track, 1 = barriers)
	trackMap := make([][]int, 50)
	for i := range trackMap {
		trackMap[i] = make([]int, 50)
		for j := range trackMap[i] {
			trackMap[i][j] = 0 // All track by default
		}
	}

	// Create outer barrier
	for i := 0; i < 50; i++ {
		trackMap[i][0] = 1
		trackMap[i][49] = 1
		trackMap[0][i] = 1
		trackMap[49][i] = 1
	}

	// Create inner barriers (make it a circuit)
	for i := 10; i < 40; i++ {
		trackMap[i][10] = 1
		trackMap[i][39] = 1
		trackMap[10][i] = 1
		trackMap[39][i] = 1
	}

	// Create some obstacles/chicanes
	for i := 15; i < 20; i++ {
		trackMap[i][25] = 1
		trackMap[30][i] = 1
	}

	eng.SetWorldMap(trackMap)

	// Start position on the track
	eng.SetPlayer(
		5.0, 25.0,   // Position (X, Y)
		1.0, 0.0,    // Direction (facing right)
		0.0, 0.66,   // Camera plane
		0.1,         // Move speed (not used in vehicle mode)
		0.05,        // Rotation speed (not used in vehicle mode)
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

	// Set game mode to first person (racing view)
	eng.SetGameMode(rengine.GameModeFirstPerson)

	// Setup vehicle physics
	// Parameters: maxSpeed, acceleration, turnRate
	eng.SetupVehicle(0.3, 0.01, 0.03)

	// Enable debug mode to see speed
	eng.EnableDebug()

	log.Println("Racing Game")
	log.Println("Controls:")
	log.Println("  W - Accelerate")
	log.Println("  S - Brake")
	log.Println("  A - Turn Left")
	log.Println("  D - Turn Right")
	log.Println("  F3 - Toggle Debug")
	log.Println("  ESC - Quit")
	log.Println("")
	log.Println("Try to complete laps around the track!")

	// Main game loop
	for eng.IsRunning() {
		eng.UpdateDebugInfo()

		// Get vehicle controls
		accelerate := eng.IsKeyPressed(glfw.KeyW)
		brake := eng.IsKeyPressed(glfw.KeyS)
		turnLeft := eng.IsKeyPressed(glfw.KeyA)
		turnRight := eng.IsKeyPressed(glfw.KeyD)

		// Update vehicle physics
		eng.UpdateVehiclePhysics(accelerate, brake, turnLeft, turnRight)

		eng.Render()
		eng.DrawFrameBuffer()
		eng.Window.SwapBuffers()
		glfw.PollEvents()
		time.Sleep(time.Duration(eng.FPS) * time.Millisecond)
	}
}
