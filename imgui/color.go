package imgui

import (
	"fmt"
	"image/color"
	"math"

	"github.com/qeedquan/go-media/image/chroma"
	"github.com/qeedquan/go-media/math/f64"
	"github.com/qeedquan/go-media/math/mathutil"
)

type ColorEditFlags int

const (
	ColorEditFlagsNoAlpha        ColorEditFlags = 1 << 1 //              // ColorEdit ColorPicker ColorButton: ignore Alpha component (read 3 components from the input pointer).
	ColorEditFlagsNoPicker       ColorEditFlags = 1 << 2 //              // ColorEdit: disable picker when clicking on colored square.
	ColorEditFlagsNoOptions      ColorEditFlags = 1 << 3 //              // ColorEdit: disable toggling options menu when right-clicking on inputs/small preview.
	ColorEditFlagsNoSmallPreview ColorEditFlags = 1 << 4 //              // ColorEdit ColorPicker: disable colored square preview next to the inputs. (e.g. to show only the inputs)
	ColorEditFlagsNoInputs       ColorEditFlags = 1 << 5 //              // ColorEdit ColorPicker: disable inputs sliders/text widgets (e.g. to show only the small preview colored square).
	ColorEditFlagsNoTooltip      ColorEditFlags = 1 << 6 //              // ColorEdit ColorPicker ColorButton: disable tooltip when hovering the preview.
	ColorEditFlagsNoLabel        ColorEditFlags = 1 << 7 //              // ColorEdit ColorPicker: disable display of inline text label (the label is still forwarded to the tooltip and picker).
	ColorEditFlagsNoSidePreview  ColorEditFlags = 1 << 8 //              // ColorPicker: disable bigger color preview on right side of the picker use small colored square preview instead.

	// User Options (right-click on widget to change some of them). You can set application defaults using SetColorEditOptions(). The idea is that you probably don't want to override them in most of your calls let the user choose and/or call SetColorEditOptions() during startup.
	ColorEditFlagsAlphaBar         ColorEditFlags = 1 << 9  //              // ColorEdit ColorPicker: show vertical alpha bar/gradient in picker.
	ColorEditFlagsAlphaPreview     ColorEditFlags = 1 << 10 //              // ColorEdit ColorPicker ColorButton: display preview as a transparent color over a checkerboard instead of opaque.
	ColorEditFlagsAlphaPreviewHalf ColorEditFlags = 1 << 11 //              // ColorEdit ColorPicker ColorButton: display half opaque / half checkerboard instead of opaque.
	ColorEditFlagsHDR              ColorEditFlags = 1 << 12 //              // (WIP) ColorEdit: Currently only disable 0.0f..1.0f limits in RGBA edition (note: you probably want to use ColorEditFlagsFloat flag as well).
	ColorEditFlagsRGB              ColorEditFlags = 1 << 13 // [Inputs]     // ColorEdit: choose one among RGB/HSV/HEX. ColorPicker: choose any combination using RGB/HSV/HEX.
	ColorEditFlagsHSV              ColorEditFlags = 1 << 14 // [Inputs]     // "
	ColorEditFlagsHEX              ColorEditFlags = 1 << 15 // [Inputs]     // "
	ColorEditFlagsUint8            ColorEditFlags = 1 << 16 // [DataType]   // ColorEdit ColorPicker ColorButton: _display_ values formatted as 0..255.
	ColorEditFlagsFloat            ColorEditFlags = 1 << 17 // [DataType]   // ColorEdit ColorPicker ColorButton: _display_ values formatted as 0.0f..1.0f floats instead of 0..255 integers. No round-trip of value via integers.
	ColorEditFlagsPickerHueBar     ColorEditFlags = 1 << 18 // [PickerMode] // ColorPicker: bar for Hue rectangle for Sat/Value.
	ColorEditFlagsPickerHueWheel   ColorEditFlags = 1 << 19 // [PickerMode] // ColorPicker: wheel for Hue triangle for Sat/Value.

	// Internals/Masks
	ColorEditFlags_InputsMask     ColorEditFlags = ColorEditFlagsRGB | ColorEditFlagsHSV | ColorEditFlagsHEX
	ColorEditFlags_DataTypeMask   ColorEditFlags = ColorEditFlagsUint8 | ColorEditFlagsFloat
	ColorEditFlags_PickerMask     ColorEditFlags = ColorEditFlagsPickerHueWheel | ColorEditFlagsPickerHueBar
	ColorEditFlags_OptionsDefault ColorEditFlags = ColorEditFlagsUint8 | ColorEditFlagsRGB | ColorEditFlagsPickerHueBar // Change application default using SetColorEditOptions()
)

func (c *Context) SetColorEditOptions(flags ColorEditFlags) {
	if flags&ColorEditFlags_InputsMask == 0 {
		flags |= ColorEditFlags_OptionsDefault & ColorEditFlags_InputsMask
	}
	if flags&ColorEditFlags_DataTypeMask == 0 {
		flags |= ColorEditFlags_OptionsDefault & ColorEditFlags_DataTypeMask
	}
	if flags&ColorEditFlags_PickerMask == 0 {
		flags |= ColorEditFlags_OptionsDefault & ColorEditFlags_PickerMask
	}
	assert(mathutil.IsPow2(int(flags & ColorEditFlags_InputsMask)))   // Check only 1 option is selected
	assert(mathutil.IsPow2(int(flags & ColorEditFlags_DataTypeMask))) // Check only 1 option is selected
	assert(mathutil.IsPow2(int(flags & ColorEditFlags_PickerMask)))   // Check only 1 option is selected
	c.ColorEditOptions = flags
}

func (c *Context) ColorTooltip(text string, col color.RGBA, flags ColorEditFlags) {
	c.ColorTooltipV(text, chroma.RGBA2VEC4(col), flags)
}

