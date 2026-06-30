#include "gfx.h"

static int *gfxPrimitivesPolyIntsGlobal;
static int gfxPrimitivesPolyAllocatedGlobal;

int _aalineRGBA(SDL_Renderer *renderer, Sint16 x1, Sint16 y1, Sint16 x2, Sint16 y2, Uint8 r, Uint8 g, Uint8 b, Uint8 a, int draw_endpoint);
int _gfxPrimitivesCompareInt(const void *a, const void *b);
int hline(SDL_Renderer *renderer, Sint16 x1, Sint16 x2, Sint16 y);
int _HLineTextured(SDL_Renderer *renderer, Sint16 x1, Sint16 x2, Sint16 y, SDL_Texture *texture, int texture_w, int texture_h, int texture_dx, int texture_dy);
double _evaluateBezier(double *data, int ndata, double t);
int line(SDL_Renderer *renderer, Sint16 x1, Sint16 y1, Sint16 x2, Sint16 y2);

int
goPolygon(SDL_Renderer *renderer, const SDL_Point *pts, int n)
{
	int result = 0;
	/*
	* Draw 
	*/
	result |= SDL_RenderDrawLines(renderer, pts, n);
	result |= SDL_RenderDrawLine(renderer, pts[n - 1].x, pts[n - 1].y, pts[0].x, pts[0].y);

	return (result);
}

int
goPolygonRGBA(SDL_Renderer *renderer, const SDL_Point *pts, int n, Uint8 r, Uint8 g, Uint8 b, Uint8 a)
{
	int result;

	/*
	* Set color 
	*/
	result = 0;
	result |= SDL_SetRenderDrawBlendMode(renderer, (a == 255) ? SDL_BLENDMODE_NONE : SDL_BLENDMODE_BLEND);
	result |= SDL_SetRenderDrawColor(renderer, r, g, b, a);

	/*
	* Draw 
	*/
	result |= goPolygon(renderer, pts, n);

	return (result);
}

int
goAAPolygonRGBA(SDL_Renderer *renderer, const SDL_Point *pts, int n, Uint8 r, Uint8 g, Uint8 b, Uint8 a)
{
	int i, result;

	/*
	* Draw 
	*/
	result = 0;
	for (i = 1; i < n; i++)
		result |= _aalineRGBA(renderer, pts[i - 1].x, pts[i - 1].y, pts[i].x, pts[i].y, r, g, b, a, 0);

	result |= _aalineRGBA(renderer, pts[n - 1].x, pts[n - 1].y, pts[0].x, pts[0].y, r, g, b, a, 0);

	return (result);
}

