package xcb

/*
#include <stdlib.h>
#include <xcb/xcb.h>

#cgo pkg-config: xcb
*/
import "C"
import "unsafe"

type (
	AuthInfo       C.xcb_auth_info_t
	Connection     C.xcb_connection_t
	Drawable       C.xcb_drawable_t
	Extension      C.xcb_extension_t
	GC             C.xcb_gc_t
	GContext       C.xcb_gcontext_t
	Screen         C.xcb_screen_t
	ScreenIterator C.xcb_screen_iterator_t
	Window         C.xcb_window_t
	GenericEvent   C.xcb_generic_event_t
	GenericError   C.xcb_generic_error_t
	VoidCookie     C.xcb_void_cookie_t
	Error          int
)

const (
	GC_FOREGROUND = C.XCB_GC_FOREGROUND
	GC_BACKGROUND = C.XCB_GC_BACKGROUND
	GC_FONT       = C.XCB_GC_FONT
)

const (
	EXPOSE      = C.XCB_EXPOSE
	KEY_RELEASE = C.XCB_KEY_RELEASE
)

const (
	CW_EVENT_MASK = C.XCB_CW_EVENT_MASK
	CW_BACK_PIXEL = C.XCB_CW_BACK_PIXEL
)

const (
	EVENT_MASK_KEY_RELEASE    = C.XCB_EVENT_MASK_KEY_RELEASE
	EVENT_MASK_BUTTON_PRESS   = C.XCB_EVENT_MASK_BUTTON_PRESS
	EVENT_MASK_EXPOSURE       = C.XCB_EVENT_MASK_EXPOSURE
	EVENT_MASK_POINTER_MOTION = C.XCB_EVENT_MASK_POINTER_MOTION
)

const (
	CONN_ERROR                   Error = C.XCB_CONN_ERROR
	CONN_CLOSED_EXT_NOTSUPPORTED Error = C.XCB_CONN_CLOSED_EXT_NOTSUPPORTED
	CONN_CLOSED_MEM_INSUFFICIENT Error = C.XCB_CONN_CLOSED_MEM_INSUFFICIENT
	CONN_CLOSED_REQ_LEN_EXCEED   Error = C.XCB_CONN_CLOSED_REQ_LEN_EXCEED
	CONN_CLOSED_PARSE_ERR        Error = C.XCB_CONN_CLOSED_PARSE_ERR
	CONN_CLOSED_INVALID_SCREEN   Error = C.XCB_CONN_CLOSED_INVALID_SCREEN
)

func Connect(displayname string) (*Connection, int) {
	var cdisplayname *C.char
	if displayname != "" {
		cdisplayname = C.CString(displayname)
		defer C.free(unsafe.Pointer(cdisplayname))
	}
	var cscreen C.int
	conn := C.xcb_connect(cdisplayname, &cscreen)
	return (*Connection)(conn), int(cscreen)
}

func Disconnect(c *Connection) {
	C.xcb_disconnect((*C.xcb_connection_t)(c))
}

func Flush(c *Connection) int {
	return int(C.xcb_flush((*C.xcb_connection_t)(c)))
}

func ConnectionHasError(c *Connection) error {
	return nil
}

func GenerateID(c *Connection) uint32 {
	return uint32(C.xcb_generate_id((*C.xcb_connection_t)(c)))
}

func GetFileDescriptor(c *Connection) int {
	return int(C.xcb_get_file_descriptor((*C.xcb_connection_t)(c)))
}

func GetMaximumRequestLength(c *Connection) int {
	return int(C.xcb_get_maximum_request_length((*C.xcb_connection_t)(c)))
}

func WaitForEvent(c *Connection) *GenericEvent {
	return (*GenericEvent)(C.xcb_wait_for_event((*C.xcb_connection_t)(c)))
}

func MapWindow(c *Connection, w Window) VoidCookie {
	return VoidCookie(C.xcb_map_window((*C.xcb_connection_t)(c), C.xcb_window_t(w)))
}
