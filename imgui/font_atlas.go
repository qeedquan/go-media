package imgui

import (
	"math"

	"github.com/qeedquan/go-media/math/f64"
	"github.com/qeedquan/go-media/math/mathutil"
	"github.com/qeedquan/go-media/stb/stbtt"
)

type CustomRect struct {
	ID            uint     // Input    // User ID. Use <0x10000 to map into a font glyph, >=0x10000 for other/internal/custom texture data.
	Width, Height uint     // Input    // Desired rectangle dimension
	X, Y          uint     // Output   // Packed position in Atlas
	GlyphAdvanceX float64  // Input    // For custom font glyphs only (ID<0x10000): glyph xadvance
	GlyphOffset   f64.Vec2 // Input    // For custom font glyphs only (ID<0x10000): glyph display offset
	Font          *Font    // Input    // For custom font glyphs only (ID<0x10000): target font
}

type FontAtlasFlags int

const (
	FontAtlasFlagsNoPowerOfTwoHeight FontAtlasFlags = 1 << 0 // Don't round the height to next power of two
	FontAtlasFlagsNoMouseCursors     FontAtlasFlags = 1 << 1 // Don't build software mouse cursors into the atlas
)

type FontAtlas struct {
	Flags           FontAtlasFlags // Build flags (see ImFontAtlasFlags_)
	TexID           TextureID      // User data to refer to the texture once it has been uploaded to user's graphic systems. It is passed back to you during rendering via the ImDrawCmd structure.
	TexDesiredWidth int            // Texture width desired by user before Build(). Must be a power-of-two. If have many glyphs your graphics API have texture size restrictions you may want to increase texture width to decrease height.
	TexGlyphPadding int            // Padding between glyphs within texture in pixels. Defaults to 1.

	// [Internal]
	// NB: Access texture data via GetTexData*() calls! Which will setup a default font for you.
	TexPixelsAlpha8 []uint8      // 1 component per pixel, each component is unsigned 8-bit. Total size = TexWidth * TexHeight
	TexPixelsRGBA32 []uint8      // 4 component per pixel, each component is unsigned 8-bit. Total size = TexWidth * TexHeight * 4
	TexWidth        int          // Texture width calculated during Build().
	TexHeight       int          // Texture height calculated during Build().
	TexUvScale      f64.Vec2     // = (1.0f/TexWidth, 1.0f/TexHeight)
	TexUvWhitePixel f64.Vec2     // Texture coordinates to a white pixel
	Fonts           []*Font      // Hold all the fonts returned by AddFont*. Fonts[0] is the default font upon calling ImGui::NewFrame(), use ImGui::PushFont()/PopFont() to change the current font.
	CustomRects     []CustomRect // Rectangles for packing custom texture data into the atlas.
	ConfigData      []FontConfig // Internal data
	CustomRectIds   [1]int       // Identifiers of custom texture rectangle used by ImFontAtlas/ImDrawList
}

func (c *CustomRect) Init() {
	c.ID = 0xFFFFFFFF
	c.Width = 0
	c.Height = 0
	c.X = 0xFFFF
	c.Y = 0xFFFF
	c.GlyphAdvanceX = 0.0
	c.GlyphOffset = f64.Vec2{0, 0}
	c.Font = nil
}

func (c *CustomRect) IsPacked() bool {
	return c.X != 0xFFFF
}

func NewFontAtlas() *FontAtlas {
	f := &FontAtlas{}
	f.Init()
	return f
}

func (f *FontAtlas) Init() {
	f.Flags = 0x0
	f.TexID = nil
	f.TexDesiredWidth = 0
	f.TexGlyphPadding = 1

	f.TexPixelsAlpha8 = nil
	f.TexPixelsRGBA32 = nil
	f.TexWidth = 0
	f.TexHeight = 0
	f.TexUvScale = f64.Vec2{0, 0}
	f.TexUvWhitePixel = f64.Vec2{0, 0}
	for n := range f.CustomRectIds {
		f.CustomRectIds[n] = -1
	}
}

func (f *FontAtlas) GetMouseCursorTexData(cursor_type MouseCursor, out_offset, out_size *f64.Vec2, out_uv_border, out_uv_fill []f64.Vec2) bool {
	if cursor_type <= MouseCursorNone || cursor_type >= MouseCursorCOUNT {
		return false
	}

	if f.Flags&FontAtlasFlagsNoMouseCursors != 0 {
		return false
	}

	assert(f.CustomRectIds[0] != -1)
	r := &f.CustomRects[f.CustomRectIds[0]]
	assert(r.ID == FONT_ATLAS_DEFAULT_TEX_DATA_ID)
	pos := FONT_ATLAS_DEFAULT_TEX_CURSOR_DATA[cursor_type][0]
	pos = pos.Add(f64.Vec2{float64(r.X), float64(r.Y)})
	size := FONT_ATLAS_DEFAULT_TEX_CURSOR_DATA[cursor_type][1]
	*out_size = size
	*out_offset = FONT_ATLAS_DEFAULT_TEX_CURSOR_DATA[cursor_type][2]
	out_uv_border[0] = pos.Scale2(f.TexUvScale)
	out_uv_border[1] = pos.Add(size).Scale2(f.TexUvScale)
	pos.X += FONT_ATLAS_DEFAULT_TEX_DATA_W_HALF + 1
	out_uv_fill[0] = pos.Scale2(f.TexUvScale)
	out_uv_fill[1] = pos.Add(size).Scale2(f.TexUvScale)
	return true
}

func (f *FontAtlas) GetTexDataAsAlpha8() (out_pixels []byte, out_width, out_height, out_bytes_per_pixel int) {
	// Build atlas on demand
	if f.TexPixelsAlpha8 == nil {
		if len(f.ConfigData) == 0 {
			f.AddFontDefault(nil)
		}
		f.Build()
	}
	out_pixels = f.TexPixelsAlpha8
	out_width = f.TexWidth
	out_height = f.TexHeight
	out_bytes_per_pixel = 1
	return
}

func (f *FontAtlas) GetTexDataAsRGBA32() (out_pixels []byte, out_width, out_height, out_bytes_per_pixel int) {
	// Convert to RGBA32 format on demand
	// Although it is likely to be the most commonly used format, our font rendering is 1 channel / 8 bpp
	if f.TexPixelsRGBA32 == nil {
		pixels, _, _, _ := f.GetTexDataAsAlpha8()
		f.TexPixelsRGBA32 = make([]byte, f.TexWidth*f.TexHeight*4)
		for i := 0; i < f.TexWidth*f.TexHeight; i++ {
			f.TexPixelsRGBA32[i*4+0] = 255
			f.TexPixelsRGBA32[i*4+1] = 255
			f.TexPixelsRGBA32[i*4+2] = 255
			f.TexPixelsRGBA32[i*4+3] = pixels[i]
		}
	}

	out_pixels = f.TexPixelsRGBA32
	out_width = f.TexWidth
	out_height = f.TexHeight
	out_bytes_per_pixel = 4

	return
}

// Load embedded ProggyClean.ttf at size 13, disable oversampling
func (f *FontAtlas) AddFontDefault(font_cfg_template *FontConfig) *Font {
	var font_cfg FontConfig
	if font_cfg_template != nil {
		font_cfg = *font_cfg_template
	} else {
		font_cfg.Init()
	}

	if font_cfg_template == nil {
		font_cfg.OversampleH = 1
		font_cfg.OversampleV = 1
		font_cfg.PixelSnapH = true
	}

	if font_cfg.Name == "" {
		font_cfg.Name = "ProggyClean.ttf, 13px"
	}
	if font_cfg.SizePixels <= 0 {
		font_cfg.SizePixels = 13
	}

	ttf_compressed_base85 := GetDefaultCompressedFontDataTTFBase85()
	font := f.AddFontFromMemoryCompressedBase85TTF(ttf_compressed_base85, font_cfg.SizePixels, &font_cfg, f.GetGlyphRangesDefault())
	font.DisplayOffset.Y = 1
	return font
}

