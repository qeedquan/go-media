package imgui

import (
	"image/color"
	"math"

	"github.com/qeedquan/go-media/math/f64"
	"github.com/qeedquan/go-media/math/mathutil"
)

type PlotType int

const (
	PlotTypeLines = iota
	PlotTypeHistogram
)

func (c *Context) PlotLines(label string, values []float64) {
	c.PlotLinesEx(label, values, 0, "", math.MaxFloat32, math.MaxFloat32, f64.Vec2{0, 0})
}

func (c *Context) PlotLinesEx(label string, values []float64, values_offset int, overlay_text string, scale_min, scale_max float64, graph_size f64.Vec2) {
	values_getter := func(idx int) float64 {
		if 0 <= idx && idx < len(values) {
			return values[idx]
		}
		return 0
	}

	c.PlotEx(PlotTypeLines, label, values_getter, len(values), values_offset, overlay_text, scale_min, scale_max, graph_size)
}

func (c *Context) PlotLinesItem(label string, values_getter func(idx int) float64, values_count int) {
	c.PlotLinesItemEx(label, values_getter, values_count, 0, "", math.MaxFloat32, math.MaxFloat32, f64.Vec2{0, 0})
}

func (c *Context) PlotLinesItemEx(label string, values_getter func(idx int) float64, values_count, values_offset int, overlay_text string, scale_min, scale_max float64, graph_size f64.Vec2) {
	c.PlotEx(PlotTypeLines, label, values_getter, values_count, values_offset, overlay_text, scale_min, scale_max, graph_size)
}

func (c *Context) PlotHistogram(label string, values []float64) {
	c.PlotHistogramEx(label, values, 0, "", math.MaxFloat32, math.MaxFloat32, f64.Vec2{0, 0})
}

func (c *Context) PlotHistogramEx(label string, values []float64, values_offset int, overlay_text string, scale_min, scale_max float64, graph_size f64.Vec2) {
	values_getter := func(idx int) float64 {
		if 0 <= idx && idx < len(values) {
			return values[idx]
		}
		return 0
	}
	c.PlotEx(PlotTypeHistogram, label, values_getter, len(values), values_offset, overlay_text, scale_min, scale_max, graph_size)
}

func (c *Context) PlotHistogramItem(label string, values_getter func(idx int) float64, values_count int) {
	c.PlotHistogramItemEx(label, values_getter, values_count, 0, "", math.MaxFloat32, math.MaxFloat32, f64.Vec2{0, 0})
}

func (c *Context) PlotHistogramItemEx(label string, values_getter func(idx int) float64, values_count, values_offset int, overlay_text string, scale_min, scale_max float64, graph_size f64.Vec2) {
	c.PlotEx(PlotTypeHistogram, label, values_getter, values_count, values_offset, overlay_text, scale_min, scale_max, graph_size)
}

