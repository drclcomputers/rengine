// interacrive.go

package engine

import "time"

type DoorState int
const (
	DoorClosed DoorState = iota
	DoorOpening
	DoorOpen
	DoorClosing
)

type Door struct {
	MapX        int
	MapY        int
	State       DoorState
	OpenAmount  float64 // 0.0 = closed, 1.0 = fully open
	Locked      bool
	KeyID       string
	OpenSpeed   float64
	CloseDelay  time.Duration
	OpenTime    time.Time
}

func NewDoor(mapX, mapY int, locked bool, keyID string) *Door {
	return &Door{
		MapX:       mapX,
		MapY:       mapY,
		State:      DoorClosed,
		OpenAmount: 0.0,
		Locked:     locked,
		KeyID:      keyID,
		OpenSpeed:  2.0,
		CloseDelay: time.Second * 3,
	}
}

func (e *Engine) AddDoor(door *Door) {
	if e.Doors == nil {
		e.Doors = make([]*Door, 0)
	}
	e.Doors = append(e.Doors, door)
}

func (e *Engine) UpdateDoors(deltaTime float64) {
	for _, door := range e.Doors {
		switch door.State {
		case DoorOpening:
			door.OpenAmount += door.OpenSpeed * deltaTime
			if door.OpenAmount >= 1.0 {
				door.OpenAmount = 1.0
				door.State = DoorOpen
				door.OpenTime = time.Now()
			}
			
		case DoorOpen:
			if time.Since(door.OpenTime) >= door.CloseDelay {
				// Check if player is in doorway
				playerInDoor := int(e.Player.PosX) == door.MapX && int(e.Player.PosY) == door.MapY
				if !playerInDoor {
					door.State = DoorClosing
				}
			}
			
		case DoorClosing:
			door.OpenAmount -= door.OpenSpeed * deltaTime
			if door.OpenAmount <= 0.0 {
				door.OpenAmount = 0.0
				door.State = DoorClosed
			}
		}
	}
}

func (e *Engine) TryOpenDoor(mapX, mapY int) bool {
	for _, door := range e.Doors {
		if door.MapX == mapX && door.MapY == mapY {
			if door.Locked {
				if e.Player.HasKey(door.KeyID) {
					door.Locked = false
					door.State = DoorOpening
					return true
				}
				return false
			}
			
			if door.State == DoorClosed {
				door.State = DoorOpening
				return true
			}
		}
	}
	return false
}

func (e *Engine) GetDoorAt(mapX, mapY int) *Door {
	for _, door := range e.Doors {
		if door.MapX == mapX && door.MapY == mapY {
			return door
		}
	}
	return nil
}

func (e *Engine) IsDoorOpen(mapX, mapY int) bool {
	door := e.GetDoorAt(mapX, mapY)
	if door == nil {
		return false
	}
	return door.OpenAmount > 0.5
}

func (e *Engine) HandleUseKey() {
	rayDist := 0.0
	stepSize := 0.1
	maxDist := 2.0
	
	for rayDist < maxDist {
		checkX := int(e.Player.PosX + e.Player.DirX*rayDist)
		checkY := int(e.Player.PosY + e.Player.DirY*rayDist)
		
		if checkX < 0 || checkX >= e.MapWidth || checkY < 0 || checkY >= e.MapHeight {
			break
		}
		
		if e.WorldMap[checkX][checkY] == 4 {
			e.TryOpenDoor(checkX, checkY)
			return
		}
		
		rayDist += stepSize
	}
}
