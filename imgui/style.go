package imgui

import (
	"image/color"

	"github.com/qeedquan/go-media/image/chroma"
	"github.com/qeedquan/go-media/math/f64"
)

type Col int

const (
	ColText Col = iota
	ColTextDisabled
	ColWindowBg // Background of normal windows
	ColChildBg  // Background of child windows
	ColPopupBg  // Background of popups menus tooltips windows
	ColBorder
	ColBorderShadow
	ColFrameBg // Background of checkbox radio button plot slider text input
	ColFrameBgHovered
	ColFrameBgActive
	ColTitleBg
	ColTitleBgActive
	ColTitleBgCollapsed
	ColMenuBarBg
	ColScrollbarBg
	ColScrollbarGrab
	ColScrollbarGrabHovered
	ColScrollbarGrabActive
	ColCheckMark
	ColSliderGrab
	ColSliderGrabActive
	ColButton
	ColButtonHovered
	ColButtonActive
	ColHeader
	ColHeaderHovered
	ColHeaderActive
	ColSeparator
	ColSeparatorHovered
	ColSeparatorActive
	ColResizeGrip
	ColResizeGripHovered
	ColResizeGripActive
	ColPlotLines
	ColPlotLinesHovered
	ColPlotHistogram
	ColPlotHistogramHovered
	ColTextSelectedBg
	ColModalWindowDarkening // Darken/colorize entire screen behind a modal window, when one is active
	ColDragDropTarget
	ColNavHighlight          // Gamepad/keyboard: current highlighted item
	ColNavWindowingHighlight // Gamepad/keyboard: when holding NavMenu to focus/move/resize windows
	ColCOUNT
)

type ColMod struct {
	Col         Col
	BackupValue f64.Vec4
}

type StyleMod struct {
	VarIdx StyleVar
	Value  interface{}
}

type Style struct {
	Alpha                  float64  // Global alpha applies to everything in ImGui.
	WindowPadding          f64.Vec2 // Padding within a window.
	WindowRounding         float64  // Radius of window corners rounding. Set to 0.0f to have rectangular windows.
	WindowBorderSize       float64  // Thickness of border around windows. Generally set to 0.0f or 1.0f. (Other values are not well tested and more CPU/GPU costly).
	WindowMinSize          f64.Vec2 // Minimum window size. This is a global setting. If you want to constraint individual windows, use SetNextWindowSizeConstraints().
	WindowTitleAlign       f64.Vec2 // Alignment for title bar text. Defaults to (0.0,0.5f) for left-aligned,vertically centered.
	ChildRounding          float64  // Radius of child window corners rounding. Set to 0.0f to have rectangular windows.
	ChildBorderSize        float64  // Thickness of border around child windows. Generally set to 0.0f or 1.0f. (Other values are not well tested and more CPU/GPU costly).
	PopupRounding          float64  // Radius of popup window corners rounding.
	PopupBorderSize        float64  // Thickness of border around popup windows. Generally set to 0.0f or 1.0f. (Other values are not well tested and more CPU/GPU costly).
	FramePadding           f64.Vec2 // Padding within a framed rectangle (used by most widgets).
	FrameRounding          float64  // Radius of frame corners rounding. Set to 0.0f to have rectangular frame (used by most widgets).
	FrameBorderSize        float64  // Thickness of border around frames. Generally set to 0.0f or 1.0f. (Other values are not well tested and more CPU/GPU costly).
	ItemSpacing            f64.Vec2 // Horizontal and vertical spacing between widgets/lines.
	ItemInnerSpacing       f64.Vec2 // Horizontal and vertical spacing between within elements of a composed widget (e.g. a slider and its label).
	TouchExtraPadding      f64.Vec2 // Expand reactive bounding box for touch-based system where touch position is not accurate enough. Unfortunately we don't sort widgets so priority on overlap will always be given to the first widget. So don't grow this too much!
	IndentSpacing          float64  // Horizontal indentation when e.g. entering a tree node. Generally == (FontSize + FramePadding.x*2).
	ColumnsMinSpacing      float64  // Minimum horizontal spacing between two columns.
	ScrollbarSize          float64  // Width of the vertical scrollbar, Height of the horizontal scrollbar.
	ScrollbarRounding      float64  // Radius of grab corners for scrollbar.
	GrabMinSize            float64  // Minimum width/height of a grab box for slider/scrollbar.
	GrabRounding           float64  // Radius of grabs corners rounding. Set to 0.0f to have rectangular slider grabs.
	ButtonTextAlign        f64.Vec2 // Alignment of button text when button is larger than text. Defaults to (0.5,0.5f) for horizontally+vertically centered.
	DisplayWindowPadding   f64.Vec2 // Window positions are clamped to be visible within the display area by at least this amount. Only covers regular windows.
	DisplaySafeAreaPadding f64.Vec2 // If you cannot see the edge of your screen (e.g. on a TV) increase the safe area padding. Covers popups/tooltips as well regular windows. NB: Prefer configuring your TV sets correctly!
	MouseCursorScale       float64  // Scale software rendered mouse cursor (when io.MouseDrawCursor is enabled). May be removed later.
	AntiAliasedLines       bool     // Enable anti-aliasing on lines/borders. Disable if you are really tight on CPU/GPU.
	AntiAliasedFill        bool     // Enable anti-aliasing on filled shapes (rounded rectangles, circles, etc.)
	CurveTessellationTol   float64  // Tessellation tolerance when using PathBezierCurveTo() without a specific number of segments. Decrease for highly tessellated curves (higher quality, more polygons), increase to reduce quality.
	Colors                 [ColCOUNT]f64.Vec4
}

