package imageutil

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/qeedquan/go-media/image/chroma"
	"github.com/qeedquan/go-media/image/pnm"
	_ "github.com/qeedquan/go-media/image/psd"
	"github.com/qeedquan/go-media/image/tga"
	"github.com/qeedquan/go-media/math/ga"
	"golang.org/x/image/bmp"
)

func LoadRGBADir(dir string) ([]*image.RGBA, error) {
	exts := []string{"png", "jpg", "jpeg", "tga", "gif"}

	var m []*image.RGBA
	for _, ext := range exts {
		glob := fmt.Sprintf("%v/*.%v", dir, ext)
		p, _ := LoadRGBAGlob(glob)
		m = append(m, p...)
	}
	return m, nil
}

func LoadRGBAGlob(glob string) ([]*image.RGBA, error) {
	files, err := filepath.Glob(glob)
	if err != nil {
		return nil, err
	}

	var img []*image.RGBA
	for _, name := range files {
		m, err := LoadRGBAFile(name)
		if err == nil {
			img = append(img, m)
		}
	}
	return img, nil
}

func LoadRGBAFS(fs fs.FS, name string) (*image.RGBA, error) {
	f, err := fs.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return LoadRGBAReader(f)
}

func LoadRGBAFile(name string) (*image.RGBA, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	m, err := LoadRGBAReader(f)
	if err != nil {
		if strings.HasSuffix(strings.ToLower(name), "tga") {
			f.Seek(0, io.SeekStart)
			m, xerr := tga.Decode(f)
			if xerr == nil {
				return m.(*image.RGBA), nil
			}
		}

		return nil, &os.PathError{Op: "decode", Path: name, Err: err}
	}
	return m, nil
}

func LoadRGBAReader(rd io.Reader) (*image.RGBA, error) {
	m, _, err := image.Decode(rd)
	if err != nil {
		return nil, err
	}

	if p, _ := m.(*image.RGBA); p != nil {
		return p, nil
	}

	r := m.Bounds()
	p := image.NewRGBA(r)
	draw.Draw(p, p.Bounds(), m, r.Min, draw.Src)
	return p, nil
}

func LoadGrayFile(name string) (*image.Gray, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	m, err := LoadGrayReader(f)
	if err != nil {
		return nil, &os.PathError{Op: "decode", Path: name, Err: err}
	}
	return m, nil
}

func LoadGrayReader(rd io.Reader) (*image.Gray, error) {
	m, _, err := image.Decode(rd)
	if err != nil {
		return nil, err
	}

	if p, _ := m.(*image.Gray); p != nil {
		return p, nil
	}

	r := m.Bounds()
	p := image.NewGray(r)
	draw.Draw(p, p.Bounds(), m, r.Min, draw.Src)
	return p, nil
}

func LoadFloatFile(name string) (*Float, error) {
	m, err := LoadRGBAFile(name)
	if err != nil {
		return nil, err
	}
	return ImageToFloat(m), nil
}

func LoadFloatReader(rd io.Reader) (*Float, error) {
	m, err := LoadRGBAReader(rd)
	if err != nil {
		return nil, err
	}
	return ImageToFloat(m), nil
}

func WriteRGBAFile(name string, img image.Image) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}

	ext := filepath.Ext(name)
	switch strings.ToLower(ext) {
	case ".jpg", ".jpeg":
		err = jpeg.Encode(f, img, &jpeg.Options{Quality: 100})
	case ".pbm":
		err = pnm.Encode(f, img, &pnm.Options{Format: 1})
	case ".pgm":
		err = pnm.Encode(f, img, &pnm.Options{Format: 2})
	case ".ppm":
		err = pnm.Encode(f, img, &pnm.Options{Format: 3})
	case ".gif":
		err = gif.Encode(f, img, &gif.Options{
			NumColors: 256,
		})
	case ".tga":
		err = tga.Encode(f, img)
	case ".bmp":
		err = bmp.Encode(f, img)
	case ".png":
		fallthrough
	default:
		err = png.Encode(f, img)
	}

	xerr := f.Close()
	if err == nil {
		err = xerr
	}
	return err
}

