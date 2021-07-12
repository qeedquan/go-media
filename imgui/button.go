package imgui

import (
	"fmt"
	"image/color"
	"math"

	"github.com/qeedquan/go-media/image/chroma"
	"github.com/qeedquan/go-media/math/f64"
)

type ButtonFlags int

const (
	ButtonFlagsRepeat                ButtonFlags = 1 << 0  // hold to repeat
	ButtonFlagsPressedOnClickRelease ButtonFlags = 1 << 1  // return true on click + release on same item [DEFAULT if no PressedOn* flag is set]
	ButtonFlagsPressedOnClick        ButtonFlags = 1 << 2  // return true on click (default requires click+release)
	ButtonFlagsPressedOnRelease      ButtonFlags = 1 << 3  // return true on release (default requires click+release)
	ButtonFlagsPressedOnDoubleClick  ButtonFlags = 1 << 4  // return true on double-click (default requires click+release)
	ButtonFlagsFlattenChildren       ButtonFlags = 1 << 5  // allow interactions even if a child window is overlapping
	ButtonFlagsAllowItemOverlap      ButtonFlags = 1 << 6  // require previous frame HoveredId to either match id or be null before being usable use along with SetItemAllowOverlap()
	ButtonFlagsDontClosePopups       ButtonFlags = 1 << 7  // disable automatically closing parent popup on press // [UNUSED]
	ButtonFlagsDisabled              ButtonFlags = 1 << 8  // disable interactions
	ButtonFlagsAlignTextBaseLine     ButtonFlags = 1 << 9  // vertically align button to match text baseline - ButtonEx() only // FIXME: Should be removed and handled by SmallButton() not possible currently because of DC.CursorPosPrevLine
	ButtonFlagsNoKeyModifiers        ButtonFlags = 1 << 10 // disable interaction if a key modifier is held
	ButtonFlagsNoHoldingActiveID     ButtonFlags = 1 << 11 // don't set ActiveId while holding the mouse (ButtonFlagsPressedOnClick only)
	ButtonFlagsPressedOnDragDropHold ButtonFlags = 1 << 12 // press when held into while we are drag and dropping another item (used by e.g. tree nodes collapsing headers)
	ButtonFlagsNoNavFocus            ButtonFlags = 1 << 13 // don't override navigation focus when activated
)

func (c *Context) ColorButton(desc_id string, col color.RGBA) bool {
	return c.ColorButtonV(desc_id, chroma.RGBA2VEC4(col))
}

func (c *Context) ColorButtonEx(desc_id string, col color.RGBA, flags ColorEditFlags, size f64.Vec2) bool {
	return c.ColorButtonVEx(desc_id, chroma.RGBA2VEC4(col), flags, size)
}

func (c *Context) ColorButtonV(desc_id string, colv f64.Vec4) bool {
	return c.ColorButtonVEx(desc_id, colv, 0, f64.Vec2{0, 0})
}

