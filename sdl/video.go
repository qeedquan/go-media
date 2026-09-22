package sdl

/*
#include "gosdl.h"
*/
import "C"
import (
	"reflect"
	"unsafe"
)

type DisplayMode struct {
	Format     uint32
	W, H       int
	Rate       int
	driverdata unsafe.Pointer
}

type (
	Window               C.SDL_Window
	WindowFlags          C.SDL_WindowFlags
	WindowEventID        C.SDL_WindowEventID
	GLattr               C.SDL_GLattr
	GLprofile            C.SDL_GLprofile
	GLContext            C.SDL_GLContext
	GLcontextFlag        C.SDL_GLcontextFlag
	GLcontextReleaseFlag C.SDL_GLcontextReleaseFlag
)

const (
	WINDOW_FULLSCREEN         WindowFlags = C.SDL_WINDOW_FULLSCREEN
	WINDOW_OPENGL             WindowFlags = C.SDL_WINDOW_OPENGL
	WINDOW_SHOWN              WindowFlags = C.SDL_WINDOW_SHOWN
	WINDOW_HIDDEN             WindowFlags = C.SDL_WINDOW_HIDDEN
	WINDOW_BORDERLESS         WindowFlags = C.SDL_WINDOW_BORDERLESS
	WINDOW_RESIZABLE          WindowFlags = C.SDL_WINDOW_RESIZABLE
	WINDOW_MINIMIZED          WindowFlags = C.SDL_WINDOW_MINIMIZED
	WINDOW_MAXIMIZED          WindowFlags = C.SDL_WINDOW_MAXIMIZED
	WINDOW_INPUT_GRABBED      WindowFlags = C.SDL_WINDOW_INPUT_GRABBED
	WINDOW_INPUT_FOCUS        WindowFlags = C.SDL_WINDOW_INPUT_FOCUS
	WINDOW_MOUSE_FOCUS        WindowFlags = C.SDL_WINDOW_MOUSE_FOCUS
	WINDOW_FULLSCREEN_DESKTOP WindowFlags = C.SDL_WINDOW_FULLSCREEN_DESKTOP
	WINDOW_FOREIGN            WindowFlags = C.SDL_WINDOW_FOREIGN
	WINDOW_ALLOW_HIGHDPI      WindowFlags = C.SDL_WINDOW_ALLOW_HIGHDPI
	WINDOW_MOUSE_CAPTURE      WindowFlags = C.SDL_WINDOW_MOUSE_CAPTURE
	WINDOW_ALWAYS_ON_TOP      WindowFlags = C.SDL_WINDOW_ALWAYS_ON_TOP
	WINDOW_SKIP_TASKBAR       WindowFlags = C.SDL_WINDOW_SKIP_TASKBAR
	WINDOW_UTILITY            WindowFlags = C.SDL_WINDOW_UTILITY
	WINDOW_TOOLTIP            WindowFlags = C.SDL_WINDOW_TOOLTIP
	WINDOW_POPUP_MENU         WindowFlags = C.SDL_WINDOW_POPUP_MENU
)

const (
	WINDOWPOS_UNDEFINED = C.SDL_WINDOWPOS_UNDEFINED
	WINDOWPOS_CENTERED  = C.SDL_WINDOWPOS_CENTERED
)

