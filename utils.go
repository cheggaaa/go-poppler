package poppler

// #include <glib.h>
// #include <unistd.h>
// #include <stdlib.h>
import "C"

func toString(in *C.gchar) string {
	return C.GoString((*C.char)(in))
}

func toBool(in C.gboolean) bool {
	return  int(in) > 0
}