// Enumeration for PushStyleVar() / PopStyleVar() to temporarily modify the ImGuiStyle structure.
// NB: the enum only refers to fields of ImGuiStyle which makes sense to be pushed/popped inside UI code. During initialization, feel free to just poke into ImGuiStyle directly.
// NB: if changing this enum, you need to update the associated internal table GStyleVarInfo[] accordingly. This is where we link enum values to members offset/type.
type StyleVar int

const (
	// Enum name ......................// Member in ImGuiStyle structure (see ImGuiStyle for descriptions)
	StyleVarAlpha             StyleVar = iota // float     Alpha
	StyleVarWindowPadding                     // ImVec2    WindowPadding
	StyleVarWindowRounding                    // float     WindowRounding
	StyleVarWindowBorderSize                  // float     WindowBorderSize
	StyleVarWindowMinSize                     // ImVec2    WindowMinSize
	StyleVarWindowTitleAlign                  // ImVec2    WindowTitleAlign
	StyleVarChildRounding                     // float     ChildRounding
	StyleVarChildBorderSize                   // float     ChildBorderSize
	StyleVarPopupRounding                     // float     PopupRounding
	StyleVarPopupBorderSize                   // float     PopupBorderSize
	StyleVarFramePadding                      // ImVec2    FramePadding
	StyleVarFrameRounding                     // float     FrameRounding
	StyleVarFrameBorderSize                   // float     FrameBorderSize
	StyleVarItemSpacing                       // ImVec2    ItemSpacing
	StyleVarItemInnerSpacing                  // ImVec2    ItemInnerSpacing
	StyleVarIndentSpacing                     // float     IndentSpacing
	StyleVarScrollbarSize                     // float     ScrollbarSize
	StyleVarScrollbarRounding                 // float     ScrollbarRounding
	StyleVarGrabMinSize                       // float     GrabMinSize
	StyleVarGrabRounding                      // float     GrabRounding
	StyleVarButtonTextAlign                   // ImVec2    ButtonTextAlign
	StyleVarCOUNT
)

