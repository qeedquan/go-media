package imgui

import (
	"fmt"
	"image/color"
	"math"

	"github.com/qeedquan/go-media/math/f64"
	"github.com/qeedquan/go-media/math/mathutil"
)

type DragDropFlags int

const (
	// BeginDragDropSource() flags
	DragDropFlagsSourceNoPreviewTooltip   DragDropFlags = 1 << 0 // By default a successful call to BeginDragDropSource opens a tooltip so you can display a preview or description of the source contents. This flag disable this behavior.
	DragDropFlagsSourceNoDisableHover     DragDropFlags = 1 << 1 // By default when dragging we clear data so that IsItemHovered() will return true to avoid subsequent user code submitting tooltips. This flag disable this behavior so you can still call IsItemHovered() on the source item.
	DragDropFlagsSourceNoHoldToOpenOthers DragDropFlags = 1 << 2 // Disable the behavior that allows to open tree nodes and collapsing header by holding over them while dragging a source item.
	DragDropFlagsSourceAllowNullID        DragDropFlags = 1 << 3 // Allow items such as Text() Image() that have no unique identifier to be used as drag source by manufacturing a temporary identifier based on their window-relative position. This is extremely unusual within the dear imgui ecosystem and so we made it explicit.
	DragDropFlagsSourceExtern             DragDropFlags = 1 << 4 // External source (from outside of imgui) won't attempt to read current item/window info. Will always return true. Only one Extern source can be active simultaneously.
	// AcceptDragDropPayload() flags
	DragDropFlagsAcceptBeforeDelivery    DragDropFlags = 1 << 10                                                                  // AcceptDragDropPayload() will returns true even before the mouse button is released. You can then call IsDelivery() to test if the payload needs to be delivered.
	DragDropFlagsAcceptNoDrawDefaultRect DragDropFlags = 1 << 11                                                                  // Do not draw the default highlight rectangle when hovering over target.
	DragDropFlagsAcceptPeekOnly          DragDropFlags = DragDropFlagsAcceptBeforeDelivery | DragDropFlagsAcceptNoDrawDefaultRect // For peeking ahead and inspecting the payload before delivery.
)

type Payload struct {
	Data           interface{}
	SourceId       ID   // Source item id
	SourceParentId ID   // Source parent id (if available)
	DataFrameCount int  // Data timestamp
	Preview        bool // Set when AcceptDragDropPayload() was called and mouse has been hovering the target item (nb: handle overlapping drag targets)
	Delivery       bool // Set when AcceptDragDropPayload() was called and mouse button is released over the target item.
}

// We don't use BeginDragDropTargetCustom() and duplicate its code because:
// 1) we use LastItemRectHoveredRect which handles items that pushes a temporarily clip rectangle in their code. Calling BeginDragDropTargetCustom(LastItemRect) would not handle them.
// 2) and it's faster. as this code may be very frequently called, we want to early out as fast as we can.
// Also note how the HoveredWindow test is positioned differently in both functions (in both functions we optimize for the cheapest early out case)
func (c *Context) BeginDragDropTarget() bool {
	if !c.DragDropActive {
		return false
	}

	window := c.CurrentWindow
	if window.DC.LastItemStatusFlags&ItemStatusFlagsHoveredRect == 0 {
		return false
	}
	if c.HoveredWindow == nil || window.RootWindow != c.HoveredWindow.RootWindow {
		return false
	}

	display_rect := window.DC.LastItemRect
	if window.DC.LastItemStatusFlags&ItemStatusFlagsHasDisplayRect != 0 {
		display_rect = window.DC.LastItemDisplayRect
	}
	id := window.DC.LastItemId
	if id == 0 {
		id = window.GetIDFromRectangle(display_rect)
	}
	if c.DragDropPayload.SourceId == id {
		return false
	}

	c.DragDropTargetRect = display_rect
	c.DragDropTargetId = id
	return true
}

func (c *Context) BeginDragDropTargetCustom(bb f64.Rectangle, id ID) bool {
	if !c.DragDropActive {
		return false
	}

	window := c.CurrentWindow
	if c.HoveredWindow == nil || window.RootWindow != c.HoveredWindow.RootWindow {
		return false
	}
	assert(id != 0)
	if !c.IsMouseHoveringRect(bb.Min, bb.Max) || (id == c.DragDropPayload.SourceId) {
		return false
	}

	c.DragDropTargetRect = bb
	c.DragDropTargetId = id
	return true
}

