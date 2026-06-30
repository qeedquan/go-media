#include "SDL_ttf.h"

extern DECLSPEC int SDLCALL TTF_SizeUTF8Ex(TTF_Font *font, const char *text, size_t textlen, int *w, int *h);
extern DECLSPEC SDL_Surface *SDLCALL TTF_RenderUTF8_BlendedEx(TTF_Font *font, SDL_Surface *textbuf, SDL_Rect *bounds,
                                                              const char *text, size_t textlen, SDL_Color fg);
extern DECLSPEC SDL_Surface *SDLCALL TTF_RenderUTF8_SolidEx(TTF_Font *font, SDL_Surface *textbuf, SDL_Rect *bounds,
                                                            const char *text, size_t textlen, SDL_Color fg);