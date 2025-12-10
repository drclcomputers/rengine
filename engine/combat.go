// combat.go - Combat and shooting mechanics
package engine

import (
	"math"
	"time"
)

type Combat struct {
	PlayerHealth    float64
	PlayerMaxHealth float64
	Ammo            int
	MaxAmmo         int
	LastShotTime    time.Time
	ShotCooldown    time.Duration
	Damage          float64
}

func NewCombat() *Combat {
	return &Combat{
		PlayerHealth:    100,
		PlayerMaxHealth: 100,
		Ammo:            50,
		MaxAmmo:         100,
		ShotCooldown:    time.Millisecond * 250, 
		Damage:          25.0,
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

	hitSprite := e.raycastForSprite()
	
	if hitSprite != nil && hitSprite.Type == SpriteTypeEnemy {
		hitSprite.Health -= e.Combat.Damage
		
		if hitSprite.Health <= 0 {
			hitSprite.State = StateDead
			// You could add death animation or delay before removing
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
		spriteX /= dist
		spriteY /= dist

		dot := spriteX*e.Player.DirX + spriteY*e.Player.DirY
		
		// dot > 0.95 means within ~18 degree cone
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
		// Handle player death
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
			// Handle pickup based on item type (you could add ItemSubType field)
			// For now, assume all items are health packs
			e.HealPlayer(25)
			e.RemoveSprite(i)
		}
	}
}

func (e *Engine) CleanupDeadEnemies() {
	for i := len(e.Sprites) - 1; i >= 0; i-- {
		sprite := e.Sprites[i]
		
		if sprite.Type == SpriteTypeEnemy && sprite.State == StateDead && sprite.Health <= -100 {
			e.RemoveSprite(i)
		}
	}
}
