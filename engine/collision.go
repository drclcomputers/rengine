// collision.go

package engine

import (
	"math"
)

type ColliderType int
const (
	ColliderTypeCircle ColliderType = iota
	ColliderTypeAABB
	ColliderTypePoint
)

type Collider struct {
	PosX     float64
	PosY     float64
	Radius   float64 
	Width    float64
	Height   float64 
	Type     ColliderType
	Solid    bool   
}

type CollisionResult struct {
	Hit         bool
	Distance    float64
	Normal      [2]float64 
	Penetration float64   
	ContactX    float64  
	ContactY    float64
}

func CheckCircleCollision(c1, c2 *Collider) bool {
	dx := c1.PosX - c2.PosX
	dy := c1.PosY - c2.PosY
	distSq := dx*dx + dy*dy
	radiusSum := c1.Radius + c2.Radius
	return distSq < radiusSum*radiusSum
}

func GetCircleCollisionResult(c1, c2 *Collider) CollisionResult {
	result := CollisionResult{Hit: false}
	
	dx := c2.PosX - c1.PosX
	dy := c2.PosY - c1.PosY
	dist := math.Sqrt(dx*dx + dy*dy)
	radiusSum := c1.Radius + c2.Radius
	
	if dist < radiusSum {
		result.Hit = true
		result.Distance = dist
		result.Penetration = radiusSum - dist
		
		if dist > 0.0001 {
			result.Normal[0] = dx / dist
			result.Normal[1] = dy / dist
		} else {
			result.Normal[0] = 1.0
			result.Normal[1] = 0.0
		}
		
		result.ContactX = c1.PosX + result.Normal[0]*(c1.Radius-result.Penetration/2)
		result.ContactY = c1.PosY + result.Normal[1]*(c1.Radius-result.Penetration/2)
	}
	
	return result
}

func ResolveCircleCollision(c1, c2 *Collider) {
	result := GetCircleCollisionResult(c1, c2)
	if !result.Hit {
		return
	}
	
	halfPenetration := result.Penetration / 2.0
	c1.PosX -= result.Normal[0] * halfPenetration
	c1.PosY -= result.Normal[1] * halfPenetration
	c2.PosX += result.Normal[0] * halfPenetration
	c2.PosY += result.Normal[1] * halfPenetration
}

func CheckPointInCircle(px, py float64, circle *Collider) bool {
	dx := px - circle.PosX
	dy := py - circle.PosY
	distSq := dx*dx + dy*dy
	return distSq < circle.Radius*circle.Radius
}

type RayWallIntersection struct {
	Hit       bool
	Distance  float64
	MapX      int
	MapY      int
	Side      int   
	WallX     float64
}

func CheckRayWall(startX, startY, dirX, dirY float64, worldMap [][]int, mapWidth, mapHeight int, maxDist float64) RayWallIntersection {
	result := RayWallIntersection{Hit: false}
	
	if dirX == 0 && dirY == 0 {
		return result
	}
	
	mapX := int(startX)
	mapY := int(startY)
	
	deltaDistX := math.Abs(1 / dirX)
	deltaDistY := math.Abs(1 / dirY)
	
	var sideDistX, sideDistY float64
	var stepX, stepY int
	
	if dirX < 0 {
		stepX = -1
		sideDistX = (startX - float64(mapX)) * deltaDistX
	} else {
		stepX = 1
		sideDistX = (float64(mapX) + 1.0 - startX) * deltaDistX
	}
	
	if dirY < 0 {
		stepY = -1
		sideDistY = (startY - float64(mapY)) * deltaDistY
	} else {
		stepY = 1
		sideDistY = (float64(mapY) + 1.0 - startY) * deltaDistY
	}
	
	hit := false
	side := 0
	
	for !hit {
		if sideDistX < sideDistY {
			sideDistX += deltaDistX
			mapX += stepX
			side = 0
		} else {
			sideDistY += deltaDistY
			mapY += stepY
			side = 1
		}
		
		if mapX < 0 || mapX >= mapWidth || mapY < 0 || mapY >= mapHeight {
			return result
		}
		
		if worldMap[mapX][mapY] > 0 {
			hit = true
		}
		
		var perpWallDist float64
		if side == 0 {
			perpWallDist = (float64(mapX) - startX + (1-float64(stepX))/2) / dirX
		} else {
			perpWallDist = (float64(mapY) - startY + (1-float64(stepY))/2) / dirY
		}
		
		if perpWallDist > maxDist {
			return result
		}
	}
	
	var perpWallDist float64
	if side == 0 {
		perpWallDist = (float64(mapX) - startX + (1-float64(stepX))/2) / dirX
	} else {
		perpWallDist = (float64(mapY) - startY + (1-float64(stepY))/2) / dirY
	}
	
	result.Hit = true
	result.Distance = perpWallDist
	result.MapX = mapX
	result.MapY = mapY
	result.Side = side
	
	if side == 0 {
		result.WallX = startY + perpWallDist*dirY
	} else {
		result.WallX = startX + perpWallDist*dirX
	}
	result.WallX -= math.Floor(result.WallX)
	
	return result
}

