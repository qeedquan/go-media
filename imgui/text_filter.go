package imgui

import (
	"strings"
)

// Helper: Parse and apply text filters. In format "aaaaa[,bbbb][,ccccc]"
type TextFilter struct {
	Ctx       *Context
	InputBuf  []byte
	Filters   []string
	CountGrep int
}

func (t *TextFilter) Init(ctx *Context, default_filter string) {
	t.Ctx = ctx
	if default_filter != "" {
		t.InputBuf = []byte(default_filter)
		t.Build()
	} else {
		t.InputBuf = t.InputBuf[:0]
		t.CountGrep = 0
	}
}

func (t *TextFilter) Build() {
	t.Filters = strings.Split(string(t.InputBuf), ",")

	t.CountGrep = 0
	for i := range t.Filters {
		t.Filters[i] = strings.TrimSpace(t.Filters[i])
		if t.Filters[i] == "" {
			continue
		}
		if !strings.HasPrefix(t.Filters[i], "-") {
			t.CountGrep += 1
		}
	}
}

func (t *TextFilter) PassFilter(text string) bool {
	if len(t.Filters) == 0 {
		return true
	}

	for i := range t.Filters {
		f := t.Filters[i]
		if len(f) == 0 {
			continue
		}

		a := strings.ToLower(text)
		b := strings.ToLower(f)
		if strings.HasPrefix(f, "-") {
			// Subtract
			if strings.Index(a, b) >= 0 {
				return false
			}
		} else {
			// Grep
			if strings.Index(a, b) >= 0 {
				return true
			}
		}
	}

	// Implicit * grep
	if t.CountGrep == 0 {
		return true
	}
	return false
}

func (t *TextFilter) IsActive() bool {
	return len(t.Filters) > 0
}

func (t *TextFilter) Draw() bool {
	return t.DrawEx("Filter (inc,-exc)", 0)
}

func (t *TextFilter) DrawEx(label string, width float64) bool {
	c := t.Ctx
	if width != 0.0 {
		c.PushItemWidth(width)
	}
	value_changed := c.InputText(label, t.InputBuf)
	if width != 0.0 {
		c.PopItemWidth()
	}
	if value_changed {
		t.Build()
	}
	return value_changed
}