func ColorKey(m image.Image, c color.Color) *image.RGBA {
	p := image.NewRGBA(m.Bounds())
	b := p.Bounds()

	cr, cg, cb, _ := c.RGBA()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			col := m.At(x, y)
			r, g, b, _ := col.RGBA()
			if cr == r && cg == g && cb == b {
				p.Set(x, y, color.Transparent)
			} else {
				p.Set(x, y, col)
			}
		}
	}
	return p
}

type CompareOption struct {
	Distance   func(a, b color.Color) float64
	Threshold  float64
	SubRegionA *image.Rectangle
	SubRegionB *image.Rectangle
	SameSize   bool
}

func Equals(a, b image.Image, o *CompareOption) bool {
	if o == nil {
		o = &CompareOption{
			Distance:  chroma.DistanceL2RGB,
			Threshold: 1e-3,
		}
	}

	r := a.Bounds()
	s := b.Bounds()

	if o.SubRegionA != nil {
		r = *o.SubRegionA
	}
	if o.SubRegionB != nil {
		s = *o.SubRegionB
	}

	if o.SameSize && (r.Dx() != s.Dx() || r.Dy() != s.Dy()) {
		return false
	}

	for y := 0; y < r.Dy(); y++ {
		for x := 0; x < r.Dx(); x++ {
			ax := x + r.Min.X
			ay := y + r.Min.Y
			bx := x + s.Min.X
			by := y + s.Min.Y

			if !image.Pt(ax, ay).In(r) || !image.Pt(bx, by).In(s) {
				continue
			}

			u := a.At(ax, ay)
			v := b.At(bx, by)

			if o.Distance(u, v) > o.Threshold {
				return false
			}
		}
	}
	return true
}

func IsTransparent(m image.Image) bool {
	r := m.Bounds()
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			_, _, _, ca := m.At(x, y).RGBA()
			if ca != 0 {
				return false
			}
		}
	}
	return true
}

func clamp(x, a, b int) int {
	if x < a {
		x = a
	}
	if x > b {
		x = b
	}
	return x
}

func CombineRGBA(imgs ...*image.RGBA) *image.RGBA {
	var mw, mh int
	for _, m := range imgs {
		r := m.Bounds()
		if mw < r.Dx() {
			mw = r.Dx()
		}
		if mh < r.Dy() {
			mh = r.Dy()
		}
	}

	d := max(mw, mh)
	n := int(math.Ceil(math.Sqrt(float64(len(imgs)))))
	d = ga.Ceil2(d * n)

	mw, mh = 0, 0
	p := image.NewRGBA(image.Rect(0, 0, d, d))
	pt := image.ZP
	for i, m := range imgs {
		s := m.Bounds()
		r := p.Bounds()
		r.Min = pt
		draw.Draw(p, r, m, image.ZP, draw.Src)

		pt.X += s.Dx()
		mw = max(mw, pt.X)

		if i > 0 && i%n == 0 {
			pt.Y += s.Dy()
			pt.X = 0
			mh = max(mh, pt.Y)
		}
	}

	mw = ga.Ceil2(mw)
	mh = ga.Ceil2(mh)
	p = p.SubImage(image.Rect(0, 0, mw, mh)).(*image.RGBA)

	return p
}

func SplitRGBA(m *image.RGBA, r image.Rectangle, sp image.Point) []*image.RGBA {
	var imgs []*image.RGBA

	b := m.Bounds()
	for y := b.Min.Y + sp.Y; y < b.Max.Y; {
		for x := b.Min.X + sp.X; x < b.Max.X; {
			p := image.NewRGBA(r)
			draw.Draw(p, p.Bounds(), m, image.Pt(x, y), draw.Src)
			imgs = append(imgs, p)
			x += r.Dx()
		}
		y += r.Dy()
	}
	return imgs
}
