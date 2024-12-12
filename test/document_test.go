package poppler

import (
	"testing"

	"github.com/cheggaaa/go-poppler"

)

func TestClosePageAndCloseDoc(t *testing.T) {
	doc, err := poppler.Open("test.pdf")
	defer doc.Close()
	if err != nil {
		return
	}


	n_pages := doc.GetNPages()

	for i := 0; i < n_pages; i++ {
		page := doc.GetPage(i)
		page.Close()
	}

}

