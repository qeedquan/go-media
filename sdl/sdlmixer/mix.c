#include <SDL.h>
#include <SDL_mixer.h>
#include "mix.h"
#include "_cgo_export.h"

void setError(char *s) {
    Mix_SetError("%s", s);
}

Mix_Chunk *loadWAV(char *file) {
    return Mix_LoadWAV(file);
}

void hookMusicFinished(void) {
    Mix_HookMusicFinished(goMusicFinished);
}