func (c *Context) ColorButtonVEx(desc_id string, colv f64.Vec4, flags ColorEditFlags, size f64.Vec2) bool {
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	id := window.GetID(desc_id)
	default_size := c.GetFrameHeight()
	if size.X == 0.0 {
		size.X = default_size
	}
	if size.Y == 0.0 {
		size.Y = default_size
	}

	bb := f64.Rectangle{window.DC.CursorPos, window.DC.CursorPos.Add(size)}
	bb_size := 0.0
	if size.Y >= default_size {
		bb_size = c.Style.FramePadding.Y
	}
	c.ItemSizeBBEx(bb, bb_size)
	if !c.ItemAdd(bb, id) {
		return false
	}

	hovered, _, pressed := c.ButtonBehavior(bb, id, 0)
	if flags&ColorEditFlagsNoAlpha != 0 {
		flags &^= ColorEditFlagsAlphaPreview | ColorEditFlagsAlphaPreviewHalf
	}

	col := chroma.VEC42RGBA(colv)
	col_without_alpha := col
	col_without_alpha.A = 255
	grid_step := math.Min(size.X, size.Y) / 2.99
	rounding := math.Min(c.Style.FrameRounding, grid_step*0.5)
	bb_inner := bb

	// The border (using Col_FrameBg) tends to look off when color is near-opaque and rounding is enabled. This offset seemed like a good middle ground to reduce those artifacts.
	off := -0.75
	bb_inner = bb_inner.Expand(off, off)
	if flags&ColorEditFlagsAlphaPreviewHalf != 0 && col.A < 255 {
		mid_x := float64(int((bb_inner.Min.X+bb_inner.Max.X)*0.5 + 0.5))
		c.RenderColorRectWithAlphaCheckerboardEx(
			f64.Vec2{bb_inner.Min.X + grid_step, bb_inner.Min.Y},
			bb_inner.Max,
			col,
			grid_step,
			f64.Vec2{-grid_step + off, off},
			rounding, DrawCornerFlagsTopRight|DrawCornerFlagsBotRight,
		)
		window.DrawList.AddRectFilledEx(bb_inner.Min, f64.Vec2{mid_x, bb_inner.Max.Y}, col_without_alpha, rounding, DrawCornerFlagsTopLeft|DrawCornerFlagsBotLeft)
	} else {
		// Because GetColorU32() multiplies by the global style Alpha and we don't want to display a checkerboard if the source code had no alpha
		col_source := col_without_alpha
		if flags&ColorEditFlagsAlphaPreview != 0 {
			col_source = col
		}
		if col_source.A < 255 {
			c.RenderColorRectWithAlphaCheckerboardDx(bb_inner.Min, bb_inner.Max, col_source, grid_step, f64.Vec2{off, off}, rounding)
		} else {
			window.DrawList.AddRectFilledEx(bb_inner.Min, bb_inner.Max, col_source, rounding, DrawCornerFlagsAll)
		}
	}
	c.RenderNavHighlight(bb, id)
	if c.Style.FrameBorderSize > 0.0 {
		c.RenderFrameBorder(bb.Min, bb.Max, rounding)
	} else {
		// Color button are often in need of some sort of border
		window.DrawList.AddRectDx(bb.Min, bb.Max, c.GetColorFromStyle(ColFrameBg), rounding)
	}

	// Drag and Drop Source
	if c.ActiveId == id && c.BeginDragDropSource() {
		if flags&ColorEditFlagsNoAlpha != 0 {
			c.SetDragDropPayload(col)
		} else {
			c.SetDragDropPayload(col)
		}
		c.ColorButtonEx(desc_id, col, flags, f64.Vec2{})
		c.SameLine()
		c.TextUnformatted("Color")
		c.EndDragDropSource()
		hovered = false
	}

	// Tooltip
	if flags&ColorEditFlagsNoTooltip == 0 && hovered {
		c.ColorTooltip(desc_id, col, flags&(ColorEditFlagsNoAlpha|ColorEditFlagsAlphaPreview|ColorEditFlagsAlphaPreviewHalf))
	}

	return pressed
}

func (c *Context) Button(label string) bool {
	return c.ButtonEx(label, f64.Vec2{}, 0)
}

func (c *Context) SmallButton(label string) bool {
	backup_padding_y := c.Style.FramePadding.Y
	c.Style.FramePadding.Y = 0
	pressed := c.ButtonEx(label, f64.Vec2{0, 0}, ButtonFlagsAlignTextBaseLine)
	c.Style.FramePadding.Y = backup_padding_y
	return pressed
}

func (c *Context) ButtonEx(label string, size_arg f64.Vec2, flags ButtonFlags) bool {
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	style := &c.Style
	id := window.GetID(label)
	label_size := c.CalcTextSizeEx(label, true, -1)

	pos := window.DC.CursorPos
	if flags&ButtonFlagsAlignTextBaseLine != 0 && style.FramePadding.Y < window.DC.CurrentLineTextBaseOffset {
		pos.Y += window.DC.CurrentLineTextBaseOffset - style.FramePadding.Y
	}
	size := c.CalcItemSize(size_arg, label_size.X+style.FramePadding.X*2, label_size.Y+style.FramePadding.Y*2)

	bb := f64.Rectangle{pos, pos.Add(size)}
	c.ItemSizeBBEx(bb, style.FramePadding.Y)
	if !c.ItemAdd(bb, id) {
		return false
	}

	if window.DC.ItemFlags&ItemFlagsButtonRepeat != 0 {
		flags |= ButtonFlagsRepeat
	}

	hovered, held, pressed := c.ButtonBehavior(bb, id, flags)

	var col color.RGBA
	switch {
	case hovered && held:
		col = c.GetColorFromStyle(ColButtonActive)
	case hovered:
		col = c.GetColorFromStyle(ColButtonHovered)
	default:
		col = c.GetColorFromStyle(ColButton)
	}

	// Render
	c.RenderNavHighlight(bb, id)
	c.RenderFrameEx(bb.Min, bb.Max, col, true, style.FrameRounding)
	c.RenderTextClippedEx(
		bb.Min.Add(style.FramePadding),
		bb.Max.Sub(style.FramePadding),
		label,
		&label_size,
		style.ButtonTextAlign,
		&bb,
	)

	return pressed
}

