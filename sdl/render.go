package sdl

/*
#include "gosdl.h"
*/
import "C"
import (
	"image"
	"image/draw"
	"unsafe"

	"github.com/qeedquan/go-media/math/ga"
)

type (
	Renderer        C.SDL_Renderer
	RendererFlags   C.SDL_RendererFlags
	RendererFlip    C.SDL_RendererFlip
	Texture         C.SDL_Texture
	TextureAccess   C.SDL_TextureAccess
	TextureModulate C.SDL_TextureModulate
	ScaleMode       C.SDL_ScaleMode
)

type RendererInfo struct {
	Name             string
	Flags            uint32
	TextureFormats   []uint32
	MaxTextureWidth  int
	MaxTextureHeight int
}

type Vertex struct {
	Position ga.Vec2f
	Color    Color
	Texcoord ga.Vec2f
}

const (
	RENDERER_SOFTWARE      RendererFlags = C.SDL_RENDERER_SOFTWARE
	RENDERER_ACCELERATED   RendererFlags = C.SDL_RENDERER_ACCELERATED
	RENDERER_PRESENTVSYNC  RendererFlags = C.SDL_RENDERER_PRESENTVSYNC
	RENDERER_TARGETTEXTURE RendererFlags = C.SDL_RENDERER_TARGETTEXTURE
)

const (
	TEXTUREACCESS_STATIC    TextureAccess = C.SDL_TEXTUREACCESS_STATIC
	TEXTUREACCESS_STREAMING TextureAccess = C.SDL_TEXTUREACCESS_STREAMING
	TEXTUREACCESS_TARGET    TextureAccess = C.SDL_TEXTUREACCESS_TARGET
)

const (
	TEXTUREMODULATE_NONE  TextureModulate = C.SDL_TEXTUREMODULATE_NONE
	TEXTUREMODULATE_COLOR TextureModulate = C.SDL_TEXTUREMODULATE_COLOR
	TEXTUREMODULATE_ALPHA TextureModulate = C.SDL_TEXTUREMODULATE_ALPHA
)

const (
	FLIP_NONE       RendererFlip = C.SDL_FLIP_NONE
	FLIP_HORIZONTAL RendererFlip = C.SDL_FLIP_HORIZONTAL
	FLIP_VERTICAL   RendererFlip = C.SDL_FLIP_VERTICAL
)

const (
	ScaleModeNearest ScaleMode = C.SDL_ScaleModeNearest
	ScaleModeLinear  ScaleMode = C.SDL_ScaleModeLinear
	ScaleModeBest    ScaleMode = C.SDL_ScaleModeBest
)

func makeRendererInfo(info *C.SDL_RendererInfo) RendererInfo {
	r := RendererInfo{
		C.GoString(info.name),
		uint32(info.flags),
		nil,
		int(info.max_texture_width),
		int(info.max_texture_height),
	}
	r.TextureFormats = make([]uint32, info.num_texture_formats)
	for i := range r.TextureFormats {
		r.TextureFormats[i] = uint32(info.texture_formats[i])
	}
	return r
}

func GetNumRenderDrivers() int {
	return int(C.SDL_GetNumRenderDrivers())
}

func GetRenderDriverInfo(index int) (RendererInfo, error) {
	var info C.SDL_RendererInfo
	rc := C.SDL_GetRenderDriverInfo(C.int(index), &info)
	if rc < 0 {
		return RendererInfo{}, GetError()
	}
	return makeRendererInfo(&info), nil
}

func CreateWindowAndRenderer(width, height int, windowFlags WindowFlags) (*Window, *Renderer, error) {
	var window *C.SDL_Window
	var renderer *C.SDL_Renderer
	err := ek(C.SDL_CreateWindowAndRenderer(C.int(width), C.int(height), C.Uint32(windowFlags), &window, &renderer))
	return (*Window)(window), (*Renderer)(renderer), err
}