func (c *Context) StyleColorsClassic(style *Style) {
	if style == nil {
		style = c.GetStyle()
	}
	colors := style.Colors[:]

	colors[ColText] = f64.Vec4{0.90, 0.90, 0.90, 1.00}
	colors[ColTextDisabled] = f64.Vec4{0.60, 0.60, 0.60, 1.00}
	colors[ColWindowBg] = f64.Vec4{0.00, 0.00, 0.00, 0.70}
	colors[ColChildBg] = f64.Vec4{0.00, 0.00, 0.00, 0.00}
	colors[ColPopupBg] = f64.Vec4{0.11, 0.11, 0.14, 0.92}
	colors[ColBorder] = f64.Vec4{0.50, 0.50, 0.50, 0.50}
	colors[ColBorderShadow] = f64.Vec4{0.00, 0.00, 0.00, 0.00}
	colors[ColFrameBg] = f64.Vec4{0.43, 0.43, 0.43, 0.39}
	colors[ColFrameBgHovered] = f64.Vec4{0.47, 0.47, 0.69, 0.40}
	colors[ColFrameBgActive] = f64.Vec4{0.42, 0.41, 0.64, 0.69}
	colors[ColTitleBg] = f64.Vec4{0.27, 0.27, 0.54, 0.83}
	colors[ColTitleBgActive] = f64.Vec4{0.32, 0.32, 0.63, 0.87}
	colors[ColTitleBgCollapsed] = f64.Vec4{0.40, 0.40, 0.80, 0.20}
	colors[ColMenuBarBg] = f64.Vec4{0.40, 0.40, 0.55, 0.80}
	colors[ColScrollbarBg] = f64.Vec4{0.20, 0.25, 0.30, 0.60}
	colors[ColScrollbarGrab] = f64.Vec4{0.40, 0.40, 0.80, 0.30}
	colors[ColScrollbarGrabHovered] = f64.Vec4{0.40, 0.40, 0.80, 0.40}
	colors[ColScrollbarGrabActive] = f64.Vec4{0.41, 0.39, 0.80, 0.60}
	colors[ColCheckMark] = f64.Vec4{0.90, 0.90, 0.90, 0.50}
	colors[ColSliderGrab] = f64.Vec4{1.00, 1.00, 1.00, 0.30}
	colors[ColSliderGrabActive] = f64.Vec4{0.41, 0.39, 0.80, 0.60}
	colors[ColButton] = f64.Vec4{0.35, 0.40, 0.61, 0.62}
	colors[ColButtonHovered] = f64.Vec4{0.40, 0.48, 0.71, 0.79}
	colors[ColButtonActive] = f64.Vec4{0.46, 0.54, 0.80, 1.00}
	colors[ColHeader] = f64.Vec4{0.40, 0.40, 0.90, 0.45}
	colors[ColHeaderHovered] = f64.Vec4{0.45, 0.45, 0.90, 0.80}
	colors[ColHeaderActive] = f64.Vec4{0.53, 0.53, 0.87, 0.80}
	colors[ColSeparator] = f64.Vec4{0.50, 0.50, 0.50, 1.00}
	colors[ColSeparatorHovered] = f64.Vec4{0.60, 0.60, 0.70, 1.00}
	colors[ColSeparatorActive] = f64.Vec4{0.70, 0.70, 0.90, 1.00}
	colors[ColResizeGrip] = f64.Vec4{1.00, 1.00, 1.00, 0.16}
	colors[ColResizeGripHovered] = f64.Vec4{0.78, 0.82, 1.00, 0.60}
	colors[ColResizeGripActive] = f64.Vec4{0.78, 0.82, 1.00, 0.90}
	colors[ColPlotLines] = f64.Vec4{1.00, 1.00, 1.00, 1.00}
	colors[ColPlotLinesHovered] = f64.Vec4{0.90, 0.70, 0.00, 1.00}
	colors[ColPlotHistogram] = f64.Vec4{0.90, 0.70, 0.00, 1.00}
	colors[ColPlotHistogramHovered] = f64.Vec4{1.00, 0.60, 0.00, 1.00}
	colors[ColTextSelectedBg] = f64.Vec4{0.00, 0.00, 1.00, 0.35}
	colors[ColModalWindowDarkening] = f64.Vec4{0.20, 0.20, 0.20, 0.35}
	colors[ColDragDropTarget] = f64.Vec4{1.00, 1.00, 0.00, 0.90}
	colors[ColNavHighlight] = colors[ColHeaderHovered]
	colors[ColNavWindowingHighlight] = f64.Vec4{1.00, 1.00, 1.00, 0.70}
}