const (
	WINDOWEVENT_NONE         WindowEventID = C.SDL_WINDOWEVENT_NONE
	WINDOWEVENT_SHOWN        WindowEventID = C.SDL_WINDOWEVENT_SHOWN
	WINDOWEVENT_HIDDEN       WindowEventID = C.SDL_WINDOWEVENT_HIDDEN
	WINDOWEVENT_EXPOSED      WindowEventID = C.SDL_WINDOWEVENT_EXPOSED
	WINDOWEVENT_MOVED        WindowEventID = C.SDL_WINDOWEVENT_MOVED
	WINDOWEVENT_RESIZED      WindowEventID = C.SDL_WINDOWEVENT_RESIZED
	WINDOWEVENT_SIZE_CHANGED WindowEventID = C.SDL_WINDOWEVENT_SIZE_CHANGED
	WINDOWEVENT_MINIMIZED    WindowEventID = C.SDL_WINDOWEVENT_MINIMIZED
	WINDOWEVENT_MAXIMIZED    WindowEventID = C.SDL_WINDOWEVENT_MAXIMIZED
	WINDOWEVENT_RESTORED     WindowEventID = C.SDL_WINDOWEVENT_RESTORED
	WINDOWEVENT_ENTER        WindowEventID = C.SDL_WINDOWEVENT_ENTER
	WINDOWEVENT_LEAVE        WindowEventID = C.SDL_WINDOWEVENT_LEAVE
	WINDOWEVENT_FOCUS_GAINED WindowEventID = C.SDL_WINDOWEVENT_FOCUS_GAINED
	WINDOWEVENT_FOCUS_LOST   WindowEventID = C.SDL_WINDOWEVENT_FOCUS_LOST
	WINDOWEVENT_CLOSE        WindowEventID = C.SDL_WINDOWEVENT_CLOSE
	WINDOWEVENT_TAKE_FOCUS   WindowEventID = C.SDL_WINDOWEVENT_TAKE_FOCUS
	WINDOWEVENT_HIT_TEST     WindowEventID = C.SDL_WINDOWEVENT_HIT_TEST
)

const (
	GL_RED_SIZE                   = C.SDL_GL_RED_SIZE
	GL_GREEN_SIZE                 = C.SDL_GL_GREEN_SIZE
	GL_BLUE_SIZE                  = C.SDL_GL_BLUE_SIZE
	GL_ALPHA_SIZE                 = C.SDL_GL_ALPHA_SIZE
	GL_BUFFER_SIZE                = C.SDL_GL_BUFFER_SIZE
	GL_DOUBLEBUFFER               = C.SDL_GL_DOUBLEBUFFER
	GL_DEPTH_SIZE                 = C.SDL_GL_DEPTH_SIZE
	GL_STENCIL_SIZE               = C.SDL_GL_STENCIL_SIZE
	GL_ACCUM_RED_SIZE             = C.SDL_GL_ACCUM_RED_SIZE
	GL_ACCUM_GREEN_SIZE           = C.SDL_GL_ACCUM_GREEN_SIZE
	GL_ACCUM_BLUE_SIZE            = C.SDL_GL_ACCUM_BLUE_SIZE
	GL_ACCUM_ALPHA_SIZE           = C.SDL_GL_ACCUM_ALPHA_SIZE
	GL_STEREO                     = C.SDL_GL_STEREO
	GL_MULTISAMPLEBUFFERS         = C.SDL_GL_MULTISAMPLEBUFFERS
	GL_MULTISAMPLESAMPLES         = C.SDL_GL_MULTISAMPLESAMPLES
	GL_ACCELERATED_VISUAL         = C.SDL_GL_ACCELERATED_VISUAL
	GL_RETAINED_BACKING           = C.SDL_GL_RETAINED_BACKING
	GL_CONTEXT_MAJOR_VERSION      = C.SDL_GL_CONTEXT_MAJOR_VERSION
	GL_CONTEXT_MINOR_VERSION      = C.SDL_GL_CONTEXT_MINOR_VERSION
	GL_CONTEXT_EGL                = C.SDL_GL_CONTEXT_EGL
	GL_CONTEXT_FLAGS              = C.SDL_GL_CONTEXT_FLAGS
	GL_CONTEXT_PROFILE_MASK       = C.SDL_GL_CONTEXT_PROFILE_MASK
	GL_SHARE_WITH_CURRENT_CONTEXT = C.SDL_GL_SHARE_WITH_CURRENT_CONTEXT
	GL_FRAMEBUFFER_SRGB_CAPABLE   = C.SDL_GL_FRAMEBUFFER_SRGB_CAPABLE
	GL_CONTEXT_RELEASE_BEHAVIOR   = C.SDL_GL_CONTEXT_RELEASE_BEHAVIOR
)

