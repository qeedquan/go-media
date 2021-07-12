package imgui

import (
	"fmt"
	"image/color"
	"math"

	"github.com/qeedquan/go-media/math/f64"
)

type SliderFlags int

const (
	SliderFlagsVertical SliderFlags = 1 << 0
)

func (c *Context) SliderBehaviorCalcRatioFromValue(v, v_min, v_max, power, linear_zero_pos float64) float64 {
	if v_min == v_max {
		return 0
	}

	is_non_linear := (power < 1.0-0.00001) || (power > 1.0+0.00001)
	var v_clamped float64
	if v_min < v_max {
		v_clamped = f64.Clamp(v, v_min, v_max)
	} else {
		v_clamped = f64.Clamp(v, v_max, v_min)
	}

	if is_non_linear {
		if v_clamped < 0 {
			f := 1.0 - (v_clamped-v_min)/(math.Min(0.0, v_max)-v_min)
			return (1.0 - math.Pow(f, 1.0/power)) * linear_zero_pos
		} else {
			f := (v_clamped - math.Max(0.0, v_min)) / (v_max - math.Max(0.0, v_min))
			return linear_zero_pos + math.Pow(f, 1.0/power)*(1.0-linear_zero_pos)
		}
	}

	// Linear slider
	return (v_clamped - v_min) / (v_max - v_min)
}

