#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <X11/Xlib.h>
#include <X11/Xutil.h>
#include <X11/extensions/scrnsaver.h>
#include <X11/Xatom.h>
#include <X11/Xresource.h>
#include <X11/cursorfont.h>
#include <X11/Xregion.h>

int ev_type(XEvent *ev);
XKeyEvent *ev_xkey(XEvent *ev);
XSelectionEvent *ev_xselection(XEvent *ev);
XSelectionRequestEvent *ev_xselectionrequest(XEvent *ev);
XSelectionClearEvent *ev_xselectionclear(XEvent *ev);
XCreateWindowEvent *ev_xcreatewindow(XEvent *ev);
XGenericEventCookie *ev_xcookie(XEvent *ev);
XVisibilityEvent *ev_xvisibility(XEvent *ev);
XClientMessageEvent *ev_xclient(XEvent *ev);
XButtonEvent *ev_xbutton(XEvent *ev);
XConfigureEvent *ev_xconfigure(XEvent *ev);
XFocusChangeEvent *ev_xfocus(XEvent *ev);
XPropertyEvent *ev_xproperty(XEvent *ev);
void ev_xclient_long(XClientMessageEvent *ev, long *l, int n);
Bool xregister_im_instantitate_callback(Display *display, XrmDatabase db, char *res_name, char *res_class, unsigned long id);
void xim_destroy_callback(XIM xim, XPointer client, XPointer call);
char *xset_im_values_void(XIM xim, const char *key, void *ptr);
char *xset_ic_values_void(XIC im, const char *key, void *ptr);
XIC xcreateic(XIM xim, long input_style, Window client_window, Window focus_window);
void xset_error_handler(void);
XVaNestedList va_create_nested_list1(const char *k1, void *v1);