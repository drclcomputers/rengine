package engine

import (
	"math"
	"sort"
	"time"
)

// Sprite types
type SpriteType int
const (
	SpriteTypeDecoration SpriteType = iota
	SpriteTypeEnemy
	SpriteTypeItem
)

// Enemy AI states
type EnemyState int
const (
	StateIdle EnemyState = iota
	StateChasing
	StateAttackingStill
	StateAttackingChasing
	StateDead
)

type Sprite struct {
	// Common properties
	PosX     float64
	PosY     float64
	Texture  *Texture
	Scale    float64
	VMove    float64
	Type     SpriteType
	
	// Enemy-specific properties
	Health         float64
	MaxHealth      float64
	Speed          float64
	Damage         float64
	State          EnemyState
	LastAttackTime time.Time
	AttackCooldown time.Duration
	DetectionRange float64
	AttackRange    float64
	LoseRange      float64
	
	// Item-specific properties
	ItemType       ItemType
	ItemValue      int
	ItemID         string
}

type SpriteDistance struct {
	Index    int
	Distance float64
}

func NewSprite(x, y float64, texture *Texture) *Sprite {
	return &Sprite{
		PosX:    x,
		PosY:    y,
		Texture: texture,
		Scale:   1.0,
		VMove:   0.0,
		Type:    SpriteTypeDecoration,
	}
}

func NewEnemy(x, y float64, texture *Texture, health, speed, damage float64) *Sprite {
	return &Sprite{
		PosX:           x,
		PosY:           y,
		Texture:        texture,
		Scale:          1.0,
		VMove:          0.0,
		Type:           SpriteTypeEnemy,
		Health:         health,
		MaxHealth:      health,
		Speed:          speed,
		Damage:         damage,
		State:          StateIdle,
		AttackCooldown: time.Second * 1,
		LastAttackTime: time.Now().Add(-time.Second * 2),
		DetectionRange: 10.0,
		AttackRange:    1.5,
		LoseRange:      15.0,
	}
}

func NewHealthPickup(x, y float64, texture *Texture, amount int) *Sprite {
	return &Sprite{
		PosX:      x,
		PosY:      y,
		Texture:   texture,
		Scale:     0.5,
		VMove:     -0.5,
		Type:      SpriteTypeItem,
		ItemType:  ItemTypeHealth,
		ItemValue: amount,
	}
}

func NewAmmoPickup(x, y float64, texture *Texture, amount int) *Sprite {
	return &Sprite{
		PosX:      x,
		PosY:      y,
		Texture:   texture,
		Scale:     0.5,
		VMove:     -0.5,
		Type:      SpriteTypeItem,
		ItemType:  ItemTypeAmmo,
		ItemValue: amount,
	}
}

func NewKeyPickup(x, y float64, texture *Texture, keyID string) *Sprite {
	return &Sprite{
		PosX:      x,
		PosY:      y,
		Texture:   texture,
		Scale:     0.5,
		VMove:     -0.5,
		Type:      SpriteTypeItem,
		ItemType:  ItemTypeKey,
		ItemID:    keyID,
	}
}

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
		if distToPlayer < sprite.DetectionRange {
			sprite.State = StateChasing
		}
		
	case StateChasing:
		if distToPlayer > sprite.LoseRange {
			sprite.State = StateIdle
			return
		}
		
		if distToPlayer < sprite.AttackRange {
			sprite.State = StateAttackingStill
			return
		}
		
		e.moveEnemyTowardsPlayer(sprite, deltaTime)
		
	case StateAttackingStill:
		if distToPlayer > sprite.AttackRange*1.5 {
			sprite.State = StateChasing
			return
		}
		
		e.enemyAttackPlayer(sprite)
		
	case StateAttackingChasing:
		if distToPlayer > sprite.LoseRange {
			sprite.State = StateIdle
			return
		}
		
		if distToPlayer < sprite.AttackRange {
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
	
	if IsWalkable(int(newX), int(sprite.PosY), e.MapWidth, e.MapHeight, e.WorldMap) {
		sprite.PosX = newX
	}
	if IsWalkable(int(sprite.PosX), int(newY), e.MapWidth, e.MapHeight, e.WorldMap) {
		sprite.PosY = newY
	}
}

func (e *Engine) enemyAttackPlayer(sprite *Sprite) {
	if time.Since(sprite.LastAttackTime) >= sprite.AttackCooldown {
		e.Player.Damage(sprite.Damage)
		sprite.LastAttackTime = time.Now()
	}
}

func (e *Engine) CleanupDeadEnemies() {
	for i := len(e.Sprites) - 1; i >= 0; i-- {
		sprite := e.Sprites[i]
		
		if sprite.Type == SpriteTypeEnemy && sprite.State == StateDead {
			sprite.Health -= 0.5
			if sprite.Health <= -10 {
				e.RemoveSprite(i)
			}
		}
	}
}

func (e *Engine) AddSprite(sprite *Sprite) {
	if e.Sprites == nil {
		e.Sprites = make([]*Sprite, 0)
	}
	e.Sprites = append(e.Sprites, sprite)
}