func (c *Context) SliderBehavior(frame_bb f64.Rectangle, id ID, v *float64, v_min, v_max, power float64, format string, flags SliderFlags) bool {
	window := c.GetCurrentWindow()
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

	is_non_linear := (power < 1.0-0.00001) || (power > 1.0+0.00001)
	is_horizontal := (flags & SliderFlagsVertical) == 0
	is_decimal := ParseFormatPrecision(format, 3) > 0

	grab_padding := 2.0
	var slider_sz float64
	if is_horizontal {
		slider_sz = frame_bb.Dx() - grab_padding*2.0
	} else {
		slider_sz = frame_bb.Dy() - grab_padding*2.0
	}

	grab_sz := 0.0
	if is_decimal {
		grab_sz = math.Min(style.GrabMinSize, slider_sz)
	} else {
		// Integer sliders, if possible have the grab size represent 1 unit
		v_ratio := math.Abs(v_max-v_min) + 1
		grab_sz = math.Min(math.Max(slider_sz/v_ratio, style.GrabMinSize), slider_sz)
	}
	slider_usable_sz := slider_sz - grab_sz
	var slider_usable_pos_min, slider_usable_pos_max float64
	if is_horizontal {
		slider_usable_pos_min = frame_bb.Min.X + grab_padding + grab_sz*0.5
		slider_usable_pos_max = frame_bb.Max.X - grab_padding - grab_sz*0.5
	} else {
		slider_usable_pos_min = frame_bb.Min.Y + grab_padding + grab_sz*0.5
		slider_usable_pos_max = frame_bb.Max.Y - grab_padding - grab_sz*0.5
	}

	// For logarithmic sliders that cross over sign boundary we want the exponential increase to be symmetric around 0.0f
	linear_zero_pos := 0.0 // 0.0->1.0f
	if v_min*v_max < 0.0 {
		// Different sign
		linear_dist_min_to_0 := math.Pow(math.Abs(0.0-v_min), 1.0/power)
		linear_dist_max_to_0 := math.Pow(math.Abs(v_max-0.0), 1.0/power)
		linear_zero_pos = linear_dist_min_to_0 / (linear_dist_min_to_0 + linear_dist_max_to_0)
	} else {
		// Same sign
		linear_zero_pos = 0
		if v_min < 0 {
			linear_zero_pos = 1
		}
	}

	// Process interacting with the slider
	value_changed := false
	if c.ActiveId == id {
		set_new_value := false
		clicked_t := 0.0
		if c.ActiveIdSource == InputSourceMouse {
			if !c.IO.MouseDown[0] {
				c.ClearActiveID()
			} else {
				mouse_abs_pos := c.IO.MousePos.Y
				if is_horizontal {
					mouse_abs_pos = c.IO.MousePos.X
				}
				clicked_t = 0.0
				if slider_usable_sz > 0 {
					clicked_t = f64.Clamp((mouse_abs_pos-slider_usable_pos_min)/slider_usable_sz, 0.0, 1.0)
				}
				if !is_horizontal {
					clicked_t = 1.0 - clicked_t
				}
				set_new_value = true
			}
		} else if c.ActiveIdSource == InputSourceNav {
			delta2 := c.GetNavInputAmount2dEx(NavDirSourceFlagsKeyboard|NavDirSourceFlagsPadDPad, InputReadModeRepeatFast, 0.0, 0.0)
			delta := -delta2.Y
			if is_horizontal {
				delta = delta2.X
			}
			if c.NavActivatePressedId == id && !c.ActiveIdIsJustActivated {
				c.ClearActiveID()
			} else if delta != 0.0 {
				clicked_t = c.SliderBehaviorCalcRatioFromValue(*v, v_min, v_max, power, linear_zero_pos)
				if !is_decimal && !is_non_linear {
					if math.Abs(v_max-v_min) <= 100.0 || c.IsNavInputDown(NavInputTweakSlow) {
						// Gamepad/keyboard tweak speeds in integer steps
						if delta < 0 {
							delta = -1 / (v_max - v_min)
						} else {
							delta = 1 / (v_max - v_min)
						}
					} else {
						delta /= 100.0
					}
				} else {
					delta /= 100.0 // Gamepad/keyboard tweak speeds in % of slider bounds
					if c.IsNavInputDown(NavInputTweakSlow) {
						delta /= 10.0
					}
				}

				if c.IsNavInputDown(NavInputTweakFast) {
					delta *= 10.0
				}
				set_new_value = true
				// This is to avoid applying the saturation when already past the limits
				if (clicked_t >= 1.0 && delta > 0.0) || (clicked_t <= 0.0 && delta < 0.0) {
					set_new_value = false
				} else {
					clicked_t = f64.Saturate(clicked_t + delta)
				}
			}
		}

		if set_new_value {
			var new_value float64

			if is_non_linear {
				// Account for logarithmic scale on both sides of the zero
				if clicked_t < linear_zero_pos {
					// Negative: rescale to the negative range before powering
					a := 1.0 - (clicked_t / linear_zero_pos)
					a = math.Pow(a, power)
					new_value = f64.Lerp(a, math.Min(v_max, 0.0), v_min)
				} else {
					// Positive: rescale to the positive range before powering
					var a float64
					if math.Abs(linear_zero_pos-1.0) > 1.e-6 {
						a = (clicked_t - linear_zero_pos) / (1.0 - linear_zero_pos)
					} else {
						a = clicked_t
					}
					a = math.Pow(a, power)
					new_value = f64.Lerp(a, math.Max(v_min, 0.0), v_max)
				}
			} else {
				// Linear slider
				new_value = f64.Lerp(clicked_t, v_min, v_max)
			}

			// Round past decimal precision
			new_value = RoundScalarWithFormat(format, new_value)
			if *v != new_value {
				*v = new_value
				value_changed = true
			}
		}
	}

	grab_t := c.SliderBehaviorCalcRatioFromValue(*v, v_min, v_max, power, linear_zero_pos)
	if !is_horizontal {
		grab_t = 1.0 - grab_t
	}
	grab_pos := f64.Lerp(grab_t, slider_usable_pos_min, slider_usable_pos_max)
	var grab_bb f64.Rectangle
	if is_horizontal {
		grab_bb = f64.Rectangle{
			f64.Vec2{grab_pos - grab_sz*0.5, frame_bb.Min.Y + grab_padding},
			f64.Vec2{grab_pos + grab_sz*0.5, frame_bb.Max.Y - grab_padding},
		}
	} else {
		grab_bb = f64.Rectangle{
			f64.Vec2{frame_bb.Min.X + grab_padding, grab_pos - grab_sz*0.5},
			f64.Vec2{frame_bb.Max.X - grab_padding, grab_pos + grab_sz*0.5},
		}
	}

	var col color.RGBA
	if c.ActiveId == id {
		col = c.GetColorFromStyle(ColSliderGrabActive)
	} else {
		col = c.GetColorFromStyle(ColSliderGrabActive)
	}
	window.DrawList.AddRectFilledEx(grab_bb.Min, grab_bb.Max, col, style.GrabRounding, DrawCornerFlagsAll)

	return value_changed
}

func (c *Context) SliderFloat(label string, v *float64, v_min, v_max float64) bool {
	return c.SliderFloatEx(label, v, v_min, v_max, "%.3f", 1.0)
}

func (c *Context) SliderInt(label string, v *int, v_min, v_max int) bool {
	return c.SliderIntEx(label, v, v_min, v_max, "%.0f")
}

