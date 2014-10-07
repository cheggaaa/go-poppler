package poppler

import ()

type TextEl struct {
	Text  string
	Attrs *TextAttributes
	Rect  Rectangle
}

type TextAttributes struct {
	FontName             string
	FontSize             float64
	IsUnderlined         bool
	Color                Color
	StartIndex, EndIndex int
}

type Color struct {
	R, G, B int
}
