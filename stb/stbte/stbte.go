package stbte

/*
#include "stbte.h"
#define STB_TEXTEDIT_IMPLEMENTATION
#include "stb_textedit.h"
*/
import "C"
import (
	"sync"
	"unsafe"
)

const (
	K_START     = C.STB_TEXTEDIT_K_START
	K_LEFT      = C.STB_TEXTEDIT_K_LEFT
	K_RIGHT     = C.STB_TEXTEDIT_K_RIGHT
	K_UP        = C.STB_TEXTEDIT_K_UP
	K_DOWN      = C.STB_TEXTEDIT_K_DOWN
	K_LINESTART = C.STB_TEXTEDIT_K_LINESTART
	K_LINEEND   = C.STB_TEXTEDIT_K_LINEEND
	K_TEXTSTART = C.STB_TEXTEDIT_K_TEXTSTART
	K_TEXTEND   = C.STB_TEXTEDIT_K_TEXTEND
	K_DELETE    = C.STB_TEXTEDIT_K_DELETE
	K_BACKSPACE = C.STB_TEXTEDIT_K_BACKSPACE
	K_UNDO      = C.STB_TEXTEDIT_K_UNDO
	K_REDO      = C.STB_TEXTEDIT_K_REDO
	K_WORDLEFT  = C.STB_TEXTEDIT_K_WORDLEFT
	K_WORDRIGHT = C.STB_TEXTEDIT_K_WORDRIGHT
	K_PGUP      = C.STB_TEXTEDIT_K_PGUP
	K_PGDOWN    = C.STB_TEXTEDIT_K_PGDOWN
	K_SHIFT     = C.STB_TEXTEDIT_K_SHIFT
)

type (
	State       C.STB_TexteditState
	TextEditRow C.StbTexteditRow
)

type String interface {
	Len() int
	GetChar(idx int) rune
	GetWidth(line_start_idx, char_idx int) float64
	LayoutRow(r *TextEditRow, line_start_idx int)
	MoveWordRight(idx int) int
	MoveWordLeft(idx int) int
	DeleteChars(pos, n int)
	InsertChars(pos int, new_text []rune) bool
}

type context struct {
	sync.Mutex
	str String
}

var ctx context

func (r *TextEditRow) X0() float64 {
	return float64(r.x0)
}

func (r *TextEditRow) X1() float64 {
	return float64(r.x1)
}

func (r *TextEditRow) YMin() float64 {
	return float64(r.ymin)
}

func (r *TextEditRow) YMax() float64 {
	return float64(r.ymax)
}

func (r *TextEditRow) BaselineYDelta() float64 {
	return float64(r.baseline_y_delta)
}

func (r *TextEditRow) SetX0(x0 float64) {
	r.x0 = C.float(x0)
}

func (r *TextEditRow) SetX1(x1 float64) {
	r.x1 = C.float(x1)
}

func (r *TextEditRow) SetYMin(ymin float64) {
	r.ymin = C.float(ymin)
}

func (r *TextEditRow) SetYMax(ymax float64) {
	r.ymax = C.float(ymax)
}

func (r *TextEditRow) SetBaselineYDelta(baseline_y_delta float64) {
	r.baseline_y_delta = C.float(baseline_y_delta)
}

func (r *TextEditRow) NumChars() int {
	return int(r.num_chars)
}

func (r *TextEditRow) SetNumChars(num_chars int) {
	r.num_chars = C.int(num_chars)
}

func (s *State) Init(is_single_line bool) {
	C.stb_textedit_initialize_state((*C.STB_TexteditState)(s), truth(is_single_line))
}

func (s *State) Click(str String, x, y float64) {
	ctx.Lock()
	ctx.str = str
	defer ctx.Unlock()
	C.stb_textedit_click(nil, (*C.STB_TexteditState)(s), C.float(x), C.float(y))
}

func (s *State) Drag(str String, x, y float64) {
	ctx.Lock()
	ctx.str = str
	defer ctx.Unlock()
	C.stb_textedit_drag(nil, (*C.STB_TexteditState)(s), C.float(x), C.float(y))
}

func (s *State) Cut(str String) {
	ctx.Lock()
	ctx.str = str
	defer ctx.Unlock()
	C.stb_textedit_cut(nil, (*C.STB_TexteditState)(s))
}

func (s *State) Paste(str String, text []rune) int {
	ctx.Lock()
	ctx.str = str
	defer ctx.Unlock()
	return int(C.stb_textedit_paste(nil, (*C.STB_TexteditState)(s), (*C.int)(unsafe.Pointer(&text[0])), C.int(len(text))))
}

func (s *State) Key(str String, key int) {
	ctx.Lock()
	ctx.str = str
	defer ctx.Unlock()
	C.stb_textedit_key(nil, (*C.STB_TexteditState)(s), C.int(key))
}

func (s *State) Cursor() int {
	return int(s.cursor)
}

func (s *State) SetCursor(cursor int) {
	s.cursor = C.int(cursor)
}

func (s *State) SelectStart() int {
	return int(s.select_start)
}

func (s *State) SetSelectStart(select_start int) {
	s.select_start = C.int(select_start)
}

func (s *State) SelectEnd() int {
	return int(s.select_end)
}

func (s *State) SetSelectEnd(select_end int) {
	s.select_end = C.int(select_end)
}

func (s *State) InsertMode() bool {
	return s.insert_mode != 0
}

func (s *State) SetInsertMode(insert_mode bool) {
	s.insert_mode = C.uchar(truth(insert_mode))
}

func (s *State) PreferredX() float64 {
	return float64(s.preferred_x)
}

func (s *State) SetPreferredX(preferred_x float64) {
	s.preferred_x = C.float(preferred_x)
}

func (s *State) HasPreferredX() bool {
	return s.has_preferred_x != 0
}

func (s *State) SetHasPreferredX(has_preferred_x bool) {
	s.has_preferred_x = C.uchar(truth(has_preferred_x))
}

func truth(cond bool) C.int {
	if cond {
		return 1
	}
	return 0
}

//export stringlen
func stringlen(unsafe.Pointer) int {
	p := ctx.str
	return p.Len()
}

//export getchar
func getchar(_ unsafe.Pointer, idx int) rune {
	p := ctx.str
	return p.GetChar(idx)
}

//export getwidth
func getwidth(_ unsafe.Pointer, line_start_idx, char_idx int) float64 {
	p := ctx.str
	return p.GetWidth(line_start_idx, char_idx)
}

//export layoutrow
func layoutrow(r *TextEditRow, _ unsafe.Pointer, line_start_idx int) {
	p := ctx.str
	p.LayoutRow(r, line_start_idx)
}

//export movewordright
func movewordright(_ unsafe.Pointer, idx int) int {
	p := ctx.str
	return p.MoveWordRight(idx)
}

//export movewordleft
func movewordleft(_ unsafe.Pointer, idx int) int {
	p := ctx.str
	return p.MoveWordLeft(idx)
}

//export insertchars
func insertchars(_ unsafe.Pointer, pos int, text *rune, new_text_len int) bool {
	s := ((*[1 << 27]rune)(unsafe.Pointer(text)))[:new_text_len:new_text_len]
	p := ctx.str
	return p.InsertChars(pos, s)
}

//export deletechars
func deletechars(_ unsafe.Pointer, pos, n int) {
	p := ctx.str
	p.DeleteChars(pos, n)
}
