package poppler

// #cgo pkg-config: poppler-glib
// #include <poppler.h>
// #include <glib.h>
// #include <cairo.h>
//
// /* macro wrappings */
// gboolean wrap_POPPLER_IS_ANNOT_TEXT_MARKUP(PopplerAnnot *annot){
//   return POPPLER_IS_ANNOT_TEXT_MARKUP(annot);
// }
// PopplerAnnotTextMarkup *wrap_POPPLER_ANNOT_TEXT_MARKUP(PopplerAnnot *annot) {
//	return POPPLER_ANNOT_TEXT_MARKUP(annot);
// }
import "C"

import "unsafe"
//import "github.com/ungerik/go-cairo"

//import "fmt"

type Point struct {
	X, Y float64
}
type Quad struct {
	P1, P2, P3, P4 Point
}

type Annot struct {
	am *C.struct__PopplerAnnotMapping
}

type AnnotType int

const (
	AnnotUnknown AnnotType = iota
	AnnotText
	AnnotLink
	AnnotFreeText
	AnnotLine
	AnnotSquare
	AnnotCircle
	AnnotPolygon
	AnnotPolyLine
	AnnotHighlight
	AnnotUnderline
	AnnotSquiggly
	AnnotStrikeOut
	AnnotStamp
	AnnotCaret
	AnnotInk
	AnnotPopup
	AnnotFileAttachment
	AnnotSound
	AnnotMovie
	AnnotWidget
	AnnotScreen
	AnnotPrinterMark
	AnnotTrapNet
	AnnotWatermark
	Annot3D
)

type AnnotFlag int

const (
	AnnotFlagUnknown AnnotFlag = 1 << iota
	AnnotFlagInvisible
	AnnotFlagHidden
	AnnotFlagPrint
	AnnotFlagNoZoom
	AnnotFlagNoRotate
	AnnotFlagNoView
	AnnotFlagReadOnly
	AnnotFlagLocked
	AnnotFlagToggleNoView
	AnnotFlagLockedContents
)

func (a *Annot) Type() AnnotType {
	t := C.poppler_annot_get_annot_type(a.am.annot)
	return AnnotType(t)
}

func (a *Annot) Index() int {
	i := C.poppler_annot_get_page_index(a.am.annot)
	return int(i)
}

func (a *Annot) Date() string {
	cText := C.poppler_annot_get_modified(a.am.annot)
	return C.GoString(cText)
}

func (a *Annot) Rect() Rectangle {
	var r C.PopplerRectangle
	C.poppler_annot_get_rectangle(a.am.annot, &r)

	rect := Rectangle{
		X1: float64(r.x1),
		Y1: float64(r.y1),
		X2: float64(r.x2),
		Y2: float64(r.y2),
	}

	return rect

}

func (a *Annot) Color() Color {
	c := C.poppler_annot_get_color(a.am.annot)
	defer C.poppler_color_free(c)

	color := Color{
		R: int(c.red),
		G: int(c.green),
		B: int(c.blue),
	}

	return color
}

func (a *Annot) Name() string {
	cText := C.poppler_annot_get_name(a.am.annot)
	return C.GoString(cText)
}

func (a *Annot) Contents() string {
	cText := C.poppler_annot_get_contents(a.am.annot)
	return C.GoString(cText)
}

func (a *Annot) Flags() AnnotFlag {
	f := C.poppler_annot_get_flags(a.am.annot)
	return AnnotFlag(f)
}

func (a *Annot) Quads() []Quad {
	if C.wrap_POPPLER_IS_ANNOT_TEXT_MARKUP(a.am.annot) == C.FALSE {
		return nil
	}


	textMarkup := C.wrap_POPPLER_ANNOT_TEXT_MARKUP(a.am.annot)

	q := C.poppler_annot_text_markup_get_quadrilaterals(textMarkup)

	quads := gArrayToQuads(q)

	C.g_array_free(q, 1)

	return quads
}

func (a *Annot) Close() {
	if a.am != nil {
		C.poppler_annot_mapping_free(a.am)
		a.am = nil
	}
}

func (a *Annot) SetColor(c Color){
	pColor := C.poppler_color_new()
	pColor.red = C.ushort(c.R)
	pColor.green = C.ushort(c.G)
	pColor.blue = C.ushort(c.B)
	defer gFree(pColor)

	C.poppler_annot_set_color(a.am.annot, pColor )
}

func (a *Annot) SetContents(c string){
	cStr := C.CString(c)
	defer C.free(unsafe.Pointer(cStr))

	C.poppler_annot_set_contents(a.am.annot, cStr)
}

func  (a *Annot) SetFlags(f AnnotFlag){
	pFlags := C.PopplerAnnotFlag(f)

	C.poppler_annot_set_flags(a.am.annot, pFlags)
}
