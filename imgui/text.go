package imgui

import (
	"bytes"
	"fmt"
	"image/color"
	"math"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/qeedquan/go-media/math/f64"
	"github.com/qeedquan/go-media/math/mathutil"
	"github.com/qeedquan/go-media/stb/stbte"
)

// Shared state of InputText(), passed to callback when a ImGuiInputTextFlags_Callback* flag is used and the corresponding callback is triggered.
type TextEditCallbackData struct {
	EventFlag InputTextFlags // One of ImGuiInputTextFlags_Callback* // Read-only
	Flags     InputTextFlags // What user passed to InputText()      // Read-only
	ReadOnly  bool           // Read-only mode                       // Read-only

	// CharFilter event:
	EventChar rune // Character input                      // Read-write (replace character or set to zero)

	// Completion,History,Always events:
	// If you modify the buffer contents make sure you update 'BufTextLen' and set 'BufDirty' to true.
	EventKey       Key    // Key pressed (Up/Down/TAB)            // Read-only
	Buf            []byte // Current text buffer                  // Read-write (pointed data only, can't replace the actual pointer)
	BufDirty       bool   // Set if you modify Buf/BufTextLen!!   // Write
	CursorPos      int    //                                      // Read-write
	SelectionStart int    //                                      // Read-write (== to SelectionEnd when no selection)
	SelectionEnd   int    //                                      // Read-write
}

type TextEditCallback func(*TextEditCallbackData) int

type SeparatorFlags int

const (
	SeparatorFlagsHorizontal SeparatorFlags = 1 << 0 // Axis default to current layout type, so generally Horizontal unless e.g. in a menu bar
	SeparatorFlagsVertical   SeparatorFlags = 1 << 1
)

type InputTextFlags int

const (
	InputTextFlagsCharsDecimal        InputTextFlags = 1 << 0  // Allow 0123456789.+-*/
	InputTextFlagsCharsHexadecimal    InputTextFlags = 1 << 1  // Allow 0123456789ABCDEFabcdef
	InputTextFlagsCharsUppercase      InputTextFlags = 1 << 2  // Turn a..z into A..Z
	InputTextFlagsCharsNoBlank        InputTextFlags = 1 << 3  // Filter out spaces tabs
	InputTextFlagsAutoSelectAll       InputTextFlags = 1 << 4  // Select entire text when first taking mouse focus
	InputTextFlagsEnterReturnsTrue    InputTextFlags = 1 << 5  // Return 'true' when Enter is pressed (as opposed to when the value was modified)
	InputTextFlagsCallbackCompletion  InputTextFlags = 1 << 6  // Call user function on pressing TAB (for completion handling)
	InputTextFlagsCallbackHistory     InputTextFlags = 1 << 7  // Call user function on pressing Up/Down arrows (for history handling)
	InputTextFlagsCallbackAlways      InputTextFlags = 1 << 8  // Call user function every time. User code may query cursor position modify text buffer.
	InputTextFlagsCallbackCharFilter  InputTextFlags = 1 << 9  // Call user function to filter character. Modify data->EventChar to replace/filter input or return 1 to discard character.
	InputTextFlagsAllowTabInput       InputTextFlags = 1 << 10 // Pressing TAB input a '\t' character into the text field
	InputTextFlagsCtrlEnterForNewLine InputTextFlags = 1 << 11 // In multi-line mode unfocus with Enter add new line with Ctrl+Enter (default is opposite: unfocus with Ctrl+Enter add line with Enter).
	InputTextFlagsNoHorizontalScroll  InputTextFlags = 1 << 12 // Disable following the cursor horizontally
	InputTextFlagsAlwaysInsertMode    InputTextFlags = 1 << 13 // Insert mode
	InputTextFlagsReadOnly            InputTextFlags = 1 << 14 // Read-only mode
	InputTextFlagsPassword            InputTextFlags = 1 << 15 // Password mode display all characters as '*'
	InputTextFlagsNoUndoRedo          InputTextFlags = 1 << 16 // Disable undo/redo. Note that input text owns the text data while active if you want to provide your own undo/redo stack you need e.g. to call ClearActiveID().
	InputTextFlagsCharsScientific     InputTextFlags = 1 << 17 // Allow 0123456789.+-*/eE (Scientific notation input)
	// [Internal]
	InputTextFlagsMultiline InputTextFlags = 1 << 20 // For internal use by InputTextMultiline()
)

type TextEditState struct {
	Ctx                  *Context
	Id                   ID     // widget id owning the text state
	Text                 []rune // edit buffer, we need to persist but can't guarantee the persistence of the user-provided buffer. so we copy into own buffer.
	InitialText          []byte // backup of end-user buffer at the time of focus (in UTF-8, unaltered)
	CurLenA, CurLenW     int    // we need to maintain our buffer length in both UTF-8 and wchar format.
	BufSizeA             int    // end-user buffer size
	ScrollX              float64
	StbState             stbte.State
	CursorAnim           float64
	CursorFollow         bool
	SelectedAllMouseLock bool
}

func (c *Context) GetTextLineHeight() float64 {
	return c.FontSize
}

func (c *Context) GetTextLineHeightWithSpacing() float64 {
	return c.FontSize + c.Style.ItemSpacing.Y
}

func (c *Context) CalcTextSize(text string) f64.Vec2 {
	return c.CalcTextSizeEx(text, false, -1)
}

// Calculate text size. Text can be multi-line. Optionally ignore text after a ## marker.
// CalcTextSize("") should return ImVec2(0.0f, GImGui->FontSize)
func (c *Context) CalcTextSizeEx(text string, hide_text_after_double_hash bool, wrap_width float64) f64.Vec2 {
	text_display_end := len(text)
	if hide_text_after_double_hash {
		// Hide anything after a '##' string
		text_display_end = c.FindRenderedTextEnd(text)
	}

	font := c.Font
	font_size := c.FontSize
	if text_display_end == 0 {
		return f64.Vec2{0, font_size}
	}
	text_size, _ := font.CalcTextSizeA(font_size, math.MaxFloat32, wrap_width, text[:text_display_end])

	// Cancel out character spacing for the last character of a line (it is baked into glyph->AdvanceX field)
	font_scale := font_size / font.FontSize
	character_spacing_x := 1.0 * font_scale
	if text_size.X > 0.0 {
		text_size.X -= character_spacing_x
	}
	text_size.X = float64(int(text_size.X + 0.95))

	return text_size
}

func (c *Context) FindRenderedTextEnd(text string) int {
	text_display_end := strings.Index(text, "##")
	if text_display_end == -1 {
		text_display_end = len(text)
	}
	return text_display_end
}

func (c *Context) Text(format string, args ...interface{}) {
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return
	}
	text := fmt.Sprintf(format, args...)
	c.TextUnformatted(text)
}