func (f *FontAtlas) AddFontFromMemoryCompressedBase85TTF(compressed_ttf_data_base85 []byte, size_pixels float64, font_cfg *FontConfig, glyph_ranges []rune) *Font {
	compressed_ttf_size := ((len(compressed_ttf_data_base85) + 4) / 5) * 4
	compressed_ttf := make([]byte, compressed_ttf_size)
	Decode85(compressed_ttf_data_base85, compressed_ttf)
	font := f.AddFontFromMemoryCompressedTTF(compressed_ttf, size_pixels, font_cfg, glyph_ranges)
	return font
}

func (f *FontAtlas) AddFontFromMemoryCompressedTTF(compressed_ttf_data []byte, size_pixels float64, font_cfg_template *FontConfig, glyph_ranges []rune) *Font {
	var comp stbCompress
	buf_decompressed_size := comp.DecompressLength(compressed_ttf_data)
	buf_decompressed_data := make([]byte, buf_decompressed_size)
	comp.Decompress(buf_decompressed_data, compressed_ttf_data)

	var font_cfg FontConfig
	if font_cfg_template != nil {
		font_cfg = *font_cfg_template
	} else {
		font_cfg.Init()
	}
	font_cfg.FontDataOwnedByAtlas = true

	return f.AddFontFromMemoryTTF(buf_decompressed_data, size_pixels, &font_cfg, glyph_ranges)
}

// NB: Transfer ownership of 'ttf_data' to ImFontAtlas, unless font_cfg_template->FontDataOwnedByAtlas == false. Owned TTF buffer will be deleted after Build().
func (f *FontAtlas) AddFontFromMemoryTTF(ttf_data []byte, size_pixels float64, font_cfg_template *FontConfig, glyph_ranges []rune) *Font {
	var font_cfg FontConfig
	if font_cfg_template != nil {
		font_cfg = *font_cfg_template
	} else {
		font_cfg.Init()
	}
	font_cfg.FontData = ttf_data
	font_cfg.FontDataSize = len(ttf_data)
	font_cfg.SizePixels = size_pixels
	if glyph_ranges != nil {
		font_cfg.GlyphRanges = glyph_ranges
	}
	return f.AddFont(&font_cfg)
}

func (f *FontAtlas) AddFont(font_cfg *FontConfig) *Font {
	// Create new font
	if !font_cfg.MergeMode {
		f.Fonts = append(f.Fonts, NewFont())
	}

	f.ConfigData = append(f.ConfigData, *font_cfg)
	new_font_cfg := &f.ConfigData[len(f.ConfigData)-1]
	if new_font_cfg.DstFont == nil {
		new_font_cfg.DstFont = f.Fonts[len(f.Fonts)-1]
	}
	if !new_font_cfg.FontDataOwnedByAtlas {
		new_font_cfg.FontData = make([]byte, new_font_cfg.FontDataSize)
		new_font_cfg.FontDataOwnedByAtlas = true
		copy(new_font_cfg.FontData, font_cfg.FontData)
	}

	// Invalidate texture
	f.ClearTexData()
	return new_font_cfg.DstFont
}

func (f *FontAtlas) ClearTexData() {
	f.TexPixelsAlpha8 = nil
	f.TexPixelsRGBA32 = nil
}

// Retrieve list of range (2 int per range, values are inclusive)
func (f *FontAtlas) GetGlyphRangesDefault() []rune {
	return []rune{
		0x20, 0xff, // Basic Latin + Latin Supplement
		0,
	}
}

func (f *FontAtlas) Build() error {
	return f.BuildWithStbTruetype()
}

