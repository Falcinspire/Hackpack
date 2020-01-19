package build

import (
	"strconv"
	"strings"

	"github.com/falcinspire/hackpackpdf/internal/lex"
	"github.com/jung-kurt/gofpdf"
)

const (
	TOP_MARGIN        = 14.0
	HORIZONTAL_MARGIN = 4.0
	COLUMN_PADDING    = 3.0
)

type indexElement struct {
	name  string
	index int
}

type Builder struct {
	headerText  string
	codeSize    float64
	headerSize  float64
	indexSize   float64
	columnWidth float64
	columnCount int
	pdf         *gofpdf.Fpdf
	columnIndex int
	index       []*indexElement
}

func New(headerText string, codeSize float64, headerSize float64, indexSize float64, columnCount int, paper string) *Builder {
	pdf := gofpdf.New("L", "mm", paper, "") //TODO validate paper
	pdf.SetCellMargin(0)
	builder := &Builder{headerText, codeSize, headerSize, indexSize, 0, columnCount, pdf, 0, make([]*indexElement, 0)}
	width, _ := pdf.GetPageSize()
	builder.columnWidth = (width - 2*HORIZONTAL_MARGIN) / float64(columnCount)
	pdf.SetMargins(TOP_MARGIN, HORIZONTAL_MARGIN, HORIZONTAL_MARGIN)
	pdf.SetAutoPageBreak(true, TOP_MARGIN/2)
	pdf.SetAcceptPageBreakFunc(builder.acceptPageBreak)
	pdf.SetHeaderFunc(builder.header)
	pdf.SetFooterFunc(builder.footer)
	builder.pdf.AddPage()
	builder.alignColumn(0)
	return builder
}

func (builder *Builder) AppendTitle(name string, lookup string) {
	builder.pdf.SetFont("Courier", "", builder.headerSize)
	_, unitSize := builder.pdf.GetFontSize()
	builder.pdf.SetTextColor(0, 0, 0)
	builder.pdf.Write(unitSize, name)
	builder.pdf.Write(unitSize, "\n")
	builder.pdf.SetDrawColor(100, 100, 100)
	builder.pdf.Line(builder.pdf.GetX()-0.5, builder.pdf.GetY()+1.0, builder.pdf.GetX()+builder.pdf.GetStringWidth(name)+1.0, builder.pdf.GetY()+1.0)
	builder.pdf.SetFont("Courier", "", builder.codeSize)
	_, unitSize = builder.pdf.GetFontSize()
	builder.pdf.Write(unitSize, "\n")
	builder.index = append(builder.index, &indexElement{lookup, builder.pdf.PageNo()})
}

func (builder *Builder) AppendCode(line []*lex.ColoredElement) {
	builder.pdf.SetFont("Courier", "", builder.codeSize)
	_, unitSize := builder.pdf.GetFontSize()
	for _, element := range line {
		// width, _ := builder.pdf.GetPageSize()
		// _, _, rightMargin, _ := builder.pdf.GetMargins()
		// if builder.pdf.GetX()+builder.pdf.GetStringWidth(element.Content) > (width - rightMargin) {
		// 	builder.pdf.Write(unitSize, "\n")
		// }
		builder.pdf.SetFont("Courier", joinStyles(element.Underline, element.Italic, element.Bold), builder.codeSize)
		builder.pdf.SetTextColor(int(element.Red), int(element.Green), int(element.Blue))
		builder.pdf.Write(unitSize, element.Content)
	}
	builder.pdf.Write(unitSize, "\n")
}

func (builder *Builder) AppendLine() {
	builder.pdf.SetFont("Courier", "", builder.codeSize)
	_, unitSize := builder.pdf.GetFontSize()
	builder.pdf.Write(unitSize, "\n")
}

