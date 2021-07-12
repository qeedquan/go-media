package imgui

import (
	"math"

	"github.com/qeedquan/go-media/math/f64"
	"github.com/qeedquan/go-media/math/mathutil"
)

// Helper: Manually clip large list of items.
// If you are submitting lots of evenly spaced items and you have a random access to the list, you can perform coarse clipping based on visibility to save yourself from processing those items at all.
// The clipper calculates the range of visible items and advance the cursor to compensate for the non-visible items we have skipped.
// ImGui already clip items based on their bounds but it needs to measure text size to do so. Coarse clipping before submission makes this cost and your own data fetching/submission cost null.
// Usage:
//     ImGuiListClipper clipper(1000);  // we have 1000 elements, evenly spaced.
//     while (clipper.Step())
//         for (int i = clipper.DisplayStart; i < clipper.DisplayEnd; i++)
//             ImGui::Text("line number %d", i);
// - Step 0: the clipper let you process the first element, regardless of it being visible or not, so we can measure the element height (step skipped if we passed a known height as second arg to constructor).
// - Step 1: the clipper infer height from first element, calculate the actual range of elements to display, and position the cursor before the first element.
// - (Step 2: dummy step only required if an explicit items_height was passed to constructor or Begin() and user call Step(). Does nothing and switch to Step 3.)
// - Step 3: the clipper validate that we have reached the expected Y position (corresponding to element DisplayEnd), advance the cursor to the end of the list and then returns 'false' to end the loop.
type ListClipper struct {
	Ctx                                          *Context
	StartPosY                                    float64
	ItemsHeight                                  float64
	ItemsCount, StepNo, DisplayStart, DisplayEnd int
}

func (l *ListClipper) Init(ctx *Context, items_count int, items_height float64) {
	l.Ctx = ctx
	l.Begin(items_count, items_height)
}

func (l *ListClipper) Begin(items_count int, items_height float64) {
	l.StartPosY = l.Ctx.GetCursorPosY()
	l.ItemsHeight = items_height
	l.ItemsCount = items_count
	l.StepNo = 0
	l.DisplayEnd = -1
	l.DisplayStart = -1
	if l.ItemsHeight > 0 {
		// calculate how many to clip/display
		l.DisplayStart, l.DisplayEnd = l.Ctx.CalcListClipping(l.ItemsCount, l.ItemsHeight)
		if l.DisplayStart > 0 {
			// advance cursor
			l.Ctx.SetCursorPosYAndSetupDummyPrevLine(l.StartPosY+float64(l.DisplayStart)*l.ItemsHeight, l.ItemsHeight)
		}
		l.StepNo = 2
	}
}

func (l *ListClipper) Step() bool {
	if l.ItemsCount == 0 || l.Ctx.GetCurrentWindowRead().SkipItems {
		l.ItemsCount = 0
		return false
	}

	switch l.StepNo {
	case 0: // Step 0: the clipper let you process the first element, regardless of it being visible or not, so we can measure the element height.
		l.DisplayStart = 0
		l.DisplayEnd = 1
		l.StartPosY = l.Ctx.GetCursorPosY()
		l.StepNo = 1
		return true

	case 1: // Step 1: the clipper infer height from first element, calculate the actual range of elements to display, and position the cursor before the first element.
		if l.ItemsCount == 1 {
			l.ItemsCount = -1
			return false
		}
		items_height := l.Ctx.GetCursorPosY() - l.StartPosY
		// If this triggers, it means Item 0 hasn't moved the cursor vertically
		assert(items_height > 0.0)
		l.Begin(l.ItemsCount-1, items_height)
		l.DisplayStart++
		l.DisplayEnd++
		l.StepNo = 3
		return true

	case 2: // Step 2: dummy step only required if an explicit items_height was passed to constructor or Begin() and user still call Step(). Does nothing and switch to Step 3.
		assert(l.DisplayStart >= 0 && l.DisplayEnd >= 0)
		l.StepNo = 3
		return true

	case 3: // Step 3: the clipper validate that we have reached the expected Y position (corresponding to element DisplayEnd), advance the cursor to the end of the list and then returns 'false' to end the loop.
		l.End()
	}

	return false
}