func (f *FontAtlas) BuildWithStbTruetype() error {
	assert(len(f.ConfigData) > 0)
	f.FontAtlasBuildRegisterDefaultCustomRects()

	f.TexID = nil
	f.TexWidth = 0
	f.TexHeight = 0
	f.TexUvScale = f64.Vec2{0, 0}
	f.TexUvWhitePixel = f64.Vec2{0, 0}
	f.ClearTexData()

	// Count glyphs/ranges
	total_glyphs_count := 0
	total_ranges_count := 0
	for input_i := 0; input_i < len(f.ConfigData); input_i++ {
		cfg := &f.ConfigData[input_i]
		if cfg.GlyphRanges == nil {
			cfg.GlyphRanges = f.GetGlyphRangesDefault()
		}
		for in_range := cfg.GlyphRanges; in_range[0] != 0 && in_range[1] != 0; in_range = in_range[2:] {
			total_glyphs_count += int(in_range[1]-in_range[0]) + 1
			total_ranges_count++
		}
	}

	// We need a width for the skyline algorithm. Using a dumb heuristic here to decide of width. User can override TexDesiredWidth and TexGlyphPadding if they wish.
	// Width doesn't really matter much, but some API/GPU have texture size limitations and increasing width can decrease height.
	if f.TexDesiredWidth > 0 {
		f.TexWidth = f.TexDesiredWidth
	} else {
		if total_glyphs_count > 4000 {
			f.TexWidth = 4096
		} else if total_glyphs_count > 2000 {
			f.TexWidth = 2048
		} else if total_glyphs_count > 1000 {
			f.TexWidth = 1024
		} else {
			f.TexWidth = 512
		}
	}
	f.TexHeight = 0

	// Start packing
	const max_tex_height = 1024 * 32
	spc := stbtt.NewPackContext()
	defer stbtt.FreePackContext(spc)
	err := spc.Begin(nil, f.TexWidth, max_tex_height, 0, f.TexGlyphPadding)
	if err != nil {
		return err
	}
	spc.SetOversampling(1, 1)

	// Pack our extra data rectangles first, so it will be on the upper-left corner of our texture (UV will have small values).
	f.BuildPackCustomRects(spc)

	// Initialize font information (so we can error without any cleanup)
	type FontTempBuildData struct {
		FontInfo *stbtt.FontInfo
		Rects    []stbtt.Rect
		Ranges   []stbtt.PackRange
	}
	tmp_array := make([]FontTempBuildData, len(f.ConfigData))
	for i := range tmp_array {
		tmp_array[i].FontInfo = stbtt.NewFontInfo()
		defer stbtt.FreeFontInfo(tmp_array[i].FontInfo)
	}
	for input_i := 0; input_i < len(f.ConfigData); input_i++ {
		cfg := &f.ConfigData[input_i]
		tmp := &tmp_array[input_i]
		assert(cfg.DstFont != nil && (!cfg.DstFont.IsLoaded() || cfg.DstFont.ContainerAtlas == f))

		font_offset := stbtt.GetFontOffsetForIndex(cfg.FontData, cfg.FontNo)
		// FontData is incorrect, or FontNo cannot be found.
		assert(font_offset >= 0)
		err := tmp.FontInfo.Init(cfg.FontData, font_offset)
		if err != nil {
			// Reset output on failure
			f.TexWidth = 0
			f.TexHeight = 0
			return err
		}
	}

	// Allocate packing character data and flag packed characters buffer as non-packed (x0=y0=x1=y1=0)
	buf_packedchars_n, buf_rects_n, buf_ranges_n := 0, 0, 0
	buf_packedchars := stbtt.MakePackedChars(total_glyphs_count)
	buf_rects := stbtt.MakeRects(total_glyphs_count)
	buf_ranges := stbtt.MakePackRanges(total_ranges_count)
	defer stbtt.FreePackedChars(buf_packedchars)
	defer stbtt.FreeRects(buf_rects)
	defer stbtt.FreePackRanges(buf_ranges)

	// First font pass: pack all glyphs (no rendering at this point, we are working with rectangles in an infinitely tall texture at this point)
	for input_i := 0; input_i < len(f.ConfigData); input_i++ {
		cfg := &f.ConfigData[input_i]
		tmp := &tmp_array[input_i]

		// Setup ranges
		font_glyphs_count := 0
		font_ranges_count := 0
		for in_range := cfg.GlyphRanges; in_range[0] != 0 && in_range[1] != 0; in_range = in_range[2:] {
			font_glyphs_count += int(in_range[1]-in_range[0]) + 1
			font_ranges_count++
		}
		tmp.Ranges = buf_ranges[buf_ranges_n : buf_ranges_n+font_ranges_count]
		buf_ranges_n += font_ranges_count

		for i := 0; i < font_ranges_count; i++ {
			in_range := cfg.GlyphRanges[i*2:]
			range_ := &tmp.Ranges[i]
			range_.SetFontSize(cfg.SizePixels)
			range_.SetFirstUnicodeCodepointInRange(int(in_range[0]))
			range_.SetNumChars(int(in_range[1] - in_range[0] + 1))
			range_.SetChardataForRange(buf_packedchars[buf_packedchars_n:])
			buf_packedchars_n += range_.NumChars()
		}

		// Pack
		tmp.Rects = buf_rects[buf_rects_n : buf_rects_n+font_glyphs_count]
		buf_rects_n += font_glyphs_count
		spc.SetOversampling(uint(cfg.OversampleH), uint(cfg.OversampleV))
		n := spc.FontRangesGatherRects(tmp.FontInfo, tmp.Ranges, tmp.Rects)
		assert(n == font_glyphs_count)
		spc.FontRangesPackRects(tmp.Rects[:n])

		// Extend texture height
		for i := 0; i < n; i++ {
			if tmp.Rects[i].WasPacked() != 0 {
				f.TexHeight = mathutil.Max(f.TexHeight, tmp.Rects[i].Y()+tmp.Rects[i].H())
			}
		}
	}
	assert(buf_rects_n == total_glyphs_count)
	assert(buf_packedchars_n == total_glyphs_count)
	assert(buf_ranges_n == total_ranges_count)

	// Create texture
	if f.Flags&FontAtlasFlagsNoPowerOfTwoHeight != 0 {
		f.TexHeight = f.TexHeight + 1
	} else {
		f.TexHeight = UpperPowerOfTwo(f.TexHeight)
	}
	f.TexUvScale = f64.Vec2{1.0 / float64(f.TexWidth), 1.0 / float64(f.TexHeight)}
	f.TexPixelsAlpha8 = make([]byte, f.TexWidth*f.TexHeight)
	spc.SetPixels(f.TexPixelsAlpha8)
	spc.SetHeight(f.TexHeight)

	// Second pass: render font characters
	for input_i := 0; input_i < len(f.ConfigData); input_i++ {
		cfg := &f.ConfigData[input_i]
		tmp := &tmp_array[input_i]
		spc.SetOversampling(uint(cfg.OversampleH), uint(cfg.OversampleV))
		spc.FontRangesRenderIntoRects(tmp.FontInfo, tmp.Ranges, tmp.Rects)

		if cfg.RasterizerMultiply != 1.0 {
			var multiply_table [256]byte
			f.BuildMultiplyCalcLookupTable(multiply_table[:], cfg.RasterizerMultiply)
			for i := range tmp.Rects {
				r := &tmp.Rects[i]
				if r.WasPacked() != 0 {
					f.BuildMultiplyRectAlpha8(multiply_table[:], spc.Pixels(), r.X(), r.Y(), r.W(), r.H(), spc.StrideInBytes())
				}
			}
		}
		tmp.Rects = nil
	}

	// End packing
	spc.End()

	// Third pass: setup ImFont and glyphs for runtime
	for input_i := 0; input_i < len(f.ConfigData); input_i++ {
		cfg := &f.ConfigData[input_i]
		tmp := &tmp_array[input_i]
		// We can have multiple input fonts writing into a same destination font (when using MergeMode=true)
		dst_font := cfg.DstFont
		if cfg.MergeMode {
			dst_font.BuildLookupTable()
		}
		font_scale := tmp.FontInfo.ScaleForPixelHeight(cfg.SizePixels)
		unscaled_ascent, unscaled_descent, _ := tmp.FontInfo.FontVMetrics()

		sign := -1.0
		if unscaled_ascent > 0 {
			sign = 1
		}
		ascent := math.Floor(float64(unscaled_ascent)*font_scale + sign)

		sign = -1.0
		if unscaled_descent > 0 {
			sign = 1
		}
		descent := math.Floor(float64(unscaled_descent)*font_scale + sign)

		f.BuildSetupFont(dst_font, f.ConfigData[input_i:], ascent, descent)
		off_x := cfg.GlyphOffset.X
		off_y := cfg.GlyphOffset.Y + float64(int(dst_font.Ascent+0.5))
		for i := range tmp.Ranges {
			range_ := tmp.Ranges[i]
			for char_idx := 0; char_idx < range_.NumChars(); char_idx += 1 {
				chardata_for_range := range_.CharDataForRange()
				pc := &chardata_for_range[char_idx]
				if pc.X0() == 0 && pc.X1() == 0 && pc.Y0() == 0 && pc.Y1() == 0 {
					continue
				}

				first_unicode_codepoint_in_range := range_.FirstUnicodepointInRange()
				codepoint := rune(first_unicode_codepoint_in_range + char_idx)
				if cfg.MergeMode && dst_font.FindGlyphNoFallback(codepoint) == nil {
					continue
				}

				_, _, q := stbtt.GetPackedQuad(chardata_for_range, f.TexWidth, f.TexHeight, char_idx, 0)
				dst_font.AddGlyph(codepoint, q.X0+off_x, q.Y0+off_y, q.X1+off_x, q.Y1+off_y, q.S0, q.T0, q.S1, q.T1, pc.XAdvance())
			}
		}
	}

	f.BuildFinish()
	return nil
}

func (f *FontAtlas) BuildSetupFont(font *Font, font_config []FontConfig, ascent, descent float64) {
	if !font_config[0].MergeMode {
		font.ClearOutputData()
		font.FontSize = font_config[0].SizePixels
		font.ConfigData = font_config
		font.ContainerAtlas = f
		font.Ascent = ascent
		font.Descent = descent
	}
	font.ConfigDataCount++
}

func (f *FontAtlas) BuildMultiplyCalcLookupTable(out_table []uint8, in_brighten_factor float64) {
	for i := range out_table {
		value := uint(float64(i) * in_brighten_factor)
		if value > 255 {
			out_table[i] = 255
		} else {
			out_table[i] = uint8(value)
		}
	}
}

func (f *FontAtlas) BuildMultiplyRectAlpha8(table, pixels []byte, x, y, w, h, stride int) {
	n := x + y*stride
	for j := h; j > 0; j, n = j-1, n+stride {
		for i := 0; i < w; i++ {
			pixels[n] = table[pixels[n]]
		}
	}
}

func (f *FontAtlas) BuildRenderDefaultTexData() {
	assert(f.CustomRectIds[0] >= 0)
	assert(f.TexPixelsAlpha8 != nil)
	r := &f.CustomRects[f.CustomRectIds[0]]
	assert(r.ID == FONT_ATLAS_DEFAULT_TEX_DATA_ID)
	assert(r.IsPacked())

	w := f.TexWidth
	if f.Flags&FontAtlasFlagsNoMouseCursors == 0 {
		// Render/copy pixels
		assert(r.Width == FONT_ATLAS_DEFAULT_TEX_DATA_W_HALF*2+1 && r.Height == FONT_ATLAS_DEFAULT_TEX_DATA_H)
		for y, n := 0, 0; y < FONT_ATLAS_DEFAULT_TEX_DATA_H; y++ {
			for x := 0; x < FONT_ATLAS_DEFAULT_TEX_DATA_W_HALF; x, n = x+1, n+1 {
				offset0 := int(int(r.X)+x) + int(int(r.Y)+y)*w
				offset1 := offset0 + FONT_ATLAS_DEFAULT_TEX_DATA_W_HALF + 1
				f.TexPixelsAlpha8[offset0] = 0x00
				if FONT_ATLAS_DEFAULT_TEX_DATA_PIXELS[n] == '.' {
					f.TexPixelsAlpha8[offset0] = 0xFF
				}
				f.TexPixelsAlpha8[offset1] = 0x00
				if FONT_ATLAS_DEFAULT_TEX_DATA_PIXELS[n] == 'X' {
					f.TexPixelsAlpha8[offset1] = 0xFF
				}
			}
		}
	} else {
		assert(r.Width == 2 && r.Height == 2)
		offset := int(r.X) + int(r.Y)*w
		f.TexPixelsAlpha8[offset] = 0xFF
		f.TexPixelsAlpha8[offset+1] = 0xFF
		f.TexPixelsAlpha8[offset+w] = 0xFF
		f.TexPixelsAlpha8[offset+w+1] = 0xFF
	}
	f.TexUvWhitePixel = f64.Vec2{(float64(r.X) + 0.5) * f.TexUvScale.X, (float64(r.Y) + 0.5) * f.TexUvScale.Y}
}