func ResolveCircleWallCollision(collider *Collider, worldMap [][]int, mapWidth, mapHeight int) {
	mapX := int(collider.PosX)
	mapY := int(collider.PosY)
	
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			checkX := mapX + dx
			checkY := mapY + dy
			
			if checkX < 0 || checkX >= mapWidth || checkY < 0 || checkY >= mapHeight {
				continue
			}
			
			if worldMap[checkX][checkY] > 0 {
				resolveCircleAABB(collider, float64(checkX), float64(checkY), 1.0, 1.0)
			}
		}
	}
}

func resolveCircleAABB(circle *Collider, boxX, boxY, boxW, boxH float64) {
	closestX := clamp(circle.PosX, boxX, boxX+boxW)
	closestY := clamp(circle.PosY, boxY, boxY+boxH)
	
	dx := circle.PosX - closestX
	dy := circle.PosY - closestY
	distSq := dx*dx + dy*dy
	
	if distSq < circle.Radius*circle.Radius {
		dist := math.Sqrt(distSq)
		
		if dist < 0.0001 {
			edgeDistLeft := circle.PosX - boxX
			edgeDistRight := (boxX + boxW) - circle.PosX
			edgeDistTop := circle.PosY - boxY
			edgeDistBottom := (boxY + boxH) - circle.PosY
			
			minDist := math.Min(math.Min(edgeDistLeft, edgeDistRight), math.Min(edgeDistTop, edgeDistBottom))
			
			if minDist == edgeDistLeft {
				circle.PosX = boxX - circle.Radius
			} else if minDist == edgeDistRight {
				circle.PosX = boxX + boxW + circle.Radius
			} else if minDist == edgeDistTop {
				circle.PosY = boxY - circle.Radius
			} else {
				circle.PosY = boxY + boxH + circle.Radius
			}
		} else {
			penetration := circle.Radius - dist
			nx := dx / dist
			ny := dy / dist
			circle.PosX += nx * penetration
			circle.PosY += ny * penetration
		}
	}
}

func GetNearbySprites(posX, posY, radius float64, sprites []*Sprite) []*Sprite {
	nearby := make([]*Sprite, 0)
	radiusSq := radius * radius
	
	for _, sprite := range sprites {
		dx := sprite.PosX - posX
		dy := sprite.PosY - posY
		distSq := dx*dx + dy*dy
		
		if distSq <= radiusSq {
			nearby = append(nearby, sprite)
		}
	}
	
	return nearby
}

func GetSpritesInCone(originX, originY, dirX, dirY, range_, angle float64, sprites []*Sprite) []*Sprite {
	inCone := make([]*Sprite, 0)
	cosAngle := math.Cos(angle / 2)
	
	for _, sprite := range sprites {
		dx := sprite.PosX - originX
		dy := sprite.PosY - originY
		dist := math.Sqrt(dx*dx + dy*dy)
		
		if dist > range_ || dist < 0.001 {
			continue
		}
		
		nx := dx / dist
		ny := dy / dist
		
		dot := nx*dirX + ny*dirY
		
		if dot >= cosAngle {
			inCone = append(inCone, sprite)
		}
	}
	
	return inCone
}

func GetClosestSprite(posX, posY float64, sprites []*Sprite, spriteType SpriteType) *Sprite {
	var closest *Sprite
	minDistSq := math.MaxFloat64
	
	for _, sprite := range sprites {
		if sprite.Type != spriteType {
			continue
		}
		
		dx := sprite.PosX - posX
		dy := sprite.PosY - posY
		distSq := dx*dx + dy*dy
		
		if distSq < minDistSq {
			minDistSq = distSq
			closest = sprite
		}
	}
	
	return closest
}

func CheckLineIntersection(x1, y1, x2, y2, x3, y3, x4, y4 float64) (bool, float64, float64) {
	denom := (x1-x2)*(y3-y4) - (y1-y2)*(x3-x4)
	
	if math.Abs(denom) < 0.0001 {
		return false, 0, 0
	}
	
	t := ((x1-x3)*(y3-y4) - (y1-y3)*(x3-x4)) / denom
	u := -((x1-x2)*(y1-y3) - (y1-y2)*(x1-x3)) / denom
	
	if t >= 0 && t <= 1 && u >= 0 && u <= 1 {
		ix := x1 + t*(x2-x1)
		iy := y1 + t*(y2-y1)
		return true, ix, iy
	}
	
	return false, 0, 0
}

func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

