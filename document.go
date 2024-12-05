package poppler

// #cgo pkg-config: poppler-glib
// #include <poppler.h>
// #include <stdlib.h>
// #include <glib.h>
// #include <unistd.h>
import "C"

import (
	"errors"
	"unsafe"
	"path/filepath"
)

type Document struct {
	doc                poppDoc
	openedPopplerPages []*C.struct__PopplerPage
}

type DocumentInfo struct {
	PdfVersion, Title, Author, Subject, KeyWords, Creator, Producer, Metadata string
	CreationDate, ModificationDate, Pages                                     int
	IsLinearized                                                              bool
}

func (d *Document) Info() DocumentInfo {
	return DocumentInfo{
		PdfVersion:       toString(C.poppler_document_get_pdf_version_string(d.doc)),
		Title:            toString(C.poppler_document_get_title(d.doc)),
		Author:           toString(C.poppler_document_get_author(d.doc)),
		Subject:          toString(C.poppler_document_get_subject(d.doc)),
		KeyWords:         toString(C.poppler_document_get_keywords(d.doc)),
		Creator:          toString(C.poppler_document_get_creator(d.doc)),
		Producer:         toString(C.poppler_document_get_producer(d.doc)),
		Metadata:         toString(C.poppler_document_get_metadata(d.doc)),
		CreationDate:     int(C.poppler_document_get_creation_date(d.doc)),
		ModificationDate: int(C.poppler_document_get_modification_date(d.doc)),
		Pages:            int(C.poppler_document_get_n_pages(d.doc)),
		IsLinearized:     toBool(C.poppler_document_is_linearized(d.doc)),
	}
}

func (d *Document) GetNPages() int {
	return int(C.poppler_document_get_n_pages(d.doc))
}

func (d *Document) GetPage(i int) (page *Page) {
	p := C.poppler_document_get_page(d.doc, C.int(i))
	d.openedPopplerPages = append(d.openedPopplerPages, p)

	page = &Page{
		p:                p,
		openedPopplerAnnotMappings: []*C.struct__PopplerAnnotMapping{},
	}
	return page
}

func (d *Document) HasAttachments() bool {
	return toBool(C.poppler_document_has_attachments(d.doc))
}

func (d *Document) GetNAttachments() int {
	return int(C.poppler_document_get_n_attachments(d.doc))
}

func (d *Document) Close() {
	for i := 0; i < len(d.openedPopplerPages); i++ {
		C.g_object_unref(C.gpointer(d.openedPopplerPages[i]))
	}
	d.openedPopplerPages = []*C.struct__PopplerPage{}
	C.g_object_unref(C.gpointer(d.doc))
}

func (d *Document) NewAnnot(t AnnotType, r Rectangle, q []Quad) (Annot, error) {
	am := C.poppler_annot_mapping_new();

	annot := Annot {
		am: am,
	}

	pRect := rectangleToPopplerRectangle(r)

	pQuad := quadsToGArray(q)
	defer C.g_array_free(pQuad, 1)


	switch (t){
	case AnnotHighlight:
		am.annot = C.poppler_annot_text_markup_new_highlight(d.doc, &pRect, pQuad)
	case AnnotUnderline:
		am.annot = C.poppler_annot_text_markup_new_underline(d.doc, &pRect, pQuad)
	case AnnotSquiggly:
		am.annot = C.poppler_annot_text_markup_new_squiggly(d.doc, &pRect, pQuad)
	case AnnotStrikeOut:
		am.annot = C.poppler_annot_text_markup_new_strikeout(d.doc, &pRect, pQuad)
	default:
		C.poppler_annot_mapping_free(am)
		return annot, errors.New("invalid type for new annotation")
	}


	if am.annot == nil {
		C.poppler_annot_mapping_free(am)
		return annot, errors.New("failed to create annotation")
	}

	/* Can't get real annot mapping area as done in
	 * poppler_page_get_annot_mapping() since page is
	 * needed for page->page->getCropBox() and
	 * page->page->getRotate()
	 *
	 * as a placeholder we just use the annot rect
	 */
	annot.am.area = pRect

	return annot, nil
}

func (d *Document) Save(filename string) (saved bool, err error) {
	filename, err = filepath.Abs(filename)
	if err != nil {
		return false, err
	}

	var e *C.GError
	cFilename := (*C.gchar)(C.CString(filename))
	defer C.free(unsafe.Pointer(cFilename))

	cUri := C.g_filename_to_uri(cFilename, nil, nil)
	cBool := C.poppler_document_save (d.doc, cUri, &e);
	if e != nil {
		err = errors.New(C.GoString((*C.char)(e.message)))
		return false, err
	}

	if cBool == C.TRUE {
		return true, nil
	}

	return false, nil
}

/*
func (d *Document) GetAttachments() []Attachment {
	return
}
*/
