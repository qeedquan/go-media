package imgui

import (
	"encoding/binary"
	"hash/fnv"
	"math"
	"sort"

	"github.com/qeedquan/go-media/math/f64"
	"github.com/qeedquan/go-media/math/mathutil"
)

// Resizing callback data to apply custom constraint. As enabled by SetNextWindowSizeConstraints(). Callback is called during the next Begin().
// NB: For basic min/max size constraint on each axis you don't need to use the callback! The SetNextWindowSizeConstraints() parameters are enough.
type SizeCallbackData struct {
	Pos         f64.Vec2 // Read-only.   Window position, for reference.
	CurrentSize f64.Vec2 // Read-only.   Current window size.
	DesiredSize f64.Vec2 // Read-write.  Desired size, based on user's mouse position. Write to this field to restrain resizing.
}

type WindowFlags int

const (
	WindowFlagsNoTitleBar        WindowFlags = 1 << 0 // Disable title-bar
	WindowFlagsNoResize          WindowFlags = 1 << 1 // Disable user resizing with the lower-right grip
	WindowFlagsNoMove            WindowFlags = 1 << 2 // Disable user moving the window
	WindowFlagsNoScrollbar       WindowFlags = 1 << 3 // Disable scrollbars (window can still scroll with mouse or programatically)
	WindowFlagsNoScrollWithMouse WindowFlags = 1 << 4 // Disable user vertically scrolling with mouse wheel. On child window mouse wheel will be forwarded to the parent unless NoScrollbar is also set.
	WindowFlagsNoCollapse        WindowFlags = 1 << 5 // Disable user collapsing window by double-clicking on it
	WindowFlagsAlwaysAutoResize  WindowFlags = 1 << 6 // Resize every window to its content every frame
	//WindowFlagsShowBorders   WindowFlags       = 1 << 7   // Show borders around windows and items (OBSOLETE! Use e.g. style.FrameBorderSize=1.0f to enable borders).
	WindowFlagsNoSavedSettings           WindowFlags = 1 << 8  // Never load/save settings in .ini file
	WindowFlagsNoInputs                  WindowFlags = 1 << 9  // Disable catching mouse or keyboard inputs hovering test with pass through.
	WindowFlagsMenuBar                   WindowFlags = 1 << 10 // Has a menu-bar
	WindowFlagsHorizontalScrollbar       WindowFlags = 1 << 11 // Allow horizontal scrollbar to appear (off by default). You may use SetNextWindowContentSize(ImVec2(width0.0f)); prior to calling Begin() to specify width. Read code in imgui_demo in the "Horizontal Scrolling" section.
	WindowFlagsNoFocusOnAppearing        WindowFlags = 1 << 12 // Disable taking focus when transitioning from hidden to visible state
	WindowFlagsNoBringToFrontOnFocus     WindowFlags = 1 << 13 // Disable bringing window to front when taking focus (e.g. clicking on it or programatically giving it focus)
	WindowFlagsAlwaysVerticalScrollbar   WindowFlags = 1 << 14 // Always show vertical scrollbar (even if ContentSize.y < Size.y)
	WindowFlagsAlwaysHorizontalScrollbar WindowFlags = 1 << 15 // Always show horizontal scrollbar (even if ContentSize.x < Size.x)
	WindowFlagsAlwaysUseWindowPadding    WindowFlags = 1 << 16 // Ensure child windows without border uses style.WindowPadding (ignored by default for non-bordered child windows because more convenient)
	WindowFlagsResizeFromAnySide         WindowFlags = 1 << 17 // (WIP) Enable resize from any corners and borders. Your back-end needs to honor the different values of io.MouseCursor set by imgui.
	WindowFlagsNoNavInputs               WindowFlags = 1 << 18 // No gamepad/keyboard navigation within the window
	WindowFlagsNoNavFocus                WindowFlags = 1 << 19 // No focusing toward this window with gamepad/keyboard navigation (e.g. skipped by CTRL+TAB)
	WindowFlagsNoNav                     WindowFlags = WindowFlagsNoNavInputs | WindowFlagsNoNavFocus

	// [Internal]
	WindowFlagsNavFlattened WindowFlags = 1 << 23 // (WIP) Allow gamepad/keyboard navigation to cross over parent border to this child (only use on child that have no scrolling!)
	WindowFlagsChildWindow  WindowFlags = 1 << 24 // Don't use! For internal use by BeginChild()
	WindowFlagsTooltip      WindowFlags = 1 << 25 // Don't use! For internal use by BeginTooltip()
	WindowFlagsPopup        WindowFlags = 1 << 26 // Don't use! For internal use by BeginPopup()
	WindowFlagsModal        WindowFlags = 1 << 27 // Don't use! For internal use by BeginPopupModal()
	WindowFlagsChildMenu    WindowFlags = 1 << 28 // Don't use! For internal use by BeginMenu()
)

type WindowSettings struct {
	Name      string
	Id        ID
	Pos       f64.Vec2
	Size      f64.Vec2
	Collapsed bool
}

type Window struct {
	Ctx                            *Context
	Name                           string
	ID                             ID            // == ImHash(Name)
	Flags                          WindowFlags   // See enum ImGuiWindowFlags_
	Pos                            f64.Vec2      // Position rounded-up to nearest pixel
	Size                           f64.Vec2      // Current size (==SizeFull or collapsed title bar size)
	SizeFull                       f64.Vec2      // Size when non collapsed
	SizeFullAtLastBegin            f64.Vec2      // Copy of SizeFull at the end of Begin. This is the reference value we'll use on the next frame to decide if we need scrollbars.
	SizeContents                   f64.Vec2      // Size of contents (== extents reach of the drawing cursor) from previous frame. Include decoration, window title, border, menu, etc.
	SizeContentsExplicit           f64.Vec2      // Size of contents explicitly set by the user via SetNextWindowContentSize()
	ContentsRegionRect             f64.Rectangle // Maximum visible content position in window coordinates. ~~ (SizeContentsExplicit ? SizeContentsExplicit : Size - ScrollbarSizes) - CursorStartPos, per axis
	WindowPadding                  f64.Vec2      // Window padding at the time of begin.
	WindowRounding                 float64       // Window rounding at the time of begin.
	WindowBorderSize               float64       // Window border size at the time of begin.
	MoveId                         ID            // == window->GetID("#MOVE")
	ChildId                        ID            // Id of corresponding item in parent window (for child windows)
	Scroll                         f64.Vec2
	ScrollTarget                   f64.Vec2 // target scroll position. stored as cursor position with scrolling canceled out, so the highest point is always 0.0f. (FLT_MAX for no change)
	ScrollTargetCenterRatio        f64.Vec2 // 0.0f = scroll so that target position is at top, 0.5f = scroll so that target position is centered
	ScrollbarSizes                 f64.Vec2
	ScrollbarX, ScrollbarY         bool
	Active                         bool // Set to true on Begin(), unless Collapsed
	WasActive                      bool
	WriteAccessed                  bool // Set to true when any widget access the current window
	Collapsed                      bool // Set when collapsing window to become only title-bar
	CollapseToggleWanted           bool
	SkipItems                      bool // Set when items can safely be all clipped (e.g. window not visible or collapsed)
	Appearing                      bool // Set during the frame where the window is appearing (or re-appearing)
	CloseButton                    bool // Set when the window has a close button (p_open != NULL)
	BeginOrderWithinParent         int  // Order within immediate parent window, if we are a child window. Otherwise 0.
	BeginOrderWithinContext        int  // Order within entire imgui context. This is mostly used for debugging submission order related issues.
	BeginCount                     int  // Number of Begin() during the current frame (generally 0 or 1, 1+ if appending via multiple Begin/End pairs)
	PopupId                        ID   // ID in the popup stack when this window is used as a popup/menu (because we use generic Name/ID for recycling)
	AutoFitFramesX, AutoFitFramesY int
	AutoFitOnlyGrows               bool
	AutoFitChildAxises             int
	AutoPosLastDirection           Dir
	HiddenFrames                   int
	SetWindowPosAllowFlags         Cond          // store acceptable condition flags for SetNextWindowPos() use.
	SetWindowSizeAllowFlags        Cond          // store acceptable condition flags for SetNextWindowSize() use.
	SetWindowCollapsedAllowFlags   Cond          // store acceptable condition flags for SetNextWindowCollapsed() use.
	SetWindowPosVal                f64.Vec2      // store window position when using a non-zero Pivot (position set needs to be processed when we know the window size)
	SetWindowPosPivot              f64.Vec2      // store window pivot for positioning. ImVec2(0,0) when positioning from top-left corner; ImVec2(0.5f,0.5f) for centering; ImVec2(1,1) for bottom right.
	DC                             DrawContext   // Temporary per-window data, reset at the beginning of the frame
	IDStack                        []ID          // ID stack. ID are hashes seeded with the value at the top of the stack
	ClipRect                       f64.Rectangle // = DrawList->clip_rect_stack.back(). Scissoring / clipping rectangle. x1, y1, x2, y2.
	WindowRectClipped              f64.Rectangle // = WindowRect just after setup in Begin(). == window->Rect() for root window.
	InnerRect, InnerClipRect       f64.Rectangle
	LastFrameActive                int
	ItemWidthDefault               float64
	MenuColumns                    MenuColumns // Simplified columns storage for menu items
	StateStorage                   map[ID]interface{}
	ColumnsStorage                 []ColumnsSet
	FontWindowScale                float64   // User scale multiplier per-window
	DrawList                       *DrawList // == &DrawListInst (for backward compatibility reason with code using imgui_internal.h we keep this a pointer)
	DrawListInst                   DrawList
	ParentWindow                   *Window // If we are a child _or_ popup window, this is pointing to our parent. Otherwise NULL.
	RootWindow                     *Window // Point to ourself or first ancestor that is not a child window.
	RootWindowForTitleBarHighlight *Window // Point to ourself or first ancestor which will display TitleBgActive color when this window is active.
	RootWindowForTabbing           *Window // Point to ourself or first ancestor which can be CTRL-Tabbed into.
	RootWindowForNav               *Window // Point to ourself or first ancestor which doesn't have the NavFlattened flag.

	NavLastChildNavWindow *Window          // When going to the menu bar, we remember the child window we came from. (This could probably be made implicit if we kept g.Windows sorted by last focused including child window.)
	NavLastIds            [2]ID            // Last known NavId for this window, per layer (0/1)
	NavRectRel            [2]f64.Rectangle // Reference rectangle, in window relative space

	// Navigation / Focus
	// FIXME-NAV: Merge all this with the new Nav system, at least the request variables should be moved to ImGuiContext
	FocusIdxAllCounter        int // Start at -1 and increase as assigned via FocusItemRegister()
	FocusIdxTabCounter        int // (same, but only count widgets which you can Tab through)
	FocusIdxAllRequestCurrent int // Item being requested for focus
	FocusIdxTabRequestCurrent int // Tab-able item being requested for focus
	FocusIdxAllRequestNext    int // Item being requested for focus, for next update (relies on layout to be stable between the frame pressing TAB and the next frame)
	FocusIdxTabRequestNext    int // "
}

