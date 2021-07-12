package imgui

import (
	"math"

	"github.com/qeedquan/go-media/math/f64"
)

type NavInput int

const (
	// Gamepad Mapping
	NavInputActivate    NavInput = iota // activate / open / toggle / tweak value       // e.g. Cross  (PS4), A (Xbox), A (Switch), Space (Keyboard)
	NavInputCancel                      // cancel / close / exit                        // e.g. Circle (PS4), B (Xbox), B (Switch), Escape (Keyboard)
	NavInputInput                       // text input / on-screen keyboard              // e.g. Triang.(PS4) Y (Xbox) X (Switch) Return (Keyboard)
	NavInputMenu                        // tap: toggle menu / hold: focus move resize // e.g. Square (PS4) X (Xbox) Y (Switch) Alt (Keyboard)
	NavInputDpadLeft                    // move / tweak / resize window (w/ PadMenu)    // e.g. D-pad Left/Right/Up/Down (Gamepads) Arrow keys (Keyboard)
	NavInputDpadRight                   //
	NavInputDpadUp                      //
	NavInputDpadDown                    //
	NavInputLStickLeft                  // scroll / move window (w/ PadMenu)            // e.g. Left Analog Stick Left/Right/Up/Down
	NavInputLStickRight                 //
	NavInputLStickUp                    //
	NavInputLStickDown                  //
	NavInputFocusPrev                   // next window (w/ PadMenu)                     // e.g. L1 or L2 (PS4) LB or LT (Xbox) L or ZL (Switch)
	NavInputFocusNext                   // prev window (w/ PadMenu)                     // e.g. R1 or R2 (PS4) RB or RT (Xbox) R or ZL (Switch)
	NavInputTweakSlow                   // slower tweaks                                // e.g. L1 or L2 (PS4) LB or LT (Xbox) L or ZL (Switch)
	NavInputTweakFast                   // faster tweaks                                // e.g. R1 or R2 (PS4) RB or RT (Xbox) R or ZL (Switch)

	// [Internal] Don't use directly! This is used internally to differentiate keyboard from gamepad inputs for behaviors that require to differentiate them.
	// Keyboard behavior that have no corresponding gamepad mapping (e.g. CTRL+TAB) will be directly reading from io.KeysDown[] instead of io.NavInputs[].
	NavInputKeyMenu_  // toggle menu                                  // = io.KeyAlt
	NavInputKeyLeft_  // move left                                    // = Arrow keys
	NavInputKeyRight_ // move right
	NavInputKeyUp_    // move up
	NavInputKeyDown_  // move down
	NavInputCOUNT
	NavInputInternalStart_ = NavInputKeyMenu_
)

type Dir int

const (
	DirNone  Dir = -1
	DirLeft  Dir = 0
	DirRight Dir = 1
	DirUp    Dir = 2
	DirDown  Dir = 3
	DirCOUNT Dir = 4
)

type Cond int

const (
	CondAlways       Cond = 1 << 0 // Set the variable
	CondOnce         Cond = 1 << 1 // Set the variable once per runtime session (only the first call with succeed)
	CondFirstUseEver Cond = 1 << 2 // Set the variable if the window has no saved data (if doesn't exist in the .ini file)
	CondAppearing    Cond = 1 << 3 // Set the variable if the window is appearing after being hidden/inactive (or the first time)
)

type NavHighlightFlags int

const (
	NavHighlightFlagsTypeDefault NavHighlightFlags = 1 << 0
	NavHighlightFlagsTypeThin    NavHighlightFlags = 1 << 1
	NavHighlightFlagsAlwaysDraw  NavHighlightFlags = 1 << 2
	NavHighlightFlagsNoRounding  NavHighlightFlags = 1 << 3
)

type NavDirSourceFlags int

const (
	NavDirSourceFlagsKeyboard  NavDirSourceFlags = 1 << 0
	NavDirSourceFlagsPadDPad   NavDirSourceFlags = 1 << 1
	NavDirSourceFlagsPadLStick NavDirSourceFlags = 1 << 2
)

type NavForward int

const (
	NavForwardNone NavForward = iota
	NavForwardForwardQueued
	NavForwardForwardActive
)

type NavMoveResult struct {
	ID         ID      // Best candidate
	ParentID   ID      // Best candidate window->IDStack.back() - to compare context
	Window     *Window // Best candidate window
	DistBox    float64 // Best candidate box distance to current NavId
	DistCenter float64 // Best candidate center distance to current NavId
	DistAxial  float64
	RectRel    f64.Rectangle // Best candidate bounding box in window relative space
}

func (c *Context) IsNavInputPressed(n NavInput, mode InputReadMode) bool {
	return c.GetNavInputAmount(n, mode) > 0
}

func (c *Context) GetNavInputAmount(n NavInput, mode InputReadMode) float64 {
	// Instant, read analog input (0.0f..1.0f, as provided by user)
	if mode == InputReadModeDown {
		return c.IO.NavInputs[n]
	}

	t := c.IO.NavInputsDownDuration[n]
	// Return 1.0f when just released, no repeat, ignore analog input.
	if t < 0 && mode == InputReadModeReleased {
		if c.IO.NavInputsDownDurationPrev[n] >= 0 {
			return 1
		}
		return 0
	}

	if t < 0 {
		return 0
	}

	// Return 1.0f when just pressed, no repeat, ignore analog input.
	if mode == InputReadModePressed {
		if t == 0 {
			return 1
		}
		return 0
	}

	if mode == InputReadModeRepeat {
		return float64(c.CalcTypematicPressedRepeatAmount(
			t,
			t-c.IO.DeltaTime,
			c.IO.KeyRepeatDelay*0.80,
			c.IO.KeyRepeatRate*0.80,
		))
	}

	if mode == InputReadModeRepeatSlow {
		return float64(c.CalcTypematicPressedRepeatAmount(
			t,
			t-c.IO.DeltaTime,
			c.IO.KeyRepeatDelay*1,
			c.IO.KeyRepeatRate*2,
		))
	}

	if mode == InputReadModeRepeatFast {
		return float64(c.CalcTypematicPressedRepeatAmount(
			t,
			t-c.IO.DeltaTime,
			c.IO.KeyRepeatDelay*0.80,
			c.IO.KeyRepeatRate*0.30,
		))
	}

	return 0
}

