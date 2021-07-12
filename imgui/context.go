package imgui

import (
	"math"
	"os"

	"github.com/qeedquan/go-media/math/f64"
)

type (
	ID        uint
	TextureID interface{}
)

type Context struct {
	Initialized             bool
	FontAtlasOwnedByContext bool // Io.Fonts-> is owned by the ImGuiContext and will be destructed along with it.
	IO                      IO
	Style                   Style
	Font                    *Font   // (Shortcut) == FontStack.empty() ? IO.Font : FontStack.back()
	FontSize                float64 // (Shortcut) == FontBaseSize * g.CurrentWindow->FontWindowScale == window->FontSize(). Text height for current window.
	FontBaseSize            float64 // (Shortcut) == IO.FontGlobalScale * Font->Scale * Font->FontSize. Base text height.
	DrawListSharedData      DrawListSharedData

	Time                     float64
	FrameCount               int
	FrameCountEnded          int
	FrameCountRendered       int
	Windows                  []*Window
	WindowsSortBuffer        []*Window
	CurrentWindowStack       []*Window
	WindowsById              map[string]*Window
	WindowsActiveCount       int
	CurrentWindow            *Window // Being drawn into
	HoveredWindow            *Window // Will catch mouse inputs
	HoveredRootWindow        *Window // Will catch mouse inputs (for focus/move only)
	HoveredId                ID      // Hovered widget
	HoveredIdAllowOverlap    bool
	HoveredIdPreviousFrame   ID
	HoveredIdTimer           float64
	ActiveId                 ID // Active widget
	ActiveIdPreviousFrame    ID
	ActiveIdTimer            float64
	ActiveIdIsAlive          bool     // Active widget has been seen this frame
	ActiveIdIsJustActivated  bool     // Set at the time of activation for one frame
	ActiveIdAllowOverlap     bool     // Active widget allows another widget to steal active id (generally for overlapping widgets, but not always)
	ActiveIdAllowNavDirFlags int      // Active widget allows using directional navigation (e.g. can activate a button and move away from it)
	ActiveIdClickOffset      f64.Vec2 // Clicked offset from upper-left corner, if applicable (currently only set by ButtonBehavior)
	ActiveIdWindow           *Window
	ActiveIdSource           InputSource    // Activating with mouse or nav (gamepad/keyboard)
	MovingWindow             *Window        // Track the window we clicked on (in order to preserve focus). The actually window that is moved is generally MovingWindow->RootWindow.
	ColorModifiers           []ColMod       // Stack for PushStyleColor()/PopStyleColor()
	StyleModifiers           []StyleMod     // Stack for PushStyleVar()/PopStyleVar()
	FontStack                []*Font        // Stack for PushFont()/PopFont()
	OpenPopupStack           []PopupRef     // Which popups are open (persistent)
	CurrentPopupStack        []PopupRef     // Which level of BeginPopup() we are in (reset every frame)
	NextWindowData           NextWindowData // Storage for SetNextWindow** functions
	NextTreeNodeOpenVal      bool           // Storage for SetNextTreeNode** functions
	NextTreeNodeOpenCond     Cond

	// Navigation data (for gamepad/keyboard)
	NavWindow                  *Window       // Focused window for navigation. Could be called 'FocusWindow'
	NavId                      ID            // Focused item for navigation
	NavActivateId              ID            // ~~ (g.ActiveId == 0) && IsNavInputPressed(ImGuiNavInput_Activate) ? NavId : 0, also set when calling ActivateItem()
	NavActivateDownId          ID            // ~~ IsNavInputDown(ImGuiNavInput_Activate) ? NavId : 0
	NavActivatePressedId       ID            // ~~ IsNavInputPressed(ImGuiNavInput_Activate) ? NavId : 0
	NavInputId                 ID            // ~~ IsNavInputPressed(ImGuiNavInput_Input) ? NavId : 0
	NavJustTabbedId            ID            // Just tabbed to this id.
	NavNextActivateId          ID            // Set by ActivateItem(), queued until next frame
	NavJustMovedToId           ID            // Just navigated to this id (result of a successfully MoveRequest)
	NavInputSource             InputSource   // Keyboard or Gamepad mode?
	NavScoringRectScreen       f64.Rectangle // Rectangle used for scoring, in screen space. Based of window->DC.NavRefRectRel[], modified for directional navigation scoring.
	NavScoringCount            int           // Metrics for debugging
	NavWindowingTarget         *Window       // When selecting a window (holding Menu+FocusPrev/Next, or equivalent of CTRL-TAB) this window is temporarily displayed front-most.
	NavWindowingHighlightTimer float64
	NavWindowingHighlightAlpha float64
	NavWindowingToggleLayer    bool
	NavWindowingInputSource    InputSource // Gamepad or keyboard mode
	NavLayer                   int         // Layer we are navigating on. For now the system is hard-coded for 0=main contents and 1=menu/title bar, may expose layers later.
	NavIdTabCounter            int         // == NavWindow->DC.FocusIdxTabCounter at time of NavId processing
	NavIdIsAlive               bool        // Nav widget has been seen this frame ~~ NavRefRectRel is valid
	NavMousePosDirty           bool        // When set we will update mouse position if (io.ConfigFlags & ImGuiConfigFlags_NavMoveMouse) if set (NB: this not enabled by default)
	NavDisableHighlight        bool        // When user starts using mouse, we hide gamepad/keyboard highlight (NB: but they are still available, which is why NavDisableHighlight isn't always != NavDisableMouseHover)
	NavDisableMouseHover       bool        // When user starts using gamepad/keyboard, we hide mouse hovering highlight until mouse is touched again.
	NavAnyRequest              bool        // ~~ NavMoveRequest || NavInitRequest
	NavInitRequest             bool        // Init request for appearing window to select first item
	NavInitRequestFromMove     bool
	NavInitResultId            ID
	NavInitResultRectRel       f64.Rectangle
	NavMoveFromClampedRefRect  bool          // Set by manual scrolling, if we scroll to a point where NavId isn't visible we reset navigation from visible items
	NavMoveRequest             bool          // Move request for this frame
	NavMoveRequestForward      NavForward    // None / ForwardQueued / ForwardActive (this is used to navigate sibling parent menus from a child menu)
	NavMoveDir, NavMoveDirLast Dir           // Direction of the move request (left/right/up/down), direction of the previous move request
	NavMoveResultLocal         NavMoveResult // Best move request candidate within NavWindow
	NavMoveResultOther         NavMoveResult // Best move request candidate within NavWindow's flattened hierarchy (when using the NavFlattened flag)

	// Render
	DrawData                  DrawData // Main ImDrawData instance to pass render information to the user
	DrawDataBuilder           DrawDataBuilder
	ModalWindowDarkeningRatio float64
	OverlayDrawList           DrawList // Optional software render of mouse cursors, if io.MouseDrawCursor is set + a few debug overlays
	MouseCursor               MouseCursor

	// Drag and Drop
	DragDropActive                  bool
	DragDropSourceFlags             DragDropFlags
	DragDropMouseButton             int
	DragDropPayload                 Payload
	DragDropTargetRect              f64.Rectangle
	DragDropTargetId                ID
	DragDropAcceptIdCurrRectSurface float64
	DragDropAcceptIdCurr            ID      // Target item id (set at the time of accepting the payload)
	DragDropAcceptIdPrev            ID      // Target item id from previous frame (we need to store this to allow for overlapping drag and drop targets)
	DragDropAcceptFrameCount        int     // Last time a target expressed a desire to accept the source
	DragDropPayloadBufHeap          []uint8 // We don't expose the ImVector<> directly
	DragDropPayloadBufLocal         [8]uint8

	// Widget state
	InputTextState                     TextEditState
	InputTextPasswordFont              Font
	ScalarAsInputTextId                ID             // Temporary text input when CTRL+clicking on a slider, etc.
	ColorEditOptions                   ColorEditFlags // Store user options for color edit widgets
	ColorPickerRef                     f64.Vec4
	DragCurrentValue                   float64 // Currently dragged value, always float, not rounded by end-user precision settings
	DragLastMouseDelta                 f64.Vec2
	DragSpeedDefaultRatio              float64 // If speed == 0.0f, uses (max-min) * DragSpeedDefaultRatio
	DragSpeedScaleSlow                 float64
	DragSpeedScaleFast                 float64
	ScrollbarClickDeltaToGrabCenter    f64.Vec2 // Distance between mouse and center of grab box, normalized in parent space. Use storage?
	TooltipOverrideCount               int
	PrivateClipboard                   string   // If no custom clipboard handler is defined
	PlatformImePos, PlatformImeLastPos f64.Vec2 // Cursor position request & last passed to the OS Input Method Editor

	// Settings
	SettingsLoaded     bool
	SettingsDirtyTimer float64                     // Save .ini Settings on disk when time reaches zero
	SettingsWindows    map[string]*WindowSettings  // .ini settings for ImGuiWindow
	SettingsHandlers   map[string]*SettingsHandler // List of .ini settings handlers

	// Logging
	LogEnabled            bool
	LogFile               *os.File // If != NULL log to stdout/ file
	LogClipboard          []rune   // Else log to clipboard. This is pointer so our GImGui static constructor doesn't call heap allocators.
	LogStartDepth         int
	LogAutoExpandMaxDepth int

	// Misc
	FramerateSecPerFrame         [120]float64 // Calculate estimate of framerate for user over the last 2 seconds.
	FramerateSecPerFrameIdx      int
	FramerateSecPerFrameAccum    float64
	WantCaptureMouseNextFrame    int // Explicit capture via CaptureKeyboardFromApp()/CaptureMouseFromApp() sets those flags
	WantCaptureKeyboardNextFrame int
	WantTextInputNextFrame       int
}

