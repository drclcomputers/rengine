// sprite.go - Sprite entities and rendering
package engine

import (
	"math"
	"sort"
)

type SpriteType int
const (
    SpriteTypeDecoration SpriteType = iota
    SpriteTypeEnemy
    SpriteTypeItem
)

type Sprite struct {
	PosX     float64
	PosY     float64
	Texture  *Texture
	Scale    float64
	VMove    float64
	Frame    int
	Type     SpriteType
	Health   float64
	Speed    float64
	Damage   float64
	State	 int  // 0 - idle, 1 - chasing, 2 - standing and attacking, 3 - chasing and attacking
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
		Frame:   0,
		Type:    SpriteTypeDecoration,
		Health:  100,
		Speed:   0.0,
		Damage:  0.0,
		State:   0,
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
	drawStartY = max(0, drawStartY)
	drawEndY := spriteHeight/2 + e.ScreenHeight/2 + vMoveScreen
	if drawEndY >= e.ScreenHeight {
		drawEndY = e.ScreenHeight - 1
	}

	drawStartX := -spriteWidth/2 + spriteScreenX
	drawStartX = max(0, drawStartX)
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
				
				texY = max(0, texY)
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

func (e *Engine) LoadSpriteTexture(path string, imageType int) (*Texture, error) {
	texture, err := e.loadTexture(path, imageType, "sprite")
	if err != nil {
		return nil, err
	}
	return texture, nil
}