func (c *Context) ButtonBehavior(bb f64.Rectangle, id ID, flags ButtonFlags) (hovered, held, pressed bool) {
	window := c.GetCurrentWindow()

	if flags&ButtonFlagsDisabled != 0 {
		if c.ActiveId == id {
			c.ClearActiveID()
		}
		return
	}

	// Default behavior requires click+release on same spot
	if flags&(ButtonFlagsPressedOnClickRelease|ButtonFlagsPressedOnClick|ButtonFlagsPressedOnRelease|ButtonFlagsPressedOnDoubleClick) == 0 {
		flags |= ButtonFlagsPressedOnClickRelease
	}

	backup_hovered_window := c.HoveredWindow
	if flags&ButtonFlagsFlattenChildren != 0 && c.HoveredRootWindow == window {
		c.HoveredWindow = window
	}

	hovered = c.ItemHoverable(bb, id)

	// Special mode for Drag and Drop where holding button pressed for a long time while dragging another item triggers the button
	if flags&ButtonFlagsPressedOnDragDropHold != 0 && c.DragDropActive && c.DragDropSourceFlags&DragDropFlagsSourceNoHoldToOpenOthers == 0 {
		if c.IsItemHoveredEx(HoveredFlagsAllowWhenBlockedByActiveItem) {
			hovered = true
			c.SetHoveredID(id)
			// FIXME: Our formula for CalcTypematicPressedRepeatAmount() is fishy
			if c.CalcTypematicPressedRepeatAmount(c.HoveredIdTimer+0.0001, c.HoveredIdTimer+0.0001-c.IO.DeltaTime, 0.01, 0.70) != 0 {
				pressed = true
				c.FocusWindow(window)
			}
		}
	}

	if flags&ButtonFlagsFlattenChildren != 0 && c.HoveredRootWindow == window {
		c.HoveredWindow = backup_hovered_window
	}

	// AllowOverlap mode (rarely used) requires previous frame HoveredId to be null or to match. This allows using patterns where a later submitted widget overlaps a previous one.
	if hovered && flags&ButtonFlagsAllowItemOverlap != 0 && (c.HoveredIdPreviousFrame != id && c.HoveredIdPreviousFrame != 0) {
		hovered = false
	}

	// Mouse
	if hovered {
		if flags&ButtonFlagsNoKeyModifiers == 0 || (!c.IO.KeyCtrl && !c.IO.KeyShift && !c.IO.KeyAlt) {
			//                        | CLICKING        | HOLDING with ImGuiButtonFlags_Repeat
			// PressedOnClickRelease  |  <on release>*  |  <on repeat> <on repeat> .. (NOT on release)  <-- MOST COMMON! (*) only if both click/release were over bounds
			// PressedOnClick         |  <on click>     |  <on click> <on repeat> <on repeat> ..
			// PressedOnRelease       |  <on release>   |  <on repeat> <on repeat> .. (NOT on release)
			// PressedOnDoubleClick   |  <on dclick>    |  <on dclick> <on repeat> <on repeat> ..
			// FIXME-NAV: We don't honor those different behaviors.
			if flags&ButtonFlagsPressedOnClickRelease != 0 && c.IO.MouseClicked[0] {
				c.SetActiveID(id, window)
				if flags&ButtonFlagsNoNavFocus == 0 {
					c.SetFocusID(id, window)
				}
				c.FocusWindow(window)
			}

			if (flags&ButtonFlagsPressedOnClick != 0 && c.IO.MouseClicked[0]) || (flags&ButtonFlagsPressedOnDoubleClick != 0 && c.IO.MouseDoubleClicked[0]) {
				pressed = true
				if flags&ButtonFlagsNoHoldingActiveID != 0 {
					c.ClearActiveID()
				} else {
					c.SetActiveID(id, window) // Hold on ID
				}
				c.FocusWindow(window)
			}

			if flags&ButtonFlagsPressedOnRelease != 0 && c.IO.MouseReleased[0] {
				// Repeat mode trumps <on release>
				if !(flags&ButtonFlagsRepeat == 0 && c.IO.MouseDownDurationPrev[0] >= c.IO.KeyRepeatDelay) {
					pressed = true
				}
				c.ClearActiveID()
			}

			// 'Repeat' mode acts when held regardless of _PressedOn flags (see table above).
			// Relies on repeat logic of IsMouseClicked() but we may as well do it ourselves if we end up exposing finer RepeatDelay/RepeatRate settings.
			if flags&ButtonFlagsRepeat != 0 && c.ActiveId == id && c.IO.MouseDownDuration[0] > 0 && c.IsMouseClicked(0, true) {
				pressed = true
			}
		}

		if pressed {
			c.NavDisableHighlight = true
		}
	}

	// Gamepad/Keyboard navigation
	// We report navigated item as hovered but we don't set g.HoveredId to not interfere with mouse.
	if c.NavId == id && !c.NavDisableHighlight && c.NavDisableMouseHover && (c.ActiveId == 0 || c.ActiveId == id || c.ActiveId == window.MoveId) {
		hovered = true
	}

	if c.NavActivateDownId == id {
		var nav_activated_by_inputs bool
		nav_activated_by_code := c.NavActivateId == id

		if flags&ButtonFlagsRepeat != 0 {
			nav_activated_by_inputs = c.IsNavInputPressed(NavInputActivate, InputReadModeRepeat)
		} else {
			nav_activated_by_inputs = c.IsNavInputPressed(NavInputActivate, InputReadModePressed)
		}

		if nav_activated_by_code || nav_activated_by_inputs {
			pressed = true
		}

		if nav_activated_by_code || nav_activated_by_inputs || c.ActiveId == id {
			// Set active id so it can be queried by user via IsItemActive(), equivalent of holding the mouse button.
			c.NavActivateId = id
			c.SetActiveID(id, window)
			if flags&ButtonFlagsNoNavFocus == 0 {
				c.SetFocusID(id, window)
			}
			c.ActiveIdAllowNavDirFlags = 1<<uint(DirLeft) | 1<<uint(DirRight) | 1<<uint(DirUp) | 1<<uint(DirDown)
		}
	}

	if c.ActiveId == id {
		if c.ActiveIdSource == InputSourceMouse {
			if c.ActiveIdIsJustActivated {
				c.ActiveIdClickOffset = c.IO.MousePos.Sub(bb.Min)
			}

			if c.IO.MouseDown[0] {
				held = true
			} else {
				if hovered && flags&ButtonFlagsPressedOnClickRelease != 0 {
					// Repeat mode trumps <on release>
					if !(flags&ButtonFlagsRepeat != 0 && c.IO.MouseDownDurationPrev[0] >= c.IO.KeyRepeatDelay) {
						if !c.DragDropActive {
							pressed = true
						}
					}
				}

				c.ClearActiveID()
			}

			if flags&ButtonFlagsNoNavFocus == 0 {
				c.NavDisableHighlight = true
			}
		} else if c.ActiveIdSource == InputSourceNav {
			if c.NavActivateDownId != id {
				c.ClearActiveID()
			}
		}
	}

	return
}

