package chroma

import (
	"fmt"
	"image/color"
)

func ParseRGBA(s string) (color.RGBA, error) {
	var r, g, b, a uint8
	n, _ := fmt.Sscanf(s, "rgb(%v,%v,%v)", &r, &g, &b)
	if n == 3 {
		return color.RGBA{r, g, b, 255}, nil
	}

	n, _ = fmt.Sscanf(s, "rgba(%v,%v,%v,%v)", &r, &g, &b, &a)
	if n == 4 {
		return color.RGBA{r, g, b, a}, nil
	}

	n, _ = fmt.Sscanf(s, "#%02x%02x%02x%02x", &r, &g, &b, &a)
	if n == 4 {
		return color.RGBA{r, g, b, a}, nil
	}

	n, _ = fmt.Sscanf(s, "#%02x%02x%02x", &r, &g, &b)
	if n == 3 {
		return color.RGBA{r, g, b, 255}, nil
	}

	n, _ = fmt.Sscanf(s, "#%02x", &r)
	if n == 1 {
		return color.RGBA{r, r, r, 255}, nil
	}

	var h HSV
	n, _ = fmt.Sscanf(s, "hsv(%v,%v,%v)", &h.H, &h.S, &h.V)
	if n == 3 {
		return HSVToRGBA(h), nil
	}

	return color.RGBA{}, fmt.Errorf("failed to parse color %q, unknown format", s)
}