// Those light colors are better suited with a thicker font than the default one + FrameBorder
func (c *Context) StyleColorsLight(style *Style) {
	if style == nil {
		style = c.GetStyle()
	}
	colors := style.Colors[:]

	colors[ColText] = f64.Vec4{0.00, 0.00, 0.00, 1.00}
	colors[ColTextDisabled] = f64.Vec4{0.60, 0.60, 0.60, 1.00}
	colors[ColWindowBg] = f64.Vec4{0.94, 0.94, 0.94, 1.00}
	colors[ColChildBg] = f64.Vec4{0.00, 0.00, 0.00, 0.00}
	colors[ColPopupBg] = f64.Vec4{1.00, 1.00, 1.00, 0.98}
	colors[ColBorder] = f64.Vec4{0.00, 0.00, 0.00, 0.30}
	colors[ColBorderShadow] = f64.Vec4{0.00, 0.00, 0.00, 0.00}
	colors[ColFrameBg] = f64.Vec4{1.00, 1.00, 1.00, 1.00}
	colors[ColFrameBgHovered] = f64.Vec4{0.26, 0.59, 0.98, 0.40}
	colors[ColFrameBgActive] = f64.Vec4{0.26, 0.59, 0.98, 0.67}
	colors[ColTitleBg] = f64.Vec4{0.96, 0.96, 0.96, 1.00}
	colors[ColTitleBgActive] = f64.Vec4{0.82, 0.82, 0.82, 1.00}
	colors[ColTitleBgCollapsed] = f64.Vec4{1.00, 1.00, 1.00, 0.51}
	colors[ColMenuBarBg] = f64.Vec4{0.86, 0.86, 0.86, 1.00}
	colors[ColScrollbarBg] = f64.Vec4{0.98, 0.98, 0.98, 0.53}
	colors[ColScrollbarGrab] = f64.Vec4{0.69, 0.69, 0.69, 0.80}
	colors[ColScrollbarGrabHovered] = f64.Vec4{0.49, 0.49, 0.49, 0.80}
	colors[ColScrollbarGrabActive] = f64.Vec4{0.49, 0.49, 0.49, 1.00}
	colors[ColCheckMark] = f64.Vec4{0.26, 0.59, 0.98, 1.00}
	colors[ColSliderGrab] = f64.Vec4{0.26, 0.59, 0.98, 0.78}
	colors[ColSliderGrabActive] = f64.Vec4{0.46, 0.54, 0.80, 0.60}
	colors[ColButton] = f64.Vec4{0.26, 0.59, 0.98, 0.40}
	colors[ColButtonHovered] = f64.Vec4{0.26, 0.59, 0.98, 1.00}
	colors[ColButtonActive] = f64.Vec4{0.06, 0.53, 0.98, 1.00}
	colors[ColHeader] = f64.Vec4{0.26, 0.59, 0.98, 0.31}
	colors[ColHeaderHovered] = f64.Vec4{0.26, 0.59, 0.98, 0.80}
	colors[ColHeaderActive] = f64.Vec4{0.26, 0.59, 0.98, 1.00}
	colors[ColSeparator] = f64.Vec4{0.39, 0.39, 0.39, 1.00}
	colors[ColSeparatorHovered] = f64.Vec4{0.14, 0.44, 0.80, 0.78}
	colors[ColSeparatorActive] = f64.Vec4{0.14, 0.44, 0.80, 1.00}
	colors[ColResizeGrip] = f64.Vec4{0.80, 0.80, 0.80, 0.56}
	colors[ColResizeGripHovered] = f64.Vec4{0.26, 0.59, 0.98, 0.67}
	colors[ColResizeGripActive] = f64.Vec4{0.26, 0.59, 0.98, 0.95}
	colors[ColPlotLines] = f64.Vec4{0.39, 0.39, 0.39, 1.00}
	colors[ColPlotLinesHovered] = f64.Vec4{1.00, 0.43, 0.35, 1.00}
	colors[ColPlotHistogram] = f64.Vec4{0.90, 0.70, 0.00, 1.00}
	colors[ColPlotHistogramHovered] = f64.Vec4{1.00, 0.45, 0.00, 1.00}
	colors[ColTextSelectedBg] = f64.Vec4{0.26, 0.59, 0.98, 0.35}
	colors[ColModalWindowDarkening] = f64.Vec4{0.20, 0.20, 0.20, 0.35}
	colors[ColDragDropTarget] = f64.Vec4{0.26, 0.59, 0.98, 0.95}
	colors[ColNavHighlight] = colors[ColHeaderHovered]
	colors[ColNavWindowingHighlight] = f64.Vec4{0.70, 0.70, 0.70, 0.70}
}