// Note: only access 3 floats if ImGuiColorEditFlags_NoAlpha flag is set.
func (c *Context) ColorTooltipV(text string, col f64.Vec4, flags ColorEditFlags) {
	rgba := col.ToRGBA()
	cr, cg, cb, ca := rgba.R, rgba.G, rgba.B, rgba.A
	fr, fg, fb, fa := col.X, col.Y, col.Z, col.W

	c.BeginTooltipEx(0, true)

	text_end := c.FindRenderedTextEnd(text)
	if text_end > 0 {
		c.TextUnformatted(text[:text_end])
		c.Separator()
	}
	sz := f64.Vec2{
		c.FontSize*3 + c.Style.FramePadding.Y*2,
		c.FontSize*3 + c.Style.FramePadding.Y*2,
	}
	c.ColorButtonVEx("##preview", col, flags&(ColorEditFlagsNoAlpha|ColorEditFlagsAlphaPreview|ColorEditFlagsAlphaPreviewHalf)|ColorEditFlagsNoTooltip, sz)
	c.SameLine()
	if flags&ColorEditFlagsNoAlpha != 0 {
		c.Text("#%02X%02X%02X\nR: %d, G: %d, B: %d\n(%.3f, %.3f, %.3f)", cr, cg, cb, cr, cg, cb, fr, fg, fb, fa)
	} else {
		c.Text("#%02X%02X%02X%02X\nR:%d, G:%d, B:%d, A:%d\n(%.3f, %.3f, %.3f, %.3f)", cr, cg, cb, ca, cr, cg, cb, ca, fr, fg, fb, fa)
	}
	c.EndTooltip()
}

func (c *Context) ColorEditV3(label string, col *f64.Vec4) bool {
	return c.ColorEditV4Ex(label, col, ColorEditFlagsNoAlpha)
}

func (c *Context) ColorEditV3Ex(label string, col *f64.Vec4, flags ColorEditFlags) bool {
	return c.ColorEditV4Ex(label, col, flags|ColorEditFlagsNoAlpha)
}

func (c *Context) ColorEditV4(label string, col *f64.Vec4) bool {
	return c.ColorEditV4Ex(label, col, 0)
}

