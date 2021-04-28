package poppler

// #cgo pkg-config: poppler-glib
// #include <poppler.h>
// #include <glib.h>
import "C"
import (
	"unsafe"

	"github.com/ungerik/go-cairo"
)

// Image

type Image struct {
	Id   int
	Area Rectangle
	p    *C.struct__PopplerPage
}

type Rectangle struct {
	X1, Y1, X2, Y2 float64
}

func (im *Image) GetSurface() (cs *cairo.Surface) {
	ci := C.poppler_page_get_image(im.p, C.gint(im.Id))
	ctx := C.cairo_create(ci)
	cip := (cairo.Cairo_surface)(unsafe.Pointer(ci))
	ctxp := (cairo.Cairo_context)(unsafe.Pointer(ctx))
	return cairo.NewSurfaceFromC(cip, ctxp)
}