func (c *Context) StyleColorsDark(style *Style) {
	if style == nil {
		style = c.GetStyle()
	}
	colors := style.Colors[:]

	colors[ColText] = f64.Vec4{1.00, 1.00, 1.00, 1.00}
	colors[ColTextDisabled] = f64.Vec4{0.50, 0.50, 0.50, 1.00}
	colors[ColWindowBg] = f64.Vec4{0.06, 0.06, 0.06, 0.94}
	colors[ColChildBg] = f64.Vec4{1.00, 1.00, 1.00, 0.00}
	colors[ColPopupBg] = f64.Vec4{0.08, 0.08, 0.08, 0.94}
	colors[ColBorder] = f64.Vec4{0.43, 0.43, 0.50, 0.50}
	colors[ColBorderShadow] = f64.Vec4{0.00, 0.00, 0.00, 0.00}
	colors[ColFrameBg] = f64.Vec4{0.16, 0.29, 0.48, 0.54}
	colors[ColFrameBgHovered] = f64.Vec4{0.26, 0.59, 0.98, 0.40}
	colors[ColFrameBgActive] = f64.Vec4{0.26, 0.59, 0.98, 0.67}
	colors[ColTitleBg] = f64.Vec4{0.04, 0.04, 0.04, 1.00}
	colors[ColTitleBgActive] = f64.Vec4{0.16, 0.29, 0.48, 1.00}
	colors[ColTitleBgCollapsed] = f64.Vec4{0.00, 0.00, 0.00, 0.51}
	colors[ColMenuBarBg] = f64.Vec4{0.14, 0.14, 0.14, 1.00}
	colors[ColScrollbarBg] = f64.Vec4{0.02, 0.02, 0.02, 0.53}
	colors[ColScrollbarGrab] = f64.Vec4{0.31, 0.31, 0.31, 1.00}
	colors[ColScrollbarGrabHovered] = f64.Vec4{0.41, 0.41, 0.41, 1.00}
	colors[ColScrollbarGrabActive] = f64.Vec4{0.51, 0.51, 0.51, 1.00}
	colors[ColCheckMark] = f64.Vec4{0.26, 0.59, 0.98, 1.00}
	colors[ColSliderGrab] = f64.Vec4{0.24, 0.52, 0.88, 1.00}
	colors[ColSliderGrabActive] = f64.Vec4{0.26, 0.59, 0.98, 1.00}
	colors[ColButton] = f64.Vec4{0.26, 0.59, 0.98, 0.40}
	colors[ColButtonHovered] = f64.Vec4{0.26, 0.59, 0.98, 1.00}
	colors[ColButtonActive] = f64.Vec4{0.06, 0.53, 0.98, 1.00}
	colors[ColHeader] = f64.Vec4{0.26, 0.59, 0.98, 0.31}
	colors[ColHeaderHovered] = f64.Vec4{0.26, 0.59, 0.98, 0.80}
	colors[ColHeaderActive] = f64.Vec4{0.26, 0.59, 0.98, 1.00}
	colors[ColSeparator] = colors[ColBorder]
	colors[ColSeparatorHovered] = f64.Vec4{0.10, 0.40, 0.75, 0.78}
	colors[ColSeparatorActive] = f64.Vec4{0.10, 0.40, 0.75, 1.00}
	colors[ColResizeGrip] = f64.Vec4{0.26, 0.59, 0.98, 0.25}
	colors[ColResizeGripHovered] = f64.Vec4{0.26, 0.59, 0.98, 0.67}
	colors[ColResizeGripActive] = f64.Vec4{0.26, 0.59, 0.98, 0.95}
	colors[ColPlotLines] = f64.Vec4{0.61, 0.61, 0.61, 1.00}
	colors[ColPlotLinesHovered] = f64.Vec4{1.00, 0.43, 0.35, 1.00}
	colors[ColPlotHistogram] = f64.Vec4{0.90, 0.70, 0.00, 1.00}
	colors[ColPlotHistogramHovered] = f64.Vec4{1.00, 0.60, 0.00, 1.00}
	colors[ColTextSelectedBg] = f64.Vec4{0.26, 0.59, 0.98, 0.35}
	colors[ColModalWindowDarkening] = f64.Vec4{0.80, 0.80, 0.80, 0.35}
	colors[ColDragDropTarget] = f64.Vec4{1.00, 1.00, 0.00, 0.90}
	colors[ColNavHighlight] = f64.Vec4{0.26, 0.59, 0.98, 1.00}
	colors[ColNavWindowingHighlight] = f64.Vec4{1.00, 1.00, 1.00, 0.70}
}