func (c *Context) GetNavInputAmount2d(dir_sources NavDirSourceFlags, mode InputReadMode) f64.Vec2 {
	return c.GetNavInputAmount2dEx(dir_sources, mode, 0, 0)
}

func (c *Context) GetNavInputAmount2dEx(dir_sources NavDirSourceFlags, mode InputReadMode, slow_factor, fast_factor float64) f64.Vec2 {
	delta := f64.Vec2{}
	if dir_sources&NavDirSourceFlagsKeyboard != 0 {
		right := c.GetNavInputAmount(NavInputKeyRight_, mode)
		left := c.GetNavInputAmount(NavInputKeyLeft_, mode)
		down := c.GetNavInputAmount(NavInputKeyRight_, mode)
		up := c.GetNavInputAmount(NavInputKeyRight_, mode)
		dir := f64.Vec2{right - left, down - up}
		delta = delta.Add(dir)
	}
	if dir_sources&NavDirSourceFlagsPadDPad != 0 {
		right := c.GetNavInputAmount(NavInputDpadRight, mode)
		left := c.GetNavInputAmount(NavInputDpadLeft, mode)
		down := c.GetNavInputAmount(NavInputDpadDown, mode)
		up := c.GetNavInputAmount(NavInputDpadUp, mode)
		dir := f64.Vec2{right - left, down - up}
		delta = delta.Add(dir)
	}
	if dir_sources&NavDirSourceFlagsPadLStick != 0 {
		right := c.GetNavInputAmount(NavInputLStickRight, mode)
		left := c.GetNavInputAmount(NavInputLStickLeft, mode)
		down := c.GetNavInputAmount(NavInputLStickDown, mode)
		up := c.GetNavInputAmount(NavInputLStickUp, mode)
		dir := f64.Vec2{right - left, down - up}
		delta = delta.Add(dir)
	}
	if slow_factor != 0.0 && c.IsNavInputDown(NavInputTweakSlow) {
		delta = delta.Scale(slow_factor)
	}
	if fast_factor != 0.0 && c.IsNavInputDown(NavInputTweakFast) {
		delta = delta.Scale(fast_factor)
	}
	return delta
}

// FIXME-OPT O(N)
func (c *Context) FindWindowIndex(window *Window) int {
	for i := len(c.Windows) - 1; i >= 0; i-- {
		if c.Windows[i] == window {
			return i
		}
	}
	return -1
}

// FIXME-OPT O(N)
func (c *Context) FindWindowNavigable(i_start, i_stop, dir int) *Window {
	for i := i_start; i >= 0 && i < len(c.Windows) && i != i_stop; i += dir {
		if c.IsWindowNavFocusable(c.Windows[i]) {
			return c.Windows[i]
		}
	}
	return nil
}

func (c *Context) IsWindowNavFocusable(window *Window) bool {
	return window.Active && window == window.RootWindowForTabbing && window.Flags&WindowFlagsNoNavFocus == 0 || window == c.NavWindow
}

func (c *Context) NavUpdateWindowingHighlightWindow(focus_change_dir int) {
	if c.NavWindowingTarget.Flags&WindowFlagsModal != 0 {
		return
	}

	i_current := c.FindWindowIndex(c.NavWindowingTarget)
	window_target := c.FindWindowNavigable(i_current+focus_change_dir, -1e6, focus_change_dir)
	if window_target == nil {
		if focus_change_dir < 0 {
			window_target = c.FindWindowNavigable(len(c.Windows)-1, i_current, focus_change_dir)
		} else {
			window_target = c.FindWindowNavigable(0, i_current, focus_change_dir)
		}
	}
	c.NavWindowingTarget = window_target
	c.NavWindowingToggleLayer = false
}

// Equivalent of IsKeyDown() for NavInputs[]
func (c *Context) IsNavInputDown(n NavInput) bool {
	return c.IO.NavInputs[n] > 0.0
}

func (n *NavMoveResult) Clear() {
	n.ID = 0
	n.ParentID = 0
	n.Window = nil
	n.DistBox = math.MaxFloat32
	n.DistCenter = math.MaxFloat32
	n.DistAxial = math.MaxFloat32
	n.RectRel = f64.Rectangle{}
}