func (c *Context) PlotEx(plot_type PlotType, label string, values_getter func(idx int) float64, values_count, values_offset int, overlay_text string, scale_min, scale_max float64, graph_size f64.Vec2) {
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return
	}

	style := &c.Style
	label_size := c.CalcTextSizeEx(label, true, -1)
	if graph_size.X == 0.0 {
		graph_size.X = c.CalcItemWidth()
	}
	if graph_size.Y == 0.0 {
		graph_size.Y = label_size.Y + (style.FramePadding.Y * 2)
	}

	frame_bb := f64.Rectangle{
		window.DC.CursorPos,
		window.DC.CursorPos.Add(f64.Vec2{graph_size.X, graph_size.Y}),
	}
	inner_bb := f64.Rectangle{
		frame_bb.Min.Add(style.FramePadding),
		frame_bb.Max.Sub(style.FramePadding),
	}
	total_bb_spacing := f64.Vec2{0, 0}
	if label_size.X > 0.0 {
		total_bb_spacing.X = style.ItemInnerSpacing.X + label_size.X
	}
	total_bb := f64.Rectangle{
		frame_bb.Min,
		frame_bb.Max.Add(total_bb_spacing),
	}
	c.ItemSizeBBEx(total_bb, style.FramePadding.Y)
	if !c.ItemAddEx(total_bb, 0, &frame_bb) {
		return
	}

	hovered := c.ItemHoverable(inner_bb, 0)

	// Determine scale from values if not specified
	if scale_min == math.MaxFloat32 || scale_max == math.MaxFloat32 {
		v_min := math.MaxFloat32
		v_max := -math.MaxFloat32
		for i := 0; i < values_count; i++ {
			v := values_getter(i)
			v_min = math.Min(v_min, v)
			v_max = math.Max(v_max, v)
		}
		if scale_min == math.MaxFloat32 {
			scale_min = v_min
		}
		if scale_max == math.MaxFloat32 {
			scale_max = v_max
		}
	}

	c.RenderFrameEx(frame_bb.Min, frame_bb.Max, c.GetColorFromStyle(ColFrameBg), true, style.FrameRounding)

	if values_count > 0 {
		res_w := mathutil.Min(int(graph_size.X), values_count)
		item_count := values_count
		if plot_type == PlotTypeLines {
			res_w += -1
			item_count += -1
		}

		// Tooltip on hover
		v_hovered := -1
		if hovered {
			t := f64.Clamp((c.IO.MousePos.X-inner_bb.Min.X)/(inner_bb.Max.X-inner_bb.Min.X), 0.0, 0.9999)
			v_idx := int(t * float64(item_count))
			assert(v_idx >= 0 && v_idx < values_count)

			v0 := values_getter((v_idx + values_offset) % values_count)
			v1 := values_getter((v_idx + 1 + values_offset) % values_count)
			if plot_type == PlotTypeLines {
				c.SetTooltip("%d: %8.4g\n%d: %8.4g", v_idx, v0, v_idx+1, v1)
			} else if plot_type == PlotTypeHistogram {
				c.SetTooltip("%d: %8.4g", v_idx, v0)
			}
			v_hovered = v_idx
		}

		t_step := 1.0 / float64(res_w)
		inv_scale := 0.0
		if scale_min != scale_max {
			inv_scale = 1.0 / (scale_max - scale_min)
		}

		v0 := values_getter((0 + values_offset) % values_count)
		t0 := 0.0
		// Point in the normalized space of our target rectangle
		tp0 := f64.Vec2{t0, 1.0 - f64.Saturate((v0-scale_min)*inv_scale)}

		// Where does the zero line stands
		var histogram_zero_line_t float64
		if scale_min*scale_max < 0.0 {
			histogram_zero_line_t = -scale_min * inv_scale
		} else if scale_min < 0.0 {
			histogram_zero_line_t = 0.0
		} else {
			histogram_zero_line_t = 1.0
		}

		var col_base, col_hovered color.RGBA
		if plot_type == PlotTypeLines {
			col_base = c.GetColorFromStyle(ColPlotLines)
			col_hovered = c.GetColorFromStyle(ColPlotLines)
		} else {
			col_base = c.GetColorFromStyle(ColPlotHistogram)
			col_hovered = c.GetColorFromStyle(ColPlotHistogram)
		}

		for n := 0; n < res_w; n++ {
			t1 := t0 + t_step
			v1_idx := int(t0*float64(item_count) + 0.5)
			assert(v1_idx >= 0 && v1_idx < values_count)
			v1 := values_getter((v1_idx + values_offset + 1) % values_count)
			tp1 := f64.Vec2{t1, 1.0 - f64.Saturate((v1-scale_min)*inv_scale)}

			// NB: Draw calls are merged together by the DrawList system. Still, we should render our batch are lower level to save a bit of CPU.
			pos0 := inner_bb.Min.Lerp2(tp0, inner_bb.Max)
			lp1 := tp1
			if plot_type == PlotTypeLines {
				lp1 = f64.Vec2{tp1.X, histogram_zero_line_t}
			}
			pos1 := inner_bb.Min.Lerp2(lp1, inner_bb.Max)

			col := col_base
			if v_hovered == v1_idx {
				col = col_hovered
			}
			if plot_type == PlotTypeLines {
				window.DrawList.AddLine(pos0, pos1, col)
			} else if plot_type == PlotTypeHistogram {
				if pos1.X >= pos0.X+2.0 {
					pos1.X -= 1.0
				}
				window.DrawList.AddRectFilled(pos0, pos1, col)
			}
			t0 = t1
			tp0 = tp1
		}
	}

	// Text overlay
	if overlay_text != "" {
		c.RenderTextClippedEx(f64.Vec2{frame_bb.Min.X, frame_bb.Min.Y + style.FramePadding.Y}, frame_bb.Max, overlay_text, nil, f64.Vec2{0.5, 0.0}, nil)
	}

	if label_size.X > 0.0 {
		c.RenderText(f64.Vec2{frame_bb.Max.X + style.ItemInnerSpacing.X, inner_bb.Min.Y}, label)
	}
}