// Tip: use ImGui::PushID()/PopID() to push indices or pointers in the ID stack.
// Then you can keep 'str_id' empty or the same for all your buttons (instead of creating a string based on a non-string id)
func (c *Context) InvisibleButton(str_id string, size_arg f64.Vec2) bool {
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	id := window.GetID(str_id)
	size := c.CalcItemSize(size_arg, 0, 0)
	bb := f64.Rectangle{
		window.DC.CursorPos,
		window.DC.CursorPos.Add(size),
	}
	c.ItemSizeBB(bb)
	if !c.ItemAdd(bb, id) {
		return false
	}

	_, _, pressed := c.ButtonBehavior(bb, id, 0)
	return pressed
}

// Button to close a window
func (c *Context) CloseButton(id ID, pos f64.Vec2, radius float64) bool {
	window := c.CurrentWindow

	// We intentionally allow interaction when clipped so that a mechanical Alt,Right,Validate sequence close a window.
	// (this isn't the regular behavior of buttons, but it doesn't affect the user much because navigation tends to keep items visible).
	rad := f64.Vec2{radius, radius}
	bb := f64.Rectangle{pos.Sub(rad), pos.Add(rad)}
	is_clipped := !c.ItemAdd(bb, id)

	hovered, held, pressed := c.ButtonBehavior(bb, id, 0)
	if is_clipped {
		return pressed
	}

	// Render
	center := bb.Center()
	if hovered {
		var col color.RGBA
		switch {
		case held && hovered:
			col = c.GetColorFromStyle(ColButtonActive)
		default:
			col = c.GetColorFromStyle(ColButtonHovered)
		}
		window.DrawList.AddCircleFilledEx(center, math.Max(2, radius), col, 9)
	}

	cross_extent := (radius * 0.7071) - 1.0
	cross_col := c.GetColorFromStyle(ColText)
	center = center.Sub(f64.Vec2{0.5, 0.5})
	window.DrawList.AddLineEx(center.Add(f64.Vec2{+cross_extent, +cross_extent}), center.Add(f64.Vec2{-cross_extent, -cross_extent}), cross_col, 1.0)
	window.DrawList.AddLineEx(center.Add(f64.Vec2{+cross_extent, -cross_extent}), center.Add(f64.Vec2{-cross_extent, +cross_extent}), cross_col, 1.0)

	return pressed
}

