package imgui

import "fmt"

func (c *Context) BeginTooltip() {
	c.BeginTooltipEx(0, false)
}

func (c *Context) BeginTooltipEx(extra_flags WindowFlags, override_previous_tooltip bool) {
	window_name := fmt.Sprintf("##Tooltip_%02d", c.TooltipOverrideCount)
	if override_previous_tooltip {
		window := c.FindWindowByName(window_name)
		if window != nil && window.Active {
			// Hide previous tooltips. We can't easily "reset" the content of a window so we create a new one.
			window.HiddenFrames = 1
			c.TooltipOverrideCount++
			window_name = fmt.Sprintf("##Tooltip_%02d", c.TooltipOverrideCount)
		}
	}
	flags := WindowFlagsTooltip | WindowFlagsNoInputs | WindowFlagsNoTitleBar |
		WindowFlagsNoMove | WindowFlagsNoResize | WindowFlagsNoSavedSettings |
		WindowFlagsAlwaysAutoResize | WindowFlagsNoNav
	c.BeginEx(window_name, nil, flags|extra_flags)
}

func (c *Context) SetTooltip(format string, args ...interface{}) {
	text := fmt.Sprintf(format, args...)
	c.BeginTooltipEx(0, true)
	c.Text(text)
	c.EndTooltip()
}

func (c *Context) EndTooltip() {
	// Mismatched BeginTooltip()/EndTooltip() calls
	assert(c.GetCurrentWindowRead().Flags&WindowFlagsTooltip != 0)
	c.End()
}