const (
	GL_CONTEXT_PROFILE_CORE          = C.SDL_GL_CONTEXT_PROFILE_CORE
	GL_CONTEXT_PROFILE_COMPATIBILITY = C.SDL_GL_CONTEXT_PROFILE_COMPATIBILITY
	GL_CONTEXT_PROFILE_ES            = C.SDL_GL_CONTEXT_PROFILE_ES
)

const (
	GL_CONTEXT_DEBUG_FLAG              = C.SDL_GL_CONTEXT_DEBUG_FLAG
	GL_CONTEXT_FORWARD_COMPATIBLE_FLAG = C.SDL_GL_CONTEXT_FORWARD_COMPATIBLE_FLAG
	GL_CONTEXT_ROBUST_ACCESS_FLAG      = C.SDL_GL_CONTEXT_ROBUST_ACCESS_FLAG
	GL_CONTEXT_RESET_ISOLATION_FLAG    = C.SDL_GL_CONTEXT_RESET_ISOLATION_FLAG
)

const (
	GL_CONTEXT_RELEASE_BEHAVIOR_NONE  = C.SDL_GL_CONTEXT_RELEASE_BEHAVIOR_NONE
	GL_CONTEXT_RELEASE_BEHAVIOR_FLUSH = C.SDL_GL_CONTEXT_RELEASE_BEHAVIOR_FLUSH
)

func GetNumVideoDrivers() int {
	return int(C.SDL_GetNumVideoDrivers())
}

func GetVideoDriver(index int) string {
	return C.GoString(C.SDL_GetVideoDriver(C.int(index)))
}

func VideoInit(driverName string) error {
	cdriverName := append([]byte(driverName), 0)
	return ek(C.SDL_VideoInit((*C.char)(unsafe.Pointer(&cdriverName[0]))))
}

func VideoQuit() {
	C.SDL_VideoQuit()
}

func GetCurrentVideoDriver() string {
	return C.GoString(C.SDL_GetCurrentVideoDriver())
}

func GetNumVideoDisplays() int {
	return int(C.SDL_GetNumVideoDisplays())
}

func GetDisplayName(displayIndex int) string {
	return C.GoString(C.SDL_GetDisplayName(C.int(displayIndex)))
}

func GetDisplayBounds(displayIndex int) (Rect, error) {
	var cr C.SDL_Rect
	rc := C.SDL_GetDisplayBounds(C.int(displayIndex), &cr)
	if rc < 0 {
		return Rect{}, GetError()
	}
	return Rect{int32(cr.x), int32(cr.y), int32(cr.w), int32(cr.h)}, nil
}

func CreateWindow(title string, x, y, w, h int, flags WindowFlags) (*Window, error) {
	ctitle := append([]byte(title), 0)
	window := C.SDL_CreateWindow((*C.char)(unsafe.Pointer(&ctitle[0])), C.int(x), C.int(y), C.int(w), C.int(h), C.Uint32(flags))
	if window == nil {
		return nil, GetError()
	}
	return (*Window)(window), nil
}

func (w *Window) Flags() WindowFlags {
	return WindowFlags(C.SDL_GetWindowFlags((*C.SDL_Window)(w)))
}

func (w *Window) SetTitle(title string) {
	ctitle := append([]byte(title), 0)
	C.SDL_SetWindowTitle((*C.SDL_Window)(w), (*C.char)(unsafe.Pointer(&ctitle[0])))
}

func (w *Window) Title() string {
	return C.GoString(C.SDL_GetWindowTitle((*C.SDL_Window)(w)))
}

func GetGrabbedWindow() *Window {
	return (*Window)(C.SDL_GetGrabbedWindow())
}

func (w *Window) SetGrab(grabbed bool) {
	C.SDL_SetWindowGrab((*C.SDL_Window)(w), truth(grabbed))
}

func (w *Window) SetSize(width, height int) {
	C.SDL_SetWindowSize((*C.SDL_Window)(w), C.int(width), C.int(height))
}