func (c *Context) SliderIntEx(label string, v *int, v_min, v_max int, format string) bool {
	if format == "" {
		format = "%.0f"
	}

	v_f := float64(*v)
	value_changed := c.SliderFloatEx(label, &v_f, float64(v_min), float64(v_max), format, 1)
	*v = int(v_f)
	return value_changed
}

// Use power!=1.0 for logarithmic sliders.
// Adjust format to decorate the value with a prefix or a suffix.
//   "%.3f"         1.234
//   "%5.2f secs"   01.23 secs
//   "Gold: %.0f"   Gold: 1
func (c *Context) SliderFloatEx(label string, v *float64, v_min, v_max float64, format string, power float64) bool {
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	style := c.Style
	id := window.GetID(label)
	w := c.CalcItemWidth()

	label_size := c.CalcTextSizeEx(label, true, -1.0)
	frame_bb := f64.Rectangle{
		window.DC.CursorPos,
		window.DC.CursorPos.Add(f64.Vec2{w, label_size.Y + style.FramePadding.Y*2.0}),
	}
	x := 0.0
	if label_size.X > 0 {
		x = style.ItemInnerSpacing.X + label_size.X
	}
	total_bb := f64.Rectangle{
		frame_bb.Min,
		frame_bb.Max.Add(f64.Vec2{x, 0}),
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

	// Tabbing or CTRL-clicking on Slider turns it into an input box
	start_text_input := false
	tab_focus_requested := c.FocusableItemRegister(window, id)
	if tab_focus_requested || (hovered && c.IO.MouseClicked[0]) || c.NavActivateId == id || (c.NavInputId == id && c.ScalarAsInputTextId != id) {
		c.SetActiveID(id, window)
		c.SetFocusID(id, window)
		c.FocusWindow(window)
		c.ActiveIdAllowNavDirFlags = 1<<uint(DirUp) | 1<<uint(DirDown)
		if tab_focus_requested || c.IO.KeyCtrl || c.NavInputId == id {
			start_text_input = true
			c.ScalarAsInputTextId = 0
		}
	}

	if start_text_input || (c.ActiveId == id && c.ScalarAsInputTextId == id) {
		return c.InputScalarAsWidgetReplacement(frame_bb, label, *v, id, format)
	}

	// Actual slider behavior + render grab
	c.ItemSizeBBEx(total_bb, style.FramePadding.Y)
	value_changed := c.SliderBehavior(frame_bb, id, v, v_min, v_max, power, format, 0)

	// Display value using user-provided display format so user can add prefix/suffix/decorations to the value.
	value := fmt.Sprintf(format, *v)
	c.RenderTextClippedEx(frame_bb.Min, frame_bb.Max, value, nil, f64.Vec2{0.5, 0.5}, nil)
	if label_size.X > 0.0 {
		c.RenderText(f64.Vec2{
			frame_bb.Max.X + style.ItemInnerSpacing.X,
			frame_bb.Min.Y + style.FramePadding.Y},
			label,
		)
	}

	return value_changed
}

func (c *Context) VSliderInt(label string, size f64.Vec2, v *int, v_min, v_max int) bool {
	return c.VSliderIntEx(label, size, v, v_min, v_max, "%.0f")
}

func (c *Context) VSliderIntEx(label string, size f64.Vec2, v *int, v_min, v_max int, format string) bool {
	if format == "" {
		format = "%.0f"
	}
	v_f := float64(*v)
	value_changed := c.VSliderFloatEx(label, size, &v_f, float64(v_min), float64(v_max), format, 1.0)
	*v = int(v_f)
	return value_changed
}

func (c *Context) VSliderFloat(label string, size f64.Vec2, v *float64, v_min, v_max float64) bool {
	return c.VSliderFloatEx(label, size, v, v_min, v_max, "%.3f", 1)
}

func (c *Context) VSliderFloatEx(label string, size f64.Vec2, v *float64, v_min, v_max float64, format string, power float64) bool {
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	style := &c.Style
	id := window.GetID(label)

	label_size := c.CalcTextSizeEx(label, true, -1)
	frame_bb := f64.Rectangle{window.DC.CursorPos, window.DC.CursorPos.Add(size)}
	x := 0.0
	if label_size.X > 0.0 {
		x += style.ItemInnerSpacing.X + label_size.X
	}
	bb := f64.Rectangle{frame_bb.Min, frame_bb.Max.Add(f64.Vec2{x, 0})}

	c.ItemSizeBBEx(bb, style.FramePadding.Y)
	if !c.ItemAdd(frame_bb, id) {
		return false
	}
	hovered := c.ItemHoverable(frame_bb, id)

	if format == "" {
		format = "%.3f"
	}

	if (hovered && c.IO.MouseClicked[0]) || c.NavActivateId == id || c.NavInputId == id {
		c.SetActiveID(id, window)
		c.SetFocusID(id, window)
		c.FocusWindow(window)
		c.ActiveIdAllowNavDirFlags = (1 << uint(DirLeft)) | (1 << uint(DirRight))
	}

	// Actual slider behavior + render grab
	value_changed := c.SliderBehavior(frame_bb, id, v, v_min, v_max, power, format, SliderFlagsVertical)

	// Display value using user-provided display format so user can add prefix/suffix/decorations to the value.
	// For the vertical slider we allow centered text to overlap the frame padding
	value := fmt.Sprintf(format, *v)
	c.RenderTextClippedEx(f64.Vec2{frame_bb.Min.X, frame_bb.Min.Y + style.FramePadding.Y}, frame_bb.Max, value, nil, f64.Vec2{0.5, 0.0}, nil)
	if label_size.X > 0.0 {
		c.RenderText(f64.Vec2{frame_bb.Max.X + style.ItemInnerSpacing.X, frame_bb.Min.Y + style.FramePadding.Y}, label)
	}

	return value_changed
}

// Create text input in place of a slider (when CTRL+Clicking on slider)
// FIXME: Logic is messy and confusing.
func (c *Context) InputScalarAsWidgetReplacement(bb f64.Rectangle, label string, data interface{}, id ID, format string) bool {
	window := c.GetCurrentWindow()

	// Our replacement widget will override the focus ID (registered previously to allow for a TAB focus to happen)
	// On the first frame, g.ScalarAsInputTextId == 0, then on subsequent frames it becomes == id
	c.SetActiveID(c.ScalarAsInputTextId, window)
	c.ActiveIdAllowNavDirFlags = (1 << uint(DirUp)) | (1 << uint(DirDown))
	c.SetHoveredID(0)
	c.FocusableItemUnregister(window)

	format = ParseFormatTrimDecorations(format)
	data_buf := DataTypeFormatString(data, format)
	flags := InputTextFlagsAutoSelectAll | InputTextFlagsCharsDecimal
	switch data.(type) {
	case float32, float64:
		flags |= InputTextFlagsCharsScientific
	}
	text_value_changed := c.InputTextEx(label, data_buf, bb.Size(), flags, nil)
	// First frame we started displaying the InputText widget
	if c.ScalarAsInputTextId == 0 {
		// InputText ID expected to match the Slider ID
		assert(c.ActiveId == id)
		c.ScalarAsInputTextId = c.ActiveId
		c.SetHoveredID(id)
	}
	if text_value_changed {
		return DataTypeApplyOpFromText(data_buf, string(c.InputTextState.InitialText), data, "")
	}

	return false
}

func (c *Context) SliderAngle(label string, v_rad *float64) bool {
	return c.SliderAngleEx(label, v_rad, -360, 360)
}

func (c *Context) SliderAngleEx(label string, v_rad *float64, v_degrees_min, v_degrees_max float64) bool {
	v_deg := (*v_rad) * 360.0 / (2 * math.Pi)
	value_changed := c.SliderFloatEx(label, &v_deg, v_degrees_min, v_degrees_max, "%.0f deg", 1.0)
	*v_rad = v_deg * (2 * math.Pi) / 360.0
	return value_changed
}

// Add multiple sliders on 1 line for compact edition of multiple components
func (c *Context) SliderFloatN(label string, v []float64, v_min, v_max float64, format string, power float64) bool {
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
		changed := c.SliderFloatEx("##v", &v[i], v_min, v_max, format, power)
		if changed {
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

func (c *Context) PushMultiItemsWidths(components int) {
	c.PushMultiItemsWidthsEx(components, 0)
}

func (c *Context) PushMultiItemsWidthsEx(components int, w_full float64) {
	window := c.GetCurrentWindow()
	style := &c.Style
	if w_full <= 0.0 {
		w_full = c.CalcItemWidth()
	}
	w_item_one := math.Max(1.0, float64(int(w_full-style.ItemInnerSpacing.X*float64(components-1)))/float64(components))
	w_item_last := math.Max(1, float64(int(w_full-(w_item_one+style.ItemInnerSpacing.X)*float64(components-1))))
	window.DC.ItemWidthStack = append(window.DC.ItemWidthStack, w_item_last)
	for i := 0; i < components-1; i++ {
		window.DC.ItemWidthStack = append(window.DC.ItemWidthStack, w_item_one)
	}
	window.DC.ItemWidth = window.DC.ItemWidthStack[len(window.DC.ItemWidthStack)-1]
}

func (c *Context) SliderIntN(label string, v []int, v_min, v_max int, format string) bool {
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
		if c.SliderIntEx("##v", &v[i], v_min, v_max, format) {
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

func (c *Context) SliderInt2(label string, v []int, v_min, v_max int) bool {
	return c.SliderInt2Ex(label, v, v_min, v_max, "%.0f")
}

func (c *Context) SliderInt2Ex(label string, v []int, v_min, v_max int, format string) bool {
	return c.SliderIntN(label, v[:2], v_min, v_max, format)
}

func (c *Context) SliderInt3(label string, v []int, v_min, v_max int) bool {
	return c.SliderInt3Ex(label, v, v_min, v_max, "%.0f")
}

func (c *Context) SliderInt3Ex(label string, v []int, v_min, v_max int, format string) bool {
	return c.SliderIntN(label, v[:3], v_min, v_max, format)
}

func (c *Context) SliderInt4(label string, v []int, v_min, v_max int) bool {
	return c.SliderInt4Ex(label, v, v_min, v_max, "%.0f")
}

func (c *Context) SliderInt4Ex(label string, v []int, v_min, v_max int, format string) bool {
	return c.SliderIntN(label, v[:4], v_min, v_max, format)
}

func (c *Context) SliderFloat2(label string, v []float64, v_min, v_max float64) bool {
	return c.SliderFloat2Ex(label, v, v_min, v_max, "%.3f", 1.0)
}

func (c *Context) SliderFloat2Ex(label string, v []float64, v_min, v_max float64, format string, power float64) bool {
	return c.SliderFloatN(label, v[:2], v_min, v_max, format, power)
}

func (c *Context) SliderFloat3(label string, v []float64, v_min, v_max float64) bool {
	return c.SliderFloat3Ex(label, v, v_min, v_max, "%.3f", 1.0)
}

func (c *Context) SliderFloat3Ex(label string, v []float64, v_min, v_max float64, format string, power float64) bool {
	return c.SliderFloatN(label, v[:3], v_min, v_max, format, power)
}

func (c *Context) SliderFloat4(label string, v []float64, v_min, v_max float64) bool {
	return c.SliderFloat4Ex(label, v, v_min, v_max, "%.3f", 1.0)
}

func (c *Context) SliderFloat4Ex(label string, v []float64, v_min, v_max float64, format string, power float64) bool {
	return c.SliderFloatN(label, v[:4], v_min, v_max, format, power)
}

func (c *Context) SliderV2(label string, v *f64.Vec2, v_min, v_max float64) bool {
	return c.SliderV2Ex(label, v, v_min, v_max, "%.3f", 1.0)
}

func (c *Context) SliderV2Ex(label string, v *f64.Vec2, v_min, v_max float64, format string, power float64) bool {
	f := [...]float64{v.X, v.Y}
	r := c.SliderFloatN(label, f[:2], v_min, v_max, format, power)
	v.X, v.Y = f[0], f[1]
	return r
}

func (c *Context) SliderV3(label string, v *f64.Vec3, v_min, v_max float64) bool {
	return c.SliderV3Ex(label, v, v_min, v_max, "%.3f", 1.0)
}

func (c *Context) SliderV3Ex(label string, v *f64.Vec3, v_min, v_max float64, format string, power float64) bool {
	f := [...]float64{v.X, v.Y, v.Z}
	r := c.SliderFloatN(label, f[:3], v_min, v_max, format, power)
	v.X, v.Y, v.Z = f[0], f[1], f[2]
	return r
}

func (c *Context) SliderV4(label string, v *f64.Vec4, v_min, v_max float64) bool {
	return c.SliderV4Ex(label, v, v_min, v_max, "%.3f", 1.0)
}

func (c *Context) SliderV4Ex(label string, v *f64.Vec4, v_min, v_max float64, format string, power float64) bool {
	f := [...]float64{v.X, v.Y, v.Z, v.W}
	r := c.SliderFloatN(label, f[:4], v_min, v_max, format, power)
	v.X, v.Y, v.Z, v.W = f[0], f[1], f[2], f[3]
	return r
}