// Edit colors components (each component in 0.0f..1.0f range).
// See enum ImGuiColorEditFlags_ for available options. e.g. Only access 3 floats if ImGuiColorEditFlags_NoAlpha flag is set.
// With typical options: Left-click on colored square to open color picker. Right-click to open option menu. CTRL-Click over input fields to edit them and TAB to go to next item.
func (c *Context) ColorEditV4Ex(label string, col *f64.Vec4, flags ColorEditFlags) bool {
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	style := &c.Style
	square_sz := c.GetFrameHeight()
	w_extra := 0.0
	if flags&ColorEditFlagsNoSmallPreview == 0 {
		w_extra = square_sz + style.ItemInnerSpacing.X
	}
	w_items_all := c.CalcItemWidth() - w_extra
	label_display_end := c.FindRenderedTextEnd(label)

	alpha := (flags & ColorEditFlagsNoAlpha) == 0
	hdr := (flags & ColorEditFlagsHDR) != 0
	components := 3
	if alpha {
		components = 4
	}
	flags_untouched := flags

	c.BeginGroup()
	c.PushStringID(label)
	label = label[:label_display_end]

	// If we're not showing any slider there's no point in doing any HSV conversions
	if flags&ColorEditFlagsNoInputs != 0 {
		flags = (flags &^ (ColorEditFlags_InputsMask)) | ColorEditFlagsRGB | ColorEditFlagsNoOptions
	}

	// Context menu: display and modify options (before defaults are applied)
	if flags&ColorEditFlagsNoOptions == 0 {
		c.ColorEditOptionsPopupV(*col, flags)
	}

	// Read stored options
	if flags&ColorEditFlags_InputsMask == 0 {
		flags |= (c.ColorEditOptions & ColorEditFlags_InputsMask)
	}
	if flags&ColorEditFlags_DataTypeMask == 0 {
		flags |= (c.ColorEditOptions & ColorEditFlags_DataTypeMask)
	}
	if flags&ColorEditFlags_PickerMask == 0 {
		flags |= (c.ColorEditOptions & ColorEditFlags_PickerMask)
	}
	flags |= (c.ColorEditOptions &^ (ColorEditFlags_InputsMask | ColorEditFlags_DataTypeMask | ColorEditFlags_PickerMask))

	// Convert to the formats we need
	v := *col
	f := [4]float64{v.X, v.Y, v.Z, 1}
	if alpha {
		f[3] = col.Z
	}
	if flags&ColorEditFlagsHSV != 0 {
		hsv := chroma.VEC42HSV(f64.Vec4{f[0], f[1], f[2], f[3]})
		f[0], f[1], f[2] = hsv.H, hsv.S, hsv.V
	}
	i := [4]int{F32_TO_INT8_UNBOUND(f[0]), F32_TO_INT8_UNBOUND(f[1]), F32_TO_INT8_UNBOUND(f[2]), F32_TO_INT8_UNBOUND(f[3])}

	value_changed := false
	value_changed_as_float := false

	if flags&(ColorEditFlagsRGB|ColorEditFlagsHSV) != 0 && flags&ColorEditFlagsNoInputs == 0 {
		// RGB/HSV 0..255 Sliders
		w_item_one := math.Max(1.0, float64(int((w_items_all-(style.ItemInnerSpacing.X)*float64(components-1))/float64(components))))
		w_item_last := math.Max(1.0, float64(int(w_items_all-(w_item_one+style.ItemInnerSpacing.X)*float64(components-1))))

		var text_size f64.Vec2
		if flags&ColorEditFlagsFloat != 0 {
			text_size = c.CalcTextSize("M:0.000")
		} else {
			text_size = c.CalcTextSize("M:000")
		}
		hide_prefix := w_item_one <= text_size.X
		ids := []string{"##X", "##Y", "##Z", "##W"}
		fmt_table_int := [][]string{
			{"%3.0f", "%3.0f", "%3.0f", "%3.0f"},         // Short display
			{"R:%3.0f", "G:%3.0f", "B:%3.0f", "A:%3.0f"}, // Long display for RGBA
			{"H:%3.0f", "S:%3.0f", "V:%3.0f", "A:%3.0f"}, // Long display for HSVA
		}

		fmt_table_float := [][]string{
			{"%0.3f", "%0.3f", "%0.3f", "%0.3f"},         // Short display
			{"R:%0.3f", "G:%0.3f", "B:%0.3f", "A:%0.3f"}, // Long display for RGBA
			{"H:%0.3f", "S:%0.3f", "V:%0.3f", "A:%0.3f"}, // Long display for HSVA
		}

		var fmt_idx int
		if !hide_prefix {
			if flags&ColorEditFlagsHSV != 0 {
				fmt_idx = 2
			} else {
				fmt_idx = 1
			}
		}

		c.PushItemWidth(w_item_one)
		for n := 0; n < components; n++ {
			if n > 0 {
				c.SameLineEx(0, style.ItemInnerSpacing.X)
			}
			if n+1 == components {
				c.PushItemWidth(w_item_last)
			}
			if flags&ColorEditFlagsFloat != 0 {
				v_max := 1.0
				if hdr {
					v_max = 0.0
				}
				if c.DragFloatEx(ids[n], &f[n], 1.0/255.0, 0.0, v_max, fmt_table_float[fmt_idx][n], 1) {
					value_changed = true
					value_changed_as_float = true
				}
			} else {
				v_max := 255
				if hdr {
					v_max = 0
				}
				if c.DragIntEx(ids[n], &i[n], 1.0, 0, v_max, fmt_table_int[fmt_idx][n]) {
					value_changed = true
				}
			}
			if flags&ColorEditFlagsNoOptions == 0 {
				c.OpenPopupOnItemClick("context", 1)
			}
		}
		c.PopItemWidth()
		c.PopItemWidth()
	} else if (flags&ColorEditFlagsHEX) != 0 && (flags&ColorEditFlagsNoInputs) == 0 {
		// RGB Hexadecimal Input
		var buf [64]byte
		var str string
		if alpha {
			str = fmt.Sprintf("#%02X%02X%02X%02X", mathutil.Clamp(i[0], 0, 255), mathutil.Clamp(i[1], 0, 255), mathutil.Clamp(i[2], 0, 255), mathutil.Clamp(i[3], 0, 255))
		} else {
			str = fmt.Sprintf("#%02X%02X%02X", mathutil.Clamp(i[0], 0, 255), mathutil.Clamp(i[1], 0, 255), mathutil.Clamp(i[2], 0, 255))
		}
		copy(buf[:], str)

		c.PushItemWidth(w_items_all)
		if c.InputTextEx("##Text", buf[:], f64.Vec2{0, 0}, InputTextFlagsCharsHexadecimal|InputTextFlagsCharsUppercase, nil) {
			value_changed = true
			p := buf[:]
			for len(p) > 0 && p[0] == '#' && CharIsSpace(rune(p[0])) {
				p = p[1:]
			}
			i[0], i[1], i[2], i[3] = 0, 0, 0, 0
			if alpha {
				fmt.Sscanf(string(p), "%02X%02X%02X%02X", &i[0], &i[1], &i[2], &i[3])
			} else {
				fmt.Sscanf(string(p), "%02X%02X%02X", &i[0], &i[1], &i[2])
			}
		}
		if flags&ColorEditFlagsNoOptions == 0 {
			c.OpenPopupOnItemClick("context", 1)
		}
		c.PopItemWidth()
	}

	var picker_active_window *Window
	if flags&ColorEditFlagsNoSmallPreview == 0 {
		if flags&ColorEditFlagsNoInputs == 0 {
			c.SameLineEx(0, style.ItemInnerSpacing.X)
		}

		if c.ColorButtonVEx("##ColorButton", *col, flags, f64.Vec2{0, 0}) {
			if flags&ColorEditFlagsNoPicker == 0 {
				// Store current color and open a picker
				c.ColorPickerRef = *col
				c.OpenPopup("picker")
				c.SetNextWindowPos(window.DC.LastItemRect.BL().Add(f64.Vec2{-1, style.ItemSpacing.Y}), 0, f64.Vec2{0, 0})
			}
		}
		if flags&ColorEditFlagsNoOptions == 0 {
			c.OpenPopupOnItemClick("context", 1)
		}

		if c.BeginPopup("picker") {
			picker_active_window = c.CurrentWindow
			if label != "" {
				c.TextUnformatted(label)
				c.Separator()
			}
			picker_flags_to_forward := ColorEditFlags_DataTypeMask | ColorEditFlags_PickerMask | ColorEditFlagsHDR | ColorEditFlagsNoAlpha | ColorEditFlagsAlphaBar
			picker_flags := (flags_untouched & picker_flags_to_forward) | ColorEditFlags_InputsMask | ColorEditFlagsNoLabel | ColorEditFlagsAlphaPreviewHalf
			c.PushItemWidth(square_sz * 12.0) // Use 256 + bar sizes?
			if c.ColorPickerV4Ex("##picker", col, picker_flags, &c.ColorPickerRef) {
				value_changed = true
			}
			c.PopItemWidth()
			c.EndPopup()
		}
	}

	if label != "" && flags&ColorEditFlagsNoLabel == 0 {
		c.SameLineEx(0, style.ItemInnerSpacing.X)
		c.TextUnformatted(label)
	}

	// Convert back
	if picker_active_window == nil {
		if !value_changed_as_float {
			for n := range i {
				f[n] = float64(i[n]) / 255.0
			}
		}
		if flags&ColorEditFlagsHSV != 0 {
			hsv := chroma.HSV2RGB(chroma.HSV{f[0], f[1], f[2]})
			vec := chroma.RGBA2VEC4(hsv)
			f[0], f[1], f[2], f[3] = vec.X, vec.Y, vec.Z, vec.W
		}
		if value_changed {
			*col = f64.Vec4{f[0], f[1], f[2], f[3]}
		}
	}

	c.PopID()
	c.EndGroup()

	// Drag and Drop Target
	// NB: The flag test is merely an optional micro-optimization, BeginDragDropTarget() does the same test.
	if window.DC.LastItemStatusFlags&ItemStatusFlagsHoveredRect != 0 && c.BeginDragDropTarget() {
		c.EndDragDropTarget()
	}

	// When picker is being actively used, use its active id so IsItemActive() will function on ColorEdit4().
	if picker_active_window != nil && c.ActiveId != 0 && c.ActiveIdWindow == picker_active_window {
		window.DC.LastItemId = c.ActiveId
	}

	return value_changed
}

