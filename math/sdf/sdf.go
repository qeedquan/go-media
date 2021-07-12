// https://www.iquilezles.org/www/articles/distfunctions2d/distfunctions2d.htm

package sdf

import "github.com/qeedquan/go-media/math/f64"

func Circle(p f64.Vec2, r float64) float64 {
	return p.Len() - r
}