type DrawContext struct {
	CursorPos                     f64.Vec2
	CursorPosPrevLine             f64.Vec2
	CursorStartPos                f64.Vec2
	CursorMaxPos                  f64.Vec2 // Used to implicitly calculate the size of our contents, always growing during the frame. Turned into window->SizeContents at the beginning of next frame
	CurrentLineHeight             float64
	CurrentLineTextBaseOffset     float64
	PrevLineHeight                float64
	PrevLineTextBaseOffset        float64
	LogLinePosY                   float64
	TreeDepth                     int
	TreeDepthMayJumpToParentOnPop uint32 // Store a copy of !g.NavIdIsAlive for TreeDepth 0..31
	LastItemId                    ID
	LastItemStatusFlags           ItemStatusFlags
	LastItemRect                  f64.Rectangle // Interaction rect
	LastItemDisplayRect           f64.Rectangle // End-user display rect (only valid if LastItemStatusFlags & ImGuiItemStatusFlags_HasDisplayRect)
	NavHideHighlightOneFrame      bool
	NavHasScroll                  bool     // Set when scrolling can be used (ScrollMax > 0.0f)
	NavLayerCurrent               int      // Current layer, 0..31 (we currently only use 0..1)
	NavLayerCurrentMask           int      // = (1 << NavLayerCurrent) used by ItemAdd prior to clipping.
	NavLayerActiveMask            int      // Which layer have been written to (result from previous frame)
	NavLayerActiveMaskNext        int      // Which layer have been written to (buffer for current frame)
	MenuBarAppending              bool     // FIXME: Remove this
	MenuBarOffset                 f64.Vec2 // MenuBarOffset.x is sort of equivalent of a per-layer CursorPos.x, saved/restored as we switch to the menu bar. The only situation when MenuBarOffset.y is > 0 if when (SafeAreaPadding.y > FramePadding.y), often used on TVs.
	ChildWindows                  []*Window
	StateStorage                  map[ID]interface{}
	LayoutType                    LayoutType
	ParentLayoutType              LayoutType // Layout type of parent window at the time of Begin()

	// We store the current settings outside of the vectors to increase memory locality (reduce cache misses). The vectors are rarely modified. Also it allows us to not heap allocate for short-lived windows which are not using those settings.
	ItemFlags        ItemFlags // == ItemFlagsStack.back() [empty == ImGuiItemFlags_Default]
	ItemWidth        float64   // == ItemWidthStack.back(). 0.0: default, >0.0: width in pixels, <0.0: align xx pixels to the right of window
	TextWrapPos      float64   // == TextWrapPosStack.back() [empty == -1.0f]
	ItemFlagsStack   []ItemFlags
	ItemWidthStack   []float64
	TextWrapPosStack []float64
	GroupStack       []GroupData
	StackSizesBackup [6]int // Store size of various stacks for asserting

	IndentX        float64 // Indentation / start position from left of window (increased by TreePush/TreePop, etc.)
	GroupOffsetX   float64
	ColumnsOffsetX float64     // Offset to the current column (if ColumnsCurrent > 0). FIXME: This and the above should be a stack to allow use cases like Tree->Column->Tree. Need revamp columns API.
	ColumnsSet     *ColumnsSet // Current columns set
}

type GroupData struct {
	BackupCursorPos                 f64.Vec2
	BackupCursorMaxPos              f64.Vec2
	BackupIndentX                   float64
	BackupGroupOffsetX              float64
	BackupCurrentLineHeight         float64
	BackupCurrentLineTextBaseOffset float64
	BackupLogLinePosY               float64
	BackupActiveIdIsAlive           bool
	AdvanceCursor                   bool
}

type ItemFlags int

const (
	ItemFlagsAllowKeyboardFocus       ItemFlags = 1 << 0 // true
	ItemFlagsButtonRepeat             ItemFlags = 1 << 1 // false    // Button() will return true multiple times based on io.KeyRepeatDelay and io.KeyRepeatRate settings.
	ItemFlagsDisabled                 ItemFlags = 1 << 2 // false    // FIXME-WIP: Disable interactions but doesn't affect visuals. Should be: grey out and disable interactions with widgets that affect data + view widgets (WIP)
	ItemFlagsNoNav                    ItemFlags = 1 << 3 // false
	ItemFlagsNoNavDefaultFocus        ItemFlags = 1 << 4 // false
	ItemFlagsSelectableDontClosePopup ItemFlags = 1 << 5 // false    // MenuItem/Selectable() automatically closes current Popup window
	ItemFlagsDefault_                 ItemFlags = ItemFlagsAllowKeyboardFocus
)

type LayoutType int

const (
	LayoutTypeVertical LayoutType = iota
	LayoutTypeHorizontal
)

type Axis int

const (
	AxisNone Axis = -1
	AxisX    Axis = 0
	AxisY    Axis = 1
)

type PopupRef struct {
	PopupId        ID       // Set on OpenPopup()
	Window         *Window  // Resolved on BeginPopup() - may stay unresolved if user never calls OpenPopup()
	ParentWindow   *Window  // Set on OpenPopup()
	OpenFrameCount int      // Set on OpenPopup()
	OpenParentId   ID       // Set on OpenPopup(), we need this to differenciate multiple menu sets from each others (e.g. inside menu bar vs loose menu items)
	OpenPopupPos   f64.Vec2 // Set on OpenPopup(), preferred popup position (typically == OpenMousePos when using mouse)
	OpenMousePos   f64.Vec2 // Set on OpenPopup(), copy of mouse position at the time of opening popup
}

type NextWindowData struct {
	PosCond             Cond
	SizeCond            Cond
	ContentSizeCond     Cond
	CollapsedCond       Cond
	SizeConstraintCond  Cond
	FocusCond           Cond
	BgAlphaCond         Cond
	PosVal              f64.Vec2
	PosPivotVal         f64.Vec2
	SizeVal             f64.Vec2
	ContentSizeVal      f64.Vec2
	CollapsedVal        bool
	SizeConstraintRect  f64.Rectangle
	SizeCallback        func(*SizeCallbackData)
	BgAlphaVal          float64
	MenuBarOffsetMinVal f64.Vec2 // This is not exposed publicly, so we don't clear it.
}

type HoveredFlags int

const (
	HoveredFlagsDefault                 HoveredFlags = 0      // Return true if directly over the item/window not obstructed by another window not obstructed by an active popup or modal blocking inputs under them.
	HoveredFlagsChildWindows            HoveredFlags = 1 << 0 // IsWindowHovered() only: Return true if any children of the window is hovered
	HoveredFlagsRootWindow              HoveredFlags = 1 << 1 // IsWindowHovered() only: Test from root window (top most parent of the current hierarchy)
	HoveredFlagsAnyWindow               HoveredFlags = 1 << 2 // IsWindowHovered() only: Return true if any window is hovered
	HoveredFlagsAllowWhenBlockedByPopup HoveredFlags = 1 << 3 // Return true even if a popup window is normally blocking access to this item/window
	//HoveredFlagsAllowWhenBlockedByModal     HoveredFlags= 1 << 4   // Return true even if a modal popup window is normally blocking access to this item/window. FIXME-TODO: Unavailable yet.
	HoveredFlagsAllowWhenBlockedByActiveItem HoveredFlags = 1 << 5 // Return true even if an active item is blocking access to this item/window. Useful for Drag and Drop patterns.
	HoveredFlagsAllowWhenOverlapped          HoveredFlags = 1 << 6 // Return true even if the position is overlapped by another window
	HoveredFlagsRectOnly                     HoveredFlags = HoveredFlagsAllowWhenBlockedByPopup | HoveredFlagsAllowWhenBlockedByActiveItem | HoveredFlagsAllowWhenOverlapped
	HoveredFlagsRootAndChildWindows          HoveredFlags = HoveredFlagsRootWindow | HoveredFlagsChildWindows
)

type ItemStatusFlags int

const (
	ItemStatusFlagsHoveredRect    ItemStatusFlags = 1 << 0
	ItemStatusFlagsHasDisplayRect ItemStatusFlags = 1 << 1
)

func (w *Window) GetID(str string) ID {
	seed := w.IDStack[len(w.IDStack)-1]
	h := fnv.New32()
	binary.Write(h, binary.LittleEndian, uint32(seed))
	h.Write([]byte(str))
	id := ID(h.Sum32())
	w.Ctx.KeepAliveID(id)
	return id
}

func (w *Window) GetIntID(n int) ID {
	seed := w.IDStack[len(w.IDStack)-1]
	h := fnv.New32()
	binary.Write(h, binary.LittleEndian, uint32(seed))
	binary.Write(h, binary.LittleEndian, uint32(n))
	id := ID(h.Sum32())
	w.Ctx.KeepAliveID(id)
	return id
}

func (c *Context) ItemSize(size f64.Vec2) {
	c.ItemSizeEx(size, 0)
}

func (c *Context) ItemSizeEx(size f64.Vec2, text_offset_y float64) {
	window := c.CurrentWindow
	if window.SkipItems {
		return
	}

	// Always align ourselves on pixel boundaries
	line_height := math.Max(window.DC.CurrentLineHeight, size.Y)
	text_base_offset := math.Max(window.DC.CurrentLineTextBaseOffset, text_offset_y)
	window.DC.CursorPosPrevLine = f64.Vec2{window.DC.CursorPos.X + size.X, window.DC.CursorPos.Y}
	window.DC.CursorPos = f64.Vec2{
		float64(int(window.Pos.X + window.DC.IndentX + window.DC.ColumnsOffsetX)),
		float64(int(window.DC.CursorPos.Y + line_height + c.Style.ItemSpacing.Y)),
	}
	window.DC.CursorMaxPos.X = math.Max(window.DC.CursorMaxPos.X, window.DC.CursorPosPrevLine.X)
	window.DC.CursorMaxPos.Y = math.Max(window.DC.CursorMaxPos.Y, window.DC.CursorPos.Y-c.Style.ItemSpacing.Y)

	window.DC.PrevLineHeight = line_height
	window.DC.PrevLineTextBaseOffset = text_base_offset
	window.DC.CurrentLineHeight = 0
	window.DC.CurrentLineTextBaseOffset = 0

	// Horizontal layout mode
	if window.DC.LayoutType == LayoutTypeHorizontal {
		c.SameLine()
	}
}

func (c *Context) ItemSizeBB(bb f64.Rectangle) {
	c.ItemSizeBBEx(bb, 0)
}

func (c *Context) ItemSizeBBEx(bb f64.Rectangle, text_offset_y float64) {
	c.ItemSizeEx(bb.Size(), text_offset_y)
}

func (c *Context) ItemAdd(bb f64.Rectangle, id ID) bool {
	return c.ItemAddEx(bb, id, nil)
}

// Declare item bounding box for clipping and interaction.
// Note that the size can be different than the one provided to ItemSize(). Typically, widgets that spread over available surface
// declare their minimum size requirement to ItemSize() and then use a larger region for drawing/interaction, which is passed to ItemAdd().
func (c *Context) ItemAddEx(bb f64.Rectangle, id ID, nav_bb_arg *f64.Rectangle) bool {
	window := c.CurrentWindow

	if id != 0 {
		// Navigation processing runs prior to clipping early-out
		//  (a) So that NavInitRequest can be honored, for newly opened windows to select a default widget
		//  (b) So that we can scroll up/down past clipped items. This adds a small O(N) cost to regular navigation requests unfortunately, but it is still limited to one window.
		//      it may not scale very well for windows with ten of thousands of item, but at least NavMoveRequest is only set on user interaction, aka maximum once a frame.
		//      We could early out with "if (is_clipped && !g.NavInitRequest) return false;" but when we wouldn't be able to reach unclipped widgets. This would work if user had explicit scrolling control (e.g. mapped on a stick)
		window.DC.NavLayerActiveMaskNext |= window.DC.NavLayerCurrentMask
		if c.NavId == id || c.NavAnyRequest {
			if c.NavWindow.RootWindowForNav == window.RootWindowForNav {
				if window == c.NavWindow || (window.Flags|c.NavWindow.Flags)&WindowFlagsNavFlattened != 0 {
					if nav_bb_arg != nil {
						c.NavProcessItem(window, *nav_bb_arg, id)
					} else {
						c.NavProcessItem(window, bb, id)
					}

				}
			}
		}
	}

	window.DC.LastItemId = id
	window.DC.LastItemRect = bb
	window.DC.LastItemStatusFlags = 0

	// Clipping test
	is_clipped := c.IsClippedEx(bb, id, false)
	if is_clipped {
		return false
	}

	// We need to calculate this now to take account of the current clipping rectangle (as items like Selectable may change them)
	if c.IsMouseHoveringRect(bb.Min, bb.Max) {
		window.DC.LastItemStatusFlags |= ItemStatusFlagsHoveredRect
	}
	return true
}

func (c *Context) SameLine() {
	c.SameLineEx(0, -1)
}

