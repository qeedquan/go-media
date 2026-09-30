package sdl

/*
#include "gosdl.h"
*/
import "C"

import (
	"sync"
	"unsafe"
)

type (
	AudioFormat   C.SDL_AudioFormat
	AudioDeviceID C.SDL_AudioDeviceID
	AudioStatus   C.SDL_AudioStatus
)

const (
	AUDIO_STOPPED AudioStatus = C.SDL_AUDIO_STOPPED
	AUDIO_PLAYING AudioStatus = C.SDL_AUDIO_PLAYING
	AUDIO_PAUSED  AudioStatus = C.SDL_AUDIO_PAUSED
)

const (
	AUDIO_MASK_BITSIZE  = C.SDL_AUDIO_MASK_BITSIZE
	AUDIO_MASK_DATATYPE = C.SDL_AUDIO_MASK_DATATYPE
	AUDIO_MASK_ENDIAN   = C.SDL_AUDIO_MASK_ENDIAN
	AUDIO_MASK_SIGNED   = C.SDL_AUDIO_MASK_SIGNED
)

const (
	AUDIO_U8     = C.AUDIO_U8
	AUDIO_S8     = C.AUDIO_S8
	AUDIO_U16LSB = C.AUDIO_U16LSB
	AUDIO_S16LSB = C.AUDIO_S16LSB
	AUDIO_U16MSB = C.AUDIO_U16MSB
	AUDIO_S16MSB = C.AUDIO_S16MSB
	AUDIO_U16    = C.AUDIO_U16
	AUDIO_S16    = C.AUDIO_S16
)

const (
	AUDIO_S32LSB = C.AUDIO_S32LSB
	AUDIO_S32MSB = C.AUDIO_S32MSB
	AUDIO_S32    = C.AUDIO_S32
)

const (
	AUDIO_F32LSB = C.AUDIO_F32LSB
	AUDIO_F32MSB = C.AUDIO_F32MSB
	AUDIO_F32    = C.AUDIO_F32
)

const (
	AUDIO_U16SYS = C.AUDIO_U16SYS
	AUDIO_S16SYS = C.AUDIO_S16SYS
	AUDIO_S32SYS = C.AUDIO_S32SYS
	AUDIO_F32SYS = C.AUDIO_F32SYS
)

const (
	AUDIO_ALLOW_FREQUENCY_CHANGE = C.SDL_AUDIO_ALLOW_FREQUENCY_CHANGE
	AUDIO_ALLOW_FORMAT_CHANGE    = C.SDL_AUDIO_ALLOW_FORMAT_CHANGE
	AUDIO_ALLOW_CHANNELS_CHANGE  = C.SDL_AUDIO_ALLOW_CHANNELS_CHANGE
	AUDIO_ALLOW_ANY_CHANGE       = C.SDL_AUDIO_ALLOW_ANY_CHANGE
)

type AudioCallback func(stream []uint8)

type AudioSpec struct {
	Freq     int
	Format   AudioFormat
	Channels uint8
	Silence  uint8
	Samples  uint16
	Size     uint32
	Callback AudioCallback
}

type AudioFilter func(cvt *AudioCVT, format AudioFormat)

type AudioCVT struct {
	Needed      bool
	SrcFormat   AudioFormat
	DstFormat   AudioFormat
	RateIncr    float64
	Buf         []uint8
	Length      int
	LenCVT      int
	LenMult     int
	LenRatio    float64
	Filters     [10]AudioFilter
	FilterIndex int
}

func GetNumAudioDrivers() int {
	return int(C.SDL_GetNumAudioDrivers())
}

func GetAudioDriver(index int) string {
	return C.GoString(C.SDL_GetAudioDriver(C.int(index)))
}

func GetCurrentAudioDriver() string {
	return C.GoString(C.SDL_GetCurrentAudioDriver())
}

func cAudioSpec(a *AudioSpec) C.SDL_AudioSpec {
	spec := C.SDL_AudioSpec{
		freq:     C.int(a.Freq),
		format:   C.SDL_AudioFormat(a.Format),
		channels: C.Uint8(a.Channels),
		silence:  C.Uint8(a.Silence),
		samples:  C.Uint16(a.Samples),
		size:     C.Uint32(a.Size),
	}
	if a.Callback != nil {
		C.setAudioCallback(&spec, C.int(xaudio.AddCallback(a.Callback)))
	}
	return spec
}