func (c *Context) ColorPicker3(label string, col *color.RGBA) bool {
	colv := chroma.RGBA2VEC4(*col)
	value_changed := c.ColorPickerV3(label, &colv)
	*col = chroma.VEC42RGBA(colv)
	return value_changed
}

func (c *Context) ColorPicker4(label string, col *color.RGBA) bool {
	colv := chroma.RGBA2VEC4(*col)
	value_changed := c.ColorPickerV4(label, &colv)
	*col = chroma.VEC42RGBA(colv)
	return value_changed
}

func (c *Context) ColorPicker3Ex(label string, col *color.RGBA, flags ColorEditFlags) bool {
	colv := chroma.RGBA2VEC4(*col)
	value_changed := c.ColorPickerV3Ex(label, &colv, flags)
	*col = chroma.VEC42RGBA(colv)
	return value_changed
}

func (c *Context) ColorPicker4Ex(label string, col *color.RGBA, flags ColorEditFlags, ref_col *color.RGBA) bool {
	var ref_colv f64.Vec4
	var ref_colv_ptr *f64.Vec4
	if ref_col != nil {
		ref_colv = chroma.RGBA2VEC4(*col)
		ref_colv_ptr = &ref_colv
	}
	colv := chroma.RGBA2VEC4(*col)
	value_changed := c.ColorPickerV4Ex(label, &colv, flags, ref_colv_ptr)
	*col = chroma.VEC42RGBA(colv)
	if ref_colv_ptr != nil {
		*ref_col = chroma.VEC42RGBA(ref_colv)
	}
	return value_changed
}

func (c *Context) ColorPickerV3(label string, col *f64.Vec4) bool {
	return c.ColorPickerV3Ex(label, col, 0)
}

func (c *Context) ColorPickerV4(label string, col *f64.Vec4) bool {
	return c.ColorPickerV4Ex(label, col, 0, nil)
}

func (c *Context) ColorPickerV3Ex(label string, col *f64.Vec4, flags ColorEditFlags) bool {
	return c.ColorPickerV4Ex(label, col, flags|ColorEditFlagsNoAlpha, nil)
}

