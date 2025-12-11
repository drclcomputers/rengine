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
        {1, 1, 1, 1, 1, 1, 1, 1},
        {1, 0, 0, 0, 4, 0, 0, 1},  // 4 = door
        {1, 0, 2, 0, 0, 0, 0, 1},
        {1, 0, 0, 0, 0, 3, 0, 1},
        {1, 0, 0, 0, 4, 0, 0, 1},  // Another door
        {1, 0, 0, 0, 0, 0, 0, 1},
        {1, 0, 0, 0, 0, 0, 0, 1},
        {1, 1, 1, 1, 1, 1, 1, 1},
    }
    
    eng := engine.NewEngine(640, 480, len(worldMap), len(worldMap[0]), 60)
    eng.SetWorldMap(worldMap)
    
    // Create player
    player := engine.NewPlayer(1.5, 1.5, -1.0, 0.0, 0.0, 0.66)
    eng.SetPlayer(player)
    
    if err := eng.Initialize(true); err != nil {
        log.Fatal(err)
    }
    defer eng.Cleanup()
    
    // Load textures
    eng.LoadFloorTexture("assets/floor.jpg", engine.TextureTypeJPEG)
    eng.LoadCeilingTexture("assets/ceiling.jpg", engine.TextureTypeJPEG)
    eng.LoadTexture("assets/wall1.jpg", engine.TextureTypeJPEG)
    eng.LoadTexture("assets/wall2.jpg", engine.TextureTypeJPEG)
    eng.LoadTexture("assets/wall3.jpg", engine.TextureTypeJPEG)
    eng.LoadTexture("assets/barrel.jpg", engine.TextureTypeJPEG)
    
    // Add doors
    door1 := engine.NewDoor(4, 1, true, "red_key")
    eng.AddDoor(door1)
    
    door2 := engine.NewDoor(4, 4, false, "")  // Unlocked
    eng.AddDoor(door2)
    
    // Load sprite textures
    enemyTex, _ := eng.LoadTexture("assets/enemy.png", engine.TextureTypePNG)
    itemTex, _ := eng.LoadTexture("assets/item.png", engine.TextureTypePNG)
    keyTex, _ := eng.LoadTexture("assets/key.png", engine.TextureTypePNG)
    
    // Add enemies
    enemy1 := engine.NewEnemy(5.5, 5.5, enemyTex, 100, 0.8, 10)
    eng.AddSprite(enemy1)
    
    // Add items
    health := engine.NewHealthPickup(3.5, 3.5, itemTex, 25)
    eng.AddSprite(health)
    
    ammo := engine.NewAmmoPickup(6.5, 2.5, itemTex, 30)
    eng.AddSprite(ammo)
    
    key := engine.NewKeyPickup(2.5, 5.5, keyTex, "red_key")
    eng.AddSprite(key)
    
    eng.EnableDebug()
    
    log.Println("=== CONTROLS ===")
    log.Println("WASD - Move")
    log.Println("Q/E - Rotate")
    log.Println("SPACE - Shoot")
    log.Println("R - Reload")
    log.Println("F - Use/Open Door")
    log.Println("M - Toggle Minimap")
    log.Println("N - Change Minimap Corner")
    log.Println("1/2/3 - Switch Weapon")
    log.Println("SHIFT - Sprint")
    log.Println("F3 - Toggle Debug")
    
    lastTime := time.Now()
    
    for eng.IsRunning() {
        currentTime := time.Now()
        deltaTime := currentTime.Sub(lastTime).Seconds()
        lastTime = currentTime
        
        eng.UpdateDebugInfo()
        eng.UpdateEnemies(deltaTime)
        eng.UpdateDoors(deltaTime)
        eng.HandleInput()
        eng.Player.UpdateReload()
        eng.CleanupDeadEnemies()
        
        eng.Render()
        eng.DrawFrameBuffer()
        
        eng.Window.SwapBuffers()
        glfw.PollEvents()
    }
}
