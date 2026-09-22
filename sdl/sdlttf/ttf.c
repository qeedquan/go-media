#include "ttf.h"
#include "SDL_ttf.c.inc"

int
TTF_SizeUTF8Ex(TTF_Font *font, const char *text, size_t textlen, int *w, int *h)
{
	int      status;
	int      x, z;
	int      minx, maxx;
	int      miny, maxy;
	c_glyph *glyph;
	FT_Error error;
	FT_Long  use_kerning;
	FT_UInt  prev_index    = 0;
	int      outline_delta = 0;

	TTF_CHECKPOINTER(text, -1);

	/* Initialize everything to 0 */
	status = 0;
	minx = maxx = 0;
	miny = maxy = 0;

	/* check kerning */
	use_kerning = FT_HAS_KERNING(font->face) && font->kerning;

	/* Init outline handling */
	if (font->outline > 0) {
		outline_delta = font->outline * 2;
	}

	/* Load each character and sum it's bounding box */
	x = 0;
	while (textlen > 0) {
		Uint16 c = UTF8_getch(&text, &textlen);
		if (c == UNICODE_BOM_NATIVE || c == UNICODE_BOM_SWAPPED) {
			continue;
		}

		error = Find_Glyph(font, c, CACHED_METRICS);
		if (error) {
			TTF_SetFTError("Couldn't find glyph", error);
			return -1;
		}
		glyph = font->current;

		/* handle kerning */
		if (use_kerning && prev_index && glyph->index) {
			FT_Vector delta;
			FT_Get_Kerning(font->face, prev_index, glyph->index,
			               ft_kerning_default, &delta);
			x += delta.x >> 6;
		}
#if 0
		if ((ch == text) && (glyph->minx < 0)) {
			/* Fixes the texture wrapping bug when the first letter
			 * has a negative minx value or horibearing value.  The entire
			 * bounding box must be adjusted to be bigger so the entire
			 * letter can fit without any texture corruption or wrapping.
			 *
			 * Effects: First enlarges bounding box.
			 * Second, xstart has to start ahead of its normal spot in the
			 * negative direction of the negative minx value.
			 * (pushes everything to the right).
			 *
			 * This will make the memory copy of the glyph bitmap data
			 * work out correctly.
			 * */
			z -= glyph->minx;
		}
#endif

		z = x + glyph->minx;
		if (minx > z) {
			minx = z;
		}
		if (TTF_HANDLE_STYLE_BOLD(font)) {
			x += font->glyph_overhang;
		}
		if (glyph->advance > glyph->maxx) {
			z = x + glyph->advance;
		} else {
			z = x + glyph->maxx;
		}
		if (maxx < z) {
			maxx = z;
		}
		x += glyph->advance;

		if (glyph->miny < miny) {
			miny = glyph->miny;
		}
		if (glyph->maxy > maxy) {
			maxy = glyph->maxy;
		}
		prev_index = glyph->index;
	}

	/* Fill the bounds rectangle */
	if (w) {
		/* Add outline extra width */
		*w = (maxx - minx) + outline_delta;
	}
	if (h) {
		/* Some fonts descend below font height (FletcherGothicFLF) */
		/* Add outline extra height */
		*h = (font->ascent - miny) + outline_delta;
		if (*h < font->height) {
			*h = font->height;
		}
		/* Update height according to the needs of the underline style */
		if (TTF_HANDLE_STYLE_UNDERLINE(font)) {
			int bottom_row = TTF_underline_bottom_row(font);
			if (*h < bottom_row) {
				*h = bottom_row;
			}
		}
	}
	return status;
}