func (c *Context) GetColorFromStyle(idx Col) color.RGBA {
	return c.GetColorFromStyleWithAlpha(idx, 1)
}

func (c *Context) GetColorFromStyleWithAlpha(idx Col, alpha_mul float64) color.RGBA {
	style := &c.Style
	col := style.Colors[idx]
	col.W *= style.Alpha * alpha_mul
	return col.ToRGBA()
}

func (s *Style) Init() {
	s.Alpha = 1.0                             // Global alpha applies to everything in ImGui
	s.WindowPadding = f64.Vec2{8, 8}          // Padding within a window
	s.WindowRounding = 7.0                    // Radius of window corners rounding. Set to 0.0f to have rectangular windows
	s.WindowBorderSize = 1.0                  // Thickness of border around windows. Generally set to 0.0f or 1.0f. Other values not well tested.
	s.WindowMinSize = f64.Vec2{32, 32}        // Minimum window size
	s.WindowTitleAlign = f64.Vec2{0.0, 0.5}   // Alignment for title bar text
	s.ChildRounding = 0.0                     // Radius of child window corners rounding. Set to 0.0f to have rectangular child windows
	s.ChildBorderSize = 1.0                   // Thickness of border around child windows. Generally set to 0.0f or 1.0f. Other values not well tested.
	s.PopupRounding = 0.0                     // Radius of popup window corners rounding. Set to 0.0f to have rectangular child windows
	s.PopupBorderSize = 1.0                   // Thickness of border around popup or tooltip windows. Generally set to 0.0f or 1.0f. Other values not well tested.
	s.FramePadding = f64.Vec2{4, 3}           // Padding within a framed rectangle (used by most widgets)
	s.FrameRounding = 0.0                     // Radius of frame corners rounding. Set to 0.0f to have rectangular frames (used by most widgets).
	s.FrameBorderSize = 0.0                   // Thickness of border around frames. Generally set to 0.0f or 1.0f. Other values not well tested.
	s.ItemSpacing = f64.Vec2{8, 4}            // Horizontal and vertical spacing between widgets/lines
	s.ItemInnerSpacing = f64.Vec2{4, 4}       // Horizontal and vertical spacing between within elements of a composed widget (e.g. a slider and its label)
	s.TouchExtraPadding = f64.Vec2{0, 0}      // Expand reactive bounding box for touch-based system where touch position is not accurate enough. Unfortunately we don't sort widgets so priority on overlap will always be given to the first widget. So don't grow this too much!
	s.IndentSpacing = 21.0                    // Horizontal spacing when e.g. entering a tree node. Generally == (FontSize + FramePadding.x*2).
	s.ColumnsMinSpacing = 6.0                 // Minimum horizontal spacing between two columns
	s.ScrollbarSize = 16.0                    // Width of the vertical scrollbar, Height of the horizontal scrollbar
	s.ScrollbarRounding = 9.0                 // Radius of grab corners rounding for scrollbar
	s.GrabMinSize = 10.0                      // Minimum width/height of a grab box for slider/scrollbar
	s.GrabRounding = 0.0                      // Radius of grabs corners rounding. Set to 0.0f to have rectangular slider grabs.
	s.ButtonTextAlign = f64.Vec2{0.5, 0.5}    // Alignment of button text when button is larger than text.
	s.DisplayWindowPadding = f64.Vec2{20, 20} // Window positions are clamped to be visible within the display area by at least this amount. Only covers regular windows.
	s.DisplaySafeAreaPadding = f64.Vec2{3, 3} // If you cannot see the edge of your screen (e.g. on a TV) increase the safe area padding. Covers popups/tooltips as well regular windows.
	s.MouseCursorScale = 1.0                  // Scale software rendered mouse cursor (when io.MouseDrawCursor is enabled). May be removed later.
	s.AntiAliasedLines = true                 // Enable anti-aliasing on lines/borders. Disable if you are really short on CPU/GPU.
	s.AntiAliasedFill = true                  // Enable anti-aliasing on filled shapes (rounded rectangles, circles, etc.)
	s.CurveTessellationTol = 1.25             // Tessellation tolerance when using PathBezierCurveTo() without a specific number of segments. Decrease for highly tessellated curves (higher quality, more polygons), increase to reduce quality.
}

