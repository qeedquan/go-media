package debug

import (
	"fmt"
	"runtime/debug"
)

type DiffReadWriter struct {
	Data  []byte
	Off   int
	Stack [][]byte
}

func (c *DiffReadWriter) Write(b []byte) (int, error) {
	n := min(len(c.Data)-c.Off, len(b))
	r := fmt.Sprintf("WRITE(OFFSET=%d, LENGTH=%d)\n{\n%s}\n\n", c.Off, n, debug.Stack())
	if n == 0 {
		c.Stack = append(c.Stack, []byte(r))
		return n, nil
	}

	for i := range b {
		if c.Data[c.Off+i] != b[i] {
			c.Stack = append(c.Stack, []byte(r))
			break
		}
	}
	c.Off += n

	return n, nil
}