// ColorPicker
// Note: only access 3 floats if ImGuiColorEditFlags_NoAlpha flag is set.
// FIXME: we adjust the big color square height based on item width, which may cause a flickering feedback loop (if automatic height makes a vertical scrollbar appears, affecting automatic width..)
func (c *Context) ColorPickerV4Ex(label string, col *f64.Vec4, flags ColorEditFlags, ref_col *f64.Vec4) bool {
	window := c.GetCurrentWindow()
	draw_list := window.DrawList
	style := &c.Style
	io := &c.IO

	c.PushStringID(label)
	c.BeginGroup()
	if flags&ColorEditFlagsNoSidePreview == 0 {
		flags |= ColorEditFlagsNoSmallPreview
	}

	// Context menu: display and store options.
	if flags&ColorEditFlagsNoOptions == 0 {
		c.ColorPickerOptionsPopupV(flags, col)
	}

	// Read stored options
	if flags&ColorEditFlags_PickerMask == 0 {
		if c.ColorEditOptions&ColorEditFlags_PickerMask != 0 {
			flags |= c.ColorEditOptions & ColorEditFlags_PickerMask
		} else {
			flags |= ColorEditFlags_OptionsDefault & ColorEditFlags_PickerMask
		}
	}
	// Check that only 1 is selected
	assert(mathutil.IsPow2(int(flags & ColorEditFlags_PickerMask)))
	if flags&ColorEditFlagsNoOptions == 0 {
		flags |= (c.ColorEditOptions & ColorEditFlagsAlphaBar)
	}

	// Setup
	alpha_bar := flags&ColorEditFlagsAlphaBar != 0 && flags&ColorEditFlagsNoAlpha == 0
	picker_pos := window.DC.CursorPos
	square_sz := c.GetFrameHeight()
	// Arbitrary smallish width of Hue/Alpha picking bars
	bars_width := square_sz
	// Saturation/Value picking box
	alpha_width := 1.0
	if alpha_bar {
		alpha_width = 2
	}
	sv_picker_size := math.Max(bars_width*1, c.CalcItemWidth()-alpha_width*(bars_width+style.ItemInnerSpacing.X))
	bar0_pos_x := picker_pos.X + sv_picker_size + style.ItemInnerSpacing.X
	bar1_pos_x := bar0_pos_x + bars_width + style.ItemInnerSpacing.X
	bars_triangles_half_sz := float64(int(bars_width * 0.20))

	backup_initial_col := *col

	wheel_thickness := sv_picker_size * 0.08
	wheel_r_outer := sv_picker_size * 0.50
	wheel_r_inner := wheel_r_outer - wheel_thickness
	wheel_center := f64.Vec2{
		picker_pos.X + (sv_picker_size+bars_width)*0.5,
		picker_pos.Y + sv_picker_size*0.5,
	}

	// Note: the triangle is displayed rotated with triangle_pa pointing to Hue, but most coordinates stays unrotated for logic.
	triangle_r := wheel_r_inner - float64(int(sv_picker_size*0.027))
	triangle_pa := f64.Vec2{triangle_r, 0.0}                           // Hue point.
	triangle_pb := f64.Vec2{triangle_r * -0.5, triangle_r * -0.866025} // Black point.
	triangle_pc := f64.Vec2{triangle_r * -0.5, triangle_r * +0.866025} // White point.

	hsv := chroma.VEC42HSV(*col)
	value_changed := false
	value_changed_h := false
	value_changed_sv := false

	c.PushItemFlag(ItemFlagsNoNav, true)
	if flags&ColorEditFlagsPickerHueWheel != 0 {
		// Hue wheel + SV triangle logic
		c.InvisibleButton("hsv", f64.Vec2{sv_picker_size + style.ItemInnerSpacing.X + bars_width, sv_picker_size})
		if c.IsItemActive() {
			initial_off := c.IO.MouseClickedPos[0].Sub(wheel_center)
			current_off := c.IO.MousePos.Sub(wheel_center)
			initial_dist2 := initial_off.LenSquared()
			if initial_dist2 >= (wheel_r_inner-1)*(wheel_r_inner-1) && initial_dist2 <= (wheel_r_outer+1)*(wheel_r_outer+1) {
				// Interactive with Hue wheel
				hsv.H = math.Atan2(current_off.Y, current_off.X) / math.Pi * 0.5
				if hsv.H < 0 {
					hsv.H += 1
				}
				value_changed = true
				value_changed_h = true
			}
			cos_hue_angle := math.Cos(-hsv.H * 2.0 * math.Pi)
			sin_hue_angle := math.Sin(-hsv.H * 2.0 * math.Pi)
			if TriangleContainsPoint(triangle_pa, triangle_pb, triangle_pc, Rotate(initial_off, cos_hue_angle, sin_hue_angle)) {
				// Interacting with SV triangle
				current_off_unrotated := Rotate(current_off, cos_hue_angle, sin_hue_angle)
				if !TriangleContainsPoint(triangle_pa, triangle_pb, triangle_pc, current_off_unrotated) {
					current_off_unrotated = TriangleClosestPoint(triangle_pa, triangle_pb, triangle_pc, current_off_unrotated)
				}
				uu, vv, _ := TriangleBarycentricCoords(triangle_pa, triangle_pb, triangle_pc, current_off_unrotated)
				hsv.V = f64.Clamp(1.0-vv, 0.0001, 1.0)
				hsv.S = f64.Clamp(uu/hsv.V, 0.0001, 1.0)
				value_changed = true
				value_changed_sv = true
			}
		}

		if flags&ColorEditFlagsNoOptions == 0 {
			c.OpenPopupOnItemClick("context", 1)
		}
	} else if flags&ColorEditFlagsPickerHueBar != 0 {
		// SV rectangle logic
		c.InvisibleButton("sv", f64.Vec2{sv_picker_size, sv_picker_size})
		if c.IsItemActive() {
			hsv.S = f64.Saturate((io.MousePos.X - picker_pos.X) / (sv_picker_size - 1))
			hsv.V = 1.0 - f64.Saturate((io.MousePos.Y-picker_pos.Y)/(sv_picker_size-1))
			value_changed = true
			value_changed_sv = true
		}
		if flags&ColorEditFlagsNoOptions == 0 {
			c.OpenPopupOnItemClick("context", 1)
		}

		// Hue bar logic
		c.SetCursorScreenPos(f64.Vec2{bar0_pos_x, picker_pos.Y})
		c.InvisibleButton("hue", f64.Vec2{bars_width, sv_picker_size})
		if c.IsItemActive() {
			hsv.H = f64.Saturate((io.MousePos.Y - picker_pos.Y) / (sv_picker_size - 1))
			value_changed = true
			value_changed_h = true
		}
	}

	// Alpha bar logic
	if alpha_bar {
		c.SetCursorScreenPos(f64.Vec2{bar1_pos_x, picker_pos.Y})
		c.InvisibleButton("alpha", f64.Vec2{bars_width, sv_picker_size})
		if c.IsItemActive() {
			col.W = 1.0 - f64.Saturate((io.MousePos.Y-picker_pos.Y)/(sv_picker_size-1))
			value_changed = true
		}
	}
	// ImGuiItemFlags_NoNav
	c.PopItemFlag()

	if flags&ColorEditFlagsNoSidePreview == 0 {
		c.SameLineEx(0, style.ItemInnerSpacing.X)
		c.BeginGroup()
	}

	if flags&ColorEditFlagsNoLabel == 0 {
		label_display_end := c.FindRenderedTextEnd(label)
		if text := label[:label_display_end]; text != "" {
			if flags&ColorEditFlagsNoSidePreview != 0 {
				c.SameLineEx(0, style.ItemInnerSpacing.X)
			}
			c.TextUnformatted(text)
		}
	}

	if flags&ColorEditFlagsNoSidePreview == 0 {
		c.PushItemFlag(ItemFlagsNoNavDefaultFocus, true)
		if flags&ColorEditFlagsNoLabel != 0 {
			c.Text("Current")
		}
		col_v4 := *col
		if flags&ColorEditFlagsNoAlpha != 0 {
			col_v4.W = 1
		}
		c.ColorButtonVEx("##current", col_v4, (flags & (ColorEditFlagsHDR | ColorEditFlagsAlphaPreview | ColorEditFlagsAlphaPreviewHalf | ColorEditFlagsNoTooltip)), f64.Vec2{square_sz * 3, square_sz * 2})
		if ref_col != nil {
			c.Text("Original")
			ref_col_v4 := *ref_col
			if flags&ColorEditFlagsNoAlpha != 0 {
				ref_col_v4.W = 1
			}
			if c.ColorButtonVEx("##original", ref_col_v4, (flags & (ColorEditFlagsHDR | ColorEditFlagsAlphaPreview | ColorEditFlagsAlphaPreviewHalf | ColorEditFlagsNoTooltip)), f64.Vec2{square_sz * 3, square_sz * 2}) {
				*col = *ref_col
				value_changed = true
			}
		}
		c.PopItemFlag()
		c.EndGroup()
	}

	// Convert back color to RGB
	if value_changed_h || value_changed_sv {
		h, s, v := hsv.H, hsv.S, hsv.V
		if h >= 1.0 {
			h -= (10 * 1e-6)
		}
		if s < 0 {
			s = 10 * 1e-6
		}
		if v < 0 {
			v = 1e-6
		}
		rgb := chroma.HSV2VEC4(chroma.HSV{h, s, v})
		*col = f64.Vec4{rgb.X, rgb.Y, rgb.Z, col.W}
	}

	// R,G,B and H,S,V slider color editor
	if flags&ColorEditFlagsNoInputs == 0 {
		bar_x := bar0_pos_x
		if alpha_bar {
			bar_x = bar1_pos_x
		}
		c.PushItemWidth(bar_x + bars_width - picker_pos.X)
		sub_flags_to_forward := ColorEditFlags_DataTypeMask | ColorEditFlagsHDR | ColorEditFlagsNoAlpha | ColorEditFlagsNoOptions | ColorEditFlagsNoSmallPreview | ColorEditFlagsAlphaPreview | ColorEditFlagsAlphaPreviewHalf
		sub_flags := (flags & sub_flags_to_forward) | ColorEditFlagsNoPicker
		if flags&ColorEditFlagsRGB != 0 || flags&ColorEditFlags_InputsMask == 0 {
			if c.ColorEditV4Ex("##rgb", col, sub_flags|ColorEditFlagsRGB) {
				value_changed = true
			}
		}
		if flags&ColorEditFlagsHSV != 0 || flags&ColorEditFlags_InputsMask == 0 {
			if c.ColorEditV4Ex("##hsv", col, sub_flags|ColorEditFlagsHSV) {
				value_changed = true
			}
		}
		if flags&ColorEditFlagsHEX != 0 || flags&ColorEditFlags_InputsMask == 0 {
			if c.ColorEditV4Ex("##hex", col, sub_flags|ColorEditFlagsHEX) {
				value_changed = true
			}
		}
		c.PopItemWidth()
	}

	// Try to cancel hue wrap (after ColorEdit), if any
	if value_changed {
		new_hsv := chroma.VEC42HSV(*col)
		if new_hsv.V <= 0 && hsv.H > 0 {
			if new_hsv.V <= 0 && hsv.V != new_hsv.V {
				if new_hsv.V <= 0 {
					new_hsv.V = hsv.V * 0.5
				}
			} else if new_hsv.S <= 0 {
				if new_hsv.S <= 0 {
					new_hsv.S = hsv.S * 0.5
				}
			}
			*col = chroma.HSV2VEC4(new_hsv)
		}
	}

	hue_color32v := chroma.HSV2VEC4(chroma.HSV{hsv.H, 1, 1})
	hue_color32 := chroma.VEC42RGBA(hue_color32v)
	col32_no_alpha := *col
	col32_no_alpha.W = 1

	hue_colors := [...]color.RGBA{
		color.RGBA{255, 0, 0, 255},
		color.RGBA{255, 255, 0, 255},
		color.RGBA{0, 255, 0, 255},
		color.RGBA{0, 255, 255, 255},
		color.RGBA{0, 0, 255, 255},
		color.RGBA{255, 0, 255, 255},
		color.RGBA{255, 0, 0, 255},
	}
	var sv_cursor_pos f64.Vec2

	white := color.RGBA{255, 255, 255, 255}
	black := color.RGBA{0, 0, 0, 255}
	black_trans := color.RGBA{0, 0, 0, 0}
	// Render cursor/preview circle (clamp S/V within 0..1 range because floating points colors may lead HSV values to be out of range)
	if flags&ColorEditFlagsPickerHueWheel != 0 {
		// Render Hue Wheel

		// Half a pixel arc length in radians (2pi cancels out).
		aeps := 1.5 / wheel_r_outer
		segment_per_arc := mathutil.Max(4, int(wheel_r_outer)/12)
		for n := 0; n < 6; n++ {
			a0 := float64(n)/6.0*2.0*math.Pi - aeps
			a1 := float64(n+1.0)/6.0*2.0*math.Pi + aeps
			vert_start_idx := len(draw_list.VtxBuffer)
			draw_list.PathArcTo(wheel_center, (wheel_r_inner+wheel_r_outer)*0.5, a0, a1, segment_per_arc)
			draw_list.PathStrokeEx(color.RGBA{255, 255, 255, 255}, false, wheel_thickness)
			vert_end_idx := len(draw_list.VtxBuffer)

			// Paint colors over existing vertices
			gradient_p0 := f64.Vec2{wheel_center.X + math.Cos(a0)*wheel_r_inner, wheel_center.Y + math.Sin(a0)*wheel_r_inner}
			gradient_p1 := f64.Vec2{wheel_center.X + math.Cos(a1)*wheel_r_inner, wheel_center.Y + math.Sin(a1)*wheel_r_inner}
			c.ShadeVertsLinearColorGradientKeepAlpha(draw_list.VtxBuffer, vert_start_idx, vert_end_idx, gradient_p0, gradient_p1, hue_colors[n], hue_colors[n+1])
		}

		// Render Cursor + preview on Hue Wheel
		cos_hue_angle := math.Cos(hsv.H * 2.0 * math.Pi)
		sin_hue_angle := math.Sin(hsv.H * 2.0 * math.Pi)
		hue_cursor_pos := f64.Vec2{
			wheel_center.X + cos_hue_angle*(wheel_r_inner+wheel_r_outer)*0.5,
			wheel_center.Y + sin_hue_angle*(wheel_r_inner+wheel_r_outer)*0.5,
		}
		hue_cursor_rad := wheel_thickness * 0.55
		if value_changed_h {
			hue_cursor_rad = wheel_thickness * 0.65
		}
		hue_cursor_segments := mathutil.Clamp(int(hue_cursor_rad/1.4), 9, 32)

		draw_list.AddCircleFilledEx(hue_cursor_pos, hue_cursor_rad, hue_color32, hue_cursor_segments)
		draw_list.AddCircleEx(hue_cursor_pos, hue_cursor_rad+1, color.RGBA{128, 128, 128, 255}, hue_cursor_segments, 1.0)
		draw_list.AddCircleEx(hue_cursor_pos, hue_cursor_rad, color.RGBA{255, 255, 255, 255}, hue_cursor_segments, 1.0)

		// Render SV triangle (rotated according to hue)
		tra := wheel_center.Add(Rotate(triangle_pa, cos_hue_angle, sin_hue_angle))
		trb := wheel_center.Add(Rotate(triangle_pb, cos_hue_angle, sin_hue_angle))
		trc := wheel_center.Add(Rotate(triangle_pc, cos_hue_angle, sin_hue_angle))
		uv_white := c.GetFontTexUvWhitePixel()
		draw_list.PrimReserve(6, 6)
		draw_list.PrimVtx(tra, uv_white, hue_color32)
		draw_list.PrimVtx(trb, uv_white, hue_color32)
		draw_list.PrimVtx(trc, uv_white, white)
		draw_list.PrimVtx(tra, uv_white, black_trans)
		draw_list.PrimVtx(trb, uv_white, black)
		draw_list.PrimVtx(trc, uv_white, black_trans)
		draw_list.AddTriangleEx(tra, trb, trc, color.RGBA{128, 128, 128, 255}, 1.5)
		t0 := trc.Lerp(f64.Saturate(hsv.S), tra)
		t1 := t0.Lerp(f64.Saturate(1-hsv.V), trb)
		sv_cursor_pos = t1
	} else if flags&ColorEditFlagsPickerHueBar != 0 {
		// Render SV Square
		draw_list.AddRectFilledMultiColor(picker_pos, picker_pos.Add(f64.Vec2{sv_picker_size, sv_picker_size}), white, hue_color32, hue_color32, white)
		draw_list.AddRectFilledMultiColor(picker_pos, picker_pos.Add(f64.Vec2{sv_picker_size, sv_picker_size}), black_trans, black_trans, black, black)
		c.RenderFrameBorder(picker_pos, picker_pos.Add(f64.Vec2{sv_picker_size, sv_picker_size}), 0.0)
		// Sneakily prevent the circle to stick out too much
		sv_cursor_pos.X = f64.Clamp(float64(int(picker_pos.X+f64.Saturate(hsv.S)*sv_picker_size+0.5)), picker_pos.X+2, picker_pos.X+sv_picker_size-2)
		sv_cursor_pos.Y = f64.Clamp(float64(int(picker_pos.Y+f64.Saturate(1-hsv.V)*sv_picker_size+0.5)), picker_pos.Y+2, picker_pos.Y+sv_picker_size-2)

		// Render Hue Bar
		for i := 0; i < 6; i++ {
			draw_list.AddRectFilledMultiColor(
				f64.Vec2{bar0_pos_x, picker_pos.Y + float64(i)*(sv_picker_size/6)},
				f64.Vec2{bar0_pos_x + bars_width, picker_pos.Y + float64(i+1)*(sv_picker_size/6)},
				hue_colors[i], hue_colors[i], hue_colors[i+1], hue_colors[i+1],
			)
		}
		bar0_line_y := float64(int(picker_pos.Y + hsv.H*sv_picker_size + 0.5))
		c.RenderFrameBorder(f64.Vec2{bar0_pos_x, picker_pos.Y}, f64.Vec2{bar0_pos_x + bars_width, picker_pos.Y + sv_picker_size}, 0.0)
		c.RenderArrowsForVerticalBar(draw_list, f64.Vec2{bar0_pos_x - 1, bar0_line_y}, f64.Vec2{bars_triangles_half_sz + 1, bars_triangles_half_sz}, bars_width+2.0)
	}

	// Render cursor/preview circle (clamp S/V within 0..1 range because floating points colors may lead HSV values to be out of range)
	sv_cursor_rad := 6.0
	if value_changed_sv {
		sv_cursor_rad = 10
	}
	draw_list.AddCircleFilledEx(sv_cursor_pos, sv_cursor_rad, chroma.VEC42RGBA(col32_no_alpha), 12)
	draw_list.AddCircleEx(sv_cursor_pos, sv_cursor_rad+1, color.RGBA{128, 128, 128, 255}, 12, 1)
	draw_list.AddCircleEx(sv_cursor_pos, sv_cursor_rad, color.RGBA{255, 255, 255, 255}, 12, 1)

	// Render alpha bar
	if alpha_bar {
		alpha := col.W
		bar1_bb := f64.Rect(bar1_pos_x, picker_pos.Y, bar1_pos_x+bars_width, picker_pos.Y+sv_picker_size)
		c.RenderColorRectWithAlphaCheckerboard(bar1_bb.Min, bar1_bb.Max, color.RGBA{}, bar1_bb.Dx()/2.0, f64.Vec2{0.0, 0.0})
		col32_trans_alpha := col32_no_alpha
		col32_trans_alpha.W = 0

		rgba32_no_alpha := chroma.VEC42RGBA(col32_no_alpha)
		rgba32_trans_alpha := chroma.VEC42RGBA(col32_trans_alpha)

		draw_list.AddRectFilledMultiColor(bar1_bb.Min, bar1_bb.Max, rgba32_no_alpha, rgba32_no_alpha, rgba32_trans_alpha, rgba32_trans_alpha)
		bar1_line_y := float64(int(picker_pos.Y + (1.0-alpha)*sv_picker_size + 0.5))
		c.RenderFrameBorder(bar1_bb.Min, bar1_bb.Max, 0.0)
		c.RenderArrowsForVerticalBar(draw_list, f64.Vec2{bar1_pos_x - 1, bar1_line_y}, f64.Vec2{bars_triangles_half_sz + 1, bars_triangles_half_sz}, bars_width+2.0)
	}

	c.EndGroup()
	c.PopID()
	return value_changed && backup_initial_col != *col
}