type ConfigFlags uint

const (
	ConfigFlagsNavEnableKeyboard    ConfigFlags = 1 << 0 // Master keyboard navigation enable flag. NewFrame() will automatically fill io.NavInputs[] based on io.KeysDown[].
	ConfigFlagsNavEnableGamepad     ConfigFlags = 1 << 1 // Master gamepad navigation enable flag. This is mostly to instruct your imgui back-end to fill io.NavInputs[].
	ConfigFlagsNavEnableSetMousePos ConfigFlags = 1 << 2 // Instruct navigation to move the mouse cursor. May be useful on TV/console systems where moving a virtual mouse is awkward. Will update io.MousePos and set io.WantSetMousePos=true. If enabled you MUST honor io.WantSetMousePos requests in your binding, otherwise ImGui will react as if the mouse is jumping around back and forth.
	ConfigFlagsNavNoCaptureKeyboard ConfigFlags = 1 << 3 // Instruct navigation to not set the io.WantCaptureKeyboard flag with io.NavActive is set.
	ConfigFlagsNoMouse              ConfigFlags = 1 << 4 // Instruct imgui to clear mouse position/buttons in NewFrame(). This allows ignoring the mouse information back-end
	ConfigFlagsNoMouseCursorChange  ConfigFlags = 1 << 5 // Instruct back-end to not alter mouse cursor shape and visibility.

	// User storage (to allow your back-end/engine to communicate to code that may be shared between multiple projects. Those flags are not used by core ImGui)
	ConfigFlagsIsSRGB        ConfigFlags = 1 << 20 // Back-end is SRGB-aware.
	ConfigFlagsIsTouchScreen ConfigFlags = 1 << 21 // Back-end is using a touch screen instead of a mouse.
)