int
goFilledPolygonRGBAMT(SDL_Renderer *renderer, const SDL_Point *pts, int n, Uint8 r, Uint8 g, Uint8 b, Uint8 a, int **polyInts, int *polyAllocated)
{
	int result;
	int i;
	int y, xa, xb;
	int miny, maxy;
	int x1, y1;
	int x2, y2;
	int ind1, ind2;
	int ints;
	int *gfxPrimitivesPolyInts = NULL;
	int *gfxPrimitivesPolyIntsNew = NULL;
	int gfxPrimitivesPolyAllocated = 0;

	/*
	* Map polygon cache  
	*/
	if ((polyInts == NULL) || (polyAllocated == NULL)) {
		/* Use global cache */
		gfxPrimitivesPolyInts = gfxPrimitivesPolyIntsGlobal;
		gfxPrimitivesPolyAllocated = gfxPrimitivesPolyAllocatedGlobal;
	} else {
		/* Use local cache */
		gfxPrimitivesPolyInts = *polyInts;
		gfxPrimitivesPolyAllocated = *polyAllocated;
	}

	/*
	* Allocate temp array, only grow array 
	*/
	if (!gfxPrimitivesPolyAllocated) {
		gfxPrimitivesPolyInts = (int *)malloc(sizeof(int) * n);
		gfxPrimitivesPolyAllocated = n;
	} else {
		if (gfxPrimitivesPolyAllocated < n) {
			gfxPrimitivesPolyIntsNew = (int *)realloc(gfxPrimitivesPolyInts, sizeof(int) * n);
			if (!gfxPrimitivesPolyIntsNew) {
				if (!gfxPrimitivesPolyInts) {
					free(gfxPrimitivesPolyInts);
					gfxPrimitivesPolyInts = NULL;
				}
				gfxPrimitivesPolyAllocated = 0;
			} else {
				gfxPrimitivesPolyInts = gfxPrimitivesPolyIntsNew;
				gfxPrimitivesPolyAllocated = n;
			}
		}
	}

	/*
	* Check temp array
	*/
	if (gfxPrimitivesPolyInts == NULL) {
		gfxPrimitivesPolyAllocated = 0;
	}

	/*
	* Update cache variables
	*/
	if ((polyInts == NULL) || (polyAllocated == NULL)) {
		gfxPrimitivesPolyIntsGlobal = gfxPrimitivesPolyInts;
		gfxPrimitivesPolyAllocatedGlobal = gfxPrimitivesPolyAllocated;
	} else {
		*polyInts = gfxPrimitivesPolyInts;
		*polyAllocated = gfxPrimitivesPolyAllocated;
	}

	/*
	* Check temp array again
	*/
	if (gfxPrimitivesPolyInts == NULL) {
		return (-1);
	}

	/*
	* Determine Y maxima 
	*/
	miny = pts[0].y;
	maxy = pts[0].y;
	for (i = 1; (i < n); i++) {
		if (pts[i].y < miny) {
			miny = pts[i].y;
		} else if (pts[i].y > maxy) {
			maxy = pts[i].y;
		}
	}

	/*
	* Draw, scanning y 
	*/
	result = 0;
	for (y = miny; (y <= maxy); y++) {
		ints = 0;
		for (i = 0; (i < n); i++) {
			if (!i) {
				ind1 = n - 1;
				ind2 = 0;
			} else {
				ind1 = i - 1;
				ind2 = i;
			}
			y1 = pts[ind1].y;
			y2 = pts[ind2].y;
			if (y1 < y2) {
				x1 = pts[ind1].x;
				x2 = pts[ind2].x;
			} else if (y1 > y2) {
				y2 = pts[ind1].y;
				y1 = pts[ind2].y;
				x2 = pts[ind1].x;
				x1 = pts[ind2].x;
			} else {
				continue;
			}
			if (((y >= y1) && (y < y2)) || ((y == maxy) && (y > y1) && (y <= y2))) {
				gfxPrimitivesPolyInts[ints++] = ((65536 * (y - y1)) / (y2 - y1)) * (x2 - x1) + (65536 * x1);
			}
		}

		qsort(gfxPrimitivesPolyInts, ints, sizeof(int), _gfxPrimitivesCompareInt);

		/*
		* Set color 
		*/
		result = 0;
		result |= SDL_SetRenderDrawBlendMode(renderer, (a == 255) ? SDL_BLENDMODE_NONE : SDL_BLENDMODE_BLEND);
		result |= SDL_SetRenderDrawColor(renderer, r, g, b, a);

		for (i = 0; (i < ints); i += 2) {
			xa = gfxPrimitivesPolyInts[i] + 1;
			xa = (xa >> 16) + ((xa & 32768) >> 15);
			xb = gfxPrimitivesPolyInts[i + 1] - 1;
			xb = (xb >> 16) + ((xb & 32768) >> 15);
			result |= hline(renderer, xa, xb, y);
		}
	}

	return (result);
}

int
goFilledPolygonRGBA(SDL_Renderer *renderer, const SDL_Point *pts, int n, Uint8 r, Uint8 g, Uint8 b, Uint8 a)
{
	return goFilledPolygonRGBAMT(renderer, pts, n, r, g, b, a, NULL, NULL);
}