func (c *Context) GetStyleIndex(style *Style, idx StyleVar) interface{} {
	if style == nil {
		style = &c.Style
	}
	switch idx {
	case StyleVarAlpha:
		return &style.Alpha
	case StyleVarWindowPadding:
		return &style.WindowPadding
	case StyleVarWindowRounding:
		return &style.WindowRounding
	case StyleVarWindowBorderSize:
		return &style.WindowBorderSize
	case StyleVarWindowMinSize:
		return &style.WindowMinSize
	case StyleVarWindowTitleAlign:
		return &style.WindowTitleAlign
	case StyleVarChildRounding:
		return &style.ChildRounding
	case StyleVarChildBorderSize:
		return &style.ChildBorderSize
	case StyleVarPopupRounding:
		return &style.PopupRounding
	case StyleVarPopupBorderSize:
		return &style.PopupBorderSize
	case StyleVarFramePadding:
		return &style.FramePadding
	case StyleVarFrameRounding:
		return &style.FrameRounding
	case StyleVarFrameBorderSize:
		return &style.FrameBorderSize
	case StyleVarItemSpacing:
		return &style.ItemSpacing
	case StyleVarItemInnerSpacing:
		return &style.ItemInnerSpacing
	case StyleVarIndentSpacing:
		return &style.IndentSpacing
	case StyleVarScrollbarSize:
		return &style.ScrollbarSize
	case StyleVarScrollbarRounding:
		return &style.ScrollbarRounding
	case StyleVarGrabMinSize:
		return &style.GrabMinSize
	case StyleVarGrabRounding:
		return &style.GrabRounding
	case StyleVarButtonTextAlign:
		return &style.ButtonTextAlign
	default:
		panic("unreachable")
	}
}

func (c *Context) PushStyleVar(idx StyleVar, val interface{}) {
	v := c.GetStyleIndex(nil, idx)
	switch v := v.(type) {
	case *float64:
		c.StyleModifiers = append(c.StyleModifiers, StyleMod{idx, *v})
		*v = val.(float64)
	case *f64.Vec2:
		c.StyleModifiers = append(c.StyleModifiers, StyleMod{idx, *v})
		*v = val.(f64.Vec2)
	default:
		panic("unreachable")
	}
}

func (c *Context) PopStyleVar() {
	c.PopStyleVarN(1)
}