func (c *Context) TextUnformatted(text string) {
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return
	}

	text_pos := f64.Vec2{
		window.DC.CursorPos.X,
		window.DC.CursorPos.Y + window.DC.CurrentLineTextBaseOffset,
	}
	wrap_pos_x := window.DC.TextWrapPos
	wrap_enabled := wrap_pos_x >= 0.0

	if len(text) > 2000 && !wrap_enabled {
		// Long text!
		// Perform manual coarse clipping to optimize for long multi-line text
		// From this point we will only compute the width of lines that are visible. Optimization only available when word-wrapping is disabled.
		// We also don't vertically center the text within the line full height, which is unlikely to matter because we are likely the biggest and only item on the line.
		line := 0
		line_height := c.GetTextLineHeight()
		clip_rect := window.ClipRect
		text_size := f64.Vec2{0, 0}

		if text_pos.Y <= clip_rect.Max.Y {
			pos := text_pos

			// Lines to skip (can't skip when logging text)
			if !c.LogEnabled {
				lines_skippable := int((clip_rect.Min.Y - text_pos.Y) / line_height)
				if lines_skippable > 0 {
					lines_skipped := 0
					for line < len(text) && lines_skipped < lines_skippable {
						line_end := strings.IndexRune(text[line:], '\n')
						if line_end < 0 {
							line_end = len(text) - line - 1
						}
						line += line_end + 1
						lines_skipped++
					}
					pos.Y += float64(lines_skipped) * line_height
				}
			}

			// Lines to render
			if line < len(text) {
				line_rect := f64.Rectangle{pos, pos.Add(f64.Vec2{math.MaxFloat32, line_height})}
				for line < len(text) {
					line_end := strings.IndexRune(text[line:], '\n')
					if c.IsClippedEx(line_rect, 0, false) {
						break
					}

					line_size := c.CalcTextSizeEx(text[line:], false, -1)
					text_size.X = math.Max(text_size.X, line_size.X)
					c.RenderTextEx(pos, text[line:], false)
					if line_end < 0 {
						line_end = len(text) - line - 1
					}
					line += line_end + 1
					line_rect.Min.Y += line_height
					line_rect.Max.Y += line_height
					pos.Y += line_height
				}

				// Count remaining lines
				lines_skipped := 0
				for line < len(text) {
					line_end := strings.IndexRune(text[line:], '\n')
					if line_end < 0 {
						line_end = len(text) - line - 1
					}
					line = line_end + 1
					lines_skipped++
				}
				pos.Y += float64(lines_skipped) * line_height
			}

			text_size.Y += (pos.Sub(text_pos)).Y
		}

		bb := f64.Rectangle{text_pos, text_pos.Add(text_size)}
		c.ItemSizeBB(bb)
		c.ItemAdd(bb, 0)
	} else {
		wrap_width := 0.0
		if wrap_enabled {
			wrap_width = c.CalcWrapWidthForPos(window.DC.CursorPos, wrap_pos_x)
		}
		text_size := c.CalcTextSizeEx(text, false, wrap_width)

		// Account of baseline offset
		bb := f64.Rectangle{text_pos, text_pos.Add(text_size)}
		c.ItemSize(text_size)
		if !c.ItemAdd(bb, 0) {
			return
		}

		// Render (we don't hide text after ## in this end-user function)
		c.RenderTextWrapped(bb.Min, text, wrap_width)
	}
}

func (c *Context) RenderTextWrapped(pos f64.Vec2, text string, wrap_width float64) {
	window := c.CurrentWindow

	if len(text) > 0 {
		window.DrawList.AddTextEx(c.Font, c.FontSize, pos, c.GetColorFromStyle(ColText), text, wrap_width, nil)
		if c.LogEnabled {
			c.LogRenderedText(&pos, text)
		}
	}
}

// Horizontal separating line.
func (c *Context) Separator() {
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return
	}

	var flags SeparatorFlags
	if flags&(SeparatorFlagsHorizontal|SeparatorFlagsVertical) == 0 {
		if window.DC.LayoutType == LayoutTypeHorizontal {
			flags |= SeparatorFlagsVertical
		} else {
			flags |= SeparatorFlagsHorizontal
		}
	}
	// Check that only 1 option is selected
	assert(mathutil.IsPow2(int(flags & (SeparatorFlagsHorizontal | SeparatorFlagsVertical))))

	if flags&SeparatorFlagsVertical != 0 {
		c.VerticalSeparator()
		return
	}

	// Horizontal Separator
	if window.DC.ColumnsSet != nil {
		c.PopClipRect()
	}

	x1 := window.Pos.X
	x2 := window.Pos.X + window.Size.X
	if len(window.DC.GroupStack) > 0 {
		x1 += window.DC.IndentX
	}

	bb := f64.Rectangle{
		f64.Vec2{x1, window.DC.CursorPos.Y},
		f64.Vec2{x2, window.DC.CursorPos.Y + 1.0},
	}

	// NB: we don't provide our width so that it doesn't get feed back into AutoFit, we don't provide height to not alter layout
	c.ItemSize(f64.Vec2{0, 0})
	if !c.ItemAdd(bb, 0) {
		if window.DC.ColumnsSet != nil {
			c.PushColumnClipRect()
		}
		return
	}

	window.DrawList.AddLine(bb.Min, f64.Vec2{bb.Max.X, bb.Min.Y}, c.GetColorFromStyle(ColSeparator))
	if c.LogEnabled {
		c.LogRenderedText(nil, "--------------------------------\n")
	}

	if window.DC.ColumnsSet != nil {
		c.PushColumnClipRect()
		window.DC.ColumnsSet.LineMinY = window.DC.CursorPos.Y
	}
}

func (c *Context) VerticalSeparator() {
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return
	}

	y1 := window.DC.CursorPos.Y
	y2 := window.DC.CursorPos.Y + window.DC.CurrentLineHeight
	bb := f64.Rectangle{
		f64.Vec2{window.DC.CursorPos.X, y1},
		f64.Vec2{window.DC.CursorPos.X + 1.0, y2},
	}
	c.ItemSize(f64.Vec2{bb.Dx(), 0})
	if !c.ItemAdd(bb, 0) {
		return
	}

	window.DrawList.AddLine(f64.Vec2{bb.Min.X, bb.Min.Y}, f64.Vec2{bb.Min.X, bb.Max.Y}, c.GetColorFromStyle(ColSeparator))
	if c.LogEnabled {
		c.LogText(" |")
	}
}

func (c *Context) InputTextExCallback(label, buf string, size_arg f64.Vec2, flags InputTextFlags, callback func()) bool {
	return false
}

func (c *Context) AlignTextToFramePadding() {
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return
	}

	window.DC.CurrentLineHeight = math.Max(window.DC.CurrentLineHeight, c.FontSize+c.Style.FramePadding.Y*2)
	window.DC.CurrentLineTextBaseOffset = math.Max(window.DC.CurrentLineTextBaseOffset, c.Style.FramePadding.Y)
}

func (c *Context) TextWrapped(format string, args ...interface{}) {
	// Keep existing wrap position is one ia already set
	need_wrap := c.CurrentWindow.DC.TextWrapPos < 0.0
	if need_wrap {
		c.PushTextWrapPos(0.0)
	}
	c.Text(format, args...)
	if need_wrap {
		c.PopTextWrapPos()
	}
}

func (c *Context) PushTextWrapPos(wrap_pos_x float64) {
	window := c.GetCurrentWindow()
	window.DC.TextWrapPos = wrap_pos_x
	window.DC.TextWrapPosStack = append(window.DC.TextWrapPosStack, wrap_pos_x)
}

func (c *Context) PopTextWrapPos() {
	window := c.GetCurrentWindow()
	window.DC.TextWrapPosStack = window.DC.TextWrapPosStack[:len(window.DC.TextWrapPosStack)-1]
	window.DC.TextWrapPos = -1.0
	if len(window.DC.TextWrapPosStack) > 0 {
		window.DC.TextWrapPos = window.DC.TextWrapPosStack[len(window.DC.TextWrapPosStack)-1]
	}
}

