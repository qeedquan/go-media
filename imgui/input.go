package imgui

import (
	"math"

	"github.com/qeedquan/go-media/math/f64"
)

type Key int

const (
	KeyTab Key = iota
	KeyLeftArrow
	KeyRightArrow
	KeyUpArrow
	KeyDownArrow
	KeyPageUp
	KeyPageDown
	KeyHome
	KeyEnd
	KeyInsert
	KeyDelete
	KeyBackspace
	KeySpace
	KeyEnter
	KeyEscape
	KeyA // for text edit CTRL+A: select all
	KeyC // for text edit CTRL+C: copy
	KeyV // for text edit CTRL+V: paste
	KeyX // for text edit CTRL+X: cut
	KeyY // for text edit CTRL+Y: redo
	KeyZ // for text edit CTRL+Z: undo
	KeyCOUNT
)

type InputSource int

const (
	InputSourceNone InputSource = iota
	InputSourceMouse
	InputSourceNav
	InputSourceNavKeyboard // Only used occasionally for storage, not tested/handled by most code
	InputSourceNavGamepad  // "
	InputSourceCOUNT
)

type InputReadMode int

const (
	InputReadModeDown InputReadMode = iota
	InputReadModePressed
	InputReadModeReleased
	InputReadModeRepeat
	InputReadModeRepeatSlow
	InputReadModeRepeatFast
)

type MouseCursor int

const (
	MouseCursorNone MouseCursor = -1 + iota
	MouseCursorArrow
	MouseCursorTextInput  // When hovering over InputText, etc.
	MouseCursorResizeAll  // Unused
	MouseCursorResizeNS   // When hovering over an horizontal border
	MouseCursorResizeEW   // When hovering over a vertical border or a column
	MouseCursorResizeNESW // When hovering over the bottom-left corner of a window
	MouseCursorResizeNWSE // When hovering over the bottom-right corner of a window
	MouseCursorCOUNT
)

func (c *Context) CalcTypematicPressedRepeatAmount(t, t_prev, repeat_delay, repeat_rate float64) int {
	if t == 0 {
		return 1
	}
	if t <= repeat_delay || repeat_rate <= 0 {
		return 0
	}

	count := int((t-repeat_delay)/repeat_rate) - int((t_prev-repeat_delay)/repeat_rate)

	if count > 0 {
		return count
	}
	return 0
}

// Note that imgui doesn't know the semantic of each entry of io.KeysDown[]. Use your own indices/enums according to how your back-end/engine stored them into io.KeysDown[]!
func (c *Context) IsKeyDown(user_key_index int) bool {
	if user_key_index < 0 {
		return false
	}
	return c.IO.KeysDown[user_key_index]
}

func (c *Context) IsKeyPressed(user_key_index int, repeat bool) bool {
	if user_key_index < 0 {
		return false
	}
	t := c.IO.KeysDownDuration[user_key_index]
	if t == 0 {
		return true
	}
	if repeat && t > c.IO.KeyRepeatDelay {
		return c.GetKeyPressedAmount(user_key_index, c.IO.KeyRepeatDelay, c.IO.KeyRepeatRate) > 0
	}
	return false
}

func (c *Context) IsKeyReleased(user_key_index int) bool {
	if user_key_index < 0 {
		return false
	}
	return c.IO.KeysDownDurationPrev[user_key_index] >= 0 && !c.IO.KeysDown[user_key_index]
}

func (c *Context) IsMouseDown(button int) bool {
	return c.IO.MouseDown[button]
}

func (c *Context) GetKeyPressedAmount(key_index int, repeat_delay, repeat_rate float64) int {
	if key_index < 0 {
		return 0
	}
	t := c.IO.KeysDownDuration[key_index]
	return c.CalcTypematicPressedRepeatAmount(t, t-c.IO.DeltaTime, repeat_delay, repeat_rate)
}

func (c *Context) PushAllowKeyboardFocus(allow_keyboard_focus bool) {
	c.PushItemFlag(ItemFlagsAllowKeyboardFocus, allow_keyboard_focus)
}

func (c *Context) PopAllowKeyboardFocus() {
	c.PopItemFlag()
}

