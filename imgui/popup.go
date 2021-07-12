package imgui

import (
	"fmt"
	"math"

	"github.com/qeedquan/go-media/math/f64"
)

type PopupPositionPolicy int

const (
	PopupPositionPolicyDefault PopupPositionPolicy = iota
	PopupPositionPolicyComboBox
)

func (c *Context) OpenPopup(str_id string) {
	c.OpenPopupEx(c.CurrentWindow.GetID(str_id))
}

// Mark popup as open (toggle toward open state).
// Popups are closed when user click outside, or activate a pressable item, or CloseCurrentPopup() is called within a BeginPopup()/EndPopup() block.
// Popup identifiers are relative to the current ID-stack (so OpenPopup and BeginPopup needs to be at the same level).
// One open popup per level of the popup hierarchy (NB: when assigning we reset the Window member of ImGuiPopupRef to NULL)
func (c *Context) OpenPopupEx(id ID) {
	parent_window := c.CurrentWindow
	current_stack_size := len(c.CurrentPopupStack)
	// Tagged as new ref as Window will be set back to NULL if we write this into OpenPopupStack.
	var popup_ref PopupRef
	popup_ref.PopupId = id
	popup_ref.Window = nil
	popup_ref.ParentWindow = parent_window
	popup_ref.OpenFrameCount = c.FrameCount
	popup_ref.OpenParentId = parent_window.IDStack[len(parent_window.IDStack)-1]
	popup_ref.OpenMousePos = c.IO.MousePos
	popup_ref.OpenPopupPos = c.NavCalcPreferredRefPos()

	if len(c.OpenPopupStack) < current_stack_size+1 {
		c.OpenPopupStack = append(c.OpenPopupStack, popup_ref)
	} else {
		// Close child popups if any
		c.OpenPopupStack = c.OpenPopupStack[:current_stack_size+1]

		// Gently handle the user mistakenly calling OpenPopup() every frame. It is a programming mistake! However, if we were to run the regular code path, the ui
		// would become completely unusable because the popup will always be in hidden-while-calculating-size state _while_ claiming focus. Which would be a very confusing
		// situation for the programmer. Instead, we silently allow the popup to proceed, it will keep reappearing and the programming error will be more obvious to understand.
		if c.OpenPopupStack[current_stack_size].PopupId == id && c.OpenPopupStack[current_stack_size].OpenFrameCount == c.FrameCount-1 {
			c.OpenPopupStack[current_stack_size].OpenFrameCount = popup_ref.OpenFrameCount
		} else {
			c.OpenPopupStack[current_stack_size] = popup_ref
		}

		// When reopening a popup we first refocus its parent, otherwise if its parent is itself a popup it would get closed by ClosePopupsOverWindow().
		// This is equivalent to what ClosePopupToLevel() does.
		//if (g.OpenPopupStack[current_stack_size].PopupId == id)
		//    FocusWindow(parent_window);
	}
}

func (c *Context) IsPopupOpenName(str_id string) bool {
	return len(c.OpenPopupStack) > len(c.CurrentPopupStack) && c.OpenPopupStack[len(c.CurrentPopupStack)].PopupId == c.CurrentWindow.GetID(str_id)
}

func (c *Context) IsPopupOpenID(id ID) bool {
	return len(c.OpenPopupStack) > len(c.CurrentPopupStack) && c.OpenPopupStack[len(c.CurrentPopupStack)].PopupId == id
}

func (c *Context) EndPopup() {
	assert(c.CurrentWindow.Flags&WindowFlagsPopup != 0)
	assert(len(c.CurrentPopupStack) > 0)

	// Make all menus and popups wrap around for now, may need to expose that policy.
	c.NavProcessMoveRequestWrapAround(c.CurrentWindow)
	c.End()
}

