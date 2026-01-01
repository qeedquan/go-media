#ifndef _GFX_H_
#define _GFX_H_

#include "SDL2_gfxPrimitives.h"

int goPolygonRGBA(SDL_Renderer *renderer, const SDL_Point *pts, int n, Uint8 r, Uint8 g, Uint8 b, Uint8 a);
int goAAPolygonRGBA(SDL_Renderer *renderer, const SDL_Point *pts, int n, Uint8 r, Uint8 g, Uint8 b, Uint8 a);
int goFilledPolygonRGBA(SDL_Renderer *renderer, const SDL_Point *pts, int n, Uint8 r, Uint8 g, Uint8 b, Uint8 a);
int goTexturedPolygon(SDL_Renderer *renderer, const SDL_Point *pts, int n, SDL_Surface *texture, int texture_dx, int texture_dy);
int goBezierRGBA(SDL_Renderer *renderer, const SDL_Point *pts, int n, int s, Uint8 r, Uint8 g, Uint8 b, Uint8 a);
void goGfxPrimitivesSetFont(const void *fontdata, Uint32 cw, Uint32 ch);
void goGfxPrimitivesSetFontRotation(Uint32 rotation);

int goCircle(SDL_Renderer *renderer, int xm, int ym, int rad, Uint8 r, Uint8 g, Uint8 b, Uint8 a);
int goFilledCircle(SDL_Renderer *renderer, int xm, int ym, int rad, Uint8 r, Uint8 g, Uint8 b, Uint8 a);

extern Uint32 goCharWidth;
extern Uint32 goCharHeight;
extern Uint32 goCharRotation;

#endif