func (c *Context) ColorPickerOptionsPopup(flags ColorEditFlags, ref_col *color.RGBA) {
	ref_colv := chroma.RGBA2VEC4(*ref_col)
	c.ColorPickerOptionsPopupV(flags, &ref_colv)
	*ref_col = chroma.VEC42RGBA(ref_colv)
}

func (c *Context) ColorPickerOptionsPopupV(flags ColorEditFlags, ref_col *f64.Vec4) {
	allow_opt_picker := flags&ColorEditFlags_PickerMask == 0
	allow_opt_alpha_bar := flags&ColorEditFlagsNoAlpha == 0 && flags&ColorEditFlagsAlphaBar == 0
	if (!allow_opt_picker && !allow_opt_alpha_bar) || !c.BeginPopup("context") {
		return
	}
	if allow_opt_picker {
		// FIXME: Picker size copied from main picker function
		picker_size := f64.Vec2{
			c.FontSize * 8,
			math.Max(c.FontSize*8-(c.GetFrameHeight()+c.Style.ItemInnerSpacing.X), 1.0),
		}
		c.PushItemWidth(picker_size.X)
		for picker_type := 0; picker_type < 2; picker_type++ {
			// Draw small/thumbnail version of each picker type (over an invisible button for selection)
			if picker_type > 0 {
				c.Separator()
			}
			c.PushID(ID(picker_type))
			picker_flags := ColorEditFlagsNoInputs | ColorEditFlagsNoOptions | ColorEditFlagsNoLabel | ColorEditFlagsNoSidePreview | (flags & ColorEditFlagsNoAlpha)
			if picker_type == 0 {
				picker_flags |= ColorEditFlagsPickerHueBar
			}
			if picker_type == 1 {
				picker_flags |= ColorEditFlagsPickerHueWheel
			}
			backup_pos := c.GetCursorScreenPos()
			// By default, Selectable() is closing popup
			if c.SelectableEx("##selectable", false, 0, picker_size) {
				c.ColorEditOptions = (c.ColorEditOptions &^ ColorEditFlags_PickerMask) | (picker_flags & ColorEditFlags_PickerMask)
			}
			c.SetCursorScreenPos(backup_pos)
			dummy_ref_col := *ref_col
			c.ColorPickerV4Ex("##dummypicker", &dummy_ref_col, picker_flags, nil)
			c.PopID()
		}
		c.PopItemWidth()
	}
	if allow_opt_alpha_bar {
		if allow_opt_picker {
			c.Separator()
		}
		colorEditOptions := uint(c.ColorEditOptions)
		c.CheckboxFlags("Alpha Bar", &colorEditOptions, uint(ColorEditFlagsAlphaBar))
		c.ColorEditOptions = ColorEditFlags(colorEditOptions)
	}
	c.EndPopup()
}