func CreateContext() *Context {
	return CreateContextEx(nil)
}

func CreateContextEx(shared_font_atlas *FontAtlas) *Context {
	c := &Context{}
	c.Init(shared_font_atlas)
	return c
}

func (c *Context) Init(shared_font_atlas *FontAtlas) {
	c.DrawListSharedData.Init()
	c.OverlayDrawList.Init(nil)
	c.Style.Init()
	c.StyleColorsDark(nil)
	c.IO.Init(c)
	io := c.GetIO()
	c.Font = nil
	c.FontSize = 0
	c.FontBaseSize = 0
	c.FontAtlasOwnedByContext = false
	io.Fonts = nil
	if shared_font_atlas == nil {
		c.FontAtlasOwnedByContext = true
		io.Fonts = NewFontAtlas()
	}
	c.Time = 0
	c.FrameCount = 0
	c.FrameCountEnded = -1
	c.FrameCountRendered = -1
	c.WindowsActiveCount = 0
	c.CurrentWindow = nil
	c.HoveredWindow = nil
	c.HoveredRootWindow = nil
	c.HoveredId = 0
	c.HoveredIdAllowOverlap = false
	c.HoveredIdPreviousFrame = 0
	c.HoveredIdTimer = 0.0
	c.ActiveId = 0
	c.ActiveIdPreviousFrame = 0
	c.ActiveIdTimer = 0.0
	c.ActiveIdIsAlive = false
	c.ActiveIdIsJustActivated = false
	c.ActiveIdAllowOverlap = false
	c.ActiveIdAllowNavDirFlags = 0
	c.ActiveIdClickOffset = f64.Vec2{-1, -1}
	c.ActiveIdWindow = nil
	c.ActiveIdSource = InputSourceNone
	c.MovingWindow = nil
	c.NextTreeNodeOpenVal = false
	c.NextTreeNodeOpenCond = 0

	c.WindowsById = make(map[string]*Window)
	c.NavWindow = nil
	c.NavId = 0
	c.NavActivateId = 0
	c.NavActivateDownId = 0
	c.NavActivatePressedId = 0
	c.NavInputId = 0
	c.NavJustTabbedId = 0
	c.NavJustMovedToId = 0
	c.NavNextActivateId = 0
	c.NavInputSource = InputSourceNone
	c.NavScoringRectScreen = f64.Rectangle{}
	c.NavScoringCount = 0
	c.NavWindowingTarget = nil
	c.NavWindowingHighlightTimer = 0
	c.NavWindowingHighlightAlpha = 0
	c.NavWindowingToggleLayer = false
	c.NavLayer = 0
	c.NavIdTabCounter = math.MaxInt32
	c.NavIdIsAlive = false
	c.NavMousePosDirty = false
	c.NavDisableHighlight = true
	c.NavDisableMouseHover = false
	c.NavAnyRequest = false
	c.NavInitRequest = false
	c.NavInitRequestFromMove = false
	c.NavInitResultId = 0
	c.NavMoveFromClampedRefRect = false
	c.NavMoveRequest = false
	c.NavMoveRequestForward = NavForwardNone
	c.NavMoveDir = DirNone
	c.NavMoveDirLast = DirNone

	c.ModalWindowDarkeningRatio = 0.0
	c.OverlayDrawList.Data = &c.DrawListSharedData
	c.OverlayDrawList.OwnerName = "##Overlay" // Give it a name for debugging
	c.MouseCursor = MouseCursorArrow

	c.DragDropActive = false
	c.DragDropSourceFlags = 0
	c.DragDropMouseButton = -1
	c.DragDropTargetId = 0
	c.DragDropAcceptIdCurrRectSurface = 0.0
	c.DragDropAcceptIdPrev = 0
	c.DragDropAcceptIdCurr = 0
	c.DragDropAcceptFrameCount = -1

	c.ScalarAsInputTextId = 0
	c.ColorEditOptions = ColorEditFlags_OptionsDefault
	c.DragCurrentValue = 0.0
	c.DragLastMouseDelta = f64.Vec2{0.0, 0.0}
	c.DragSpeedDefaultRatio = 1.0 / 100.0
	c.DragSpeedScaleSlow = 1.0 / 100.0
	c.DragSpeedScaleFast = 10.0
	c.ScrollbarClickDeltaToGrabCenter = f64.Vec2{0.0, 0.0}
	c.TooltipOverrideCount = 0
	c.PlatformImePos = f64.Vec2{-1.0, -1.0}
	c.PlatformImeLastPos = f64.Vec2{-1.0, -1.0}

	c.NextWindowData.Init()

	c.SettingsWindows = make(map[string]*WindowSettings)
	c.SettingsHandlers = make(map[string]*SettingsHandler)
	c.SettingsLoaded = false
	c.SettingsDirtyTimer = 0.0

	c.LogEnabled = false
	c.LogFile = nil
	c.LogClipboard = nil
	c.LogStartDepth = 0
	c.LogAutoExpandMaxDepth = 2

	for i := range c.FramerateSecPerFrame {
		c.FramerateSecPerFrame[i] = 0
	}
	c.FramerateSecPerFrameIdx = 0
	c.FramerateSecPerFrameAccum = 0.0
	c.WantCaptureMouseNextFrame = -1
	c.WantCaptureKeyboardNextFrame = -1
	c.WantTextInputNextFrame = -1
	c.InputTextState.Init(c)

	// Add .ini handle for ImGuiWindow type
	ini_handler := &SettingsHandler{
		TypeName:   "Window",
		ReadOpenFn: c.SettingsHandlerWindow_ReadOpen,
		ReadLineFn: c.SettingsHandlerWindow_ReadLine,
		WriteAllFn: c.SettingsHandlerWindow_WriteAll,
	}
	c.SettingsHandlers[ini_handler.TypeName] = ini_handler

	c.Initialized = true
}