// We don't really use/need this now, but added it for the sake of consistency and because we might need it later.
func (c *Context) EndDragDropTarget() {
	assert(c.DragDropActive)
}

func (c *Context) BeginDragDropSource() bool {
	return false
}

func (c *Context) EndDragDropSource() {
	assert(c.DragDropActive)
	if c.DragDropSourceFlags&DragDropFlagsSourceNoPreviewTooltip == 0 {
		c.EndTooltip()
		c.PopStyleColor()
		// c.PopStyleVar()
	}

	// Discard the drag if have not called SetDragDropPayload()
	if c.DragDropPayload.DataFrameCount == -1 {
		c.ClearDragDrop()
	}
}

func (c *Context) SetDragDropPayload(a interface{}) {
}

func (c *Context) DragInt(label string, v *int) bool {
	return c.DragIntEx(label, v, 1, 0, 0, "%.0f")
}

// NB: v_speed is float to allow adjusting the drag speed with more precision
func (c *Context) DragIntEx(label string, v *int, v_speed float64, v_min, v_max int, format string) bool {
	if format == "" {
		format = "%.0f"
	}
	v_f := float64(*v)
	value_changed := c.DragFloatEx(label, &v_f, v_speed, float64(v_min), float64(v_max), format, 1)
	*v = int(v_f)
	return value_changed
}

func (c *Context) DragIntN(label string, v []int, v_speed float64, v_min, v_max int, format string) bool {
	components := len(v)
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	value_changed := false
	c.BeginGroup()
	c.PushStringID(label)
	c.PushMultiItemsWidths(components)
	for i := 0; i < components; i++ {
		c.PushID(ID(i))
		if c.DragIntEx("##v", &v[i], v_speed, v_min, v_max, format) {
			value_changed = true
		}
		c.SameLineEx(0, c.Style.ItemInnerSpacing.X)
		c.PopID()
		c.PopItemWidth()
	}
	c.PopID()

	n := c.FindRenderedTextEnd(label)
	c.TextUnformatted(label[:n])
	c.EndGroup()

	return value_changed
}

func (c *Context) DragFloat(label string, v *float64) bool {
	return c.DragFloatEx(label, v, 1, 0, 0, "%.3f", 1)
}

func (c *Context) DragFloatEx(label string, v *float64, v_speed, v_min, v_max float64, format string, power float64) bool {
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	style := c.Style
	id := window.GetID(label)
	w := c.CalcItemWidth()

	label_size := c.CalcTextSizeEx(label, true, -1)
	frame_bb := f64.Rectangle{
		window.DC.CursorPos,
		window.DC.CursorPos.Add(f64.Vec2{w, label_size.Y + style.FramePadding.Y*2.0}),
	}
	inner_bb := f64.Rectangle{
		frame_bb.Min.Add(style.FramePadding),
		frame_bb.Max.Sub(style.FramePadding),
	}
	total_bb_x := 0.0
	if label_size.X > 0 {
		total_bb_x = style.ItemInnerSpacing.X + label_size.X
	}
	total_bb := f64.Rectangle{
		frame_bb.Min,
		frame_bb.Max.Add(f64.Vec2{total_bb_x, 0}),
	}

	// NB- we don't call ItemSize() yet because we may turn into a text edit box below
	if !c.ItemAddEx(total_bb, id, &frame_bb) {
		c.ItemSizeBBEx(total_bb, style.FramePadding.Y)
		return false
	}
	hovered := c.ItemHoverable(frame_bb, id)

	if format == "" {
		format = "%.3f"
	}

	// Tabbing or CTRL-clicking on Drag turns it into an input box
	start_text_input := false
	tab_focus_requested := c.FocusableItemRegister(window, id)
	if tab_focus_requested || (hovered && (c.IO.MouseClicked[0] || c.IO.MouseDoubleClicked[0])) || c.NavActivateId == id || (c.NavInputId == id && c.ScalarAsInputTextId != id) {
		c.SetActiveID(id, window)
		c.SetFocusID(id, window)
		c.FocusWindow(window)
		c.ActiveIdAllowNavDirFlags = (1 << uint(DirUp)) | (1 << uint(DirDown))
		if tab_focus_requested || c.IO.KeyCtrl || c.IO.MouseDoubleClicked[0] || c.NavInputId == id {
			start_text_input = true
			c.ScalarAsInputTextId = 0
		}
	}
	if start_text_input || (c.ActiveId == id && c.ScalarAsInputTextId == id) {
		return c.InputScalarAsWidgetReplacement(frame_bb, label, v, id, format)
	}

	// Actual drag behavior
	c.ItemSizeBBEx(total_bb, style.FramePadding.Y)
	value_changed := c.DragBehavior(frame_bb, id, v, v_speed, v_min, v_max, format, power)

	// Display value using user-provided display format so user can add prefix/suffix/decorations to the value.
	value := fmt.Sprintf(format, *v)
	c.RenderTextClippedEx(frame_bb.Min, frame_bb.Max, value, nil, f64.Vec2{0.5, 0.5}, nil)

	if label_size.X > 0.0 {
		c.RenderText(f64.Vec2{frame_bb.Max.X + style.ItemInnerSpacing.X, inner_bb.Min.Y}, label)
	}

	return value_changed
}

