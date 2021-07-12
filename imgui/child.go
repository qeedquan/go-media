package imgui

import (
	"fmt"
	"math"

	"github.com/qeedquan/go-media/math/f64"
)

func (c *Context) BeginChild(str_id string) bool {
	return c.BeginChildEx(str_id, f64.Vec2{0, 0}, false, 0)
}

func (c *Context) BeginChildEx(str_id string, size_arg f64.Vec2, border bool, extra_flags WindowFlags) bool {
	window := c.GetCurrentWindow()
	return c.beginChildEx(str_id, window.GetID(str_id), size_arg, border, extra_flags)
}

func (c *Context) BeginChildID(id ID) bool {
	return c.BeginChildIDEx(id, f64.Vec2{0, 0}, false, 0)
}

func (c *Context) BeginChildIDEx(id ID, size_arg f64.Vec2, border bool, extra_flags WindowFlags) bool {
	assert(id != 0)
	return c.beginChildEx("", id, size_arg, border, extra_flags)
}

func (c *Context) beginChildEx(name string, id ID, size_arg f64.Vec2, border bool, extra_flags WindowFlags) bool {
	parent_window := c.GetCurrentWindow()
	flags := WindowFlagsNoTitleBar | WindowFlagsNoResize | WindowFlagsNoSavedSettings | WindowFlagsChildWindow
	// Inherit the NoMove flag
	flags |= (parent_window.Flags & WindowFlagsNoMove)

	content_avail := c.GetContentRegionAvail()
	size := size_arg.Floor()
	auto_fit_axises := 0
	if size.X == 0 {
		auto_fit_axises |= (1 << uint(AxisX))
	}
	if size.Y == 0 {
		auto_fit_axises |= (1 << uint(AxisY))
	}
	if size.X <= 0.0 {
		// Arbitrary minimum child size (0.0f causing too much issues)
		size.X = math.Max(content_avail.X+size.X, 4.0)
	}
	if size.Y <= 0.0 {
		size.Y = math.Max(content_avail.Y+size.Y, 4.0)
	}

	backup_border_size := c.Style.ChildBorderSize
	if !border {
		c.Style.ChildBorderSize = 0
	}
	flags |= extra_flags

	var title string
	if name != "" {
		title = fmt.Sprintf("%s/%s", parent_window.Name, name)
	} else {
		title = fmt.Sprintf("%s/%08X", parent_window.Name, id)
	}

	c.SetNextWindowSize(size, 0)
	ret := c.BeginEx(title, nil, flags)
	child_window := c.GetCurrentWindow()
	child_window.ChildId = id
	child_window.AutoFitChildAxises = auto_fit_axises
	c.Style.ChildBorderSize = backup_border_size

	// Process navigation-in immediately so NavInit can run on first frame
	if flags&WindowFlagsNavFlattened == 0 && (child_window.DC.NavLayerActiveMask != 0 || child_window.DC.NavHasScroll) && c.NavActivateId == id {
		c.FocusWindow(child_window)
		c.NavInitWindow(child_window, false)
		// Steal ActiveId with a dummy id so that key-press won't activate child item
		c.SetActiveID(id+1, child_window)
		c.ActiveIdSource = InputSourceNav
	}

	return ret
}

func (c *Context) EndChild() {
	window := c.CurrentWindow
	// Mismatched BeginChild()/EndChild() callss
	assert(window.Flags&WindowFlagsChildWindow != 0)
	if window.BeginCount > 1 {
		c.End()
	} else {
		// When using auto-filling child window, we don't provide full width/height to ItemSize so that it doesn't feed back into automatic size-fitting.
		sz := c.GetWindowSize()
		if window.AutoFitChildAxises&(1<<uint(AxisX)) != 0 {
			sz.X = math.Max(4.0, sz.X)
		}
		if window.AutoFitChildAxises&(1<<uint(AxisY)) != 0 {
			sz.Y = math.Max(4.0, sz.Y)
		}
		c.End()

		parent_window := c.CurrentWindow
		bb := f64.Rectangle{
			parent_window.DC.CursorPos,
			parent_window.DC.CursorPos.Add(sz),
		}
		c.ItemSize(sz)

		if (window.DC.NavLayerActiveMask != 0 || window.DC.NavHasScroll) && window.Flags&WindowFlagsNavFlattened == 0 {
			c.ItemAdd(bb, window.ChildId)
			c.RenderNavHighlight(bb, window.ChildId)

			// When browsing a window that has no activable items (scroll only) we keep a highlight on the child
			if window.DC.NavLayerActiveMask == 0 && window == c.NavWindow {
				c.RenderNavHighlightEx(f64.Rectangle{bb.Min.Sub(f64.Vec2{2, 2}), bb.Max.Add(f64.Vec2{2, 2})}, c.NavId, NavHighlightFlagsTypeThin)
			}
		} else {
			// Not navigable into
			c.ItemAdd(bb, 0)
		}
	}
}

func (c *Context) BeginChildFrame(id ID, size f64.Vec2, extra_flags WindowFlags) bool {
	style := &c.Style
	c.PushStyleColorV4(ColChildBg, style.Colors[ColFrameBg])
	c.PushStyleVar(StyleVarChildRounding, style.FrameRounding)
	c.PushStyleVar(StyleVarChildBorderSize, style.FrameBorderSize)
	c.PushStyleVar(StyleVarWindowPadding, style.FramePadding)
	return c.beginChildEx("", id, size, true, WindowFlagsNoMove|WindowFlagsAlwaysUseWindowPadding|extra_flags)
}

func (c *Context) EndChildFrame() {
	c.EndChild()
	c.PopStyleVarN(3)
	c.PopStyleColor()
}
