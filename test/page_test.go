package poppler

import (
	"testing"
	"fmt"

	"github.com/cheggaaa/go-poppler"
)

func TestCloseAnnots(t *testing.T) {
	doc, err := poppler.Open("test.pdf")
	defer doc.Close()
	if err != nil {
		return
	}


	n_pages := doc.GetNPages()

	for i := 0; i < n_pages; i++ {
		page := doc.GetPage(i)
		annots := page.GetAnnots()

		fmt.Println("page: ", i)
		for _, a := range annots {
			a.Close()
		}

		page.Close()
	}

}

func TestGetQuads(t *testing.T) {
	doc, err := poppler.Open("test.pdf")
	defer doc.Close()
	if err != nil {
		return
	}


	n_pages := doc.GetNPages()

	for i := 0; i < n_pages; i++ {
		page := doc.GetPage(i)
		annots := page.GetAnnots()

		for _, a := range annots {
			fmt.Println(a.Quads())
		}

		page.Close()
	}

}

func TestPageText(t *testing.T) {
	doc, err := poppler.Open("test.pdf")
	defer doc.Close()
	if err != nil {
		return
	}


	n_pages := doc.GetNPages()

	for i := 0; i < n_pages; i++ {
		page := doc.GetPage(i)

		fmt.Println(page.Text())

		page.Close()
	}

}

func TestAnnotName(t *testing.T) {
	doc, err := poppler.Open("test.pdf")
	defer doc.Close()
	if err != nil {
		return
	}


	n_pages := doc.GetNPages()

	for i := 0; i < n_pages; i++ {
		page := doc.GetPage(i)

		annots := page.GetAnnots()

		for _, a := range annots {
			if a.Type() == poppler.AnnotHighlight {
				fmt.Println(a.Name())
			}
		}


		page.Close()
	}

}

