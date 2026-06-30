package xlib

/*
#include "go_xlib.h"

#cgo pkg-config: x11
#cgo LDFLAGS: -lXss

*/
import "C"
import (
	"fmt"
	"math"
	"unsafe"

	"github.com/qeedquan/go-media/x11/xlib/xrm"
)

type (
	Atom                  C.Atom
	ButtonEvent           C.XButtonEvent
	ClassHint             C.XClassHint
	ClientMessageEvent    C.XClientMessageEvent
	Color                 C.XColor
	Colormap              C.Colormap
	ColormapEvent         C.XColormapEvent
	ConfigureEvent        C.XConfigureEvent
	ConfigureRequestEvent C.XConfigureRequestEvent
	CreateWindowEvent     C.XCreateWindowEvent
	CrossingEvent         C.XCrossingEvent
	Cursor                C.Cursor
	DestroyWindowEvent    C.XDestroyWindowEvent
	Display               C.Display
	Drawable              C.Drawable
	ErrorEvent            C.XErrorEvent
	Event                 C.XEvent
	ExposeEvent           C.XExposeEvent
	FocusChangeEvent      C.XFocusChangeEvent
	Font                  C.Font
	GC                    C.GC
	GCValues              C.XGCValues
	GenericEventCookie    C.XGenericEventCookie
	GravityEvent          C.XGravityEvent
	ICCEncodingStyle      C.XICCEncodingStyle
	IC                    C.XIC
	ID                    C.XID
	Image                 C.XImage
	IM                    C.XIM
	KeyEvent              C.XKeyEvent
	KeyPressedEvent       C.XKeyPressedEvent
	KeySym                C.KeySym
	MapEvent              C.XMapEvent
	MapRequestEvent       C.XMapRequestEvent
	MotionEvent           C.XMotionEvent
	Pixmap                C.Pixmap
	Point                 C.XPoint
	Pointer               C.XPointer
	PropertyEvent         C.XPropertyEvent
	Rectangle             C.XRectangle
	Region                C.Region
	ReparentEvent         C.XReparentEvent
	ResizeRequestEvent    C.XResizeRequestEvent
	Screen                C.Screen
	ScreenSaverInfo       C.XScreenSaverInfo
	SelectionEvent        C.XSelectionEvent
	SelectionClearEvent   C.XSelectionClearEvent
	SelectionRequestEvent C.XSelectionRequestEvent
	SetWindowAttributes   C.XSetWindowAttributes
	SizeHints             C.XSizeHints
	Status                C.Status
	TextProperty          C.XTextProperty
	Time                  C.Time
	TimeCoord             C.XTimeCoord
	UnmapEvent            C.XUnmapEvent
	VaNestedList          struct {
		list C.XVaNestedList
		key  []*C.char
		val  []unsafe.Pointer
	}
	VisibilityEvent  C.XVisibilityEvent
	Visual           C.Visual
	WindowAttributes C.XWindowAttributes
	Window           C.Window
	WMHints          C.XWMHints
)

type IDProc func(*Display, Pointer)

type IMValue struct {
	Key   string
	Value interface{}
}

type ICValue = IMValue

const (
	NVaNestedList                = C.XNVaNestedList
	NSeparatorofNestedList       = C.XNSeparatorofNestedList
	NQueryInputStyle             = C.XNQueryInputStyle
	NFocusWindow                 = C.XNFocusWindow
	NResourceName                = C.XNResourceName
	NResourceClass               = C.XNResourceClass
	NGeometryCallback            = C.XNGeometryCallback
	NDestroyCallback             = C.XNDestroyCallback
	NFilterEvents                = C.XNFilterEvents
	NPreeditStartCallback        = C.XNPreeditStartCallback
	NPreeditDoneCallback         = C.XNPreeditDoneCallback
	NPreeditDrawCallback         = C.XNPreeditDrawCallback
	NPreeditCaretCallback        = C.XNPreeditCaretCallback
	NPreeditStateNotifyCallback  = C.XNPreeditStateNotifyCallback
	NPreeditAttributes           = C.XNPreeditAttributes
	NStatusStartCallback         = C.XNStatusStartCallback
	NStatusDoneCallback          = C.XNStatusDoneCallback
	NStatusDrawCallback          = C.XNStatusDrawCallback
	NStatusAttributes            = C.XNStatusAttributes
	NArea                        = C.XNArea
	NAreaNeeded                  = C.XNAreaNeeded
	NSpotLocation                = C.XNSpotLocation
	NColormap                    = C.XNColormap
	NStdColormap                 = C.XNStdColormap
	NForeground                  = C.XNForeground
	NBackground                  = C.XNBackground
	NBackgroundPixmap            = C.XNBackgroundPixmap
	NFontSet                     = C.XNFontSet
	NLineSpace                   = C.XNLineSpace
	NCursor                      = C.XNCursor
	NQueryIMValuesList           = C.XNQueryIMValuesList
	NQueryICValuesList           = C.XNQueryICValuesList
	NStringConversionCallback    = C.XNStringConversionCallback
	NStringConversion            = C.XNStringConversion
	NResetState                  = C.XNResetState
	NHotKey                      = C.XNHotKey
	NHotKeyState                 = C.XNHotKeyState
	NPreeditState                = C.XNPreeditState
	NVisiblePosition             = C.XNVisiblePosition
	NRequiredCharSet             = C.XNRequiredCharSet
	NQueryOrientation            = C.XNQueryOrientation
	NDirectionalDependentDrawing = C.XNDirectionalDependentDrawing
	NContextualDrawing           = C.XNContextualDrawing
	NBaseFontName                = C.XNBaseFontName
	NMissingCharSet              = C.XNMissingCharSet
	NDefaultString               = C.XNDefaultString
	NOrientation                 = C.XNOrientation
	NFontInfo                    = C.XNFontInfo
	NOMAutomatic                 = C.XNOMAutomatic
	NInputStyle                  = C.XNInputStyle
	NClientWindow                = C.XNClientWindow
)

const (
	GCFunction          = C.GCFunction
	GCPlaneMask         = C.GCPlaneMask
	GCForeground        = C.GCForeground
	GCBackground        = C.GCBackground
	GCLineWidth         = C.GCLineWidth
	GCLineStyle         = C.GCLineStyle
	GCCapStyle          = C.GCCapStyle
	GCJoinStyle         = C.GCJoinStyle
	GCFillStyle         = C.GCFillStyle
	GCFillRule          = C.GCFillRule
	GCTile              = C.GCTile
	GCStipple           = C.GCStipple
	GCTileStipXOrigin   = C.GCTileStipXOrigin
	GCTileStipYOrigin   = C.GCTileStipYOrigin
	GCFont              = C.GCFont
	GCSubwindowMode     = C.GCSubwindowMode
	GCGraphicsExposures = C.GCGraphicsExposures
	GCClipXOrigin       = C.GCClipXOrigin
	GCClipYOrigin       = C.GCClipYOrigin
	GCClipMask          = C.GCClipMask
	GCDashOffset        = C.GCDashOffset
	GCDashList          = C.GCDashList
	GCArcMode           = C.GCArcMode
)