func (c *Context) InputText(label string, buf []byte) bool {
	// call InputTextMultiline()
	return c.InputTextEx(label, buf, f64.Vec2{0, 0}, 0, nil)
}

func (c *Context) InputTextMultiline(label string, buf []byte, size f64.Vec2, flags InputTextFlags, callback TextEditCallback) bool {
	return c.InputTextEx(label, buf, size, flags|InputTextFlagsMultiline, callback)
}

// Edit a string of text
// NB: when active, hold on a privately held copy of the text (and apply back to 'buf'). So changing 'buf' while active has no effect.
// FIXME: Rather messy function partly because we are doing UTF8 > u16 > UTF8 conversions on the go to more easily handle stb_textedit calls. Ideally we should stay in UTF-8 all the time. See https://github.com/nothings/stb/issues/188
func (c *Context) InputTextEx(label string, buf []byte, size_arg f64.Vec2, flags InputTextFlags, callback TextEditCallback) bool {
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	var temp_text_buffer []byte
	style := &c.Style
	io := &c.IO

	// Can't use both together (they both use up/down keys)
	assert(!(flags&InputTextFlagsCallbackHistory != 0 && flags&InputTextFlagsMultiline != 0))
	// Can't use both together (they both use tab key)
	assert(!(flags&InputTextFlagsCallbackCompletion != 0 && flags&InputTextFlagsAllowTabInput != 0))

	is_multiline := flags&InputTextFlagsMultiline != 0
	is_editable := flags&InputTextFlagsReadOnly == 0
	is_password := flags&InputTextFlagsPassword != 0
	is_undoable := flags&InputTextFlagsNoUndoRedo == 0

	// Open group before calling GetID() because groups tracks id created during their spawn
	if is_multiline {
		c.BeginGroup()
	}
	id := window.GetID(label)
	label_size := c.CalcTextSizeEx(label, true, -1)

	// Arbitrary default of 8 lines high for multi-line
	item_size := label_size.Y
	if is_multiline {
		item_size = c.GetTextLineHeight() * 8.0
	}
	size := c.CalcItemSize(size_arg, c.CalcItemWidth(), item_size+style.FramePadding.Y*2.0)
	frame_bb := f64.Rectangle{window.DC.CursorPos, window.DC.CursorPos.Add(size)}
	total_bb_x := 0.0
	if label_size.X > 0 {
		total_bb_x = style.ItemInnerSpacing.X + label_size.X
	}
	total_bb := f64.Rectangle{frame_bb.Min, frame_bb.Max.Add(f64.Vec2{total_bb_x, 0.0})}

	draw_window := window
	if is_multiline {
		c.ItemAddEx(total_bb, id, &frame_bb)
		if !c.BeginChildFrame(id, frame_bb.Size(), 0) {
			c.EndChildFrame()
			c.EndGroup()
			return false
		}
		draw_window = c.GetCurrentWindow()
		size.X -= draw_window.ScrollbarSizes.X
	} else {
		c.ItemSizeBBEx(total_bb, style.FramePadding.Y)
		if !c.ItemAddEx(total_bb, id, &frame_bb) {
			return false
		}
	}

	hovered := c.ItemHoverable(frame_bb, id)
	if hovered {
		c.MouseCursor = MouseCursorTextInput
	}

	// Password pushes a temporary font with only a fallback glyph
	if is_password {
		glyph := c.Font.FindGlyph('*')
		password_font := &c.InputTextPasswordFont
		password_font.FontSize = c.Font.FontSize
		password_font.Scale = c.Font.Scale
		password_font.DisplayOffset = c.Font.DisplayOffset
		password_font.Ascent = c.Font.Ascent
		password_font.Descent = c.Font.Descent
		password_font.ContainerAtlas = c.Font.ContainerAtlas
		password_font.FallbackGlyph = glyph
		password_font.FallbackAdvanceX = glyph.AdvanceX
		assert(len(password_font.Glyphs) == 0 && len(password_font.IndexAdvanceX) == 0 && len(password_font.IndexLookup) == 0)
		c.PushFont(password_font)
	}

	// NB: we are only allowed to access 'edit_state' if we are the active widget.
	edit_state := &c.InputTextState

	// Using completion callback disable keyboard tabbing
	focus_requested := c.FocusableItemRegisterEx(window, id, flags&(InputTextFlagsCallbackCompletion|InputTextFlagsAllowTabInput) != 0)
	focus_requested_by_code := focus_requested && (window.FocusIdxAllCounter == window.FocusIdxAllRequestCurrent)
	focus_requested_by_tab := focus_requested && !focus_requested_by_code

	user_clicked := hovered && io.MouseClicked[0]
	user_scrolled := is_multiline && c.ActiveId == 0 && edit_state.Id == id && c.ActiveIdPreviousFrame == draw_window.GetIDNoKeepAlive("#SCROLLY")
	user_nav_input_start := (c.ActiveId != id) && ((c.NavInputId == id) || (c.NavActivateId == id && c.NavInputSource == InputSourceNavKeyboard))

	clear_active_id := false

	select_all := (c.ActiveId != id) && ((flags&InputTextFlagsAutoSelectAll) != 0 || user_nav_input_start) && (!is_multiline)
	if focus_requested || user_clicked || user_scrolled || user_nav_input_start {
		if c.ActiveId != id {
			// Start edition
			// Take a copy of the initial buffer value (both in original UTF-8 format and converted to wchar)
			// From the moment we focused we are ignoring the content of 'buf' (unless we are in read-only mode)
			prev_len_w := edit_state.CurLenW
			edit_state.Text = []rune(string(buf))
			edit_state.InitialText = []byte(buf)
			edit_state.CurLenW = len(edit_state.Text)
			edit_state.CurLenA = len(buf)
			edit_state.CursorAnimReset()

			// Preserve cursor position and undo/redo stack if we come back to same widget
			// FIXME: We should probably compare the whole buffer to be on the safety side. Comparing buf (utf8) and edit_state.Text (wchar).
			recycle_state := (edit_state.Id == id) && (prev_len_w == edit_state.CurLenW)
			if recycle_state {
				// Recycle existing cursor/selection/undo stack but clamp position
				// Note a single mouse click will override the cursor/position immediately by calling stb_textedit_click handler.
				edit_state.CursorClamp()
			} else {
				edit_state.Id = id
				edit_state.ScrollX = 0.0
				edit_state.StbState.Init(!is_multiline)
				if !is_multiline && focus_requested_by_code {
					select_all = true
				}
			}

			if flags&InputTextFlagsAlwaysInsertMode != 0 {
				edit_state.StbState.SetInsertMode(true)
			}
			if !is_multiline && (focus_requested_by_tab || (user_clicked && io.KeyCtrl)) {
				select_all = true
			}
		}
		c.SetActiveID(id, window)
		c.SetFocusID(id, window)
		c.FocusWindow(window)
		if !is_multiline && flags&InputTextFlagsCallbackHistory == 0 {
			c.ActiveIdAllowNavDirFlags |= ((1 << uint(DirUp)) | (1 << uint(DirDown)))
		}
	} else if io.MouseClicked[0] {
		// Release focus when we click outside
		clear_active_id = true
	}

	value_changed := false
	enter_pressed := false
	if c.ActiveId == id {
		if !is_editable && !c.ActiveIdIsJustActivated {
			// When read-only we always use the live data passed to the function
			if len(buf) > len(edit_state.Text) {
				edit_state.Text = append(edit_state.Text, make([]rune, len(buf)-len(edit_state.Text))...)
			}
			edit_state.CurLenW = len(edit_state.Text)
			edit_state.CurLenA = len(buf)
			edit_state.CursorClamp()
		}

		edit_state.BufSizeA = len(buf)

		// Although we are active we don't prevent mouse from hovering other elements unless we are interacting right now with the widget.
		// Down the line we should have a cleaner library-wide concept of Selected vs Active.
		c.ActiveIdAllowOverlap = !io.MouseDown[0]
		c.WantTextInputNextFrame = 1

		// Edit in progress
		mouse_x := (io.MousePos.X - frame_bb.Min.X - style.FramePadding.X) + edit_state.ScrollX
		mouse_y := c.FontSize * 0.5
		if is_multiline {
			mouse_y = io.MousePos.Y - draw_window.DC.CursorPos.Y - style.FramePadding.Y
		}

		is_osx := io.OptMacOSXBehaviors
		if select_all || (hovered && !is_osx && io.MouseDoubleClicked[0]) {
			edit_state.SelectAll()
			edit_state.SelectedAllMouseLock = true
		} else if hovered && is_osx && io.MouseDoubleClicked[0] {
			// Double-click select a word only, OS X style (by simulating keystrokes)
			edit_state.OnKeyPressed(stbte.K_WORDLEFT)
			edit_state.OnKeyPressed(stbte.K_WORDRIGHT | stbte.K_SHIFT)
		} else if io.MouseClicked[0] && !edit_state.SelectedAllMouseLock {
			if hovered {
				edit_state.StbState.Click(edit_state, mouse_x, mouse_y)
				edit_state.CursorAnimReset()
			}
		} else if io.MouseDown[0] && !edit_state.SelectedAllMouseLock && (io.MouseDelta.X != 0.0 || io.MouseDelta.Y != 0.0) {
			edit_state.StbState.Drag(edit_state, mouse_x, mouse_y)
			edit_state.CursorAnimReset()
			edit_state.CursorFollow = true
		}

		if io.InputCharacters[0] != 0 {
			// Process text input (before we check for Return because using some IME will effectively send a Return?)
			// We ignore CTRL inputs, but need to allow ALT+CTRL as some keyboards (e.g. German) use AltGR (which _is_ Alt+Ctrl) to input certain characters.
			ignore_inputs := (io.KeyCtrl && !io.KeyAlt) || (is_osx && io.KeySuper)
			if !ignore_inputs && is_editable && !user_nav_input_start {
				for n := 0; n < len(io.InputCharacters); n++ {
					// Insert character if they pass filtering
					ch := io.InputCharacters[n]
					if c.InputTextFilterCharacter(&ch, flags, callback) {
						edit_state.OnKeyPressed(int(ch))
					}
				}
			}
			// Consume characters
			for i := range c.IO.InputCharacters {
				c.IO.InputCharacters[i] = 0
			}
		}
	}

	cancel_edit := false
	if c.ActiveId == id && !c.ActiveIdIsJustActivated && !clear_active_id {
		// Handle key-presses
		k_mask := 0
		if io.KeyShift {
			k_mask = stbte.K_SHIFT
		}
		is_osx := io.OptMacOSXBehaviors
		// OS X style: Shortcuts using Cmd/Super instead of Ctrl
		is_shortcut_key := io.KeyCtrl && !io.KeySuper && !io.KeyAlt && !io.KeyShift
		if is_osx {
			is_shortcut_key = io.KeySuper && !io.KeyCtrl && !io.KeyAlt && !io.KeyShift
		}
		is_osx_shift_shortcut := is_osx && io.KeySuper && io.KeyShift && !io.KeyCtrl && !io.KeyAlt
		// OS X style: Text editing cursor movement using Alt instead of Ctrl
		is_wordmove_key_down := io.KeyCtrl
		if is_osx {
			is_wordmove_key_down = io.KeyAlt
		}
		// OS X style: Line/Text Start and End using Cmd+Arrows instead of Home/End
		is_startend_key_down := is_osx && io.KeySuper && !io.KeyCtrl && !io.KeyAlt
		is_ctrl_key_only := io.KeyCtrl && !io.KeyShift && !io.KeyAlt && !io.KeySuper
		is_shift_key_only := io.KeyShift && !io.KeyCtrl && !io.KeyAlt && !io.KeySuper

		is_cut := ((is_shortcut_key && c.IsKeyPressedMap(KeyX)) || (is_shift_key_only && c.IsKeyPressedMap(KeyDelete))) && is_editable && !is_password && (!is_multiline || edit_state.HasSelection())
		is_copy := ((is_shortcut_key && c.IsKeyPressedMap(KeyC)) || (is_ctrl_key_only && c.IsKeyPressedMap(KeyInsert))) && !is_password && (!is_multiline || edit_state.HasSelection())
		is_paste := ((is_shortcut_key && c.IsKeyPressedMap(KeyV)) || (is_shift_key_only && c.IsKeyPressedMap(KeyInsert))) && is_editable
		is_undo := ((is_shortcut_key && c.IsKeyPressedMap(KeyZ)) && is_editable && is_undoable)
		is_redo := ((is_shortcut_key && c.IsKeyPressedMap(KeyY)) || (is_osx_shift_shortcut && c.IsKeyPressedMap(KeyZ))) && is_editable && is_undoable

		if c.IsKeyPressedMap(KeyLeftArrow) {
			switch {
			case is_startend_key_down:
				edit_state.OnKeyPressed(stbte.K_LINESTART | k_mask)
			case is_wordmove_key_down:
				edit_state.OnKeyPressed(stbte.K_WORDLEFT | k_mask)
			default:
				edit_state.OnKeyPressed(stbte.K_LEFT | k_mask)
			}
		} else if c.IsKeyPressedMap(KeyRightArrow) {
			switch {
			case is_startend_key_down:
				edit_state.OnKeyPressed(stbte.K_LINEEND | k_mask)
			case is_wordmove_key_down:
				edit_state.OnKeyPressed(stbte.K_WORDRIGHT | k_mask)
			default:
				edit_state.OnKeyPressed(stbte.K_RIGHT | k_mask)
			}
		} else if c.IsKeyPressedMap(KeyUpArrow) && is_multiline {
			if io.KeyCtrl {
				c.SetWindowScrollY(draw_window, math.Max(draw_window.Scroll.Y-c.FontSize, 0.0))
			} else {
				if is_startend_key_down {
					edit_state.OnKeyPressed(stbte.K_TEXTSTART | k_mask)
				} else {
					edit_state.OnKeyPressed(stbte.K_UP | k_mask)
				}
			}
		} else if c.IsKeyPressedMap(KeyDownArrow) && is_multiline {
			if io.KeyCtrl {
				c.SetWindowScrollY(draw_window, math.Min(draw_window.Scroll.Y+c.FontSize, c.GetScrollMaxY()))
			} else {
				if is_startend_key_down {
					edit_state.OnKeyPressed(stbte.K_TEXTEND | k_mask)
				} else {
					edit_state.OnKeyPressed(stbte.K_DOWN | k_mask)
				}
			}
		} else if c.IsKeyPressedMap(KeyHome) {
			if io.KeyCtrl {
				edit_state.OnKeyPressed(stbte.K_TEXTSTART | k_mask)
			} else {
				edit_state.OnKeyPressed(stbte.K_LINESTART | k_mask)
			}
		} else if c.IsKeyPressedMap(KeyEnd) {
			if io.KeyCtrl {
				edit_state.OnKeyPressed(stbte.K_TEXTEND | k_mask)
			} else {
				edit_state.OnKeyPressed(stbte.K_LINEEND | k_mask)
			}
		} else if c.IsKeyPressedMap(KeyDelete) && is_editable {
			edit_state.OnKeyPressed(stbte.K_DELETE | k_mask)
		} else if c.IsKeyPressedMap(KeyBackspace) && is_editable {
			if !edit_state.HasSelection() {
				if is_wordmove_key_down {
					edit_state.OnKeyPressed(stbte.K_WORDLEFT | stbte.K_SHIFT)
				} else if is_osx && io.KeySuper && !io.KeyAlt && !io.KeyCtrl {
					edit_state.OnKeyPressed(stbte.K_LINESTART | stbte.K_SHIFT)
				}
			}
			edit_state.OnKeyPressed(stbte.K_BACKSPACE | k_mask)
		} else if c.IsKeyPressedMap(KeyEnter) {
			ctrl_enter_for_new_line := flags&InputTextFlagsCtrlEnterForNewLine != 0
			if !is_multiline || (ctrl_enter_for_new_line && !io.KeyCtrl) || (!ctrl_enter_for_new_line && io.KeyCtrl) {
				enter_pressed, clear_active_id = true, true
			} else if is_editable {
				ch := '\n' // Insert new line
				if c.InputTextFilterCharacter(&ch, flags, callback) {
					edit_state.OnKeyPressed(int(ch))
				}
			}
		} else if flags&InputTextFlagsAllowTabInput != 0 && c.IsKeyPressedMap(KeyTab) && !io.KeyCtrl && !io.KeyShift && !io.KeyAlt && is_editable {
			ch := '\t' // Insert TAB
			if c.InputTextFilterCharacter(&ch, flags, callback) {
				edit_state.OnKeyPressed(int(ch))
			}
		} else if c.IsKeyPressedMap(KeyEscape) {
			clear_active_id, cancel_edit = true, true
		} else if is_undo || is_redo {
			if is_undo {
				edit_state.OnKeyPressed(stbte.K_UNDO)
			} else {
				edit_state.OnKeyPressed(stbte.K_REDO)
			}
			edit_state.ClearSelection()
		} else if is_shortcut_key && c.IsKeyPressedMap(KeyA) {
			edit_state.SelectAll()
			edit_state.CursorFollow = true
		} else if is_cut || is_copy {
			// Cut, Copy
			if io.SetClipboardTextFn != nil {
				ib, ie := 0, edit_state.CurLenW
				if edit_state.HasSelection() {
					ib = mathutil.Min(edit_state.StbState.SelectStart(), edit_state.StbState.SelectEnd())
					ie = mathutil.Max(edit_state.StbState.SelectStart(), edit_state.StbState.SelectEnd())
					str := string(edit_state.Text[ib:ie])
					c.SetClipboardText(str)
				}
			}

			if is_cut {
				if !edit_state.HasSelection() {
					edit_state.SelectAll()
				}
				edit_state.CursorFollow = true
				edit_state.StbState.Cut(edit_state)
			}
		} else if is_paste {
			clipboard := c.GetClipboardText()

			// Filter pasted buffer
			var clipboard_filtered []rune
			for _, ch := range clipboard {
				if ch == 0 {
					break
				}
				if ch >= 0x10000 || !c.InputTextFilterCharacter(&ch, flags, callback) {
					continue
				}
				clipboard_filtered = append(clipboard_filtered, ch)
			}
			// If everything was filtered, ignore the pasting operation
			if len(clipboard_filtered) > 0 {
				edit_state.StbState.Paste(edit_state, clipboard_filtered)
				edit_state.CursorFollow = true
			}
		}
	}

	if c.ActiveId == id {
		if cancel_edit {
			// Restore initial value. Only return true if restoring to the initial value changes the current buffer contents.
			if is_editable && string(buf) != string(edit_state.InitialText) {
				buf = edit_state.InitialText
				value_changed = true
			}
		}

		// When using 'ImGuiInputTextFlags_EnterReturnsTrue' as a special case we reapply the live buffer back to the input buffer before clearing ActiveId, even though strictly speaking it wasn't modified on this frame.
		// If we didn't do that, code like InputInt() with ImGuiInputTextFlags_EnterReturnsTrue would fail. Also this allows the user to use InputText() with ImGuiInputTextFlags_EnterReturnsTrue without maintaining any user-side storage.
		apply_edit_back_to_user_buffer := !cancel_edit || (enter_pressed && (flags&InputTextFlagsEnterReturnsTrue) != 0)
		if apply_edit_back_to_user_buffer {
			// Apply new value immediately - copy modified buffer back
			// Note that as soon as the input box is active, the in-widget value gets priority over any underlying modification of the input buffer
			// FIXME: We actually always render 'buf' when calling DrawList->AddText, making the comment above incorrect.
			// FIXME-OPT: CPU waste to do this every time the widget is active, should mark dirty state from the stb_textedit callbacks.
			if is_editable {
				temp_text_buffer = []byte(buf)
			}

			// User callback
			if flags&(InputTextFlagsCallbackCompletion|InputTextFlagsCallbackHistory|InputTextFlagsCallbackAlways) != 0 {
				assert(callback != nil)

				// The reason we specify the usage semantic (Completion/History) is that Completion needs to disable keyboard TABBING at the moment.
				event_flag := InputTextFlags(0)
				event_key := KeyCOUNT
				if (flags&InputTextFlagsCallbackCompletion) != 0 && c.IsKeyPressedMap(KeyTab) {
					event_flag = InputTextFlagsCallbackCompletion
					event_key = KeyTab
				} else if (flags&InputTextFlagsCallbackHistory) != 0 && c.IsKeyPressedMap(KeyUpArrow) {
					event_flag = InputTextFlagsCallbackHistory
					event_key = KeyUpArrow
				} else if (flags&InputTextFlagsCallbackHistory) != 0 && c.IsKeyPressedMap(KeyDownArrow) {
					event_flag = InputTextFlagsCallbackHistory
					event_key = KeyDownArrow
				} else if flags&InputTextFlagsCallbackAlways != 0 {
					event_flag = InputTextFlagsCallbackAlways
				}

				if event_flag != 0 {
					callback_data := TextEditCallbackData{
						EventFlag: event_flag,
						Flags:     flags,
						ReadOnly:  !is_editable,
						EventKey:  event_key,
						Buf:       temp_text_buffer,
						BufDirty:  false,
					}

					// We have to convert from wchar-positions to UTF-8-positions, which can be pretty slow (an incentive to ditch the ImWchar buffer, see https://github.com/nothings/stb/issues/188)
					text := edit_state.Text
					utf8_cursor_pos := TextCountUtf8BytesFromStr(text[:edit_state.StbState.Cursor()])
					utf8_selection_start := TextCountUtf8BytesFromStr(text[:edit_state.StbState.SelectStart()])
					utf8_selection_end := TextCountUtf8BytesFromStr(text[:edit_state.StbState.SelectEnd()])
					callback_data.CursorPos = utf8_cursor_pos
					callback_data.SelectionStart = utf8_selection_start
					callback_data.SelectionEnd = utf8_selection_end

					// Call user code
					callback(&callback_data)

					// Read back what user may have modified
					assert(len(callback_data.Buf) == edit_state.BufSizeA)
					assert(callback_data.Flags == flags)
					if callback_data.CursorPos != utf8_cursor_pos {
						edit_state.StbState.SetCursor(utf8.RuneCount(callback_data.Buf[:callback_data.CursorPos]))
					}
					if callback_data.SelectionStart != utf8_selection_start {
						edit_state.StbState.SetSelectStart(utf8.RuneCount(callback_data.Buf[:callback_data.SelectionStart]))
					}
					if callback_data.SelectionEnd != utf8_selection_end {
						edit_state.StbState.SetSelectEnd(utf8.RuneCount(callback_data.Buf[:callback_data.SelectionEnd]))
					}

					if callback_data.BufDirty {
						// Assume correct length and valid UTF-8 from user, saves us an extra strlen()
						edit_state.Text = []rune(string(callback_data.Buf))
						edit_state.CurLenW = len(edit_state.Text)
						edit_state.CurLenA = len(callback_data.Buf)
						edit_state.CursorAnimReset()
					}
				}
			}

			// Copy back to user buffer
			if is_editable && bytes.Compare(temp_text_buffer, buf) != 0 {
				copy(buf, temp_text_buffer)
				value_changed = true
			}
		}
	}

	// Release active ID at the end of the function (so e.g. pressing Return still does a final application of the value)
	if clear_active_id && c.ActiveId == id {
		c.ClearActiveID()
	}

	// Render
	// Select which buffer we are going to display. When ImGuiInputTextFlags_NoLiveEdit is set 'buf' might still be the old value. We set buf to NULL to prevent accidental usage from now on.
	var buf_display []byte
	if c.ActiveId == id && is_editable {
		buf_display = temp_text_buffer
	} else {
		buf_display = []byte(buf)
	}

	c.RenderNavHighlight(frame_bb, id)
	if !is_multiline {
		c.RenderFrameEx(frame_bb.Min, frame_bb.Max, c.GetColorFromStyle(ColFrameBg), true, style.FrameRounding)
	}

	// Not using frame_bb.Max because we have adjusted size
	clip_rect := f64.Vec4{frame_bb.Min.X, frame_bb.Min.Y, frame_bb.Min.X + size.X, frame_bb.Min.Y + size.Y}

	var render_pos f64.Vec2
	if is_multiline {
		render_pos = draw_window.DC.CursorPos
	} else {
		render_pos = frame_bb.Min.Add(style.FramePadding)
	}
	text_size := f64.Vec2{0, 0}
	is_currently_scrolling := (edit_state.Id == id && is_multiline && c.ActiveId == draw_window.GetIDNoKeepAlive("#SCROLLY"))
	if c.ActiveId == id || is_currently_scrolling {
		edit_state.CursorAnim += io.DeltaTime

		// This is going to be messy. We need to:
		// - Display the text (this alone can be more easily clipped)
		// - Handle scrolling, highlight selection, display cursor (those all requires some form of 1d->2d cursor position calculation)
		// - Measure text height (for scrollbar)
		// We are attempting to do most of that in **one main pass** to minimize the computation cost (non-negligible for large amount of text) + 2nd pass for selection rendering (we could merge them by an extra refactoring effort)
		// FIXME: This should occur on buf_display but we'd need to maintain cursor/select_start/select_end for UTF-8.
		var cursor_offset, select_start_offset f64.Vec2

		text_begin := edit_state.Text
		// Count lines + find lines numbers straddling 'cursor' and 'select_start' position.
		searches_input_ptr := [2]int{
			edit_state.StbState.Cursor(),
			-1,
		}
		searches_remaining := 1
		searches_result_line_number := [2]int{-1, -999}
		if edit_state.StbState.SelectStart() != edit_state.StbState.SelectEnd() {
			searches_input_ptr[1] = mathutil.Min(edit_state.StbState.SelectStart(), edit_state.StbState.SelectEnd())
			searches_result_line_number[1] = -1
			searches_remaining++
		}

		// Iterate all lines to find our line numbers
		// In multi-line mode, we never exit the loop until all lines are counted, so add one extra to the searches_remaining counter.
		if is_multiline {
			searches_remaining += 1
		}
		line_count := 0
		for i, s := range text_begin {
			if s == '\n' {
				line_count++
				if searches_result_line_number[0] == -1 && i >= searches_input_ptr[0] {
					searches_result_line_number[0] = line_count
					if searches_remaining--; searches_remaining <= 0 {
						break
					}
				}
				if searches_result_line_number[1] == -1 && i >= searches_input_ptr[1] {
					searches_result_line_number[1] = line_count
					if searches_remaining--; searches_remaining <= 0 {
						break
					}
				}
			}
		}
		line_count++
		if searches_result_line_number[0] == -1 {
			searches_result_line_number[0] = line_count
		}
		if searches_result_line_number[1] == -1 {
			searches_result_line_number[1] = line_count
		}

		// Calculate 2d position by finding the beginning of the line and measuring distance
		start := StrbolW(text_begin, searches_input_ptr[0], 0)
		end := searches_input_ptr[0]
		sz, _, _ := c.InputTextCalcTextSizeW(text_begin[start:end], false)
		cursor_offset.X = sz.X
		cursor_offset.Y = float64(searches_result_line_number[0]) * c.FontSize
		if searches_result_line_number[1] >= 0 {
			start := StrbolW(text_begin, searches_input_ptr[1], 0)
			end := searches_input_ptr[1]
			sz, _, _ := c.InputTextCalcTextSizeW(text_begin[start:end], false)
			cursor_offset.X = sz.X
			cursor_offset.Y = float64(searches_result_line_number[1]) * c.FontSize
		}

		// Store text height (note that we haven't calculated text width at all, see GitHub issues #383, #1224)
		if is_multiline {
			text_size = f64.Vec2{size.X, float64(line_count) * c.FontSize}
		}

		// Scroll
		if edit_state.CursorFollow {
			// Horizontal scroll in chunks of quarter width
			if flags&InputTextFlagsNoHorizontalScroll == 0 {
				scroll_increment_x := size.X * 0.25
				if cursor_offset.X < edit_state.ScrollX {
					edit_state.ScrollX = float64(int(math.Max(0.0, cursor_offset.X-scroll_increment_x)))
				} else if cursor_offset.X-size.X >= edit_state.ScrollX {
					edit_state.ScrollX = float64(int(cursor_offset.X - size.X + scroll_increment_x))
				}
			} else {
				edit_state.ScrollX = 0.0
			}

			// Vertical scroll
			if is_multiline {
				scroll_y := draw_window.Scroll.Y
				if cursor_offset.Y-c.FontSize < scroll_y {
					scroll_y = math.Max(0.0, cursor_offset.Y-c.FontSize)
				} else if cursor_offset.Y-size.Y >= scroll_y {
					scroll_y = cursor_offset.Y - size.Y
				}
				// To avoid a frame of lag
				draw_window.DC.CursorPos.Y += (draw_window.Scroll.Y - scroll_y)
				draw_window.Scroll.Y = scroll_y
				render_pos.Y = draw_window.DC.CursorPos.Y
			}
		}

		edit_state.CursorFollow = false
		render_scroll := f64.Vec2{edit_state.ScrollX, 0.0}

		// Draw selection
		if edit_state.StbState.SelectStart() != edit_state.StbState.SelectEnd() {
			text_selected_begin := mathutil.Min(edit_state.StbState.SelectStart(), edit_state.StbState.SelectEnd())
			text_selected_end := mathutil.Max(edit_state.StbState.SelectStart(), edit_state.StbState.SelectEnd())

			// FIXME: those offsets should be part of the style? they don't play so well with multi-line selection.
			bg_offy_up := -1.0
			bg_offy_dn := 2.0
			if is_multiline {
				bg_offy_up = 0
				bg_offy_dn = 0
			}
			bg_color := c.GetColorFromStyle(ColTextSelectedBg)
			rect_pos := render_pos.Add(select_start_offset).Sub(render_scroll)
			for p := text_selected_begin; p < text_selected_end; {
				if rect_pos.Y > clip_rect.W+c.FontSize {
					break
				}
				if rect_pos.Y < clip_rect.Y {
					for p < text_selected_end {
						if text_begin[p] == '\n' {
							break
						}
						p++
					}
				} else {
					rect_size, remaining, _ := c.InputTextCalcTextSizeW(text_begin[p:text_selected_end], true)
					p += remaining
					if rect_size.X <= 0.0 {
						// So we can see selected empty lines
						rect_size.X = float64(int((c.Font.GetCharAdvance(' ') * 0.50)))
						rect := f64.Rectangle{
							rect_pos.Add(f64.Vec2{0.0, bg_offy_up - c.FontSize}),
							rect_pos.Add(f64.Vec2{rect_size.X, bg_offy_dn}),
						}
						clip_rect_ := f64.Rect(clip_rect.X, clip_rect.Y, clip_rect.Z, clip_rect.W)
						rect = rect.Intersect(clip_rect_)
						if rect.Overlaps(clip_rect_) {
							draw_window.DrawList.AddRectFilled(rect.Min, rect.Max, bg_color)
						}
					}
				}
				rect_pos.X = render_pos.X - render_scroll.X
				rect_pos.Y += c.FontSize
			}
		}

		clip := &clip_rect
		if is_multiline {
			clip = nil
		}
		draw_window.DrawList.AddTextEx(c.Font, c.FontSize, render_pos.Sub(render_scroll), c.GetColorFromStyle(ColText), string(buf_display), 0.0, clip)

		// Draw blinking cursor
		cursor_is_visible := (!c.IO.OptCursorBlink) || (c.InputTextState.CursorAnim <= 0.0) || math.Mod(c.InputTextState.CursorAnim, 1.20) <= 0.80
		cursor_screen_pos := render_pos.Add(cursor_offset).Sub(render_scroll)
		cursor_screen_rect := f64.Rect(cursor_screen_pos.X, cursor_screen_pos.Y-c.FontSize+0.5, cursor_screen_pos.X+1.0, cursor_screen_pos.Y-1.5)
		clip_rect_ := f64.Rect(clip_rect.X, clip_rect.Y, clip_rect.Z, clip_rect.W)
		if cursor_is_visible && cursor_screen_rect.Overlaps(clip_rect_) {
			draw_window.DrawList.AddLine(cursor_screen_rect.Min, cursor_screen_rect.BL(), c.GetColorFromStyle(ColText))
		}

		// Notify OS of text input position for advanced IME (-1 x offset so that Windows IME can cover our cursor. Bit of an extra nicety.)
		if is_editable {
			c.PlatformImePos = f64.Vec2{cursor_screen_pos.X - 1, cursor_screen_pos.Y - c.FontSize}
		}
	} else {
		// Render text only
		if is_multiline {
			// We don't need width
			text_size = f64.Vec2{size.X, float64(InputTextCalcTextLenAndLineCount(string(buf_display))) * c.FontSize}
			clip := &clip_rect
			if is_multiline {
				clip = nil
			}
			draw_window.DrawList.AddTextEx(c.Font, c.FontSize, render_pos, c.GetColorFromStyle(ColText), string(buf_display), 0.0, clip)
		}
	}

	if is_multiline {
		// Always add room to scroll an extra line
		c.Dummy(text_size.Add(f64.Vec2{0.0, c.FontSize}))
		c.EndChildFrame()
		c.EndGroup()
	}

	if is_password {
		c.PopFont()
	}

	// Log as text
	if c.LogEnabled && !is_password {
		c.LogRenderedText(&render_pos, string(buf_display))
	}

	if label_size.X > 0 {
		c.RenderText(f64.Vec2{frame_bb.Max.X + style.ItemInnerSpacing.X, frame_bb.Min.Y + style.FramePadding.Y}, label)
	}

	if flags&InputTextFlagsEnterReturnsTrue != 0 {
		return enter_pressed
	}
	return value_changed
}

