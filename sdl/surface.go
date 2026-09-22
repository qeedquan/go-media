package sdl

/*
#include "gosdl.h"

Uint32 getPixel(SDL_Surface *surface, size_t x, size_t y) {
	SDL_LockSurface(surface);
	int bpp = surface->format->BytesPerPixel;
	Uint8 *p = (Uint8 *)surface->pixels + y * surface->pitch + x * bpp;
	Uint32 c;

	switch(bpp) {
	case 1:
		c = *p;
		break;

	case 2:
		c = *(Uint16 *)p;
		break;

	case 3:
		if(SDL_BYTEORDER == SDL_BIG_ENDIAN)
			c = p[0] << 16 | p[1] << 8 | p[2];
		else
			c = p[0] | p[1] << 8 | p[2] << 16;
		break;

	case 4:
		c = *(Uint32 *)p;
		break;

	default:
		c = 0;
	}
	SDL_UnlockSurface(surface);
	return c;
}

void setPixel(SDL_Surface *surface, size_t x, size_t y, Uint32 pixel) {
	SDL_LockSurface(surface);
	int bpp = surface->format->BytesPerPixel;
	Uint8 *p = (Uint8 *)surface->pixels + y * surface->pitch + x * bpp;

	switch(bpp) {
	case 1:
		*p = pixel;
		break;

	case 2:
		*(Uint16 *)p = pixel;
		break;

	case 3:
		if(SDL_BYTEORDER == SDL_BIG_ENDIAN) {
			p[0] = (pixel >> 16) & 0xff;
			p[1] = (pixel >> 8) & 0xff;
			p[2] = pixel & 0xff;
		} else {
			p[0] = pixel & 0xff;
			p[1] = (pixel >> 8) & 0xff;
			p[2] = (pixel >> 16) & 0xff;
		}
		break;

	case 4:
		*(Uint32 *)p = pixel;
		break;
	}
	SDL_UnlockSurface(surface);
}

int blitSurface(SDL_Surface *src, SDL_Rect *srcrect, SDL_Surface *dst, SDL_Rect *dstrect) {
	return SDL_BlitSurface(src, srcrect, dst, dstrect);
}
*/
import "C"
import (
	"image"
	"image/color"
	"unsafe"
)

const (
	SWSURFACE = C.SDL_SWSURFACE
	PREALLOC  = C.SDL_PREALLOC
	RLEACCEL  = C.SDL_RLEACCEL
	DONTFREE  = C.SDL_DONTFREE
)

type Surface C.SDL_Surface

func (s *Surface) Flags() uint32             { return uint32(s.flags) }
func (s *Surface) Size() (width, height int) { return int(s.w), int(s.h) }
func (s *Surface) PixelFormat() *PixelFormat { return &PixelFormat{s.format} }
func (s *Surface) Pixels() []byte {
	n := s.w * s.h * C.int(s.format.BytesPerPixel)
	return ((*[1 << 30]uint8)(s.pixels))[:n:n]
}

func CreateRGBSurface(flags uint32, width, height, depth int, rmask, gmask, bmask, amask uint32) (*Surface, error) {
	s := (*Surface)(C.SDL_CreateRGBSurface(C.Uint32(flags), C.int(width), C.int(height), C.int(depth),
		C.Uint32(rmask), C.Uint32(gmask), C.Uint32(bmask), C.Uint32(amask)))
	if s == nil {
		return nil, GetError()
	}
	return s, nil
}

func CreateRGBSurfaceWithFormat(flags uint32, width, height, depth int, format uint32) (*Surface, error) {
	s := (*Surface)(C.SDL_CreateRGBSurfaceWithFormat(C.Uint32(flags), C.int(width), C.int(height), C.int(depth), C.Uint32(format)))
	if s == nil {
		return nil, GetError()
	}
	return s, nil
}

