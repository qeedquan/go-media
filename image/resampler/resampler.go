// ported from Separable filtering image rescaler v2.21, Rich Geldreich - richgel99@gmail.com
// https://github.com/richgel999/imageresampler
package resampler

// outline of how the thing works:
// its used to sample a 2D plane of values, ie, an image
// the two main function are putline and getline
// on init, the function will build a 2d table of all the weights
// to be convolved with the source for each pixel
// on every putline it waits it fills up the y source samples
// it waits until there is enough samples before it does a convolution
// it keeps a scanline buffer to hold the samples

import (
	"image"
	"math"

	"github.com/qeedquan/go-media/math/ga"
)

const (
	BOUNDARY_WRAP = iota
	BOUNDARY_REFLECT
	BOUNDARY_CLAMP
)

type contrib struct {
	Pixel  int
	Weight float64
}

type contribBound struct {
	Center      float64
	Left, Right int
}

type scanline struct {
	y int
	l []float64
}

type Options struct {
	BoundaryOp  int
	Filter      Filter
	SampleRange ga.Vec2d
	FilterScale ga.Vec2d
	SourceOff   ga.Vec2d
}

type Resampler struct {
	opt           *Options
	samples       []float64
	tmpSamples    []float64
	scanlines     []scanline
	pcx           [][]contrib
	pcy           [][]contrib
	ycount        []int
	yflag         []bool
	xdelay        bool
	xintermediate int
	dn, dc        image.Point
	sn, sc        image.Point
}

func New(dn, sn image.Point, opt *Options) *Resampler {
	if opt == nil {
		opt = &Options{
			BoundaryOp:  BOUNDARY_CLAMP,
			Filter:      GetFilter("blackman"),
			SampleRange: ga.Vec2d{0, 1},
			FilterScale: ga.Vec2d{1, 1},
			SourceOff:   ga.Vec2d{0, 0},
		}
	}

	r := &Resampler{
		opt:     opt,
		dn:      dn,
		sn:      sn,
		samples: make([]float64, dn.X),
	}
	r.pcx = r.makeList(dn.X, sn.X, opt.FilterScale.X, opt.SourceOff.X)
	r.pcy = r.makeList(dn.Y, sn.Y, opt.FilterScale.Y, opt.SourceOff.Y)

	r.ycount = make([]int, sn.Y)
	r.yflag = make([]bool, sn.Y)

	// count how many times each source line
	// contributes to destination line
	for i := range r.pcy {
		for j := range r.pcy[i] {
			p := r.pcy[i][j].Pixel
			r.ycount[p]++
		}
	}

	r.chooseSampleAxis()
	return r
}

// chooseSampleAxis chooses the order of which
// axis to sample first, based on the which axis
// has the smaller amount of multiply operations
func (r *Resampler) chooseSampleAxis() {
	xops := r.countOps(r.pcx)
	yops := r.countOps(r.pcy)

	// check which resample order is better
	// in case of a tie, choose order which buffers
	// least amount of data
	if xops > yops || (xops == yops && r.sn.X < r.dn.X) {
		r.xdelay = true
		r.xintermediate = r.sn.X
	} else {
		r.xdelay = false
		r.xintermediate = r.dn.X
	}

	if r.xdelay {
		r.tmpSamples = make([]float64, r.xintermediate)
	}
}

func (r *Resampler) Restart() {
	r.sc.Y = 0
	r.dc.Y = 0

	for i := range r.ycount {
		r.ycount[i] = 0
	}
	for i := range r.yflag {
		r.yflag[i] = false
	}

	for i := range r.pcy {
		for _, p := range r.pcy[i] {
			r.ycount[p.Pixel]++
		}
	}

	for i := range r.scanlines {
		s := &r.scanlines[i]
		s.y = -1
	}
}

