// gamemode.go - Game mode system for different game types
package rengine

type GameMode int

const (
	GameModeFirstPerson GameMode = iota // FPS, walking sim
	GameModeTopDown                     // Racing, strategy
	GameModeIsometric                   // RPG, RTS
	GameModeSideScroll                  // Platformer
	GameModeCustom                      // User-defined
)

type Camera struct {
	Mode          GameMode
	FOV           float64  
	Height        float64  
	Pitch         float64 
	Roll          float64 
	ZoomLevel     float64 
	FollowDistance float64
	FollowHeight   float64
}

func (e *Engine) SetGameMode(mode GameMode) {
	if e.Camera == nil {
		e.Camera = &Camera{}
	}
	e.Camera.Mode = mode
	
	switch mode {
	case GameModeFirstPerson:
		e.Camera.FOV = 0.66
		e.Camera.Height = 0.0
		e.Camera.Pitch = 0.0
		e.Camera.Roll = 0.0
	case GameModeTopDown:
		e.Camera.ZoomLevel = 1.0
		e.Camera.Height = 10.0
	case GameModeIsometric:
		e.Camera.ZoomLevel = 1.0
		e.Camera.Height = 8.0
		e.Camera.Pitch = -0.785 // 45 degrees down
	case GameModeSideScroll:
		e.Camera.ZoomLevel = 1.0
		e.Camera.Height = 0.0
	}
}

func (e *Engine) GetGameMode() GameMode {
	if e.Camera == nil {
		return GameModeFirstPerson
	}
	return e.Camera.Mode
}

func (e *Engine) SetCameraHeight(height float64) {
	if e.Camera == nil {
		e.Camera = &Camera{}
	}
	e.Camera.Height = height
}

func (e *Engine) SetCameraPitch(pitch float64) {
	if e.Camera == nil {
		e.Camera = &Camera{}
	}
	e.Camera.Pitch = pitch
}

func (e *Engine) SetCameraRoll(roll float64) {
	if e.Camera == nil {
		e.Camera = &Camera{}
	}
	e.Camera.Roll = roll
}

func (e *Engine) SetZoomLevel(zoom float64) {
	if e.Camera == nil {
		e.Camera = &Camera{}
	}
	e.Camera.ZoomLevel = zoom
}

type Vehicle struct {
	Speed         float64
	MaxSpeed      float64
	Acceleration  float64
	Deceleration  float64
	TurnRate      float64
	Drift         float64
}

func (e *Engine) SetupVehicle(maxSpeed, acceleration, turnRate float64) *Vehicle {
	if e.Vehicle == nil {
		e.Vehicle = &Vehicle{
			Speed:        0.0,
			MaxSpeed:     maxSpeed,
			Acceleration: acceleration,
			Deceleration: 0.05,
			TurnRate:     turnRate,
			Drift:        0.0,
		}
	}
	return e.Vehicle
}

func (e *Engine) UpdateVehiclePhysics(accelerate, brake, turnLeft, turnRight bool) {
	if e.Vehicle == nil || e.Player == nil {
		return
	}

	if accelerate {
		e.Vehicle.Speed += e.Vehicle.Acceleration
		if e.Vehicle.Speed > e.Vehicle.MaxSpeed {
			e.Vehicle.Speed = e.Vehicle.MaxSpeed
		}
	} else if brake {
		e.Vehicle.Speed -= e.Vehicle.Deceleration * 2
		if e.Vehicle.Speed < 0 {
			e.Vehicle.Speed = 0
		}
	} else {
		e.Vehicle.Speed -= e.Vehicle.Deceleration
		if e.Vehicle.Speed < 0 {
			e.Vehicle.Speed = 0
		}
	}

	if e.Vehicle.Speed > 0.01 {
		if turnLeft {
			e.rotatePlayer(e.Vehicle.TurnRate * (e.Vehicle.Speed / e.Vehicle.MaxSpeed))
		}
		if turnRight {
			e.rotatePlayer(-e.Vehicle.TurnRate * (e.Vehicle.Speed / e.Vehicle.MaxSpeed))
		}
	}

	newPosX := e.Player.PosX + e.Player.DirX*e.Vehicle.Speed
	newPosY := e.Player.PosY + e.Player.DirY*e.Vehicle.Speed

	if e.isWalkable(int(newPosX), int(e.Player.PosY)) {
		e.Player.PosX = newPosX
	} else {
		e.Vehicle.Speed = 0
	}

	if e.isWalkable(int(e.Player.PosX), int(newPosY)) {
		e.Player.PosY = newPosY
	} else {
		e.Vehicle.Speed = 0
	}

	if e.Camera != nil {
		if turnLeft {
			e.Camera.Roll = -0.2 * (e.Vehicle.Speed / e.Vehicle.MaxSpeed)
		} else if turnRight {
			e.Camera.Roll = 0.2 * (e.Vehicle.Speed / e.Vehicle.MaxSpeed)
		} else {
			e.Camera.Roll *= 0.9
		}
	}
}

func (e *Engine) GetVehicleSpeed() float64 {
	if e.Vehicle == nil {
		return 0
	}
	return e.Vehicle.Speed
}