const (
	NoEventMask      = C.NoEventMask
	FocusOut         = C.FocusOut
	KeymapNotify     = C.KeymapNotify
	Expose           = C.Expose
	GraphicsExpose   = C.GraphicsExpose
	NoExpose         = C.NoExpose
	VisibilityNotify = C.VisibilityNotify
	CreateNotify     = C.CreateNotify
	DestroyNotify    = C.DestroyNotify
	UnmapNotify      = C.UnmapNotify
	MapNotify        = C.MapNotify
	KeyPress         = C.KeyPress
	MapRequest       = C.MapRequest
	ReparentNotify   = C.ReparentNotify
	ConfigureNotify  = C.ConfigureNotify
	ConfigureRequest = C.ConfigureRequest
	GravityNotify    = C.GravityNotify
	ResizeRequest    = C.ResizeRequest
	CirculateNotify  = C.CirculateNotify
	CirculateRequest = C.CirculateRequest
	PropertyNotify   = C.PropertyNotify
	SelectionClear   = C.SelectionClear
	KeyRelease       = C.KeyRelease
	SelectionRequest = C.SelectionRequest
	SelectionNotify  = C.SelectionNotify
	ColormapNotify   = C.ColormapNotify
	ClientMessage    = C.ClientMessage
	MappingNotify    = C.MappingNotify
	ButtonPress      = C.ButtonPress
	ButtonRelease    = C.ButtonRelease
	MotionNotify     = C.MotionNotify
	EnterNotify      = C.EnterNotify
	LeaveNotify      = C.LeaveNotify
	FocusIn          = C.FocusIn
	GenericEvent     = C.GenericEvent
	LASTEvent        = C.LASTEvent
)

const (
	Button1Mask = C.Button1Mask
	ShiftMask   = C.ShiftMask
	LockMask    = C.LockMask
	ControlMask = C.ControlMask
	Mod1Mask    = C.Mod1Mask
	Mod2Mask    = C.Mod2Mask
	Mod3Mask    = C.Mod3Mask
	Mod4Mask    = C.Mod4Mask
	Mod5Mask    = C.Mod5Mask
)

const (
	VisibilityUnobscured        = C.VisibilityUnobscured
	VisibilityPartiallyObscured = C.VisibilityPartiallyObscured
	VisibilityFullyObscured     = C.VisibilityFullyObscured
)

const (
	ForgetGravity    = C.ForgetGravity
	UnmapGravity     = C.UnmapGravity
	NorthWestGravity = C.NorthWestGravity
	StaticGravity    = C.StaticGravity
	NorthGravity     = C.NorthGravity
	NorthEastGravity = C.NorthEastGravity
	WestGravity      = C.WestGravity
	CenterGravity    = C.CenterGravity
	EastGravity      = C.EastGravity
	SouthWestGravity = C.SouthWestGravity
	SouthGravity     = C.SouthGravity
	SouthEastGravity = C.SouthEastGravity
)

const (
	CWBackPixmap       = C.CWBackPixmap
	CWBackPixel        = C.CWBackPixel
	CWSaveUnder        = C.CWSaveUnder
	CWEventMask        = C.CWEventMask
	CWDontPropagate    = C.CWDontPropagate
	CWColormap         = C.CWColormap
	CWCursor           = C.CWCursor
	CWBorderPixmap     = C.CWBorderPixmap
	CWBorderPixel      = C.CWBorderPixel
	CWBitGravity       = C.CWBitGravity
	CWWinGravity       = C.CWWinGravity
	CWBackingStore     = C.CWBackingStore
	CWBackingPlanes    = C.CWBackingPlanes
	CWBackingPixel     = C.CWBackingPixel
	CWOverrideRedirect = C.CWOverrideRedirect
)

const (
	PSize       = C.PSize
	PResizeInc  = C.PResizeInc
	PBaseSize   = C.PBaseSize
	PMinSize    = C.PMinSize
	PMaxSize    = C.PMaxSize
	PWinGravity = C.PWinGravity
)

const (
	USPosition = C.USPosition
)

const (
	InputOutput = C.InputOutput
	InputOnly   = C.InputOnly
)

const (
	CWX           = C.CWX
	CWY           = C.CWY
	CWWidth       = C.CWWidth
	CWHeight      = C.CWHeight
	CWBorderWidth = C.CWBorderWidth
	CWSibling     = C.CWSibling
	CWStackMode   = C.CWStackMode
)

const (
	Above    = C.Above
	Below    = C.Below
	TopIf    = C.TopIf
	BottomIf = C.BottomIf
	Opposite = C.Opposite
)

const (
	RaiseLowest  = C.RaiseLowest
	LowerHighest = C.LowerHighest
)

const (
	LineSolid      = C.LineSolid
	LineOnOffDash  = C.LineOnOffDash
	LineDoubleDash = C.LineDoubleDash
)

const (
	CapNotLast    = C.CapNotLast
	CapButt       = C.CapButt
	CapRound      = C.CapRound
	CapProjecting = C.CapProjecting
)

const (
	JoinMiter = C.JoinMiter
	JoinRound = C.JoinRound
	JoinBevel = C.JoinBevel
)

const (
	FillSolid          = C.FillSolid
	FillTiled          = C.FillTiled
	FillStippled       = C.FillStippled
	FillOpaqueStippled = C.FillOpaqueStippled
)

const (
	EvenOddRule = C.EvenOddRule
	WindingRule = C.WindingRule
)

const (
	NoValue     = C.NoValue
	XValue      = C.XValue
	YValue      = C.YValue
	WidthValue  = C.WidthValue
	HeightValue = C.HeightValue
	AllValues   = C.AllValues
	XNegative   = C.XNegative
	YNegative   = C.YNegative
)

const (
	DisableScreenInterval = C.DisableScreenInterval
	DisableScreenSaver    = C.DisableScreenSaver
	DontAllowExposures    = C.DontAllowExposures
	DontPreferBlanking    = C.DontPreferBlanking
	AllowExposures        = C.AllowExposures
	PreferBlanking        = C.PreferBlanking
	DefaultBlanking       = C.DefaultBlanking
	DefaultExposures      = C.DefaultExposures
)

const (
	InputHint        = C.InputHint
	StateHint        = C.StateHint
	IconPixmapHint   = C.IconPixmapHint
	IconWindowHint   = C.IconWindowHint
	IconPositionHint = C.IconPositionHint
	IconMaskHint     = C.IconMaskHint
	WindowGroupHint  = C.WindowGroupHint
	UrgencyHint      = C.XUrgencyHint
	AllHints         = C.AllHints
)

const (
	WithdrawnState = C.WithdrawnState
	NormalState    = C.NormalState
	IconicState    = C.IconicState
)

const (
	IsUnmapped   = C.IsUnmapped
	IsUnviewable = C.IsUnviewable
	IsViewable   = C.IsViewable
)