func (c *Context) DragFloatN(label string, v []float64, v_speed, v_min, v_max float64, format string, power float64) bool {
	components := len(v)
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	value_changed := false
	c.BeginGroup()
	c.PushStringID(label)
	c.PushMultiItemsWidths(components)
	for i := 0; i < components; i++ {
		c.PushID(ID(i))
		if c.DragFloatEx("##v", &v[i], v_speed, v_min, v_max, format, power) {
			value_changed = true
		}
		c.SameLineEx(0, c.Style.ItemInnerSpacing.X)
		c.PopID()
		c.PopItemWidth()
	}
	c.PopID()

	n := c.FindRenderedTextEnd(label)
	c.TextUnformatted(label[:n])
	c.EndGroup()

	return value_changed
}

func (c *Context) DragBehavior(frame_bb f64.Rectangle, id ID, v *float64, v_speed, v_min, v_max float64, format string, power float64) bool {
	style := &c.Style

	// Draw frame
	var frame_col color.RGBA
	switch {
	case c.ActiveId == id:
		frame_col = c.GetColorFromStyle(ColFrameBgActive)
	case c.HoveredId == id:
		frame_col = c.GetColorFromStyle(ColFrameBgHovered)
	default:
		frame_col = c.GetColorFromStyle(ColFrameBg)
	}
	c.RenderNavHighlight(frame_bb, id)
	c.RenderFrameEx(frame_bb.Min, frame_bb.Max, frame_col, true, style.FrameRounding)

	// Process interacting with the drag
	if c.ActiveId == id {
		if c.ActiveIdSource == InputSourceMouse && !c.IO.MouseDown[0] {
			c.ClearActiveID()
		} else if c.ActiveIdSource == InputSourceNav && c.NavActivatePressedId == id && !c.ActiveIdIsJustActivated {
			c.ClearActiveID()
		}
	}
	if c.ActiveId != id {
		return false
	}

	// Default tweak speed
	if v_speed == 0.0 && (v_max-v_min) != 0.0 && (v_max-v_min) < math.MaxFloat32 {
		v_speed = (v_max - v_min) * c.DragSpeedDefaultRatio
	}

	if c.ActiveIdIsJustActivated {
		// Lock current value on click
		c.DragCurrentValue = *v
		c.DragLastMouseDelta = f64.Vec2{0, 0}
	}

	mouse_drag_delta := c.GetMouseDragDelta(0, 1.0)
	adjust_delta := 0.0
	if c.ActiveIdSource == InputSourceMouse && c.IsMousePosValid() {
		adjust_delta := mouse_drag_delta.X - c.DragLastMouseDelta.X
		if c.IO.KeyShift && c.DragSpeedScaleFast >= 0.0 {
			adjust_delta *= c.DragSpeedScaleFast
		}
		if c.IO.KeyAlt && c.DragSpeedScaleSlow >= 0.0 {
			adjust_delta *= c.DragSpeedScaleSlow
		}
		c.DragLastMouseDelta.X = mouse_drag_delta.X
	}

	if c.ActiveIdSource == InputSourceNav {
		decimal_precision := ParseFormatPrecision(format, 3)
		adjust_delta = c.GetNavInputAmount2dEx(NavDirSourceFlagsKeyboard|NavDirSourceFlagsPadDPad, InputReadModeRepeatFast, 1.0/10.0, 10.0).X
		v_speed = math.Max(v_speed, GetMinimumStepAtDecimalPrecision(decimal_precision))
	}
	adjust_delta *= v_speed

	// Avoid applying the saturation when we are _already_ past the limits and heading in the same direction, so e.g. if range is 0..255, current value is 300 and we are pushing to the right side, keep the 300
	v_cur := c.DragCurrentValue
	if v_min < v_max && ((v_cur >= v_max && adjust_delta > 0.0) || (v_cur <= v_min && adjust_delta < 0.0)) {
		adjust_delta = 0.0
	}

	if math.Abs(adjust_delta) > 0.0 {
		if math.Abs(power-1.0) > 0.001 {
			// Logarithmic curve on both side of 0.0
			v0_abs := math.Abs(v_cur)
			v0_sign := f64.SignStrict(v_cur)
			v1 := math.Pow(v0_abs, 1.0/power) + (adjust_delta * v0_sign)
			v1_abs := math.Abs(v1)
			v1_sign := f64.SignStrict(v1)                       // Crossed sign line
			v_cur = math.Pow(v1_abs, power) * v0_sign * v1_sign // Reapply sign
		} else {
			v_cur += adjust_delta
		}

		// Clamp
		if v_min < v_max {
			v_cur = f64.Clamp(v_cur, v_min, v_max)
		}
		c.DragCurrentValue = v_cur
	}

	// Round to user desired precision, then apply
	value_changed := false
	v_cur = RoundScalarWithFormat(format, v_cur)
	if *v != v_cur {
		*v = v_cur
		value_changed = true
	}

	return value_changed
}