// Gets back to previous line and continue with horizontal layout
//      pos_x == 0      : follow right after previous item
//      pos_x != 0      : align to specified x position (relative to window/group left)
//      spacing_w < 0   : use default spacing if pos_x == 0, no spacing if pos_x != 0
//      spacing_w >= 0  : enforce spacing amount
func (c *Context) SameLineEx(pos_x, spacing_w float64) {
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return
	}

	if pos_x != 0 {
		if spacing_w < 0 {
			spacing_w = 0
		}
		window.DC.CursorPos.X = window.Pos.X - window.Scroll.X + pos_x + spacing_w + window.DC.GroupOffsetX + window.DC.ColumnsOffsetX
		window.DC.CursorPos.Y = window.DC.CursorPosPrevLine.Y
	} else {
		if spacing_w < 0 {
			spacing_w = c.Style.ItemSpacing.X
		}
		window.DC.CursorPos.X = window.DC.CursorPosPrevLine.X + spacing_w
		window.DC.CursorPos.Y = window.DC.CursorPosPrevLine.Y
	}
	window.DC.CurrentLineHeight = window.DC.PrevLineHeight
	window.DC.CurrentLineTextBaseOffset = window.DC.PrevLineTextBaseOffset
}

func (c *Context) NewLine() {
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return
	}

	backup_layout_type := window.DC.LayoutType
	window.DC.LayoutType = LayoutTypeVertical
	if window.DC.CurrentLineHeight > 0.0 { // In the event that we are on a line with items that is smaller that FontSize high, we will preserve its height.
		c.ItemSizeEx(f64.Vec2{0, 0}, -1)
	} else {
		c.ItemSizeEx(f64.Vec2{0, c.FontSize}, -1)
	}
	window.DC.LayoutType = backup_layout_type
}

func (c *Context) NextColumn() {
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return
	}

	c.PopItemWidth()
	c.PopClipRect()

	columns := window.DC.ColumnsSet
	columns.CellMaxY = math.Max(columns.CellMaxY, window.DC.CursorPos.Y)
	if columns.Current++; columns.Current < columns.Count {
		// Columns 1+ cancel out IndentX
		window.DC.ColumnsOffsetX = c.GetColumnOffset(columns.Current) - window.DC.IndentX + c.Style.ItemSpacing.X
		window.DrawList.ChannelsSetCurrent(columns.Current)
	} else {
		window.DC.ColumnsOffsetX = 0
		window.DrawList.ChannelsSetCurrent(0)
		columns.Current = 0
		columns.CellMinY = columns.CellMaxY
	}

	window.DC.CursorPos.X = float64(int(window.Pos.X + window.DC.IndentX + window.DC.ColumnsOffsetX))
	window.DC.CursorPos.Y = columns.CellMinY
	window.DC.CurrentLineHeight = 0
	window.DC.CurrentLineTextBaseOffset = 0

	c.PushColumnClipRect()
	c.PushItemWidth(c.GetColumnWidth() * 0.65) // FIXME: Move on columns setup
}

func (c *Context) PushColumnClipRect() {
	c.PushColumnClipRectEx(-1)
}

func (c *Context) PushColumnClipRectEx(column_index int) {
	window := c.GetCurrentWindowRead()
	columns := window.DC.ColumnsSet
	if column_index < 0 {
		column_index = columns.Current
	}

	c.PushClipRect(
		columns.Columns[column_index].ClipRect.Min,
		columns.Columns[column_index].ClipRect.Max,
		false,
	)
}

// When using this function it is sane to ensure that float are perfectly rounded to integer values, to that e.g. (int)(max.x-min.x) in user's render produce correct result.
func (c *Context) PushClipRect(clip_rect_min, clip_rect_max f64.Vec2, intersect_with_current_clip_rect bool) {
	window := c.GetCurrentWindow()
	window.DrawList.PushClipRectEx(clip_rect_min, clip_rect_max, intersect_with_current_clip_rect)
	length := len(window.DrawList.ClipRectStack)
	clipRect := window.DrawList.ClipRectStack[length-1]
	window.ClipRect = f64.Rectangle{f64.Vec2{clipRect.X, clipRect.Y}, f64.Vec2{clipRect.Z, clipRect.W}}
}

func (c *Context) PushItemWidth(item_width float64) {
	window := c.GetCurrentWindow()
	window.DC.ItemWidth = item_width
	if window.DC.ItemWidth == 0 {
		window.DC.ItemWidth = window.ItemWidthDefault
	}
	window.DC.ItemWidthStack = append(window.DC.ItemWidthStack, window.DC.ItemWidth)
}

func (c *Context) GetColumnWidth() float64 {
	return c.GetColumnWidthDx(-1)
}

func (c *Context) GetColumnWidthDx(column_index int) float64 {
	window := c.GetCurrentWindowRead()
	columns := window.DC.ColumnsSet
	if column_index < 0 {
		column_index = columns.Current
	}
	return c.OffsetNormToPixels(
		columns,
		columns.Columns[column_index+1].OffsetNorm-columns.Columns[column_index].OffsetNorm,
	)
}

func (c *Context) OffsetNormToPixels(columns *ColumnsSet, offset_norm float64) float64 {
	return offset_norm * (columns.MaxX - columns.MinX)
}

func (c *Context) GetColumnWidthEx(columns *ColumnsSet, column_index int, before_resize bool) float64 {
	if column_index < 0 {
		column_index = columns.Current
	}

	var offset_norm float64
	if before_resize {
		offset_norm = columns.Columns[column_index+1].OffsetNormBeforeResize - columns.Columns[column_index].OffsetNormBeforeResize
	} else {
		offset_norm = columns.Columns[column_index+1].OffsetNorm - columns.Columns[column_index].OffsetNorm
	}

	return c.OffsetNormToPixels(columns, offset_norm)
}

func (c *Context) GetContentRegionAvail() f64.Vec2 {
	window := c.GetCurrentWindowRead()
	regionMax := c.GetContentRegionMax()
	windowRegion := window.DC.CursorPos.Sub(window.Pos)
	return regionMax.Sub(windowRegion)
}

func (c *Context) GetContentRegionMax() f64.Vec2 {
	window := c.GetCurrentWindowRead()
	mx := window.ContentsRegionRect.Max
	if window.DC.ColumnsSet != nil {
		mx.X = c.GetColumnOffset(window.DC.ColumnsSet.Current+1) - window.WindowPadding.X
	}
	return mx
}

func (c *Context) GetContentRegionAvailWidth() float64 {
	return c.GetContentRegionAvail().X
}

func (c *Context) GetWindowContentRegionMin() f64.Vec2 {
	window := c.GetCurrentWindowRead()
	return window.ContentsRegionRect.Min
}

func (c *Context) GetWindowContentRegionMax() f64.Vec2 {
	window := c.GetCurrentWindowRead()
	return window.ContentsRegionRect.Max
}

func (c *Context) GetWindowContentRegionWidth() float64 {
	window := c.GetCurrentWindowRead()
	return window.ContentsRegionRect.Max.X - window.ContentsRegionRect.Min.X
}

func (c *Context) GetColumnOffset(column_index int) float64 {
	window := c.GetCurrentWindowRead()
	columns := window.DC.ColumnsSet

	if column_index < 0 {
		column_index = columns.Current
	}

	t := columns.Columns[column_index].OffsetNorm
	x_offset := f64.Lerp(t, columns.MinX, columns.MaxX)
	return x_offset
}

func (c *Context) GetColumnIndex() int {
	window := c.GetCurrentWindowRead()
	if window.DC.ColumnsSet != nil {
		return window.DC.ColumnsSet.Current
	}
	return 0
}

func (c *Context) GetColumnsCount() int {
	window := c.GetCurrentWindowRead()
	if window.DC.ColumnsSet != nil {
		return window.DC.ColumnsSet.Count
	}
	return 0
}

func (c *Context) PopItemWidth() {
	window := c.GetCurrentWindow()
	length := len(window.DC.ItemWidthStack)
	window.DC.ItemWidthStack = window.DC.ItemWidthStack[:length-1]
	if length--; length == 0 {
		window.DC.ItemWidth = window.ItemWidthDefault
	} else {
		window.DC.ItemWidth = window.DC.ItemWidthStack[length-1]
	}
}

func (c *Context) CalcItemWidth() float64 {
	window := c.GetCurrentWindowRead()
	w := window.DC.ItemWidth
	if w < 0 {
		// Align to a right-side limit. We include 1 frame padding in the calculation because this is how the width is always used (we add 2 frame padding to it), but we could move that responsibility to the widget as well.
		width_to_right_edge := c.GetContentRegionAvail().X
		w = math.Max(1, width_to_right_edge+w)
	}
	w = float64(int(w))
	return w
}

func (c *Context) PopClipRect() {
	window := c.GetCurrentWindow()
	window.DrawList.PopClipRect()
	length := len(window.DrawList.ClipRectStack)
	clipRect := window.DrawList.ClipRectStack[length-1]
	window.ClipRect = f64.Rectangle{f64.Vec2{clipRect.X, clipRect.Y}, f64.Vec2{clipRect.Z, clipRect.W}}
}

func (c *Context) CalcItemSize(size f64.Vec2, default_x, default_y float64) f64.Vec2 {
	var content_max f64.Vec2
	if size.X < 0 || size.Y < 0 {
		content_max = c.CurrentWindow.Pos.Add(c.GetContentRegionMax())
	}
	if size.X == 0 {
		size.X = default_x
	} else {
		size.X += math.Max(content_max.X-c.CurrentWindow.DC.CursorPos.X, 4)
	}

	if size.Y == 0 {
		size.Y = default_y
	} else {
		size.Y += math.Max(content_max.Y-c.CurrentWindow.DC.CursorPos.Y, 4)
	}

	return size
}

func (c *Context) CalcWrapWidthForPos(pos f64.Vec2, wrap_pos_x float64) float64 {
	if wrap_pos_x < 0 {
		return 0
	}

	window := c.GetCurrentWindowRead()
	if wrap_pos_x == 0 {
		wrap_pos_x = c.GetContentRegionMax().X + window.Pos.X
	} else if wrap_pos_x > 0 {
		wrap_pos_x += window.Pos.X - window.Scroll.X // wrap_pos_x is provided is window local space
	}

	return math.Max(wrap_pos_x-pos.X, 1)
}

func (w *Window) CalcFontSize() float64 {
	return w.Ctx.FontBaseSize * w.FontWindowScale
}

func (c *Context) ItemHoverable(bb f64.Rectangle, id ID) bool {
	if c.HoveredId != 0 && c.HoveredId != id && !c.HoveredIdAllowOverlap {
		return false
	}

	window := c.CurrentWindow
	if c.HoveredWindow != window {
		return false
	}

	if c.ActiveId != 0 && c.ActiveId != id && !c.ActiveIdAllowOverlap {
		return false
	}

	if !c.IsMouseHoveringRect(bb.Min, bb.Max) {
		return false
	}

	if c.NavDisableMouseHover || !c.IsWindowContentHoverable(window, HoveredFlagsDefault) {
		return false
	}

	if window.DC.ItemFlags&ItemFlagsDisabled != 0 {
		return false
	}

	c.SetHoveredID(id)
	return true
}

func (c *Context) IsItemHovered() bool {
	return c.IsItemHoveredEx(0)
}