func (c *Context) NavUpdate() {
	c.IO.WantSetMousePos = false

	if c.IO.ConfigFlags&ConfigFlagsNavEnableGamepad != 0 && c.IO.BackendFlags&BackendFlagsHasGamepad != 0 {
		if c.IO.NavInputs[NavInputActivate] > 0.0 || c.IO.NavInputs[NavInputInput] > 0.0 || c.IO.NavInputs[NavInputCancel] > 0.0 || c.IO.NavInputs[NavInputMenu] > 0.0 {
			c.NavInputSource = InputSourceNavGamepad
		}
	}

	// Update Keyboard->Nav inputs mapping
	for i := int(NavInputInternalStart_); i < len(c.IO.NavInputs); i++ {
		c.IO.NavInputs[i] = 0
	}
	if c.IO.ConfigFlags&ConfigFlagsNavEnableKeyboard != 0 {
		c.navMapKey(KeySpace, NavInputActivate)
		c.navMapKey(KeyEnter, NavInputInput)
		c.navMapKey(KeyEscape, NavInputCancel)
		c.navMapKey(KeyLeftArrow, NavInputKeyLeft_)
		c.navMapKey(KeyRightArrow, NavInputKeyRight_)
		c.navMapKey(KeyUpArrow, NavInputKeyUp_)
		c.navMapKey(KeyDownArrow, NavInputKeyDown_)
		if c.IO.KeyCtrl {
			c.IO.NavInputs[NavInputTweakSlow] = 1
		}
		if c.IO.KeyShift {
			c.IO.NavInputs[NavInputTweakFast] = 1
		}
		if c.IO.KeyAlt {
			c.IO.NavInputs[NavInputKeyMenu_] = 1
		}
	}

	copy(c.IO.NavInputsDownDurationPrev[:], c.IO.NavInputsDownDuration[:])
	for i := range c.IO.NavInputs {
		if c.IO.NavInputs[i] > 0.0 {
			if c.IO.NavInputsDownDuration[i] < 0.0 {
				c.IO.NavInputsDownDuration[i] = 0
			} else {
				c.IO.NavInputsDownDuration[i] = c.IO.NavInputsDownDuration[i] + c.IO.DeltaTime
			}
		} else {
			c.IO.NavInputsDownDuration[i] = -1
		}
	}

	// Process navigation init request (select first/default focus)
	if c.NavInitResultId != 0 && (!c.NavDisableHighlight || c.NavInitRequestFromMove) {
		// Apply result from previous navigation init request (will typically select the first item, unless SetItemDefaultFocus() has been called)
		assert(c.NavWindow != nil)
		if c.NavInitRequestFromMove {
			c.SetNavIDWithRectRel(c.NavInitResultId, c.NavLayer, c.NavInitResultRectRel)
		} else {
			c.SetNavID(c.NavInitResultId, c.NavLayer)
		}
		c.NavWindow.NavRectRel[c.NavLayer] = c.NavInitResultRectRel
	}
	c.NavInitRequest = false
	c.NavInitRequestFromMove = false
	c.NavInitResultId = 0
	c.NavJustMovedToId = 0

	// Process navigation move request
	if c.NavMoveRequest && (c.NavMoveResultLocal.ID != 0 || c.NavMoveResultOther.ID != 0) {
		// Select which result to use
		result := &c.NavMoveResultOther
		if c.NavMoveResultLocal.ID != 0 {
			result = &c.NavMoveResultLocal
		}
		// Maybe entering a flattened child? In this case solve the tie using the regular scoring rules
		if c.NavMoveResultOther.ID != 0 && c.NavMoveResultOther.Window.ParentWindow == c.NavWindow {
			if (c.NavMoveResultOther.DistBox < c.NavMoveResultLocal.DistBox) || (c.NavMoveResultOther.DistBox == c.NavMoveResultLocal.DistBox && c.NavMoveResultOther.DistCenter < c.NavMoveResultLocal.DistCenter) {
				result = &c.NavMoveResultOther
			}
		}

		assert(c.NavWindow != nil && result.Window != nil)

		// Scroll to keep newly navigated item fully into view
		if c.NavLayer == 0 {
			c.NavScrollToBringItemIntoView(result.Window, result.RectRel)
		}

		// Apply result from previous frame navigation directional move request
		c.ClearActiveID()
		c.NavWindow = result.Window
		c.SetNavIDWithRectRel(result.ID, c.NavLayer, result.RectRel)
		c.NavJustMovedToId = result.ID
		c.NavMoveFromClampedRefRect = false
	}

	// When a forwarded move request failed, we restore the highlight that we disabled during the forward frame
	if c.NavMoveRequestForward == NavForwardForwardActive {
		assert(c.NavMoveRequest)
		if c.NavMoveResultLocal.ID == 0 && c.NavMoveResultOther.ID == 0 {
			c.NavDisableHighlight = false
		}
		c.NavMoveRequestForward = NavForwardNone
	}

	// Apply application mouse position movement, after we had a chance to process move request result.
	if c.NavMousePosDirty && c.NavIdIsAlive {
		// Set mouse position given our knowledge of the navigated item position from last frame
		if c.IO.ConfigFlags&ConfigFlagsNavEnableSetMousePos != 0 && c.IO.BackendFlags&BackendFlagsHasSetMousePos != 0 {
			assert(!c.NavDisableHighlight && c.NavDisableMouseHover)
			c.IO.MousePos = c.NavCalcPreferredRefPos()
			c.IO.MousePosPrev = c.IO.MousePos
			c.IO.WantSetMousePos = true
		}
		c.NavMousePosDirty = false
	}

	c.NavIdIsAlive = false
	c.NavJustTabbedId = 0
	assert(c.NavLayer == 0 || c.NavLayer == 1)

	// Store our return window (for returning from Layer 1 to Layer 0) and clear it as soon as we step back in our own Layer 0
	if c.NavWindow != nil {
		c.NavSaveLastChildNavWindow(c.NavWindow)
	}
	if c.NavWindow != nil && c.NavWindow.NavLastChildNavWindow != nil && c.NavLayer == 0 {
		c.NavWindow.NavLastChildNavWindow = nil
	}

	c.NavUpdateWindowing()

	// Set output flags for user application
	nav_keyboard_active := (c.IO.ConfigFlags & ConfigFlagsNavEnableKeyboard) != 0
	nav_gamepad_active := (c.IO.ConfigFlags&ConfigFlagsNavEnableGamepad) != 0 && (c.IO.BackendFlags&BackendFlagsHasGamepad) != 0
	c.IO.NavActive = (nav_keyboard_active || nav_gamepad_active) && c.NavWindow != nil && c.NavWindow.Flags&WindowFlagsNoNavInputs == 0
	c.IO.NavVisible = (c.IO.NavActive && c.NavId != 0 && !c.NavDisableHighlight) || (c.NavWindowingTarget != nil) || c.NavInitRequest

	// Process NavCancel input (to close a popup, get back to parent, clear focus)
	if c.IsNavInputPressed(NavInputCancel, InputReadModePressed) {
		if c.ActiveId != 0 {
			c.ClearActiveID()
		} else if c.NavWindow != nil && c.NavWindow.Flags&WindowFlagsChildWindow != 0 && c.NavWindow.Flags&WindowFlagsPopup == 0 && c.NavWindow.ParentWindow != nil {
			// Exit child window
			child_window := c.NavWindow
			parent_window := c.NavWindow.ParentWindow
			assert(child_window.ChildId != 0)
			c.FocusWindow(parent_window)
			c.SetNavID(child_window.ChildId, 0)
			c.NavIdIsAlive = false
			if c.NavDisableMouseHover {
				c.NavMousePosDirty = true
			}
		} else if len(c.OpenPopupStack) > 0 {
			// Close open popup/menu
			if c.OpenPopupStack[len(c.OpenPopupStack)-1].Window.Flags&WindowFlagsModal == 0 {
				c.ClosePopupToLevel(len(c.OpenPopupStack) - 1)
			}
		} else if c.NavLayer != 0 {
			// Leave the "menu" layer
			c.NavRestoreLayer(0)
		} else {
			// Clear NavLastId for popups but keep it for regular child window so we can leave one and come back where we were
			if c.NavWindow != nil && ((c.NavWindow.Flags&WindowFlagsPopup) != 0 || c.NavWindow.Flags&WindowFlagsChildWindow == 0) {
				c.NavWindow.NavLastIds[0] = 0
			}
			c.NavId = 0
		}
	}

	// Process manual activation request
	c.NavActivateId = 0
	c.NavActivateDownId = 0
	c.NavActivatePressedId = 0
	c.NavInputId = 0
	if c.NavId != 0 && !c.NavDisableHighlight && c.NavWindowingTarget == nil && c.NavWindow != nil && c.NavWindow.Flags&WindowFlagsNoNavInputs == 0 {
		activate_down := c.IsNavInputDown(NavInputActivate)
		activate_pressed := activate_down && c.IsNavInputPressed(NavInputActivate, InputReadModePressed)
		if c.ActiveId == 0 && activate_pressed {
			c.NavActivateId = c.NavId
		}
		if (c.ActiveId == 0 || c.ActiveId == c.NavId) && activate_down {
			c.NavActivateDownId = c.NavId
		}
		if (c.ActiveId == 0 || c.ActiveId == c.NavId) && activate_pressed {
			c.NavActivatePressedId = c.NavId
		}
		if (c.ActiveId == 0 || c.ActiveId == c.NavId) && c.IsNavInputPressed(NavInputInput, InputReadModePressed) {
			c.NavInputId = c.NavId
		}
	}

	if c.NavWindow != nil && c.NavWindow.Flags&WindowFlagsNoNavInputs != 0 {
		c.NavDisableHighlight = true
	}
	if c.NavActivateId != 0 {
		assert(c.NavActivateDownId == c.NavActivateId)
	}
	c.NavMoveRequest = false

	// Process programmatic activation request
	if c.NavNextActivateId != 0 {
		c.NavActivateId = c.NavNextActivateId
		c.NavActivateDownId = c.NavNextActivateId
		c.NavActivatePressedId = c.NavNextActivateId
		c.NavInputId = c.NavNextActivateId
	}
	c.NavNextActivateId = 0

	// Initiate directional inputs request
	allowed_dir_flags := c.ActiveIdAllowNavDirFlags
	if c.ActiveId == 0 {
		allowed_dir_flags = ^0
	}
	if c.NavMoveRequestForward == NavForwardNone {
		c.NavMoveDir = DirNone
		if c.NavWindow != nil && c.NavWindowingTarget == nil && allowed_dir_flags != 0 && c.NavWindow.Flags&WindowFlagsNoNavInputs == 0 {
			if (allowed_dir_flags&(1<<uint(DirLeft))) != 0 && c.IsNavInputPressedAnyOfTwo(NavInputDpadLeft, NavInputKeyLeft_, InputReadModeRepeat) {
				c.NavMoveDir = DirLeft
			}
			if (allowed_dir_flags&(1<<uint(DirRight))) != 0 && c.IsNavInputPressedAnyOfTwo(NavInputDpadRight, NavInputKeyRight_, InputReadModeRepeat) {
				c.NavMoveDir = DirRight
			}
			if (allowed_dir_flags&(1<<uint(DirUp))) != 0 && c.IsNavInputPressedAnyOfTwo(NavInputDpadUp, NavInputKeyUp_, InputReadModeRepeat) {
				c.NavMoveDir = DirUp
			}
			if (allowed_dir_flags&(1<<uint(DirDown))) != 0 && c.IsNavInputPressedAnyOfTwo(NavInputDpadDown, NavInputKeyDown_, InputReadModeRepeat) {
				c.NavMoveDir = DirDown
			}
		}
	} else {
		// Forwarding previous request (which has been modified, e.g. wrap around menus rewrite the requests with a starting rectangle at the other side of the window)
		assert(c.NavMoveDir != DirNone)
		assert(c.NavMoveRequestForward == NavForwardForwardQueued)
		c.NavMoveRequestForward = NavForwardForwardActive
	}

	if c.NavMoveDir != DirNone {
		c.NavMoveRequest = true
		c.NavMoveDirLast = c.NavMoveDir
	}

	// If we initiate a movement request and have no current NavId, we initiate a InitDefautRequest that will be used as a fallback if the direction fails to find a match
	if c.NavMoveRequest && c.NavId == 0 {
		c.NavInitRequest = true
		c.NavInitRequestFromMove = true
		c.NavInitResultId = 0
		c.NavDisableHighlight = false
	}

	c.NavUpdateAnyRequestFlag()

	// Scrolling
	if c.NavWindow != nil && c.NavWindow.Flags&WindowFlagsNoNavInputs == 0 && c.NavWindowingTarget == nil {
		// *Fallback* manual-scroll with NavUp/NavDown when window has no navigable item
		window := c.NavWindow
		// We need round the scrolling speed because sub-pixel scroll isn't reliably supported.
		scroll_speed := math.Floor(window.CalcFontSize()*100*c.IO.DeltaTime + 0.5)
		if window.DC.NavLayerActiveMask == 0x00 && window.DC.NavHasScroll && c.NavMoveRequest {
			if c.NavMoveDir == DirLeft || c.NavMoveDir == DirRight {
				dir := 1.0
				if c.NavMoveDir == DirLeft {
					dir = -1.0
				}
				c.SetWindowScrollX(window, math.Floor(window.Scroll.X+dir*scroll_speed))
			}
			if c.NavMoveDir == DirUp || c.NavMoveDir == DirDown {
				dir := 1.0
				if c.NavMoveDir == DirUp {
					dir = -1.0
				}
				c.SetWindowScrollY(window, math.Floor(window.Scroll.Y+dir*scroll_speed))
			}
		}

		// *Normal* Manual scroll with NavScrollXXX keys
		// Next movement request will clamp the NavId reference rectangle to the visible area, so navigation will resume within those bounds.
		scroll_dir := c.GetNavInputAmount2dEx(NavDirSourceFlagsPadLStick, InputReadModeDown, 1.0/10.0, 10.0)
		if scroll_dir.X != 0.0 && window.ScrollbarX {
			c.SetWindowScrollX(window, math.Floor(window.Scroll.X+scroll_dir.X*scroll_speed))
			c.NavMoveFromClampedRefRect = true
		}
		if scroll_dir.Y != 0.0 {
			c.SetWindowScrollY(window, math.Floor(window.Scroll.Y+scroll_dir.Y*scroll_speed))
			c.NavMoveFromClampedRefRect = true
		}
	}

	// Reset search results
	c.NavMoveResultLocal.Clear()
	c.NavMoveResultOther.Clear()

	// When we have manually scrolled (without using navigation) and NavId becomes out of bounds, we project its bounding box to the visible area to restart navigation within visible items
	if c.NavMoveRequest && c.NavMoveFromClampedRefRect && c.NavLayer == 0 {
		window := c.NavWindow
		window_rect_rel := f64.Rectangle{
			window.InnerRect.Min.Sub(window.Pos.Sub(f64.Vec2{1, 1})),
			window.InnerRect.Max.Sub(window.Pos.Add(f64.Vec2{1, 1})),
		}

		if !window.NavRectRel[c.NavLayer].In(window_rect_rel) {
			pad := window.CalcFontSize() * 0.5
			// Terrible approximation for the intent of starting navigation from first fully visible item
			window_rect_rel = window_rect_rel.Expand2(f64.Vec2{
				-math.Min(window_rect_rel.Dx(), pad),
				-math.Min(window_rect_rel.Dy(), pad),
			})
			window.NavRectRel[c.NavLayer] = window.NavRectRel[c.NavLayer].Intersect(window_rect_rel)
			c.NavId = 0
		}
		c.NavMoveFromClampedRefRect = false
	}

	// For scoring we use a single segment on the left side our current item bounding box (not touching the edge to avoid box overlap with zero-spaced items)
	var nav_rect_rel f64.Rectangle
	if c.NavWindow != nil && !c.NavWindow.NavRectRel[c.NavLayer].Inverted() {
		nav_rect_rel = c.NavWindow.NavRectRel[c.NavLayer]
	}
	c.NavScoringRectScreen = c.GetViewportRect()
	if c.NavWindow != nil {
		c.NavScoringRectScreen = f64.Rectangle{
			c.NavWindow.Pos.Add(nav_rect_rel.Min),
			c.NavWindow.Pos.Add(nav_rect_rel.Max),
		}
	}
	c.NavScoringRectScreen.Min.X = math.Min(c.NavScoringRectScreen.Min.X+1.0, c.NavScoringRectScreen.Max.X)
	c.NavScoringRectScreen.Max.X = c.NavScoringRectScreen.Min.X
	// Ensure if we have a finite, non-inverted bounding box here will allows us to remove extraneous fabsf() calls in NavScoreItem().
	assert(!c.NavScoringRectScreen.Inverted())
}

