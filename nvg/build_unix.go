// +build linux darwin openbsd netbsd freebsd dragonfly solaris

package nvg

// #cgo LDFLAGS: -lm -lGL -ldl
// #cgo pkg-config: glew
import "C"
