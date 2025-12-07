// texture.go - Texture loading and management
package engine

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
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

// LoadTexture loads a texture from a file and adds it to the engine
func (e *Engine) LoadTexture(path string, imageType int) (*Texture, error) {
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

	e.Textures = append(e.Textures, texture)
	fmt.Printf("Loaded texture: %s (%dx%d)\n", path, texture.Width, texture.Height)

	return texture, nil
}

// GetPixel retrieves a pixel color from the texture
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

// convertJPEGToPNG converts from JPEG to PNG.
func (t *Texture) ConvertJPEGToPNG(w io.Writer, r io.Reader) error {
	img, err := jpeg.Decode(r)
	if err != nil {
		return err
	}
	return png.Encode(w, img)
}