func (c *Context) DragFloatRange2(label string, v_current_min, v_current_max *float64) bool {
	return c.DragFloatRange2Ex(label, v_current_min, v_current_max, 1, 0, 0, "%.3f", "", 1)
}

func (c *Context) DragFloatRange2Ex(label string, v_current_min, v_current_max *float64, v_speed, v_min, v_max float64, format, format_max string, power float64) bool {
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	c.PushStringID(label)
	c.BeginGroup()
	c.PushMultiItemsWidths(2)

	min := v_min
	max := math.Min(v_max, *v_current_max)
	if v_min >= v_max {
		min = -math.MaxFloat32
		max = *v_current_max
	}

	value_changed := c.DragFloatEx("##min", v_current_min, v_speed, min, max, format, power)
	c.PopItemWidth()
	c.SameLineEx(0, c.Style.ItemInnerSpacing.X)

	min = math.Max(v_min, *v_current_min)
	max = v_max
	if v_min >= v_max {
		min = *v_current_min
		max = math.MaxFloat32
	}
	if format_max == "" {
		format_max = format
	}
	if c.DragFloatEx("##max", v_current_max, v_speed, min, max, format_max, power) {
		value_changed = true
	}

	c.PopItemWidth()
	c.SameLineEx(0, c.Style.ItemInnerSpacing.X)

	n := c.FindRenderedTextEnd(label)
	c.TextUnformatted(label[:n])
	c.EndGroup()
	c.PopID()

	return value_changed
}

func (c *Context) DragIntRange2(label string, v_current_min, v_current_max *int) bool {
	return c.DragIntRange2Ex(label, v_current_min, v_current_max, 1.0, 0, 0, "%.0f", "")
}

func (c *Context) DragIntRange2Ex(label string, v_current_min, v_current_max *int, v_speed float64, v_min, v_max int, format, format_max string) bool {
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	c.PushStringID(label)
	c.BeginGroup()
	c.PushMultiItemsWidths(2)

	min := v_min
	max := mathutil.Min(v_max, *v_current_max)
	if v_min >= v_max {
		min = math.MinInt32
		max = *v_current_max
	}
	value_changed := c.DragIntEx("##min", v_current_min, v_speed, min, max, format)
	c.PopItemWidth()
	c.SameLineEx(0, c.Style.ItemInnerSpacing.X)

	min = mathutil.Max(v_min, *v_current_min)
	max = v_max
	if v_min >= v_max {
		min = *v_current_min
		max = math.MaxInt32
	}
	if format_max == "" {
		format_max = format
	}
	if c.DragIntEx("##max", v_current_max, v_speed, min, max, format) {
		value_changed = true
	}

	c.PopItemWidth()
	c.SameLineEx(0, c.Style.ItemInnerSpacing.X)

	n := c.FindRenderedTextEnd(label)
	c.TextUnformatted(label[:n])
	c.EndGroup()
	c.PopID()

	return value_changed
}