// NB: We modify rect_rel by the amount we scrolled for, so it is immediately updated.
func (c *Context) NavScrollToBringItemIntoView(window *Window, item_rect_rel f64.Rectangle) {
	// Scroll to keep newly navigated item fully into view
	window_pos := window.InnerRect.Min.Sub(window.Pos)
	window_rect_rel := f64.Rectangle{
		window_pos.Sub(f64.Vec2{1, 1}),
		window_pos.Add(f64.Vec2{1, 1}),
	}

	if item_rect_rel.In(window_rect_rel) {
		return
	}

	if window.ScrollbarX && item_rect_rel.Min.X < window_rect_rel.Min.X {
		window.ScrollTarget.X = item_rect_rel.Min.X + window.Scroll.X - c.Style.ItemSpacing.X
		window.ScrollTargetCenterRatio.X = 0.0
	} else if window.ScrollbarX && item_rect_rel.Max.X >= window_rect_rel.Max.X {
		window.ScrollTarget.X = item_rect_rel.Max.X + window.Scroll.X + c.Style.ItemSpacing.X
		window.ScrollTargetCenterRatio.X = 1.0
	}

	if item_rect_rel.Min.Y < window_rect_rel.Min.Y {
		window.ScrollTarget.Y = item_rect_rel.Min.Y + window.Scroll.Y - c.Style.ItemSpacing.Y
		window.ScrollTargetCenterRatio.Y = 0.0
	} else if item_rect_rel.Max.Y >= window_rect_rel.Max.Y {
		window.ScrollTarget.Y = item_rect_rel.Max.Y + window.Scroll.Y + c.Style.ItemSpacing.Y
		window.ScrollTargetCenterRatio.Y = 1.0
	}

	// Estimate upcoming scroll so we can offset our relative mouse position so mouse position can be applied immediately (under this block)
	next_scroll := c.CalcNextScrollFromScrollTargetAndClamp(window)
	item_rect_rel = item_rect_rel.Add(window.Scroll.Sub(next_scroll))
}