const (
	Success             Status = C.Success
	BadAtom             Status = C.BadAtom
	BadRequest          Status = C.BadRequest
	BadAccess           Status = C.BadAccess
	BadColor            Status = C.BadColor
	BadImplementation   Status = C.BadImplementation
	BadAlloc            Status = C.BadAlloc
	BadFont             Status = C.BadFont
	BadGC               Status = C.BadGC
	BadMatch            Status = C.BadMatch
	BadPixmap           Status = C.BadPixmap
	BadValue            Status = C.BadValue
	BadWindow           Status = C.BadWindow
	BadDrawable         Status = C.BadDrawable
	FirstExtensionError Status = C.FirstExtensionError
)

const (
	NoMemory           Status = C.XNoMemory
	LocaleNotSupported Status = C.XLocaleNotSupported
	ConverterNotFound  Status = C.XConverterNotFound
)

const (
	StringStyle       = C.XStringStyle
	CompoundTextStyle = C.XCompoundTextStyle
	TextStyle         = C.XTextStyle
	StdICCTextStyle   = C.XStdICCTextStyle
)

const (
	None = C.None
)

const (
	CurrentTime = C.CurrentTime
)

const (
	PointerRoot         = C.PointerRoot
	RevertToPointerRoot = C.RevertToPointerRoot
)

const (
	ScreenSaverOn       = C.ScreenSaverOn
	ScreenSaverOff      = C.ScreenSaverOff
	ScreenSaverDisabled = C.ScreenSaverDisabled
)

const (
	CopyFromParent = C.CopyFromParent
)

const (
	XA_PRIMARY   = C.XA_PRIMARY
	XA_SECONDARY = C.XA_SECONDARY
	XA_CARDINAL  = C.XA_CARDINAL
	XA_STRING    = C.XA_STRING
	XA_ATOM      = C.XA_ATOM
)

const (
	UTF8StringStyle = C.XUTF8StringStyle
)

const (
	PointerMotionMask      = C.PointerMotionMask
	KeyReleaseMask         = C.KeyReleaseMask
	KeyPressMask           = C.KeyPressMask
	ExposureMask           = C.ExposureMask
	VisibilityChangeMask   = C.VisibilityChangeMask
	Button1MotionMask      = C.Button1MotionMask
	Button2MotionMask      = C.Button2MotionMask
	Button3MotionMask      = C.Button3MotionMask
	Button4MotionMask      = C.Button4MotionMask
	Button5MotionMask      = C.Button5MotionMask
	ButtonMotionMask       = C.ButtonMotionMask
	ButtonPressMask        = C.ButtonPressMask
	ButtonReleaseMask      = C.ButtonReleaseMask
	FocusChangeMask        = C.FocusChangeMask
	StructureNotifyMask    = C.StructureNotifyMask
	SubstructureNotifyMask = C.SubstructureNotifyMask
)

const (
	PropModeReplace = C.PropModeReplace
	PropModePrepend = C.PropModePrepend
	PropModeAppend  = C.PropModeAppend
)

const (
	QueuedAlready      = C.QueuedAlready
	QueuedAfterFlush   = C.QueuedAfterFlush
	QueuedAfterReading = C.QueuedAfterReading
)

const (
	IMPreeditNothing = C.XIMPreeditNothing
	IMStatusNothing  = C.XIMStatusNothing
)

const (
	NotifyNormal = C.NotifyNormal
	NotifyGrab   = C.NotifyGrab
	NotifyUngrab = C.NotifyUngrab
)

const (
	PropertyNewValue   = C.PropertyNewValue
	PropertyDelete     = C.PropertyDelete
	AnyPropertyType    = C.AnyPropertyType
	PropertyChangeMask = C.PropertyChangeMask
)

const (
	Button1 = C.Button1
	Button2 = C.Button2
	Button3 = C.Button3
	Button4 = C.Button4
	Button5 = C.Button5
)

func (s Status) Error() string {
	switch s {
	case Success:
		return "success"
	case NoMemory:
		return "no memory"
	case LocaleNotSupported:
		return "locale not supported"
	case ConverterNotFound:
		return "converter not found"
	case BadAtom:
		return "bad atom"
	case BadRequest:
		return "bad request"
	case BadAccess:
		return "bad access"
	case BadColor:
		return "bad color"
	case BadImplementation:
		return "bad implementation"
	case BadAlloc:
		return "bad alloc"
	case BadFont:
		return "bad font"
	case BadGC:
		return "bad gc"
	case BadMatch:
		return "bad match"
	case BadPixmap:
		return "bad pixmap"
	case BadValue:
		return "bad value"
	case BadWindow:
		return "bad window"
	case BadDrawable:
		return "bad drawable"
	default:
		return "unknown error"
	}
}

func OpenDisplay(name string) *Display {
	var cname *C.char
	if name != "" {
		cname = C.CString(name)
		defer C.free(unsafe.Pointer(cname))
	}
	return (*Display)(C.XOpenDisplay(cname))
}

func CloseDisplay(display *Display) error {
	return xerr(C.XCloseDisplay((*C.Display)(display)))
}

func ConnectionNumber(display *Display) int {
	return int(C.XConnectionNumber((*C.Display)(display)))
}

func CreateGC(display *Display, drawable Drawable, valuemask uint64, values *GCValues) GC {
	return GC(C.XCreateGC((*C.Display)(display), C.Drawable(drawable), C.ulong(valuemask), (*C.XGCValues)(values)))
}

func CreatePixmap(display *Display, drawable Drawable, width, height, depth int) Pixmap {
	return Pixmap(C.XCreatePixmap((*C.Display)(display), C.Drawable(drawable), C.uint(width), C.uint(height), C.uint(depth)))
}

func DisplayName(str string) string {
	var cstr *C.char
	if str == "" {
		cstr = C.CString(str)
		defer C.free(unsafe.Pointer(cstr))
	}
	return C.GoString(C.XDisplayName(cstr))
}

func DefaultRootWindow(display *Display) Window {
	return Window(C.XDefaultRootWindow((*C.Display)(display)))
}

func DefaultScreenOfDisplay(display *Display) *Screen {
	return (*Screen)(C.XDefaultScreenOfDisplay((*C.Display)(display)))
}

func DefaultScreen(display *Display) int {
	return int(C.XDefaultScreen((*C.Display)(display)))
}

func DefaultDepth(display *Display, screen int) int {
	return int(C.XDefaultDepth((*C.Display)(display), C.int(screen)))
}

func DefaultVisual(display *Display, screen_number int) *Visual {
	return (*Visual)(C.XDefaultVisual((*C.Display)(display), C.int(screen_number)))
}

func DefaultCells(display *Display, screen_number int) int {
	return int(C.XDisplayCells((*C.Display)(display), C.int(screen_number)))
}

func DisplayString(display *Display) string {
	return C.GoString(C.XDisplayString((*C.Display)(display)))
}

func ExtendedMaxRequestSize(display *Display) int {
	return int(C.XExtendedMaxRequestSize((*C.Display)(display)))
}

func MaxRequestSize(display *Display) int {
	return int(C.XMaxRequestSize((*C.Display)(display)))
}

func CreateWindow(display *Display, parent Window, x, y, width, height, border_width, depth int, class uint, visual *Visual, valuemask uint64, attribute *SetWindowAttributes) Window {
	return Window(C.XCreateWindow((*C.Display)(display), C.ulong(parent), C.int(x), C.int(y), C.uint(width), C.uint(height), C.uint(border_width), C.int(depth), C.uint(class), (*C.Visual)(unsafe.Pointer(visual)), C.ulong(valuemask), (*C.XSetWindowAttributes)(attribute)))
}