func (e *Engine) RemoveSprite(index int) {
	if index >= 0 && index < len(e.Sprites) {
		e.Sprites = append(e.Sprites[:index], e.Sprites[index+1:]...)
	}
}

func (e *Engine) RenderSprites() {
	if len(e.Sprites) == 0 {
		return
	}
	
	spriteDistances := make([]SpriteDistance, len(e.Sprites))
	for i, sprite := range e.Sprites {
		spriteDistances[i] = SpriteDistance{
			Index:    i,
			Distance: (e.Player.PosX-sprite.PosX)*(e.Player.PosX-sprite.PosX) + 
			          (e.Player.PosY-sprite.PosY)*(e.Player.PosY-sprite.PosY),
		}
	}
	
	sort.Slice(spriteDistances, func(i, j int) bool {
		return spriteDistances[i].Distance > spriteDistances[j].Distance
	})
	
	for _, sd := range spriteDistances {
		sprite := e.Sprites[sd.Index]
		e.renderSprite(sprite)
	}
}

func (e *Engine) renderSprite(sprite *Sprite) {
	spriteX := sprite.PosX - e.Player.PosX
	spriteY := sprite.PosY - e.Player.PosY
	
	invDet := 1.0 / (e.Player.PlaneX*e.Player.DirY - e.Player.DirX*e.Player.PlaneY)
	
	transformX := invDet * (e.Player.DirY*spriteX - e.Player.DirX*spriteY)
	transformY := invDet * (-e.Player.PlaneY*spriteX + e.Player.PlaneX*spriteY)
	
	if transformY <= 0 {
		return
	}
	
	spriteScreenX := int(float64(e.ScreenWidth/2) * (1 + transformX/transformY))
	
	vMoveScreen := int(-sprite.VMove*100 / transformY)
	spriteHeight := int(math.Abs(float64(e.ScreenHeight)/transformY) * sprite.Scale)
	
	aspectRatio := float64(sprite.Texture.Width) / float64(sprite.Texture.Height)
	spriteWidth := int(math.Abs(float64(e.ScreenHeight)/transformY) * sprite.Scale * aspectRatio)
	
	drawStartY := -spriteHeight/2 + e.ScreenHeight/2 + vMoveScreen
	drawStartY = int(max(0, float64(drawStartY)))
	drawEndY := spriteHeight/2 + e.ScreenHeight/2 + vMoveScreen
	if drawEndY >= e.ScreenHeight {
		drawEndY = e.ScreenHeight - 1
	}
	
	drawStartX := -spriteWidth/2 + spriteScreenX
	drawStartX = int(max(0, float64(drawStartX)))
	drawEndX := spriteWidth/2 + spriteScreenX
	if drawEndX >= e.ScreenWidth {
		drawEndX = e.ScreenWidth - 1
	}
	
	for stripe := drawStartX; stripe < drawEndX; stripe++ {
		texX := int(256*(stripe-(-spriteWidth/2+spriteScreenX))*sprite.Texture.Width/spriteWidth) / 256
		
		if transformY < e.ZBuffer[stripe] && stripe >= 0 && stripe < e.ScreenWidth {
			for y := drawStartY; y < drawEndY; y++ {
				d := (y-vMoveScreen)*256 - e.ScreenHeight*128 + spriteHeight*128
				texY := ((d * sprite.Texture.Height) / spriteHeight) / 256
				
				texY = int(max(0, float64(texY)))
				if texY >= sprite.Texture.Height {
					texY = sprite.Texture.Height - 1
				}
				
				r, g, b, a := sprite.Texture.GetPixelWithAlpha(texX, texY)
				
				if a > 128 {
					factor := 1.0 / (1.0 + transformY*0.1)
					r = byte(float64(r) * factor)
					g = byte(float64(g) * factor)
					b = byte(float64(b) * factor)
					
					idx := (y*e.ScreenWidth + stripe) * 4
					e.FrameBuffer[idx] = r
					e.FrameBuffer[idx+1] = g
					e.FrameBuffer[idx+2] = b
					e.FrameBuffer[idx+3] = 255
				}
			}
		}
	}
}

func (e *Engine) CheckItemPickups() {
	for i := len(e.Sprites) - 1; i >= 0; i-- {
		sprite := e.Sprites[i]
		
		if sprite.Type != SpriteTypeItem {
			continue
		}
		
		dx := e.Player.PosX - sprite.PosX
		dy := e.Player.PosY - sprite.PosY
		dist := math.Sqrt(dx*dx + dy*dy)
		
		if dist < 0.5 {
			e.handleItemPickup(sprite)
			e.RemoveSprite(i)
		}
	}
}

func (e *Engine) handleItemPickup(sprite *Sprite) {
	switch sprite.ItemType {
	case ItemTypeHealth:
		e.Player.Heal(float64(sprite.ItemValue))
		
	case ItemTypeAmmo:
		e.Player.AddAmmo(sprite.ItemValue)
		
	case ItemTypeKey:
		item := &Item{
			ID:   sprite.ItemID,
			Name: "Key: " + sprite.ItemID,
			Type: ItemTypeKey,
		}
		e.Player.AddItem(item)
	}
}

func (e *Engine) RaycastForSprite() *Sprite {
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
