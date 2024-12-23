package poppler

// #cgo pkg-config: poppler-glib
// #include <poppler.h>
// #include <glib.h>
// #include <cairo.h>
import "C"
import "unsafe"
import "github.com/ungerik/go-cairo"

//import "fmt"

type Page struct {
	p *C.struct__PopplerPage
	openedAnnots []*Annot
}

func (p *Page) Text() string {
	return C.GoString(C.poppler_page_get_text(p.p))
}

func (p *Page) TextAttributes() (results []TextAttributes) {
	a := C.poppler_page_get_text_attributes(p.p)
	defer C.poppler_page_free_text_attributes(a)
	var attr *C.PopplerTextAttributes
	results = make([]TextAttributes, 0)
	el := C.g_list_first(a)
	for el != nil {
		attr = (*C.PopplerTextAttributes)(el.data)
		fn := attr.font_name
		result := TextAttributes{
			FontName:     toString(fn),
			FontSize:     float64(attr.font_size),
			IsUnderlined: toBool(attr.is_underlined),
			StartIndex:   int(attr.start_index),
			EndIndex:     int(attr.end_index),
			Color: Color{
				R: int(attr.color.red),
				G: int(attr.color.green),
				B: int(attr.color.blue),
			},
		}
		results = append(results, result)
		el = el.next
	}
	return
}

func (p *Page) Size() (width, height float64) {
	var w, h C.double
	C.poppler_page_get_size(p.p, &w, &h)
	return float64(w), float64(h)
}

func (p *Page) Index() int {
	return int(C.poppler_page_get_index(p.p))
}

func (p *Page) Label() string {
	return toString(C.poppler_page_get_label(p.p))
}

func (p *Page) Duration() float64 {
	return float64(C.poppler_page_get_duration(p.p))
}

func (p *Page) Images() (results []Image) {
	l := C.poppler_page_get_image_mapping(p.p)
	defer C.poppler_page_free_image_mapping(l)
	results = make([]Image, 0)
	var im *C.PopplerImageMapping
	for el := C.g_list_first(l); el != nil; el = el.next {
		im = (*C.PopplerImageMapping)(el.data)
		result := Image{
			Id: int(im.image_id),
			Area: Rectangle{
				X1: float64(im.area.x1),
				Y1: float64(im.area.y1),
				X2: float64(im.area.x2),
				Y2: float64(im.area.y2),
			},
			p: p.p,
		}
		results = append(results, result)
	}
	return
}

func (p *Page) TextLayout() (layouts []Rectangle) {
	var rect *C.PopplerRectangle
	var n C.guint
	if toBool(C.poppler_page_get_text_layout(p.p, &rect, &n)) {
		defer C.g_free((C.gpointer)(rect))
		layouts = make([]Rectangle, int(n))
		r := (*[1 << 30]C.PopplerRectangle)(unsafe.Pointer(rect))[:n:n]
		for i := 0; i < int(n); i++ {
			layouts[i] = Rectangle{
				X1: float64(r[i].x1),
				Y1: float64(r[i].y1),
				X2: float64(r[i].x2),
				Y2: float64(r[i].y2),
			}
		}
	}
	return
}

func (p *Page) TextLayoutAndAttrs() (result []TextEl) {
	text := p.Text()
	attrs := p.TextAttributes()
	layout := p.TextLayout()
	result = make([]TextEl, len(layout))
	attrsRef := make([]*TextAttributes, len(attrs))
	for i, a := range attrs {
		attr := a
		attrsRef[i] = &attr
	}
	i := 0
	for _, t := range text {
		var a *TextAttributes
		for _, a = range attrsRef {
			if i >= a.StartIndex && i <= a.EndIndex {
				break
			}
		}
		result[i] = TextEl{
			Text:  string(t),
			Attrs: a,
			Rect:  layout[i],
		}
		i++
	}
	return
}

func (p *Page) Close() {
	p.closeAnnotMappings()

	if p.p != nil {
		C.g_object_unref(C.gpointer(p.p))
		/* avoid double free */
		p.p = nil
	}
}

// Converts a page into SVG and saves to file.
// Inspired by https://github.com/dawbarton/pdf2svg
func (p *Page) ConvertToSVG(filename string){
	width, height := p.Size()

	// Open the SVG file
	surface := cairo.NewSVGSurface( filename, width, height, cairo.SVG_VERSION_1_2 )

	// TODO Can be improved by using cairo_svg_surface_create_for_stream() instead of
	//      cairo_svg_surface_create() for stream processing instead of file processing.
	//      However, this needs to be changed in github.com/ungerik/go-cairo/surface.go

	// Get cairo context pointer
	_, drawcontext :=  surface.Native()

	// Render the PDF file into the SVG file
	C.poppler_page_render_for_printing(p.p, (*C.cairo_t)(unsafe.Pointer(drawcontext)) );

	// Close the SVG file
	surface.ShowPage()
	surface.Destroy()
}

func (p *Page) closeAnnotMappings(){
	for i := 0; i < len(p.openedAnnots); i++ {
		p.openedAnnots[i].Close()
	}

	p.openedAnnots = nil

}

func (p *Page) GetAnnots() (Annots []*Annot) {
	var annots []*Annot

	annotGlist := C.poppler_page_get_annot_mapping(p.p)
	defer C.g_list_free(annotGlist)

	p.closeAnnotMappings()

	for annotGlist != nil {
		popplerAnnot := (*C.PopplerAnnotMapping)(annotGlist.data)


		annot := &Annot{
			am: popplerAnnot,
		}

		/* Maybe we can used openedAnnots instead of annots + openedAnnots
		 */

		annots = append(annots, annot)
		p.openedAnnots = append(p.openedAnnots, annot)


		annotGlist = annotGlist.next
	}

	return annots
}

func (p *Page) AnnotText(a Annot) string {
	cText := C.poppler_page_get_text_for_area(p.p, &a.am.area)
	return C.GoString(cText)
}

func (p *Page) AddAnnot(a Annot) {
	C.poppler_page_add_annot(p.p, a.am.annot)
}
func (p *Page) RemoveAnnot(a Annot) {
	C.poppler_page_remove_annot(p.p, a.am.annot)
}
