package poppler

// #cgo pkg-config: poppler-glib
// #include <poppler.h>
// #include <glib.h>
// #include <cairo.h>
import "C"
//import "unsafe"
//import "github.com/ungerik/go-cairo"

//import "fmt"

type Annot struct {
	am *C.struct__PopplerAnnotMapping 
}

func (a *Annot) Close() {
	C.poppler_annot_mapping_free((*C.struct__PopplerAnnotMapping)(a.am))
}
