package sdlmixer

/*
#include <SDL.h>
#include <SDL_mixer.h>
#include "mix.h"
*/
import "C"

import (
	"errors"
	"fmt"
	"os"
	"unsafe"
)

const (
	CHANNELS = C.MIX_CHANNELS

	DEFAULT_FREQUENCY = C.MIX_DEFAULT_FREQUENCY
	DEFAULT_FORMAT    = C.MIX_DEFAULT_FORMAT
	DEFAULT_CHANNELS  = C.MIX_DEFAULT_CHANNELS
	MAX_VOLUME        = C.MIX_MAX_VOLUME
)

type InitFlags C.MIX_InitFlags
type Fading C.Mix_Fading
type MusicType C.Mix_MusicType
type Chunk C.Mix_Chunk
type Music C.Mix_Music

const (
	INIT_FLAC InitFlags = C.MIX_INIT_FLAC
	INIT_MOD  InitFlags = C.MIX_INIT_MOD
	INIT_MP3  InitFlags = C.MIX_INIT_MP3
	INIT_OGG  InitFlags = C.MIX_INIT_OGG
)

const (
	NO_FADING  Fading = C.MIX_NO_FADING
	FADING_OUT Fading = C.MIX_FADING_OUT
	FADING_IN  Fading = C.MIX_FADING_IN
)

const (
	MUS_NONE MusicType = C.MUS_NONE
	MUS_CMD  MusicType = C.MUS_CMD
	MUS_WAV  MusicType = C.MUS_WAV
	MUS_MOD  MusicType = C.MUS_MOD
	MUS_MID  MusicType = C.MUS_MID
	MUS_OGG  MusicType = C.MUS_OGG
	MUS_MP3  MusicType = C.MUS_MP3
	MUS_FLAC MusicType = C.MUS_FLAC
)

func SetError(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	cs := C.CString(s)
	C.setError(cs)
	C.free(unsafe.Pointer(cs))
}

func GetError() error {
	return errors.New(C.GoString(C.Mix_GetError()))
}

func Init(flags InitFlags) (InitFlags, error) {
	initted := InitFlags(C.Mix_Init(C.int(flags)))
	if initted&flags != flags {
		return initted, GetError()
	}
	return initted, nil
}

func Quit() {
	C.Mix_Quit()
}

func Pause(channel int) {
	C.Mix_Pause(C.int(channel))
}

func PausedMusic() bool {
	return C.Mix_PausedMusic() == 1
}

func Resume(channel int) {
	C.Mix_Resume(C.int(channel))
}

func ResumeMusic() {
	C.Mix_ResumeMusic()
}

func Volume(channel, volume int) int {
	return int(C.Mix_Volume(C.int(channel), C.int(volume)))
}

func VolumeMusic(volume int) int {
	return int(C.Mix_VolumeMusic(C.int(volume)))
}

func HaltChannel(channel int) {
	C.Mix_HaltChannel(C.int(channel))
}

func ExpireChannel(channel, ticks int) int {
	return int(C.Mix_ExpireChannel(C.int(channel), C.int(ticks)))
}

func FadeOutChannel(channel, ms int) int {
	return int(C.Mix_FadeOutChannel(C.int(channel), C.int(ms)))
}

func Playing(channel int) int {
	return int(C.Mix_Playing(C.int(channel)))
}

func PlayingMusic() int {
	return int(C.Mix_PlayingMusic())
}

func Paused(channel int) int {
	return int(C.Mix_Paused(C.int(channel)))
}

func FadingChannel(which int) Fading {
	return Fading(C.Mix_FadingChannel(C.int(which)))
}

func HaltMusic() {
	C.Mix_HaltMusic()
}

func OpenAudio(freq int, format uint16, channels, chunksize int) error {
	err := C.Mix_OpenAudio(C.int(freq), C.Uint16(format), C.int(channels), C.int(chunksize))
	if err < 0 {
		return GetError()
	}
	return nil
}

func CloseAudio() {
	C.Mix_CloseAudio()
}

func AllocateChannels(num int) int {
	return int(C.Mix_AllocateChannels(C.int(num)))
}

func QuerySpec() (freq int, format uint16, channels int, err error) {
	var cfreq, cchannels C.int
	var cformat C.Uint16

	rc := C.Mix_QuerySpec(&cfreq, &cformat, &cchannels)
	freq, format, channels = int(cfreq), uint16(cformat), int(cchannels)
	if rc == 0 {
		err = GetError()
		return
	}
	return
}

