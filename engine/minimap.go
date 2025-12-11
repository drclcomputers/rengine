package engine
/*
func (e *Engine) RenderMinimap() {
    if !e.MinimapEnabled {
        return
    }
    
    size := 150
    padding := 10
    x, y := e.getMinimapPosition(size, padding)
    
    e.drawMinimapBackground(x, y, size)
    
    scale := float64(size) / float64(max(e.MapWidth, e.MapHeight))
    e.drawMinimapTiles(x, y, scale)
    
    e.drawMinimapSprites(x, y, scale)
    
    e.drawMinimapPlayer(x, y, scale)
    
    e.drawMinimapFOV(x, y, scale)
}

func (e *Engine) getMinimapPosition(size, padding int) (int, int) {
    switch e.MinimapCorner {
    case 0: return padding, padding  // Top-left
    case 1: return e.ScreenWidth - size - padding, padding  // Top-right
    case 2: return e.ScreenWidth - size - padding, e.ScreenHeight - size - padding  // Bottom-right
    case 3: return padding, e.ScreenHeight - size - padding  // Bottom-left
    }
    return 0, 0
}

func (e *Engine) ToggleMinimap() {
    e.MinimapEnabled = !e.MinimapEnabled
}

func (e *Engine) SetMinimapCorner(corner int) {
    if corner >= 0 && corner <= 3 {
        e.MinimapCorner = corner
    }
}*/