func CreateRenderer(window *Window, index int, rendererFlags RendererFlags) (*Renderer, error) {
	renderer := C.SDL_CreateRenderer((*C.SDL_Window)(window), C.int(index), C.Uint32(rendererFlags))
	if renderer == nil {
		return nil, GetError()
	}
	return (*Renderer)(renderer), nil
}

func CreateSoftwareRenderer(surface *Surface) (*Renderer, error) {
	re := C.SDL_CreateSoftwareRenderer((*C.SDL_Surface)(surface))
	if re == nil {
		return nil, GetError()
	}
	return (*Renderer)(re), nil
}

func (w *Window) Renderer() *Renderer {
	return (*Renderer)(C.SDL_GetRenderer((*C.SDL_Window)(w)))
}

func (re *Renderer) Info() (RendererInfo, error) {
	var info C.SDL_RendererInfo
	rc := C.SDL_GetRendererInfo((*C.SDL_Renderer)(re), &info)
	return makeRendererInfo(&info), ek(rc)
}

func (re *Renderer) OutputSize() (width, height int, err error) {
	var cw, ch C.int
	rc := C.SDL_GetRendererOutputSize((*C.SDL_Renderer)(re), &cw, &ch)
	return int(cw), int(ch), ek(rc)
}

func (re *Renderer) CreateTexture(format uint32, access TextureAccess, width, height int) (*Texture, error) {
	t := C.SDL_CreateTexture((*C.SDL_Renderer)(re), C.Uint32(format), C.int(access), C.int(width), C.int(height))
	if t == nil {
		return nil, GetError()
	}
	return (*Texture)(t), nil
}

func (re *Renderer) CreateTextureFromSurface(surface *Surface) (*Texture, error) {
	t := C.SDL_CreateTextureFromSurface((*C.SDL_Renderer)(re), (*C.SDL_Surface)(surface))
	if t == nil {
		return nil, GetError()
	}
	return (*Texture)(t), nil
}

func (re *Renderer) CreateTextureFromImage(img image.Image) (*Texture, error) {
	r := img.Bounds()
	s, err := CreateRGBSurfaceWithFormat(SWSURFACE, r.Dx(), r.Dy(), 32, PIXELFORMAT_ABGR8888)
	if err != nil {
		return nil, err
	}
	defer s.Free()
	draw.Draw(s, r, img, image.ZP, draw.Src)
	return re.CreateTextureFromSurface(s)
}

func (t *Texture) Query() (format uint32, access TextureAccess, width, height int, err error) {
	var cformat C.Uint32
	var caccess, cw, ch C.int
	rc := C.SDL_QueryTexture((*C.SDL_Texture)(t), &cformat, &caccess, &cw, &ch)
	return uint32(cformat), TextureAccess(caccess), int(cw), int(ch), ek(rc)
}

func (t *Texture) SetColorMod(c Color) error {
	return ek(C.SDL_SetTextureColorMod((*C.SDL_Texture)(t), C.Uint8(c.R), C.Uint8(c.G), C.Uint8(c.B)))
}

func (t *Texture) ColorMod() (Color, error) {
	var cr, cg, cb C.Uint8
	rc := C.SDL_GetTextureColorMod((*C.SDL_Texture)(t), &cr, &cg, &cb)
	return Color{uint8(cr), uint8(cg), uint8(cb), 255}, ek(rc)
}

func (t *Texture) SetAlphaMod(alpha uint8) error {
	return ek(C.SDL_SetTextureAlphaMod((*C.SDL_Texture)(t), C.Uint8(alpha)))
}

func (t *Texture) AlphaMod() (uint8, error) {
	var calpha C.Uint8
	rc := C.SDL_GetTextureAlphaMod((*C.SDL_Texture)(t), &calpha)
	return uint8(calpha), ek(rc)
}