func (c *Context) IsKeyPressedMap(key Key) bool {
	return c.IsKeyPressedMapEx(key, true)
}

func (c *Context) IsKeyPressedMapEx(key Key, repeat bool) bool {
	key_index := c.IO.KeyMap[key]
	if key_index >= 0 {
		return c.IsKeyPressed(key_index, repeat)
	}
	return false
}

func (c *Context) InputInt(label string, v *int) bool {
	return c.InputIntEx(label, v, 1, 100, 0)
}

func (c *Context) InputIntEx(label string, v *int, step, step_fast int, extra_flags InputTextFlags) bool {
	var step_ptr, step_fast_ptr *int
	if step > 0 {
		step_ptr = &step
	}
	if step_fast > 0 {
		step_fast_ptr = &step_fast
	}

	// Hexadecimal input provided as a convenience but the flag name is awkward. Typically you'd use InputText() to parse your own data, if you want to handle prefixes.
	format := "%d"
	if extra_flags&InputTextFlagsCharsHexadecimal != 0 {
		format = "%08X"
	}
	return c.InputScalarEx(label, v, step_ptr, step_fast_ptr, format, extra_flags)
}

func (c *Context) InputIntN(label string, v []int, extra_flags InputTextFlags) bool {
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
		if c.InputIntEx("##v", &v[i], 0, 0, extra_flags) {
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

func (c *Context) InputFloat(label string, v *float64, step float64) bool {
	return c.InputFloatEx(label, v, step, 0, "", 0)
}

func (c *Context) InputFloatEx(label string, v *float64, step, step_fast float64, format string, extra_flags InputTextFlags) bool {
	var step_ptr, step_fast_ptr *float64
	if step > 0 {
		step_ptr = &step
	}
	if step_fast > 0 {
		step_fast_ptr = &step_fast
	}
	extra_flags |= InputTextFlagsCharsScientific
	return c.InputScalarEx(label, v, step_ptr, step_fast_ptr, format, extra_flags)
}

// NB: scalar_format here must be a simple "%xx" format string with no prefix/suffix (unlike the Drag/Slider functions "format" argument)
func (c *Context) InputScalarEx(label string, data, step_ptr, step_fast_ptr interface{}, scalar_format string, extra_flags InputTextFlags) bool {
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	style := &c.Style

	buf := DataTypeFormatString(data, scalar_format)

	value_changed := false
	if extra_flags&(InputTextFlagsCharsHexadecimal|InputTextFlagsCharsScientific) == 0 {
		extra_flags |= InputTextFlagsCharsDecimal
	}
	extra_flags |= InputTextFlagsAutoSelectAll

	if step_ptr != nil {
		button_size := c.GetFrameHeight()

		// The only purpose of the group here is to allow the caller to query item data e.g. IsItemActive()
		c.BeginGroup()
		c.PushStringID(label)
		c.PushItemWidth(math.Max(1.0, c.CalcItemWidth()-(button_size+style.ItemInnerSpacing.X)*2))
		// PushId(label) + "" gives us the expected ID from outside point of view
		if c.InputTextEx("", buf, f64.Vec2{0, 0}, extra_flags, nil) {
			value_changed = DataTypeApplyOpFromText(buf, string(c.InputTextState.InitialText), data, scalar_format)
		}
		c.PopItemWidth()

		// Step buttons
		c.SameLineEx(0, style.ItemInnerSpacing.X)
		if c.ButtonEx("-", f64.Vec2{button_size, button_size}, ButtonFlagsRepeat|ButtonFlagsDontClosePopups) {
			step := step_fast_ptr
			if c.IO.KeyCtrl && step_fast_ptr != nil {
				step = step_ptr
			}
			DataTypeApplyOp('-', data, data, step)
			value_changed = true
		}
		c.SameLineEx(0, style.ItemInnerSpacing.X)
		if c.ButtonEx("+", f64.Vec2{button_size, button_size}, ButtonFlagsRepeat|ButtonFlagsDontClosePopups) {
			step := step_fast_ptr
			if c.IO.KeyCtrl && step_fast_ptr != nil {
				step = step_ptr
			}
			DataTypeApplyOp('+', data, data, step)
			value_changed = true
		}
		c.SameLineEx(0, style.ItemInnerSpacing.X)
		n := c.FindRenderedTextEnd(label)
		c.TextUnformatted(label[:n])

		c.PopID()
		c.EndGroup()
	} else {
		if c.InputTextEx(label, buf, f64.Vec2{0, 0}, extra_flags, nil) {
			value_changed = DataTypeApplyOpFromText(buf, string(c.InputTextState.InitialText), data, scalar_format)
		}
	}

	return value_changed
}

func (c *Context) InputFloatN(label string, v []float64) bool {
	return c.InputFloatNEx(label, v, "", 0)
}

func (c *Context) InputFloatNEx(label string, v []float64, format string, extra_flags InputTextFlags) bool {
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
		if c.InputFloatEx("##v", &v[i], 0, 0, format, extra_flags) {
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

func (c *Context) InputV2(label string, v *f64.Vec2) bool {
	return c.InputV2Ex(label, v, "%.3f", 0)
}

func (c *Context) InputV2Ex(label string, v *f64.Vec2, format string, extra_flags InputTextFlags) bool {
	f := [...]float64{v.X, v.Y}
	r := c.InputFloatNEx(label, f[:2], format, extra_flags)
	v.X, v.Y = f[0], f[1]
	return r
}

func (c *Context) InputV3(label string, v *f64.Vec3) bool {
	return c.InputV3Ex(label, v, "%.3f", 0)
}

func (c *Context) InputV3Ex(label string, v *f64.Vec3, format string, extra_flags InputTextFlags) bool {
	f := [...]float64{v.X, v.Y, v.Z}
	r := c.InputFloatNEx(label, f[:3], format, extra_flags)
	v.X, v.Y, v.Z = f[0], f[1], f[2]
	return r
}

func (c *Context) InputV4(label string, v *f64.Vec4) bool {
	return c.InputV4Ex(label, v, "%.3f", 0)
}

func (c *Context) InputV4Ex(label string, v *f64.Vec4, format string, extra_flags InputTextFlags) bool {
	f := [...]float64{v.X, v.Y, v.Z, v.W}
	r := c.InputFloatNEx(label, f[:4], format, extra_flags)
	v.X, v.Y, v.Z, v.W = f[0], f[1], f[2], f[3]
	return r
}

func (c *Context) InputFloat2(label string, v []float64) bool {
	return c.InputFloat2Ex(label, v, "%.3f", 0)
}

func (c *Context) InputFloat2Ex(label string, v []float64, format string, extra_flags InputTextFlags) bool {
	return c.InputFloatNEx(label, v[:2], format, extra_flags)
}

func (c *Context) InputFloat3(label string, v []float64) bool {
	return c.InputFloat3Ex(label, v, "%.3f", 0)
}

func (c *Context) InputFloat3Ex(label string, v []float64, format string, extra_flags InputTextFlags) bool {
	return c.InputFloatNEx(label, v[:3], format, extra_flags)
}

func (c *Context) InputFloat4(label string, v []float64) bool {
	return c.InputFloat4Ex(label, v, "%.3f", 0)
}

func (c *Context) InputFloat4Ex(label string, v []float64, format string, extra_flags InputTextFlags) bool {
	return c.InputFloatNEx(label, v[:4], format, extra_flags)
}

func (c *Context) InputInt2(label string, v []int) bool {
	return c.InputInt2Ex(label, v, 0)
}

func (c *Context) InputInt2Ex(label string, v []int, extra_flags InputTextFlags) bool {
	return c.InputIntN(label, v[:2], extra_flags)
}

func (c *Context) InputInt3(label string, v []int) bool {
	return c.InputInt3Ex(label, v, 0)
}

func (c *Context) InputInt3Ex(label string, v []int, extra_flags InputTextFlags) bool {
	return c.InputIntN(label, v[:3], extra_flags)
}

func (c *Context) InputInt4(label string, v []int) bool {
	return c.InputInt4Ex(label, v, 0)
}

func (c *Context) InputInt4Ex(label string, v []int, extra_flags InputTextFlags) bool {
	return c.InputIntN(label, v[:4], extra_flags)
}