func (c *Context) DragInt2(label string, v []int) bool {
	return c.DragInt2Ex(label, v, 1.0, 0, 0, "%.0f")
}

func (c *Context) DragInt2Ex(label string, v []int, v_speed float64, v_min, v_max int, format string) bool {
	return c.DragIntN(label, v[:2], v_speed, v_min, v_max, format)
}

func (c *Context) DragInt3(label string, v []int) bool {
	return c.DragInt2Ex(label, v, 1.0, 0, 0, "%.0f")
}

func (c *Context) DragInt3Ex(label string, v []int, v_speed float64, v_min, v_max int, format string) bool {
	return c.DragIntN(label, v[:3], v_speed, v_min, v_max, format)
}

func (c *Context) DragInt4(label string, v []int) bool {
	return c.DragInt4Ex(label, v, 1.0, 0, 0, "%.0f")
}

func (c *Context) DragInt4Ex(label string, v []int, v_speed float64, v_min, v_max int, format string) bool {
	return c.DragIntN(label, v[:4], v_speed, v_min, v_max, format)
}

func (c *Context) DragFloat2(label string, v []float64) bool {
	return c.DragFloat2Ex(label, v, 1.0, 0.0, 0.0, "%.3f", 1.0)
}

func (c *Context) DragFloat2Ex(label string, v []float64, v_speed, v_min, v_max float64, format string, power float64) bool {
	return c.DragFloatN(label, v[:2], v_speed, v_min, v_max, format, 1)
}

func (c *Context) DragFloat3(label string, v []float64) bool {
	return c.DragFloat3Ex(label, v, 1.0, 0.0, 0.0, "%.3f", 1.0)
}

func (c *Context) DragFloat3Ex(label string, v []float64, v_speed, v_min, v_max float64, format string, power float64) bool {
	return c.DragFloatN(label, v[:3], v_speed, v_min, v_max, format, 1)
}

func (c *Context) DragFloat4(label string, v []float64) bool {
	return c.DragFloat4Ex(label, v, 1.0, 0.0, 0.0, "%.3f", 1.0)
}

func (c *Context) DragFloat4Ex(label string, v []float64, v_speed, v_min, v_max float64, format string, power float64) bool {
	return c.DragFloatN(label, v[:4], v_speed, v_min, v_max, format, 1)
}

func (c *Context) DragV2(label string, v *f64.Vec2) bool {
	return c.DragV2Ex(label, v, 1.0, 0.0, 0.0, "%.3f", 1.0)
}

func (c *Context) DragV2Ex(label string, v *f64.Vec2, v_speed, v_min, v_max float64, format string, power float64) bool {
	f := [...]float64{v.X, v.Y}
	r := c.DragFloatN(label, f[:2], v_speed, v_min, v_max, format, 1)
	v.X, v.Y = f[0], f[1]
	return r
}

func (c *Context) DragV3(label string, v *f64.Vec3) bool {
	return c.DragV3Ex(label, v, 1.0, 0.0, 0.0, "%.3f", 1.0)
}

func (c *Context) DragV3Ex(label string, v *f64.Vec3, v_speed, v_min, v_max float64, format string, power float64) bool {
	f := [...]float64{v.X, v.Y, v.Z}
	r := c.DragFloatN(label, f[:3], v_speed, v_min, v_max, format, 1)
	v.X, v.Y, v.Z = f[0], f[1], f[2]
	return r
}

func (c *Context) DragV4(label string, v *f64.Vec4) bool {
	return c.DragV4Ex(label, v, 1.0, 0.0, 0.0, "%.3f", 1.0)
}

func (c *Context) DragV4Ex(label string, v *f64.Vec4, v_speed, v_min, v_max float64, format string, power float64) bool {
	f := [...]float64{v.X, v.Y, v.Z, v.W}
	r := c.DragFloatN(label, f[:4], v_speed, v_min, v_max, format, 1)
	v.X, v.Y, v.Z, v.W = f[0], f[1], f[2], f[3]
	return r
}
