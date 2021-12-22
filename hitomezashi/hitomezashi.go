package hitomezashi

import (
	"image/color"
	"math"
	"math/rand"
	"strconv"

	"github.com/fogleman/gg"
	"github.com/muesli/gamut"
)

type Cell struct {
	color      color.Color
	rightBound bool
	topBound   bool
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
	boxColorA   color.Color
	boxColorB   color.Color
}

var colors = []color.Color{
	color.NRGBA64{57818, 35930, 12168, 65535},
	color.NRGBA64{19216, 48673, 16633, 65535},
	color.NRGBA64{57090, 34948, 47660, 65535},
	color.NRGBA64{36003, 45435, 13660, 65535},
	color.NRGBA64{20404, 39496, 20721, 65535},
	color.NRGBA64{52060, 15473, 58198, 65535},
	color.NRGBA64{46641, 36500, 57269, 65535},
	color.NRGBA64{57480, 15399, 30954, 65535},
	color.NRGBA64{42716, 42029, 24161, 65535},
	color.NRGBA64{39231, 22794, 37294, 65535},
	color.NRGBA64{55575, 38184, 27239, 65535},
	color.NRGBA64{16859, 42870, 53708, 65535},
	color.NRGBA64{41379, 27522, 9102, 65535},
	color.NRGBA64{49700, 24602, 29228, 65535},
	color.NRGBA64{55766, 29478, 56914, 65535},
	color.NRGBA64{42260, 24550, 16624, 65535},
	color.NRGBA64{28496, 23990, 58742, 65535},
	color.NRGBA64{47353, 17967, 10505, 65535},
	color.NRGBA64{58946, 25652, 10518, 65535},
	color.NRGBA64{20285, 45089, 37379, 65535},
	color.NRGBA64{59253, 28945, 25262, 65535},
	color.NRGBA64{55792, 14417, 47691, 65535},
	color.NRGBA64{57658, 13672, 10787, 65535},
	color.NRGBA64{27779, 29712, 10939, 65535},
	color.NRGBA64{27078, 27904, 54256, 65535},
	color.NRGBA64{33328, 15023, 59635, 65535},
	color.NRGBA64{40308, 19184, 46614, 65535},
	color.NRGBA64{27566, 39157, 59264, 65535},
	color.NRGBA64{55481, 14389, 20760, 65535},
	color.NRGBA64{50442, 41923, 13432, 65535},
	color.NRGBA64{23731, 27990, 45070, 65535},
	color.NRGBA64{52365, 18426, 38505, 65535},
}

func New(horizontal []string, vertical []string, scale int, hash string) Hitomezashi {
	h := Hitomezashi{across: horizontal, down: vertical, width: len(horizontal), height: len(vertical), scale: scale, borderWidth: 24}

	// Set up colors
	seed, _ := strconv.ParseInt(hash[:7], 16, 64)
	r := rand.New(rand.NewSource(seed))
	h.boxColorA = colors[r.Intn(len(colors)-1)+1]
	h.boxColorB = gamut.Complementary(h.boxColorA)

	h.cells = make([]Cell, h.width*h.height)

	h.setupCanvas()

	h.calculateCells()

	h.drawLines()

	h.drawBorder()

	return h
}

func (h *Hitomezashi) setupCanvas() {
	pixelWidth := h.width*h.scale + h.borderWidth
	pixelHeight := h.height*h.scale + h.borderWidth
	h.Canvas = gg.NewContext(pixelWidth, pixelHeight)
	h.Canvas.SetHexColor("#ffffff")
	h.Canvas.Clear()
}

func (h *Hitomezashi) calculateCells() {
	for row := 0; row < h.height; row++ {
		for col := 0; col < h.width; col++ {
			position := (row * h.width) + col
			cell := Cell{color: h.boxColorA, topBound: false, rightBound: false}

			cell.determineCellBounds(row, col, h)
			cell.determineCellColor(position, h)

			// Draw the boxes...
			x := col*h.scale + (h.borderWidth / 2)
			y := row*h.scale + (h.borderWidth / 2)

			h.Canvas.SetColor(cell.color)
			h.Canvas.DrawRectangle(float64(x), float64(y), float64(h.scale), float64(h.scale))
			h.Canvas.Fill()

			h.cells[position] = cell
		}
	}
}

func (h *Hitomezashi) drawLines() {
	// Set up line drawing
	h.Canvas.SetLineWidth(math.RoundToEven(float64(h.scale / 6)))

	// Draw lines
	for pos, cell := range h.cells {
		if cell.topBound || cell.rightBound {
			col := pos % h.width
			row := pos / h.width
			x := int(col)*h.scale + (h.borderWidth / 2)
			y := int(row)*h.scale + (h.borderWidth / 2)

			h.Canvas.SetHexColor("#000000")

			// Draw the horizontal lines...
			if cell.topBound && pos >= h.width {
				h.Canvas.DrawLine(float64(x), float64(y), float64(x+h.scale), float64(y))
				h.Canvas.Stroke()
			}

			// Draw the vertical lines...
			if cell.rightBound && pos%h.width < h.width-1 {
				h.Canvas.DrawLine(float64(x+h.scale), float64(y), float64(x+h.scale), float64(y+h.scale))
				h.Canvas.Stroke()
			}
		}
	}
}

func (h *Hitomezashi) drawBorder() {
	h.Canvas.SetLineWidth(float64(h.borderWidth))
	h.Canvas.SetHexColor("#ffffff")
	h.Canvas.DrawRectangle(0, 0, float64(h.Canvas.Width()), float64(h.Canvas.Height()))
	h.Canvas.Stroke()
}

func (h *Hitomezashi) oppositeColor(color color.Color) color.Color {
	if color == h.boxColorA {
		return h.boxColorB
	}
	return h.boxColorA
}

func (cell *Cell) determineCellBounds(row int, col int, h *Hitomezashi) {
	// Figure out if bounded on top
	if (h.down[row] == "0" && col%2 != 0) || (h.down[row] == "1" && col%2 == 0) {
		cell.topBound = true
	}

	// Figure out if bounded on right
	if (h.across[col] == "0" && row%2 != 0) || (h.across[col] == "1" && row%2 == 0) {
		cell.rightBound = true
	}
}

func (cell *Cell) determineCellColor(position int, h *Hitomezashi) {
	if position > 0 {
		// If we are on the cell starting a new row then look to the row above
		if position%h.width == 0 {
			// is our cell top bound, if so use opposite color
			if cell.topBound {
				cell.color = h.oppositeColor(h.cells[position-h.width].color)
				// our cell isnt top bound so match the cell aboves color
			} else {
				cell.color = h.cells[position-h.width].color
			}
			// otherwise look to the previous cell to see if its right bound
		} else if h.cells[position-1].rightBound {
			cell.color = h.oppositeColor(h.cells[position-1].color)
			// previous cell isnt bound so match its color
		} else {
			cell.color = h.cells[position-1].color
		}
	}
}