func (c *Context) PushButtonRepeat(repeat bool) {
	c.PushItemFlag(ItemFlagsButtonRepeat, repeat)
}

func (c *Context) PopButtonRepeat() {
	c.PopItemFlag()
}

func (c *Context) RadioButton(label string, active bool) bool {
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	style := &c.Style
	id := window.GetID(label)
	label_size := c.CalcTextSizeEx(label, true, -1)
	check_bb := f64.Rectangle{
		window.DC.CursorPos,
		window.DC.CursorPos.Add(f64.Vec2{
			label_size.Y + style.FramePadding.Y*2 - 1,
			label_size.Y + style.FramePadding.Y*2 - 1,
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
		c.ItemSizeEx(f64.Vec2{text_bb.Dx(), check_bb.Dy()}, style.FramePadding.Y)
		total_bb = total_bb.Union(text_bb)
	}

	if !c.ItemAdd(total_bb, id) {
		return false
	}

	center := check_bb.Center()
	center.X = float64(int(center.X + 0.5))
	center.Y = float64(int(center.Y + 0.5))
	radius := check_bb.Dy() * 0.5

	hovered, held, pressed := c.ButtonBehavior(total_bb, id, 0)

	var col color.RGBA
	switch {
	case hovered && held:
		col = c.GetColorFromStyle(ColFrameBgActive)
	case hovered:
		col = c.GetColorFromStyle(ColFrameBgHovered)
	default:
		col = c.GetColorFromStyle(ColFrameBg)
	}

	c.RenderNavHighlight(total_bb, id)
	window.DrawList.AddCircleFilledEx(center, radius, col, 16)
	if active {
		check_sz := math.Min(check_bb.Dx(), check_bb.Dy())
		pad := math.Max(1.0, float64(int(check_sz/6.0)))
		window.DrawList.AddCircleFilledEx(center, radius-pad, c.GetColorFromStyle(ColCheckMark), 16)
	}

	if style.FrameBorderSize > 0.0 {
		window.DrawList.AddCircleEx(center.Add(f64.Vec2{1, 1}), radius, c.GetColorFromStyle(ColBorderShadow), 16, style.FrameBorderSize)
		window.DrawList.AddCircleEx(center, radius, c.GetColorFromStyle(ColBorder), 16, style.FrameBorderSize)
	}

	if c.LogEnabled {
		if active {
			c.LogRenderedText(&text_bb.Min, "(x)")
		} else {
			c.LogRenderedText(&text_bb.Min, "( )")
		}
	}
	if label_size.X > 0.0 {
		c.RenderText(text_bb.Min, label)
	}

	return pressed
}

func (c *Context) RadioButtonEx(label string, v *int, v_button int) bool {
	pressed := c.RadioButton(label, *v == v_button)
	if pressed {
		*v = v_button
	}
	return pressed
}

func (c *Context) Bullet() {
	style := &c.Style
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return
	}

	line_height := math.Max(math.Min(window.DC.CurrentLineHeight, c.FontSize+c.Style.FramePadding.Y*2), c.FontSize)
	bb := f64.Rectangle{window.DC.CursorPos, window.DC.CursorPos.Add(f64.Vec2{c.FontSize, line_height})}
	c.ItemSizeBB(bb)
	if !c.ItemAdd(bb, 0) {
		c.SameLineEx(0, style.FramePadding.X*2)
		return
	}

	// Render and stay on same line
	c.RenderBullet(bb.Min.Add(f64.Vec2{style.FramePadding.X + c.FontSize*0.5, line_height * 0.5}))
	c.SameLineEx(0, style.FramePadding.X*2)
}

// Text with a little bullet aligned to the typical tree node.
func (c *Context) BulletText(format string, args ...interface{}) {
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return
	}
	style := &c.Style
	text := fmt.Sprintf(format, args...)
	label_size := c.CalcTextSizeEx(text, false, -1)
	// Latch before ItemSize changes it
	text_base_offset_y := math.Max(0.0, window.DC.CurrentLineTextBaseOffset)
	line_height := math.Max(math.Min(window.DC.CurrentLineHeight, c.FontSize+c.Style.FramePadding.Y*2), c.FontSize)
	// Empty text doesn't add padding
	offset := f64.Vec2{c.FontSize, math.Max(line_height, label_size.Y)}
	if label_size.X > 0.0 {
		offset.X += label_size.X + style.FramePadding.X*2
	}
	bb := f64.Rectangle{
		window.DC.CursorPos,
		window.DC.CursorPos.Add(offset),
	}

	c.ItemSizeBB(bb)
	if !c.ItemAdd(bb, 0) {
		return
	}

	// Render
	c.RenderBullet(bb.Min.Add(f64.Vec2{style.FramePadding.X + c.FontSize*0.5, line_height * 0.5}))
	c.RenderTextEx(bb.Min.Add(f64.Vec2{c.FontSize + style.FramePadding.X*2, text_base_offset_y}), text, false)
}