func (t *TextEditState) Init(ctx *Context) {
	*t = TextEditState{
		Ctx: ctx,
	}
}

func (t *TextEditState) CursorAnimReset() {
	// After a user-input the cursor stays on for a while without blinking
	t.CursorAnim = -0.30
}

func (t *TextEditState) CursorClamp() {
	t.StbState.SetCursor(mathutil.Min(t.StbState.Cursor(), t.CurLenW))
	t.StbState.SetSelectStart(mathutil.Min(t.StbState.SelectStart(), t.CurLenW))
	t.StbState.SetSelectEnd(mathutil.Min(t.StbState.SelectEnd(), t.CurLenW))
}

func (t *TextEditState) SelectAll() {
	t.StbState.SetSelectStart(0)
	t.StbState.SetCursor(t.CurLenW)
	t.StbState.SetSelectEnd(t.CurLenW)
	t.StbState.SetHasPreferredX(false)
}

func (t *TextEditState) OnKeyPressed(key int) {
	t.StbState.Key(t, key)
	t.CursorFollow = true
	t.CursorAnimReset()
}

func (t *TextEditState) GetChar(idx int) rune {
	return t.Text[idx]
}

func (t *TextEditState) GetWidth(line_start_idx, char_idx int) float64 {
	ctx := t.Ctx

	const STB_TEXTEDIT_GETWIDTH_NEWLINE = -1
	c := t.Text[line_start_idx+char_idx]
	if c == '\n' {
		return STB_TEXTEDIT_GETWIDTH_NEWLINE
	}
	return ctx.Font.GetCharAdvance(c) * (ctx.FontSize / ctx.Font.FontSize)
}