func CreateRGBSurfaceFrom(pixels []byte, width, height, depth, pitch int, rmask, gmask, bmask, amask uint32) (*Surface, error) {
	s := (*Surface)(C.SDL_CreateRGBSurfaceFrom(unsafe.Pointer(&pixels[0]), C.int(width), C.int(height), C.int(depth), C.int(pitch),
		C.Uint32(rmask), C.Uint32(gmask), C.Uint32(bmask), C.Uint32(amask)))
	if s == nil {
		return nil, GetError()
	}
	return s, nil
}

func CreateRGBSurfaceWithFormatFrom(pixels []byte, width, height, depth, pitch int, format uint32) (*Surface, error) {
	s := (*Surface)(C.SDL_CreateRGBSurfaceWithFormatFrom(unsafe.Pointer(&pixels[0]), C.int(width), C.int(height), C.int(depth), C.int(pitch), C.Uint32(format)))
	if s == nil {
		return nil, GetError()
	}
	return s, nil
}

func (s *Surface) SetClipRect(rect Rect) {
	C.SDL_SetClipRect((*C.SDL_Surface)(s), (*C.SDL_Rect)(unsafe.Pointer(&rect)))
}

func (s *Surface) ClipRect() Rect {
	var rect Rect
	C.SDL_GetClipRect((*C.SDL_Surface)(s), (*C.SDL_Rect)(unsafe.Pointer(&rect)))
	return rect
}

func (s *Surface) FillRect(rect *Rect, c Color) {
	col := uint32(c.R) | uint32(c.G)<<8 | uint32(c.B)<<16 | uint32(c.A)<<24
	C.SDL_FillRect((*C.SDL_Surface)(s), (*C.SDL_Rect)(unsafe.Pointer(rect)), C.Uint32(col))
}

func (s *Surface) FillRects(rects []Rect, c Color) {
	col := uint32(c.R) | uint32(c.G)<<8 | uint32(c.B)<<16 | uint32(c.A)<<24
	C.SDL_FillRects((*C.SDL_Surface)(s), (*C.SDL_Rect)(unsafe.Pointer(&rects[0])), C.int(len(rects)), C.Uint32(col))
}

func BlitSurface(src *Surface, srcrect *Rect, dst *Surface, dstrect *Rect) error {
	return ek(C.blitSurface((*C.SDL_Surface)(src), (*C.SDL_Rect)(unsafe.Pointer(srcrect)), (*C.SDL_Surface)(dst), (*C.SDL_Rect)(unsafe.Pointer(dstrect))))
}

func (s *Surface) Lock() error {
	return ek(C.SDL_LockSurface((*C.SDL_Surface)(s)))
}

func (s *Surface) Unlock() {
	C.SDL_UnlockSurface((*C.SDL_Surface)(s))
}

func (s *Surface) Free() {
	C.SDL_FreeSurface((*C.SDL_Surface)(s))
}

func (s *Surface) Pitch() int {
	return int(s.pitch)
}

func (s *Surface) ColorModel() color.Model {
	return color.NRGBAModel
}

func (s *Surface) Bounds() image.Rectangle {
	return image.Rect(0, 0, int(s.w), int(s.h))
}

func (s *Surface) At(x, y int) color.Color {
	var cr, cg, cb, ca C.Uint8
	pixel := C.getPixel((*C.SDL_Surface)(s), C.size_t(x), C.size_t(y))
	C.SDL_GetRGBA(pixel, s.format, &cr, &cg, &cb, &ca)
	return color.NRGBA{uint8(cr), uint8(cg), uint8(cb), uint8(ca)}
}

func (s *Surface) Set(x, y int, c color.Color) {
	p := color.NRGBAModel.Convert(c).(color.NRGBA)
	C.setPixel((*C.SDL_Surface)(s), C.size_t(x), C.size_t(y), C.SDL_MapRGBA(s.format, C.Uint8(p.R), C.Uint8(p.G), C.Uint8(p.B), C.Uint8(p.A)))
}