func (c *Context) GetVersion() string {
	return "1.61 WIP"
}

func (c *Context) GetIO() *IO {
	return &c.IO
}

func (c *Context) GetStyle() *Style {
	return &c.Style
}

func (c *Context) GetCurrentWindowRead() *Window {
	c.CurrentWindow.WriteAccessed = true
	return c.CurrentWindow
}

func (c *Context) GetCurrentWindow() *Window {
	c.CurrentWindow.WriteAccessed = true
	return c.CurrentWindow
}

func (c *Context) SetFocusID(id ID, window *Window) {
	// Assume that SetFocusID() is called in the context where its NavLayer is the current layer, which is the case everywhere we call it.
	nav_layer := window.DC.NavLayerCurrent
	if c.NavWindow != window {
		c.NavInitRequest = false
	}
	c.NavId = id
	c.NavWindow = window
	c.NavLayer = nav_layer
	window.NavLastIds[nav_layer] = id
	if window.DC.LastItemId == id {
		window.NavRectRel[nav_layer] = f64.Rectangle{
			window.DC.LastItemRect.Min.Sub(window.Pos),
			window.DC.LastItemRect.Max.Sub(window.Pos),
		}
	}

	if c.ActiveIdSource == InputSourceNav {
		c.NavDisableMouseHover = true
	} else {
		c.NavDisableHighlight = true
	}
}