func (t *TextEditState) InsertChars(pos int, new_text []rune) bool {
	t.Text = append(t.Text[:pos], append(new_text, t.Text[pos:]...)...)
	t.CurLenW += len(new_text)
	t.CurLenA += TextCountUtf8BytesFromStr(new_text)
	return true
}

func (t *TextEditState) DeleteChars(pos, n int) {
	// We maintain our buffer length in both UTF-8 and wchar formats
	t.CurLenA -= TextCountUtf8BytesFromStr(t.Text[pos:])
	t.CurLenW -= n

	// Offset remaining text
	copy(t.Text[pos:pos+n], t.Text[pos+n:])
	t.Text = t.Text[:len(t.Text)-n]
}

func (t *TextEditState) LayoutRow(r *stbte.TextEditRow, line_start_idx int) {
	ctx := t.Ctx
	text := t.Text
	size, text_remaining, _ := ctx.InputTextCalcTextSizeW(text[line_start_idx:], true)
	r.SetX0(0.0)
	r.SetX1(size.X)
	r.SetBaselineYDelta(size.Y)
	r.SetYMin(0)
	r.SetYMax(size.Y)
	r.SetNumChars(text_remaining)
}

func (t *TextEditState) Len() int {
	return t.CurLenW
}

func (t *TextEditState) isSeparator(c rune) bool {
	switch c {
	case ',', ';', '(', ')', '{', '}', '[', ']', '|':
		return true
	}
	return CharIsSpace(c)
}

