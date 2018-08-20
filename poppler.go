package poppler

// #cgo pkg-config: --cflags-only-I poppler-glib
// #include <poppler.h>
// #include <stdlib.h>
// #include <glib.h>
// #include <unistd.h>
import "C"

func Version() string {
	return C.GoString(C.poppler_get_version())
}
