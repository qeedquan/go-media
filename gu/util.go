package gu

import "log"

func Assert(x bool) {
	if !x {
		panic("assertion failed")
	}
}

func Check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
