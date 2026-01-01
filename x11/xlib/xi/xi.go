package xi

/*
#include <stdlib.h>
#include <X11/Xlib.h>
#include <X11/Xutil.h>
#include <X11/extensions/XInput2.h>
#include <X11/extensions/XInput.h>

#cgo pkg-config: x11
#cgo LDFLAGS: -lXfixes -lXi

void xi_set_mask(unsigned char *mask, unsigned long event) {
	XISetMask(mask, event);
}

#define DEVICE_NOTIFY_FUNC(name, func) \
void name(XDevice *device, int *type, XEventClass *event) { \
	int ttype = *type; \
	XEventClass tevent = *event; \
	func(device, ttype, tevent); \
	*type = ttype; \
	*event = tevent; \
}

DEVICE_NOTIFY_FUNC(device_key_press, DeviceKeyPress)
DEVICE_NOTIFY_FUNC(device_key_release, DeviceKeyRelease)
DEVICE_NOTIFY_FUNC(device_button_press, DeviceButtonPress)
DEVICE_NOTIFY_FUNC(device_button_release, DeviceButtonRelease)
DEVICE_NOTIFY_FUNC(device_motion_notify, DeviceMotionNotify)
*/
import "C"
import (
	"math"
	"unsafe"

	"github.com/qeedquan/go-media/x11/xlib"
)

type (
	Device         C.XDevice
	DeviceInfo     C.XDeviceInfo
	EventClass     C.XEventClass
	InputClassInfo C.XInputClassInfo
	DeviceEvent    C.XIDeviceEvent
)

const (
	AllMasterDevices = C.XIAllMasterDevices
)

const (
	RawMotion        = C.XI_RawMotion
	RawButtonPress   = C.XI_RawButtonPress
	RawButtonRelease = C.XI_RawButtonRelease
)

const (
	LASTEVENT = C.XI_LASTEVENT
)

const (
	IsXExtensionKeyboard = C.IsXExtensionKeyboard
	IsXExtensionPointer  = C.IsXExtensionPointer
)

const (
	KeyClass      = C.KeyClass
	ButtonClass   = C.ButtonClass
	ValuatorClass = C.ValuatorClass
)

type EventMask struct {
	DeviceID int
	Mask     []uint8
}

func QueryVersion(display *xlib.Display, major, minor int) (int, int, error) {
	cmajor := C.int(major)
	cminor := C.int(minor)
	rc := C.XIQueryVersion((*C.Display)(display), &cmajor, &cminor)
	return int(cmajor), int(cminor), xerr(rc)
}

func SelectEvents(display *xlib.Display, window xlib.Window, masks []EventMask) {
	cmasks := make([]C.XIEventMask, len(masks))
	for i := range cmasks {
		cmasks[i].deviceid = C.int(masks[i].DeviceID)
		cmasks[i].mask_len = C.int(len(masks[i].Mask))
		cmasks[i].mask = (*C.uchar)(C.CBytes(masks[i].Mask))
	}
	C.XISelectEvents((*C.Display)(display), C.Window(window), &cmasks[0], C.int(len(cmasks)))
	for i := range cmasks {
		C.free(unsafe.Pointer(cmasks[i].mask))
	}
}

func SetMask(mask []uint8, event uint64) {
	C.xi_set_mask((*C.uchar)(unsafe.Pointer(&mask[0])), C.ulong(event))
}

func ListInputDevices(display *xlib.Display) []DeviceInfo {
	var ndevices C.int
	devinfo := C.XListInputDevices((*C.Display)(display), &ndevices)
	return (*[math.MaxInt32]DeviceInfo)(unsafe.Pointer(devinfo))[:ndevices:ndevices]
}

func OpenDevice(display *xlib.Display, device_id xlib.ID) *Device {
	return (*Device)(C.XOpenDevice((*C.Display)(display), C.XID(device_id)))
}

func CloseDevice(display *xlib.Display, device *Device) {
	C.XCloseDevice((*C.Display)(display), (*C.XDevice)(device))
}

func DeviceKeyPress(device *Device, typ int, event EventClass) (int, EventClass) {
	ctyp := C.int(typ)
	cevent := C.XEventClass(event)
	C.device_key_press((*C.XDevice)(device), &ctyp, &cevent)
	return int(ctyp), EventClass(cevent)
}

func DeviceKeyRelease(device *Device, typ int, event EventClass) (int, EventClass) {
	ctyp := C.int(typ)
	cevent := C.XEventClass(event)
	C.device_key_release((*C.XDevice)(device), &ctyp, &cevent)
	return int(ctyp), EventClass(cevent)
}

func DeviceButtonPress(device *Device, typ int, event EventClass) (int, EventClass) {
	ctyp := C.int(typ)
	cevent := C.XEventClass(event)
	C.device_button_press((*C.XDevice)(device), &ctyp, &cevent)
	return int(ctyp), EventClass(cevent)
}

func DeviceButtonRelease(device *Device, typ int, event EventClass) (int, EventClass) {
	ctyp := C.int(typ)
	cevent := C.XEventClass(event)
	C.device_button_release((*C.XDevice)(device), &ctyp, &cevent)
	return int(ctyp), EventClass(cevent)
}

func DeviceMotionNotify(device *Device, typ int, event EventClass) (int, EventClass) {
	ctyp := C.int(typ)
	cevent := C.XEventClass(event)
	C.device_motion_notify((*C.XDevice)(device), &ctyp, &cevent)
	return int(ctyp), EventClass(cevent)
}

func SelectExtensionEvent(display *xlib.Display, window xlib.Window, event_list []EventClass) error {
	if len(event_list) > 0 {
		return xerr(C.XSelectExtensionEvent((*C.Display)(display), C.Window(window), (*C.XEventClass)(unsafe.Pointer(&event_list[0])), C.int(len(event_list))))
	}
	return xerr(C.XSelectExtensionEvent((*C.Display)(display), C.Window(window), nil, 0))
}

func (d *Device) Classes() []InputClassInfo {
	return (*[math.MaxInt32]InputClassInfo)(unsafe.Pointer(d.classes))[:d.num_classes:d.num_classes]
}

func (d *DeviceInfo) Use() int {
	return int(d.use)
}

func (d *DeviceInfo) ID() xlib.ID {
	return xlib.ID(d.id)
}

func (d *DeviceInfo) Name() string {
	return C.GoString(d.name)
}

func (c *InputClassInfo) InputClass() int {
	return int(c.input_class)
}

func (d *DeviceEvent) EvType() int {
	return int(d.evtype)
}

func xerr(code C.int) error {
	if code == 0 {
		return nil
	}
	return xlib.Status(code)
}