func (c *Context) navMapKey(key Key, nav_input NavInput) {
	if c.IO.KeyMap[key] != -1 && c.IsKeyDown(c.IO.KeyMap[key]) {
		c.IO.NavInputs[nav_input] = 1
	}
}

// We get there when either NavId == id, or when g.NavAnyRequest is set (which is updated by NavUpdateAnyRequestFlag above)
func (c *Context) NavProcessItem(window *Window, nav_bb f64.Rectangle, id ID) {
	item_flags := window.DC.ItemFlags
	nav_bb_rel := f64.Rectangle{
		nav_bb.Min.Sub(window.Pos),
		nav_bb.Max.Sub(window.Pos),
	}
	if c.NavInitRequest && c.NavLayer == window.DC.NavLayerCurrent {
		// Even if 'ImGuiItemFlags_NoNavDefaultFocus' is on (typically collapse/close button) we record the first ResultId so they can be used as a fallback
		if item_flags&ItemFlagsNoNavDefaultFocus == 0 || c.NavInitResultId == 0 {
			c.NavInitResultId = id
			c.NavInitResultRectRel = nav_bb_rel
		}

		if item_flags&ItemFlagsNoNavDefaultFocus == 0 {
			c.NavInitRequest = false // Found a match, clear request
			c.NavUpdateAnyRequestFlag()
		}
	}

	// Scoring for navigation
	if c.NavId != id && item_flags&ItemFlagsNoNav == 0 {
		var result *NavMoveResult
		if window == c.NavWindow {
			result = &c.NavMoveResultLocal
		} else {
			result = &c.NavMoveResultOther
		}

		new_best := c.NavMoveRequest && c.NavScoreItem(result, nav_bb)
		if new_best {
			result.ID = id
			result.ParentID = window.IDStack[len(window.IDStack)-1]
			result.Window = window
			result.RectRel = nav_bb_rel
		}
	}

	// Update window-relative bounding box of navigated item
	if c.NavId == id {
		// Always refresh g.NavWindow, because some operations such as FocusItem() don't have a window.
		c.NavWindow = window
		c.NavLayer = window.DC.NavLayerCurrent
		c.NavIdIsAlive = true
		c.NavIdTabCounter = window.FocusIdxTabCounter
		// Store item bounding box (relative to window position)
		window.NavRectRel[window.DC.NavLayerCurrent] = nav_bb_rel
	}
}

