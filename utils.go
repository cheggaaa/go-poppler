package poppler

// #cgo pkg-config: poppler-glib
// #include <poppler.h>
// #include <glib.h>
// #include <unistd.h>
// #include <stdlib.h>
import "C"

import "unsafe"

func toString(in *C.gchar) string {
	return C.GoString((*C.char)(in))
}

func toBool(in C.gboolean) bool {
	return  int(in) > 0
}

/* convert a Quad struct to a GArray */
func quadsToGArray(quads []Quad) *C.GArray {
	garray := C.g_array_new(C.FALSE, C.FALSE, C.sizeof_PopplerQuadrilateral)

	for _, quad := range quads {
		item := C.PopplerQuadrilateral{
			p1: C.PopplerPoint{
				x: C.double(quad.P1.X),
				y: C.double(quad.P1.Y),
			},
			p2: C.PopplerPoint{
				x: C.double(quad.P2.X),
				y: C.double(quad.P2.Y),
			},
			p3: C.PopplerPoint{
				x: C.double(quad.P3.X),
				y: C.double(quad.P3.Y),
			},
			p4: C.PopplerPoint{
				x: C.double(quad.P4.X),
				y: C.double(quad.P4.Y),
			},
		}

		C.g_array_append_vals(garray, C.gconstpointer(&item),1)
	}

	return garray
}

/* convert a GArray to a quad */
func gArrayToQuads(q *C.GArray) []Quad {
	length := int(q.len)

	quads := make([]Quad, length)

	for i := 0; i < length; i++ {
		item := (*C.PopplerQuadrilateral)(unsafe.Pointer(uintptr(unsafe.Pointer(q.data)) + uintptr(i)*unsafe.Sizeof(C.PopplerQuadrilateral{})))
		quads[i] = Quad{
			P1: Point{X: float64(item.p1.x), Y: float64(item.p1.y)},
			P2: Point{X: float64(item.p2.x), Y: float64(item.p2.y)},
			P3: Point{X: float64(item.p3.x), Y: float64(item.p3.y)},
			P4: Point{X: float64(item.p4.x), Y: float64(item.p4.y)},
		}
	}

	return quads
}

func rectangleToPopplerRectangle (r Rectangle) C.PopplerRectangle {
	var pRect C.PopplerRectangle

	pRect.x1 = C.double(r.X1)
	pRect.y1 = C.double(r.Y1)
	pRect.x2 = C.double(r.X2)
	pRect.y2 = C.double(r.Y2)

	return pRect
}
