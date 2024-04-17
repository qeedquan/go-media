package readline

/*
#include <stdio.h>
#include <stdlib.h>
#include "go_readline.h"

*/
// #cgo LDFLAGS: -lreadline
import "C"

import (
	"sync"
	"unsafe"
)

type handles struct {
	sync.Mutex
	key map[int]func(int, int) int
}

func newHandles() *handles {
	return &handles{
		key: make(map[int]func(int, int) int),
	}
}

var (
	hands = newHandles()
)

//export goKeyCallback
func goKeyCallback(arg0, invokingKey C.int) C.int {
	hands.Lock()
	defer hands.Unlock()
	f := hands.key[int(invokingKey)]
	return C.int(f(int(arg0), int(invokingKey)))
}

func ResetTerminal(terminal string) int {
	cterminal := C.CString(terminal)
	defer C.free(unsafe.Pointer(cterminal))
	return int(C.rl_reset_terminal(cterminal))
}

func BindKey(key int, fun func(int, int) int) {
	hands.Lock()
	defer hands.Unlock()
	hands.key[key] = fun
	C.bindKeyCallback(C.int(key))
}

func Read(prompt string) (input string, eof bool) {
	cprompt := C.CString(prompt)
	defer C.free(unsafe.Pointer(cprompt))
	cinput := C.readline(cprompt)
	defer C.free(unsafe.Pointer(cinput))
	if cinput == nil {
		return "", true
	}
	return C.GoString(cinput), false
}

func AddHistory(input string) {
	cinput := C.CString(input)
	defer C.free(unsafe.Pointer(cinput))
	C.add_history(cinput)
}

func Complete(ignore, invokingKey int) int {
	return int(C.rl_complete(C.int(ignore), C.int(invokingKey)))
}