func (f *FontAtlas) CalcCustomRectUV(rect *CustomRect) (out_uv_min, out_uv_max f64.Vec2) {
	// Font atlas needs to be built before we can calculate UV coordinates
	assert(f.TexWidth > 0 && f.TexHeight > 0)
	// Make sure the rectangle has been packed
	assert(rect.IsPacked())
	out_uv_min = f64.Vec2{
		float64(rect.X) * f.TexUvScale.X,
		float64(rect.Y) * f.TexUvScale.Y,
	}
	out_uv_max = f64.Vec2{
		float64(rect.X+rect.Width) * f.TexUvScale.X,
		float64(rect.Y+rect.Height) * f.TexUvScale.Y,
	}
	return
}

func (f *FontAtlas) BuildFinish() {
	// Render into our custom data block
	f.BuildRenderDefaultTexData()

	// Register custom rectangle glyphs
	for i := range f.CustomRects {
		r := f.CustomRects[i]
		if r.Font == nil || r.ID > 0x10000 {
			continue
		}

		assert(r.Font.ContainerAtlas == f)
		f.CalcCustomRectUV(&r)
	}

	// Build all fonts lookup tables
	for i := range f.Fonts {
		if f.Fonts[i].DirtyLookupTables {
			f.Fonts[i].BuildLookupTable()
		}
	}
}

func (f *FontAtlas) BuildPackCustomRects(spc *stbtt.PackContext) {
	user_rects := f.CustomRects
	// We expect at least the default custom rects to be registered, else something went wrong.
	assert(len(user_rects) >= 1)

	pack_rects := make([]stbtt.Rect, len(user_rects))
	for i := range pack_rects {
		pack_rects[i].SetW(int(user_rects[i].Width))
		pack_rects[i].SetH(int(user_rects[i].Height))
	}

	spc.FontRangesPackRects(pack_rects)

	for i := range pack_rects {
		if pack_rects[i].WasPacked() != 0 {
			user_rects[i].X = uint(pack_rects[i].X())
			user_rects[i].Y = uint(pack_rects[i].Y())
			assert(pack_rects[i].W() == int(user_rects[i].Width) && pack_rects[i].H() == int(user_rects[i].Height))
			f.TexHeight = mathutil.Max(f.TexHeight, pack_rects[i].Y()+pack_rects[i].H())
		}
	}
}

func (f *FontAtlas) FontAtlasBuildRegisterDefaultCustomRects() {
	if f.CustomRectIds[0] >= 0 {
		return
	}
	if f.Flags&FontAtlasFlagsNoMouseCursors == 0 {
		f.CustomRectIds[0] = f.AddCustomRectRegular(FONT_ATLAS_DEFAULT_TEX_DATA_ID, FONT_ATLAS_DEFAULT_TEX_DATA_W_HALF*2+1, FONT_ATLAS_DEFAULT_TEX_DATA_H)
	} else {
		f.CustomRectIds[0] = f.AddCustomRectRegular(FONT_ATLAS_DEFAULT_TEX_DATA_ID, 2, 2)
	}
}

func (f *FontAtlas) AddCustomRectRegular(id, width, height uint) int {
	assert(id >= 0x10000)
	assert(width > 0 && width <= 0xFFFF)
	assert(height > 0 && height <= 0xFFFF)
	r := CustomRect{}
	r.Init()
	r.ID = id
	r.Width = width
	r.Height = height
	f.CustomRects = append(f.CustomRects, r)
	return len(f.CustomRects) - 1
}

// A work of art lies ahead! (. = white layer, X = black layer, others are blank)
// The white texels on the top left are the ones we'll use everywhere in ImGui to render filled shapes.
const (
	FONT_ATLAS_DEFAULT_TEX_DATA_W_HALF = 90
	FONT_ATLAS_DEFAULT_TEX_DATA_H      = 27
	FONT_ATLAS_DEFAULT_TEX_DATA_ID     = 0x80000000
)

var FONT_ATLAS_DEFAULT_TEX_DATA_PIXELS = []byte(
	"..-         -XXXXXXX-    X    -           X           -XXXXXXX          -          XXXXXXX" +
		"..-         -X.....X-   X.X   -          X.X          -X.....X          -          X.....X" +
		"---         -XXX.XXX-  X...X  -         X...X         -X....X           -           X....X" +
		"X           -  X.X  - X.....X -        X.....X        -X...X            -            X...X" +
		"XX          -  X.X  -X.......X-       X.......X       -X..X.X           -           X.X..X" +
		"X.X         -  X.X  -XXXX.XXXX-       XXXX.XXXX       -X.X X.X          -          X.X X.X" +
		"X..X        -  X.X  -   X.X   -          X.X          -XX   X.X         -         X.X   XX" +
		"X...X       -  X.X  -   X.X   -    XX    X.X    XX    -      X.X        -        X.X      " +
		"X....X      -  X.X  -   X.X   -   X.X    X.X    X.X   -       X.X       -       X.X       " +
		"X.....X     -  X.X  -   X.X   -  X..X    X.X    X..X  -        X.X      -      X.X        " +
		"X......X    -  X.X  -   X.X   - X...XXXXXX.XXXXXX...X -         X.X   XX-XX   X.X         " +
		"X.......X   -  X.X  -   X.X   -X.....................X-          X.X X.X-X.X X.X          " +
		"X........X  -  X.X  -   X.X   - X...XXXXXX.XXXXXX...X -           X.X..X-X..X.X           " +
		"X.........X -XXX.XXX-   X.X   -  X..X    X.X    X..X  -            X...X-X...X            " +
		"X..........X-X.....X-   X.X   -   X.X    X.X    X.X   -           X....X-X....X           " +
		"X......XXXXX-XXXXXXX-   X.X   -    XX    X.X    XX    -          X.....X-X.....X          " +
		"X...X..X    ---------   X.X   -          X.X          -          XXXXXXX-XXXXXXX          " +
		"X..X X..X   -       -XXXX.XXXX-       XXXX.XXXX       ------------------------------------" +
		"X.X  X..X   -       -X.......X-       X.......X       -    XX           XX    -           " +
		"XX    X..X  -       - X.....X -        X.....X        -   X.X           X.X   -           " +
		"      X..X          -  X...X  -         X...X         -  X..X           X..X  -           " +
		"       XX           -   X.X   -          X.X          - X...XXXXXXXXXXXXX...X -           " +
		"------------        -    X    -           X           -X.....................X-           " +
		"                    ----------------------------------- X...XXXXXXXXXXXXX...X -           " +
		"                                                      -  X..X           X..X  -           " +
		"                                                      -   X.X           X.X   -           " +
		"                                                      -    XX           XX    -           ")

var FONT_ATLAS_DEFAULT_TEX_CURSOR_DATA = [MouseCursorCOUNT][3]f64.Vec2{
	// Pos ........ Size ......... Offset ......
	{f64.Vec2{0, 3}, f64.Vec2{12, 19}, f64.Vec2{0, 0}},    // ImGuiMouseCursor_Arrow
	{f64.Vec2{13, 0}, f64.Vec2{7, 16}, f64.Vec2{4, 8}},    // ImGuiMouseCursor_TextInput
	{f64.Vec2{31, 0}, f64.Vec2{23, 23}, f64.Vec2{11, 11}}, // ImGuiMouseCursor_ResizeAll
	{f64.Vec2{21, 0}, f64.Vec2{9, 23}, f64.Vec2{5, 11}},   // ImGuiMouseCursor_ResizeNS
	{f64.Vec2{55, 18}, f64.Vec2{23, 9}, f64.Vec2{11, 5}},  // ImGuiMouseCursor_ResizeEW
	{f64.Vec2{73, 0}, f64.Vec2{17, 17}, f64.Vec2{9, 9}},   // ImGuiMouseCursor_ResizeNESW
	{f64.Vec2{55, 0}, f64.Vec2{17, 17}, f64.Vec2{9, 9}},   // ImGuiMouseCursor_ResizeNWSE
}

