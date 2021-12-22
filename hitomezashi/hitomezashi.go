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

type Cell struct {
	colorChoice ColorChoice
	rightBound  bool
	topBound    bool
}

type Hitomezashi struct {
	Canvas      *gg.Context
	width       int
	height      int
	scale       int
	borderWidth int
	across      []string
	down        []string
	cells       []Cell
	colorA      color.Color
	colorB      color.Color
}

func New(horizontal []string, vertical []string, scale int) Hitomezashi {
	h := Hitomezashi{
		across:      horizontal,
		down:        vertical,
		width:       len(horizontal),
		height:      len(vertical),
		scale:       scale,
		borderWidth: 24.0,
		colorA:      color.Black,
		colorB:      color.White,
	}

	pixelWidth := h.width*h.scale + h.borderWidth
	pixelHeight := h.height*h.scale + h.borderWidth
	h.Canvas = gg.NewContext(pixelWidth, pixelHeight)

	h.setupCells()

	return h
}

func (h *Hitomezashi) SetColors(a color.Color, b color.Color) {
	h.colorA = a
	h.colorB = b
}

func (h *Hitomezashi) Draw() {
	h.drawCanvas()

	h.fillCells()
	h.drawLines()

	h.drawCanvasBorder()
}

func (h *Hitomezashi) drawCanvas() {
	h.Canvas.SetColor(color.White)
	h.Canvas.Clear()
}

func (h *Hitomezashi) setupCells() {
	for row := 0; row < h.height; row++ {
		for col := 0; col < h.width; col++ {
			cell := Cell{colorChoice: COLOR_A, topBound: false, rightBound: false}

			cell.setCellBounds(row, col, h)
			cell.setCellColor(row, col, h)

			h.cells = append(h.cells, cell)
		}
	}
}

func (h *Hitomezashi) positionToCoords(position int) (int, int) {
	col := position % h.width
	row := position / h.width
	x := col*h.scale + (h.borderWidth / 2)
	y := row*h.scale + (h.borderWidth / 2)

	return x, y
}

func (h *Hitomezashi) fillCells() {
	for index, cell := range h.cells {
		x, y := h.positionToCoords(index)

		h.Canvas.SetColor(h.getColorFromChoice(cell.colorChoice))
		h.Canvas.DrawRectangle(float64(x), float64(y), float64(h.scale), float64(h.scale))
		h.Canvas.Fill()
	}
}

func (h *Hitomezashi) drawLines() {
	// Set line width based on the scale
	h.Canvas.SetLineWidth(math.RoundToEven(float64(h.scale / 6)))

	for index, cell := range h.cells {
		if cell.topBound || cell.rightBound {
			x, y := h.positionToCoords(index)

			h.Canvas.SetRGB(0, 0, 0)

			// Horizontal line
			if cell.topBound && index >= h.width {
				h.Canvas.DrawLine(float64(x), float64(y), float64(x+h.scale), float64(y))
				h.Canvas.Stroke()
			}

			// Vertical line
			if cell.rightBound && index%h.width < h.width-1 {
				h.Canvas.DrawLine(float64(x+h.scale), float64(y), float64(x+h.scale), float64(y+h.scale))
				h.Canvas.Stroke()
			}
		}
	}
}

func (h *Hitomezashi) drawCanvasBorder() {
	h.Canvas.SetLineWidth(float64(h.borderWidth))
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

func (cell *Cell) setCellBounds(row int, col int, h *Hitomezashi) {
	if (h.down[row] == "0" && col%2 != 0) || (h.down[row] == "1" && col%2 == 0) {
		cell.topBound = true
	}

	if (h.across[col] == "0" && row%2 != 0) || (h.across[col] == "1" && row%2 == 0) {
		cell.rightBound = true
	}
}

func (cell *Cell) setCellColor(row int, col int, h *Hitomezashi) {
	position := (row * h.width) + col
	if position > 0 {
		// Starting cell in a row, use color from row above
		if position%h.width == 0 {
			// is our cell top bound, if so use opposite color
			if cell.topBound {
				cell.colorChoice = ^h.cells[position-h.width].colorChoice
				// our cell isnt top bound so match the cell aboves color
			} else {
				cell.colorChoice = h.cells[position-h.width].colorChoice
			}
			// otherwise look to the previous cell to see if its right bound
		} else if h.cells[position-1].rightBound {
			cell.colorChoice = ^h.cells[position-1].colorChoice
			// previous cell isnt bound so match its color
		} else {
			cell.colorChoice = h.cells[position-1].colorChoice
		}
	}
}