int
goTexturedPolygonMT(SDL_Renderer *renderer, const SDL_Point *pts, int n,
                    SDL_Surface *texture, int texture_dx, int texture_dy, int **polyInts, int *polyAllocated)
{
	int result;
	int i;
	int y, xa, xb;
	int minx, maxx, miny, maxy;
	int x1, y1;
	int x2, y2;
	int ind1, ind2;
	int ints;
	int *gfxPrimitivesPolyInts = NULL;
	int *gfxPrimitivesPolyIntsTemp = NULL;
	int gfxPrimitivesPolyAllocated = 0;
	SDL_Texture *textureAsTexture = NULL;

	/*
	* Map polygon cache  
	*/
	if ((polyInts == NULL) || (polyAllocated == NULL)) {
		/* Use global cache */
		gfxPrimitivesPolyInts = gfxPrimitivesPolyIntsGlobal;
		gfxPrimitivesPolyAllocated = gfxPrimitivesPolyAllocatedGlobal;
	} else {
		/* Use local cache */
		gfxPrimitivesPolyInts = *polyInts;
		gfxPrimitivesPolyAllocated = *polyAllocated;
	}

	/*
	* Allocate temp array, only grow array 
	*/
	if (!gfxPrimitivesPolyAllocated) {
		gfxPrimitivesPolyInts = (int *)malloc(sizeof(int) * n);
		gfxPrimitivesPolyAllocated = n;
	} else {
		if (gfxPrimitivesPolyAllocated < n) {
			gfxPrimitivesPolyIntsTemp = (int *)realloc(gfxPrimitivesPolyInts, sizeof(int) * n);
			if (gfxPrimitivesPolyIntsTemp == NULL) {
				/* Realloc failed - keeps original memory block, but fails this operation */
				return (-1);
			}
			gfxPrimitivesPolyInts = gfxPrimitivesPolyIntsTemp;
			gfxPrimitivesPolyAllocated = n;
		}
	}

	/*
	* Check temp array
	*/
	if (gfxPrimitivesPolyInts == NULL) {
		gfxPrimitivesPolyAllocated = 0;
	}

	/*
	* Update cache variables
	*/
	if ((polyInts == NULL) || (polyAllocated == NULL)) {
		gfxPrimitivesPolyIntsGlobal = gfxPrimitivesPolyInts;
		gfxPrimitivesPolyAllocatedGlobal = gfxPrimitivesPolyAllocated;
	} else {
		*polyInts = gfxPrimitivesPolyInts;
		*polyAllocated = gfxPrimitivesPolyAllocated;
	}

	/*
	* Check temp array again
	*/
	if (gfxPrimitivesPolyInts == NULL) {
		return (-1);
	}

	/*
	* Determine X,Y minima,maxima 
	*/
	miny = pts[0].y;
	maxy = pts[0].y;
	minx = pts[0].x;
	maxx = pts[0].x;
	for (i = 1; (i < n); i++) {
		if (pts[i].y < miny) {
			miny = pts[i].y;
		} else if (pts[i].y > maxy) {
			maxy = pts[i].y;
		}
		if (pts[i].x < minx) {
			minx = pts[i].x;
		} else if (pts[i].x > maxx) {
			maxx = pts[i].x;
		}
	}

	/* Create texture for drawing */
	textureAsTexture = SDL_CreateTextureFromSurface(renderer, texture);
	if (textureAsTexture == NULL) {
		return -1;
	}
	SDL_SetTextureBlendMode(textureAsTexture, SDL_BLENDMODE_BLEND);

	/*
	* Draw, scanning y 
	*/
	result = 0;
	for (y = miny; (y <= maxy); y++) {
		ints = 0;
		for (i = 0; (i < n); i++) {
			if (!i) {
				ind1 = n - 1;
				ind2 = 0;
			} else {
				ind1 = i - 1;
				ind2 = i;
			}
			y1 = pts[ind1].y;
			y2 = pts[ind2].y;
			if (y1 < y2) {
				x1 = pts[ind1].x;
				x2 = pts[ind2].x;
			} else if (y1 > y2) {
				y2 = pts[ind1].y;
				y1 = pts[ind2].y;
				x2 = pts[ind1].x;
				x1 = pts[ind2].x;
			} else {
				continue;
			}
			if (((y >= y1) && (y < y2)) || ((y == maxy) && (y > y1) && (y <= y2))) {
				gfxPrimitivesPolyInts[ints++] = ((65536 * (y - y1)) / (y2 - y1)) * (x2 - x1) + (65536 * x1);
			}
		}

		qsort(gfxPrimitivesPolyInts, ints, sizeof(int), _gfxPrimitivesCompareInt);

		for (i = 0; (i < ints); i += 2) {
			xa = gfxPrimitivesPolyInts[i] + 1;
			xa = (xa >> 16) + ((xa & 32768) >> 15);
			xb = gfxPrimitivesPolyInts[i + 1] - 1;
			xb = (xb >> 16) + ((xb & 32768) >> 15);
			result |= _HLineTextured(renderer, xa, xb, y, textureAsTexture, texture->w, texture->h, texture_dx, texture_dy);
		}
	}

	SDL_RenderPresent(renderer);
	SDL_DestroyTexture(textureAsTexture);

	return (result);
}

int
goTexturedPolygon(SDL_Renderer *renderer, const SDL_Point *pts, int n, SDL_Surface *texture, int texture_dx, int texture_dy)
{
	/*
	* Draw
	*/
	return (goTexturedPolygonMT(renderer, pts, n, texture, texture_dx, texture_dy, NULL, NULL));
}