func (t *Texture) Lock(rect *Rect) ([]byte, error) {
	_, _, _, height, err := t.Query()
	if err != nil {
		return nil, err
	}

	var pixels unsafe.Pointer
	var pitch C.int
	rc := C.SDL_LockTexture((*C.SDL_Texture)(t), (*C.SDL_Rect)(unsafe.Pointer(rect)), &pixels, &pitch)
	if rc < 0 {
		return nil, GetError()
	}

	var length C.int
	if rect == nil {
		length = pitch * C.int(height)
	} else {
		length = pitch * C.int(rect.H)
	}
	return ((*[1 << 30]uint8)(unsafe.Pointer(pixels)))[:length:length], nil
}

func (t *Texture) Unlock() {
	C.SDL_UnlockTexture((*C.SDL_Texture)(t))
}

func (t *Texture) SetScaleMode(mode ScaleMode) {
	C.SDL_SetTextureScaleMode((*C.SDL_Texture)(t), C.SDL_ScaleMode(mode))
}

func (t *Texture) ScaleMode() (mode ScaleMode, err error) {
	var cmode C.SDL_ScaleMode
	rc := C.SDL_GetTextureScaleMode((*C.SDL_Texture)(t), &cmode)
	return ScaleMode(cmode), ek(rc)
}

func (re *Renderer) SetDrawColor(c Color) error {
	return ek(C.SDL_SetRenderDrawColor((*C.SDL_Renderer)(re), C.Uint8(c.R), C.Uint8(c.G), C.Uint8(c.B), C.Uint8(c.A)))
}

func (re *Renderer) DrawColor() (Color, error) {
	var cr, cg, cb, ca C.Uint8
	rc := C.SDL_GetRenderDrawColor((*C.SDL_Renderer)(re), &cr, &cg, &cb, &ca)
	return Color{uint8(cr), uint8(cg), uint8(cb), uint8(ca)}, ek(rc)
}

func (re *Renderer) DrawPoint(x, y int) error {
	return ek(C.SDL_RenderDrawPoint((*C.SDL_Renderer)(re), C.int(x), C.int(y)))
}

func (re *Renderer) DrawPointF(x, y float32) error {
	return ek(C.SDL_RenderDrawPointF((*C.SDL_Renderer)(re), C.float(x), C.float(y)))
}

func (re *Renderer) DrawPoints(pts []Point) error {
	return ek(C.SDL_RenderDrawPoints((*C.SDL_Renderer)(re), (*C.SDL_Point)(unsafe.Pointer(&pts[0])), C.int(len(pts))))
}

func (re *Renderer) DrawPointsF(pts []FPoint) error {
	return ek(C.SDL_RenderDrawPointsF((*C.SDL_Renderer)(re), (*C.SDL_FPoint)(unsafe.Pointer(&pts[0])), C.int(len(pts))))
}

func (re *Renderer) DrawLine(x1, y1, x2, y2 int) error {
	return ek(C.SDL_RenderDrawLine((*C.SDL_Renderer)(re), C.int(x1), C.int(y1), C.int(x2), C.int(y2)))
}

func (re *Renderer) DrawLineF(x1, y1, x2, y2 float32) error {
	return ek(C.SDL_RenderDrawLineF((*C.SDL_Renderer)(re), C.float(x1), C.float(y1), C.float(x2), C.float(y2)))
}

func (re *Renderer) DrawLines(points []Point) error {
	return ek(C.SDL_RenderDrawLines((*C.SDL_Renderer)(re), (*C.SDL_Point)(unsafe.Pointer(&points[0])), C.int(len(points))))
}

func (re *Renderer) DrawLinesF(points []FPoint) error {
	return ek(C.SDL_RenderDrawLinesF((*C.SDL_Renderer)(re), (*C.SDL_FPoint)(unsafe.Pointer(&points[0])), C.int(len(points))))
}

func (re *Renderer) DrawRect(rect *Rect) error {
	return ek(C.SDL_RenderDrawRect((*C.SDL_Renderer)(re), (*C.SDL_Rect)(unsafe.Pointer(rect))))
}