func (c *Context) NavUpdateAnyRequestFlag() {
	c.NavAnyRequest = c.NavMoveRequest || c.NavInitRequest
	if c.NavAnyRequest {
		assert(c.NavWindow != nil)
	}
}

// Scoring function for directional navigation. Based on https://gist.github.com/rygorous/6981057
func (c *Context) NavScoreItem(result *NavMoveResult, cand f64.Rectangle) bool {
	window := c.CurrentWindow
	if c.NavLayer != window.DC.NavLayerCurrent {
		return false
	}

	// Current modified source rect (NB: we've applied Max.x = Min.x in NavUpdate() to inhibit the effect of having varied item width)
	curr := &c.NavScoringRectScreen
	c.NavScoringCount++

	// We perform scoring on items bounding box clipped by their parent window on the other axis (clipping on our movement axis would give us equal scores for all clipped items)
	if c.NavMoveDir == DirLeft || c.NavMoveDir == DirRight {
		cand.Min.Y = f64.Clamp(cand.Min.Y, window.ClipRect.Min.Y, window.ClipRect.Max.Y)
		cand.Max.Y = f64.Clamp(cand.Max.Y, window.ClipRect.Min.Y, window.ClipRect.Max.Y)
	} else {
		cand.Min.X = f64.Clamp(cand.Min.X, window.ClipRect.Min.X, window.ClipRect.Max.X)
		cand.Max.X = f64.Clamp(cand.Max.X, window.ClipRect.Min.X, window.ClipRect.Max.X)
	}

	// Compute distance between boxes
	// FIXME-NAV: Introducing biases for vertical navigation, needs to be removed.
	dbx := c.NavScoreItemDistInterval(cand.Min.X, cand.Max.X, curr.Min.X, curr.Max.X)
	// Scale down on Y to keep using box-distance for vertically touching items
	dby := c.NavScoreItemDistInterval(
		f64.Lerp(0.2, cand.Min.Y, cand.Max.Y),
		f64.Lerp(0.8, cand.Min.Y, cand.Max.Y),
		f64.Lerp(0.2, curr.Min.Y, curr.Max.Y),
		f64.Lerp(0.8, curr.Min.Y, curr.Max.Y),
	)
	if dby != 0 && dbx != 0 {
		if dbx > 0 {
			dbx = dbx/1000 + 1
		} else {
			dbx = dbx/1000 - 1
		}
	}
	dist_box := math.Abs(dbx) + math.Abs(dby)

	// Compute distance between centers (this is off by a factor of 2, but we only compare center distances with each other so it doesn't matter)
	dcx := (cand.Min.X + cand.Max.X) - (curr.Min.X + curr.Max.X)
	dcy := (cand.Min.Y + cand.Max.Y) - (curr.Min.Y + curr.Max.Y)
	dist_center := math.Abs(dcx) + math.Abs(dcy) // L1 metric (need this for our connectedness guarantee)

	// Determine which quadrant of 'curr' our candidate item 'cand' lies in based on distance
	var quadrant Dir
	var dax, day, dist_axial float64
	if dbx != 0 || dby != 0 {
		// For non-overlapping boxes, use distance between boxes
		dax = dbx
		day = dby
		dist_axial = dist_box
		quadrant = c.NavScoreItemGetQuadrant(dbx, dby)
	} else if dcx != 0 || dcy != 0 {
		// For overlapping boxes with different centers, use distance between centers
		dax = dcx
		day = dcy
		dist_axial = dist_center
		quadrant = c.NavScoreItemGetQuadrant(dcx, dcy)
	} else {
		// Degenerate case: two overlapping buttons with same center, break ties arbitrarily (note that LastItemId here is really the _previous_ item order, but it doesn't matter)
		if window.DC.LastItemId < c.NavId {
			quadrant = DirLeft
		} else {
			quadrant = DirRight
		}
	}

	// Is it in the quadrant we're interesting in moving to?
	new_best := false
	if quadrant == c.NavMoveDir {
		// Does it beat the current best candidate?
		if dist_box < result.DistBox {
			result.DistBox = dist_box
			result.DistCenter = dist_center
			return true
		}

		if dist_box == result.DistBox {
			// Try using distance between center points to break ties
			if dist_center < result.DistCenter {
				result.DistCenter = dist_center
				new_best = true
			} else if dist_center == result.DistCenter {
				// Still tied! we need to be extra-careful to make sure everything gets linked properly. We consistently break ties by symbolically moving "later" items
				// (with higher index) to the right/downwards by an infinitesimal amount since we the current "best" button already (so it must have a lower index),
				// this is fairly easy. This rule ensures that all buttons with dx==dy==0 will end up being linked in order of appearance along the x axis.

				// moving bj to the right/down decreases distance
				if c.NavMoveDir == DirUp || c.NavMoveDir == DirDown {
					new_best = dby < 0
				} else {
					new_best = dbx < 0
				}
			}
		}
	}

	// Axial check: if 'curr' has no link at all in some direction and 'cand' lies roughly in that direction, add a tentative link. This will only be kept if no "real" matches
	// are found, so it only augments the graph produced by the above method using extra links. (important, since it doesn't guarantee strong connectedness)
	// This is just to avoid buttons having no links in a particular direction when there's a suitable neighbor. you get good graphs without this too.
	// 2017/09/29: FIXME: This now currently only enabled inside menu bars, ideally we'd disable it everywhere. Menus in particular need to catch failure. For general navigation it feels awkward.
	// Disabling it may however lead to disconnected graphs when nodes are very spaced out on different axis. Perhaps consider offering this as an option?

	// Check axial match
	if result.DistBox == math.MaxFloat32 && dist_axial < result.DistAxial {
		if c.NavLayer == 1 && c.NavWindow.Flags&WindowFlagsChildMenu == 0 {
			if (c.NavMoveDir == DirRight && dax > 0) || (c.NavMoveDir == DirRight && dax > 0) ||
				(c.NavMoveDir == DirUp && day < 0) || (c.NavMoveDir == DirDown && day > 0) {
				result.DistAxial = dist_axial
				new_best = true
			}
		}
	}

	return new_best
}

