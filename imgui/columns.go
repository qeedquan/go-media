package imgui

import (
	"image/color"
	"math"

	"github.com/qeedquan/go-media/math/f64"
)

type ColumnData struct {
	OffsetNorm             float64 // Column start offset, normalized 0.0 (far left) -> 1.0 (far right)
	OffsetNormBeforeResize float64
	Flags                  ColumnsFlags // Not exposed
	ClipRect               f64.Rectangle
}

type ColumnsSet struct {
	ID                 ID
	Flags              ColumnsFlags
	IsFirstFrame       bool
	IsBeingResized     bool
	Current            int
	Count              int
	MinX, MaxX         float64
	LineMinY, LineMaxY float64
	StartPosY          float64
	StartMaxPosX       float64 // Backup of CursorMaxPos
	CellMinY, CellMaxY float64
	Columns            []ColumnData
}

type MenuColumns struct {
	Count            int
	Spacing          float64
	Width, NextWidth float64
	Pos, NextWidths  [4]float64
}

type ColumnsFlags int

const (
	// Default: 0
	ColumnsFlagsNoBorder               ColumnsFlags = 1 << 0 // Disable column dividers
	ColumnsFlagsNoResize               ColumnsFlags = 1 << 1 // Disable resizing columns when clicking on the dividers
	ColumnsFlagsNoPreserveWidths       ColumnsFlags = 1 << 2 // Disable column width preservation when adjusting columns
	ColumnsFlagsNoForceWithinWindow    ColumnsFlags = 1 << 3 // Disable forcing columns to fit within window
	ColumnsFlagsGrowParentContentsSize ColumnsFlags = 1 << 4 // (WIP) Restore pre-1.51 behavior of extending the parent window contents size but _without affecting the columns width at all_. Will eventually remove.
)

// [2018-03: This is currently the only public API, while we are working on making BeginColumns/EndColumns user-facing]
func (c *Context) Columns(columns_count int, id string, border bool) {
	window := c.GetCurrentWindow()
	assert(columns_count >= 1)

	var flags ColumnsFlags
	if border {
		flags |= ColumnsFlagsNoBorder
	}
	if window.DC.ColumnsSet != nil && window.DC.ColumnsSet.Count == columns_count && window.DC.ColumnsSet.Flags == flags {
		return
	}

	if window.DC.ColumnsSet != nil {
		c.EndColumns()
	}

	if columns_count != 1 {
		c.BeginColumns(id, columns_count, flags)
	}
}

func (c *Context) BeginColumns(str_id string, columns_count int, flags ColumnsFlags) {
	window := c.GetCurrentWindow()

	assert(columns_count > 1)
	assert(window.DC.ColumnsSet == nil) // Nested columns are currently not supported

	// Differentiate column ID with an arbitrary prefix for cases where users name their columns set the same as another widget.
	// In addition, when an identifier isn't explicitly provided we include the number of columns in the hash to make it uniquer.
	column_id := ID(0x11223347)
	if str_id != "" {
		column_id += ID(columns_count)
	}
	c.PushID(column_id)

	var id ID
	if str_id != "" {
		id = window.GetID(str_id)
	} else {
		id = window.GetID("columns")
	}
	c.PopID()

	// Acquire storage for the columns set
	columns := c.FindOrAddColumnsSet(window, id)
	assert(columns.ID == id)
	columns.Current = 0
	columns.Count = columns_count
	columns.Flags = flags
	window.DC.ColumnsSet = columns

	// Set state for first column
	content_region_width := window.InnerClipRect.Max.X - window.Pos.X
	if window.SizeContentsExplicit.X != 0.0 {
		content_region_width = window.SizeContentsExplicit.X
	}
	// Lock our horizontal range
	columns.MinX = window.DC.IndentX - c.Style.ItemSpacing.X
	columns.MaxX = math.Max(content_region_width-window.Scroll.X, columns.MinX+1.0)
	columns.StartPosY = window.DC.CursorPos.Y
	columns.StartMaxPosX = window.DC.CursorMaxPos.X
	columns.LineMinY = window.DC.CursorPos.Y
	columns.LineMaxY = columns.LineMinY
	window.DC.ColumnsOffsetX = 0.0
	window.DC.CursorPos.X = float64(int(window.Pos.X + window.DC.IndentX + window.DC.ColumnsOffsetX))

	// Clear data if columns count changed
	if len(columns.Columns) != 0 && len(columns.Columns) != columns_count+1 {
		columns.Columns = columns.Columns[:0]
	}

	// Initialize defaults
	columns.IsFirstFrame = len(columns.Columns) == 0
	if len(columns.Columns) == 0 {
		for n := 0; n < columns_count+1; n++ {
			var column ColumnData
			column.Init()
			column.OffsetNorm = float64(n) / float64(columns_count)
			columns.Columns = append(columns.Columns, column)
		}
	}

	for n := 0; n < columns_count; n++ {
		// Compute clipping rectangle
		column := &columns.Columns[n]
		clip_x1 := math.Floor(0.5 + window.Pos.X + c.GetColumnOffset(n) - 1.0)
		clip_x2 := math.Floor(0.5 + window.Pos.X + c.GetColumnOffset(n+1) - 1.0)
		column.ClipRect = f64.Rect(clip_x1, -math.MaxFloat32, clip_x2, +math.MaxFloat32)
		column.ClipRect = column.ClipRect.Intersect(window.ClipRect)
	}

	window.DrawList.ChannelsSplit(columns.Count)
	c.PushColumnClipRect()
	c.PushItemWidth(c.GetColumnWidth() * 0.65)
}

