package ttf

type Family map[string][]byte

var VGA437 = Family{
	"default": vga437,
}

var ProggyClean = Family{
	"default": proggyClean,
}

var Roboto = Family{
	"default": robotoNormal,
	"light":   robotoLight,
	"bold":    robotoBold,
}

var DejaVuSerif = Family{
	"default": dejaVuSerif,
}

var DejaVuSans = Family{
	"default": dejaVuSans,
}
