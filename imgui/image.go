package imgui

import (
	"image/color"
	"math"

	"github.com/qeedquan/go-media/math/f64"
)

func (c *Context) Image(user_texture_id TextureID, size, uv0, uv1 f64.Vec2, tint_col, border_col color.RGBA) {
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return
	}

	bb := f64.Rectangle{
		window.DC.CursorPos,
		window.DC.CursorPos.Add(size),
	}
	if border_col.A > 0 {
		bb.Max = bb.Max.Add(f64.Vec2{2, 2})
	}
	c.ItemSizeBB(bb)
	if !c.ItemAdd(bb, 0) {
		return
	}

	if border_col.A > 0 {
		window.DrawList.AddRect(bb.Min, bb.Max, border_col)
		window.DrawList.AddImageEx(user_texture_id, bb.Min.Add(f64.Vec2{1, 1}), bb.Max.Sub(f64.Vec2{1, 1}), uv0, uv1, tint_col)
	} else {
		window.DrawList.AddImageEx(user_texture_id, bb.Min, bb.Max, uv0, uv1, tint_col)
	}
}

func (c *Context) ImageButton(user_texture_id TextureID, size f64.Vec2) bool {
	return c.ImageButtonEx(user_texture_id, size, f64.Vec2{0, 0}, f64.Vec2{1, 1}, -1, color.RGBA{0, 0, 0, 0}, color.RGBA{255, 255, 255, 255})
}

// frame_padding < 0: uses FramePadding from style (default)
// frame_padding = 0: no framing
// frame_padding > 0: set framing size
// The color used are the button colors.
func (c *Context) ImageButtonEx(user_texture_id TextureID, size, uv0, uv1 f64.Vec2, frame_padding int, bg_col, tint_col color.RGBA) bool {
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	style := &c.Style

	// Default to using texture ID as ID. User can still push string/integer prefixes.
	// We could hash the size/uv to create a unique ID but that would prevent the user from animating UV.
	c.PushID(ID(user_texture_id.(uint32)))
	id := window.GetID("#image")
	c.PopID()

	padding := style.FramePadding
	if frame_padding >= 0 {
		padding = f64.Vec2{float64(frame_padding), float64(frame_padding)}
	}
	bb := f64.Rectangle{
		window.DC.CursorPos,
		window.DC.CursorPos.Add(size).Add(padding.Scale(2)),
	}
	image_bb := f64.Rectangle{
		window.DC.CursorPos.Add(padding),
		window.DC.CursorPos.Add(padding).Add(size),
	}
	c.ItemSizeBB(bb)
	if !c.ItemAdd(bb, id) {
		return false
	}

	hovered, held, pressed := c.ButtonBehavior(bb, id, 0)

	// Render
	var col color.RGBA
	switch {
	case hovered && held:
		col = c.GetColorFromStyle(ColButtonActive)
	case hovered:
		col = c.GetColorFromStyle(ColButtonHovered)
	default:
		col = c.GetColorFromStyle(ColButton)
	}
	c.RenderNavHighlight(bb, id)
	c.RenderFrameEx(bb.Min, bb.Max, col, true, f64.Clamp(math.Min(padding.X, padding.Y), 0.0, style.FrameRounding))
	if bg_col.A > 0 {
		window.DrawList.AddRectFilled(image_bb.Min, image_bb.Max, bg_col)
	}
	window.DrawList.AddImageEx(user_texture_id, image_bb.Min, image_bb.Max, uv0, uv1, tint_col)

	return pressed
}