func ParseGeometry(str string) (gm, x, y, w, h int) {
	var cx, cy C.int
	var cw, ch C.uint
	cstr := C.CString(str)
	defer C.free(unsafe.Pointer(cstr))
	gm = int(C.XParseGeometry(cstr, &cx, &cy, &cw, &ch))
	return gm, int(cx), int(cy), int(cw), int(ch)
}

func SetWMHints(display *Display, window Window, wmhints *WMHints) error {
	return xerr(C.XSetWMHints((*C.Display)(display), C.Window(window), (*C.XWMHints)(wmhints)))
}

func GetWMHints(display *Display, window Window) *WMHints {
	return (*WMHints)(C.XGetWMHints((*C.Display)(display), C.Window(window)))
}

func SetWMName(display *Display, window Window, textprop *TextProperty) {
	C.XSetWMName((*C.Display)(display), C.Window(window), (*C.XTextProperty)(textprop))
}

func InternAtom(display *Display, atom_name string, only_if_exists bool) Atom {
	catom_name := C.CString(atom_name)
	defer C.free(unsafe.Pointer(catom_name))
	return Atom(C.XInternAtom((*C.Display)(display), catom_name, xbool(only_if_exists)))
}

func SetSelectionOwner(display *Display, selection Atom, owner Window, time Time) error {
	return xerr(C.XSetSelectionOwner((*C.Display)(display), C.Atom(selection), C.Window(owner), C.Time(time)))
}

func GetSelectionOwner(display *Display, selection Atom) Window {
	return Window(C.XGetSelectionOwner((*C.Display)(display), C.Atom(selection)))
}

func QueryTree(display *Display, window Window) (root, parent Window, children []Window, err error) {
	var croot, cparent C.Window
	var cchildren *C.Window
	var nchildren C.uint
	rc := C.XQueryTree((*C.Display)(display), C.Window(window), &croot, &cparent, &cchildren, &nchildren)
	defer C.XFree(unsafe.Pointer(cchildren))

	if Status(rc) == BadWindow {
		err = xerr(rc)
	}

	root = Window(croot)
	parent = Window(cparent)
	if nchildren > 0 {
		children = make([]Window, nchildren)
		pchildren := (*[1<<27]C.Window)(unsafe.Pointer(cchildren))[:nchildren:nchildren]
		for i := range children {
			children[i] = Window(pchildren[i])
		}
	}
	return
}

func GetWindowAttributes(display *Display, window Window, window_attributes *WindowAttributes) error {
	rc := C.XGetWindowAttributes((*C.Display)(display), C.Window(window), (*C.XWindowAttributes)(window_attributes))
	if Status(rc) == BadDrawable || Status(rc) == BadWindow {
		return xerr(rc)
	}
	return nil
}

func GetWMName(display *Display, window Window, text_prop *TextProperty) error {
	rc := C.XGetWMName((*C.Display)(display), C.Window(window), (*C.XTextProperty)(text_prop))
	if rc == 0 {
		return fmt.Errorf("failed to get wm name property")
	}
	return nil
}

func SetTextProperty(display *Display, window Window, text_prop *TextProperty, property Atom) {
	C.XSetTextProperty((*C.Display)(display), C.Window(window), (*C.XTextProperty)(text_prop), C.Atom(property))
}

func GetTextProperty(display *Display, window Window, text_prop *TextProperty, property Atom) error {
	rc := C.XGetTextProperty((*C.Display)(display), C.Window(window), (*C.XTextProperty)(text_prop), C.Atom(property))
	switch Status(rc) {
	case BadAtom, BadWindow:
		return xerr(rc)
	case 0:
		return fmt.Errorf("failed to get text property")
	}
	return nil
}

func UTF8TextListToTextProperty(display *Display, list []string, style ICCEncodingStyle, text_prop *TextProperty) error {
	var plist **C.char
	clist := make([]*C.char, len(list))
	for i := range clist {
		clist[i] = C.CString(list[i])
		defer C.free(unsafe.Pointer(clist[i]))
	}
	if len(clist) > 0 {
		plist = &clist[0]
	}
	return xerr(C.Xutf8TextListToTextProperty((*C.Display)(display), plist, C.int(len(clist)), C.XICCEncodingStyle(style), (*C.XTextProperty)(text_prop)))
}

func (w *WindowAttributes) OverrideRedirect() bool {
	return w.override_redirect != 0
}

func (w *WindowAttributes) MapState() int {
	return int(w.map_state)
}

func (t *TextProperty) Free() {
	C.XFree(unsafe.Pointer(t.value))
}

func (t *TextProperty) NumItems() int {
	return int(t.nitems)
}

func XmbTextPropertyToTextList(display *Display, text_prop *TextProperty) ([]string, error) {
	var list **C.char
	var count C.int
	rc := C.XmbTextPropertyToTextList((*C.Display)(display), (*C.XTextProperty)(text_prop), &list, &count)
	if rc == 0 {
		return nil, fmt.Errorf("failed to get text list")
	}

	switch Status(rc) {
	case NoMemory, LocaleNotSupported, ConverterNotFound:
		return nil, xerr(rc)
	}
	plist := (*[1 << 27]*C.char)(unsafe.Pointer(list))[:count:count]

	str := make([]string, count)
	for i := range str {
		str[i] = C.GoString(plist[i])
	}
	C.XFreeStringList(list)

	return str, nil
}

func WarpPointer(display *Display, src_w, dest_w Window, src_x, src_y, src_width, src_height, dest_x, dest_y int) error {
	rc := C.XWarpPointer((*C.Display)(display), C.Window(src_w), C.Window(dest_w), C.int(src_x), C.int(src_y), C.uint(src_width), C.uint(src_height), C.int(dest_x), C.int(dest_y))
	return xerr(rc)
}

func RootWindow(display *Display, screen int) Window {
	return Window(C.XRootWindow((*C.Display)(display), C.int(screen)))
}

func DisplayWidth(display *Display, screen int) int {
	return int(C.XDisplayWidth((*C.Display)(display), C.int(screen)))
}

func DisplayHeight(display *Display, screen int) int {
	return int(C.XDisplayHeight((*C.Display)(display), C.int(screen)))
}

func SetInputFocus(display *Display, focus Window, revert_to int, time Time) {
	C.XSetInputFocus((*C.Display)(display), C.Window(focus), C.int(revert_to), C.Time(time))
}

func ScreenSaverQueryExtension(display *Display) (supported bool, base, errbase int) {
	var cbase, cerrbase C.int
	rc := C.XScreenSaverQueryExtension((*C.Display)(display), &cbase, &cerrbase)
	return rc != 0, int(cbase), int(cerrbase)
}

func ScreenSaverAllocInfo() *ScreenSaverInfo {
	return (*ScreenSaverInfo)(C.XScreenSaverAllocInfo())
}

