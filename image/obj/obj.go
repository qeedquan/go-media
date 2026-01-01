package obj

import (
	"bufio"
	"fmt"
	"image"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/scanner"

	"github.com/qeedquan/go-media/image/imageutil"
	"github.com/qeedquan/go-media/math/ga"
)

type Model struct {
	Verts   []ga.Vec4d
	Coords  []ga.Vec4d
	Normals []ga.Vec4d
	Params  []ga.Vec4d
	Faces   [][3][3]int
	Lines   [][]int
	Mats    []Material
}

type Material struct {
	Name           string
	Colors         [3]ga.Vec3d
	SpecularFactor float64
	DissolveFactor float64
	Illumination   int
	Texture        struct {
		Displacement    *Texture
		Diffuse         *Texture
		Ambient         *Texture
		SpecularColor   *Texture
		SpecularHilight *Texture
		Alpha           *Texture
		Bump            *Texture
		Dissolve        *Texture
	}
}

type Texture struct {
	Blend       [2]bool
	MipmapBoost float64
	Origin      ga.Vec3d
	Scale       ga.Vec3d
	Turbulence  ga.Vec3d
	Clamp       bool
	BumpFactor  bool
	IMF         uint64
	Map         *image.RGBA
}

func LoadFile(name string) (*Model, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Load(name, f)
}

func Load(name string, r io.Reader) (*Model, error) {
	var err error
	m := &Model{}
	s := bufio.NewScanner(r)
	dir := filepath.Dir(name)
	for s.Scan() {
		line := s.Text()
		switch {
		case strings.HasPrefix(line, "v "):
			m.Verts = addVert(m.Verts, line, 1)
		case strings.HasPrefix(line, "vt "):
			m.Coords = addVert(m.Coords, line, 0)
		case strings.HasPrefix(line, "vn "):
			m.Normals = addVert(m.Normals, line, 0)
		case strings.HasPrefix(line, "vp "):
			m.Params = addVert(m.Params, line, 0)
		case strings.HasPrefix(line, "f "):
			m.Faces = addFace(m.Faces, line)
		case strings.HasPrefix(line, "l "):
			m.Lines = addLines(m.Lines, line)
		case strings.HasPrefix(line, "mtllib "):
			m.Mats, _ = addMat(dir, m.Mats, line)
		}
	}

	var mat Material
	texs := []**Texture{&mat.Texture.Diffuse, &mat.Texture.Bump, &mat.Texture.SpecularColor}
	files := []string{"diffuse", "nm_tangent", "spec"}
	exts := []string{".tga", ".png", ".jpg", ".jpeg", ".bmp", ".gif"}
	found := false
	for i, file := range files {
		for _, ext := range exts {
			base := filepath.Base(name)
			baseExt := filepath.Ext(base)
			filename := fmt.Sprintf("%s_%s%s", base[:len(base)-len(baseExt)], file, ext)
			(*texs[i]), err = loadTexture(dir, filename)
			if err == nil {
				found = true
				break
			}
		}
	}
	if found {
		m.Mats = append(m.Mats, mat)
	}

	return m, nil
}

func getFloats(line string) []float64 {
	var f []float64
	toks := strings.Split(line, " ")
	for i := 1; i < len(toks); i++ {
		v, _ := strconv.ParseFloat(toks[i], 64)
		f = append(f, v)
	}
	return f
}

func addVert(verts []ga.Vec4d, line string, w float64) []ga.Vec4d {
	p := [4]float64{0, 0, 0, w}
	copy(p[:], getFloats(line))
	return append(verts, ga.Vec4d{p[0], p[1], p[2], p[3]})
}

func addFace(faces [][3][3]int, line string) [][3][3]int {
	var f [][3]int
	toks := strings.Split(line, " ")
	for i := 1; i < len(toks); i++ {
		p := [3]int{-1, -1, -1}
		xtoks := strings.Split(toks[i], "/")
		for j := range p {
			if j >= len(xtoks) {
				break
			}
			n, _ := strconv.Atoi(xtoks[j])
			p[j] = n - 1
		}

		f = append(f, p)
	}

	for i := 0; i < len(f)-2; i++ {
		j := max(0, i-2)
		faces = append(faces, [3][3]int{
			{f[j][0], f[j][1], f[j][2]},
			{f[i+1][0], f[i+1][1], f[i+1][2]},
			{f[i+2][0], f[i+2][1], f[i+2][2]},
		})
	}

	return faces
}

