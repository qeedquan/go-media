#include "gosdl.h"
#include "_cgo_export.h"

void setAudioCallback(SDL_AudioSpec *spec, int handle) {
	spec->callback = goAudioCallback;
	spec->userdata = (void*)(uintptr_t)handle;
}