// This is roughly matching the behavior of internal-facing ItemHoverable()
// - we allow hovering to be true when ActiveId==window->MoveID, so that clicking on non-interactive items such as a Text() item still returns true with IsItemHovered()
// - this should work even for non-interactive items that have no ID, so we cannot use LastItemId
func (c *Context) IsItemHoveredEx(flags HoveredFlags) bool {
	window := c.CurrentWindow
	if c.NavDisableMouseHover && !c.NavDisableHighlight {
		return c.IsItemFocused()
	}

	// Test for bounding box overlap, as updated as ItemAdd()
	if window.DC.LastItemStatusFlags&ItemStatusFlagsHoveredRect == 0 {
		return false
	}
	// Flags not supported by this function
	assert((flags & (HoveredFlagsRootWindow | HoveredFlagsChildWindows)) == 0)

	// Test if we are hovering the right window (our window could be behind another window)
	// [2017/10/16] Reverted commit 344d48be3 and testing RootWindow instead. I believe it is correct to NOT test for RootWindow but this leaves us unable to use IsItemHovered() after EndChild() itself.
	// Until a solution is found I believe reverting to the test from 2017/09/27 is safe since this was the test that has been running for a long while.
	//if (c.HoveredWindow != window) {
	//    return false;
	//}
	if c.HoveredRootWindow != window.RootWindow && flags&HoveredFlagsAllowWhenOverlapped == 0 {
		return false
	}

	// Test if another item is active (e.g. being dragged)
	if flags&HoveredFlagsAllowWhenBlockedByActiveItem == 0 {
		if c.ActiveId != 0 && c.ActiveId != window.DC.LastItemId && !c.ActiveIdAllowOverlap && c.ActiveId != window.MoveId {
			return false
		}
	}

	// Test if interactions on this window are blocked by an active popup or modal
	if !c.IsWindowContentHoverable(window, flags) {
		return false
	}

	// Test if the item is disabled
	if window.DC.ItemFlags&ItemFlagsDisabled != 0 {
		return false
	}

	// Special handling for the 1st item after Begin() which represent the title bar. When the window is collapsed (SkipItems==true) that last item will never be overwritten so we need to detect tht case.
	if window.DC.LastItemId == window.MoveId && window.WriteAccessed {
		return false
	}

	return true
}

func (c *Context) IsMouseHoveringRect(r_min, r_max f64.Vec2) bool {
	return c.IsMouseHoveringRectEx(r_min, r_max, true)
}

// Test if mouse cursor is hovering given rectangle
// NB- Rectangle is clipped by our current clip setting
// NB- Expand the rectangle to be generous on imprecise inputs systems (g.Style.TouchExtraPadding)
func (c *Context) IsMouseHoveringRectEx(r_min, r_max f64.Vec2, clip bool) bool {
	window := c.CurrentWindow

	rect_clipped := f64.Rectangle{r_min, r_max}
	if clip {
		rect_clipped.Intersect(window.ClipRect)
	}

	// Expand for touch input
	rect_for_touch := f64.Rectangle{
		rect_clipped.Min.Sub(c.Style.TouchExtraPadding),
		rect_clipped.Max.Add(c.Style.TouchExtraPadding),
	}
	return c.IO.MousePos.In(rect_for_touch)
}

func (c *Context) IsWindowContentHoverable(window *Window, flags HoveredFlags) bool {
	if c.NavWindow != nil {
		focused_root_window := c.NavWindow.RootWindow
		if focused_root_window != nil {
			if focused_root_window.WasActive && focused_root_window != window.RootWindow {
				// For the purpose of those flags we differentiate "standard popup" from "modal popup"
				// NB: The order of those two tests is important because Modal windows are also Popups.
				if focused_root_window.Flags&WindowFlagsModal != 0 {
					return false
				}
				if focused_root_window.Flags&WindowFlagsPopup != 0 && flags&HoveredFlagsAllowWhenBlockedByPopup == 0 {
					return false
				}
			}
		}
	}
	return true
}

func (c *Context) IsItemFocused() bool {
	return c.NavId != 0 && !c.NavDisableHighlight && c.NavId == c.CurrentWindow.DC.LastItemId
}

func (c *Context) IsItemClicked(mouse_button int) bool {
	return c.IsMouseClicked(mouse_button, false) && c.IsItemHoveredEx(HoveredFlagsDefault)
}

func (c *Context) IsAnyItemHovered() bool {
	return c.HoveredId != 0 || c.HoveredIdPreviousFrame != 0
}

func (c *Context) IsMouseClicked(button int, repeat bool) bool {
	t := c.IO.MouseDownDuration[button]
	if t == 0 {
		return true
	}

	if repeat && t > c.IO.KeyRepeatDelay {
		delay := c.IO.KeyRepeatDelay
		rate := c.IO.KeyRepeatRate
		mod1 := math.Mod(t-delay, rate) > rate*0.5
		mod2 := math.Mod(t-delay-c.IO.DeltaTime, rate) > rate*0.5
		if mod1 != mod2 {
			return true
		}
	}

	return false
}

func (c *Context) IsMouseReleased(button int) bool {
	return c.IO.MouseReleased[button]
}

func (c *Context) IsMouseDoubleClicked(button int) bool {
	return c.IO.MouseDoubleClicked[button]
}

func (c *Context) IsMouseDragging() bool {
	return c.IsMouseDraggingEx(0, -1)
}

func (c *Context) IsMouseDraggingEx(button int, lock_threshold float64) bool {
	if !c.IO.MouseDown[button] {
		return false
	}
	if lock_threshold < 0.0 {
		lock_threshold = c.IO.MouseDragThreshold
	}
	return c.IO.MouseDragMaxDistanceSqr[button] >= lock_threshold*lock_threshold
}

func (c *Context) GetMousePos() f64.Vec2 {
	return c.IO.MousePos
}

func (c *Context) IsMousePosValid() bool {
	return c.IsMousePosValidEx(nil)
}

func (c *Context) IsMousePosValidEx(mouse_pos *f64.Vec2) bool {
	if mouse_pos == nil {
		mouse_pos = &c.IO.MousePos
	}
	const MOUSE_INVALID = -256000.0
	return mouse_pos.X >= MOUSE_INVALID && mouse_pos.Y >= MOUSE_INVALID
}

func (c *Context) GetMouseDragDelta(button int, lock_threshold float64) f64.Vec2 {
	if lock_threshold < 0 {
		lock_threshold = c.IO.MouseDragThreshold
	}
	if c.IO.MouseDown[button] {
		if c.IO.MouseDragMaxDistanceSqr[button] >= lock_threshold*lock_threshold {
			return c.IO.MousePos.Sub(c.IO.MouseClickedPos[button]) // Assume we can only get active with left-mouse button (at the moment).
		}
	}
	return f64.Vec2{0, 0}
}

func (c *Context) ResetMouseDragDelta(button int) {
	c.IO.MouseClickedPos[button] = c.IO.MousePos
}

func (c *Context) GetMouseCursor() MouseCursor {
	return c.MouseCursor
}

func (c *Context) SetMouseCursor(cursor_type MouseCursor) {
	c.MouseCursor = cursor_type
}

func (c *Context) IsItemActive() bool {
	if c.ActiveId != 0 {
		window := c.CurrentWindow
		return c.ActiveId == window.DC.LastItemId
	}
	return false
}

func (c *Context) IsAnyItemFocused() bool {
	return c.NavId != 0 && !c.NavDisableHighlight
}

func (c *Context) IsItemVisible() bool {
	window := c.GetCurrentWindowRead()
	return window.ClipRect.Overlaps(window.DC.LastItemRect)
}

// Allow last item to be overlapped by a subsequent item. Both may be activated during the same frame before the later one takes priority.
func (c *Context) SetItemAllowOverlap() {
	if c.HoveredId == c.CurrentWindow.DC.LastItemId {
		c.HoveredIdAllowOverlap = true
	}
	if c.ActiveId == c.CurrentWindow.DC.LastItemId {
		c.ActiveIdAllowOverlap = true
	}
}

func (c *Context) GetViewportRect() f64.Rectangle {
	if c.IO.DisplayVisibleMin.X != c.IO.DisplayVisibleMax.X && c.IO.DisplayVisibleMin.Y != c.IO.DisplayVisibleMax.Y {
		return f64.Rectangle{c.IO.DisplayVisibleMin, c.IO.DisplayVisibleMax}
	}
	return f64.Rect(0, 0, c.IO.DisplaySize.X, c.IO.DisplaySize.Y)
}

// Moving window to front of display and set focus (which happens to be back of our sorted list)
func (c *Context) FocusWindow(window *Window) {
	if c.NavWindow != window {
		c.NavWindow = window
		if window != nil && c.NavDisableMouseHover {
			c.NavMousePosDirty = true
		}
		c.NavInitRequest = false
		c.NavId = 0
		if window != nil {
			c.NavId = window.NavLastIds[0]
		}
		c.NavIdIsAlive = false
		c.NavLayer = 0
	}

	// Passing NULL allow to disable keyboard focus
	if window == nil {
		return
	}

	// Move the root window to the top of the pile
	if window.RootWindow != nil {
		window = window.RootWindow
	}

	// Steal focus on active widgets
	// FIXME: This statement should be unnecessary. Need further testing before removing it..
	if window.Flags&WindowFlagsPopup != 0 {
		if c.ActiveId != 0 && c.ActiveIdWindow != nil && c.ActiveIdWindow.RootWindow != window {
			c.ClearActiveID()
		}
	}

	// Bring to front
	if window.Flags&WindowFlagsNoBringToFrontOnFocus == 0 {
		c.BringWindowToFront(window)
	}
}

func (c *Context) BringWindowToFront(window *Window) {
	current_front_window := c.Windows[len(c.Windows)-1]
	if current_front_window == window || current_front_window.RootWindow == window {
		return
	}
	// We can ignore the front most window
	for i := len(c.Windows) - 2; i >= 0; i-- {
		if c.Windows[i] == window {
			c.Windows = append(c.Windows[:i], c.Windows[i+1:]...)
			c.Windows = append(c.Windows, window)
			break
		}
	}
}

func (c *Context) SetCurrentWindow(window *Window) {
	c.CurrentWindow = window
	if window != nil {
		c.FontSize = window.CalcFontSize()
		c.DrawListSharedData.FontSize = c.FontSize
	}
}

func (c *Context) GetColumnsRectHalfWidth() float64 {
	return 4
}

func (c *Context) IsClippedEx(bb f64.Rectangle, id ID, clip_even_when_logged bool) bool {
	window := c.CurrentWindow
	if !bb.Overlaps(window.ClipRect) {
		if id == 0 || id != c.ActiveId {
			if clip_even_when_logged || !c.LogEnabled {
				return true
			}
		}
	}
	return false
}

func (c *Context) GetDraggedColumnOffset(columns *ColumnsSet, column_index int) float64 {
	// Active (dragged) column always follow mouse. The reason we need this is that dragging a column to the right edge of an auto-resizing
	// window creates a feedback loop because we store normalized positions. So while dragging we enforce absolute positioning.
	window := c.CurrentWindow
	x := c.IO.MousePos.X - c.ActiveIdClickOffset.X + c.GetColumnsRectHalfWidth() - window.Pos.X
	x = math.Max(x, c.GetColumnOffset(column_index-1)+c.Style.ColumnsMinSpacing)
	if columns.Flags&ColumnsFlagsNoPreserveWidths != 0 {
		x = math.Min(x, c.GetColumnOffset(column_index+1)-c.Style.ColumnsMinSpacing)
	}

	return x
}

func (c *Context) SetColumnOffset(column_index int, offset float64) {
	window := c.CurrentWindow
	columns := window.DC.ColumnsSet

	if column_index < 0 {
		column_index = columns.Current
	}

	preserve_width := columns.Flags&ColumnsFlagsNoPreserveWidths == 0 && column_index < columns.Count-1
	width := 0.0
	if preserve_width {
		width = c.GetColumnWidthEx(columns, column_index, columns.IsBeingResized)
	}

	if columns.Flags&ColumnsFlagsNoForceWithinWindow == 0 {
		offset = math.Min(offset, columns.MaxX-c.Style.ColumnsMinSpacing*float64(columns.Count-column_index))
	}
	columns.Columns[column_index].OffsetNorm = c.PixelsToOffsetNorm(columns, offset-columns.MinX)

	if preserve_width {
		c.SetColumnOffset(column_index+1, offset+math.Max(c.Style.ColumnsMinSpacing, width))
	}
}

func (c *Context) PixelsToOffsetNorm(columns *ColumnsSet, offset float64) float64 {
	return offset / (columns.MaxX - columns.MinX)
}

func (c *Context) GetFrontMostPopupModal() *Window {
	for n := len(c.OpenPopupStack) - 1; n >= 0; n-- {
		popup := c.OpenPopupStack[n].Window
		if popup != nil && popup.Flags&WindowFlagsModal != 0 {
			return popup
		}
	}
	return nil
}

// Call when we are expected to land on Layer 0 after FocusWindow()
func (c *Context) NavRestoreLastChildNavWindow(window *Window) *Window {
	if window.NavLastChildNavWindow != nil {
		return window.NavLastChildNavWindow
	}
	return window
}