func (c *Context) ClosePopupsOverWindow(ref_window *Window) {
	if len(c.OpenPopupStack) == 0 {
		return
	}

	// When popups are stacked, clicking on a lower level popups puts focus back to it and close popups above it.
	// Don't close our own child popup windows.
	var n int
	if ref_window != nil {
		for ; n < len(c.OpenPopupStack); n++ {
			popup := &c.OpenPopupStack[n]
			if popup.Window == nil {
				continue
			}
			assert((popup.Window.Flags & WindowFlagsPopup) != 0)
			if popup.Window.Flags&WindowFlagsChildWindow != 0 {
				continue
			}

			// Trim the stack if popups are not direct descendant of the reference window (which is often the NavWindow)
			has_focus := false
			for m := n; m < len(c.OpenPopupStack) && !has_focus; m++ {
				has_focus = c.OpenPopupStack[m].Window != nil && c.OpenPopupStack[m].Window.RootWindow == ref_window.RootWindow
			}
			if !has_focus {
				break
			}
		}
	}

	// This test is not required but it allows to set a convenient breakpoint on the block below
	if n < len(c.OpenPopupStack) {
		c.ClosePopupToLevel(n)
	}

}

func (c *Context) ClosePopupToLevel(remaining int) {
	var focus_window *Window
	if remaining > 0 {
		focus_window = c.OpenPopupStack[remaining-1].Window
	} else {
		focus_window = c.OpenPopupStack[0].ParentWindow
	}

	if c.NavLayer == 0 {
		focus_window = c.NavRestoreLastChildNavWindow(focus_window)
	}
	c.FocusWindow(focus_window)
	focus_window.DC.NavHideHighlightOneFrame = true
	c.OpenPopupStack = c.OpenPopupStack[:remaining]
}

func (c *Context) CloseCurrentPopup() {
	popup_idx := len(c.CurrentPopupStack) - 1
	if popup_idx < 0 || popup_idx >= len(c.OpenPopupStack) || c.CurrentPopupStack[popup_idx].PopupId != c.OpenPopupStack[popup_idx].PopupId {
		return
	}
	for popup_idx > 0 && c.OpenPopupStack[popup_idx].Window != nil && c.OpenPopupStack[popup_idx].Window.Flags&WindowFlagsChildMenu != 0 {
		popup_idx--
	}
	c.ClosePopupToLevel(popup_idx)
}

func (c *Context) BeginPopupEx(id ID, extra_flags WindowFlags) bool {
	if !c.IsPopupOpenID(id) {
		// We behave like Begin() and need to consume those values
		c.NextWindowData.Clear()
		return false
	}

	var name string
	if extra_flags&WindowFlagsChildMenu != 0 {
		// Recycle windows based on depth
		name = fmt.Sprintf("##Menu_%02d", len(c.CurrentPopupStack))
	} else {
		// Not recycling, so we can close/open during the same frame
		name = fmt.Sprintf("##Popup_%08x", id)
	}

	is_open := c.BeginEx(name, nil, extra_flags|WindowFlagsPopup)
	// NB: Begin can return false when the popup is completely clipped (e.g. zero size display)
	if !is_open {
		c.EndPopup()
	}
	return is_open
}

func (c *Context) IsPopupOpen(id ID) bool {
	return len(c.OpenPopupStack) > len(c.CurrentPopupStack) && c.OpenPopupStack[len(c.CurrentPopupStack)].PopupId == id
}

func (c *Context) CalcMaxPopupHeightFromItemCount(items_count int) float64 {
	if items_count <= 0 {
		return math.MaxFloat32
	}
	return (c.FontSize+c.Style.ItemSpacing.Y)*float64(items_count) - c.Style.ItemSpacing.Y + (c.Style.WindowPadding.Y * 2)
}

func (c *Context) BeginPopup(str_id string) bool {
	return c.BeginPopupWindow(str_id, 0)
}

func (c *Context) BeginPopupWindow(str_id string, flags WindowFlags) bool {
	// Early out for performance
	if len(c.OpenPopupStack) <= len(c.CurrentPopupStack) {
		// We behave like Begin() and need to consume those values
		c.NextWindowData.Clear()
		return false
	}
	return c.BeginPopupEx(c.CurrentWindow.GetID(str_id), flags|WindowFlagsAlwaysAutoResize|WindowFlagsNoTitleBar|WindowFlagsNoSavedSettings)
}