func (c *Context) EndColumns() {
	window := c.GetCurrentWindow()
	columns := window.DC.ColumnsSet

	c.PopItemWidth()
	c.PopClipRect()
	window.DrawList.ChannelsMerge()

	columns.CellMaxY = math.Max(columns.CellMaxY, window.DC.CursorPos.Y)
	window.DC.CursorPos.Y = columns.CellMaxY
	if columns.Flags&ColumnsFlagsGrowParentContentsSize == 0 {
		// Restore cursor max pos, as columns don't grow parent
		window.DC.CursorMaxPos.X = columns.StartMaxPosX
	}

	// Draw columns borders and handle resize
	is_being_resized := false
	if columns.Flags&ColumnsFlagsNoBorder == 0 && !window.SkipItems {
		y1 := columns.StartPosY
		y2 := window.DC.CursorPos.Y
		dragging_column := -1
		for n := 1; n < columns.Count; n++ {
			x := window.Pos.X + c.GetColumnOffset(n)
			column_id := columns.ID + ID(n)
			column_hw := c.GetColumnsRectHalfWidth() // Half-width for interaction
			column_rect := f64.Rectangle{f64.Vec2{x - column_hw, y1}, f64.Vec2{x + column_hw, y2}}
			c.KeepAliveID(column_id)
			if c.IsClippedEx(column_rect, column_id, false) {
				continue
			}

			var hovered, held bool
			if columns.Flags&ColumnsFlagsNoResize == 0 {
				hovered, held, _ = c.ButtonBehavior(column_rect, column_id, 0)
				if hovered || held {
					c.MouseCursor = MouseCursorResizeEW
				}
				if held && columns.Columns[n].Flags&ColumnsFlagsNoResize == 0 {
					dragging_column = n
				}
			}

			// Draw column (we clip the Y boundaries CPU side because very long triangles are mishandled by some GPU drivers.)
			var col color.RGBA
			switch {
			case held:
				col = c.GetColorFromStyle(ColSeparatorActive)
			case hovered:
				col = c.GetColorFromStyle(ColSeparatorHovered)
			default:
				col = c.GetColorFromStyle(ColSeparator)
			}
			xi := float64(int(x))
			window.DrawList.AddLine(
				f64.Vec2{xi, math.Max(y1+1, window.ClipRect.Min.Y)},
				f64.Vec2{xi, math.Min(y2, window.ClipRect.Max.Y)},
				col,
			)
		}

		// Apply dragging after drawing the column lines, so our rendered lines are in sync
		// with how items were displayed during the frame.
		if dragging_column != -1 {
			if !columns.IsBeingResized {
				for n := 0; n < columns.Count+1; n++ {
					columns.Columns[n].OffsetNormBeforeResize = columns.Columns[n].OffsetNorm
				}
				columns.IsBeingResized = true
				is_being_resized = columns.IsBeingResized
				x := c.GetDraggedColumnOffset(columns, dragging_column)
				c.SetColumnOffset(dragging_column, x)
			}
		}
	}

	columns.IsBeingResized = is_being_resized

	window.DC.ColumnsSet = nil
	window.DC.ColumnsOffsetX = 0
	window.DC.CursorPos.X = float64(int(window.Pos.X + window.DC.IndentX + window.DC.ColumnsOffsetX))

}

func (c *ColumnData) Init() {
	c.OffsetNorm = 0
	c.OffsetNormBeforeResize = 0.0
	c.Flags = 0
}

func (c *Context) FindOrAddColumnsSet(window *Window, id ID) *ColumnsSet {
	for n := range window.ColumnsStorage {
		if window.ColumnsStorage[n].ID == id {
			return &window.ColumnsStorage[n]
		}
	}

	var columns_set ColumnsSet
	columns_set.Init()
	window.ColumnsStorage = append(window.ColumnsStorage, columns_set)
	columns := &window.ColumnsStorage[len(window.ColumnsStorage)-1]
	columns.ID = id
	return columns
}

func (c *ColumnsSet) Init() {
	c.Clear()
}

func (c *ColumnsSet) Clear() {
	c.ID = 0
	c.Flags = 0
	c.IsFirstFrame = false
	c.IsBeingResized = false
	c.Current = 0
	c.Count = 1
	c.MinX = 0.0
	c.MaxX = 0.0
	c.LineMinY = 0.0
	c.LineMaxY = 0.0
	c.StartPosY = 0.0
	c.StartMaxPosX = 0.0
	c.Columns = c.Columns[:0]
}

func (m *MenuColumns) DeclColumns(w0, w1, w2 float64) float64 {
	m.NextWidth = 0.0
	m.NextWidths[0] = math.Max(m.NextWidths[0], w0)
	m.NextWidths[1] = math.Max(m.NextWidths[1], w1)
	m.NextWidths[2] = math.Max(m.NextWidths[2], w2)
	for i := 0; i < 3; i++ {
		m.NextWidth += m.NextWidths[i]
		if i > 0 && m.NextWidths[i] > 0.0 {
			m.NextWidth += m.Spacing
		}
	}
	return math.Max(m.Width, m.NextWidth)
}

func (m *MenuColumns) Update(count int, spacing float64, clear bool) {
	m.Count = count
	m.Width = 0
	m.NextWidth = 0
	m.Spacing = spacing
	if clear {
		for i := range m.NextWidths {
			m.NextWidths[i] = 0
		}
	}
	for i := 0; i < m.Count; i++ {
		if i > 0 && m.NextWidths[i] > 0 {
			m.Width += m.Spacing
		}
		m.Pos[i] = float64(int(m.Width))
		m.Width += m.NextWidths[i]
		m.NextWidths[i] = 0
	}
}