func ScreenSaverQueryInfo(display *Display, drawable Drawable, saver_info *ScreenSaverInfo) error {
	rc := C.XScreenSaverQueryInfo((*C.Display)(display), C.Drawable(drawable), (*C.XScreenSaverInfo)(saver_info))
	if rc == 0 {
		return fmt.Errorf("screensaver extension not supported")
	}
	return nil
}

func (s *ScreenSaverInfo) State() int {
	return int(s.state)
}

func (s *ScreenSaverInfo) Idle() int {
	return int(s.idle)
}

func (s *ScreenSaverInfo) TilOrSince() int {
	return int(s.til_or_since)
}

func CreateSimpleWindow(display *Display, parent Window, x, y, width, height, border_width int, border, background uint) Window {
	return Window(C.XCreateSimpleWindow((*C.Display)(display), C.Window(parent), C.int(x), C.int(y), C.uint(width), C.uint(height), C.uint(border_width), C.ulong(border), C.ulong(background)))
}

func ConvertSelection(display *Display, selection, target, property Atom, requestor Window, time Time) error {
	return xerr(C.XConvertSelection((*C.Display)(display), C.Atom(selection), C.Atom(target), C.Atom(property), C.Window(requestor), C.Time(time)))
}

func NextEvent(display *Display, ev *Event) {
	C.XNextEvent((*C.Display)(display), (*C.XEvent)(ev))
}

func (ev *Event) Type() int {
	return int(C.ev_type((*C.XEvent)(ev)))
}

func (ev *Event) Visibility() *VisibilityEvent {
	return (*VisibilityEvent)(C.ev_xvisibility((*C.XEvent)(ev)))
}

func (ev *Event) Selection() *SelectionEvent {
	return (*SelectionEvent)(C.ev_xselection((*C.XEvent)(ev)))
}

func (ev *Event) SelectionClear() *SelectionClearEvent {
	return (*SelectionClearEvent)(C.ev_xselectionclear((*C.XEvent)(ev)))
}

func (ev *Event) SelectionRequest() *SelectionRequestEvent {
	return (*SelectionRequestEvent)(C.ev_xselectionrequest((*C.XEvent)(ev)))
}

func (ev *Event) Configure() *ConfigureEvent {
	return (*ConfigureEvent)(C.ev_xconfigure((*C.XEvent)(ev)))
}

func (ev *SelectionEvent) Property() Atom {
	return Atom(ev.property)
}

func (ev *VisibilityEvent) State() int {
	return int(ev.state)
}

func (ev *ConfigureEvent) Width() int {
	return int(ev.width)
}

func (ev *ConfigureEvent) Height() int {
	return int(ev.height)
}

func DefaultColormap(display *Display, scr int) Colormap {
	return Colormap(C.XDefaultColormap((*C.Display)(display), C.int(scr)))
}

func GetWindowProperty(display *Display, window Window, property Atom, offset, length int, delete_ bool, req_type Atom) (actual_type Atom, actual_format, nitems, bytes_after int, prop []byte, err error) {
	var cactual_type C.Atom
	var cactual_format C.int
	var cnitems, cbytes_after C.ulong
	var cprop *C.uchar
	rc := C.XGetWindowProperty((*C.Display)(display), C.Window(window), C.Atom(property), C.long(offset), C.long(length), xbool(delete_), C.Atom(req_type), &cactual_type, &cactual_format, &cnitems, &cbytes_after, &cprop)
	actual_type = Atom(cactual_type)
	actual_format = int(cactual_format)
	nitems = int(cnitems)
	bytes_after = int(cbytes_after)
	err = xerr(rc)
	if cprop != nil {
		pprop := (*[math.MaxInt32]C.uchar)(unsafe.Pointer(cprop))[:nitems:nitems]
		prop = make([]byte, nitems)
		for i := range prop {
			prop[i] = byte(pprop[i])
		}
	}
	return
}

func SupportsLocale() bool {
	return C.XSupportsLocale() != 0
}

func SetLocaleModifiers(str string) string {
	cstr := C.CString(str)
	defer C.free(unsafe.Pointer(cstr))
	return C.GoString(C.XSetLocaleModifiers(cstr))
}

func ListExtensions(display *Display) []string {
	var nextensions C.int
	list := C.XListExtensions((*C.Display)(display), &nextensions)
	plist := (*[1 << 27]*C.char)(unsafe.Pointer(list))[:nextensions:nextensions]
	defer C.XFreeExtensionList(list)

	exts := make([]string, nextensions)
	for i := range exts {
		exts[i] = C.GoString(plist[i])
	}
	return exts
}

func QueryExtension(display *Display, name string) (supported bool, major_opcode, first_event, first_error int) {
	var cmajor_opcode, cfirst_event, cfirst_error C.int
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	rc := C.XQueryExtension((*C.Display)(display), cname, &cmajor_opcode, &cfirst_event, &cfirst_error)
	return rc != 0, int(cmajor_opcode), int(cfirst_event), int(cfirst_error)
}

func SelectInput(display *Display, window Window, event_mask int64) error {
	return xerr(C.XSelectInput((*C.Display)(display), C.Window(window), C.long(event_mask)))
}

func (p *SetWindowAttributes) EventMask() int64 {
	return int64(p.event_mask)
}

func (p *SetWindowAttributes) SetEventMask(event_mask int64) {
	p.event_mask = C.long(event_mask)
}

func (p *SetWindowAttributes) SetBackgroundPixel(background_pixel uint64) {
	p.background_pixel = C.ulong(background_pixel)
}

func (p *SetWindowAttributes) SetBorderPixel(border_pixel uint64) {
	p.border_pixel = C.ulong(border_pixel)
}

func (p *SetWindowAttributes) SetBitGravity(bit_gravity int) {
	p.bit_gravity = C.int(bit_gravity)
}

func (p *SetWindowAttributes) SetColormap(colormap Colormap) {
	p.colormap = C.ulong(colormap)
}

func ChangeWindowAttributes(display *Display, window Window, valuemask uint64, attributes *SetWindowAttributes) error {
	return xerr(C.XChangeWindowAttributes((*C.Display)(display), C.Window(window), C.ulong(valuemask), (*C.XSetWindowAttributes)(attributes)))
}

func CopyArea(display *Display, src, dest Drawable, gc GC, src_x, src_y, width, height, dest_x, dest_y int) {
	C.XCopyArea((*C.Display)(display), C.Drawable(src), C.Drawable(dest), C.GC(gc), C.int(src_x), C.int(src_y), C.uint(width), C.uint(height), C.int(dest_x), C.int(dest_y))
}

func SetState(display *Display, gc GC, foreground, background uint64, function int, plane_mask uint64) error {
	return xerr(C.XSetState((*C.Display)(display), C.GC(gc), C.ulong(foreground), C.ulong(background), C.int(function), C.ulong(plane_mask)))
}

func SetFunction(display *Display, gc GC, function int) error {
	return xerr(C.XSetFunction((*C.Display)(display), C.GC(gc), C.int(function)))
}

func SetPlaneMask(display *Display, gc GC, plane_mask uint64) error {
	return xerr(C.XSetPlaneMask((*C.Display)(display), C.GC(gc), C.ulong(plane_mask)))
}