//-----------------------------------------------------------------------------
// DEFAULT FONT DATA
//-----------------------------------------------------------------------------
// Compressed with stb_compress() then converted to a C array.
// Use the program in misc/fonts/binary_to_compressed_c.cpp to create the array from a TTF file.
// Decompression from stb.h (public domain) by Sean Barrett https://github.com/nothings/stb/blob/master/stb.h
//-----------------------------------------------------------------------------

func GetDefaultCompressedFontDataTTFBase85() []byte {
	return []byte(proggy_clean_ttf_compressed_data_base85)
}

//-----------------------------------------------------------------------------
// ProggyClean.ttf
// Copyright (c) 2004, 2005 Tristan Grimmer
// MIT license (see License.txt in http://www.upperbounds.net/download/ProggyClean.ttf.zip)
// Download and more information at http://upperbounds.net
//-----------------------------------------------------------------------------
// File: 'ProggyClean.ttf' (41208 bytes)
// Exported using binary_to_compressed_c.cpp
//-----------------------------------------------------------------------------
const proggy_clean_ttf_compressed_data_base85 = "7])#######hV0qs'/###[),##/l:$#Q6>##5[n42>c-TH`->>#/e>11NNV=Bv(*:.F?uu#(gRU.o0XGH`$vhLG1hxt9?W`#,5LsCp#-i>.r$<$6pD>Lb';9Crc6tgXmKVeU2cD4Eo3R/2" +
	"*>]b(MC;$jPfY.;h^`IWM9<Lh2TlS+f-s$o6Q<BWH`YiU.xfLq$N;$0iR/GX:U(jcW2p/W*q?-qmnUCI;jHSAiFWM.R*kU@C=GH?a9wp8f$e.-4^Qg1)Q-GL(lf(r/7GrRgwV%MS=C#`" +
	"8ND>Qo#t'X#(v#Y9w0#1D$CIf;W'#pWUPXOuxXuU(H9M(1<q-UE31#^-V'8IRUo7Qf./L>=Ke$$'5F%)]0^#0X@U.a<r:QLtFsLcL6##lOj)#.Y5<-R&KgLwqJfLgN&;Q?gI^#DY2uLi" +
	"@^rMl9t=cWq6##weg>$FBjVQTSDgEKnIS7EM9>ZY9w0#L;>>#Mx&4Mvt//L[MkA#W@lK.N'[0#7RL_&#w+F%HtG9M#XL`N&.,GM4Pg;-<nLENhvx>-VsM.M0rJfLH2eTM`*oJMHRC`Nk" +
	"fimM2J,W-jXS:)r0wK#@Fge$U>`w'N7G#$#fB#$E^$#:9:hk+eOe--6x)F7*E%?76%^GMHePW-Z5l'&GiF#$956:rS?dA#fiK:)Yr+`&#0j@'DbG&#^$PG.Ll+DNa<XCMKEV*N)LN/N*" +
	"b=%Q6pia-Xg8I$<MR&,VdJe$<(7G;Ckl'&hF;;$<_=X(b.RS%%)###MPBuuE1V:v&cX&#2m#(&cV]`k9OhLMbn%s$G2,B$BfD3X*sp5#l,$R#]x_X1xKX%b5U*[r5iMfUo9U`N99hG)t" +
	"m+/Us9pG)XPu`<0s-)WTt(gCRxIg(%6sfh=ktMKn3j)<6<b5Sk_/0(^]AaN#(p/L>&VZ>1i%h1S9u5o@YaaW$e+b<TWFn/Z:Oh(Cx2$lNEoN^e)#CFY@@I;BOQ*sRwZtZxRcU7uW6CXo" +
	"w0i(?$Q[cjOd[P4d)]>ROPOpxTO7Stwi1::iB1q)C_=dV26J;2,]7op$]uQr@_V7$q^%lQwtuHY]=DX,n3L#0PHDO4f9>dC@O>HBuKPpP*E,N+b3L#lpR/MrTEH.IAQk.a>D[.e;mc.x" +
	"]Ip.PH^'/aqUO/$1WxLoW0[iLA<QT;5HKD+@qQ'NQ(3_PLhE48R.qAPSwQ0/WK?Z,[x?-J;jQTWA0X@KJ(_Y8N-:/M74:/-ZpKrUss?d#dZq]DAbkU*JqkL+nwX@@47`5>w=4h(9.`GC" +
	"RUxHPeR`5Mjol(dUWxZa(>STrPkrJiWx`5U7F#.g*jrohGg`cg:lSTvEY/EV_7H4Q9[Z%cnv;JQYZ5q.l7Zeas:HOIZOB?G<Nald$qs]@]L<J7bR*>gv:[7MI2k).'2($5FNP&EQ(,)U" +
	"]W]+fh18.vsai00);D3@4ku5P?DP8aJt+;qUM]=+b'8@;mViBKx0DE[-auGl8:PJ&Dj+M6OC]O^((##]`0i)drT;-7X`=-H3[igUnPG-NZlo.#k@h#=Ork$m>a>$-?Tm$UV(?#P6YY#'" +
	"/###xe7q.73rI3*pP/$1>s9)W,JrM7SN]'/4C#v$U`0#V.[0>xQsH$fEmPMgY2u7Kh(G%siIfLSoS+MK2eTM$=5,M8p`A.;_R%#u[K#$x4AG8.kK/HSB==-'Ie/QTtG?-.*^N-4B/ZM_" +
	"3YlQC7(p7q)&](`6_c)$/*JL(L-^(]$wIM`dPtOdGA,U3:w2M-0<q-]L_?^)1vw'.,MRsqVr.L;aN&#/EgJ)PBc[-f>+WomX2u7lqM2iEumMTcsF?-aT=Z-97UEnXglEn1K-bnEO`guF" +
	"t(c%=;Am_Qs@jLooI&NX;]0#j4#F14;gl8-GQpgwhrq8'=l_f-b49'UOqkLu7-##oDY2L(te+Mch&gLYtJ,MEtJfLh'x'M=$CS-ZZ%P]8bZ>#S?YY#%Q&q'3^Fw&?D)UDNrocM3A76//" +
	"oL?#h7gl85[qW/NDOk%16ij;+:1a'iNIdb-ou8.P*w,v5#EI$TWS>Pot-R*H'-SEpA:g)f+O$%%`kA#G=8RMmG1&O`>to8bC]T&$,n.LoO>29sp3dt-52U%VM#q7'DHpg+#Z9%H[K<L%" +
	"a2E-grWVM3@2=-k22tL]4$##6We'8UJCKE[d_=%wI;'6X-GsLX4j^SgJ$##R*w,vP3wK#iiW&#*h^D&R?jp7+/u&#(AP##XU8c$fSYW-J95_-Dp[g9wcO&#M-h1OcJlc-*vpw0xUX&#O" +
	"QFKNX@QI'IoPp7nb,QU//MQ&ZDkKP)X<WSVL(68uVl&#c'[0#(s1X&xm$Y%B7*K:eDA323j998GXbA#pwMs-jgD$9QISB-A_(aN4xoFM^@C58D0+Q+q3n0#3U1InDjF682-SjMXJK)(h" +
	"$hxua_K]ul92%'BOU&#BRRh-slg8KDlr:%L71Ka:.A;%YULjDPmL<LYs8i#XwJOYaKPKc1h:'9Ke,g)b),78=I39B;xiY$bgGw-&.Zi9InXDuYa%G*f2Bq7mn9^#p1vv%#(Wi-;/Z5ho" +
	";#2:;%d&#x9v68C5g?ntX0X)pT`;%pB3q7mgGN)3%(P8nTd5L7GeA-GL@+%J3u2:(Yf>et`e;)f#Km8&+DC$I46>#Kr]]u-[=99tts1.qb#q72g1WJO81q+eN'03'eM>&1XxY-caEnOj" +
	"%2n8)),?ILR5^.Ibn<-X-Mq7[a82Lq:F&#ce+S9wsCK*x`569E8ew'He]h:sI[2LM$[guka3ZRd6:t%IG:;$%YiJ:Nq=?eAw;/:nnDq0(CYcMpG)qLN4$##&J<j$UpK<Q4a1]MupW^-s" +
	"j_$%[HK%'F####QRZJ::Y3EGl4'@%FkiAOg#p[##O`gukTfBHagL<LHw%q&OV0##F=6/:chIm0@eCP8X]:kFI%hl8hgO@RcBhS-@Qb$%+m=hPDLg*%K8ln(wcf3/'DW-$.lR?n[nCH-e" +
	"XOONTJlh:.RYF%3'p6sq:UIMA945&^HFS87@$EP2iG<-lCO$%c`uKGD3rC$x0BL8aFn--`ke%#HMP'vh1/R&O_J9'um,.<tx[@%wsJk&bUT2`0uMv7gg#qp/ij.L56'hl;.s5CUrxjOM" +
	"7-##.l+Au'A&O:-T72L]P`&=;ctp'XScX*rU.>-XTt,%OVU4)S1+R-#dg0/Nn?Ku1^0f$B*P:Rowwm-`0PKjYDDM'3]d39VZHEl4,.j']Pk-M.h^&:0FACm$maq-&sgw0t7/6(^xtk%L" +
	"uH88Fj-ekm>GA#_>568x6(OFRl-IZp`&b,_P'$M<Jnq79VsJW/mWS*PUiq76;]/NM_>hLbxfc$mj`,O;&%W2m`Zh:/)Uetw:aJ%]K9h:TcF]u_-Sj9,VK3M.*'&0D[Ca]J9gp8,kAW]%" +
	"(?A%R$f<->Zts'^kn=-^@c4%-pY6qI%J%1IGxfLU9CP8cbPlXv);C=b),<2mOvP8up,UVf3839acAWAW-W?#ao/^#%KYo8fRULNd2.>%m]UK:n%r$'sw]J;5pAoO_#2mO3n,'=H5(etH" +
	"g*`+RLgv>=4U8guD$I%D:W>-r5V*%j*W:Kvej.Lp$<M-SGZ':+Q_k+uvOSLiEo(<aD/K<CCc`'Lx>'?;++O'>()jLR-^u68PHm8ZFWe+ej8h:9r6L*0//c&iH&R8pRbA#Kjm%upV1g:a" +
	"_#Ur7FuA#(tRh#.Y5K+@?3<-8m0$PEn;J:rh6?I6uG<-`wMU'ircp0LaE_OtlMb&1#6T.#FDKu#1Lw%u%+GM+X'e?YLfjM[VO0MbuFp7;>Q&#WIo)0@F%q7c#4XAXN-U&VB<HFF*qL($" +
	"/V,;(kXZejWO`<[5??ewY(*9=%wDc;,u<'9t3W-(H1th3+G]ucQ]kLs7df($/*JL]@*t7Bu_G3_7mp7<iaQjO@.kLg;x3B0lqp7Hf,^Ze7-##@/c58Mo(3;knp0%)A7?-W+eI'o8)b<n" +
	"Knw'Ho8C=Y>pqB>0ie&jhZ[?iLR@@_AvA-iQC(=ksRZRVp7`.=+NpBC%rh&3]R:8XDmE5^V8O(x<<aG/1N$#FX$0V5Y6x'aErI3I$7x%E`v<-BY,)%-?Psf*l?%C3.mM(=/M0:JxG'?7" +
	"WhH%o'a<-80g0NBxoO(GH<dM]n.+%q@jH?f.UsJ2Ggs&4<-e47&Kl+f//9@`b+?.TeN_&B8Ss?v;^Trk;f#YvJkl&w$]>-+k?'(<S:68tq*WoDfZu';mM?8X[ma8W%*`-=;D.(nc7/;)" +
	"g:T1=^J$&BRV(-lTmNB6xqB[@0*o.erM*<SWF]u2=st-*(6v>^](H.aREZSi,#1:[IXaZFOm<-ui#qUq2$##Ri;u75OK#(RtaW-K-F`S+cF]uN`-KMQ%rP/Xri.LRcB##=YL3BgM/3MD" +
	"?@f&1'BW-)Ju<L25gl8uhVm1hL$##*8###'A3/LkKW+(^rWX?5W_8g)a(m&K8P>#bmmWCMkk&#TR`C,5d>g)F;t,4:@_l8G/5h4vUd%&%950:VXD'QdWoY-F$BtUwmfe$YqL'8(PWX(P" +
	"?^@Po3$##`MSs?DWBZ/S>+4%>fX,VWv/w'KD`LP5IbH;rTV>n3cEK8U#bX]l-/V+^lj3;vlMb&[5YQ8#pekX9JP3XUC72L,,?+Ni&co7ApnO*5NK,((W-i:$,kp'UDAO(G0Sq7MVjJsb" +
	"Iu)'Z,*[>br5fX^:FPAWr-m2KgL<LUN098kTF&#lvo58=/vjDo;.;)Ka*hLR#/k=rKbxuV`>Q_nN6'8uTG&#1T5g)uLv:873UpTLgH+#FgpH'_o1780Ph8KmxQJ8#H72L4@768@Tm&Qh" +
	"4CB/5OvmA&,Q&QbUoi$a_%3M01H)4x7I^&KQVgtFnV+;[Pc>[m4k//,]1?#`VY[Jr*3&&slRfLiVZJ:]?=K3Sw=[$=uRB?3xk48@aeg<Z'<$#4H)6,>e0jT6'N#(q%.O=?2S]u*(m<-V" +
	"8J'(1)G][68hW$5'q[GC&5j`TE?m'esFGNRM)j,ffZ?-qx8;->g4t*:CIP/[Qap7/9'#(1sao7w-.qNUdkJ)tCF&#B^;xGvn2r9FEPFFFcL@.iFNkTve$m%#QvQS8U@)2Z+3K:AKM5is" +
	"Z88+dKQ)W6>J%CL<KE>`.d*(B`-n8D9oK<Up]c$X$(,)M8Zt7/[rdkqTgl-0cuGMv'?>-XV1q['-5k'cAZ69e;D_?$ZPP&s^+7])$*$#@QYi9,5P&#9r+$%CE=68>K8r0=dSC%%(@p7." +
	"m7jilQ02'0-VWAg<a/''3u.=4L$Y)6k/K:_[3=&jvL<L0C/2'v:^;-DIBW,B4E68:kZ;%?8(Q8BH=kO65BW?xSG&#@uU,DS*,?.+(o(#1vCS8#CHF>TlGW'b)Tq7VT9q^*^$$.:&N@@$" +
	"&)WHtPm*5_rO0&e%K&#-30j(E4#'Zb.o/(Tpm$>K'f@[PvFl,hfINTNU6u'0pao7%XUp9]5.>%h`8_=VYbxuel.NTSsJfLacFu3B'lQSu/m6-Oqem8T+oE--$0a/k]uj9EwsG>%veR*h" +
	"v^BFpQj:K'#SJ,sB-'#](j.Lg92rTw-*n%@/;39rrJF,l#qV%OrtBeC6/,;qB3ebNW[?,Hqj2L.1NP&GjUR=1D8QaS3Up&@*9wP?+lo7b?@%'k4`p0Z$22%K3+iCZj?XJN4Nm&+YF]u@" +
	"-W$U%VEQ/,,>>#)D<h#`)h0:<Q6909ua+&VU%n2:cG3FJ-%@Bj-DgLr`Hw&HAKjKjseK</xKT*)B,N9X3]krc12t'pgTV(Lv-tL[xg_%=M_q7a^x?7Ubd>#%8cY#YZ?=,`Wdxu/ae&#w" +
	"6)R89tI#6@s'(6Bf7a&?S=^ZI_kS&ai`&=tE72L_D,;^R)7[$s<Eh#c&)q.MXI%#v9ROa5FZO%sF7q7Nwb&#ptUJ:aqJe$Sl68%.D###EC><?-aF&#RNQv>o8lKN%5/$(vdfq7+ebA#u" +
	"1p]ovUKW&Y%q]'>$1@-[xfn$7ZTp7mM,G,Ko7a&Gu%G[RMxJs[0MM%wci.LFDK)(<c`Q8N)jEIF*+?P2a8g%)$q]o2aH8C&<SibC/q,(e:v;-b#6[$NtDZ84Je2KNvB#$P5?tQ3nt(0d" +
	"=j.LQf./Ll33+(;q3L-w=8dX$#WF&uIJ@-bfI>%:_i2B5CsR8&9Z&#=mPEnm0f`<&c)QL5uJ#%u%lJj+D-r;BoF&#4DoS97h5g)E#o:&S4weDF,9^Hoe`h*L+_a*NrLW-1pG_&2UdB86" +
	"e%B/:=>)N4xeW.*wft-;$'58-ESqr<b?UI(_%@[P46>#U`'6AQ]m&6/`Z>#S?YY#Vc;r7U2&326d=w&H####?TZ`*4?&.MK?LP8Vxg>$[QXc%QJv92.(Db*B)gb*BM9dM*hJMAo*c&#b" +
	"0v=Pjer]$gG&JXDf->'StvU7505l9$AFvgYRI^&<^b68?j#q9QX4SM'RO#&sL1IM.rJfLUAj221]d##DW=m83u5;'bYx,*Sl0hL(W;;$doB&O/TQ:(Z^xBdLjL<Lni;''X.`$#8+1GD:" +
	"k$YUWsbn8ogh6rxZ2Z9]%nd+>V#*8U_72Lh+2Q8Cj0i:6hp&$C/:p(HK>T8Y[gHQ4`4)'$Ab(Nof%V'8hL&#<NEdtg(n'=S1A(Q1/I&4([%dM`,Iu'1:_hL>SfD07&6D<fp8dHM7/g+t" +
	"lPN9J*rKaPct&?'uBCem^jn%9_K)<,C5K3s=5g&GmJb*[SYq7K;TRLGCsM-$$;S%:Y@r7AK0pprpL<Lrh,q7e/%KWK:50I^+m'vi`3?%Zp+<-d+$L-Sv:@.o19n$s0&39;kn;S%BSq*$" +
	"3WoJSCLweV[aZ'MQIjO<7;X-X;&+dMLvu#^UsGEC9WEc[X(wI7#2.(F0jV*eZf<-Qv3J-c+J5AlrB#$p(H68LvEA'q3n0#m,[`*8Ft)FcYgEud]CWfm68,(aLA$@EFTgLXoBq/UPlp7:" +
	"d[/;r_ix=:TF`S5H-b<LI&HY(K=h#)]Lk$K14lVfm:x$H<3^Ql<M`$OhapBnkup'D#L$Pb_`N*g]2e;X/Dtg,bsj&K#2[-:iYr'_wgH)NUIR8a1n#S?Yej'h8^58UbZd+^FKD*T@;6A7" +
	"aQC[K8d-(v6GI$x:T<&'Gp5Uf>@M.*J:;$-rv29'M]8qMv-tLp,'886iaC=Hb*YJoKJ,(j%K=H`K.v9HggqBIiZu'QvBT.#=)0ukruV&.)3=(^1`o*Pj4<-<aN((^7('#Z0wK#5GX@7u" +
	"][`*S^43933A4rl][`*O4CgLEl]v$1Q3AeF37dbXk,.)vj#x'd`;qgbQR%FW,2(?LO=s%Sc68%NP'##Aotl8x=BE#j1UD([3$M(]UI2LX3RpKN@;/#f'f/&_mt&F)XdF<9t4)Qa.*kTL" +
	"wQ'(TTB9.xH'>#MJ+gLq9-##@HuZPN0]u:h7.T..G:;$/Usj(T7`Q8tT72LnYl<-qx8;-HV7Q-&Xdx%1a,hC=0u+HlsV>nuIQL-5<N?)NBS)QN*_I,?&)2'IM%L3I)X((e/dl2&8'<M:" +
	"^#M*Q+[T.Xri.LYS3v%fF`68h;b-X[/En'CR.q7E)p'/kle2HM,u;^%OKC-N+Ll%F9CF<Nf'^#t2L,;27W:0O@6##U6W7:$rJfLWHj$#)woqBefIZ.PK<b*t7ed;p*_m;4ExK#h@&]>_" +
	">@kXQtMacfD.m-VAb8;IReM3$wf0''hra*so568'Ip&vRs849'MRYSp%:t:h5qSgwpEr$B>Q,;s(C#$)`svQuF$##-D,##,g68@2[T;.XSdN9Qe)rpt._K-#5wF)sP'##p#C0c%-Gb%h" +
	"d+<-j'Ai*x&&HMkT]C'OSl##5RG[JXaHN;d'uA#x._U;.`PU@(Z3dt4r152@:v,'R.Sj'w#0<-;kPI)FfJ&#AYJ&#//)>-k=m=*XnK$>=)72L]0I%>.G690a:$##<,);?;72#?x9+d;^" +
	"V'9;jY@;)br#q^YQpx:X#Te$Z^'=-=bGhLf:D6&bNwZ9-ZD#n^9HhLMr5G;']d&6'wYmTFmL<LD)F^%[tC'8;+9E#C$g%#5Y>q9wI>P(9mI[>kC-ekLC/R&CH+s'B;K-M6$EB%is00:+" +
	"A4[7xks.LrNk0&E)wILYF@2L'0Nb$+pv<(2.768/FrY&h$^3i&@+G%JT'<-,v`3;_)I9M^AE]CN?Cl2AZg+%4iTpT3<n-&%H%b<FDj2M<hH=&Eh<2Len$b*aTX=-8QxN)k11IM1c^j%9" +
	"s<L<NFSo)B?+<-(GxsF,^-Eh@$4dXhN$+#rxK8'je'D7k`e;)2pYwPA'_p9&@^18ml1^[@g4t*[JOa*[=Qp7(qJ_oOL^('7fB&Hq-:sf,sNj8xq^>$U4O]GKx'm9)b@p7YsvK3w^YR-C" +
	"dQ*:Ir<($u&)#(&?L9Rg3H)4fiEp^iI9O8KnTj,]H?D*r7'M;PwZ9K0E^k&-cpI;.p/6_vwoFMV<->#%Xi.LxVnrU(4&8/P+:hLSKj$#U%]49t'I:rgMi'FL@a:0Y-uA[39',(vbma*h" +
	"U%<-SRF`Tt:542R_VV$p@[p8DV[A,?1839FWdF<TddF<9Ah-6&9tWoDlh]&1SpGMq>Ti1O*H&#(AL8[_P%.M>v^-))qOT*F5Cq0`Ye%+$B6i:7@0IX<N+T+0MlMBPQ*Vj>SsD<U4JHY8" +
	"kD2)2fU/M#$e.)T4,_=8hLim[&);?UkK'-x?'(:siIfL<$pFM`i<?%W(mGDHM%>iWP,##P`%/L<eXi:@Z9C.7o=@(pXdAO/NLQ8lPl+HPOQa8wD8=^GlPa8TKI1CjhsCTSLJM'/Wl>-S" +
	"(qw%sf/@%#B6;/U7K]uZbi^Oc^2n<bhPmUkMw>%t<)'mEVE''n`WnJra$^TKvX5B>;_aSEK',(hwa0:i4G?.Bci.(X[?b*($,=-n<.Q%`(X=?+@Am*Js0&=3bh8K]mL<LoNs'6,'85`0" +
	"?t/'_U59@]ddF<#LdF<eWdF<OuN/45rY<-L@&#+fm>69=Lb,OcZV/);TTm8VI;?%OtJ<(b4mq7M6:u?KRdF<gR@2L=FNU-<b[(9c/ML3m;Z[$oF3g)GAWqpARc=<ROu7cL5l;-[A]%/+" +
	"fsd;l#SafT/f*W]0=O'$(Tb<[)*@e775R-:Yob%g*>l*:xP?Yb.5)%w_I?7uk5JC+FS(m#i'k.'a0i)9<7b'fs'59hq$*5Uhv##pi^8+hIEBF`nvo`;'l0.^S1<-wUK2/Coh58KKhLjM" +
	"=SO*rfO`+qC`W-On.=AJ56>>i2@2LH6A:&5q`?9I3@@'04&p2/LVa*T-4<-i3;M9UvZd+N7>b*eIwg:CC)c<>nO&#<IGe;__.thjZl<%w(Wk2xmp4Q@I#I9,DF]u7-P=.-_:YJ]aS@V?" +
	"6*C()dOp7:WL,b&3Rg/.cmM9&r^>$(>.Z-I&J(Q0Hd5Q%7Co-b`-c<N(6r@ip+AurK<m86QIth*#v;-OBqi+L7wDE-Ir8K['m+DDSLwK&/.?-V%U_%3:qKNu$_b*B-kp7NaD'QdWQPKY" +
	"q[@>P)hI;*_F]u`Rb[.j8_Q/<&>uu+VsH$sM9TA%?)(vmJ80),P7E>)tjD%2L=-t#fK[%`v=Q8<FfNkgg^oIbah*#8/Qt$F&:K*-(N/'+1vMB,u()-a.VUU*#[e%gAAO(S>WlA2);Sa>" +
	"gXm8YB`1d@K#n]76-a$U,mF<fX]idqd)<3,]J7JmW4`6]uks=4-72L(jEk+:bJ0M^q-8Dm_Z?0olP1C9Sa&H[d&c$ooQUj]Exd*3ZM@-WGW2%s',B-_M%>%Ul:#/'xoFM9QX-$.QN'>[" +
	"%$Z$uF6pA6Ki2O5:8w*vP1<-1`[G,)-m#>0`P&#eb#.3i)rtB61(o'$?X3B</R90;eZ]%Ncq;-Tl]#F>2Qft^ae_5tKL9MUe9b*sLEQ95C&`=G?@Mj=wh*'3E>=-<)Gt*Iw)'QG:`@Iw" +
	"Of7&]1i'S01B+Ev/Nac#9S;=;YQpg_6U`*kVY39xK,[/6Aj7:'1Bm-_1EYfa1+o&o4hp7KN_Q(OlIo@S%;jVdn0'1<Vc52=u`3^o-n1'g4v58Hj&6_t7$##?M)c<$bgQ_'SY((-xkA#Y" +
	"(,p'H9rIVY-b,'%bCPF7.J<Up^,(dU1VY*5#WkTU>h19w,WQhLI)3S#f$2(eb,jr*b;3Vw]*7NH%$c4Vs,eD9>XW8?N]o+(*pgC%/72LV-u<Hp,3@e^9UB1J+ak9-TN/mhKPg+AJYd$M" +
	"lvAF_jCK*.O-^(63adMT->W%iewS8W6m2rtCpo'RS1R84=@paTKt)>=%&1[)*vp'u+x,VrwN;&]kuO9JDbg=pO$J*.jVe;u'm0dr9l,<*wMK*Oe=g8lV_KEBFkO'oU]^=[-792#ok,)i" +
	"]lR8qQ2oA8wcRCZ^7w/Njh;?.stX?Q1>S1q4Bn$)K1<-rGdO'$Wr.Lc.CG)$/*JL4tNR/,SVO3,aUw'DJN:)Ss;wGn9A32ijw%FL+Z0Fn.U9;reSq)bmI32U==5ALuG&#Vf1398/pVo1" +
	"*c-(aY168o<`JsSbk-,1N;$>0:OUas(3:8Z972LSfF8eb=c-;>SPw7.6hn3m`9^Xkn(r.qS[0;T%&Qc=+STRxX'q1BNk3&*eu2;&8q$&x>Q#Q7^Tf+6<(d%ZVmj2bDi%.3L2n+4W'$Pi" +
	"DDG)g,r%+?,$@?uou5tSe2aN_AQU*<h`e-GI7)?OK2A.d7_c)?wQ5AS@DL3r#7fSkgl6-++D:'A,uq7SvlB$pcpH'q3n0#_%dY#xCpr-l<F0NR@-##FEV6NTF6##$l84N1w?AO>'IAOU" +
	"RQ##V^Fv-XFbGM7Fl(N<3DhLGF%q.1rC$#:T__&Pi68%0xi_&[qFJ(77j_&JWoF.V735&T,[R*:xFR*K5>>#`bW-?4Ne_&6Ne_&6Ne_&n`kr-#GJcM6X;uM6X;uM(.a..^2TkL%oR(#;" +
	"u.T%fAr%4tJ8&><1=GHZ_+m9/#H1F^R#SC#*N=BA9(D?v[UiFY>>^8p,KKF.W]L29uLkLlu/+4T<XoIB&hx=T1PcDaB&;HH+-AFr?(m9HZV)FKS8JCw;SD=6[^/DZUL`EUDf]GGlG&>w" +
	"$)F./^n3+rlo+DB;5sIYGNk+i1t-69Jg--0pao7Sm#K)pdHW&;LuDNH@H>#/X-TI(;P>#,Gc>#0Su>#4`1?#8lC?#<xU?#@.i?#D:%@#HF7@#LRI@#P_[@#Tkn@#Xw*A#]-=A#a9OA#d" +
	"<F&#*;G##.GY##2Sl##6`($#:l:$#>xL$#B.`$#F:r$#JF.%#NR@%#R_R%#Vke%#Zww%#_-4&#3^Rh%Sflr-k'MS.o?.5/sWel/wpEM0%3'/1)K^f1-d>G21&v(35>V`39V7A4=onx4A" +
	"1OY5EI0;6Ibgr6M$HS7Q<)58C5w,;WoA*#[%T*#`1g*#d=#+#hI5+#lUG+#pbY+#tnl+#x$),#&1;,#*=M,#.I`,#2Ur,#6b.-#;w[H#iQtA#m^0B#qjBB#uvTB##-hB#'9$C#+E6C#/" +
	"QHC#3^ZC#7jmC#;v)D#?,<D#C8ND#GDaD#KPsD#O]/E#g1A5#KA*1#gC17#MGd;#8(02#L-d3#rWM4#Hga1#,<w0#T.j<#O#'2#CYN1#qa^:#_4m3#o@/=#eG8=#t8J5#`+78#4uI-#m" +
	"3B2#SB[8#Q0@8#i[*9#iOn8#1Nm;#^sN9#qh<9#:=x-#P;K2#$%X9#bC+.#Rg;<#mN=.#MTF.#RZO.#2?)4#Y#(/#[)1/#b;L/#dAU/#0Sv;#lY$0#n`-0#sf60#(F24#wrH0#%/e0#T" +
	"mD<#%JSMFove:CTBEXI:<eh2g)B,3h2^G3i;#d3jD>)4kMYD4lVu`4m`:&5niUA5@(A5BA1]PBB:xlBCC=2CDLXMCEUtiCf&0g2'tN?PGT4CPGT4CPGT4CPGT4CPGT4CPGT4CPGT4CPG" +
	"T4CPGT4CPGT4CPGT4CPGT4CPGT4CP-qekC`.9kEg^+F$kwViFJTB&5KTB&5KTB&5KTB&5KTB&5KTB&5KTB&5KTB&5KTB&5KTB&5KTB&5KTB&5KTB&5KTB&5KTB&5o,^<-28ZI'O?;xpO" +
	"?;xpO?;xpO?;xpO?;xpO?;xpO?;xpO?;xpO?;xpO?;xpO?;xpO?;xpO?;xpO?;xp;7q-#lLYI:xvD=#"
