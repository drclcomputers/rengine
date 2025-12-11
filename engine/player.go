// player.go 

package engine

import (
	"math"
	"time"
)

type Player struct {
    PosX, PosY, DirX, DirY, PlaneX, PlaneY float64
    MoveSpeed, RotationSpeed, SprintMultiplier float64
    
    Health, MaxHealth float64
    CurrentWeapon *Weapon
    Weapons []*Weapon
    
    Inventory []*Item
}

type Weapon struct {
	Name              string
	Damage            float64
	FireRate          time.Duration
	MagSize           int
	CurrentAmmo       int
	ReserveAmmo       int
	MaxReserveAmmo    int
	LastShotTime      time.Time
	ReloadTime        time.Duration
	IsReloading       bool
	ReloadStartTime   time.Time
	MuzzleFlashTime   time.Time
	MuzzleFlashDuration time.Duration
}

type Item struct {
	ID       string
	Name     string
	Type     ItemType
	Value    int
	Effect   func(*Player, *Engine)
}

type ItemType int
const (
	ItemTypeHealth ItemType = iota
	ItemTypeAmmo
	ItemTypeWeapon
	ItemTypeKey
	ItemTypeArmor
)

func NewPlayer(posX, posY, dirX, dirY, planeX, planeY float64) *Player {
	pistol := &Weapon{
		Name:                "Pistol",
		Damage:              25.0,
		FireRate:            time.Millisecond * 250,
		MagSize:             12,
		CurrentAmmo:         12,
		ReserveAmmo:         48,
		MaxReserveAmmo:      96,
		MuzzleFlashDuration: time.Millisecond * 50,
		ReloadTime:          time.Second * 2,
	}
	
	return &Player{
		PosX:             posX,
		PosY:             posY,
		DirX:             dirX,
		DirY:             dirY,
		PlaneX:           planeX,
		PlaneY:           planeY,
		MoveSpeed:        0.05,
		RotationSpeed:    0.05,
		SprintMultiplier: 1.0,
		Health:           100,
		MaxHealth:        100,
		CurrentWeapon:    pistol,
		Weapons:          []*Weapon{pistol},
		Inventory:        make([]*Item, 0),
	}
}

// Movement methods
func (p *Player) Move(dirX, dirY, deltaTime float64, worldMap [][]int, mapWidth, mapHeight int) {
	speed := p.MoveSpeed * p.SprintMultiplier
	newPosX := p.PosX + dirX*speed
	newPosY := p.PosY + dirY*speed

	if IsWalkable(int(newPosX), int(p.PosY), mapWidth, mapHeight, worldMap) {
		p.PosX = newPosX
	}
	if IsWalkable(int(p.PosX), int(newPosY), mapWidth, mapHeight, worldMap) {
		p.PosY = newPosY
	}
}

func (p *Player) Rotate(angle float64) {
	oldDirX := p.DirX
	p.DirX = p.DirX*math.Cos(angle) - p.DirY*math.Sin(angle)
	p.DirY = oldDirX*math.Sin(angle) + p.DirY*math.Cos(angle)

	oldPlaneX := p.PlaneX
	p.PlaneX = p.PlaneX*math.Cos(angle) - p.PlaneY*math.Sin(angle)
	p.PlaneY = oldPlaneX*math.Sin(angle) + p.PlaneY*math.Cos(angle)
}

// Combat methods
func (p *Player) Shoot() bool {
	if p.CurrentWeapon == nil {
		return false
	}
	
	if p.CurrentWeapon.IsReloading {
		return false
	}
	
	if time.Since(p.CurrentWeapon.LastShotTime) < p.CurrentWeapon.FireRate {
		return false
	}
	
	if p.CurrentWeapon.CurrentAmmo <= 0 {
		p.Reload()
		return false
	}
	
	p.CurrentWeapon.CurrentAmmo--
	p.CurrentWeapon.LastShotTime = time.Now()
	p.CurrentWeapon.MuzzleFlashTime = time.Now()
	
	return true
}

func (p *Player) Reload() {
	if p.CurrentWeapon == nil || p.CurrentWeapon.IsReloading {
		return
	}
	
	if p.CurrentWeapon.CurrentAmmo >= p.CurrentWeapon.MagSize {
		return
	}
	
	if p.CurrentWeapon.ReserveAmmo <= 0 {
		return
	}
	
	p.CurrentWeapon.IsReloading = true
	p.CurrentWeapon.ReloadStartTime = time.Now()
}

func (p *Player) UpdateReload() {
	if !p.CurrentWeapon.IsReloading {
		return
	}
	
	if time.Since(p.CurrentWeapon.ReloadStartTime) >= p.CurrentWeapon.ReloadTime {
		ammoNeeded := p.CurrentWeapon.MagSize - p.CurrentWeapon.CurrentAmmo
		ammoToReload := min(float64(ammoNeeded), float64(p.CurrentWeapon.ReserveAmmo))
		
		p.CurrentWeapon.CurrentAmmo += int(ammoToReload)
		p.CurrentWeapon.ReserveAmmo -= int(ammoToReload)
		p.CurrentWeapon.IsReloading = false
	}
}

func (p *Player) IsMuzzleFlashActive() bool {
	if p.CurrentWeapon == nil {
		return false
	}
	return time.Since(p.CurrentWeapon.MuzzleFlashTime) < p.CurrentWeapon.MuzzleFlashDuration
}

func (p *Player) Damage(amount float64) {
	p.Health -= amount
	if p.Health < 0 {
		p.Health = 0
	}
}

func (p *Player) Heal(amount float64) {
	p.Health += amount
	if p.Health > p.MaxHealth {
		p.Health = p.MaxHealth
	}
}

func (p *Player) AddAmmo(amount int) {
	if p.CurrentWeapon == nil {
		return
	}
	p.CurrentWeapon.ReserveAmmo += amount
	if p.CurrentWeapon.ReserveAmmo > p.CurrentWeapon.MaxReserveAmmo {
		p.CurrentWeapon.ReserveAmmo = p.CurrentWeapon.MaxReserveAmmo
	}
}

func (p *Player) AddItem(item *Item) {
	p.Inventory = append(p.Inventory, item)
}

func (p *Player) HasKey(keyID string) bool {
	for _, item := range p.Inventory {
		if item.Type == ItemTypeKey && item.ID == keyID {
			return true
		}
	}
	return false
}

func (p *Player) SwitchWeapon(index int) {
	if index >= 0 && index < len(p.Weapons) {
		p.CurrentWeapon = p.Weapons[index]
	}
}
