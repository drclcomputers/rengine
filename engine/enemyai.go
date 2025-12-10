// enemyai.go - Enemy AI and behavior system
package engine

import (
	"math"
)

const (
	StateIdle             = 0
	StateChasing          = 1
	StateAttackingStill   = 2
	StateAttackingChasing = 3
	StateDead             = 4
)

const (
	EnemyDetectionRange = 10.0
	EnemyAttackRange    = 1.5
	EnemyLosePlayerRange = 15.0
)

func (e *Engine) UpdateEnemies(deltaTime float64) {
	for _, sprite := range e.Sprites {
		if sprite.Type != SpriteTypeEnemy {
			continue
		}

		if sprite.Health <= 0 {
			sprite.State = StateDead
			continue
		}

		e.updateEnemyBehavior(sprite, deltaTime)
	}
}

func (e *Engine) updateEnemyBehavior(sprite *Sprite, deltaTime float64) {
	dx := e.Player.PosX - sprite.PosX
	dy := e.Player.PosY - sprite.PosY
	distToPlayer := math.Sqrt(dx*dx + dy*dy)

	switch sprite.State {
	case StateIdle:
		if distToPlayer < EnemyDetectionRange {
			sprite.State = StateChasing
		}

	case StateChasing:
		if distToPlayer > EnemyLosePlayerRange {
			sprite.State = StateIdle
			return
		}

		if distToPlayer < EnemyAttackRange {
			sprite.State = StateAttackingStill
			return
		}

		e.moveEnemyTowardsPlayer(sprite, deltaTime)

	case StateAttackingStill:
		if distToPlayer > EnemyAttackRange*1.5 {
			sprite.State = StateChasing
			return
		}

		e.enemyAttackPlayer(sprite)

	case StateAttackingChasing:
		if distToPlayer > EnemyLosePlayerRange {
			sprite.State = StateIdle
			return
		}

		if distToPlayer < EnemyAttackRange {
			e.enemyAttackPlayer(sprite)
		}

		e.moveEnemyTowardsPlayer(sprite, deltaTime)
	}
}

func (e *Engine) moveEnemyTowardsPlayer(sprite *Sprite, deltaTime float64) {
	dx := e.Player.PosX - sprite.PosX
	dy := e.Player.PosY - sprite.PosY
	dist := math.Sqrt(dx*dx + dy*dy)

	if dist < 0.1 {
		return 
	}

	dirX := dx / dist
	dirY := dy / dist

	moveSpeed := sprite.Speed * deltaTime
	newX := sprite.PosX + dirX*moveSpeed
	newY := sprite.PosY + dirY*moveSpeed

	if e.isWalkable(int(newX), int(sprite.PosY)) {
		sprite.PosX = newX
	}
	if e.isWalkable(int(sprite.PosX), int(newY)) {
		sprite.PosY = newY
	}
}

func (e *Engine) enemyAttackPlayer(sprite *Sprite) {
	// TODO: Add attack cooldown to prevent constant damage
	// For now, you'd implement damage dealing to player here
	// Example: e.Player.Health -= sprite.Damage * deltaTime
}

func NewEnemy(x, y float64, texture *Texture, health, speed, damage float64) *Sprite {
	return &Sprite{
		PosX:    x,
		PosY:    y,
		Texture: texture,
		Scale:   1.0,
		VMove:   0.0,
		Frame:   0,
		Type:    SpriteTypeEnemy,
		Health:  health,
		Speed:   speed,
		Damage:  damage,
		State:   StateIdle,
	}
}

func NewItem(x, y float64, texture *Texture) *Sprite {
	return &Sprite{
		PosX:    x,
		PosY:    y,
		Texture: texture,
		Scale:   0.5,
		VMove:   -0.5,
		Frame:   0,
		Type:    SpriteTypeItem,
		Health:  0,
		Speed:   0,
		Damage:  0,
		State:   0,
	}
}