// reflect ensures that contributing sample
// is within bounds, if not, clamp/wrap/reflect
// based on op
func (r *Resampler) reflect(x, w, op int) int {
	var n int
	switch {
	case x < 0:
		switch op {
		case BOUNDARY_REFLECT:
			n = -x
			if n >= w {
				n = w - 1
			}
		case BOUNDARY_WRAP:
			n = posmod(x, w)
		default:
			n = 0
		}

	case x >= w:
		switch op {
		case BOUNDARY_REFLECT:
			n = (w - x) + (w - 1)
			if n < 0 {
				n = 0
			}
		case BOUNDARY_WRAP:
			n = posmod(x, w)
		default:
			n = w - 1
		}

	default:
		n = x
	}
	return n
}

// makeList generates, for all destination samples,
// the list of all source samples with non-zero
// weighted contributions
func (r *Resampler) makeList(dn, sn int, filterScale float64, sourceOff float64) [][]contrib {
	const NUDGE = 0.5

	filter := r.opt.Filter
	contribs := make([][]contrib, dn)
	contribBounds := make([]contribBound, dn)

	scale := float64(dn) / float64(sn)
	ooFilterScale := 1 / filterScale
	halfWidth := filter.Support * filterScale
	scaleMul := 1.0
	if scale < 1 {
		scaleMul = scale
		halfWidth = (filter.Support / scale) * filterScale
	}

	// find source samples that contribute to each destination sample
	for i := 0; i < dn; i++ {
		// convert discrete to continuous, scale, and then back to discrete
		center := (float64(i) + NUDGE) / scale
		center -= NUDGE
		center += sourceOff

		left := math.Floor(center - halfWidth)
		right := math.Ceil(center + halfWidth)

		contribBounds[i] = contribBound{
			Center: center,
			Left:   int(left),
			Right:  int(right),
		}
	}

	// create the list of each source samples which
	// contribute to each destination sample
	totalWeight := 0.0
	for i := 0; i < dn; i++ {
		maxIndex := -1
		maxWeight := -1e20

		center := contribBounds[i].Center
		left := contribBounds[i].Left
		right := contribBounds[i].Right
		contribs[i] = make([]contrib, right-left+1)

		totalWeight = 0
		for j := left; j <= right; j++ {
			totalWeight += filter.Sample((center - float64(j)) * scaleMul * ooFilterScale)
		}
		norm := 1 / totalWeight

		totalWeight = 0
		index := 0
		for j := left; j <= right; j++ {
			weight := filter.Sample((center-float64(j))*scaleMul*ooFilterScale) * norm
			if weight == 0 {
				continue
			}
			contribs[i][index].Pixel = r.reflect(j, sn, r.opt.BoundaryOp)
			contribs[i][index].Weight = weight

			// increment the number of source samples which
			// contribute to the current destination sample
			totalWeight += weight
			if weight > maxWeight {
				maxWeight = weight
				maxIndex = index
			}
			index++
		}

		if totalWeight != 1 {
			contribs[i][maxIndex].Weight += 1 - totalWeight
		}
		contribs[i] = contribs[i][:index]
	}

	return contribs
}

func (r *Resampler) allocScanline() *scanline {
	for i := range r.scanlines {
		s := &r.scanlines[i]
		if s.y == -1 {
			return s
		}
	}
	r.scanlines = append(r.scanlines, scanline{y: -1})
	return &r.scanlines[len(r.scanlines)-1]
}

// resampleX convolves the x axis for a destination
// with source and weights
func (r *Resampler) resampleX(dst, src []float64) {
	for i := 0; i < r.dn.X; i++ {
		total := 0.0
		for _, p := range r.pcx[i] {
			total += src[p.Pixel] * p.Weight
		}
		dst[i] = total
	}
}

