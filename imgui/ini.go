package imgui

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/qeedquan/go-media/math/f64"
)

type SettingsHandler struct {
	TypeName   string // Short description stored in .ini file. Disallowed characters: '[' ']'
	ReadOpenFn func(ctx *Context, handler *SettingsHandler, name string) interface{}
	ReadLineFn func(ctx *Context, handler *SettingsHandler, entry interface{}, line string)
	WriteAllFn func(ctx *Context, handler *SettingsHandler, out_buf io.Writer)
}

func (c *Context) LoadIniSettingsFromDisk(filename string) error {
	file_data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	c.LoadIniSettingsFromMemory(file_data)
	return nil
}

// Zero-tolerance, no error reporting, cheap .ini parsing
func (c *Context) LoadIniSettingsFromMemory(buf []byte) error {
	var (
		entry_handler *SettingsHandler
		entry_data    interface{}
	)

	scan := bufio.NewScanner(bytes.NewBuffer(buf))
	for scan.Scan() {
		line := scan.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			// Parse "[Type][Name]". Note that 'Name' can itself contains [] characters, which is acceptable with the current format and parsing code.

			n := strings.IndexRune(line, ']')
			typ := line[1:n]
			name := line[n+1:]

			n = strings.IndexRune(line, '[')
			if n < 0 {
				name = typ
				typ = "Window"
			} else {
				name = strings.TrimPrefix(name, "[")
				name = strings.TrimSuffix(name, "]")
			}

			entry_handler = c.FindSettingsHandler(typ)
			if entry_handler != nil {
				entry_data = entry_handler.ReadOpenFn(c, entry_handler, name)
			}
		} else if entry_handler != nil {
			// Let type handler parse the line
			entry_handler.ReadLineFn(c, entry_handler, entry_data, line)
		}
	}

	c.SettingsLoaded = true
	return nil
}

func (c *Context) FindSettingsHandler(typ string) *SettingsHandler {
	return c.SettingsHandlers[typ]
}

func (c *Context) SaveIniSettingsToDisk(ini_filename string) error {
	buf := c.SaveIniSettingsToMemory()

	f, err := os.Create(ini_filename)
	if err != nil {
		return err
	}

	f.Write(buf)
	return f.Close()
}

func (c *Context) SaveIniSettingsToMemory() []byte {
	w := new(bytes.Buffer)
	c.SettingsDirtyTimer = 0.0
	for i := range c.SettingsHandlers {
		c.SettingsHandlers[i].WriteAllFn(c, c.SettingsHandlers[i], w)
	}
	return w.Bytes()
}

func (c *Context) MarkIniSettingsDirty() {
	if c.SettingsDirtyTimer <= 0 {
		c.SettingsDirtyTimer = c.IO.IniSavingRate
	}
}

func (c *Context) MarkIniSettingsDirtyForWindow(window *Window) {
	if window.Flags&WindowFlagsNoSavedSettings == 0 {
		if c.SettingsDirtyTimer <= 0 {
			c.SettingsDirtyTimer = c.IO.IniSavingRate
		}
	}
}

func (c *Context) SettingsHandlerWindow_ReadOpen(_ *Context, _ *SettingsHandler, name string) interface{} {
	settings := c.FindWindowSettings(name)
	if settings == nil {
		settings = c.AddWindowSettings(name)
	}
	return settings
}

func (c *Context) SettingsHandlerWindow_ReadLine(_ *Context, _ *SettingsHandler, entry interface{}, line string) {
	settings := entry.(*WindowSettings)

	var x, y float64
	n, _ := fmt.Sscanf(line, "Pos=%f,%f", &x, &y)
	if n == 2 {
		settings.Pos = f64.Vec2{x, y}
		return
	}

	n, _ = fmt.Sscanf(line, "Size=%f,%f", &x, &y)
	if n == 2 {
		settings.Size = f64.Vec2{x, y}.Max(c.Style.WindowMinSize)
		return
	}

	n, _ = fmt.Sscanf(line, "Collapsed=%d", &x)
	if n == 1 {
		settings.Collapsed = x != 0
		return
	}
}

func (c *Context) SettingsHandlerWindow_WriteAll(ctx *Context, handler *SettingsHandler, w io.Writer) {
}
