package stbperlin

/*

#define STB_PERLIN_IMPLEMENTATION
#include "stb_perlin.h"

*/
import "C"

import (
	"github.com/qeedquan/go-media/math/ga"
)

func Noise3[T ga.Signed](x, y, z T, x_wrap, y_wrap, z_wrap int) T {
	return T(C.stb_perlin_noise3(C.float(x), C.float(y), C.float(z), C.int(x_wrap), C.int(y_wrap), C.int(z_wrap)))
}

func Noise3s[T ga.Signed](x, y, z T, x_wrap, y_wrap, z_wrap, seed int) T {
	return T(C.stb_perlin_noise3_seed(C.float(x), C.float(y), C.float(z), C.int(x_wrap), C.int(y_wrap), C.int(z_wrap), C.int(seed)))
}

func Ridge3[T ga.Signed](x, y, z, lacunarity, gain, offset T, octaves int) T {
	return T(C.stb_perlin_ridge_noise3(C.float(x), C.float(y), C.float(z), C.float(lacunarity), C.float(gain), C.float(offset), C.int(octaves)))
}

func FBM3[T ga.Signed](x, y, z, lacunarity, gain T, octaves int) T {
	return T(C.stb_perlin_fbm_noise3(C.float(x), C.float(y), C.float(z), C.float(lacunarity), C.float(gain), C.int(octaves)))
}

func Turbulence3[T ga.Signed](x, y, z, lacunarity, gain, offset T, octaves int) T {
	return T(C.stb_perlin_turbulence_noise3(C.float(x), C.float(y), C.float(z), C.float(lacunarity), C.float(gain), C.int(octaves)))
}