func (w *Window) Size() (width, height int) {
	var cw, ch C.int
	C.SDL_GetWindowSize((*C.SDL_Window)(w), &cw, &ch)
	return int(cw), int(ch)
}

func (w *Window) SetMinimumSize(minWidth, minHeight int) {
	C.SDL_SetWindowMinimumSize((*C.SDL_Window)(w), C.int(minWidth), C.int(minHeight))
}

func (w *Window) MinimumSize() (minWidth, minHeight int) {
	var mw, mh C.int
	C.SDL_GetWindowMinimumSize((*C.SDL_Window)(w), &mw, &mh)
	return int(mw), int(mh)
}

func (w *Window) SetMaximumSize(maxWidth, maxHeight int) {
	C.SDL_SetWindowMaximumSize((*C.SDL_Window)(w), C.int(maxWidth), C.int(maxHeight))
}

func (w *Window) MaximumSize() (maxWidth, maxHeight int) {
	var mw, mh C.int
	C.SDL_GetWindowMaximumSize((*C.SDL_Window)(w), &mw, &mh)
	return int(mw), int(mh)
}

func (w *Window) SetBordered(bordered bool) {
	C.SDL_SetWindowBordered((*C.SDL_Window)(w), truth(bordered))
}

func (w *Window) Show() {
	C.SDL_ShowWindow((*C.SDL_Window)(w))
}

func (w *Window) Hide() {
	C.SDL_HideWindow((*C.SDL_Window)(w))
}

func (w *Window) Raise() {
	C.SDL_RaiseWindow((*C.SDL_Window)(w))
}

func (w *Window) Maximize() {
	C.SDL_MaximizeWindow((*C.SDL_Window)(w))
}

func (w *Window) SetBrightness(brightness float64) error {
	return ek(C.SDL_SetWindowBrightness((*C.SDL_Window)(w), C.float(brightness)))
}

func (w *Window) Brightness() float64 {
	return float64(C.SDL_GetWindowBrightness((*C.SDL_Window)(w)))
}

func (w *Window) SetOpacity(opacity float64) error {
	return ek(C.SDL_SetWindowOpacity((*C.SDL_Window)(w), C.float(opacity)))
}

func (w *Window) Opacity() (float64, error) {
	var copacity C.float
	rc := C.SDL_GetWindowOpacity((*C.SDL_Window)(w), &copacity)
	return float64(copacity), ek(rc)
}

func (w *Window) SetInputFocus() error {
	return ek(C.SDL_SetWindowInputFocus((*C.SDL_Window)(w)))
}

func IsScreenSaverEnabled() bool {
	return C.SDL_IsScreenSaverEnabled() != 0
}

func EnableScreenSaver() {
	C.SDL_EnableScreenSaver()
}

func DisableScreenSaver() {
	C.SDL_DisableScreenSaver()
}

func GLSetAttribute(attr GLattr, value int) error {
	return ek(C.SDL_GL_SetAttribute(C.SDL_GLattr(attr), C.int(value)))
}

func GLGetAttribute(attr GLattr) (int, error) {
	var cvalue C.int
	rc := C.SDL_GL_GetAttribute(C.SDL_GLattr(attr), &cvalue)
	return int(cvalue), ek(rc)
}

func (t *Texture) Update(rect *Rect, pixels interface{}, pitch int) error {
	return ek(C.SDL_UpdateTexture((*C.SDL_Texture)(t), (*C.SDL_Rect)(unsafe.Pointer(rect)),
		unsafe.Pointer(reflect.ValueOf(pixels).Pointer()), C.int(pitch)))
}

func (t *Texture) SetBlendMode(mode BlendMode) error {
	return ek(C.SDL_SetTextureBlendMode((*C.SDL_Texture)(t), C.SDL_BlendMode(mode)))
}

func (s *Surface) SetBlendMode(blendMode BlendMode) error {
	return ek(C.SDL_SetSurfaceBlendMode((*C.SDL_Surface)(s), C.SDL_BlendMode(blendMode)))
}