func (c *Context) OpenPopupOnItemClick(str_id string, mouse_button int) bool {
	window := c.CurrentWindow
	if c.IsMouseReleased(mouse_button) && c.IsItemHoveredEx(HoveredFlagsAllowWhenBlockedByPopup) {
		// If user hasn't passed an ID, we can use the LastItemID. Using LastItemID as a Popup ID won't conflict!
		id := window.DC.LastItemId
		if str_id != "" {
			id = window.GetID(str_id)
		}
		// However, you cannot pass a NULL str_id if the last item has no identifier (e.g. a Text() item)
		assert(id != 0)
		c.OpenPopupEx(id)
		return true
	}
	return false
}

func (c *Context) BeginPopupContextWindow() bool {
	return c.BeginPopupContextWindowEx("", 1, true)
}

func (c *Context) BeginPopupContextWindowEx(str_id string, mouse_button int, also_over_items bool) bool {
	if str_id == "" {
		str_id = "window_context"
	}
	id := c.CurrentWindow.GetID(str_id)
	if c.IsMouseReleased(mouse_button) && c.IsWindowHoveredEx(HoveredFlagsAllowWhenBlockedByPopup) {
		if also_over_items || !c.IsAnyItemHovered() {
			c.OpenPopupEx(id)
		}
	}
	return c.BeginPopupEx(id, WindowFlagsAlwaysAutoResize|WindowFlagsNoTitleBar|WindowFlagsNoSavedSettings)
}

func (c *Context) BeginPopupContextItem() bool {
	return c.BeginPopupContextItemEx("", 1)
}

// This is a helper to handle the simplest case of associating one named popup to one given widget.
// You may want to handle this on user side if you have specific needs (e.g. tweaking IsItemHovered() parameters).
// You can pass a NULL str_id to use the identifier of the last item.
func (c *Context) BeginPopupContextItemEx(str_id string, mouse_button int) bool {
	window := c.CurrentWindow
	id := window.DC.LastItemId
	// If user hasn't passed an ID, we can use the LastItemID. Using LastItemID as a Popup ID won't conflict!
	if str_id != "" {
		id = window.GetID(str_id)
	}
	// However, you cannot pass a NULL str_id if the last item has no identifier (e.g. a Text() item)
	assert(id != 0)
	if c.IsMouseReleased(mouse_button) && c.IsItemHoveredEx(HoveredFlagsAllowWhenBlockedByPopup) {
		c.OpenPopupEx(id)
	}
	return c.BeginPopupEx(id, WindowFlagsAlwaysAutoResize|WindowFlagsNoTitleBar|WindowFlagsNoSavedSettings)
}

func (c *Context) BeginPopupModal(name string) bool {
	return c.BeginPopupModalEx(name, nil, 0)
}

func (c *Context) BeginPopupModalEx(name string, p_open *bool, flags WindowFlags) bool {
	window := c.CurrentWindow
	id := window.GetID(name)
	if !c.IsPopupOpen(id) {
		// We behave like Begin() and need to consume those values
		c.NextWindowData.Clear()
		return false
	}
	// Center modal windows by default
	// FIXME: Should test for (PosCond & window->SetWindowPosAllowFlags) with the upcoming window.
	if c.NextWindowData.PosCond == 0 {
		c.SetNextWindowPos(c.IO.DisplaySize.Scale(0.5), CondAppearing, f64.Vec2{0.5, 0.5})
	}

	is_open := c.BeginEx(name, p_open, flags|WindowFlagsPopup|WindowFlagsModal|WindowFlagsNoCollapse|WindowFlagsNoSavedSettings)
	// NB: is_open can be 'false' when the popup is completely clipped (e.g. zero size display)
	if !is_open || (p_open != nil && !*p_open) {
		c.EndPopup()
		if is_open {
			c.ClosePopup(id)
		}
		return false
	}

	return is_open
}

func (c *Context) ClosePopup(id ID) {
	if !c.IsPopupOpen(id) {
		return
	}
	c.ClosePopupToLevel(len(c.OpenPopupStack) - 1)
}