func goAudioSpec(a *C.SDL_AudioSpec) AudioSpec {
	return AudioSpec{
		Freq:     int(a.freq),
		Format:   AudioFormat(a.format),
		Channels: uint8(a.channels),
		Silence:  uint8(a.silence),
		Samples:  uint16(a.samples),
		Size:     uint32(a.size),
	}
}

func OpenAudio(desired AudioSpec) (AudioSpec, error) {
	var spec, obtained C.SDL_AudioSpec
	spec = cAudioSpec(&desired)
	C.SDL_OpenAudio(&spec, &obtained)
	return goAudioSpec(&obtained), nil
}

func OpenAudioDevice(device string, iscapture int, desired AudioSpec, allowed_changes int) (dev AudioDeviceID, obtained AudioSpec, err error) {
	var spec, ospec C.SDL_AudioSpec
	var cdevice *C.char

	spec = cAudioSpec(&desired)
	if device != "" {
		cdevice = C.CString(device)
		defer C.free(unsafe.Pointer(cdevice))
	}

	dev = AudioDeviceID(C.SDL_OpenAudioDevice(cdevice, C.int(iscapture), &spec, &ospec, C.int(allowed_changes)))
	obtained = goAudioSpec(&ospec)
	if dev == 0 {
		err = GetError()
	}
	return
}

func (dev AudioDeviceID) GetQueuedAudioSize() int {
	return int(C.SDL_GetQueuedAudioSize(C.SDL_AudioDeviceID(dev)))
}

func (dev AudioDeviceID) ClearQueuedAudio() {
	C.SDL_ClearQueuedAudio(C.SDL_AudioDeviceID(dev))
}

func (dev AudioDeviceID) QueueAudio(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	return ek(C.SDL_QueueAudio(C.SDL_AudioDeviceID(dev), unsafe.Pointer(&data[0]), C.Uint32(len(data))))
}

func (dev AudioDeviceID) DequeueAudio(data []byte) uint32 {
	if len(data) == 0 {
		return 0
	}
	return uint32(C.SDL_DequeueAudio(C.SDL_AudioDeviceID(dev), unsafe.Pointer(&data[0]), C.Uint32(len(data))))
}

func (dev AudioDeviceID) Lock() {
	C.SDL_LockAudioDevice(C.SDL_AudioDeviceID(dev))
}

func (dev AudioDeviceID) Unlock() {
	C.SDL_UnlockAudioDevice(C.SDL_AudioDeviceID(dev))
}

func LockAudio() {
	C.SDL_LockAudio()
}

func LockAudioDevice(dev AudioDeviceID) {
	C.SDL_LockAudioDevice(C.SDL_AudioDeviceID(dev))
}

func UnlockAudio() {
	C.SDL_UnlockAudio()
}

func UnlockAudioDevice(dev AudioDeviceID) {
	C.SDL_UnlockAudioDevice(C.SDL_AudioDeviceID(dev))
}

func CloseAudio() {
	C.SDL_CloseAudio()
}

func CloseAudioDevice(dev AudioDeviceID) {
	C.SDL_CloseAudioDevice(C.SDL_AudioDeviceID(dev))
}

func GetQueueAudioSize(dev AudioDeviceID) uint32 {
	return uint32(C.SDL_GetQueuedAudioSize(C.SDL_AudioDeviceID(dev)))
}

func PauseAudio(pauseOn int) {
	C.SDL_PauseAudio(C.int(pauseOn))
}

func (dev AudioDeviceID) PauseAudio(pauseOn int) {
	C.SDL_PauseAudioDevice(C.SDL_AudioDeviceID(dev), C.int(pauseOn))
}

type audio struct {
	sync.Mutex
	callbacks []AudioCallback
}

var (
	xaudio audio
)

func (a *audio) AddCallback(f AudioCallback) int {
	a.Lock()
	defer a.Unlock()
	a.callbacks = append(a.callbacks, f)
	return len(a.callbacks) - 1
}

//export goAudioCallback
func goAudioCallback(p unsafe.Pointer, stream *C.Uint8, len C.int) {
	buf := unsafe.Slice((*byte)(stream), len)
	xaudio.callbacks[uintptr(p)](buf)
}