func SetForeground(display *Display, gc GC, foreground uint64) error {
	return xerr(C.XSetForeground((*C.Display)(display), C.GC(gc), C.ulong(foreground)))
}

func SetBackground(display *Display, gc GC, background uint64) error {
	return xerr(C.XSetBackground((*C.Display)(display), C.GC(gc), C.ulong(background)))
}

func Flush(display *Display) {
	C.XFlush((*C.Display)(display))
}

func (e *Event) Button() *ButtonEvent {
	return (*ButtonEvent)(C.ev_xbutton((*C.XEvent)(e)))
}

func (e *Event) Key() *KeyEvent {
	return (*KeyEvent)(C.ev_xkey((*C.XEvent)(e)))
}

func (e *Event) Focus() *FocusChangeEvent {
	return (*FocusChangeEvent)(C.ev_xfocus((*C.XEvent)(e)))
}

func (e *Event) CreateWindow() *CreateWindowEvent {
	return (*CreateWindowEvent)(C.ev_xcreatewindow((*C.XEvent)(e)))
}

func (e *Event) Cookie() *GenericEventCookie {
	return (*GenericEventCookie)(C.ev_xcookie((*C.XEvent)(e)))
}

func (e *Event) Client() *ClientMessageEvent {
	return (*ClientMessageEvent)(C.ev_xclient((*C.XEvent)(e)))
}

func (e *Event) Property() *PropertyEvent {
	return (*PropertyEvent)(C.ev_xproperty((*C.XEvent)(e)))
}

func (e *PropertyEvent) State() uint {
	return uint(e.state)
}

func (e *PropertyEvent) Atom() Atom {
	return Atom(e.atom)
}

func (e *KeyEvent) State() uint {
	return uint(e.state)
}

func (e *CreateWindowEvent) Parent() Window {
	return Window(e.parent)
}

func (e *GenericEventCookie) Data() unsafe.Pointer {
	return e.data
}

func (e *ErrorEvent) ErrorCode() Status {
	return Status(e.error_code)
}

func (e *ClientMessageEvent) MessageType() Atom {
	return Atom(e.message_type)
}

func (e *ClientMessageEvent) Format() int {
	return int(e.format)
}

func (e *ClientMessageEvent) Long() [5]uint64 {
	var l [5]C.long
	var r [5]uint64
	C.ev_xclient_long((*C.XClientMessageEvent)(e), &l[0], C.int(len(l)))
	for i := range l {
		r[i] = uint64(l[i])
	}
	return r
}

func (e *ButtonEvent) Type() int {
	return int(e._type)
}

func (e *ButtonEvent) Button() uint {
	return uint(e.button)
}

func (e *ButtonEvent) Time() Time {
	return Time(e.time)
}

func (e *ButtonEvent) State() uint {
	return uint(e.state)
}

func (e *ButtonEvent) X() int {
	return int(e.x)
}

func (e *ButtonEvent) Y() int {
	return int(e.y)
}

func (e *FocusChangeEvent) Mode() int {
	return int(e.mode)
}

func (e *SelectionEvent) SetType(typ int) {
	e._type = C.int(typ)
}

func (e *SelectionEvent) SetProperty(property Atom) {
	e.property = C.Atom(property)
}

func (e *SelectionEvent) SetRequestor(requestor Window) {
	e.requestor = C.Window(requestor)
}

func (e *SelectionEvent) SetSelection(selection Atom) {
	e.selection = C.Atom(selection)
}

func (e *SelectionEvent) SetTarget(target Atom) {
	e.target = C.Atom(target)
}

func (e *SelectionEvent) SetTime(time Time) {
	e.time = C.Time(time)
}

func (e *SelectionEvent) Cast() *Event {
	return (*Event)(unsafe.Pointer(e))
}

func (e *SelectionRequestEvent) SetProperty(property Atom) {
	e.property = C.Atom(property)
}

func (e *SelectionRequestEvent) SetTarget(target Atom) {
	e.target = C.Atom(target)
}

func (e *SelectionRequestEvent) Display() *Display {
	return (*Display)(e.display)
}

func (e *SelectionRequestEvent) Target() Atom {
	return Atom(e.target)
}

func (e *SelectionRequestEvent) Owner() Window {
	return Window(e.owner)
}

func (e *SelectionRequestEvent) Requestor() Window {
	return Window(e.requestor)
}

func (e *SelectionRequestEvent) Time() Time {
	return Time(e.time)
}

func (e *SelectionRequestEvent) Selection() Atom {
	return Atom(e.selection)
}

func (e *SelectionRequestEvent) Property() Atom {
	return Atom(e.property)
}

func GetEventData(display *Display, cookie *GenericEventCookie) {
	C.XGetEventData((*C.Display)(display), (*C.XGenericEventCookie)(cookie))
}

func FreeEventData(display *Display, cookie *GenericEventCookie) {
	C.XFreeEventData((*C.Display)(display), (*C.XGenericEventCookie)(cookie))
}

func SetErrorHandler(f func(d *Display, ev *ErrorEvent) int) {
	errorHandler = f
	C.xset_error_handler()
}

var errorHandler func(d *Display, ev *ErrorEvent) int

//export goErrorHandler
func goErrorHandler(display *C.Display, ev *C.XErrorEvent) C.int {
	return C.int(errorHandler((*Display)(display), (*ErrorEvent)(ev)))
}

func FillRectangle(display *Display, drawable Drawable, gc GC, x, y, width, height int) error {
	return xerr(C.XFillRectangle((*C.Display)(display), C.Drawable(drawable), C.GC(gc), C.int(x), C.int(y), C.uint(width), C.uint(height)))
}

func OpenIM(display *Display, db xrm.Database, res_name string, res_class string) IM {
	var cres_name, cres_class *C.char
	if res_name != "" {
		cres_name = C.CString(res_name)
		defer C.free(unsafe.Pointer(cres_name))
	}
	if res_class != "" {
		cres_class = C.CString(res_class)
		defer C.free(unsafe.Pointer(cres_class))
	}
	return IM(C.XOpenIM((*C.Display)(display), C.XrmDatabase(unsafe.Pointer(db)), cres_name, cres_class))
}

func CreateFontCursor(display *Display, shape uint) Cursor {
	return Cursor(C.XCreateFontCursor((*C.Display)(display), C.uint(shape)))
}

func DefineCursor(display *Display, window Window, cursor Cursor) error {
	return xerr(C.XDefineCursor((*C.Display)(display), C.Window(window), C.Cursor(cursor)))
}

func UndefineCursor(display *Display, window Window) error {
	return xerr(C.XUndefineCursor((*C.Display)(display), C.Window(window)))
}

func (c *Color) SetRed(red uint16) {
	c.red = C.ushort(red)
}

func (c *Color) SetGreen(green uint16) {
	c.green = C.ushort(green)
}

func (c *Color) SetBlue(blue uint16) {
	c.blue = C.ushort(blue)
}

func ParseColor(display *Display, colormap Colormap, spec string, color *Color) error {
	cspec := C.CString(spec)
	defer C.free(unsafe.Pointer(cspec))
	rc := C.XParseColor((*C.Display)(display), C.Colormap(colormap), cspec, (*C.XColor)(color))
	if Status(rc) == BadColor {
		return xerr(rc)
	}
	if rc != 0 {
		return nil
	}
	return fmt.Errorf("failed to resolve color")
}

