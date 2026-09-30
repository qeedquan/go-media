package tiled

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"encoding/base64"
	"encoding/binary"
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"image"
	"io"
	"io/fs"
	"sort"
	"strconv"
	"strings"

	"github.com/qeedquan/go-media/image/imageutil"
)

const (
	ORTHOGONAL = iota
	ISOMETRIC
	HEXAGONAL
)

type Map struct {
	Sets        []*Set
	Layers      []*Layer
	Orientation int
	Width       int
	Height      int
	TileWidth   int
	TileHeight  int
}

type Set struct {
	Image      *image.RGBA
	TileWidth  int
	TileHeight int
}

type Layer struct {
	Tiles []Tile
}

type Tile struct {
	ID int
}

type TMX struct {
	XMLName         xml.Name `xml:"map"`
	Version         string   `xml:"version,attr"`
	TiledVersion    string   `xml:"tiledversion,attr"`
	Orientation     string   `xml:"orientation,attr"`
	RenderOrder     string   `xml:"renderorder,attr"`
	Width           int      `xml:"width,attr"`
	Height          int      `xml:"height,attr"`
	TileWidth       int      `xml:"tilewidth,attr"`
	TileHeight      int      `xml:"tileheight,attr"`
	BackgroundColor string   `xml:"backgroundcolor,attr"`
	NextObjectID    int      `xml:"nextobjectid,attr"`
	Tileset         []TSX    `xml:"tileset"`
	Layer           []TLY    `xml:"layer"`
}

type TSX struct {
	XMLName    xml.Name `xml:"tileset"`
	FirstGID   int      `xml:"firstgid,attr"`
	Source     string   `xml:"source,attr"`
	Name       string   `xml:"name,attr"`
	TileWidth  int      `xml:"tilewidth,attr"`
	TileHeight int      `xml:"tileheight,attr"`
	TileCount  int      `xml:"tilecount,attr"`
	Columns    int      `xml:"columns,attr"`
	Margin     int      `xml:"margin,attr"`
	Spacing    int      `xml:"spacing,attr"`
	Image      struct {
		Source string `xml:"source,attr"`
		Trans  string `xml:"trans,attr"`
		Width  int    `xml:"width,attr"`
		Height int    `xml:"height,attr"`
	} `xml:"image"`
}

type TLY struct {
	XMLName xml.Name `xml:"layer"`
	Name    string   `xml:"name,attr"`
	Width   int      `xml:"width,attr"`
	Height  int      `xml:"height,attr"`
	Visible *int     `xml:"visible,attr"`
	Data    struct {
		Encoding    string `xml:"encoding,attr"`
		Compression string `xml:"compression,attr"`
		Tile        []struct {
			GID int `xml:"gid,attr"`
		} `xml:"tile"`
		Chardata string `xml:",chardata"`
	} `xml:"data"`
}

func OpenMap(fs fs.FS, name string) (*Map, error) {
	d := decoder{
		fs: fs,
		m:  &Map{},
	}
	err := d.decode(name)
	if err != nil {
		return nil, fmt.Errorf("tiled: %v", err)
	}

	return d.m, nil
}

type decoder struct {
	fs fs.FS
	tm TMX
	m  *Map
}

func (d *decoder) decode(name string) error {
	err := d.decodeXML(name, &d.tm)
	if err != nil {
		return err
	}

	switch s := strings.ToLower(d.tm.Orientation); s {
	case "orthogonal":
		d.m.Orientation = ORTHOGONAL
	case "isometric":
		d.m.Orientation = ISOMETRIC
	case "hexagonal":
		d.m.Orientation = HEXAGONAL
	default:
		return fmt.Errorf("unsupported orientation %q", s)
	}

	switch s := strings.ToLower(d.tm.RenderOrder); s {
	case "right-down",
		"right-up",
		"left-down",
		"left-up":
	default:
		return fmt.Errorf("unsupported render order %q", s)
	}

	d.m.Width = d.tm.Width
	d.m.Height = d.tm.Height
	d.m.TileWidth = d.tm.TileWidth
	d.m.TileHeight = d.tm.TileHeight

	sort.Slice(d.tm.Tileset, func(i, j int) bool {
		return d.tm.Tileset[i].FirstGID < d.tm.Tileset[j].FirstGID
	})
	for i := range d.tm.Tileset {
		s, err := d.decodeTSX(&d.tm.Tileset[i])
		if err != nil {
			return err
		}
		d.m.Sets = append(d.m.Sets, s)
	}

	for i := range d.tm.Layer {
		d.decodeTLY(&d.tm.Layer[i])
	}

	return nil
}

func (d *decoder) decodeTSX(ts *TSX) (*Set, error) {
	if ts.Source != "" {
		err := d.decodeXML(ts.Source, ts)
		if err != nil {
			return nil, err
		}
	}

	var err error
	s := &Set{
		TileWidth:  ts.TileWidth,
		TileHeight: ts.TileHeight,
	}
	s.Image, err = imageutil.LoadRGBAFS(d.fs, ts.Image.Source)
	if err != nil {
		return nil, err
	}

	return s, err
}

func (d *decoder) decodeTLY(tl *TLY) (*Layer, error) {
	var t []int

	c := &tl.Data
	l := &Layer{}
	switch c.Encoding {
	case "base64":
		var buf []byte

		c.Chardata = strings.Trim(c.Chardata, " \r\n")
		buf, err := base64.StdEncoding.DecodeString(c.Chardata)
		if err != nil {
			return nil, err
		}
		br := bufio.NewReader(bytes.NewBuffer(buf))
		var cr io.Reader
		switch c.Compression {
		case "gzip":
			cr, err = gzip.NewReader(br)
		case "zlib":
			cr, err = zlib.NewReader(br)
		case "":
			cr = br
		default:
			return nil, fmt.Errorf("unknown tile compression %q", c.Compression)
		}
		if err != nil {
			return nil, err
		}

		var v uint32
		for i := 0; i < tl.Width*tl.Height; i++ {
			err = binary.Read(cr, binary.LittleEndian, &v)
			if err != nil {
				return nil, err
			}
			t = append(t, int(v))
		}

	case "csv":
		c.Chardata = strings.Map(func(r rune) rune {
			if strings.ContainsRune(" \t\r\n", r) {
				return -1
			}
			return r
		}, c.Chardata)

		r := csv.NewReader(bytes.NewBufferString(c.Chardata))
		sp, err := r.Read()
		if err != nil {
			return nil, err
		}

		for i := range sp {
			v, err := strconv.ParseInt(sp[i], 0, 64)
			if err != nil {
				return nil, err
			}
			t = append(t, int(v))
		}

	case "":
		for _, p := range c.Tile {
			t = append(t, p.GID)
		}

	default:
		return nil, fmt.Errorf("unknown tile encoding %q", c.Encoding)
	}

	if len(t) != tl.Width*tl.Height {
		return nil, fmt.Errorf("unexpected EOF reading tiles, got %d, expected %d", len(t), tl.Width*tl.Height)
	}

	return l, nil
}

func (d *decoder) decodeXML(name string, v interface{}) error {
	buf, err := fs.ReadFile(d.fs, name)
	if err != nil {
		return err
	}
	return xml.Unmarshal(buf, v)
}
