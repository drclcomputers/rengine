// texture.go - Texture loading and management
package engine

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"

	"github.com/go-gl/gl/v2.1/gl"
)

type Texture struct {
	ID     uint32
	Data   *image.RGBA
	Width  int
	Height int
}

const (
	TextureTypePNG = iota
	TextureTypeJPEG
)

func (e *Engine) loadTexture(path string, imageType int, texType string) (*Texture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open texture %s: %v", path, err)
	}
	defer file.Close()

	var img image.Image
	switch imageType {
	case TextureTypePNG:
		img, err = png.Decode(file)
		if err != nil {
			return nil, fmt.Errorf("failed to decode PNG %s: %v", path, err)
		}
	case TextureTypeJPEG:
		img, err = jpeg.Decode(file)
		if err != nil {
			return nil, fmt.Errorf("failed to decode JPEG %s: %v", path, err)
		}
	default:
		return nil, fmt.Errorf("unsupported image type: %d", imageType)
	}

	rgbaImg := image.NewRGBA(img.Bounds())
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			rgbaImg.Set(x, y, img.At(x, y))
		}
	}

	var texID uint32
	gl.GenTextures(1, &texID)
	gl.BindTexture(gl.TEXTURE_2D, texID)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgbaImg.Bounds().Dx()),
		int32(rgbaImg.Bounds().Dy()),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgbaImg.Pix),
	)

	texture := &Texture{
		ID:     texID,
		Data:   rgbaImg,
		Width:  rgbaImg.Bounds().Dx(),
		Height: rgbaImg.Bounds().Dy(),
	}

	fmt.Printf("Loaded %s texture: %s (%dx%d)\n", texType, path, texture.Width, texture.Height)
	return texture, nil
}

func (e *Engine) LoadTexture(path string, imageType int) (*Texture, error) {
	texture, err := e.loadTexture(path, imageType, "wall")
	if err != nil {
		return nil, err
	}
	e.Textures = append(e.Textures, texture)
	return texture, nil
}

func (e *Engine) LoadFloorTexture(path string, imageType int) (*Texture, error) {
	texture, err := e.loadTexture(path, imageType, "floor")
	if err != nil {
		return nil, err
	}
	e.FloorTexture = texture
	return texture, nil
}

func (e *Engine) LoadCeilingTexture(path string, imageType int) (*Texture, error) {
	texture, err := e.loadTexture(path, imageType, "ceiling")
	if err != nil {
		return nil, err
	}
	e.CeilingTexture = texture
	return texture, nil
}

func (t *Texture) GetPixel(x, y int) (r, g, b byte) {
	if x < 0 || x >= t.Width || y < 0 || y >= t.Height {
		return 0, 0, 0
	}

	idx := (y*t.Width + x) * 4
	r = t.Data.Pix[idx]
	g = t.Data.Pix[idx+1]
	b = t.Data.Pix[idx+2]
	return
}