func RecolorCursor(display *Display, cursor Cursor, fg, bg *Color) error {
	return xerr(C.XRecolorCursor((*C.Display)(display), C.Cursor(cursor), (*C.XColor)(fg), (*C.XColor)(bg)))
}

func FreeCursor(display *Display, cursor Cursor) error {
	return xerr(C.XFreeCursor((*C.Display)(display), C.Cursor(cursor)))
}

func QueryBestCursor(display *Display, drawable Drawable, width, height int) (int, int, error) {
	var cwidth, cheight C.uint
	rc := C.XQueryBestCursor((*C.Display)(display), C.Drawable(drawable), C.uint(width), C.uint(height), &cwidth, &cheight)
	return int(cwidth), int(cheight), xerr(rc)
}

func SetWMProperties(display *Display, window Window, window_name, icon_name *TextProperty, args []string, normal_hints *SizeHints, wm_hints *WMHints, class_hints *ClassHint) {
	var parg **C.char
	cargs := make([]*C.char, len(args))
	for i := range cargs {
		cargs[i] = C.CString(args[i])
		defer C.free(unsafe.Pointer(cargs[i]))
	}
	if len(cargs) > 0 {
		parg = &cargs[0]
	}
	C.XSetWMProperties((*C.Display)(display), C.Window(window), (*C.XTextProperty)(window_name), (*C.XTextProperty)(icon_name), parg, C.int(len(cargs)), (*C.XSizeHints)(normal_hints), (*C.XWMHints)(wm_hints), (*C.XClassHint)(class_hints))
}

func SetWMProtocols(display *Display, window Window, protocols []Atom) error {
	rc := C.XSetWMProtocols((*C.Display)(display), C.Window(window), (*C.Atom)(&protocols[0]), C.int(len(protocols)))
	switch Status(rc) {
	case BadAlloc, BadWindow:
		return xerr(rc)
	case 0:
		return fmt.Errorf("failed to set wm protocol")
	}
	return nil
}

func ChangeProperty(display *Display, window Window, property, typ Atom, format, mode int, data interface{}) error {
	switch v := data.(type) {
	case int:
		vi := C.int(v)
		return xerr(C.XChangeProperty((*C.Display)(display), C.Window(window), C.Atom(property), C.Atom(typ), C.int(format), C.int(mode), (*C.uchar)(unsafe.Pointer(&vi)), C.sizeof_int))

	case []uint8:
		if len(v) > 0 {
			return xerr(C.XChangeProperty((*C.Display)(display), C.Window(window), C.Atom(property), C.Atom(typ), C.int(format), C.int(mode), (*C.uchar)(unsafe.Pointer(&v[0])), C.int(len(v))))
		}
		return nil

	default:
		panic(fmt.Errorf("unsupported type %T", v))
	}
}

func MapWindow(display *Display, window Window) error {
	return xerr(C.XMapWindow((*C.Display)(display), C.Window(window)))
}

func Sync(display *Display, discard bool) {
	C.XSync((*C.Display)(display), xbool(discard))
}

func AllocSizeHints() *SizeHints {
	return (*SizeHints)(C.XAllocSizeHints())
}

func (s *SizeHints) Free() {
	C.XFree(unsafe.Pointer(s))
}

func (s *SizeHints) SetFlags(flags int64) {
	s.flags = C.long(flags)
}

func (s *SizeHints) SetWidth(width int) {
	s.width = C.int(width)
}

func (s *SizeHints) SetHeight(height int) {
	s.height = C.int(height)
}

func (s *SizeHints) SetWidthInc(width_inc int) {
	s.width_inc = C.int(width_inc)
}

func (s *SizeHints) SetHeightInc(height_inc int) {
	s.height_inc = C.int(height_inc)
}

func (s *SizeHints) SetBaseHeight(base_height int) {
	s.base_height = C.int(base_height)
}

func (s *SizeHints) SetBaseWidth(base_width int) {
	s.base_width = C.int(base_width)
}

func (s *SizeHints) SetMinHeight(min_height int) {
	s.min_height = C.int(min_height)
}

func (s *SizeHints) SetMinWidth(min_width int) {
	s.min_width = C.int(min_width)
}

func (s *SizeHints) SetMaxHeight(max_height int) {
	s.max_height = C.int(max_height)
}

func (s *SizeHints) SetMaxWidth(max_width int) {
	s.max_width = C.int(max_width)
}

func (s *SizeHints) SetWinGravity(win_gravity int) {
	s.win_gravity = C.int(win_gravity)
}

func (s *SizeHints) Flags() int64 {
	return int64(s.flags)
}

func (s *SizeHints) SetX(x int) {
	s.x = C.int(x)
}

func (s *SizeHints) SetY(y int) {
	s.y = C.int(y)
}

func (c *ClassHint) SetResName(name string) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	c.res_name = cname
}

func (c *ClassHint) SetResClass(class string) {
	cclass := C.CString(class)
	defer C.free(unsafe.Pointer(cclass))
	c.res_class = cclass
}

func (w *WMHints) Flags() int64 {
	return int64(w.flags)
}

func (w *WMHints) Free() {
	C.XFree(unsafe.Pointer(w))
}

func (w *WMHints) SetFlags(flags int64) {
	w.flags = C.long(flags)
}

func (w *WMHints) SetInput(input int) {
	w.input = C.int(input)
}

func Pending(display *Display) int {
	return int(C.XPending((*C.Display)(display)))
}

func EventsQueued(display *Display, mode int) int {
	return int(C.XEventsQueued((*C.Display)(display), C.int(mode)))
}

func QLength(display *Display) int {
	return int(C.XQLength((*C.Display)(display)))
}

func ScreenCount(display *Display) int {
	return int(C.XScreenCount((*C.Display)(display)))
}

func ServerVendor(display *Display) string {
	return C.GoString(C.XServerVendor((*C.Display)(display)))
}

func VendorRelease(display *Display) int {
	return int(C.XVendorRelease((*C.Display)(display)))
}

func FilterEvent(event *Event, window Window) bool {
	return C.XFilterEvent((*C.XEvent)(event), C.Window(window)) != 0
}

var (
	imInstanceCallbackID    uint64
	imInstanceCallbacksName [8]string
	imInstanceCallbacks     [8]IDProc
	imCallbacks             [8]func(IM, Pointer)
)

//export goIMInstantitateCallback
func goIMInstantitateCallback(display *C.Display, client, call C.XPointer) {
	id := uintptr(unsafe.Pointer(client))
	if int(id) >= len(imInstanceCallbacks) || imInstanceCallbacks[id] == nil {
		return
	}
	imInstanceCallbacks[id]((*Display)(display), Pointer(call))
}

//export goIMDestroyCallback
func goIMDestroyCallback(xim C.XIM, call C.XPointer) {
	imCallbacks[0](IM(xim), Pointer(call))
}

