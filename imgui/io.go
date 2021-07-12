package imgui

import (
	"math"
	"runtime"

	"github.com/qeedquan/go-media/math/f64"
)

type IO struct {
	//------------------------------------------------------------------
	// Settings (fill once)                 // Default value:
	//------------------------------------------------------------------
	Ctx                     *Context
	ConfigFlags             ConfigFlags   // = 0                  // See ImGuiConfigFlags_ enum. Gamepad/keyboard navigation options, etc.
	BackendFlags            BackendFlags  // = 0                  // See ImGuiConfigFlags_ enum. Set by user/application. Gamepad/keyboard navigation options, etc.
	DisplaySize             f64.Vec2      // <unset>              // Display size, in pixels. For clamping windows positions.
	DeltaTime               float64       // = 1.0f/60.0f         // Time elapsed since last frame, in seconds.
	IniSavingRate           float64       // = 5.0f               // Maximum time between saving positions/sizes to .ini file, in seconds.
	IniFilename             string        // = "imgui.ini"        // Path to .ini file. NULL to disable .ini saving.
	LogFilename             string        // = "imgui_log.txt"    // Path to .log file (default parameter to ImGui::LogToFile when no file is specified).
	MouseDoubleClickTime    float64       // = 0.30f              // Time for a double-click, in seconds.
	MouseDoubleClickMaxDist float64       // = 6.0f               // Distance threshold to stay in to validate a double-click, in pixels.
	MouseDragThreshold      float64       // = 6.0f               // Distance threshold before considering we are dragging.
	KeyMap                  [KeyCOUNT]int // <unset>              // Map of indices into the KeysDown[512] entries array which represent your "native" keyboard state.
	KeyRepeatDelay          float64       // = 0.250f             // When holding a key/button, time before it starts repeating, in seconds (for buttons in Repeat mode, etc.).
	KeyRepeatRate           float64       // = 0.050f             // When holding a key/button, rate at which it repeats, in seconds.
	UserData                interface{}   // = NULL               // Store your own data for retrieval by callbacks.

	Fonts                   *FontAtlas // <auto>               // Load and assemble one or more fonts into a single tightly packed texture. Output to Fonts array.
	FontGlobalScale         float64    // = 1.0f               // Global scale all fonts
	FontAllowUserScaling    bool       // = false              // Allow user scaling text of individual window with CTRL+Wheel.
	FontDefault             *Font      // = NULL               // Font to use on NewFrame(). Use NULL to uses Fonts->Fonts[0].
	DisplayFramebufferScale f64.Vec2   // = (1.0f,1.0f)        // For retina display or other situations where window coordinates are different from framebuffer coordinates. User storage only, presently not used by ImGui.
	DisplayVisibleMin       f64.Vec2   // <unset> (0.0f,0.0f)  // If you use DisplaySize as a virtual space larger than your screen, set DisplayVisibleMin/Max to the visible area.
	DisplayVisibleMax       f64.Vec2   // <unset> (0.0f,0.0f)  // If the values are the same, we defaults to Min=(0.0f) and Max=DisplaySize

	// Advanced/subtle behaviors
	OptMacOSXBehaviors bool // = defined(__APPLE__) // OS X style: Text editing cursor movement using Alt instead of Ctrl, Shortcuts using Cmd/Super instead of Ctrl, Line/Text Start and End using Cmd+Arrows instead of Home/End, Double click selects by word instead of selecting whole text, Multi-selection in lists uses Cmd/Super instead of Ctrl
	OptCursorBlink     bool // = true               // Enable blinking cursor, for users who consider it annoying.

	//------------------------------------------------------------------
	// Settings (User Functions)
	//------------------------------------------------------------------

	// Optional: access OS clipboard
	// (default to use native Win32 clipboard on Windows, otherwise uses a private clipboard. Override to access OS clipboard on other architectures)
	GetClipboardTextFn func() string
	SetClipboardTextFn func(text string)
	ClipboardUserData  interface{}

	// Optional: notify OS Input Method Editor of the screen position of your cursor for text input position (e.g. when using Japanese/Chinese IME in Windows)
	// (default to use native imm32 api on Windows)
	ImeSetInputScreenPosFn func(x, y int)
	ImeWindowHandle        interface{} // (Windows) Set this to your HWND to get automatic IME cursor positioning.

	//------------------------------------------------------------------
	// Input - Fill before calling NewFrame()
	//------------------------------------------------------------------

	MousePos        f64.Vec2               // Mouse position, in pixels. Set to ImVec2(-FLT_MAX,-FLT_MAX) if mouse is unavailable (on another screen, etc.)
	MouseDown       [5]bool                // Mouse buttons: left, right, middle + extras. ImGui itself mostly only uses left button (BeginPopupContext** are using right button). Others buttons allows us to track if the mouse is being used by your application + available to user as a convenience via IsMouse** API.
	MouseWheel      float64                // Mouse wheel Vertical: 1 unit scrolls about 5 lines text.
	MouseWheelH     float64                // Mouse wheel Horizontal. Most users don't have a mouse with an horizontal wheel, may not be filled by all back-ends.
	MouseDrawCursor bool                   // Request ImGui to draw a mouse cursor for you (if you are on a platform without a mouse cursor).
	KeyCtrl         bool                   // Keyboard modifier pressed: Control
	KeyShift        bool                   // Keyboard modifier pressed: Shift
	KeyAlt          bool                   // Keyboard modifier pressed: Alt
	KeySuper        bool                   // Keyboard modifier pressed: Cmd/Super/Windows
	KeysDown        [512]bool              // Keyboard keys that are pressed (ideally left in the "native" order your engine has access to keyboard keys, so you can use your own defines/enums for keys).
	InputCharacters [16 + 1]rune           // List of characters input (translated by user from keypress+keyboard state). Fill using AddInputCharacter() helper.
	NavInputs       [NavInputCOUNT]float64 // Gamepad inputs (keyboard keys will be auto-mapped and be written here by ImGui::NewFrame)

	//------------------------------------------------------------------
	// Output - Retrieve after calling NewFrame()
	//------------------------------------------------------------------

	WantCaptureMouse      bool     // When io.WantCaptureMouse is true, imgui will use the mouse inputs, do not dispatch them to your main game/application. (e.g. unclicked mouse is hovering over an imgui window, widget is active, mouse was clicked over an imgui window, etc.).
	WantCaptureKeyboard   bool     // When io.WantCaptureKeyboard is true, imgui will use the keyboard inputs, do not dispatch them to your main game/application. (e.g. InputText active, or an imgui window is focused and navigation is enabled, etc.).
	WantTextInput         bool     // Mobile/console: when io.WantTextInput is true, you may display an on-screen keyboard. This is set by ImGui when it wants textual keyboard input to happen (e.g. when a InputText widget is active).
	WantSetMousePos       bool     // MousePos has been altered, back-end should reposition mouse on next frame. Set only when ImGuiConfigFlags_NavMoveMouse flag is enabled.
	NavActive             bool     // Directional navigation is currently allowed (will handle ImGuiKey_NavXXX events) = a window is focused and it doesn't use the ImGuiWindowFlags_NoNavInputs flag.
	NavVisible            bool     // Directional navigation is visible and allowed (will handle ImGuiKey_NavXXX events).
	Framerate             float64  // Application framerate estimation, in frame per second. Solely for convenience. Rolling average estimation based on IO.DeltaTime over 120 frames
	MetricsRenderVertices int      // Vertices output during last call to Render()
	MetricsRenderIndices  int      // Indices output during last call to Render() = number of triangles * 3
	MetricsActiveWindows  int      // Number of visible root windows (exclude child windows)
	MouseDelta            f64.Vec2 // Mouse delta. Note that this is zero if either current or previous position are invalid (-FLT_MAX,-FLT_MAX), so a disappearing/reappearing mouse won't have a huge delta.

	//------------------------------------------------------------------
	// [Internal] ImGui will maintain those fields. Forward compatibility not guaranteed!
	//------------------------------------------------------------------

	MousePosPrev              f64.Vec2     // Previous mouse position temporary storage (nb: not for public use, set to MousePos in NewFrame())
	MouseClickedPos           [5]f64.Vec2  // Position at time of clicking
	MouseClickedTime          [5]float64   // Time of last click (used to figure out double-click)
	MouseClicked              [5]bool      // Mouse button went from !Down to Down
	MouseDoubleClicked        [5]bool      // Has mouse button been double-clicked?
	MouseReleased             [5]bool      // Mouse button went from Down to !Down
	MouseDownOwned            [5]bool      // Track if button was clicked inside a window. We don't request mouse capture from the application if click started outside ImGui bounds.
	MouseDownDuration         [5]float64   // Duration the mouse button has been down (0.0f == just clicked)
	MouseDownDurationPrev     [5]float64   // Previous time the mouse button has been down
	MouseDragMaxDistanceAbs   [5]f64.Vec2  // Maximum distance, absolute, on each axis, of how much mouse has traveled from the clicking point
	MouseDragMaxDistanceSqr   [5]float64   // Squared maximum distance of how much mouse has traveled from the clicking point
	KeysDownDuration          [512]float64 // Duration the keyboard key has been down (0.0f == just pressed)
	KeysDownDurationPrev      [512]float64 // Previous duration the key has been down
	NavInputsDownDuration     [NavInputCOUNT]float64
	NavInputsDownDurationPrev [NavInputCOUNT]float64
}