func (builder *Builder) AppendIndex() {
	builder.pdf.SetFooterFunc(func() {}) //TODO this might cause bugs if unexpected
	builder.pdf.AddPage()
	builder.pdf.SetTextColor(0, 0, 0)
	builder.pdf.SetFont("Courier", "", builder.indexSize)
	_, unitSize := builder.pdf.GetFontSize()
	builder.alignColumn(0)
	for _, element := range builder.index {
		name := element.name
		page := strconv.Itoa(element.index)
		indexLine := name + (builder.makeDots(builder.fitInColumn(name, page))) + page
		// fmt.Println("==============")
		// fmt.Println(builder.columnWidth - 2*COLUMN_PADDING)
		// fmt.Println(builder.pdf.GetStringWidth(indexLine))
		builder.pdf.Write(unitSize, indexLine)
		builder.pdf.Write(unitSize, "\n")
	}
}

func (builder *Builder) alignColumn(c int) {
	builder.columnIndex = c
	leftMargin := HORIZONTAL_MARGIN + builder.columnWidth*float64(c)
	width, _ := builder.pdf.GetPageSize()
	rightMargin := width - (leftMargin + builder.columnWidth)
	builder.pdf.SetLeftMargin(leftMargin + COLUMN_PADDING)
	builder.pdf.SetRightMargin(rightMargin + COLUMN_PADDING)
	builder.pdf.SetY(TOP_MARGIN)
}

func (builder *Builder) drawColumn(c int) {
	leftMargin := HORIZONTAL_MARGIN + builder.columnWidth*float64(c)
	if c > 0 {
		builder.pdf.SetDrawColor(100, 100, 100)
		_, height := builder.pdf.GetPageSize()
		builder.pdf.Line(leftMargin, 12, leftMargin, height*0.98)
	}
}

func (builder *Builder) header() {
	pdf := builder.pdf
	pdf.SetLeftMargin(HORIZONTAL_MARGIN)
	pdf.SetRightMargin(HORIZONTAL_MARGIN)
	pdf.SetTopMargin(0)
	pdf.SetFont("Arial", "B", 12)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetXY(10, 7)
	_, unitSize := pdf.GetFontSize()
	pdf.Write(unitSize, builder.headerText)
	pdf.Write(unitSize, "\n")
	pdf.SetTopMargin(TOP_MARGIN)
	pdf.SetDrawColor(100, 100, 100)
	width, _ := pdf.GetPageSize()
	pdf.Line(width*0.02, pdf.GetY()+0.5, width*0.96, pdf.GetY()+0.5)

	builder.alignColumn(0) // TODO this needs to be here until margins are fixed
	builder.drawColumn(0)
}

func (builder *Builder) footer() {
	pdf := builder.pdf
	pdf.SetLeftMargin(HORIZONTAL_MARGIN)
	pdf.SetRightMargin(HORIZONTAL_MARGIN)
	pdf.SetY(-10)
	pdf.SetX(-10)
	pdf.SetFont("Arial", "", 7)
	_, unitSize := pdf.GetFontSize()
	pdf.SetTextColor(0, 0, 0)
	pdf.Write(unitSize, strconv.Itoa(pdf.PageNo()))

	builder.alignColumn(0) // TODO this needs to be here until margins are fixed
}

func (builder *Builder) acceptPageBreak() bool {
	if builder.columnIndex+1 == builder.columnCount {
		// builder.alignColumn(0)
		return true
	}
	builder.columnIndex++
	builder.alignColumn(builder.columnIndex)
	builder.drawColumn(builder.columnIndex)
	return false
}

func joinStyles(underline, italic, bold bool) string {
	res := ""
	if underline {
		res += "U"
	}
	if italic {
		res += "I"
	}
	if bold {
		res += "B"
	}
	return res
}

func (builder *Builder) fitInColumn(a, b string) int {
	totalSpace := builder.columnWidth - 2*COLUMN_PADDING
	aSize := builder.pdf.GetStringWidth(a)
	bSize := builder.pdf.GetStringWidth(b)
	space := totalSpace - (aSize + bSize)
	dotSize := builder.pdf.GetStringWidth(".")
	return int(space / dotSize)
}

func (builder *Builder) makeDots(count int) string {
	var sb strings.Builder
	for i := 0; i < count; i++ {
		sb.WriteRune('.')
	}
	return sb.String()
}

func (builder *Builder) SaveAndClose(path string) {
	err := builder.pdf.OutputFileAndClose(path)
	if err != nil {
		panic(err)
	}
}
