package imgui

import (
	"fmt"
	"image/color"
	"math"

	"github.com/qeedquan/go-media/math/f64"
)

type TreeNodeFlags int

const (
	TreeNodeFlagsSelected          TreeNodeFlags = 1 << 0  // Draw as selected
	TreeNodeFlagsFramed            TreeNodeFlags = 1 << 1  // Full colored frame (e.g. for CollapsingHeader)
	TreeNodeFlagsAllowItemOverlap  TreeNodeFlags = 1 << 2  // Hit testing to allow subsequent widgets to overlap this one
	TreeNodeFlagsNoTreePushOnOpen  TreeNodeFlags = 1 << 3  // Don't do a TreePush() when open (e.g. for CollapsingHeader) = no extra indent nor pushing on ID stack
	TreeNodeFlagsNoAutoOpenOnLog   TreeNodeFlags = 1 << 4  // Don't automatically and temporarily open node when Logging is active (by default logging will automatically open tree nodes)
	TreeNodeFlagsDefaultOpen       TreeNodeFlags = 1 << 5  // Default node to be open
	TreeNodeFlagsOpenOnDoubleClick TreeNodeFlags = 1 << 6  // Need double-click to open node
	TreeNodeFlagsOpenOnArrow       TreeNodeFlags = 1 << 7  // Only open when clicking on the arrow part. If TreeNodeFlagsOpenOnDoubleClick is also set single-click arrow or double-click all box to open.
	TreeNodeFlagsLeaf              TreeNodeFlags = 1 << 8  // No collapsing no arrow (use as a convenience for leaf nodes).
	TreeNodeFlagsBullet            TreeNodeFlags = 1 << 9  // Display a bullet instead of arrow
	TreeNodeFlagsFramePadding      TreeNodeFlags = 1 << 10 // Use FramePadding (even for an unframed text node) to vertically align text baseline to regular widget height. Equivalent to calling AlignTextToFramePadding().
	//ImGuITreeNodeFlags_SpanAllAvailWidth TreeNodeFlags = 1 << 11  // FIXME: TODO: Extend hit box horizontally even if not framed
	//TreeNodeFlagsNoScrollOnOpen   TreeNodeFlags  = 1 << 12  // FIXME: TODO: Disable automatic scroll on TreePop() if node got just open and contents is not visible
	TreeNodeFlagsNavLeftJumpsBackHere TreeNodeFlags = 1 << 13 // (WIP) Nav: left direction may move to this TreeNode() from any of its child (items submitted between TreeNode and TreePop)
	TreeNodeFlagsCollapsingHeader     TreeNodeFlags = TreeNodeFlagsFramed | TreeNodeFlagsNoAutoOpenOnLog
)

type ItemHoveredDataBackup struct {
	LastItemId          ID
	LastItemStatusFlags ItemStatusFlags
	LastItemRect        f64.Rectangle
	LastItemDisplayRect f64.Rectangle
}

func (b *ItemHoveredDataBackup) Backup(c *Context) {
	window := c.CurrentWindow
	b.LastItemId = window.DC.LastItemId
	b.LastItemStatusFlags = window.DC.LastItemStatusFlags
	b.LastItemRect = window.DC.LastItemRect
	b.LastItemDisplayRect = window.DC.LastItemDisplayRect
}

func (b *ItemHoveredDataBackup) Restore(c *Context) {
	window := c.CurrentWindow
	window.DC.LastItemId = b.LastItemId
	window.DC.LastItemStatusFlags = b.LastItemStatusFlags
	window.DC.LastItemRect = b.LastItemRect
	window.DC.LastItemDisplayRect = b.LastItemDisplayRect
}

func (c *Context) CollapsingHeader(label string) bool {
	return c.CollapsingHeaderEx(label, 0)
}

func (c *Context) CollapsingHeaderEx(label string, flags TreeNodeFlags) bool {
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return false
	}
	return c.TreeNodeBehavior(window.GetID(label), flags|TreeNodeFlagsCollapsingHeader|TreeNodeFlagsNoTreePushOnOpen, label)
}

func (c *Context) CollapsingHeaderOpen(label string, p_open *bool) bool {
	return c.CollapsingHeaderOpenEx(label, p_open, 0)
}