func (t *TextEditState) isWordBoundaryFromRight(idx int) bool {
	if t.isSeparator(t.Text[idx-1]) && !t.isSeparator(t.Text[idx]) {
		return idx > 0
	}
	return true
}

func (t *TextEditState) isWordBoundaryFromLeft(idx int) bool {
	if !t.isSeparator(t.Text[idx-1]) && t.isSeparator(t.Text[idx]) {
		return idx > 0
	}
	return true
}

func (t *TextEditState) MoveWordLeft(idx int) int {
	idx--
	for idx >= 0 && !t.isWordBoundaryFromRight(idx) {
		idx--
	}

	if idx < 0 {
		return 0
	}
	return idx
}

func (t *TextEditState) MoveWordRight(idx int) int {
	idx++
	length := t.CurLenW
	for idx < length && !t.isWordBoundaryFromLeft(idx) {
		idx++
	}

	if idx > length {
		return length
	}
	return idx
}

func (t *TextEditState) HasSelection() bool {
	return t.StbState.SelectStart() != t.StbState.SelectEnd()
}

func (t *TextEditState) ClearSelection() {
	cursor := t.StbState.Cursor()
	t.StbState.SetSelectStart(cursor)
	t.StbState.SetSelectEnd(cursor)
}

// Return false to discard a character.
func (c *Context) InputTextFilterCharacter(p_char *rune, flags InputTextFlags, callback TextEditCallback) bool {
	ch := *p_char
	if ch < 128 && ch != ' ' && !unicode.IsPrint(ch&0xff) {
		pass := false
		if ch == '\n' && flags&InputTextFlagsMultiline != 0 {
			pass = true
		}
		if ch == '\t' && flags&InputTextFlagsAllowTabInput != 0 {
			pass = true
		}
		if !pass {
			return false
		}
	}

	// Filter private Unicode range. I don't imagine anybody would want to input them. GLFW on OSX seems to send private characters for special keys like arrow keys.
	if ch >= 0xE000 && ch <= 0xF8FF {
		return false
	}

	if flags&(InputTextFlagsCharsDecimal|InputTextFlagsCharsHexadecimal|InputTextFlagsCharsUppercase|InputTextFlagsCharsNoBlank|InputTextFlagsCharsScientific) != 0 {
		if flags&InputTextFlagsCharsDecimal != 0 {
			if !(ch >= '0' && ch <= '9') && (ch != '.') && (ch != '-') && (ch != '+') && (ch != '*') && (ch != '/') {
				return false
			}
		}

		if flags&InputTextFlagsCharsScientific != 0 {
			if !(ch >= '0' && ch <= '9') && (ch != '.') && (ch != '-') && (ch != '+') && (ch != '*') && (ch != '/') && (ch != 'e') && (ch != 'E') {
				return false
			}
		}

		if flags&InputTextFlagsCharsHexadecimal != 0 {
			if !(ch >= '0' && ch <= '9') && !(ch >= 'a' && ch <= 'f') && !(ch >= 'A' && ch <= 'F') {
				return false
			}
		}

		if flags&InputTextFlagsCharsUppercase != 0 {
			if ch >= 'a' && ch <= 'z' {
				ch += 'A' - 'a'
				*p_char = ch
			}
		}

		if flags&InputTextFlagsCharsNoBlank != 0 {
			if CharIsSpace(ch) {
				return false
			}
		}
	}

	if flags&InputTextFlagsCallbackCharFilter != 0 {
		callback_data := TextEditCallbackData{
			EventFlag: InputTextFlagsCallbackCharFilter,
			EventChar: ch,
			Flags:     flags,
		}
		if callback(&callback_data) != 0 {
			return false
		}
		*p_char = callback_data.EventChar
		if callback_data.EventChar == 0 {
			return false
		}
	}

	return true
}