func (c *Context) PopStyleVarN(count int) {
	style := &c.Style
	for ; count > 0 && len(c.StyleModifiers) > 0; count-- {
		n := len(c.StyleModifiers) - 1
		m := c.StyleModifiers[n]
		c.StyleModifiers = c.StyleModifiers[:n]

		v := c.GetStyleIndex(style, m.VarIdx)
		switch v := v.(type) {
		case *float64:
			*v = m.Value.(float64)
		case *f64.Vec2:
			*v = m.Value.(f64.Vec2)
		default:
			panic("unreachable")
		}
	}
}

func (c *Context) PushStyleColorV4(idx Col, col f64.Vec4) {
	c.ColorModifiers = append(c.ColorModifiers, ColMod{idx, c.Style.Colors[idx]})
	c.Style.Colors[idx] = col
}

func (c *Context) PushStyleColor(idx Col, col color.RGBA) {
	c.ColorModifiers = append(c.ColorModifiers, ColMod{idx, c.Style.Colors[idx]})
	c.Style.Colors[idx] = chroma.RGBA2VEC4(col)
}

func (c *Context) PopStyleColor() {
	c.PopStyleColorN(1)
}

func (c *Context) PopStyleColorN(count int) {
	for ; count > 0 && len(c.ColorModifiers) > 0; count-- {
		backup := c.ColorModifiers[len(c.ColorModifiers)-1]
		c.Style.Colors[backup.Col] = backup.BackupValue
		c.ColorModifiers = c.ColorModifiers[:len(c.ColorModifiers)-1]
	}
}

func (c *Context) GetStyleColorName(idx Col) string {
	switch idx {
	case ColText:
		return "Text"
	case ColTextDisabled:
		return "TextDisabled"
	case ColWindowBg:
		return "WindowBg"
	case ColChildBg:
		return "ChildBg"
	case ColPopupBg:
		return "PopupBg"
	case ColBorder:
		return "Border"
	case ColBorderShadow:
		return "BorderShadow"
	case ColFrameBg:
		return "FrameBg"
	case ColFrameBgHovered:
		return "FrameBgHovered"
	case ColFrameBgActive:
		return "FrameBgActive"
	case ColTitleBg:
		return "TitleBg"
	case ColTitleBgActive:
		return "TitleBgActive"
	case ColTitleBgCollapsed:
		return "TitleBgCollapsed"
	case ColMenuBarBg:
		return "MenuBarBg"
	case ColScrollbarBg:
		return "ScrollbarBg"
	case ColScrollbarGrab:
		return "ScrollbarGrab"
	case ColScrollbarGrabHovered:
		return "ScrollbarGrabHovered"
	case ColScrollbarGrabActive:
		return "ScrollbarGrabActive"
	case ColCheckMark:
		return "CheckMark"
	case ColSliderGrab:
		return "SliderGrab"
	case ColSliderGrabActive:
		return "SliderGrabActive"
	case ColButton:
		return "Button"
	case ColButtonHovered:
		return "ButtonHovered"
	case ColButtonActive:
		return "ButtonActive"
	case ColHeader:
		return "Header"
	case ColHeaderHovered:
		return "HeaderHovered"
	case ColHeaderActive:
		return "HeaderActive"
	case ColSeparator:
		return "Separator"
	case ColSeparatorHovered:
		return "SeparatorHovered"
	case ColSeparatorActive:
		return "SeparatorActive"
	case ColResizeGrip:
		return "ResizeGrip"
	case ColResizeGripHovered:
		return "ResizeGripHovered"
	case ColResizeGripActive:
		return "ResizeGripActive"
	case ColPlotLines:
		return "PlotLines"
	case ColPlotLinesHovered:
		return "PlotLinesHovered"
	case ColPlotHistogram:
		return "PlotHistogram"
	case ColPlotHistogramHovered:
		return "PlotHistogramHovered"
	case ColTextSelectedBg:
		return "TextSelectedBg"
	case ColModalWindowDarkening:
		return "ModalWindowDarkening"
	case ColDragDropTarget:
		return "DragDropTarget"
	case ColNavHighlight:
		return "NavHighlight"
	case ColNavWindowingHighlight:
		return "NavWindowingHighlight"
	}
	assert(false)
	return "Unknown"
}