func (w *Window) SetIcon(icon *Surface) {
	C.SDL_SetWindowIcon((*C.SDL_Window)(w), (*C.SDL_Surface)(icon))
}

func (w *Window) SetFullscreen(flags WindowFlags) error {
	return ek(C.SDL_SetWindowFullscreen((*C.SDL_Window)(w), C.Uint32(flags)))
}

func (w *Window) BordersSize() (top, left, bottom, right int, err error) {
	var ctop, cleft, cbottom, cright C.int
	rc := C.SDL_GetWindowBordersSize((*C.SDL_Window)(w), &ctop, &cleft, &cbottom, &cright)
	return int(ctop), int(cleft), int(cbottom), int(cright), ek(rc)
}

func (w *Window) SetResizable(resizable bool) {
	C.SDL_SetWindowResizable((*C.SDL_Window)(w), truth(resizable))
}

func (w *Window) SetWindowOpacity(opacity float64) error {
	return ek(C.SDL_SetWindowOpacity((*C.SDL_Window)(w), C.float(opacity)))
}

func (w *Window) SetModal(parent *Window) error {
	return ek(C.SDL_SetWindowModalFor((*C.SDL_Window)(w), (*C.SDL_Window)(parent)))
}

func (w *Window) CreateContextGL() (GLContext, error) {
	ctx := GLContext(C.SDL_GL_CreateContext((*C.SDL_Window)(w)))
	if ctx == nil {
		return nil, GetError()
	}
	return ctx, nil
}

func (w *Window) SwapGL() {
	C.SDL_GL_SwapWindow((*C.SDL_Window)(w))
}

func (w *Window) DrawableSizeGL() (width, height int) {
	var cwidth, cheight C.int
	C.SDL_GL_GetDrawableSize((*C.SDL_Window)(w), &cwidth, &cheight)
	return int(cwidth), int(cheight)
}

func GLDeleteContext(c GLContext) {
	C.SDL_GL_DeleteContext(C.SDL_GLContext(c))
}

func GLSetSwapInterval(interval int) {
	C.SDL_GL_SetSwapInterval(C.int(interval))
}

func GetDisplayUsableBounds(displayIndex int, rect *Rect) error {
	return ek(C.SDL_GetDisplayUsableBounds(C.int(displayIndex), (*C.SDL_Rect)(unsafe.Pointer(rect))))
}

func GetDisplayMode(displayIndex, modeIndex int) (*DisplayMode, error) {
	var dm C.SDL_DisplayMode
	rc := C.SDL_GetDisplayMode(C.int(displayIndex), C.int(modeIndex), &dm)
	if rc != 0 {
		return nil, GetError()
	}

	return &DisplayMode{
		Format:     uint32(dm.format),
		W:          int(dm.w),
		H:          int(dm.h),
		Rate:       int(dm.refresh_rate),
		driverdata: dm.driverdata,
	}, nil
}

func GetDesktopDisplayMode(displayIndex int) (*DisplayMode, error) {
	var dm C.SDL_DisplayMode
	rc := C.SDL_GetDesktopDisplayMode(C.int(displayIndex), &dm)
	if rc != 0 {
		return nil, GetError()
	}

	return &DisplayMode{
		Format:     uint32(dm.format),
		W:          int(dm.w),
		H:          int(dm.h),
		Rate:       int(dm.refresh_rate),
		driverdata: dm.driverdata,
	}, nil
}

func GetCurrentDisplayMode(displayIndex int) (*DisplayMode, error) {
	var dm C.SDL_DisplayMode
	rc := C.SDL_GetCurrentDisplayMode(C.int(displayIndex), &dm)
	if rc != 0 {
		return nil, GetError()
	}

	return &DisplayMode{
		Format:     uint32(dm.format),
		W:          int(dm.w),
		H:          int(dm.h),
		Rate:       int(dm.refresh_rate),
		driverdata: dm.driverdata,
	}, nil
}