func (c *Context) InputTextCalcTextSizeW(text []rune, stop_on_new_line bool) (text_size f64.Vec2, remaining int, out_offset f64.Vec2) {
	font := c.Font
	line_height := c.FontSize
	scale := line_height / font.FontSize

	line_width := 0.0
	s := 0
	for ; s < len(text); s++ {
		c := text[s]
		if c == '\n' {
			text_size.X = math.Max(text_size.X, line_width)
			text_size.Y += line_height
			line_width = 0.0
			if stop_on_new_line {
				break
			}
			continue
		}
		if c == '\r' {
			continue
		}
		char_width := font.GetCharAdvance(c) * scale
		line_width += char_width
	}

	if text_size.X < line_width {
		text_size.X = line_width
	}

	// offset allow for the possibility of sitting after a trailing \n
	out_offset = f64.Vec2{line_width, text_size.Y + line_height}

	// whereas size.y will ignore the trailing \n
	if line_width > 0 || text_size.Y == 0.0 {
		text_size.Y += line_height
	}

	remaining = s
	return
}

func InputTextCalcTextLenAndLineCount(text string) int {
	line_count := 0

	// We are only matching for \n so we can ignore UTF-8 decoding
	for _, ch := range text {
		if ch == '\n' {
			line_count++
		}
	}
	ch := text[len(text)-1]
	if ch != '\n' && ch != '\r' {
		line_count++
	}
	return line_count
}