SDL_Surface *
TTF_RenderUTF8_BlendedEx(TTF_Font *font, SDL_Surface *textbuf,
                         SDL_Rect *bounds, const char *text, size_t textlen, SDL_Color fg)
{
	SDL_bool first;
	int      xstart;
	int      width, height;
	Uint32   alpha;
	Uint32   pixel;
	Uint8 *  src;
	Uint32 * dst;
	Uint32 * dst_check;
	int      row, col;
	c_glyph *glyph;
	FT_Error error;
	FT_Long  use_kerning;
	FT_UInt  prev_index = 0;

	TTF_CHECKPOINTER(text, NULL);

	if (textbuf->format->BytesPerPixel != 4) {
		TTF_SetError("Pixel buffer is not 32bpp");
		return NULL;
	}

	/* Get the dimensions of the text surface */
	if ((TTF_SizeUTF8(font, text, &width, &height) < 0) || !width) {
		TTF_SetError("Text has zero width");
		return (NULL);
	}
	bounds->x = 0;
	bounds->y = 0;
	bounds->w = width;
	bounds->h = height;

	/* Adding bound checking to avoid all kinds of memory corruption errors
	 * that may occur. */
	dst_check =
	    (Uint32 *)textbuf->pixels + textbuf->pitch / 4 * textbuf->h;

	/* check kerning */
	use_kerning = FT_HAS_KERNING(font->face) && font->kerning;

	/* Load and render each character */
	first  = SDL_TRUE;
	xstart = 0;
	pixel  = (fg.r << 16) | (fg.g << 8) | fg.b;
	SDL_FillRect(textbuf, NULL, pixel); /* Initialize with fg and 0 alpha */
	while (textlen > 0) {
		Uint16 c = UTF8_getch(&text, &textlen);
		if (c == UNICODE_BOM_NATIVE || c == UNICODE_BOM_SWAPPED) {
			continue;
		}

		error = Find_Glyph(font, c, CACHED_METRICS | CACHED_PIXMAP);
		if (error) {
			TTF_SetFTError("Couldn't find glyph", error);
			SDL_FreeSurface(textbuf);
			return NULL;
		}
		glyph = font->current;
		/* Ensure the width of the pixmap is correct. On some cases,
		 * freetype may report a larger pixmap than possible.*/
		width = glyph->pixmap.width;
		if (font->outline <= 0 && width > glyph->maxx - glyph->minx) {
			width = glyph->maxx - glyph->minx;
		}
		/* do kerning, if possible AC-Patch */
		if (use_kerning && prev_index && glyph->index) {
			FT_Vector delta;
			FT_Get_Kerning(font->face, prev_index, glyph->index,
			               ft_kerning_default, &delta);
			xstart += delta.x >> 6;
		}

		/* Compensate for the wrap around bug with negative minx's */
		if (first && (glyph->minx < 0)) {
			xstart -= glyph->minx;
		}
		first = SDL_FALSE;

		for (row = 0; row < glyph->pixmap.rows; ++row) {
			/* Make sure we don't go either over, or under the
			 * limit */
			if (row + glyph->yoffset < 0) {
				continue;
			}
			if (row + glyph->yoffset >= textbuf->h) {
				continue;
			}
			dst = (Uint32 *)textbuf->pixels +
			      (row + glyph->yoffset) * textbuf->pitch / 4 +
			      xstart + glyph->minx;

			/* Added code to adjust src pointer for pixmaps to
			 * account for pitch.
			 * */
			src =
			    (Uint8 *)(glyph->pixmap.buffer +
			              glyph->pixmap.pitch * row);
			for (col = width; col > 0 && dst < dst_check; --col) {
				alpha = *src++;
				*dst++ |= pixel | (alpha << 24);
			}
		}

		xstart += glyph->advance;
		if (TTF_HANDLE_STYLE_BOLD(font)) {
			xstart += font->glyph_overhang;
		}
		prev_index = glyph->index;
	}

	/* Handle the underline style */
	if (TTF_HANDLE_STYLE_UNDERLINE(font)) {
		row = TTF_underline_top_row(font);
		TTF_drawLine_Blended(font, textbuf, row, pixel);
	}

	/* Handle the strikethrough style */
	if (TTF_HANDLE_STYLE_STRIKETHROUGH(font)) {
		row = TTF_strikethrough_top_row(font);
		TTF_drawLine_Blended(font, textbuf, row, pixel);
	}
	return (textbuf);
}