func addMat(dir string, mat []Material, line string) ([]Material, error) {
	var (
		t string
		e string
	)
	n, _ := fmt.Sscan(line, &t, &e)
	if n != 2 {
		return mat, nil
	}

	f, err := os.Open(filepath.Join(dir, e))
	if err != nil {
		return mat, err
	}
	defer f.Close()

	p := make([]Material, 1)
	m := &p[0]

	s := bufio.NewScanner(f)
	for s.Scan() {
		line = s.Text()
		switch {
		case strings.HasPrefix(line, "newmtl "):
			var name string
			fmt.Sscan(line, &t, &name)
			if len(p) == 1 {
				m.Name = name
			} else {
				p = append(p, Material{Name: name})
				m = &p[len(p)-1]
			}
		case strings.HasPrefix(line, "illum "):
			fmt.Sscan(line, &t, &m.Illumination)
		case strings.HasPrefix(line, "Ka "):
			fmt.Sscan(line, &t, &m.Colors[0].X, &m.Colors[0].Y, &m.Colors[0].Z)
		case strings.HasPrefix(line, "Kd "):
			fmt.Sscan(line, &t, &m.Colors[1].X, &m.Colors[1].Y, &m.Colors[1].Z)
		case strings.HasPrefix(line, "Ks "):
			fmt.Sscan(line, &t, &m.Colors[2].X, &m.Colors[2].Y, &m.Colors[2].Z)
		case strings.HasPrefix(line, "Ns "):
			fmt.Sscan(line, &t, &m.SpecularFactor)
		case strings.HasPrefix(line, "d "):
			fmt.Sscan(line, &t, &m.DissolveFactor)
		case strings.HasPrefix(line, "map_disp ") || strings.HasPrefix(line, "disp "):
			m.Texture.Displacement, err = loadTexture(dir, line)
		case strings.HasPrefix(line, "map_bump ") || strings.HasPrefix(line, "bump "):
			m.Texture.Bump, err = loadTexture(dir, line)
		case strings.HasPrefix(line, "map_Ka "):
			m.Texture.Ambient, err = loadTexture(dir, line)
		case strings.HasPrefix(line, "map_Kd "):
			m.Texture.Diffuse, err = loadTexture(dir, line)
		case strings.HasPrefix(line, "map_Ks "):
			m.Texture.SpecularColor, err = loadTexture(dir, line)
		case strings.HasPrefix(line, "map_Ns "):
			m.Texture.SpecularHilight, err = loadTexture(dir, line)
		case strings.HasPrefix(line, "map_d "):
			m.Texture.Alpha, err = loadTexture(dir, line)
		}

		if err != nil {
			return mat, err
		}
	}

	if len(p) > 1 || m.Name != "" {
		mat = append(mat, p...)
	}

	return mat, nil
}

func addLines(lines [][]int, line string) [][]int {
	var l []int
	toks := strings.Split(line, " ")
	for i := 1; i < len(toks); i++ {
		n, _ := strconv.Atoi(toks[i])
		l = append(l, n)
	}
	return append(lines, l)
}

func loadTexture(dir, line string) (*Texture, error) {
	t := &Texture{
		Blend: [2]bool{true, true},
		Clamp: true,
	}

	var (
		s   scanner.Scanner
		err error
	)
	s.Init(strings.NewReader(line))
loop:
	for {
		tok := s.Scan()
		if tok == '-' {
			tok = s.Scan()
		}
		if tok != scanner.Ident {
			continue
		}

		v := s.TokenText()
		switch v {
		case "bm":
		case "clamp":
			t.Clamp = scanBool(&s)
		case "blendu":
			t.Blend[0] = scanBool(&s)
		case "blendv":
			t.Blend[1] = scanBool(&s)
		case "imfchan":
		case "mm":
		case "o":
			t.Origin.X = scanFloat(&s)
			t.Origin.Y = scanFloat(&s)
			t.Origin.Z = scanFloat(&s)
		case "s":
			t.Scale.X = scanFloat(&s)
			t.Scale.Y = scanFloat(&s)
			t.Scale.Z = scanFloat(&s)
		case "t":
			t.Turbulence.X = scanFloat(&s)
			t.Turbulence.Y = scanFloat(&s)
			t.Turbulence.Z = scanFloat(&s)
		case "texres":
		default:
			line = line[s.Position.Offset:]
			break loop

		}
	}

	if len(line) > 0 {
		filename := filepath.Join(dir, line)
		t.Map, err = imageutil.LoadRGBAFile(filename)
		if err != nil {
			return t, fmt.Errorf("%s: failed to load texture file: %v", line, err)
		}
	}

	return t, nil
}

func scanFloat(s *scanner.Scanner) float64 {
	s.Scan()
	n, _ := strconv.ParseFloat(s.TokenText(), 64)
	return n
}

func scanBool(s *scanner.Scanner) bool {
	s.Scan()
	switch strings.ToLower(s.TokenText()) {
	case "on", "true", "1":
		return true
	}
	return false
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