func (re *Renderer) DrawRectF(rect *FRect) error {
	return ek(C.SDL_RenderDrawRectF((*C.SDL_Renderer)(re), (*C.SDL_FRect)(unsafe.Pointer(rect))))
}

func (re *Renderer) DrawRects(rects []Rect) error {
	return ek(C.SDL_RenderDrawRects((*C.SDL_Renderer)(re), (*C.SDL_Rect)(unsafe.Pointer(&rects[0])), C.int(len(rects))))
}

func (re *Renderer) DrawRectsF(rects []FRect) error {
	return ek(C.SDL_RenderDrawRectsF((*C.SDL_Renderer)(re), (*C.SDL_FRect)(unsafe.Pointer(&rects[0])), C.int(len(rects))))
}

func (re *Renderer) FillRect(rect *Rect) error {
	return ek(C.SDL_RenderFillRect((*C.SDL_Renderer)(re), (*C.SDL_Rect)(unsafe.Pointer(rect))))
}

func (re *Renderer) FillRectF(rect *Rect) error {
	return ek(C.SDL_RenderFillRectF((*C.SDL_Renderer)(re), (*C.SDL_FRect)(unsafe.Pointer(rect))))
}

func (re *Renderer) FillRects(rects []Rect) error {
	return ek(C.SDL_RenderFillRects((*C.SDL_Renderer)(re), (*C.SDL_Rect)(unsafe.Pointer(&rects[0])), C.int(len(rects))))
}

func (re *Renderer) FillRectsF(rects []Rect) error {
	return ek(C.SDL_RenderFillRectsF((*C.SDL_Renderer)(re), (*C.SDL_FRect)(unsafe.Pointer(&rects[0])), C.int(len(rects))))
}

func (re *Renderer) CopyEx(texture *Texture, src, dst *Rect, angle float64, center *Point, flip RendererFlip) error {
	return ek(C.SDL_RenderCopyEx((*C.SDL_Renderer)(re), (*C.SDL_Texture)(texture), (*C.SDL_Rect)(unsafe.Pointer(src)), (*C.SDL_Rect)(unsafe.Pointer(dst)),
		C.double(angle), (*C.SDL_Point)(unsafe.Pointer(center)), C.SDL_RendererFlip(flip)))
}

func (re *Renderer) CopyExF(texture *Texture, src *Rect, dst *FRect, angle float64, center *FPoint, flip RendererFlip) error {
	return ek(C.SDL_RenderCopyExF((*C.SDL_Renderer)(re), (*C.SDL_Texture)(texture), (*C.SDL_Rect)(unsafe.Pointer(src)), (*C.SDL_FRect)(unsafe.Pointer(dst)),
		C.double(angle), (*C.SDL_FPoint)(unsafe.Pointer(center)), C.SDL_RendererFlip(flip)))
}

func (re *Renderer) Copy(texture *Texture, src, dst *Rect) error {
	return ek(C.SDL_RenderCopy((*C.SDL_Renderer)(re), (*C.SDL_Texture)(texture), (*C.SDL_Rect)(unsafe.Pointer(src)), (*C.SDL_Rect)(unsafe.Pointer(dst))))
}

func (re *Renderer) CopyF(texture *Texture, src *Rect, dst *FRect) error {
	return ek(C.SDL_RenderCopyF((*C.SDL_Renderer)(re), (*C.SDL_Texture)(texture), (*C.SDL_Rect)(unsafe.Pointer(src)), (*C.SDL_FRect)(unsafe.Pointer(dst))))
}

func (re *Renderer) ReadPixels(rect *Rect, format uint32, pixels []byte, pitch int) error {
	return ek(C.SDL_RenderReadPixels((*C.SDL_Renderer)(re), (*C.SDL_Rect)(unsafe.Pointer(rect)), C.Uint32(format), unsafe.Pointer(&pixels[0]), C.int(pitch)))
}

