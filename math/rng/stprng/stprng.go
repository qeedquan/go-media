package stprng

type Counter struct {
	I, Is int64
	F, Fs float64
}

func (c *Counter) Int() int {
	v := c.I
	c.I += c.Is
	return int(v)
}

func (c *Counter) Float64() float64 {
	v := c.F
	c.F += c.Fs
	return v
}

type Table struct {
	I      []int
	F      []float64
	Ix, Fx int
}

func (c *Table) Int() int {
	if len(c.I) == 0 {
		return 0
	}
	v := c.I[c.Ix]
	c.Ix = (c.Ix + 1) % len(c.I)
	return v
}

func (c *Table) Float64() float64 {
	if len(c.F) == 0 {
		return 0
	}
	v := c.F[c.Fx]
	c.Fx = (c.Fx + 1) % len(c.F)
	return v
}
