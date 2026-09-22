package sdl

/*
#include "gosdl.h"
*/
import "C"

type (
	AudioFormat   C.SDL_AudioFormat
	AudioDeviceID C.SDL_AudioDeviceID
)

const (
	AUDIO_U8    = C.SDL_AUDIO_U8
	AUDIO_S8    = C.SDL_AUDIO_S8
	AUDIO_S16LE = C.SDL_AUDIO_S16LE
	AUDIO_S16BE = C.SDL_AUDIO_S16BE
	AUDIO_S16   = C.SDL_AUDIO_S16
)

const (
	AUDIO_S32LE = C.SDL_AUDIO_S32LE
	AUDIO_S32BE = C.SDL_AUDIO_S32BE
	AUDIO_S32   = C.SDL_AUDIO_S32
)

const (
	AUDIO_F32LE = C.SDL_AUDIO_F32LE
	AUDIO_F32BE = C.SDL_AUDIO_F32BE
	AUDIO_F32   = C.SDL_AUDIO_F32
)

type AudioCallback func(stream []uint8)

type AudioSpec struct {
	Format   AudioFormat
	Channels int
	Freq     int
}

func GetCurrentAudioDriver() string {
	return C.GoString(C.SDL_GetCurrentAudioDriver())
}