func (c *Context) NavScoreItemDistInterval(a0, a1, b0, b1 float64) float64 {
	if a1 < b0 {
		return a1 - b0
	}
	if b1 < a0 {
		return a0 - b1
	}
	return 0
}

func (c *Context) NavScoreItemGetQuadrant(dx, dy float64) Dir {
	if math.Abs(dx) > math.Abs(dy) {
		if dx > 0 {
			return DirRight
		}
		return DirLeft
	}

	if dy > 0 {
		return DirDown
	}
	return DirUp
}

func (c *Context) NavProcessMoveRequestWrapAround(window *Window) {
	if c.NavWindow == window && c.NavMoveRequestButNoResultYet() {
		if (c.NavMoveDir == DirUp || c.NavMoveDir == DirDown) &&
			c.NavMoveRequestForward == NavForwardNone && c.NavLayer == 0 {
			c.NavMoveRequestForward = NavForwardForwardQueued
			c.NavMoveRequestCancel()

			c.NavWindow.NavRectRel[0].Min.Y = 0
			if c.NavMoveDir == DirUp {
				c.NavWindow.NavRectRel[0].Min.Y = math.Max(window.SizeFull.Y, window.SizeContents.Y) - window.Scroll.Y
			} else {
				c.NavWindow.NavRectRel[0].Min.Y = -window.Scroll.Y
			}

			c.NavWindow.NavRectRel[0].Max.Y = c.NavWindow.NavRectRel[0].Min.Y
		}
	}
}

func (c *Context) NavMoveRequestButNoResultYet() bool {
	return c.NavMoveRequest && c.NavMoveResultLocal.ID == 0 && c.NavMoveResultOther.ID == 0
}

func (c *Context) NavMoveRequestCancel() {
	c.NavMoveRequest = false
	c.NavUpdateAnyRequestFlag()
}

func (c *Context) NavRestoreLayer(layer int) {
	c.NavLayer = layer
	if layer == 0 {
		c.NavWindow = c.NavRestoreLastChildNavWindow(c.NavWindow)
	}
	if layer == 0 && c.NavWindow.NavLastIds[0] != 0 {
		c.SetNavIDAndMoveMouse(c.NavWindow.NavLastIds[0], layer, c.NavWindow.NavRectRel[0])
	} else {
		c.NavInitWindow(c.NavWindow, true)
	}
}

func (c *Context) SetNavIDAndMoveMouse(id ID, nav_layer int, rect_rel f64.Rectangle) {
	c.SetNavID(id, nav_layer)
	c.NavWindow.NavRectRel[nav_layer] = rect_rel
	c.NavMousePosDirty = true
	c.NavDisableHighlight = false
	c.NavDisableMouseHover = true
}

func (c *Context) SetNavID(id ID, nav_layer int) {
	c.NavId = id
	c.NavWindow.NavLastIds[nav_layer] = id
}

func (c *Context) NavInitWindow(window *Window, force_reinit bool) {
	var init_for_nav bool
	if window.Flags&WindowFlagsNoNavInputs == 0 {
		if window.Flags&WindowFlagsChildWindow == 0 || window.Flags&WindowFlagsPopup != 0 ||
			window.NavLastIds[0] == 0 || force_reinit {
			init_for_nav = true
		}
	}

	if init_for_nav {
		c.SetNavID(0, c.NavLayer)
		c.NavInitRequest = true
		c.NavInitRequestFromMove = false
		c.NavInitResultId = 0
		c.NavInitResultRectRel = f64.Rectangle{}
		c.NavUpdateAnyRequestFlag()
	} else {
		c.NavId = window.NavLastIds[0]
	}
}

func (c *Context) NavCalcPreferredRefPos() f64.Vec2 {
	if c.NavDisableHighlight || !c.NavDisableMouseHover || c.NavWindow == nil {
		return c.IO.MousePos.Floor()
	}

	// When navigation is active and mouse is disabled, decide on an arbitrary position around the bottom left of the currently navigated item
	rect_rel := c.NavWindow.NavRectRel[c.NavLayer]
	pos := c.NavWindow.Pos.Add(f64.Vec2{
		rect_rel.Min.X + math.Min(c.Style.FramePadding.X*4, rect_rel.Dx()),
		rect_rel.Max.Y - math.Min(c.Style.FramePadding.Y, rect_rel.Dy()),
	})
	visible_rect := c.GetViewportRect()

	// ImFloor() is important because non-integer mouse position application in back-end might be lossy and result in undesirable non-zero delta.
	pos = pos.Clamp2(visible_rect.Min, visible_rect.Max)
	pos = pos.Floor()
	return pos
}

func (c *Context) SetNavIDWithRectRel(id ID, nav_layer int, rect_rel f64.Rectangle) {
	c.SetNavID(id, nav_layer)
	c.NavWindow.NavRectRel[nav_layer] = rect_rel
	c.NavMousePosDirty = true
	c.NavDisableHighlight = false
	c.NavDisableMouseHover = true
}

func (c *Context) IsNavInputPressedAnyOfTwo(n1, n2 NavInput, mode InputReadMode) bool {
	return c.GetNavInputAmount(n1, mode)+c.GetNavInputAmount(n2, mode) > 0
}

func (c *Context) NavSaveLastChildNavWindow(child_window *Window) {
	parent_window := child_window
	if parent_window != nil && (parent_window.Flags&WindowFlagsChildWindow) != 0 && (parent_window.Flags&(WindowFlagsPopup|WindowFlagsChildMenu)) == 0 {
		parent_window = parent_window.ParentWindow
	}
	if parent_window != nil && parent_window != child_window {
		parent_window.NavLastChildNavWindow = child_window
	}
}