func (c *Context) CollapsingHeaderOpenEx(label string, p_open *bool, flags TreeNodeFlags) bool {
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	if p_open != nil && !*p_open {
		return false
	}

	id := window.GetID(label)
	if p_open != nil {
		flags |= TreeNodeFlagsAllowItemOverlap
	}
	is_open := c.TreeNodeBehavior(id, flags|TreeNodeFlagsCollapsingHeader|TreeNodeFlagsNoTreePushOnOpen, label)
	if p_open != nil {
		// Create a small overlapping close button // FIXME: We can evolve this into user accessible helpers to add extra buttons on title bars, headers, etc.
		button_sz := c.FontSize * 0.5
		pos := f64.Vec2{
			math.Min(window.DC.LastItemRect.Max.X, window.ClipRect.Max.X) - c.Style.FramePadding.X - button_sz,
			window.DC.LastItemRect.Min.Y + c.Style.FramePadding.Y + button_sz,
		}
		var last_item_backup ItemHoveredDataBackup
		last_item_backup.Backup(c)
		if c.CloseButton(window.GetIntID(int(id+1)), pos, button_sz) {
			*p_open = false
		}
		last_item_backup.Restore(c)
	}
	return is_open
}

func (c *Context) TreeNode(label string) bool {
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return false
	}
	return c.TreeNodeBehavior(window.GetID(label), 0, label)
}

func (c *Context) TreeNodeStringID(str_id string, format string, args ...interface{}) bool {
	return c.TreeNodeStringIDEx(str_id, 0, format, args...)
}

func (c *Context) TreeNodeID(id ID, format string, args ...interface{}) bool {
	return c.TreeNodeIDEx(id, 0, format, args...)
}

func (c *Context) TreeNodeStringIDEx(str_id string, flags TreeNodeFlags, format string, args ...interface{}) bool {
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return false
	}
	label := fmt.Sprintf(format, args...)
	return c.TreeNodeBehavior(window.GetID(str_id), flags, label)
}

func (c *Context) TreeNodeIDEx(id ID, flags TreeNodeFlags, format string, args ...interface{}) bool {
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return false
	}
	label := fmt.Sprintf(format, args...)
	return c.TreeNodeBehavior(window.GetIntID(int(id)), flags, label)
}

