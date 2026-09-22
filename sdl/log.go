package sdl

/*
#include "gosdl.h"

static void xlog(const char *str) {
	SDL_Log("%s", str);
}

static void xlogVerbose(int category, const char *str) {
	SDL_LogVerbose(category, "%s", str);
}

static void xlogDebug(int category, const char *str) {
	SDL_LogDebug(category, "%s", str);
}

static void xlogInfo(int category, const char *str) {
	SDL_LogInfo(category, "%s", str);
}

static void xlogWarn(int category, const char *str) {
	SDL_LogWarn(category, "%s", str);
}

static void xlogError(int category, const char *str) {
	SDL_LogError(category, "%s", str);
}

static void xlogCritical(int category, const char *str) {
	SDL_LogCritical(category, "%s", str);
}

static void xlogMessage(int category, SDL_LogPriority priority, const char *str) {
	SDL_LogMessage(category, priority, "%s", str);
}
*/
import "C"
import (
	"fmt"
	"unsafe"
)

const (
	LOG_CATEGORY_APPLICATION = C.SDL_LOG_CATEGORY_APPLICATION
	LOG_CATEGORY_ERROR       = C.SDL_LOG_CATEGORY_ERROR
	LOG_CATEGORY_ASSERT      = C.SDL_LOG_CATEGORY_ASSERT
	LOG_CATEGORY_SYSTEM      = C.SDL_LOG_CATEGORY_SYSTEM
	LOG_CATEGORY_AUDIO       = C.SDL_LOG_CATEGORY_AUDIO
	LOG_CATEGORY_VIDEO       = C.SDL_LOG_CATEGORY_VIDEO
	LOG_CATEGORY_RENDER      = C.SDL_LOG_CATEGORY_RENDER
	LOG_CATEGORY_INPUT       = C.SDL_LOG_CATEGORY_INPUT
	LOG_CATEGORY_TEST        = C.SDL_LOG_CATEGORY_TEST

	LOG_CATEGORY_CUSTOM = C.SDL_LOG_CATEGORY_CUSTOM
)

type LogPriority C.SDL_LogPriority

const (
	LOG_PRIORITY_VERBOSE  LogPriority = C.SDL_LOG_PRIORITY_VERBOSE
	LOG_PRIORITY_DEBUG    LogPriority = C.SDL_LOG_PRIORITY_DEBUG
	LOG_PRIORITY_INFO     LogPriority = C.SDL_LOG_PRIORITY_INFO
	LOG_PRIORITY_WARN     LogPriority = C.SDL_LOG_PRIORITY_WARN
	LOG_PRIORITY_ERROR    LogPriority = C.SDL_LOG_PRIORITY_ERROR
	LOG_PRIORITY_CRITICAL LogPriority = C.SDL_LOG_PRIORITY_CRITICAL
	NUM_LOG_PRIORITIES                = C.SDL_NUM_LOG_PRIORITIES
)

func LogSetAllPriority(prio LogPriority) {
	C.SDL_LogSetAllPriority(C.SDL_LogPriority(prio))
}

func LogSetPriority(category int, prio LogPriority) {
	C.SDL_LogSetPriority(C.int(category), C.SDL_LogPriority(prio))
}

func LogResetPriorities() {
	C.SDL_LogResetPriorities()
}

func Log(format string, args ...interface{}) {
	cstr := C.CString(fmt.Sprintf(format, args...))
	defer C.free(unsafe.Pointer(cstr))
	C.xlog(cstr)
}

func LogVerbose(category int, format string, args ...interface{}) {
	cstr := C.CString(fmt.Sprintf(format, args...))
	defer C.free(unsafe.Pointer(cstr))
	C.xlogVerbose(C.int(category), cstr)
}

func LogDebug(category int, format string, args ...interface{}) {
	cstr := C.CString(fmt.Sprintf(format, args...))
	defer C.free(unsafe.Pointer(cstr))
	C.xlogDebug(C.int(category), cstr)
}

func LogInfo(category int, format string, args ...interface{}) {
	cstr := C.CString(fmt.Sprintf(format, args...))
	defer C.free(unsafe.Pointer(cstr))
	C.xlogInfo(C.int(category), cstr)
}

func LogWarn(category int, format string, args ...interface{}) {
	cstr := C.CString(fmt.Sprintf(format, args...))
	defer C.free(unsafe.Pointer(cstr))
	C.xlogWarn(C.int(category), cstr)
}

func LogError(category int, format string, args ...interface{}) {
	cstr := C.CString(fmt.Sprintf(format, args...))
	defer C.free(unsafe.Pointer(cstr))
	C.xlogError(C.int(category), cstr)
}

func LogCritical(category int, format string, args ...interface{}) {
	cstr := C.CString(fmt.Sprintf(format, args...))
	defer C.free(unsafe.Pointer(cstr))
	C.xlogCritical(C.int(category), cstr)
}

func LogMessage(category int, prio LogPriority, format string, args ...interface{}) {
	cstr := C.CString(fmt.Sprintf(format, args...))
	defer C.free(unsafe.Pointer(cstr))
	C.xlogMessage(C.int(category), C.SDL_LogPriority(prio), cstr)
}