// PutLine puts the source samples onto a queue
// for later processing by GetLine
func (r *Resampler) PutLine(src []float64) bool {
	if r.sc.Y >= r.sn.Y {
		return false
	}

	// does this source line contribute to any
	// of the destination line? if not, exit early
	if r.ycount[r.sc.Y] == 0 {
		r.sc.Y++
		return true
	}

	// find empty slot in scanline buffer
	s := r.allocScanline()
	s.y = r.sc.Y
	s.grow(max(r.xintermediate, r.dn.X))
	r.yflag[r.sc.Y] = true

	// resample on x axis first?
	if r.xdelay {
		// y-x
		copy(s.l, src)
	} else {
		// x-y
		r.resampleX(s.l, src)
	}

	r.sc.Y++
	return true
}

// GetLine returns the output after convolving for
// a scanline
func (r *Resampler) GetLine() []float64 {
	// if all destination have been generated
	// always return nil
	if r.dc.Y >= r.dn.Y {
		return nil
	}

	// check to see if all the required contributors
	// are present, if not return nil
	for _, c := range r.pcy[r.dc.Y] {
		if !r.yflag[c.Pixel] {
			return nil
		}
	}

	r.resampleY(r.samples)
	r.dc.Y++
	return r.samples
}

// resampleY does a convolution operation
// for a scanline, it is called when all
// the contributors to a pixel are present
func (r *Resampler) resampleY(samples []float64) {
	tmp := r.samples
	if r.xdelay {
		tmp = r.tmpSamples
	}

	// process each contributor
	pc := r.pcy[r.dc.Y]
	for i := range pc {
		// locate the contributor location
		// in the scan buffer, must always
		// be found!
		var (
			src []float64
			s   *scanline
		)
		for j := range r.scanlines {
			s = &r.scanlines[j]
			if s.y == pc[i].Pixel {
				src = s.l
				break
			}
		}

		if i == 0 {
			r.scaleYMov(tmp, src, pc[i].Weight, r.xintermediate)
		} else {
			r.scaleYAdd(tmp, src, pc[i].Weight, r.xintermediate)
		}

		// if this source line doesn't contribute
		// anymore to the destination slots, mark as free
		// the max number of slots used depends
		// on the y axis sampling factor and scaled filter width
		p := pc[i].Pixel
		if r.ycount[p]--; r.ycount[p] == 0 {
			r.yflag[p] = false
			s.y = -1
		}
	}

	// now generate destination line

	// was x delayed until after y resampling?
	// if not we already placed the samples we want inside samples
	if r.xdelay {
		r.resampleX(samples, tmp)
	}

	// clamp if needed
	if r.opt.SampleRange.X < r.opt.SampleRange.Y {
		r.clampSamples(samples)
	}
}

// scaleYMov initializes a convolution operation against
// source and a weight value it sets it equal instead of adding
// to destination to initialize the table
func (r *Resampler) scaleYMov(dst, src []float64, weight float64, n int) {
	for i := 0; i < n; i++ {
		dst[i] = src[i] * weight
	}
}

// scaleYAdd does a convolution operation against source and
// a weight value
func (r *Resampler) scaleYAdd(dst, src []float64, weight float64, n int) {
	for i := 0; i < n; i++ {
		dst[i] += src[i] * weight
	}
}

// clampSamples clamps the sample points between user
// defined low and high
func (r *Resampler) clampSamples(samples []float64) {
	lo := r.opt.SampleRange.X
	hi := r.opt.SampleRange.Y
	for i := range samples {
		samples[i] = ga.Clamp(samples[i], lo, hi)
	}
}

// countOps gives a rough estimate of how many
// operation an axis takes so we can choose
func (r *Resampler) countOps(contribs [][]contrib) int {
	n := 0
	for i := range contribs {
		n += len(contribs[i])
	}
	return n
}

func (s *scanline) grow(n int) {
	if len(s.l) < n {
		s.l = make([]float64, n)
	}
}

// posmod computes x%y, it wraps around in
// the case where x is negative
func posmod(x, y int) int {
	if x >= 0 {
		return x % y
	}

	m := -x % y
	if m != 0 {
		m = y - m
	}
	return m
}