func (c *Context) TreeNodeBehavior(id ID, flags TreeNodeFlags, label string) bool {
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	style := &c.Style
	display_frame := (flags & TreeNodeFlagsFramed) != 0
	padding := f64.Vec2{style.FramePadding.X, 0.0}
	if display_frame || flags&TreeNodeFlagsFramePadding != 0 {
		padding = style.FramePadding
	}

	label_end := c.FindRenderedTextEnd(label)
	label = label[:label_end]
	label_size := c.CalcTextSizeEx(label, false, -1)

	// We vertically grow up to current line height up the typical widget height.
	// Latch before ItemSize changes it
	text_base_offset_y := math.Max(padding.Y, window.DC.CurrentLineTextBaseOffset)
	frame_height := math.Max(math.Min(window.DC.CurrentLineHeight, c.FontSize+style.FramePadding.Y*2), label_size.Y+padding.Y*2)
	frame_bb := f64.Rectangle{
		window.DC.CursorPos,
		f64.Vec2{window.Pos.X + c.GetContentRegionMax().X, window.DC.CursorPos.Y + frame_height},
	}
	if display_frame {
		// Framed header expand a little outside the default padding
		frame_bb.Min.X -= float64(int(window.WindowPadding.X*0.5)) - 1
		frame_bb.Max.X += float64(int(window.WindowPadding.X*0.5)) - 1
	}

	// Collapser arrow width + Spacing
	text_offset_x := c.FontSize
	if display_frame {
		text_offset_x += padding.X * 3
	} else {
		text_offset_x += padding.X * 2
	}
	// Include collapser
	text_width := c.FontSize
	if label_size.X > 0.0 {
		text_width += label_size.X + padding.X*2
	}
	c.ItemSizeEx(f64.Vec2{text_width, frame_height}, text_base_offset_y)

	// For regular tree nodes, we arbitrary allow to click past 2 worth of ItemSpacing
	// (Ideally we'd want to add a flag for the user to specify if we want the hit test to be done up to the right side of the content or not)
	var interact_bb f64.Rectangle
	if display_frame {
		interact_bb = frame_bb
	} else {
		interact_bb = f64.Rect(frame_bb.Min.X, frame_bb.Min.Y, frame_bb.Min.X+text_width+style.ItemSpacing.X*2, frame_bb.Max.Y)
	}
	is_open := c.TreeNodeBehaviorIsOpen(id, flags)

	// Store a flag for the current depth to tell if we will allow closing this node when navigating one of its child.
	// For this purpose we essentially compare if g.NavIdIsAlive went from 0 to 1 between TreeNode() and TreePop().
	// This is currently only support 32 level deep and we are fine with (1 << Depth) overflowing into a zero.
	if is_open && !c.NavIdIsAlive && flags&TreeNodeFlagsNavLeftJumpsBackHere != 0 && flags&TreeNodeFlagsNoTreePushOnOpen == 0 {
		window.DC.TreeDepthMayJumpToParentOnPop |= (1 << uint(window.DC.TreeDepth))
	}

	item_add := c.ItemAdd(interact_bb, id)
	window.DC.LastItemStatusFlags |= ItemStatusFlagsHasDisplayRect
	window.DC.LastItemDisplayRect = frame_bb

	if !item_add {
		if is_open && flags&TreeNodeFlagsNoTreePushOnOpen == 0 {
			c.TreePushRawID(id)
		}
		return is_open
	}

	// Flags that affects opening behavior:
	// - 0(default) ..................... single-click anywhere to open
	// - OpenOnDoubleClick .............. double-click anywhere to open
	// - OpenOnArrow .................... single-click on arrow to open
	// - OpenOnDoubleClick|OpenOnArrow .. single-click on arrow or double-click anywhere to open
	button_flags := ButtonFlagsNoKeyModifiers
	if flags&TreeNodeFlagsAllowItemOverlap != 0 {
		button_flags |= ButtonFlagsAllowItemOverlap
	}

	if flags&TreeNodeFlagsLeaf == 0 {
		button_flags |= ButtonFlagsPressedOnDragDropHold
	}
	if flags&TreeNodeFlagsOpenOnDoubleClick != 0 {
		button_flags |= ButtonFlagsPressedOnDoubleClick
		if flags&TreeNodeFlagsOpenOnArrow != 0 {
			button_flags |= ButtonFlagsPressedOnClickRelease
		}
	}
	hovered, held, pressed := c.ButtonBehavior(interact_bb, id, button_flags)
	if flags&TreeNodeFlagsLeaf == 0 {
		toggled := false
		if pressed {
			toggled = (flags&(TreeNodeFlagsOpenOnArrow|TreeNodeFlagsOpenOnDoubleClick)) == 0 || (c.NavActivateId == id)
			if flags&TreeNodeFlagsOpenOnArrow != 0 {
				if c.IsMouseHoveringRect(interact_bb.Min, f64.Vec2{interact_bb.Min.X + text_offset_x, interact_bb.Max.Y}) && !c.NavDisableMouseHover {
					toggled = true
				}
			}
			if flags&TreeNodeFlagsOpenOnDoubleClick != 0 {
				if c.IO.MouseDoubleClicked[0] {
					toggled = true
				}
			}
			// When using Drag and Drop "hold to open" we keep the node highlighted after opening, but never close it again.
			if c.DragDropActive && is_open {
				toggled = false
			}
		}

		if c.NavId == id && c.NavMoveRequest && c.NavMoveDir == DirLeft && is_open {
			toggled = true
			c.NavMoveRequestCancel()
		}
		// If there's something upcoming on the line we may want to give it the priority?
		if c.NavId == id && c.NavMoveRequest && c.NavMoveDir == DirRight && !is_open {
			toggled = true
			c.NavMoveRequestCancel()
		}
		if toggled {
			is_open = !is_open
			window.DC.StateStorage[id] = is_open
		}
	}
	if flags&TreeNodeFlagsAllowItemOverlap != 0 {
		c.SetItemAllowOverlap()
	}

	// Render
	var col color.RGBA
	switch {
	case held && hovered:
		col = c.GetColorFromStyle(ColHeaderActive)
	case hovered:
		col = c.GetColorFromStyle(ColHeaderHovered)
	default:
		col = c.GetColorFromStyle(ColHeader)
	}
	text_pos := frame_bb.Min.Add(f64.Vec2{text_offset_x, text_base_offset_y})
	if display_frame {
		// Framed type
		c.RenderFrameEx(frame_bb.Min, frame_bb.Max, col, true, style.FrameRounding)
		c.RenderNavHighlightEx(frame_bb, id, NavHighlightFlagsTypeThin)
		dir := DirRight
		if is_open {
			dir = DirDown
		}
		c.RenderArrowEx(frame_bb.Min.Add(f64.Vec2{padding.X, text_base_offset_y}), dir, 1.0)
		if c.LogEnabled {
			// NB: '##' is normally used to hide text (as a library-wide feature), so we need to specify the text range to make sure the ## aren't stripped out here.
			c.LogRenderedText(&text_pos, "\n##")
			c.RenderTextClipped(text_pos, frame_bb.Max, label, &label_size)
			c.LogRenderedText(&text_pos, "#")
		} else {
			c.RenderTextClipped(text_pos, frame_bb.Max, label, &label_size)
		}
	} else {
		// Unframed typed for tree nodes
		if hovered || flags&TreeNodeFlagsSelected != 0 {
			c.RenderFrameEx(frame_bb.Min, frame_bb.Max, col, false, 0)
			c.RenderNavHighlightEx(frame_bb, id, NavHighlightFlagsTypeThin)
		}

		if flags&TreeNodeFlagsBullet != 0 {
			c.RenderBullet(frame_bb.Min.Add(f64.Vec2{text_offset_x * 0.5, c.FontSize*0.50 + text_base_offset_y}))
		} else if flags&TreeNodeFlagsLeaf == 0 {
			dir := DirRight
			if is_open {
				dir = DirDown
			}
			c.RenderArrowEx(frame_bb.Min.Add(f64.Vec2{padding.X, c.FontSize*0.15 + text_base_offset_y}), dir, 0.70)
		}
		if c.LogEnabled {
			c.LogRenderedText(&text_pos, ">")
		}
		c.RenderTextEx(text_pos, label, false)
	}

	if is_open && flags&TreeNodeFlagsNoTreePushOnOpen == 0 {
		c.TreePushRawID(id)
	}

	return is_open
}