func TextCountUtf8BytesFromStr(r []rune) int {
	n := 0
	for i := range r {
		n += utf8.RuneLen(r[i])
	}
	return n
}

func (c *Context) TextDisabled(format string, args ...interface{}) {
	c.PushStyleColorV4(ColText, c.Style.Colors[ColTextDisabled])
	c.Text(format, args...)
	c.PopStyleColor()
}

func (c *Context) TextColored(col color.RGBA, format string, args ...interface{}) {
	c.PushStyleColor(ColText, col)
	c.Text(format, args...)
	c.PopStyleColor()
}

// Add a label+text combo aligned to other label+value widgets
func (c *Context) LabelText(label, format string, args ...interface{}) {
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return
	}

	style := &c.Style
	w := c.CalcItemWidth()

	label_size := c.CalcTextSizeEx(label, true, -1)
	value_bb := f64.Rectangle{
		window.DC.CursorPos,
		window.DC.CursorPos.Add(f64.Vec2{w, label_size.Y + style.FramePadding.Y*2}),
	}
	label_bb := f64.Vec2{w, style.FramePadding.Y * 2}
	if label_size.X > 0 {
		label_bb.X += style.ItemInnerSpacing.X
	}
	total_bb := f64.Rectangle{
		window.DC.CursorPos,
		window.DC.CursorPos.Add(label_bb).Add(label_size),
	}
	c.ItemSizeBBEx(total_bb, style.FramePadding.Y)
	if !c.ItemAdd(total_bb, 0) {
		return
	}

	// Render
	value_text := fmt.Sprintf(format, args...)
	c.RenderTextClippedEx(value_bb.Min, value_bb.Max, value_text, nil, f64.Vec2{0.0, 0.5}, nil)
	if label_size.X > 0.0 {
		c.RenderText(f64.Vec2{value_bb.Max.X + style.ItemInnerSpacing.X, value_bb.Min.Y + style.FramePadding.Y}, label)
	}
}