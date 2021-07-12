package imgui

import (
	"fmt"
	"os"
	"strings"

	"github.com/qeedquan/go-media/math/f64"
)

// Pass text data straight to log (without being displayed)
func (c *Context) LogText(format string, args ...interface{}) {
	if !c.LogEnabled {
		return
	}

	text := fmt.Sprintf(format, args...)
	if c.LogFile != nil {
		fmt.Fprint(c.LogFile, text)
	} else {
		c.LogClipboard = append(c.LogClipboard, []rune(text)...)
	}
}

// Internal version that takes a position to decide on newline placement and pad items according to their depth.
// We split text into individual lines to add current tree level padding
func (c *Context) LogRenderedText(ref_pos *f64.Vec2, text string) {
	window := c.CurrentWindow

	log_new_line := ref_pos != nil && (ref_pos.Y > window.DC.LogLinePosY+1)
	if ref_pos != nil {
		window.DC.LogLinePosY = ref_pos.Y
	}

	// Re-adjust padding if we have popped out of our starting depth
	if c.LogStartDepth > window.DC.TreeDepth {
		c.LogStartDepth = window.DC.TreeDepth
	}
	tree_depth := (window.DC.TreeDepth - c.LogStartDepth)
	lines := strings.Split(text, "\n")
	for i := range lines {
		if log_new_line || i > 0 {
			c.LogText("%.*s%s", tree_depth*4, "", lines[i])
		} else {
			c.LogText(" %s", lines[i])
		}
	}
}

func (c *Context) LogToClipboard() {
	c.LogToClipboardEx(-1)
}

// Start logging ImGui output to clipboard
func (c *Context) LogToClipboardEx(max_depth int) {
	if c.LogEnabled {
		return
	}
	window := c.CurrentWindow

	assert(c.LogFile == nil)
	c.LogFile = nil
	c.LogEnabled = true
	c.LogStartDepth = window.DC.TreeDepth
	if max_depth >= 0 {
		c.LogAutoExpandMaxDepth = max_depth
	}
}

func (c *Context) LogToTTY() {
	c.LogToTTYEx(-1)
}

func (c *Context) LogToTTYEx(max_depth int) {
	if c.LogEnabled {
		return
	}
	window := c.CurrentWindow

	assert(c.LogFile == nil)
	c.LogFile = os.Stdout
	c.LogEnabled = true
	c.LogStartDepth = window.DC.TreeDepth
	if max_depth >= 0 {
		c.LogAutoExpandMaxDepth = max_depth
	}
}

// Start logging ImGui output to given file
func (c *Context) LogToFile(max_depth int, filename string) {
	if c.LogEnabled {
		return
	}

	window := c.CurrentWindow
	if filename == "" {
		filename = c.IO.LogFilename
	}
	if filename == "" {
		return
	}

	assert(c.LogFile == nil)
	var err error
	c.LogFile, err = os.OpenFile(filename, os.O_APPEND, 0644)
	if err != nil {
		// Consider this an error
		assert(c.LogFile != nil)
		return
	}
	c.LogEnabled = true
	c.LogStartDepth = window.DC.TreeDepth
	if max_depth >= 0 {
		c.LogAutoExpandMaxDepth = max_depth
	}
}

func (c *Context) LogFinish() {
	if !c.LogEnabled {
		return
	}

	c.LogText("\n")
	if c.LogFile != nil {
		c.LogFile.Close()
		c.LogFile = nil
	}
	if len(c.LogClipboard) > 0 {
		c.SetClipboardText(string(c.LogClipboard))
		c.LogClipboard = c.LogClipboard[:0]
	}
	c.LogEnabled = false
}