func (c *Context) IsRectVisible(rect_min, rect_max f64.Vec2) bool {
	window := c.GetCurrentWindowRead()
	return window.ClipRect.Overlaps(f64.Rectangle{rect_min, rect_max})
}

// Lock horizontal starting position + capture group bounding box into one "item" (so you can use IsItemHovered() or layout primitives such as SameLine() on whole group, etc.)
func (c *Context) BeginGroup() {
	window := c.GetCurrentWindow()

	group_data := GroupData{
		BackupCursorPos:                 window.DC.CursorPos,
		BackupCursorMaxPos:              window.DC.CursorMaxPos,
		BackupIndentX:                   window.DC.IndentX,
		BackupGroupOffsetX:              window.DC.GroupOffsetX,
		BackupCurrentLineHeight:         window.DC.CurrentLineHeight,
		BackupCurrentLineTextBaseOffset: window.DC.CurrentLineTextBaseOffset,
		BackupLogLinePosY:               window.DC.LogLinePosY,
		BackupActiveIdIsAlive:           c.ActiveIdIsAlive,
		AdvanceCursor:                   true,
	}
	window.DC.GroupStack = append(window.DC.GroupStack, group_data)

	window.DC.GroupOffsetX = window.DC.CursorPos.X - window.Pos.X - window.DC.ColumnsOffsetX
	window.DC.IndentX = window.DC.GroupOffsetX
	window.DC.CursorMaxPos = window.DC.CursorPos
	window.DC.CurrentLineHeight = 0
	window.DC.LogLinePosY = window.DC.CursorPos.Y - 9999
}

func (c *Context) EndGroup() {
	window := c.GetCurrentWindow()
	length := len(window.DC.GroupStack)
	group_data := &window.DC.GroupStack[length-1]

	group_bb := f64.Rectangle{group_data.BackupCursorPos, window.DC.CursorMaxPos}
	group_bb.Max = f64.Vec2{
		math.Max(group_bb.Min.X, group_bb.Max.X),
		math.Max(group_bb.Min.Y, group_bb.Max.Y),
	}

	window.DC.CursorPos = group_data.BackupCursorPos
	window.DC.CursorMaxPos = f64.Vec2{
		math.Max(group_data.BackupCursorMaxPos.X, window.DC.CursorMaxPos.X),
		math.Max(group_data.BackupCursorMaxPos.Y, window.DC.CursorMaxPos.Y),
	}
	window.DC.CurrentLineHeight = group_data.BackupCurrentLineHeight
	window.DC.CurrentLineTextBaseOffset = group_data.BackupCurrentLineTextBaseOffset
	window.DC.IndentX = group_data.BackupIndentX
	window.DC.GroupOffsetX = group_data.BackupGroupOffsetX
	window.DC.LogLinePosY = window.DC.CursorPos.Y - 9999

	if group_data.AdvanceCursor {
		// FIXME: Incorrect, we should grab the base offset from the *first line* of the group but it is hard to obtain now.
		window.DC.CurrentLineTextBaseOffset = math.Max(
			window.DC.PrevLineTextBaseOffset,
			group_data.BackupCurrentLineTextBaseOffset,
		)
		c.ItemSizeEx(group_bb.Size(), group_data.BackupCurrentLineTextBaseOffset)
		c.ItemAdd(group_bb, 0)
	}

	// If the current ActiveId was declared within the boundary of our group, we copy it to LastItemId so IsItemActive() will be functional on the entire group.
	// It would be be neater if we replaced window.DC.LastItemId by e.g. 'bool LastItemIsActive', but if you search for LastItemId you'll notice it is only used in that context.
	active_id_within_group := !group_data.BackupActiveIdIsAlive && c.ActiveIdIsAlive && c.ActiveId != 0 && c.ActiveIdWindow.RootWindow == window.RootWindow
	if active_id_within_group {
		window.DC.LastItemId = c.ActiveId
	}
	window.DC.LastItemRect = group_bb

	window.DC.GroupStack = window.DC.GroupStack[:len(window.DC.GroupStack)-1]
}

func (c *Context) Indent() {
	c.IndentEx(0)
}

func (c *Context) Unindent() {
	c.UnindentEx(0)
}

func (c *Context) IndentEx(indent_w float64) {
	window := c.GetCurrentWindow()
	if indent_w != 0 {
		window.DC.IndentX += indent_w
	} else {
		window.DC.IndentX += c.Style.IndentSpacing
	}
	window.DC.CursorPos.X = window.Pos.X + window.DC.IndentX + window.DC.ColumnsOffsetX
}

func (c *Context) UnindentEx(indent_w float64) {
	window := c.GetCurrentWindow()
	if indent_w != 0 {
		window.DC.IndentX -= indent_w
	} else {
		window.DC.IndentX -= c.Style.IndentSpacing
	}
	window.DC.CursorPos.X = window.Pos.X + window.DC.IndentX + window.DC.ColumnsOffsetX
}

func (c *Context) TreePush(str_id string) {
	window := c.GetCurrentWindow()
	c.Indent()
	window.DC.TreeDepth++
	if str_id == "" {
		str_id = "#TreePush"
	}
	c.PushStringID(str_id)
}

func (c *Context) AddWindowToSortedBuffer(out_sorted_windows *[]*Window, window *Window) {
	*out_sorted_windows = append(*out_sorted_windows, window)
	if window.Active {
		if len(window.DC.ChildWindows) > 1 {
			sort.Slice(window.DC.ChildWindows, func(i, j int) bool {
				// FIXME: Add a more explicit sort order in the window structure.
				a := window.DC.ChildWindows[i]
				b := window.DC.ChildWindows[j]
				d := (a.Flags & WindowFlagsPopup) - (b.Flags & WindowFlagsPopup)
				if d != 0 {
					return d < 0
				}
				d = (a.Flags & WindowFlagsTooltip) - (b.Flags & WindowFlagsTooltip)
				if d != 0 {
					return d < 0
				}
				return a.BeginOrderWithinParent < b.BeginOrderWithinParent
			})
		}

		for _, child := range window.DC.ChildWindows {
			if child.Active {
				c.AddWindowToSortedBuffer(out_sorted_windows, child)
			}
		}
	}
}

func (c *Context) ClearDragDrop() {
	c.DragDropActive = false
	c.DragDropPayload.Clear()
	c.DragDropAcceptIdCurr = 0
	c.DragDropAcceptIdPrev = 0
	c.DragDropAcceptIdCurrRectSurface = math.MaxFloat32
	c.DragDropAcceptFrameCount = -1
}

func (p *Payload) Clear() {
	p.SourceId = 0
	p.SourceParentId = 0
	p.Data = nil
	p.DataFrameCount = -1
	p.Preview = false
	p.Delivery = false
}

func (d *DrawData) Clear() {
	d.Valid = false
	d.CmdLists = nil
	d.TotalVtxCount = 0
	d.TotalIdxCount = 0
}

func (w *Window) Init(ctx *Context, name string) {
	h := fnv.New32()
	h.Write([]byte(name))
	id := ID(h.Sum32())

	w.Ctx = ctx
	w.Name = name
	w.IDStack = append(w.IDStack, id)
	w.Flags = 0
	w.Pos = f64.Vec2{0, 0}
	w.Size = f64.Vec2{0, 0}
	w.SizeFull = f64.Vec2{0, 0}
	w.SizeContents = f64.Vec2{0, 0}
	w.SizeContentsExplicit = f64.Vec2{0, 0}
	w.WindowPadding = f64.Vec2{0, 0}
	w.WindowRounding = 0.0
	w.WindowBorderSize = 0.0
	w.MoveId = w.GetID("#MOVE")
	w.ChildId = 0
	w.Scroll = f64.Vec2{0, 0}
	w.ScrollTarget = f64.Vec2{math.MaxFloat32, math.MaxFloat32}
	w.ScrollTargetCenterRatio = f64.Vec2{0.5, 0.5}
	w.ScrollbarSizes = f64.Vec2{0.0, 0.0}
	w.ScrollbarX, w.ScrollbarY = false, false
	w.Active = false
	w.WasActive = false
	w.WriteAccessed = false
	w.Collapsed = false
	w.CollapseToggleWanted = false
	w.SkipItems = false
	w.Appearing = false
	w.CloseButton = false
	w.BeginOrderWithinParent = -1
	w.BeginOrderWithinContext = -1
	w.BeginCount = 0
	w.PopupId = 0
	w.AutoFitFramesX = -1
	w.AutoFitFramesY = -1
	w.AutoFitOnlyGrows = false
	w.AutoFitChildAxises = 0x00
	w.AutoPosLastDirection = DirNone
	w.HiddenFrames = 0
	w.SetWindowPosAllowFlags = CondAlways | CondOnce | CondFirstUseEver | CondAppearing
	w.SetWindowSizeAllowFlags = w.SetWindowPosAllowFlags
	w.SetWindowCollapsedAllowFlags = w.SetWindowPosAllowFlags
	w.SetWindowPosVal = f64.Vec2{math.MaxFloat32, math.MaxFloat32}
	w.SetWindowPosPivot = f64.Vec2{math.MaxFloat32, math.MaxFloat32}

	w.StateStorage = make(map[ID]interface{})

	w.LastFrameActive = -1
	w.ItemWidthDefault = 0.0
	w.FontWindowScale = 1.0

	w.DrawListInst.Init(&w.Ctx.DrawListSharedData)
	w.DrawList = &w.DrawListInst
	w.DrawList.OwnerName = w.Name
	w.ParentWindow = nil
	w.RootWindow = nil
	w.RootWindowForTitleBarHighlight = nil
	w.RootWindowForTabbing = nil
	w.RootWindowForNav = nil

	w.NavLastIds[0] = 0
	w.NavLastIds[1] = 0
	w.NavRectRel[0] = f64.Rectangle{}
	w.NavRectRel[1] = f64.Rectangle{}
	w.NavLastChildNavWindow = nil

	w.FocusIdxAllCounter = -1
	w.FocusIdxTabCounter = -1
	w.FocusIdxAllRequestCurrent = math.MaxInt32
	w.FocusIdxTabRequestCurrent = math.MaxInt32
	w.FocusIdxAllRequestNext = math.MaxInt32
	w.FocusIdxTabRequestNext = math.MaxInt32
}

func (c *Context) CreateNewWindow(name string, size f64.Vec2, flags WindowFlags) *Window {
	window := &Window{}
	window.Init(c, name)
	window.Flags = flags
	c.WindowsById[name] = window

	// Use SetWindowPos() or SetNextWindowPos() with the appropriate condition flag to change the initial position of a window.
	window.Pos = f64.Vec2{60, 60}

	// User can disable loading and saving of settings. Tooltip and child windows also don't store settings.
	if flags&WindowFlagsNoSavedSettings == 0 {
		// Retrieve settings from .ini file

		settings := c.FindWindowSettings(window.Name)
		if settings != nil {
			c.SetWindowConditionAllowFlags(window, CondFirstUseEver, false)
			window.Pos = settings.Pos.Floor()
			window.Collapsed = settings.Collapsed
			if settings.Size.LenSquared() > 0.00001 {
				size = settings.Size
			}
		}
	}
	window.Size = size
	window.SizeFull = size
	window.SizeFullAtLastBegin = size

	if flags&WindowFlagsAlwaysAutoResize != 0 {
		window.AutoFitFramesX = 2
		window.AutoFitFramesY = 2
		window.AutoFitOnlyGrows = false
	} else {
		if window.Size.X <= 0 {
			window.AutoFitFramesX = 2
		}
		if window.Size.Y <= 0 {
			window.AutoFitFramesY = 2
		}
		window.AutoFitOnlyGrows = window.AutoFitFramesX > 0 || window.AutoFitFramesY > 0
	}

	if flags&WindowFlagsNoBringToFrontOnFocus != 0 {
		// Quite slow but rare and only once
		c.Windows = append([]*Window{window}, c.Windows...)
	} else {
		c.Windows = append(c.Windows, window)
	}

	return window
}

