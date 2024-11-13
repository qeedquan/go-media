package cairo

/*
#include <cairo.h>
*/
import "C"

type Status C.cairo_status_t

const (
	STATUS_SUCCESS Status = C.CAIRO_STATUS_SUCCESS
)

func (s Status) Error() string {
	return C.GoString(C.cairo_status_to_string((C.cairo_status_t)(s)))
}

func xk(n C.cairo_status_t) error {
	if Status(n) == STATUS_SUCCESS {
		return nil
	}
	return Status(n)
}