SDL_Surface *
TTF_RenderUTF8_SolidEx(TTF_Font *font, SDL_Surface *textbuf, SDL_Rect *bounds,
                       const char *text, size_t textlen, SDL_Color fg)
{
	SDL_bool     first;
	int          xstart;
	int          width;
	int          height;
	SDL_Palette *palette;
	Uint8 *      src;
	Uint8 *      dst;
	Uint8 *      dst_check;
	int          row, col;
	c_glyph *    glyph;

	FT_Bitmap *current;
	FT_Error   error;
	FT_Long    use_kerning;
	FT_UInt    prev_index = 0;

	TTF_CHECKPOINTER(text, NULL);

	/* Get the dimensions of the text surface */
	if ((TTF_SizeUTF8(font, text, &width, &height) < 0) || !width) {
		TTF_SetError("Text has zero width");
		return NULL;
	}
	bounds->x = 0;
	bounds->y = 0;
	bounds->w = width;
	bounds->h = height;

	/* Adding bound checking to avoid all kinds of memory corruption errors
	 * that may occur. */
	dst_check = (Uint8 *)textbuf->pixels + textbuf->pitch * textbuf->h;

	/* Fill the palette with the foreground color */
	palette = textbuf->format->palette;
	if (!palette) {
		TTF_SetError("Surface has no palette");
		return NULL;
	}
	palette->colors[0].r = 255 - fg.r;
	palette->colors[0].g = 255 - fg.g;
	palette->colors[0].b = 255 - fg.b;
	palette->colors[1].r = fg.r;
	palette->colors[1].g = fg.g;
	palette->colors[1].b = fg.b;
	SDL_SetColorKey(textbuf, SDL_TRUE, 0);

	/* check kerning */
	use_kerning = FT_HAS_KERNING(font->face) && font->kerning;

	/* Load and render each character */
	first  = SDL_TRUE;
	xstart = 0;
	while (textlen > 0) {
		Uint16 c = UTF8_getch(&text, &textlen);
		if (c == UNICODE_BOM_NATIVE || c == UNICODE_BOM_SWAPPED) {
			continue;
		}

		error = Find_Glyph(font, c, CACHED_METRICS | CACHED_BITMAP);
		if (error) {
			TTF_SetFTError("Couldn't find glyph", error);
			SDL_FreeSurface(textbuf);
			return NULL;
		}
		glyph   = font->current;
		current = &glyph->bitmap;
		/* Ensure the width of the pixmap is correct. On some cases,
		 * freetype may report a larger pixmap than possible.*/
		width = current->width;
		if (font->outline <= 0 && width > glyph->maxx - glyph->minx) {
			width = glyph->maxx - glyph->minx;
		}
		/* do kerning, if possible AC-Patch */
		if (use_kerning && prev_index && glyph->index) {
			FT_Vector delta;
			FT_Get_Kerning(font->face, prev_index, glyph->index,
			               ft_kerning_default, &delta);
			xstart += delta.x >> 6;
		}
		/* Compensate for wrap around bug with negative minx's */
		if (first && (glyph->minx < 0)) {
			xstart -= glyph->minx;
		}
		first = SDL_FALSE;

		for (row = 0; row < current->rows; ++row) {
			/* Make sure we don't go either over, or under the
			 * limit */
			if (row + glyph->yoffset < 0) {
				continue;
			}
			if (row + glyph->yoffset >= textbuf->h) {
				continue;
			}
			dst = (Uint8 *)textbuf->pixels +
			      (row + glyph->yoffset) * textbuf->pitch + xstart +
			      glyph->minx;
			src = current->buffer + row * current->pitch;

			for (col = width; col > 0 && dst < dst_check; --col) {
				*dst++ |= *src++;
			}
		}

		xstart += glyph->advance;
		if (TTF_HANDLE_STYLE_BOLD(font)) {
			xstart += font->glyph_overhang;
		}
		prev_index = glyph->index;
	}

	/* Handle the underline style */
	if (TTF_HANDLE_STYLE_UNDERLINE(font)) {
		row = TTF_underline_top_row(font);
		TTF_drawLine_Solid(font, textbuf, row);
	}

	/* Handle the strikethrough style */
	if (TTF_HANDLE_STYLE_STRIKETHROUGH(font)) {
		row = TTF_strikethrough_top_row(font);
		TTF_drawLine_Solid(font, textbuf, row);
	}
	return textbuf;
}