func (l *ListClipper) End() {
	if l.ItemsCount < 0 {
		return
	}
	// In theory here we should assert that ImGui::GetCursorPosY() == StartPosY + DisplayEnd * ItemsHeight, but it feels saner to just seek at the end and not assert/crash the user.
	if l.ItemsCount < math.MaxInt32 {
		// advance cursor
		l.Ctx.SetCursorPosYAndSetupDummyPrevLine(l.StartPosY+float64(l.ItemsCount)*l.ItemsHeight, l.ItemsHeight)
	}
	l.ItemsCount = -1
	l.StepNo = 3
}

func (c *Context) SetCursorPosYAndSetupDummyPrevLine(pos_y, line_height float64) {
	// Set cursor position and a few other things so that SetScrollHere() and Columns() can work when seeking cursor.
	// FIXME: It is problematic that we have to do that here, because custom/equivalent end-user code would stumble on the same issue.
	// The clipper should probably have a 4th step to display the last item in a regular manner.
	c.SetCursorPosY(pos_y)
	window := c.GetCurrentWindow()
	// Setting those fields so that SetScrollHere() can properly function after the end of our clipper usage.
	window.DC.CursorPosPrevLine.Y = window.DC.CursorPos.Y - line_height

	// If we end up needing more accurate data (to e.g. use SameLine) we may as well make the clipper have a fourth step to let user process and display the last item in their list.
	window.DC.PrevLineHeight = (line_height - c.Style.ItemSpacing.Y)

	if window.DC.ColumnsSet != nil {
		// Setting this so that cell Y position are set properly
		window.DC.ColumnsSet.LineMinY = window.DC.CursorPos.Y
	}
}

func (c *Context) ListBox(label string, current_item *int, items []string) bool {
	return c.ListBoxEx(label, current_item, items, -1)
}

func (c *Context) ListBoxEx(label string, current_item *int, items []string, height_in_items int) bool {
	items_getter := func(idx int) (string, bool) {
		if 0 <= idx && idx < len(items) {
			return items[idx], true
		}
		return "", false
	}
	return c.ListBoxItems(label, current_item, items_getter, len(items), height_in_items)
}

func (c *Context) ListBoxItems(label string, current_item *int, items_getter func(idx int) (string, bool), items_count int, height_in_items int) bool {
	if !c.ListBoxHeaderItems(label, items_count, height_in_items) {
		return false
	}

	// Assume all items have even height (= 1 line of text). If you need items of different or variable sizes you can create a custom version of ListBox() in your code without using the clipper.
	value_changed := false
	var clipper ListClipper
	clipper.Init(c, items_count, c.GetTextLineHeightWithSpacing())
	for clipper.Step() {
		for i := clipper.DisplayStart; i < clipper.DisplayEnd; i++ {
			item_selected := i == *current_item
			item_text, found := items_getter(i)
			if !found {
				item_text = "*Unknown item*"
			}

			c.PushID(ID(i))
			if c.SelectableEx(item_text, item_selected, 0, f64.Vec2{}) {
				*current_item = i
				value_changed = true
			}
			if item_selected {
				c.SetItemDefaultFocus()
			}
			c.PopID()
		}
	}

	c.ListBoxFooter()
	return value_changed
}

