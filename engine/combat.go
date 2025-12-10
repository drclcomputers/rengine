// combat.go - Fixed combat and shooting mechanics with effects
package engine

import (
	"math"
	"time"
)

type Combat struct {
	PlayerHealth      float64
	PlayerMaxHealth   float64
	Ammo              int
	MaxAmmo           int
	LastShotTime      time.Time
	ShotCooldown      time.Duration
	Damage            float64
	MuzzleFlashTime   time.Time
	MuzzleFlashDuration time.Duration
	HitMarkerTime     time.Time
	HitMarkerDuration time.Duration
}

func NewCombat() *Combat {
	return &Combat{
		PlayerHealth:        100,
		PlayerMaxHealth:     100,
		Ammo:                50,
		MaxAmmo:             100,
		ShotCooldown:        time.Millisecond * 250,
		Damage:              25.0,
		MuzzleFlashDuration: time.Millisecond * 50,
		HitMarkerDuration:   time.Millisecond * 200,
	}
}

func (e *Engine) PlayerShoot() {
	if e.Combat == nil {
		return
	}

	if time.Since(e.Combat.LastShotTime) < e.Combat.ShotCooldown {
		return
	}

	if e.Combat.Ammo <= 0 {
		return
	}

	e.Combat.Ammo--
	e.Combat.LastShotTime = time.Now()
	e.Combat.MuzzleFlashTime = time.Now()

	hitSprite := e.raycastForSprite()
	
	if hitSprite != nil && hitSprite.Type == SpriteTypeEnemy {
		hitSprite.Health -= e.Combat.Damage
		e.Combat.HitMarkerTime = time.Now() 
		
		if hitSprite.Health <= 0 {
			hitSprite.State = StateDead
		}
	}
}

func (e *Engine) raycastForSprite() *Sprite {
	var closestSprite *Sprite
	closestDist := math.MaxFloat64

	for _, sprite := range e.Sprites {
		if sprite.Type != SpriteTypeEnemy || sprite.State == StateDead {
			continue
		}

		spriteX := sprite.PosX - e.Player.PosX
		spriteY := sprite.PosY - e.Player.PosY
		
		dist := math.Sqrt(spriteX*spriteX + spriteY*spriteY)
		
		if dist < 0.001 {
			continue
		}
		
		normSpriteX := spriteX / dist
		normSpriteY := spriteY / dist

		dot := normSpriteX*e.Player.DirX + normSpriteY*e.Player.DirY
		
		if dot > 0.95 && dist < closestDist {
			closestDist = dist
			closestSprite = sprite
		}
	}

	return closestSprite
}

func (e *Engine) DamagePlayer(damage float64) {
	if e.Combat == nil {
		return
	}

	e.Combat.PlayerHealth -= damage
	if e.Combat.PlayerHealth < 0 {
		e.Combat.PlayerHealth = 0
		e.Running = false
	}
}

func (e *Engine) HealPlayer(amount float64) {
	if e.Combat == nil {
		return
	}

	e.Combat.PlayerHealth += amount
	if e.Combat.PlayerHealth > e.Combat.PlayerMaxHealth {
		e.Combat.PlayerHealth = e.Combat.PlayerMaxHealth
	}
}

func (e *Engine) AddAmmo(amount int) {
	if e.Combat == nil {
		return
	}

	e.Combat.Ammo += amount
	if e.Combat.Ammo > e.Combat.MaxAmmo {
		e.Combat.Ammo = e.Combat.MaxAmmo
	}
}

func (e *Engine) CheckItemPickup() {
	for i := len(e.Sprites) - 1; i >= 0; i-- {
		sprite := e.Sprites[i]
		
		if sprite.Type != SpriteTypeItem {
			continue
		}

		dx := e.Player.PosX - sprite.PosX
		dy := e.Player.PosY - sprite.PosY
		dist := math.Sqrt(dx*dx + dy*dy)

		if dist < 0.5 { 
			e.HealPlayer(25)
			e.RemoveSprite(i)
		}
	}
}

func (e *Engine) CleanupDeadEnemies() {
	for i := len(e.Sprites) - 1; i >= 0; i-- {
		sprite := e.Sprites[i]
		
		if sprite.Type == SpriteTypeEnemy && sprite.State == StateDead {
			sprite.Health -= 0.5 
			if sprite.Health <= 10 {
				e.RemoveSprite(i)
			}
		}
	}
}

func (e *Engine) IsMuzzleFlashActive() bool {
	if e.Combat == nil {
		return false
	}
	return time.Since(e.Combat.MuzzleFlashTime) < e.Combat.MuzzleFlashDuration
}

func (e *Engine) IsHitMarkerActive() bool {
	if e.Combat == nil {
		return false
	}
	return time.Since(e.Combat.HitMarkerTime) < e.Combat.HitMarkerDuration
}