func (c *IO) Init(ctx *Context) {
	*c = IO{
		Ctx: ctx,
	}

	// Settings
	c.ConfigFlags = 0x00
	c.BackendFlags = 0x00
	c.DisplaySize = f64.Vec2{-1.0, -1.0}
	c.DeltaTime = 1.0 / 60.0
	c.IniSavingRate = 5.0
	c.IniFilename = ""
	c.LogFilename = ""
	c.MouseDoubleClickTime = 0.30
	c.MouseDoubleClickMaxDist = 6.0
	for i := range c.KeyMap {
		c.KeyMap[i] = -1
	}
	c.KeyRepeatDelay = 0.250
	c.KeyRepeatRate = 0.050

	c.Fonts = nil
	c.FontGlobalScale = 1.0
	c.FontDefault = nil
	c.FontAllowUserScaling = false
	c.DisplayFramebufferScale = f64.Vec2{1.0, 1.0}
	c.DisplayVisibleMin = f64.Vec2{0.0, 0.0}
	c.DisplayVisibleMax = f64.Vec2{0.0, 0.0}

	if runtime.GOOS == "darwin" {
		c.OptMacOSXBehaviors = true
	} else {
		c.OptMacOSXBehaviors = false
	}
	c.OptCursorBlink = true

	// Settings (User Functions)
	c.GetClipboardTextFn = ctx.GetClipboardTextFn_DefaultImpl
	c.SetClipboardTextFn = ctx.SetClipboardTextFn_DefaultImpl
	c.ImeSetInputScreenPosFn = ImeSetInputScreenPosFn_DefaultImpl
	c.ImeWindowHandle = nil

	// Input (NB: we already have memset zero the entire structure)
	c.MousePos = f64.Vec2{-math.MaxFloat32, -math.MaxFloat32}
	c.MousePosPrev = f64.Vec2{-math.MaxFloat32, -math.MaxFloat32}
	c.MouseDragThreshold = 6.0
	for i := range c.MouseDownDuration {
		c.MouseDownDuration[i] = -1
	}
	for i := range c.KeysDownDuration {
		c.KeysDownDuration[i] = -1
	}
	for i := range c.NavInputsDownDuration {
		c.NavInputsDownDuration[i] = -1
	}
}

func ImeSetInputScreenPosFn_DefaultImpl(x, y int) {
}

func (c *Context) GetClipboardTextFn_DefaultImpl() string {
	return c.PrivateClipboard
}

func (c *Context) SetClipboardTextFn_DefaultImpl(text string) {
	c.PrivateClipboard = text
}
