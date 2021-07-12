package imgui

import (
	"image/color"
	"math"

	"github.com/qeedquan/go-media/math/f64"
)

type SelectableFlags int

const (
	SelectableFlagsDontClosePopups  SelectableFlags = 1 << 0 // Clicking this don't close parent popup window
	SelectableFlagsSpanAllColumns   SelectableFlags = 1 << 1 // Selectable frame can span all columns (text will still fit in current column)
	SelectableFlagsAllowDoubleClick SelectableFlags = 1 << 2 // Generate press events on double clicks too

	// NB: need to be in sync with last value of ImGuiSelectableFlags_
	SelectableFlagsMenu               SelectableFlags = 1 << 3 // -> PressedOnClick
	SelectableFlagsMenuItem           SelectableFlags = 1 << 4 // -> PressedOnRelease
	SelectableFlagsDisabled           SelectableFlags = 1 << 5
	SelectableFlagsDrawFillAvailWidth SelectableFlags = 1 << 6
)

func (c *Context) Selectable(label string) bool {
	return c.SelectableEx(label, false, 0, f64.Vec2{0, 0})
}

// Tip: pass an empty label (e.g. "##dummy") then you can use the space to draw other text or image.
// But you need to make sure the ID is unique, e.g. enclose calls in PushID/PopID.
func (c *Context) SelectableEx(label string, selected bool, flags SelectableFlags, size_arg f64.Vec2) bool {
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	style := &c.Style

	// FIXME-OPT: Avoid if vertically clipped.
	if flags&SelectableFlagsSpanAllColumns != 0 && window.DC.ColumnsSet != nil {
		c.PopClipRect()
	}

	id := window.GetID(label)
	label_size := c.CalcTextSizeEx(label, true, -1)
	size := label_size
	if size_arg.X != 0 {
		size.X = size_arg.X
	}
	if size_arg.Y != 0 {
		size.Y = size_arg.Y
	}

	pos := window.DC.CursorPos
	pos.Y += window.DC.CurrentLineTextBaseOffset
	bb := f64.Rectangle{pos, pos.Add(size)}
	c.ItemSizeBB(bb)

	// Fill horizontal space.
	window_padding := window.WindowPadding
	var max_x float64
	if flags&SelectableFlagsSpanAllColumns != 0 {
		max_x = c.GetWindowContentRegionMax().X
	} else {
		max_x = c.GetContentRegionMax().X
	}
	w_draw := math.Max(label_size.X, window.Pos.X+max_x-window_padding.X-window.DC.CursorPos.X)
	size_draw := f64.Vec2{w_draw, size.Y}
	if size_arg.X != 0 && flags&SelectableFlagsDrawFillAvailWidth == 0 {
		size_draw.X = size_arg.X
	}
	if size_arg.Y != 0 {
		size_draw.Y = size_arg.Y
	}
	bb_with_spacing := f64.Rectangle{pos, pos.Add(size_draw)}
	if size_arg.X == 0.0 || flags&SelectableFlagsDrawFillAvailWidth != 0 {
		bb_with_spacing.Max.X += window_padding.X
	}

	// Selectables are tightly packed together, we extend the box to cover spacing between selectable.
	spacing_L := float64((int)(style.ItemSpacing.X * 0.5))
	spacing_U := float64((int)(style.ItemSpacing.Y * 0.5))
	spacing_R := style.ItemSpacing.X - spacing_L
	spacing_D := style.ItemSpacing.Y - spacing_U
	bb_with_spacing.Min.X -= spacing_L
	bb_with_spacing.Min.Y -= spacing_U
	bb_with_spacing.Max.X += spacing_R
	bb_with_spacing.Max.Y += spacing_D

	select_id := id
	if flags&SelectableFlagsDisabled != 0 {
		select_id = 0
	}
	if !c.ItemAdd(bb_with_spacing, select_id) {
		if flags&SelectableFlagsSpanAllColumns != 0 && window.DC.ColumnsSet != nil {
			c.PushColumnClipRect()
		}
		return false
	}

	var button_flags ButtonFlags
	if flags&SelectableFlagsMenu != 0 {
		button_flags |= ButtonFlagsPressedOnClick | ButtonFlagsNoHoldingActiveID
	}
	if flags&SelectableFlagsMenuItem != 0 {
		button_flags |= ButtonFlagsPressedOnRelease
	}
	if flags&SelectableFlagsDisabled != 0 {
		button_flags |= ButtonFlagsDisabled
	}
	if flags&SelectableFlagsAllowDoubleClick != 0 {
		button_flags |= ButtonFlagsPressedOnClickRelease | ButtonFlagsPressedOnDoubleClick
	}
	hovered, held, pressed := c.ButtonBehavior(bb_with_spacing, id, button_flags)
	if flags&SelectableFlagsDisabled != 0 {
		selected = false
	}

	// Hovering selectable with mouse updates NavId accordingly so navigation can be resumed with gamepad/keyboard (this doesn't happen on most widgets)
	if pressed || hovered {
		if !c.NavDisableMouseHover && c.NavWindow == window && c.NavLayer == window.DC.NavLayerCurrent {
			c.NavDisableHighlight = true
			c.SetNavID(id, window.DC.NavLayerCurrent)
		}
	}

	// Render
	if hovered || selected {
		var col color.RGBA
		switch {
		case held && hovered:
			col = c.GetColorFromStyle(ColHeaderActive)
		case hovered:
			col = c.GetColorFromStyle(ColHeaderHovered)
		default:
			col = c.GetColorFromStyle(ColHeader)
		}
		c.RenderFrameEx(bb_with_spacing.Min, bb_with_spacing.Max, col, false, 0.0)
		c.RenderNavHighlightEx(bb_with_spacing, id, NavHighlightFlagsTypeThin|NavHighlightFlagsNoRounding)
	}

	if flags&SelectableFlagsSpanAllColumns != 0 && window.DC.ColumnsSet != nil {
		c.PushColumnClipRect()
		bb_with_spacing.Max.X -= (c.GetContentRegionMax().X - max_x)
	}

	if flags&SelectableFlagsDisabled != 0 {
		c.PushStyleColorV4(ColText, c.Style.Colors[ColTextDisabled])
	}
	c.RenderTextClippedEx(bb.Min, bb_with_spacing.Max, label, &label_size, f64.Vec2{0.0, 0.0}, nil)
	if flags&SelectableFlagsDisabled != 0 {
		c.PopStyleColor()
	}

	// Automatically close popups
	if pressed && (window.Flags&WindowFlagsPopup) != 0 && flags&SelectableFlagsDontClosePopups == 0 && window.DC.ItemFlags&ItemFlagsSelectableDontClosePopup == 0 {
		c.CloseCurrentPopup()
	}
	return pressed
}

func (c *Context) SelectableOpen(label string, p_selected *bool) bool {
	return c.SelectableOpenEx(label, p_selected, 0, f64.Vec2{0, 0})
}

func (c *Context) SelectableOpenEx(label string, p_selected *bool, flags SelectableFlags, size f64.Vec2) bool {
	if c.SelectableEx(label, *p_selected, flags, size) {
		*p_selected = !*p_selected
		return true
	}
	return false
}