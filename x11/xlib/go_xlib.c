#include "go_xlib.h"
#include "_cgo_export.h"

int ev_type(XEvent *ev) {
	return ev->type;
}

XKeyEvent *ev_xkey(XEvent *ev) {
	return &ev->xkey;
}

XSelectionEvent *ev_xselection(XEvent *ev) {
	return &ev->xselection;
}

XSelectionClearEvent *ev_xselectionclear(XEvent *ev) {
	return &ev->xselectionclear;
}

XSelectionRequestEvent *ev_xselectionrequest(XEvent *ev) {
	return &ev->xselectionrequest;
}

XCreateWindowEvent *ev_xcreatewindow(XEvent *ev) {
	return &ev->xcreatewindow;
}

XGenericEventCookie *ev_xcookie(XEvent *ev) {
	return &ev->xcookie;
}

XVisibilityEvent *ev_xvisibility(XEvent *ev) {
	return &ev->xvisibility;
}

XClientMessageEvent *ev_xclient(XEvent *ev) {
	return &ev->xclient;
}

XButtonEvent *ev_xbutton(XEvent *ev) {
	return &ev->xbutton;
}

XConfigureEvent *ev_xconfigure(XEvent *ev) {
	return &ev->xconfigure;
}

XFocusChangeEvent *ev_xfocus(XEvent *ev) {
	return &ev->xfocus;
}

XPropertyEvent *ev_xproperty(XEvent *ev) {
	return &ev->xproperty;
}

void ev_xclient_long(XClientMessageEvent *ev, long *l, int n) {
	memcpy(l, ev->data.l, sizeof(*l) * n);
}

void xset_error_handler(void) {
	XSetErrorHandler(goErrorHandler);
}

Bool xregister_im_instantitate_callback(Display *display, XrmDatabase db, char *res_name, char *res_class, unsigned long id) {
	return XRegisterIMInstantiateCallback(display, db, res_name, res_class, goIMInstantitateCallback, (char*)id);
}

char *xset_im_values_void(XIM im, const char *key, void *ptr) {
	return XSetIMValues(im, key, ptr, NULL);
}

char *xset_ic_values_void(XIC ic, const char *key, void *ptr) {
	return XSetICValues(ic, key, ptr, NULL);
}

XIC xcreateic(XIM im, long input_style, Window client_window, Window focus_window) {
	return XCreateIC(im, XNInputStyle, input_style, XNClientWindow, client_window, XNFocusWindow, focus_window, NULL);
}

void xim_destroy_callback(XIM im, XPointer client, XPointer call) {
	goIMDestroyCallback(im, call);
}

XVaNestedList va_create_nested_list1(const char *k1, void *v1) {
	XVaCreateNestedList(0, k1, v1, NULL);
}