func NumChunkDecoders() int         { return int(C.Mix_GetNumChunkDecoders()) }
func ChunkDecoder(index int) string { return C.GoString(C.Mix_GetChunkDecoder(C.int(index))) }

func NumMusicDecoders() int         { return int(C.Mix_GetNumMusicDecoders()) }
func MusicDecoder(index int) string { return C.GoString(C.Mix_GetMusicDecoder(C.int(index))) }

func LoadWAV(name string) (*Chunk, error) {
	cs := C.CString(name)
	defer C.free(unsafe.Pointer(cs))

	c := C.loadWAV(cs)
	if c == nil {
		return nil, &os.PathError{
			Op:   "open",
			Path: name,
			Err:  GetError(),
		}
	}
	return (*Chunk)(c), nil
}

func LoadMUS(name string) (*Music, error) {
	cs := C.CString(name)
	defer C.free(unsafe.Pointer(cs))

	mus := C.Mix_LoadMUS(cs)
	if mus == nil {
		return nil, &os.PathError{
			Op:   "open",
			Path: name,
			Err:  GetError(),
		}
	}
	return (*Music)(mus), nil
}

func QuickLoadWAV(mem []byte) (*Chunk, error) {
	c := C.Mix_QuickLoad_WAV((*C.Uint8)(unsafe.Pointer(&mem[0])))
	if c == nil {
		return nil, GetError()
	}
	return (*Chunk)(c), nil
}

func QuickLoadRAW(mem []byte) (*Chunk, error) {
	c := C.Mix_QuickLoad_RAW((*C.Uint8)(unsafe.Pointer(&mem[0])), C.Uint32(len(mem)))
	if c == nil {
		return nil, GetError()
	}
	return (*Chunk)(c), nil
}

func (c *Chunk) Free() {
	C.Mix_FreeChunk((*C.Mix_Chunk)(c))
}

func FadeInMusic(m *Music, loops, ms int) error {
	err := C.Mix_FadeInMusic((*C.Mix_Music)(m), C.int(loops), C.int(ms))
	if err != 0 {
		return GetError()
	}
	return nil
}

func FadeOutMusic(ms int) error {
	err := C.Mix_FadeOutMusic(C.int(ms))
	if err == 0 {
		return GetError()
	}
	return nil
}

func (m *Music) Play(loops int) error {
	err := C.Mix_PlayMusic((*C.Mix_Music)(m), C.int(loops))
	if err != 0 {
		return GetError()
	}
	return nil
}

func (m *Music) Type() MusicType {
	return MusicType(C.Mix_GetMusicType((*C.Mix_Music)(m)))
}

func (m *Music) Free() {
	C.Mix_FreeMusic((*C.Mix_Music)(m))
}

func (c *Chunk) PlayChannel(channel, loops int) int {
	return c.PlayChannelTimed(channel, loops, -1)
}

func (c *Chunk) PlayChannelTimed(channel, loops, ticks int) int {
	return int(C.Mix_PlayChannelTimed(C.int(channel), (*C.Mix_Chunk)(c), C.int(loops), C.int(ticks)))
}

func (c *Chunk) FadeInChannel(channel, loops, ms int) int {
	return int(C.Mix_FadeInChannelTimed(C.int(channel), (*C.Mix_Chunk)(c), C.int(loops), C.int(ms), -1))
}

func (c *Chunk) FadeInChannelTimed(channel, loops, ms, ticks int) int {
	return int(C.Mix_FadeInChannelTimed(C.int(channel), (*C.Mix_Chunk)(c), C.int(loops), C.int(ms), C.int(ticks)))
}

func (c *Chunk) Volume(volume int) int {
	return int(C.Mix_VolumeChunk((*C.Mix_Chunk)(c), C.int(volume)))
}

func PauseMusic() {
	C.Mix_PauseMusic()
}

func ReserveChannels(num int) int {
	return int(C.Mix_ReserveChannels(C.int(num)))
}

var (
	hookMusicFn func()
)

func HookMusicFinished(f func()) {
	hookMusicFn = f
	C.hookMusicFinished()
}

//export goMusicFinished
func goMusicFinished() {
	if hookMusicFn != nil {
		hookMusicFn()
	}
}
