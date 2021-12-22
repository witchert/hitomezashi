package hitomezashi

import (
	"image/color"
	"math"

	"github.com/fogleman/gg"
)

type ColorChoice int

const (
	COLOR_A ColorChoice = 0
	COLOR_B ColorChoice = 1
)

var borderWidth = 24

type Cell struct {
	colorChoice ColorChoice
	rightBound  bool
	topBound    bool
}

type Hitomezashi struct {
	Canvas     *gg.Context
	width      int
	height     int
	cellSize   int
	horizontal []bool
	vertical   []bool
	cells      []Cell
	colorA     color.Color
	colorB     color.Color
}

func New(horizontal []bool, vertical []bool, scale int) Hitomezashi {
	h := Hitomezashi{
		horizontal: horizontal,
		vertical:   vertical,
		width:      len(horizontal),
		height:     len(vertical),
		cellSize:   scale,
		colorA:     color.Black,
		colorB:     color.White,
		cells:      make([]Cell, len(horizontal)*len(vertical)),
	}

	pixelWidth := h.width*h.cellSize + borderWidth
	pixelHeight := h.height*h.cellSize + borderWidth
	h.Canvas = gg.NewContext(pixelWidth, pixelHeight)

	h.setupCells()

	return h
}

func (h *Hitomezashi) SetColors(a color.Color, b color.Color) {
	h.colorA = a
	h.colorB = b
}

func (h *Hitomezashi) Draw() {
	h.fillCells()
	h.drawLines()
	h.drawCanvasBorder()
}

func (h *Hitomezashi) setupCells() {
	for position := range h.cells {
		cell := Cell{colorChoice: COLOR_A, topBound: false, rightBound: false}

		cell.setCellBounds(position, h)
		cell.setCellColor(position, h)

		h.cells[position] = cell
	}
}

func (h *Hitomezashi) positionToCoords(position int) (int, int) {
	col := position % h.width
	row := position / h.width
	x := col*h.cellSize + (borderWidth / 2)
	y := row*h.cellSize + (borderWidth / 2)

	return x, y
}

func (h *Hitomezashi) fillCells() {
	for index, cell := range h.cells {
		x, y := h.positionToCoords(index)

		h.Canvas.SetColor(h.getColorFromChoice(cell.colorChoice))
		h.Canvas.DrawRectangle(float64(x), float64(y), float64(h.cellSize), float64(h.cellSize))
		h.Canvas.Fill()
	}
}

func (h *Hitomezashi) drawLines() {
	// Set line width based on the scale
	h.Canvas.SetLineWidth(math.RoundToEven(float64(h.cellSize / 6)))

	for index, cell := range h.cells {
		if cell.topBound || cell.rightBound {
			x, y := h.positionToCoords(index)

			h.Canvas.SetRGB(0, 0, 0)

			// Horizontal line
			if cell.topBound && index >= h.width {
				h.Canvas.DrawLine(float64(x), float64(y), float64(x+h.cellSize), float64(y))
				h.Canvas.Stroke()
			}

			// Vertical line
			if cell.rightBound && index%h.width < h.width-1 {
				h.Canvas.DrawLine(float64(x+h.cellSize), float64(y), float64(x+h.cellSize), float64(y+h.cellSize))
				h.Canvas.Stroke()
			}
		}
	}
}

func (h *Hitomezashi) drawCanvasBorder() {
	h.Canvas.SetLineWidth(float64(borderWidth))
	h.Canvas.SetRGB(1, 1, 1)
	h.Canvas.DrawRectangle(0, 0, float64(h.Canvas.Width()), float64(h.Canvas.Height()))
	h.Canvas.Stroke()
}

func (h *Hitomezashi) getColorFromChoice(choice ColorChoice) color.Color {
	if choice == COLOR_A {
		return h.colorA
	}
	return h.colorB
}

func (cell *Cell) setCellBounds(position int, h *Hitomezashi) {
	row := position / h.width
	col := position % h.width
	if (!h.vertical[row] && col%2 != 0) || (h.vertical[row] && col%2 == 0) {
		cell.topBound = true
	}

	if (!h.horizontal[col] && row%2 != 0) || (h.horizontal[col] && row%2 == 0) {
		cell.rightBound = true
	}
}

func (cell *Cell) setCellColor(position int, h *Hitomezashi) {
	if position > 0 {
		// Starting cell in a row, use color from row above
		if position%h.width == 0 {
			if cell.topBound {
				// the cell is top bound, use opposite color to the cell above
				cell.colorChoice = ^h.cells[position-h.width].colorChoice
			} else {
				// the cell is not top bound, use same color as the cell above
				cell.colorChoice = h.cells[position-h.width].colorChoice
			}
		} else if h.cells[position-1].rightBound {
			// the cell is not the starting cell and the previous cell is right bound
			// use opposote color to the previous cell
			cell.colorChoice = ^h.cells[position-1].colorChoice
		} else {
			// not bound above or by previous cell, use same color as the previous cell
			cell.colorChoice = h.cells[position-1].colorChoice
		}
	}
}