int
goBezierRGBA(SDL_Renderer *renderer, const SDL_Point *pts, int n, int s, Uint8 r, Uint8 g, Uint8 b, Uint8 a)
{
	int result;
	int i;
	double *x, *y, t, stepsize;
	Sint16 x1, y1, x2, y2;

	/*
	* Sanity check 
	*/
	if (n < 3) {
		return (-1);
	}
	if (s < 2) {
		return (-1);
	}

	/*
	* Variable setup 
	*/
	stepsize = (double)1.0 / (double)s;

	/* Transfer vertices into float arrays */
	if ((x = (double *)malloc(sizeof(double) * (n + 1))) == NULL) {
		return (-1);
	}
	if ((y = (double *)malloc(sizeof(double) * (n + 1))) == NULL) {
		free(x);
		return (-1);
	}
	for (i = 0; i < n; i++) {
		x[i] = (double)pts[i].x;
		y[i] = (double)pts[i].y;
	}
	x[n] = (double)pts[0].x;
	y[n] = (double)pts[0].y;

	/*
	* Set color 
	*/
	result = 0;
	result |= SDL_SetRenderDrawBlendMode(renderer, (a == 255) ? SDL_BLENDMODE_NONE : SDL_BLENDMODE_BLEND);
	result |= SDL_SetRenderDrawColor(renderer, r, g, b, a);

	/*
	* Draw 
	*/
	t = 0.0;
	x1 = (Sint16)lrint(_evaluateBezier(x, n + 1, t));
	y1 = (Sint16)lrint(_evaluateBezier(y, n + 1, t));
	for (i = 0; i <= (n * s); i++) {
		t += stepsize;
		x2 = (Sint16)_evaluateBezier(x, n, t);
		y2 = (Sint16)_evaluateBezier(y, n, t);
		result |= line(renderer, x1, y1, x2, y2);
		x1 = x2;
		y1 = y2;
	}

	/* Clean up temporary array */
	free(x);
	free(y);

	return (result);
}

Uint32 goCharWidth = 8;
Uint32 goCharHeight = 8;
Uint32 goCharRotation = 0;

void
goGfxPrimitivesSetFont(const void *fontdata, Uint32 cw, Uint32 ch)
{
	goCharWidth = cw;
	goCharHeight = ch;
	gfxPrimitivesSetFont(fontdata, cw, ch);
}

void
goGfxPrimitivesSetFontRotation(Uint32 rotation)
{
	goCharRotation = rotation;
	gfxPrimitivesSetFontRotation(rotation);
}

int
goCircle(SDL_Renderer *renderer, int xm, int ym, int rad, Uint8 r, Uint8 g, Uint8 b, Uint8 a)
{
	if (rad <= 0)
		return 0;

	int x = -rad, y = 0, err = 2 - 2 * rad;
	int ret = SDL_SetRenderDrawColor(renderer, r, g, b, a);
	do {
		ret |= SDL_RenderDrawPoint(renderer, xm - x, ym + y);
		ret |= SDL_RenderDrawPoint(renderer, xm - y, ym - x);
		ret |= SDL_RenderDrawPoint(renderer, xm + x, ym - y);
		ret |= SDL_RenderDrawPoint(renderer, xm + y, ym + x);

		rad = err;
		if (rad <= y)
			err += ++y * 2 + 1;
		if (rad > x || err > y)
			err += ++x * 2 + 1;
	} while (x < 0);
	return ret;
}

int
goFilledCircle(SDL_Renderer *renderer, int xm, int ym, int rad, Uint8 r, Uint8 g, Uint8 b, Uint8 a)
{
	if (rad <= 0)
		return 0;

	int x = -rad, y = 0, err = 2 - 2 * rad;
	int ret = SDL_SetRenderDrawColor(renderer, r, g, b, a);
	do {
		ret |= SDL_RenderDrawLine(renderer, xm - x, ym - y, xm + x, ym - y);
		ret |= SDL_RenderDrawLine(renderer, xm - x, ym + y, xm + x, ym + y);

		ret |= SDL_RenderDrawLine(renderer, xm - y, ym - x, xm + y, ym - x);
		ret |= SDL_RenderDrawLine(renderer, xm - y, ym + x, xm + y, ym + x);

		rad = err;
		if (rad <= y)
			err += ++y * 2 + 1;
		if (rad > x || err > y)
			err += ++x * 2 + 1;
	} while (x < 0);
	return ret;
}
