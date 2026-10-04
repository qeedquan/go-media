package sdlimage

/*
#include <SDL.h>
*/
import "C"

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"

	"github.com/qeedquan/go-media/image/imageutil"
	_ "github.com/qeedquan/go-media/image/psd"
	_ "github.com/qeedquan/go-media/image/tga"
	"github.com/qeedquan/go-media/sdl"
	_ "golang.org/x/image/bmp"
)

func LoadTextureFile(re *sdl.Renderer, name string) (*sdl.Texture, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return LoadTextureReader(re, f)
}

func LoadTextureReader(re *sdl.Renderer, r io.Reader) (*sdl.Texture, error) {
	m, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}
	return LoadTextureImage(re, m)
}

func LoadTextureImage(re *sdl.Renderer, m image.Image) (*sdl.Texture, error) {
	r := m.Bounds()
	w, h := r.Dx(), r.Dy()
	if w == 0 || h == 0 {
		return nil, fmt.Errorf("invalid dimension %dx%d", w, h)
	}

	b := C.malloc(C.size_t(w * h * 4))
	p := ((*[1 << 30]uint8)(b))[:w*h*4 : w*h*4]

	n := &image.NRGBA{p, w * 4, image.Rect(0, 0, w, h)}
	draw.Draw(n, n.Bounds(), m, r.Min, draw.Src)

	s, err := sdl.CreateRGBSurfaceFrom(n.Pix[:], w, h, 32, w*4, 0xff, 0xff00, 0xff0000, 0xff000000)
	if err != nil {
		C.free(b)
		return nil, err
	}
	defer s.Free()

	return re.CreateTextureFromSurface(s)
}

func LoadSurfaceFile(name string) (*sdl.Surface, error) {
	m, err := imageutil.LoadRGBAFile(name)
	if err != nil {
		return nil, err
	}
	return LoadSurfaceImage(m)
}

func LoadSurfaceImage(m image.Image) (*sdl.Surface, error) {
	r := m.Bounds()
	bpp, rmask, gmask, bmask, amask := sdl.PixelFormatEnumToMasks(sdl.PIXELFORMAT_ARGB8888)
	s, err := sdl.CreateRGBSurface(sdl.SWSURFACE, r.Dx(), r.Dy(), bpp, rmask, gmask, bmask, amask)
	if err != nil {
		return nil, err
	}
	draw.Draw(s, s.Bounds(), m, r.Min, draw.Src)
	return s, nil
}