// Helper to display logging buttons
func (c *Context) LogButtons() {
	c.PushStringID("LogButtons")
	log_to_tty := c.Button("Log To TTY")
	c.SameLine()
	log_to_file := c.Button("Log To File")
	c.SameLine()
	log_to_clipboard := c.Button("Log To Clipboard")
	c.SameLine()
	c.PushItemWidth(80.0)
	c.PushAllowKeyboardFocus(false)
	c.SliderIntEx("Depth", &c.LogAutoExpandMaxDepth, 0, 9, "")
	c.PopAllowKeyboardFocus()
	c.PopItemWidth()
	c.PopID()

	// Start logging at the end of the function so that the buttons don't appear in the log
	if log_to_tty {
		c.LogToTTYEx(c.LogAutoExpandMaxDepth)
	}
	if log_to_file {
		c.LogToFile(c.LogAutoExpandMaxDepth, c.IO.LogFilename)
	}
	if log_to_clipboard {
		c.LogToClipboardEx(c.LogAutoExpandMaxDepth)
	}
}

func (c *Context) ArrowButton(str_id string, dir Dir) bool {
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	id := window.GetID(str_id)
	sz := c.GetFrameHeight()
	bb := f64.Rectangle{
		window.DC.CursorPos,
		window.DC.CursorPos.Add(f64.Vec2{sz, sz}),
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
	c.RenderFrameEx(bb.Min, bb.Max, col, true, c.Style.FrameRounding)
	c.RenderArrow(bb.Min.Add(c.Style.FramePadding), dir)

	return pressed
}
