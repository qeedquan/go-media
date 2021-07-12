package imgui

func (c *Context) SetClipboardText(text string) {
	if c.IO.SetClipboardTextFn != nil {
		c.IO.SetClipboardTextFn(text)
	}
}

func (c *Context) GetClipboardText() string {
	if c.IO.GetClipboardTextFn != nil {
		return c.IO.GetClipboardTextFn()
	}
	return ""
}