func (c *Context) KeepAliveID(id ID) {
	if c.ActiveId == id {
		c.ActiveIdIsAlive = true
	}
}

func (c *Context) ClearActiveID() {
	c.SetActiveID(0, nil)
}

func (c *Context) SetActiveID(id ID, window *Window) {
	c.ActiveIdIsJustActivated = c.ActiveId != id
	if c.ActiveIdIsJustActivated {
		c.ActiveIdTimer = 0
	}
	c.ActiveId = id
	c.ActiveIdAllowNavDirFlags = 0
	c.ActiveIdAllowOverlap = false
	c.ActiveIdWindow = window
	if id != 0 {
		c.ActiveIdIsAlive = true
		c.ActiveIdSource = InputSourceMouse
		if c.NavActivateId == id || c.NavInputId == id || c.NavJustTabbedId == id || c.NavJustMovedToId == id {
			c.ActiveIdSource = InputSourceNav
		}
	}
}

func (c *Context) GetFrameHeight() float64 {
	return c.FontSize + c.Style.FramePadding.Y*2
}

func (c *Context) GetFrameHeightWithSpacing() float64 {
	return c.FontSize + c.Style.FramePadding.Y*2 + c.Style.ItemSpacing.Y
}

func (c *Context) GetTime() float64 {
	return c.Time
}

func (c *Context) GetFrameCount() int {
	return c.FrameCount
}

func (c *Context) GetOverlayDrawList() *DrawList {
	return &c.OverlayDrawList
}

func (c *Context) GetDrawListSharedData() *DrawListSharedData {
	return &c.DrawListSharedData
}

func (c *Context) SetHoveredID(id ID) {
	c.HoveredId = id
	c.HoveredIdAllowOverlap = false
	if id != 0 && c.HoveredIdPreviousFrame == id {
		c.HoveredIdTimer = c.IO.DeltaTime
	} else {
		c.HoveredIdTimer = 0
	}
}

func (c *Context) PushStringID(str_id string) {
	window := c.GetCurrentWindowRead()
	window.IDStack = append(window.IDStack, window.GetID(str_id))
}

func (c *Context) PushID(id ID) {
	window := c.GetCurrentWindowRead()
	window.IDStack = append(window.IDStack, window.GetIntID(int(id)))
}

func (c *Context) PopID() {
	window := c.GetCurrentWindowRead()
	window.IDStack = window.IDStack[:len(window.IDStack)-1]
}

func (c *Context) GetStringID(str_id string) ID {
	return c.CurrentWindow.GetID(str_id)
}

func (c *Context) GetID(id ID) ID {
	return c.CurrentWindow.GetIntID(int(id))
}

type BackendFlags int

const (
	BackendFlagsHasGamepad      BackendFlags = 1 << 0 // Back-end supports and has a connected gamepad.
	BackendFlagsHasMouseCursors BackendFlags = 1 << 1 // Back-end supports reading GetMouseCursor() to change the OS cursor shape.
	BackendFlagsHasSetMousePos  BackendFlags = 1 << 2 // Back-end supports io.WantSetMousePos requests to reposition the OS mouse position (only used if ImGuiConfigFlags_NavEnableSetMousePos is set).
)

func (c *Context) CaptureKeyboardFromApp(capture bool) {
	c.WantCaptureKeyboardNextFrame = truth(capture)
}

func (c *Context) CaptureMouseFromApp(capture bool) {
	c.WantCaptureMouseNextFrame = truth(capture)
}