func (re *Renderer) TargetSupported() bool {
	return C.SDL_RenderTargetSupported((*C.SDL_Renderer)(re)) != 0
}

func (re *Renderer) SetTarget(texture *Texture) error {
	return ek(C.SDL_SetRenderTarget((*C.SDL_Renderer)(re), (*C.SDL_Texture)(texture)))
}

func (re *Renderer) Present() {
	C.SDL_RenderPresent((*C.SDL_Renderer)(re))
}

func (t *Texture) Destroy() {
	C.SDL_DestroyTexture((*C.SDL_Texture)(t))
}

func (re *Renderer) Destroy() {
	C.SDL_DestroyRenderer((*C.SDL_Renderer)(re))
}

func (re *Renderer) SetLogicalSize(width, height int) {
	C.SDL_RenderSetLogicalSize((*C.SDL_Renderer)(re), C.int(width), C.int(height))
}

func (re *Renderer) LogicalSize() (width, height int) {
	var cw, ch C.int
	C.SDL_RenderGetLogicalSize((*C.SDL_Renderer)(re), &cw, &ch)
	return int(cw), int(ch)
}

func (re *Renderer) Clear() error {
	return ek(C.SDL_RenderClear((*C.SDL_Renderer)(re)))
}

func (re *Renderer) SetBlendMode(blendMode BlendMode) error {
	return ek(C.SDL_SetRenderDrawBlendMode((*C.SDL_Renderer)(re), C.SDL_BlendMode(blendMode)))
}

func (re *Renderer) BlendMode() (BlendMode, error) {
	var mode C.SDL_BlendMode
	rc := C.SDL_GetRenderDrawBlendMode((*C.SDL_Renderer)(re), &mode)
	return BlendMode(mode), ek(rc)
}

func (re *Renderer) SetViewport(rect *Rect) error {
	return ek(C.SDL_RenderSetViewport((*C.SDL_Renderer)(re), (*C.SDL_Rect)(unsafe.Pointer(rect))))
}

func (re *Renderer) Viewport() Rect {
	var r C.SDL_Rect
	C.SDL_RenderGetViewport((*C.SDL_Renderer)(re), &r)
	return Rect{int32(r.x), int32(r.y), int32(r.w), int32(r.h)}
}

func (re *Renderer) IsClipEnabled() bool {
	return C.SDL_RenderIsClipEnabled((*C.SDL_Renderer)(re)) != 0
}

func (re *Renderer) SetScale(scaleX, scaleY float64) error {
	return ek(C.SDL_RenderSetScale((*C.SDL_Renderer)(re), C.float(scaleX), C.float(scaleY)))
}

func (re *Renderer) SetIntegerScale(enable bool) error {
	return ek(C.SDL_RenderSetIntegerScale((*C.SDL_Renderer)(re), truth(enable)))
}

func (re *Renderer) IntegerScale() bool {
	return C.SDL_RenderGetIntegerScale((*C.SDL_Renderer)(re)) != 0
}

func (re *Renderer) Flush() error {
	return ek(C.SDL_RenderFlush((*C.SDL_Renderer)(re)))
}

func (re *Renderer) Geometry(texture *Texture, vertices []Vertex, indices []int32) error {
	var pvertex *Vertex
	var pindices *int32
	if len(vertices) > 0 {
		pvertex = &vertices[0]
	}
	if len(indices) > 0 {
		pindices = &indices[0]
	}
	return ek(C.SDL_RenderGeometry((*C.SDL_Renderer)(re), (*C.SDL_Texture)(texture), (*C.SDL_Vertex)(unsafe.Pointer(pvertex)), C.int(len(vertices)), (*C.int)(unsafe.Pointer(pindices)), C.int(len(indices))))
}

func (re *Renderer) SetVSync(vsync bool) error {
	return ek(C.SDL_RenderSetVSync((*C.SDL_Renderer)(re), C.int(truth(vsync))))
}