func (c *Context) CalcSizeContents(window *Window) f64.Vec2 {
	sz := window.SizeContentsExplicit
	if window.SizeContentsExplicit.X == 0 {
		sz.X = window.DC.CursorMaxPos.X - window.Pos.X + window.Scroll.X
	}
	if window.SizeContentsExplicit.Y == 0 {
		sz.Y = window.DC.CursorMaxPos.Y - window.Pos.Y + window.Scroll.Y
	}
	sz = sz.Add(window.WindowPadding)
	return sz
}

func (c *Context) GetWindowScrollMaxX(window *Window) float64 {
	return math.Max(0, window.SizeContents.X-(window.SizeFull.X-window.ScrollbarSizes.X))
}

func (c *Context) GetWindowScrollMaxY(window *Window) float64 {
	return math.Max(0, window.SizeContents.Y-(window.SizeFull.Y-window.ScrollbarSizes.Y))
}

func (c *Context) GetScrollMaxX() float64 {
	return c.GetWindowScrollMaxX(c.CurrentWindow)
}

func (c *Context) GetScrollMaxY() float64 {
	return c.GetWindowScrollMaxY(c.CurrentWindow)
}

func (c *Context) GetWindowBgColorIdxFromFlags(flags WindowFlags) Col {
	if flags&(WindowFlagsTooltip|WindowFlagsPopup) != 0 {
		return ColPopupBg
	}
	if flags&WindowFlagsChildWindow != 0 {
		return ColChildBg
	}
	return ColWindowBg
}

type ResizeGripDef struct {
	CornerPos              f64.Vec2
	InnerDir               f64.Vec2
	AngleMin12, AngleMax12 int
}

var resize_grip_def = [4]ResizeGripDef{
	{f64.Vec2{1, 1}, f64.Vec2{-1, -1}, 0, 3},  // Lower right
	{f64.Vec2{0, 1}, f64.Vec2{+1, -1}, 3, 6},  // Lower left
	{f64.Vec2{0, 0}, f64.Vec2{+1, +1}, 6, 9},  // Upper left
	{f64.Vec2{1, 0}, f64.Vec2{-1, +1}, 9, 12}, // Upper right
}

func (c *Context) PushItemFlag(option ItemFlags, enabled bool) {
	window := c.GetCurrentWindow()
	if enabled {
		window.DC.ItemFlags |= option
	} else {
		window.DC.ItemFlags &^= option
	}
	window.DC.ItemFlagsStack = append(window.DC.ItemFlagsStack, window.DC.ItemFlags)
}

func (c *Context) PopItemFlag() {
	window := c.GetCurrentWindow()
	length := len(window.DC.ItemFlagsStack) - 1
	window.DC.ItemFlagsStack = window.DC.ItemFlagsStack[:length]
	window.DC.ItemFlags = ItemFlagsDefault_
	if length > 0 {
		window.DC.ItemFlags = window.DC.ItemFlagsStack[length-1]
	}
}

func (c *Context) UpdateMovingWindow() {
	if c.MovingWindow != nil {
		// We actually want to move the root window. g.MovingWindow == window we clicked on (could be a child window).
		// We track it to preserve Focus and so that generally ActiveIdWindow == MovingWindow and ActiveId == MovingWindow->MoveId for consistency.
		c.KeepAliveID(c.ActiveId)
		moving_window := c.MovingWindow.RootWindow
		if c.IO.MouseDown[0] && c.IsMousePosValidEx(&c.IO.MousePos) {
			pos := c.IO.MousePos.Sub(c.ActiveIdClickOffset)
			if moving_window.Pos.X != pos.X || moving_window.Pos.Y != pos.Y {
				c.MarkIniSettingsDirtyForWindow(moving_window)
				c.SetWindowPos(moving_window, pos, CondAlways)
			}
			c.FocusWindow(c.MovingWindow)
		} else {
			c.ClearActiveID()
			c.MovingWindow = nil
		}
	} else {
		// When clicking/dragging from a window that has the _NoMove flag, we still set the ActiveId in order to prevent hovering others.
		if c.ActiveIdWindow != nil && c.ActiveIdWindow.MoveId == c.ActiveId {
			c.KeepAliveID(c.ActiveId)
			if !c.IO.MouseDown[0] {
				c.ClearActiveID()
			}
		}
	}
}

// Find window given position, search front-to-back
// FIXME: Note that we have a lag here because WindowRectClipped is updated in Begin() so windows moved by user via SetWindowPos() and not SetNextWindowPos() will have that rectangle lagging by a frame at the time FindHoveredWindow() is called, aka before the next Begin(). Moving window thankfully isn't affected.
func (c *Context) FindHoveredWindow() *Window {
	for i := len(c.Windows) - 1; i >= 0; i-- {
		window := c.Windows[i]
		if !window.Active {
			continue
		}
		if window.Flags&WindowFlagsNoInputs != 0 {
			continue
		}

		// Using the clipped AABB, a child window will typically be clipped by its parent (not always)
		bb := f64.Rectangle{
			window.WindowRectClipped.Min.Sub(c.Style.TouchExtraPadding),
			window.WindowRectClipped.Max.Add(c.Style.TouchExtraPadding),
		}
		if c.IO.MousePos.In(bb) {
			return window
		}
	}
	return nil
}

func (c *Context) IsWindowChildOf(window, potential_parent *Window) bool {
	if window.RootWindow == potential_parent {
		return true
	}
	for window != nil {
		if window == potential_parent {
			return true
		}
		window = window.ParentWindow
	}
	return false
}

func (c *Context) SetWindowScrollX(window *Window, new_scroll_x float64) {
	// SizeContents is generally computed based on CursorMaxPos which is affected by scroll position, so we need to apply our change to it.
	window.DC.CursorMaxPos.X += window.Scroll.X
	window.Scroll.X = new_scroll_x
	window.DC.CursorMaxPos.X -= window.Scroll.X
}

func (c *Context) SetWindowScrollY(window *Window, new_scroll_y float64) {
	// SizeContents is generally computed based on CursorMaxPos which is affected by scroll position, so we need to apply our change to it.
	window.DC.CursorMaxPos.Y += window.Scroll.Y
	window.Scroll.Y = new_scroll_y
	window.DC.CursorMaxPos.Y -= window.Scroll.Y
}

func (c *Context) FocusFrontMostActiveWindow(ignore_window *Window) {
	for i := len(c.Windows) - 1; i >= 0; i-- {
		if c.Windows[i] != ignore_window && c.Windows[i].WasActive && c.Windows[i].Flags&WindowFlagsChildWindow == 0 {
			focus_window := c.NavRestoreLastChildNavWindow(c.Windows[i])
			c.FocusWindow(focus_window)
			return
		}
	}
}

func (c *Context) SetNextWindowSize(size f64.Vec2, cond Cond) {
	// Make sure the user doesn't attempt to combine multiple condition flags.
	assert(cond == 0 || mathutil.IsPow2(int(cond)))
	c.NextWindowData.SizeVal = size
	c.NextWindowData.SizeCond = CondAlways
	if cond != 0 {
		c.NextWindowData.SizeCond = cond
	}
}

func (c *Context) FindWindowByName(name string) *Window {
	return c.WindowsById[name]
}

func (c *Context) SetWindowConditionAllowFlags(window *Window, flags Cond, enabled bool) {
	if enabled {
		window.SetWindowPosAllowFlags |= flags
		window.SetWindowSizeAllowFlags |= flags
		window.SetWindowCollapsedAllowFlags |= flags
	} else {
		window.SetWindowPosAllowFlags &^= flags
		window.SetWindowSizeAllowFlags &^= flags
		window.SetWindowCollapsedAllowFlags &^= flags
	}
}

func (c *Context) SetCurrentWindowPos(pos f64.Vec2, cond Cond) {
	window := c.GetCurrentWindowRead()
	c.SetWindowPos(window, pos, cond)
}

func (c *Context) SetWindowPosByName(name string, pos f64.Vec2, cond Cond) {
	window := c.FindWindowByName(name)
	if window != nil {
		c.SetWindowPos(window, pos, cond)
	}
}

func (c *Context) SetWindowPos(window *Window, pos f64.Vec2, cond Cond) {
	// Test condition (NB: bit 0 is always true) and clear flags for next time
	if cond != 0 && window.SetWindowPosAllowFlags&cond == 0 {
		return
	}

	// Make sure the user doesn't attempt to combine multiple condition flags.
	assert(cond == 0 || mathutil.IsPow2(int(cond)))
	window.SetWindowPosAllowFlags &^= (CondOnce | CondFirstUseEver | CondAppearing)
	window.SetWindowPosVal = f64.Vec2{math.MaxFloat32, math.MaxFloat32}

	// Set
	old_pos := window.Pos
	window.Pos = pos.Floor()
	// As we happen to move the window while it is being appended to (which is a bad idea - will smear) let's at least offset the cursor
	window.DC.CursorPos = window.DC.CursorPos.Add(window.Pos.Sub(old_pos))
	// And more importantly we need to adjust this so size calculation doesn't get affected.
	window.DC.CursorMaxPos = window.DC.CursorMaxPos.Add(window.Pos.Sub(old_pos))
}

func (c *Context) SetCurrentWindowSize(size f64.Vec2, cond Cond) {
	c.SetWindowSize(c.CurrentWindow, size, cond)
}

func (c *Context) SetWindowSizeByName(name string, size f64.Vec2, cond Cond) {
	window := c.FindWindowByName(name)
	if window != nil {
		c.SetWindowSize(window, size, cond)
	}
}

func (c *Context) SetWindowSize(window *Window, size f64.Vec2, cond Cond) {
	// Test condition (NB: bit 0 is always true) and clear flags for next time
	if cond != 0 && window.SetWindowSizeAllowFlags&cond == 0 {
		return
	}
	window.SetWindowSizeAllowFlags &^= (CondOnce | CondFirstUseEver | CondAppearing)

	// Set
	if size.X > 0 {
		window.AutoFitFramesX = 0
		window.SizeFull.X = size.X
	} else {
		window.AutoFitFramesX = 2
		window.AutoFitOnlyGrows = false
	}

	if size.Y > 0 {
		window.AutoFitFramesY = 0
		window.SizeFull.Y = size.Y
	} else {
		window.AutoFitFramesY = 2
		window.AutoFitOnlyGrows = false
	}
}

func (w *Window) TitleBarHeight() float64 {
	if w.Flags&WindowFlagsNoTitleBar != 0 {
		return 0
	}
	return w.CalcFontSize() + w.Ctx.Style.FramePadding.Y*2.0
}

func (w *Window) MenuBarHeight() float64 {
	if w.Flags&WindowFlagsMenuBar != 0 {
		return w.DC.MenuBarOffset.Y + w.CalcFontSize() + w.Ctx.Style.FramePadding.Y*2.0
	}
	return 0
}

func (c *Context) SetWindowFocus() {
	c.FocusWindow(c.CurrentWindow)
}

func (c *Context) SetWindowCollapsed(window *Window, collapsed bool, cond Cond) {
	// Test condition (NB: bit 0 is always true) and clear flags for next time
	if cond != 0 && (window.SetWindowCollapsedAllowFlags&cond) == 0 {
		return
	}
	window.SetWindowCollapsedAllowFlags &^= (CondOnce | CondFirstUseEver | CondAppearing)

	// Set
	window.Collapsed = collapsed
}

func (w *Window) TitleBarRect() f64.Rectangle {
	return f64.Rectangle{
		w.Pos,
		f64.Vec2{w.Pos.X + w.SizeFull.X, w.Pos.Y + w.TitleBarHeight()},
	}
}

