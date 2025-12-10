// Updated example.go with enemy system
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

	eng := engine.NewEngine(640, 480, len(worldMap), len(worldMap[0]), 120)

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

	// Load textures
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

	// Load sprite textures
	enemyTex, err := eng.LoadSpriteTexture("examples/pillar.jpg", engine.TextureTypeJPEG)
	if err != nil {
		log.Println("Warning: Could not load enemy sprite:", err)
	} else {
		// Create enemies
		enemy1 := engine.NewEnemy(5.5, 5.5, enemyTex, 100, 0.8, 10)
		eng.AddSprite(enemy1)

		enemy2 := engine.NewEnemy(10.5, 10.5, enemyTex, 75, 1.0, 15)
		eng.AddSprite(enemy2)

		enemy3 := engine.NewEnemy(15.5, 8.5, enemyTex, 150, 0.6, 20)
		eng.AddSprite(enemy3)
	}

	// Add decorative sprites
	pillarTex, err := eng.LoadSpriteTexture("examples/pillar.jpg", engine.TextureTypeJPEG)
	if err != nil {
		log.Println("Warning: Could not load pillar sprite:", err)
	} else {
		pillar := engine.NewSprite(7.5, 7.5, pillarTex)
		pillar.Scale = 1.5
		eng.AddSprite(pillar)
	}

	// Add health pickup
	itemTex, err := eng.LoadSpriteTexture("examples/pillar.jpg", engine.TextureTypeJPEG)
	if err != nil {
		log.Println("Warning: Could not load item sprite:", err)
	} else {
		healthPack := engine.NewItem(12.5, 5.5, itemTex)
		eng.AddSprite(healthPack)
	}

	log.Println("=== CONTROLS ===")
	log.Println("WASD - Move")
	log.Println("Q/E - Rotate")
	log.Println("SPACE - Shoot")
	log.Println("SHIFT - Sprint")
	log.Println("F3 - Toggle Debug")
	log.Println("ESC - Quit")

	lastTime := time.Now()

	for eng.IsRunning() {
		// Calculate delta time
		currentTime := time.Now()
		deltaTime := currentTime.Sub(lastTime).Seconds()
		lastTime = currentTime

		eng.UpdateDebugInfo()

		// Update enemy AI
		eng.UpdateEnemies(deltaTime)

		// Handle player input
		eng.HandleInput()

		// Clean up dead enemies
		eng.CleanupDeadEnemies()

		// Render everything
		eng.Render()

		eng.DrawFrameBuffer()

		eng.Window.SwapBuffers()
		glfw.PollEvents()


		time.Sleep(time.Duration(1000/eng.FPS))
	}
}