func (c *Context) ColorEditOptionsPopup(col color.RGBA, flags ColorEditFlags) {
	c.ColorEditOptionsPopupV(chroma.RGBA2VEC4(col), flags)
}

func (c *Context) ColorEditOptionsPopupV(col f64.Vec4, flags ColorEditFlags) {
	allow_opt_inputs := flags&ColorEditFlags_InputsMask == 0
	allow_opt_datatype := flags&ColorEditFlags_DataTypeMask == 0
	if (!allow_opt_inputs && !allow_opt_datatype) || !c.BeginPopup("context") {
		return
	}

	opts := c.ColorEditOptions
	if allow_opt_inputs {
		if c.RadioButton("RGB", opts&ColorEditFlagsRGB != 0) {
			opts = (opts &^ ColorEditFlags_InputsMask) | ColorEditFlagsRGB
		}
		if c.RadioButton("HSV", opts&ColorEditFlagsHSV != 0) {
			opts = (opts &^ ColorEditFlags_InputsMask) | ColorEditFlagsHSV
		}
		if c.RadioButton("HEX", opts&ColorEditFlagsHEX != 0) {
			opts = (opts &^ ColorEditFlags_InputsMask) | ColorEditFlagsHEX
		}
	}
	if allow_opt_datatype {
		if allow_opt_inputs {
			c.Separator()
		}
		if c.RadioButton("0..255", opts&ColorEditFlagsUint8 != 0) {
			opts = (opts &^ ColorEditFlags_DataTypeMask) | ColorEditFlagsUint8
		}
		if c.RadioButton("0.00..1.00", opts&ColorEditFlagsFloat != 0) {
			opts = (opts &^ ColorEditFlags_DataTypeMask) | ColorEditFlagsFloat
		}
	}

	if allow_opt_inputs || allow_opt_datatype {
		c.Separator()
	}
	if c.ButtonEx("Copy as..", f64.Vec2{-1, 0}, 0) {
		c.OpenPopup("Copy")
	}
	if c.BeginPopup("Copy") {
		if flags&ColorEditFlagsNoAlpha != 0 {
			col.W = 1
		}

		rgba := chroma.VEC42RGBA(col)
		cr, cg, cb, ca := rgba.R, rgba.G, rgba.B, rgba.A
		fr, fg, fb, fa := col.X, col.Y, col.Z, col.W

		buf := fmt.Sprintf("(%.3ff, %.3ff, %.3ff, %.3ff)", fr, fg, fb, fa)
		if c.Selectable(buf) {
			c.SetClipboardText(buf)
		}
		buf = fmt.Sprintf("(%d,%d,%d,%d)", cr, cg, cb, ca)
		if c.Selectable(buf) {
			c.SetClipboardText(buf)
		}
		if flags&ColorEditFlagsNoAlpha != 0 {
			buf = fmt.Sprintf("0x%02X%02X%02X", cr, cg, cb, ca)
		} else {
			buf = fmt.Sprintf("0x%02X%02X%02X%02X", cr, cg, cb, ca)
		}
		if c.Selectable(buf) {
			c.SetClipboardText(buf)
		}
		c.EndPopup()
	}
	c.ColorEditOptions = opts
	c.EndPopup()
}

func (c *Context) ColorEdit3(label string, col *color.RGBA) bool {
	return c.ColorEdit4Ex(label, col, ColorEditFlagsNoAlpha)
}

func (c *Context) ColorEdit3Ex(label string, col *color.RGBA, flags ColorEditFlags) bool {
	return c.ColorEdit4Ex(label, col, flags|ColorEditFlagsNoAlpha)
}

func (c *Context) ColorEdit4(label string, col *color.RGBA) bool {
	return c.ColorEdit4Ex(label, col, 0)
}

func (c *Context) ColorEdit4Ex(label string, col *color.RGBA, flags ColorEditFlags) bool {
	vec := chroma.RGBA2VEC4(*col)
	value_changed := c.ColorEditV4Ex(label, &vec, flags)
	*col = chroma.VEC42RGBA(vec)
	return value_changed
}