func (c *Context) CalcSizeAutoFit(window *Window, size_contents f64.Vec2) f64.Vec2 {
	var size_auto_fit f64.Vec2
	style := &c.Style

	if window.Flags&WindowFlagsTooltip != 0 {
		// Tooltip always resize. We keep the spacing symmetric on both axises for aesthetic purpose.
		return size_contents
	} else {
		// When the window cannot fit all contents (either because of constraints, either because screen is too small): we are growing the size on the other axis to compensate for expected scrollbar. FIXME: Might turn bigger than DisplaySize-WindowPadding.
		size_auto_fit = f64.Vec2{
			f64.Clamp(size_contents.X, style.WindowMinSize.X, math.Max(style.WindowMinSize.X, c.IO.DisplaySize.X-c.Style.DisplaySafeAreaPadding.X*2)),
			f64.Clamp(size_contents.Y, style.WindowMinSize.Y, math.Max(style.WindowMinSize.Y, c.IO.DisplaySize.Y-c.Style.DisplaySafeAreaPadding.Y*2)),
		}

		size_auto_fit_after_constraint := c.CalcSizeAfterConstraint(window, size_auto_fit)
		if size_auto_fit_after_constraint.X < size_contents.X && window.Flags&WindowFlagsNoScrollbar == 0 && window.Flags&WindowFlagsHorizontalScrollbar != 0 {
			size_auto_fit.Y += style.ScrollbarSize
		}
		if size_auto_fit_after_constraint.Y < size_contents.Y && window.Flags&WindowFlagsNoScrollbar == 0 {
			size_auto_fit.X += style.ScrollbarSize
		}
		return size_auto_fit
	}
}

func (c *Context) CalcSizeAfterConstraint(window *Window, new_size f64.Vec2) f64.Vec2 {
	if c.NextWindowData.SizeConstraintCond != 0 {
		// Using -1,-1 on either X/Y axis to preserve the current size.
		cr := c.NextWindowData.SizeConstraintRect
		if cr.Min.X >= 0 && cr.Max.X >= 0 {
			new_size.X = f64.Clamp(new_size.X, cr.Min.X, cr.Max.X)
		} else {
			new_size.X = window.SizeFull.X
		}

		if cr.Min.Y >= 0 && cr.Max.Y >= 0 {
			new_size.Y = f64.Clamp(new_size.Y, cr.Min.Y, cr.Max.Y)
		} else {
			new_size.Y = window.SizeFull.Y
		}

		if c.NextWindowData.SizeCallback != nil {
			data := SizeCallbackData{
				Pos:         window.Pos,
				CurrentSize: window.SizeFull,
				DesiredSize: new_size,
			}
			c.NextWindowData.SizeCallback(&data)
			new_size = data.DesiredSize
		}
	}

	// Minimum size
	if window.Flags&(WindowFlagsChildWindow|WindowFlagsAlwaysAutoResize) == 0 {
		new_size = new_size.Max(c.Style.WindowMinSize)

		// Reduce artifacts with very small windows
		new_size.Y = math.Max(new_size.Y, window.TitleBarHeight()+window.MenuBarHeight()+math.Max(0, c.Style.WindowRounding-1))
	}

	return new_size
}

func (c *Context) FindBestWindowPosForPopup(window *Window) f64.Vec2 {
	r_outer := c.FindAllowedExtentRectForWindow(window)
	if window.Flags&WindowFlagsChildMenu != 0 {
		// Child menus typically request _any_ position within the parent menu item, and then our FindBestWindowPosForPopup() function will move the new menu outside the parent bounds.
		// This is how we end up with child menus appearing (most-commonly) on the right of the parent menu.
		assert(c.CurrentWindow == window)
		parent_window := c.CurrentWindowStack[len(c.CurrentWindowStack)-2]
		// We want some overlap to convey the relative depth of each menu (currently the amount of overlap is hard-coded to style.ItemSpacing.x).
		horizontal_overlap := c.Style.ItemSpacing.X

		var r_avoid f64.Rectangle
		if parent_window.DC.MenuBarAppending {
			r_avoid = f64.Rect(
				-math.MaxFloat32,
				parent_window.Pos.Y+parent_window.TitleBarHeight(),
				math.MaxFloat32,
				parent_window.Pos.Y+parent_window.TitleBarHeight()+parent_window.MenuBarHeight(),
			)
		} else {
			r_avoid = f64.Rect(
				parent_window.Pos.X+horizontal_overlap,
				-math.MaxFloat32,
				parent_window.Pos.X+parent_window.Size.X-horizontal_overlap-parent_window.ScrollbarSizes.X,
				math.MaxFloat32,
			)
		}
		return c.FindBestWindowPosForPopupEx(window.Pos, window.Size, &window.AutoPosLastDirection, r_outer, r_avoid, PopupPositionPolicyDefault)
	}

	if window.Flags&WindowFlagsPopup != 0 {
		r_avoid := f64.Rect(window.Pos.X-1, window.Pos.Y-1, window.Pos.X+1, window.Pos.Y+1)
		return c.FindBestWindowPosForPopupEx(window.Pos, window.Size, &window.AutoPosLastDirection, r_outer, r_avoid, PopupPositionPolicyDefault)
	}

	if window.Flags&WindowFlagsTooltip != 0 {
		// Position tooltip (always follows mouse)
		sc := c.Style.MouseCursorScale
		ref_pos := c.NavCalcPreferredRefPos()
		var r_avoid f64.Rectangle
		if !c.NavDisableHighlight && c.NavDisableMouseHover && c.IO.ConfigFlags&ConfigFlagsNavEnableSetMousePos == 0 {
			r_avoid = f64.Rect(ref_pos.X-16, ref_pos.Y-8, ref_pos.X+16, ref_pos.Y+8)
		} else {
			// FIXME: Hard-coded based on mouse cursor shape expectation. Exact dimension not very important.
			r_avoid = f64.Rect(ref_pos.X-16, ref_pos.Y-8, ref_pos.X+24*sc, ref_pos.Y+24*sc)
		}
		pos := c.FindBestWindowPosForPopupEx(ref_pos, window.Size, &window.AutoPosLastDirection, r_outer, r_avoid, PopupPositionPolicyDefault)
		if window.AutoPosLastDirection == DirNone {
			// If there's not enough room, for tooltip we prefer avoiding the cursor at all cost even if it means that part of the tooltip won't be visible.
			pos = ref_pos.Add(f64.Vec2{2, 2})
		}
		return pos
	}

	assert(false)
	return window.Pos
}

// r_avoid = the rectangle to avoid (e.g. for tooltip it is a rectangle around the mouse cursor which we want to avoid. for popups it's a small point around the cursor.)
// r_outer = the visible area rectangle, minus safe area padding. If our popup size won't fit because of safe area padding we ignore it.
func (c *Context) FindBestWindowPosForPopupEx(ref_pos, size f64.Vec2, last_dir *Dir, r_outer, r_avoid f64.Rectangle, policy PopupPositionPolicy) f64.Vec2 {
	base_pos_clamped := ref_pos.Clamp2(r_outer.Min, r_outer.Max.Sub(size))

	// Combo Box policy (we want a connecting edge)
	if policy == PopupPositionPolicyComboBox {
		dir_prefered_order := [...]Dir{DirDown, DirRight, DirLeft, DirUp}
		n := Dir(0)
		if *last_dir != DirNone {
			n = -1
		}
		for ; n < DirCOUNT; n++ {
			var dir Dir
			if n == -1 {
				dir = *last_dir
			} else {
				dir = dir_prefered_order[n]
			}

			// Already tried this direction?
			if n != -1 && dir == *last_dir {
				continue
			}

			var pos f64.Vec2
			// Below, Toward Right (default)
			if dir == DirDown {
				pos = f64.Vec2{r_avoid.Min.X, r_avoid.Max.Y}
			}
			// Above, Toward Right
			if dir == DirRight {
				pos = f64.Vec2{r_avoid.Min.X, r_avoid.Min.Y - size.Y}
			}
			// Below, Toward Left
			if dir == DirLeft {
				pos = f64.Vec2{r_avoid.Max.X - size.X, r_avoid.Max.Y}
			}
			// Above, Toward Left
			if dir == DirUp {
				pos = f64.Vec2{r_avoid.Max.X - size.X, r_avoid.Min.Y - size.Y}
			}

			r_test := f64.Rectangle{pos, pos.Add(size)}
			if !r_test.In(r_outer) {
				continue
			}

			*last_dir = dir
			return pos
		}
	}

	// Default popup policy
	dir_prefered_order := [...]Dir{DirRight, DirDown, DirUp, DirLeft}
	n := Dir(0)
	if *last_dir != DirNone {
		n = -1
	}
	for ; n < DirCOUNT; n++ {
		var dir Dir
		if n == -1 {
			dir = *last_dir
		} else {
			dir = dir_prefered_order[n]
		}

		// Already tried this direction?
		if n != -1 && dir == *last_dir {
			continue
		}

		var avail_w, avail_h float64
		var x1, y1, x2, y2 float64
		if dir == DirLeft {
			x1 = r_avoid.Min.X
		} else {
			x1 = r_outer.Max.X
		}
		if dir == DirRight {
			x2 = r_avoid.Max.X
		} else {
			x2 = r_outer.Min.X
		}
		if dir == DirUp {
			y1 = r_avoid.Min.Y
		} else {
			y1 = r_outer.Max.Y
		}
		if dir == DirDown {
			y2 = r_avoid.Max.Y
		} else {
			y2 = r_outer.Min.Y
		}

		avail_w = x1 - x2
		avail_h = y1 - y2
		if avail_w < size.X || avail_h < size.Y {
			continue
		}

		var pos f64.Vec2
		switch dir {
		case DirLeft:
			pos.X = r_avoid.Min.X - size.X
		case DirRight:
			pos.X = r_avoid.Max.X
		default:
			pos.X = base_pos_clamped.X
		}

		switch dir {
		case DirUp:
			pos.Y = r_avoid.Min.Y - size.Y
		case DirDown:
			pos.Y = r_avoid.Max.Y
		default:
			pos.Y = base_pos_clamped.Y
		}

		*last_dir = dir
		return pos
	}

	// Fallback, try to keep within display
	*last_dir = DirNone
	pos := ref_pos
	pos.X = math.Max(math.Min(pos.X+size.X, r_outer.Max.X)-size.X, r_outer.Min.X)
	pos.Y = math.Max(math.Min(pos.Y+size.Y, r_outer.Max.Y)-size.Y, r_outer.Min.Y)
	return pos
}

func (w *Window) Rect() f64.Rectangle {
	return f64.Rectangle{
		w.Pos,
		w.Pos.Add(w.Size),
	}
}

func (w *Window) GetWindowBgColorIdxFromFlags(flags WindowFlags) Col {
	if flags&(WindowFlagsTooltip|WindowFlagsPopup) != 0 {
		return ColPopupBg
	}
	if flags&WindowFlagsChildWindow != 0 {
		return ColChildBg
	}
	return ColWindowBg
}

func (w *Window) MenuBarRect() f64.Rectangle {
	y1 := w.Pos.Y + w.TitleBarHeight()
	return f64.Rect(w.Pos.X, y1, w.Pos.X+w.SizeFull.X, y1+w.MenuBarHeight())
}

func (c *Context) FindWindowSettings(name string) *WindowSettings {
	return c.SettingsWindows[name]
}

func (c *Context) AddWindowSettings(name string) *WindowSettings {
	settings := &WindowSettings{}
	settings.Init()
	settings.Name = name
	c.SettingsWindows[name] = settings
	return settings
}

func (s *WindowSettings) Init() {
	s.Name = ""
	s.Id = 0
	s.Pos = f64.Vec2{0, 0}
	s.Size = f64.Vec2{0, 0}
	s.Collapsed = false
}

func (c *Context) FocusableItemRegister(window *Window, id ID) bool {
	return c.FocusableItemRegisterEx(window, id, true)
}

func (c *Context) FocusableItemRegisterEx(window *Window, id ID, tab_stop bool) bool {
	allow_keyboard_focus := window.DC.ItemFlags&(ItemFlagsAllowKeyboardFocus|ItemFlagsDisabled) == ItemFlagsAllowKeyboardFocus
	window.FocusIdxAllCounter++
	if allow_keyboard_focus {
		window.FocusIdxTabCounter++
	}

	// Process keyboard input at this point: TAB/Shift-TAB to tab out of the currently focused item.
	// Note that we can always TAB out of a widget that doesn't allow tabbing in.
	if tab_stop && (c.ActiveId == id) && window.FocusIdxAllRequestNext == math.MaxInt32 && window.FocusIdxTabRequestNext == math.MaxInt32 && !c.IO.KeyCtrl && c.IsKeyPressedMap(KeyTab) {
		// Modulo on index will be applied at the end of frame once we've got the total counter of items.
		shift := 1
		if allow_keyboard_focus {
			shift = -1
		} else {
			shift = 0
		}
		window.FocusIdxTabRequestNext = window.FocusIdxTabCounter + shift
	}

	if window.FocusIdxAllCounter == window.FocusIdxAllRequestCurrent {
		return true
	}

	if allow_keyboard_focus && window.FocusIdxTabCounter == window.FocusIdxTabRequestCurrent {
		c.NavJustTabbedId = id
		return true
	}

	return false
}