// Window management mode (hold to: change focus/move/resize, tap to: toggle menu layer)
func (c *Context) NavUpdateWindowing() {
	var apply_focus_window *Window
	var apply_toggle_layer bool

	start_windowing_with_gamepad := c.NavWindowingTarget == nil && c.IsNavInputPressed(NavInputMenu, InputReadModePressed)
	start_windowing_with_keyboard := c.NavWindowingTarget == nil && c.IO.KeyCtrl && c.IsKeyPressedMap(KeyTab) && c.IO.ConfigFlags&ConfigFlagsNavEnableKeyboard != 0
	if start_windowing_with_gamepad || start_windowing_with_keyboard {
		window := c.NavWindow
		if window == nil {
			window = c.FindWindowNavigable(len(c.Windows)-1, -math.MaxInt32, -1)
		}
		if window != nil {
			c.NavWindowingTarget = window.RootWindowForTabbing
			c.NavWindowingHighlightTimer = 0.0
			c.NavWindowingHighlightAlpha = 0.0
			c.NavWindowingToggleLayer = true
			c.NavInputSource = InputSourceNavGamepad
			if start_windowing_with_keyboard {
				c.NavWindowingToggleLayer = false
				c.NavInputSource = InputSourceNavKeyboard
			}
		}
	}

	// Gamepad update
	c.NavWindowingHighlightTimer += c.IO.DeltaTime
	if c.NavWindowingTarget != nil && c.NavInputSource == InputSourceNavGamepad {
		// Highlight only appears after a brief time holding the button, so that a fast tap on PadMenu (to toggle NavLayer) doesn't add visual noise
		c.NavWindowingHighlightAlpha = math.Max(c.NavWindowingHighlightAlpha, f64.Saturate((c.NavWindowingHighlightTimer-0.20)/0.05))

		// Select window to focus
		focus_change_dir := truth(c.IsNavInputPressed(NavInputFocusPrev, InputReadModeRepeatSlow)) - truth(c.IsNavInputPressed(NavInputFocusNext, InputReadModeRepeatSlow))
		if focus_change_dir != 0 {
			c.NavUpdateWindowingHighlightWindow(focus_change_dir)
			c.NavWindowingHighlightAlpha = 1.0
		}

		// Single press toggles NavLayer, long press with L/R apply actual focus on release (until then the window was merely rendered front-most)
		if !c.IsNavInputDown(NavInputMenu) {
			// Once button was held long enough we don't consider it a tap-to-toggle-layer press anymore.
			if !(c.NavWindowingToggleLayer && c.NavWindowingHighlightAlpha < 1.0) {
				c.NavWindowingToggleLayer = true
			} else if !c.NavWindowingToggleLayer {
				apply_focus_window = c.NavWindowingTarget
			}
			c.NavWindowingTarget = nil
		}
	}

	// Keyboard: Focus
	if c.NavWindowingTarget != nil && c.NavInputSource == InputSourceNavKeyboard {
		// Visuals only appears after a brief time after pressing TAB the first time, so that a fast CTRL+TAB doesn't add visual noise
		c.NavWindowingHighlightAlpha = math.Max(c.NavWindowingHighlightAlpha, f64.Saturate((c.NavWindowingHighlightTimer-0.15)/0.04))

		if c.IsKeyPressedMapEx(KeyTab, true) {
			if c.IO.KeyShift {
				c.NavUpdateWindowingHighlightWindow(1)
			} else {
				c.NavUpdateWindowingHighlightWindow(-1)
			}
		}

		if !c.IO.KeyCtrl {
			apply_focus_window = c.NavWindowingTarget
		}
	}

	// Keyboard: Press and Release ALT to toggle menu layer
	// FIXME: We lack an explicit IO variable for "is the imgui window focused", so compare mouse validity to detect the common case of back-end clearing releases all keys on ALT-TAB
	if (c.ActiveId == 0 || c.ActiveIdAllowOverlap) && c.IsNavInputPressed(NavInputKeyMenu_, InputReadModeReleased) {
		if c.IsMousePosValidEx(&c.IO.MousePos) == c.IsMousePosValidEx(&c.IO.MousePosPrev) {
			apply_toggle_layer = true
		}
	}

	// Move window
	if c.NavWindowingTarget != nil && c.NavWindowingTarget.Flags&WindowFlagsNoMove == 0 {
		var move_delta f64.Vec2
		if c.NavInputSource == InputSourceNavKeyboard && !c.IO.KeyShift {
			move_delta = c.GetNavInputAmount2d(NavDirSourceFlagsKeyboard, InputReadModeDown)
		}
		if c.NavInputSource == InputSourceNavGamepad {
			move_delta = c.GetNavInputAmount2d(NavDirSourceFlagsPadLStick, InputReadModeDown)
		}
		if move_delta.X != 0.0 || move_delta.Y != 0.0 {
			const NAV_MOVE_SPEED = 800.0
			move_speed := math.Floor(NAV_MOVE_SPEED * c.IO.DeltaTime * math.Min(c.IO.DisplayFramebufferScale.X, c.IO.DisplayFramebufferScale.Y))
			c.NavWindowingTarget.Pos = c.NavWindowingTarget.Pos.Add(move_delta.Scale(move_speed))
			c.NavDisableMouseHover = true
			c.MarkIniSettingsDirtyForWindow(c.NavWindowingTarget)
		}
	}

	// Apply final focus
	if apply_focus_window != nil && (c.NavWindow == nil || apply_focus_window != c.NavWindow.RootWindowForTabbing) {
		c.NavDisableHighlight = false
		c.NavDisableMouseHover = true
		apply_focus_window = c.NavRestoreLastChildNavWindow(apply_focus_window)
		c.ClosePopupsOverWindow(apply_focus_window)
		c.FocusWindow(apply_focus_window)
		if apply_focus_window.NavLastIds[0] == 0 {
			c.NavInitWindow(apply_focus_window, false)
		}

		// If the window only has a menu layer, select it directly
		if apply_focus_window.DC.NavLayerActiveMask == (1 << 1) {
			c.NavLayer = 1
		}
	}
	if apply_focus_window != nil {
		c.NavWindowingTarget = nil
	}

	// Apply menu/layer toggle
	if apply_toggle_layer && c.NavWindow != nil {
		new_nav_window := c.NavWindow
		for (new_nav_window.DC.NavLayerActiveMask&(1<<1)) == 0 && (new_nav_window.Flags&WindowFlagsChildWindow) != 0 && (new_nav_window.Flags&(WindowFlagsPopup|WindowFlagsChildMenu)) == 0 {
			new_nav_window = new_nav_window.ParentWindow
		}
		if new_nav_window != new_nav_window.ParentWindow {
			old_nav_window := c.NavWindow
			c.FocusWindow(new_nav_window)
			new_nav_window.NavLastChildNavWindow = old_nav_window
		}
		c.NavDisableHighlight = false
		c.NavDisableMouseHover = true
		if c.NavWindow.DC.NavLayerActiveMask&(1<<1) != 0 {
			c.NavRestoreLayer(c.NavLayer ^ 1)
		} else {
			c.NavRestoreLayer(0)
		}
	}
}
