package imgui

import (
	"image/color"
	"math"

	"github.com/qeedquan/go-media/math/f64"
)

func (c *Context) Checkbox(label string, v *bool) bool {
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	style := c.Style
	id := window.GetID(label)
	label_size := c.CalcTextSizeEx(label, true, -1)

	// We want a square shape to we use Y twice
	check_bb := f64.Rectangle{
		window.DC.CursorPos,
		window.DC.CursorPos.Add(f64.Vec2{
			label_size.Y + style.FramePadding.Y*2,
			label_size.Y + style.FramePadding.Y*2,
		}),
	}
	c.ItemSizeBBEx(check_bb, style.FramePadding.Y)

	total_bb := check_bb
	if label_size.X > 0 {
		c.SameLineEx(0, style.ItemInnerSpacing.X)
	}
	text_bb := f64.Rectangle{
		window.DC.CursorPos.Add(f64.Vec2{0, style.FramePadding.Y}),
		window.DC.CursorPos.Add(f64.Vec2{0, style.FramePadding.Y}.Add(label_size)),
	}
	if label_size.X > 0 {
		c.ItemSizeEx(
			f64.Vec2{text_bb.Dx(), check_bb.Dy()},
			style.FramePadding.Y,
		)

		total_bb = f64.Rectangle{
			check_bb.Min.Min(text_bb.Min),
			check_bb.Max.Max(text_bb.Max),
		}
	}

	if !c.ItemAdd(total_bb, id) {
		return false
	}

	hovered, held, pressed := c.ButtonBehavior(total_bb, id, 0)
	if pressed {
		*v = !*v
	}

	var col color.RGBA
	switch {
	case held && hovered:
		col = c.GetColorFromStyle(ColFrameBgActive)
	case hovered:
		col = c.GetColorFromStyle(ColFrameBgHovered)
	default:
		col = c.GetColorFromStyle(ColFrameBg)
	}

	c.RenderNavHighlight(total_bb, id)
	c.RenderFrameEx(check_bb.Min, check_bb.Max, col, true, style.FrameRounding)

	if *v {
		check_sz := math.Min(check_bb.Dx(), check_bb.Dy())
		pad := math.Max(1.0, float64(int(check_sz/6.0)))
		c.RenderCheckMark(
			check_bb.Min.Add(f64.Vec2{pad, pad}),
			c.GetColorFromStyle(ColCheckMark),
			check_bb.Dx()-pad*2.0,
		)
	}

	if c.LogEnabled {
		if *v {
			c.LogRenderedText(&text_bb.Min, "[x]")
		} else {
			c.LogRenderedText(&text_bb.Min, "[ ]")
		}
	}
	if label_size.X > 0 {
		c.RenderText(text_bb.Min, label)
	}

	return pressed
}

func (c *Context) CheckboxFlags(label string, flags *uint, flags_value uint) bool {
	v := *flags&flags_value == flags_value
	pressed := c.Checkbox(label, &v)
	if pressed {
		if v {
			*flags |= flags_value
		} else {
			*flags &^= flags_value
		}
	}
	return pressed
}