func (c *Context) FocusableItemUnregister(window *Window) {
	window.FocusIdxAllCounter--
	window.FocusIdxTabCounter--
}

func (c *Context) SetNextWindowPos(pos f64.Vec2, cond Cond, pivot f64.Vec2) {
	c.NextWindowData.PosVal = pos
	c.NextWindowData.PosPivotVal = pivot
	c.NextWindowData.PosCond = CondAlways
	if cond != 0 {
		c.NextWindowData.PosCond = cond
	}
}

func (n *NextWindowData) Init() {
	n.PosCond = 0
	n.SizeCond = 0
	n.ContentSizeCond = 0
	n.CollapsedCond = 0
	n.SizeConstraintCond = 0
	n.FocusCond = 0
	n.BgAlphaCond = 0
	n.PosVal = f64.Vec2{0, 0}
	n.PosPivotVal = f64.Vec2{0, 0}
	n.SizeVal = f64.Vec2{0, 0}
	n.ContentSizeVal = f64.Vec2{0, 0}
	n.CollapsedVal = false
	n.SizeConstraintRect = f64.Rectangle{}
	n.SizeCallback = nil
	n.BgAlphaVal = math.MaxFloat32
	n.MenuBarOffsetMinVal = f64.Vec2{0, 0}
}

func (n *NextWindowData) Clear() {
	n.PosCond = 0
	n.SizeCond = 0
	n.ContentSizeCond = 0
	n.CollapsedCond = 0
	n.SizeConstraintCond = 0
	n.FocusCond = 0
	n.BgAlphaCond = 0
}

func (c *Context) Spacing() {
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return
	}
	c.ItemSize(f64.Vec2{0, 0})
}

func (c *Context) SetStateStorage(tree map[ID]interface{}) {
	window := c.GetCurrentWindow()
	if tree == nil {
		tree = window.StateStorage
	}
	window.DC.StateStorage = tree
}

func (c *Context) GetStateStorage() map[ID]interface{} {
	window := c.GetCurrentWindowRead()
	return window.DC.StateStorage
}

func (c *Context) GetWindowSize() f64.Vec2 {
	window := c.GetCurrentWindowRead()
	return window.Size
}

func (c *Context) SetNextWindowSizeConstraints(size_min, size_max f64.Vec2, custom_callback func(*SizeCallbackData)) {
	c.NextWindowData.SizeConstraintCond = CondAlways
	c.NextWindowData.SizeConstraintRect = f64.Rectangle{size_min, size_max}
	c.NextWindowData.SizeCallback = custom_callback
}

func (c *Context) GetCursorScreenPos() f64.Vec2 {
	window := c.GetCurrentWindowRead()
	return window.DC.CursorPos
}

func (c *Context) GetWindowDrawList() *DrawList {
	window := c.GetCurrentWindow()
	return window.DrawList
}

func (c *Context) Dummy(size f64.Vec2) {
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return
	}

	bb := f64.Rectangle{window.DC.CursorPos, window.DC.CursorPos.Add(size)}
	c.ItemSizeBB(bb)
	c.ItemAdd(bb, 0)
}

func (c *Context) FindAllowedExtentRectForWindow(window *Window) f64.Rectangle {
	padding := c.Style.DisplaySafeAreaPadding
	r_screen := c.GetViewportRect()
	expand := f64.Vec2{0, 0}
	if r_screen.Dx() > padding.X*2 {
		expand.X = -padding.X
	}
	if r_screen.Dy() > padding.Y*2 {
		expand.Y = -padding.Y
	}
	return r_screen.Expand2(expand)
}

func (w *Window) GetIDNoKeepAlive(str string) ID {
	seed := w.IDStack[len(w.IDStack)-1]
	h := fnv.New32()
	binary.Write(h, binary.LittleEndian, uint32(seed))
	h.Write([]byte(str))
	return ID(h.Sum32())
}

// This is only used in rare/specific situations to manufacture an ID out of nowhere.
func (w *Window) GetIDFromRectangle(r_abs f64.Rectangle) ID {
	seed := w.IDStack[len(w.IDStack)-1]
	h := fnv.New32()
	r_rel := [4]int{
		int(r_abs.Min.X - w.Pos.X),
		int(r_abs.Min.Y - w.Pos.Y),
		int(r_abs.Max.X - w.Pos.X),
		int(r_abs.Max.Y - w.Pos.Y),
	}
	binary.Write(h, binary.LittleEndian, seed)
	for i := range r_rel {
		binary.Write(h, binary.LittleEndian, uint32(r_rel[i]))
	}
	return ID(h.Sum32())
}

func (c *Context) SetCursorScreenPos(screen_pos f64.Vec2) {
	window := c.GetCurrentWindow()
	window.DC.CursorPos = screen_pos
	window.DC.CursorMaxPos = window.DC.CursorMaxPos.Max(window.DC.CursorPos)
}

func (c *Context) GetWindowWidth() float64 {
	window := c.CurrentWindow
	return window.Size.X
}

func (c *Context) GetWindowHeight() float64 {
	window := c.CurrentWindow
	return window.Size.Y
}

// User generally sees positions in window coordinates. Internally we store CursorPos in absolute screen coordinates because it is more convenient.
// Conversion happens as we pass the value to user, but it makes our naming convention confusing because GetCursorPos() == (DC.CursorPos - window.Pos). May want to rename 'DC.CursorPos'.
func (c *Context) GetCursorPos() f64.Vec2 {
	window := c.GetCurrentWindowRead()
	return window.DC.CursorPos.Sub(window.Pos).Add(window.Scroll)

}

func (c *Context) GetItemRectMin() f64.Vec2 {
	window := c.GetCurrentWindowRead()
	return window.DC.LastItemRect.Min
}

func (c *Context) GetItemRectMax() f64.Vec2 {
	window := c.GetCurrentWindowRead()
	return window.DC.LastItemRect.Max
}

func (c *Context) SetNextWindowBgAlpha(alpha float64) {
	// Using a Cond member for consistency (may transition all of them to single flag set for fast Clear() op)
	c.NextWindowData.BgAlphaVal = alpha
	c.NextWindowData.BgAlphaCond = CondAlways
}

func (c *Context) IsWindowHovered() bool {
	return c.IsWindowHoveredEx(0)
}

func (c *Context) IsWindowHoveredEx(flags HoveredFlags) bool {
	// Flags not supported by this function
	assert(flags&HoveredFlagsAllowWhenOverlapped == 0)
	if flags&HoveredFlagsAnyWindow != 0 {
		if c.HoveredWindow == nil {
			return false
		}
	} else {
		switch flags & (HoveredFlagsRootWindow | HoveredFlagsChildWindows) {
		case HoveredFlagsRootWindow | HoveredFlagsChildWindows:
			if c.HoveredRootWindow != c.CurrentWindow.RootWindow {
				return false
			}
		case HoveredFlagsRootWindow:
			if c.HoveredWindow != c.CurrentWindow.RootWindow {
				return false
			}
		case HoveredFlagsChildWindows:
			if c.HoveredWindow == nil || !c.IsWindowChildOf(c.HoveredWindow, c.CurrentWindow) {
				return false
			}
		default:
			if c.HoveredWindow != c.CurrentWindow {
				return false
			}
		}
	}

	if !c.IsWindowContentHoverable(c.HoveredRootWindow, flags) {
		return false
	}

	if flags&HoveredFlagsAllowWhenBlockedByActiveItem == 0 {
		if c.ActiveId != 0 && !c.ActiveIdAllowOverlap && c.ActiveId != c.HoveredWindow.MoveId {
			return false
		}
	}

	return true
}

func (c *Context) GetCursorPosY() float64 {
	window := c.GetCurrentWindowRead()
	return window.DC.CursorPos.Y - window.Pos.Y + window.Scroll.Y
}

func (c *Context) SetCursorPosY(y float64) {
	window := c.GetCurrentWindow()
	window.DC.CursorPos.Y = window.Pos.Y - window.Scroll.Y + y
	window.DC.CursorMaxPos.Y = math.Max(window.DC.CursorMaxPos.Y, window.DC.CursorPos.Y)
}

func (c *Context) GetItemRectSize() f64.Vec2 {
	window := c.GetCurrentWindowRead()
	return window.DC.LastItemRect.Size()
}

func (c *Context) GetCursorStartPos() f64.Vec2 {
	window := c.GetCurrentWindowRead()
	return window.DC.CursorStartPos.Sub(window.Pos)
}

func (c *Context) GetScrollX() float64 {
	return c.CurrentWindow.Scroll.X
}

func (c *Context) GetScrollY() float64 {
	return c.CurrentWindow.Scroll.Y
}

func (c *Context) SetScrollX(scroll_x float64) {
	window := c.GetCurrentWindow()
	window.ScrollTarget.X = scroll_x
	window.ScrollTargetCenterRatio.X = 0.0
}

func (c *Context) SetScrollY(scroll_y float64) {
	window := c.GetCurrentWindow()
	// title bar height canceled out when using ScrollTargetRelY
	window.ScrollTarget.Y = scroll_y + window.TitleBarHeight() + window.MenuBarHeight()
	window.ScrollTargetCenterRatio.Y = 0.0
}

func (c *Context) SetNextWindowContentSize(size f64.Vec2) {
	// In Begin() we will add the size of window decorations (title bar, menu etc.) to that to form a SizeContents value.
	c.NextWindowData.ContentSizeVal = size
	c.NextWindowData.ContentSizeCond = CondAlways
}

func (c *Context) SetKeyboardFocusHere() {
	c.SetKeyboardFocusHereEx(-1)
}

func (c *Context) SetKeyboardFocusHereEx(offset int) {
	// -1 is allowed but not below
	assert(offset >= -1)
	window := c.GetCurrentWindow()
	window.FocusIdxAllRequestNext = window.FocusIdxAllCounter + 1 + offset
	window.FocusIdxTabRequestNext = math.MaxInt32
}

func (c *Context) IsWindowFocused() bool {
	return c.IsWindowFocusedEx(0)
}

func (c *Context) IsWindowFocusedEx(flags FocusedFlags) bool {
	// Not inside a Begin()/End()
	assert(c.CurrentWindow != nil)

	if flags&FocusedFlagsAnyWindow != 0 {
		return c.NavWindow != nil
	}

	switch flags & (FocusedFlagsRootWindow | FocusedFlagsChildWindows) {
	case FocusedFlagsRootWindow | FocusedFlagsChildWindows:
		return c.NavWindow != nil && c.NavWindow.RootWindow == c.CurrentWindow.RootWindow
	case FocusedFlagsRootWindow:
		return c.NavWindow == c.CurrentWindow.RootWindow
	case FocusedFlagsChildWindows:
		return c.NavWindow != nil && c.IsWindowChildOf(c.NavWindow, c.CurrentWindow)
	default:
		return c.NavWindow == c.CurrentWindow
	}
}

type FocusedFlags uint

const (
	FocusedFlagsChildWindows        FocusedFlags = 1 << 0 // IsWindowFocused(): Return true if any children of the window is focused
	FocusedFlagsRootWindow          FocusedFlags = 1 << 1 // IsWindowFocused(): Test from root window (top most parent of the current hierarchy)
	FocusedFlagsAnyWindow           FocusedFlags = 1 << 2 // IsWindowFocused(): Return true if any window is focused
	FocusedFlagsRootAndChildWindows FocusedFlags = FocusedFlagsRootWindow | FocusedFlagsChildWindows
)