// Helper to calculate the size of a listbox and display a label on the right.
// Tip: To have a list filling the entire window width, PushItemWidth(-1) and pass an empty label "##empty"
func (c *Context) ListBoxHeader(label string, size_arg f64.Vec2) bool {
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	style := &c.Style
	id := c.GetStringID(label)
	label_size := c.CalcTextSizeEx(label, true, -1)

	// Size default to hold ~7 items. Fractional number of items helps seeing that we can scroll down/up without looking at scrollbar.
	size := c.CalcItemSize(size_arg, c.CalcItemWidth(), c.GetTextLineHeightWithSpacing()*7.4+style.ItemSpacing.Y)
	frame_size := f64.Vec2{size.X, math.Max(size.Y, label_size.Y)}
	frame_bb := f64.Rectangle{window.DC.CursorPos, window.DC.CursorPos.Add(frame_size)}
	label_bb := f64.Vec2{}
	if label_size.X > 0 {
		label_bb.X = style.ItemInnerSpacing.X + label_size.X
	}
	bb := f64.Rectangle{
		frame_bb.Min,
		frame_bb.Max.Add(label_bb),
	}
	// Forward storage for ListBoxFooter.. dodgy.
	window.DC.LastItemRect = bb

	c.BeginGroup()
	if label_size.X > 0 {
		c.RenderText(f64.Vec2{frame_bb.Max.X + style.ItemInnerSpacing.X, frame_bb.Min.Y + style.FramePadding.Y}, label)
	}

	c.BeginChildFrame(id, frame_bb.Size(), 0)
	return true
}

func (c *Context) ListBoxHeaderItems(label string, items_count, height_in_items int) bool {
	style := &c.Style

	// Size default to hold ~7 items. Fractional number of items helps seeing that we can scroll down/up without looking at scrollbar.
	// However we don't add +0.40f if items_count <= height_in_items. It is slightly dodgy, because it means a dynamic list of items will make the widget resize occasionally when it crosses that size.
	// I am expecting that someone will come and complain about this behavior in a remote future, then we can advise on a better solution.
	if height_in_items < 0 {
		height_in_items = mathutil.Min(items_count, 7)
	}
	height_in_items_f := float64(height_in_items)
	if height_in_items < items_count {
		height_in_items_f += 0.4
	}

	// We include ItemSpacing.y so that a list sized for the exact number of items doesn't make a scrollbar appears. We could also enforce that by passing a flag to BeginChild().
	size := f64.Vec2{0, c.GetTextLineHeightWithSpacing()*height_in_items_f + style.ItemSpacing.Y}
	return c.ListBoxHeader(label, size)
}

func (c *Context) ListBoxFooter() {
	parent_window := c.GetCurrentWindow().ParentWindow
	bb := parent_window.DC.LastItemRect
	style := &c.Style

	c.EndChildFrame()

	// Redeclare item size so that it includes the label (we have stored the full size in LastItemRect)
	// We call SameLine() to restore DC.CurrentLine* data
	c.SameLine()
	parent_window.DC.CursorPos = bb.Min
	c.ItemSizeBBEx(bb, style.FramePadding.Y)
	c.EndGroup()
}

// Helper to calculate coarse clipping of large list of evenly sized items.
// NB: Prefer using the ImGuiListClipper higher-level helper if you can! Read comments and instructions there on how those use this sort of pattern.
// NB: 'items_count' is only used to clamp the result, if you don't know your count you can use INT_MAX
func (c *Context) CalcListClipping(items_count int, items_height float64) (out_items_display_start int, out_items_display_end int) {
	window := c.CurrentWindow
	if c.LogEnabled {
		// If logging is active, do not perform any clipping
		out_items_display_start = 0
		out_items_display_end = items_count
		return
	}
	if window.SkipItems {
		out_items_display_start, out_items_display_end = 0, 0
		return
	}
	pos := window.DC.CursorPos
	start := int((window.ClipRect.Min.Y - pos.Y) / items_height)
	end := int((window.ClipRect.Max.Y - pos.Y) / items_height)
	// When performing a navigation request, ensure we have one item extra in the direction we are moving to
	if c.NavMoveRequest && c.NavMoveDir == DirUp {
		start--
	}
	if c.NavMoveRequest && c.NavMoveDir == DirDown {
		end++
	}
	out_items_display_start = mathutil.Clamp(start, 0, items_count)
	out_items_display_end = mathutil.Clamp(end+1, start, items_count)

	return
}