package imgui

import (
	"fmt"

	"github.com/qeedquan/go-media/math/f64"
)

func (c *Context) ProgressBar(fraction float64) {
	c.ProgressBarEx(fraction, f64.Vec2{-1, 0}, "")
}

// size_arg (for each axis) < 0.0f: align to end, 0.0f: auto, > 0.0f: specified size
func (c *Context) ProgressBarEx(fraction float64, size_arg f64.Vec2, overlay string) {
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return
	}

	style := &c.Style

	pos := window.DC.CursorPos
	bb := f64.Rectangle{
		pos,
		pos.Add(c.CalcItemSize(size_arg, c.CalcItemWidth(), c.FontSize+style.FramePadding.Y*2.0)),
	}
	c.ItemSizeBBEx(bb, style.FramePadding.Y)
	if !c.ItemAdd(bb, 0) {
		return
	}

	// Render
	fraction = f64.Saturate(fraction)
	c.RenderFrameEx(bb.Min, bb.Max, c.GetColorFromStyle(ColFrameBg), true, style.FrameRounding)
	bb = bb.Expand2(f64.Vec2{-style.FrameBorderSize, -style.FrameBorderSize})
	fill_br := f64.Vec2{f64.Lerp(fraction, bb.Min.X, bb.Max.X), bb.Max.Y}
	c.RenderRectFilledRangeH(window.DrawList, bb, c.GetColorFromStyle(ColPlotHistogram), 0.0, fraction, style.FrameRounding)

	// Default displaying the fraction as percentage string, but user can override it
	if overlay != "" {
		overlay = fmt.Sprintf("%.0f%%", fraction*100+0.01)
	}

	overlay_size := c.CalcTextSize(overlay)
	if overlay_size.X > 0.0 {
		c.RenderTextClippedEx(f64.Vec2{f64.Clamp(fill_br.X+style.ItemSpacing.X, bb.Min.X, bb.Max.X-overlay_size.X-style.ItemInnerSpacing.X), bb.Min.Y}, bb.Max, overlay, &overlay_size, f64.Vec2{0.0, 0.5}, &bb)
	}
}