func (c *Context) TreeNodeBehaviorIsOpen(id ID, flags TreeNodeFlags) bool {
	if flags&TreeNodeFlagsLeaf != 0 {
		return true
	}

	// We only write to the tree storage if the user clicks (or explicitly use SetNextTreeNode*** functions)
	window := c.CurrentWindow
	var is_open bool
	if c.NextTreeNodeOpenCond != 0 {
		if c.NextTreeNodeOpenCond&CondAlways != 0 {
			is_open = c.NextTreeNodeOpenVal
			window.DC.StateStorage[id] = is_open
		} else {
			// We treat ImGuiCond_Once and ImGuiCond_FirstUseEver the same because tree node state are not saved persistently.
			stored_value, found := window.DC.StateStorage[id]
			if !found {
				is_open = c.NextTreeNodeOpenVal
				window.DC.StateStorage[id] = is_open
			} else {
				is_open = stored_value != 0
			}
		}
		c.NextTreeNodeOpenCond = 0
	} else {
		stored_value, found := window.DC.StateStorage[id]
		if !found {
			stored_value = false
			if flags&TreeNodeFlagsDefaultOpen != 0 {
				stored_value = true
			}
		}
		is_open = stored_value.(bool)
	}

	// When logging is enabled, we automatically expand tree nodes (but *NOT* collapsing headers.. seems like sensible behavior).
	// NB- If we are above max depth we still allow manually opened nodes to be logged.
	if c.LogEnabled && flags&TreeNodeFlagsNoAutoOpenOnLog == 0 && window.DC.TreeDepth < c.LogAutoExpandMaxDepth {
		is_open = true
	}

	return is_open
}

func (c *Context) TreePushRawID(id ID) {
	window := c.GetCurrentWindow()
	c.Indent()
	window.DC.TreeDepth++
	window.IDStack = append(window.IDStack, id)
}

func (c *Context) TreePop() {
	window := c.CurrentWindow
	c.Unindent()

	window.DC.TreeDepth--
	if c.NavMoveDir == DirLeft && c.NavWindow == window && c.NavMoveRequestButNoResultYet() {
		if c.NavIdIsAlive && (window.DC.TreeDepthMayJumpToParentOnPop&(1<<uint(window.DC.TreeDepth)) != 0) {
			c.SetNavID(window.IDStack[len(window.IDStack)-1], c.NavLayer)
			c.NavMoveRequestCancel()
		}
	}
	window.DC.TreeDepthMayJumpToParentOnPop &= (1 << uint(window.DC.TreeDepth)) - 1
	c.PopID()
}

func (c *Context) GetTreeNodeToLabelSpacing() float64 {
	return c.FontSize + (c.Style.FramePadding.X * 2.0)
}