func RegisterIMInstantiateCallback(display *Display, db xrm.Database, res_name, res_class, callback_name string, callback IDProc) bool {
	var cres_name, cres_class *C.char
	if res_name != "" {
		cres_name = C.CString(res_name)
		defer C.free(unsafe.Pointer(cres_name))
	}
	if res_class != "" {
		cres_class = C.CString(res_class)
		defer C.free(unsafe.Pointer(cres_class))
	}

	imInstanceCallbackID++
	imInstanceCallbacksName[imInstanceCallbackID] = callback_name
	return C.xregister_im_instantitate_callback((*C.Display)(display), C.XrmDatabase(unsafe.Pointer(db)), cres_name, cres_class, C.ulong(imInstanceCallbackID)) != 0
}

func UnregisterIMInstantiateCallback(display *Display, db xrm.Database, res_name, res_class, callback_name string, callback IDProc) bool {
	for i, name := range imInstanceCallbacksName {
		if name == callback_name {
			imInstanceCallbacks[i] = nil
			return true
		}
	}
	return false
}

func DisplayOfIM(im IM) *Display {
	return (*Display)(C.XDisplayOfIM(im))
}

func LocaleOfIM(im IM) string {
	return C.GoString((C.XLocaleOfIM(im)))
}

func CloseIM(im IM) error {
	return xerr(C.XCloseIM(im))
}

func SetIMValues(xim IM, values []IMValue) string {
	var res string
	for _, v := range values {
		var cres *C.char
		key := v.Key
		ckey := C.CString(key)
		switch key {
		case NDestroyCallback:
			var callback C.XIMCallback
			callback.callback = C.XIMProc(C.xim_destroy_callback)
			cres = C.xset_im_values_void(C.XIM(xim), ckey, unsafe.Pointer(&callback))
			imCallbacks[0] = v.Value.(func(IM, Pointer))
		default:
			panic(fmt.Errorf("unsupported key %q", key))
		}
		C.free(unsafe.Pointer(ckey))

		res = ""
		if cres != nil {
			res = key
			break
		}
	}
	return res
}

func (v *VaNestedList) Free() {
	C.XFree(unsafe.Pointer(v.list))
	for _, p := range v.key {
		C.free(unsafe.Pointer(p))
	}
	for _, p := range v.val {
		C.free(p)
	}
}

func VaCreateNestedList(values []ICValue) *VaNestedList {
	var key []*C.char
	var val []unsafe.Pointer
	for _, v := range values {
		key = append(key, C.CString(v.Key))
		switch p := v.Value.(type) {
		case *Point:
			cp := C.calloc(1, C.sizeof_XPoint)
			C.memcpy(cp, unsafe.Pointer(p), C.sizeof_XPoint)
			val = append(val, cp)
		default:
			panic(fmt.Errorf("unknown value type %T", v.Value))
		}
	}

	vl := &VaNestedList{
		key: key,
		val: val,
	}
	switch len(key) {
	case 1:
		vl.list = C.va_create_nested_list1(key[0], val[0])
	default:
		panic(fmt.Errorf("unsupported argument length %d", len(key)))
	}
	return vl
}

func CreateIC(xim IM, values []ICValue) IC {
	var clientWindow, focusWindow C.Window
	var inputStyle C.long
	for _, v := range values {
		switch v.Key {
		case NInputStyle:
			inputStyle = C.long(v.Value.(int))
		case NClientWindow:
			clientWindow = C.Window(v.Value.(Window))
		case NFocusWindow:
			focusWindow = C.Window(v.Value.(Window))
		}
	}
	return IC(C.xcreateic(xim, inputStyle, clientWindow, focusWindow))
}

func SetICValues(xic IC, values []ICValue) string {
	var res string
	for _, v := range values {
		var cres *C.char
		key := v.Key
		ckey := C.CString(key)
		switch key {
		case NPreeditAttributes:
			vl := v.Value.(*VaNestedList)
			cres = C.xset_ic_values_void(C.XIC(xic), ckey, unsafe.Pointer(vl.list))
		default:
			panic(fmt.Errorf("unsupported key %q", key))
		}
		C.free(unsafe.Pointer(ckey))

		res = ""
		if cres != nil {
			res = key
			break
		}
	}
	return res
}

func DestroyIC(ic IC) {
	C.XDestroyIC(C.XIC(ic))
}

func IMOfIC(ic IC) IM {
	return IM(C.XIMOfIC(ic))
}

func FreePixmap(display *Display, drawable Drawable) {
	C.XFreePixmap((*C.Display)(display), C.Drawable(drawable))
}

func (p *Point) SetX(x int) {
	p.x = C.short(x)
}

func (p *Point) SetY(y int) {
	p.y = C.short(y)
}

func (p *Point) X() int {
	return int(p.x)
}

func (p *Point) Y() int {
	return int(p.y)
}

func InitThreads() error {
	rc := C.XInitThreads()
	if rc == 0 {
		return fmt.Errorf("failed to init X threads")
	}
	return nil
}

func LockDisplay(display *Display) {
	C.XLockDisplay((*C.Display)(display))
}

func UnlockDisplay(display *Display) {
	C.XUnlockDisplay((*C.Display)(display))
}

func SetICFocus(ic IC) {
	C.XSetICFocus(C.XIC(ic))
}

func UnsetICFocus(ic IC) {
	C.XUnsetICFocus(C.XIC(ic))
}

func (r *Rectangle) SetX(x int) {
	r.x = C.short(x)
}

func (r *Rectangle) SetY(y int) {
	r.y = C.short(y)
}

func (r *Rectangle) SetWidth(width int) {
	r.width = C.ushort(width)
}

func (r *Rectangle) SetHeight(height int) {
	r.height = C.ushort(height)
}

func PutBackEvent(display *Display, event *Event) error {
	return xerr(C.XPutBackEvent((*C.Display)(display), (*C.XEvent)(event)))
}

func SendEvent(display *Display, window Window, propagate bool, event_mask uint64, event_send *Event) error {
	rc := C.XSendEvent((*C.Display)(display), C.Window(window), xbool(propagate), C.long(event_mask), (*C.XEvent)(event_send))
	switch Status(rc) {
	case 0:
		rc = C.int(BadValue)
		fallthrough
	case BadValue, BadWindow:
		return Status(rc)
	}
	return nil
}

func DisplayMotionBufferSize(display *Display) uint64 {
	return uint64(C.XDisplayMotionBufferSize((*C.Display)(display)))
}

func XmbLookupString(ic IC, event *KeyPressedEvent) (string, KeySym, Status) {
	var status C.Status
	var keysym C.KeySym
	var buf [128]C.char
	n := C.XmbLookupString(C.XIC(ic), (*C.XKeyPressedEvent)(event), &buf[0], C.int(len(buf)), &keysym, &status)
	str := ""
	if n > 0 {
		str = C.GoString(&buf[:n][0])
	}
	return str, KeySym(keysym), Status(status)
}

func DeleteProperty(display *Display, window Window, property Atom) error {
	return xerr(C.XDeleteProperty((*C.Display)(display), C.Window(window), C.Atom(property)))
}

func xbool(b bool) C.Bool {
	if b {
		return 1
	}
	return 0
}

func xerr(code C.int) error {
	if code == 0 {
		return nil
	}
	return Status(code)
}
