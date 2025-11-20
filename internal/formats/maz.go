package formats

import (
	"encoding/binary"
	"fmt"
	"image/color"
	"os"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/basicfont"
)

type Room struct {
	N byte
	S byte
	E byte
	W byte
}

type Plan struct {
	Cells [1024]Room
	size  int
	// cellwidth uint8
	border   int
	cellSize int
}

func createInfoBox(dc *gg.Context, text string, date string) {
	// Load TTF file
	x := 50.0
	y := float64(dc.Height()) - 150
	fontBytes, err := os.ReadFile("RomanUncialModern.ttf")
	if err != nil {
		panic(err)
	}

	font, err := truetype.Parse(fontBytes)
	if err != nil {
		panic(err)
	}

	face := truetype.NewFace(font, &truetype.Options{Size: 30})
	dc.SetFontFace(face)

	// Set color AFTER loading font
	dc.SetHexColor("#422000")
	dc.DrawString(text, x, y)
	dc.DrawLine(x, y+5, float64(dc.Height()/2), y+5)
	dc.Stroke()
	face = truetype.NewFace(font, &truetype.Options{Size: 15})
	dc.SetFontFace(face)
	dc.DrawString(date, x, y+30)
}

func drawRect(dc *gg.Context, x float64, y float64, w float64, h float64, hatch bool, border bool, fill bool) {
	// Fill background
	dc.SetHexColor("422000")
	// dc.SetHexColor("#c1ab6744")
	dc.SetLineWidth(1.0)
	dc.DrawRectangle(x, y, w, h)
	dc.Stroke()

	if fill {
		dc.SetHexColor("#422000aa")
		dc.DrawRectangle(x, y, w, h)
		dc.Fill()
	}

	if hatch {
		// Set clipping region to rectangle
		dc.DrawRectangle(x, y, w, h)
		dc.Clip()
		// Draw hatches (can go beyond bounds, will be clipped)
		dc.SetHexColor("#8B7355")
		dc.SetLineWidth(1)

		spacing := 10.0
		for i := -h; i < w+h; i += spacing {
			dc.DrawLine(x+i, y, x+i+h, y+h)
		}
		dc.Stroke()
	}

	if border {
		if !fill {
			dc.SetHexColor("#4A3F35") // Dark border
			dc.SetLineWidth(1)
			dc.DrawRectangle(x, y, w, h)
			dc.Stroke()
		} else {
			dc.SetHexColor("#000000") // Black border
			dc.SetLineWidth(1)
			dc.DrawRectangle(x, y, w, h)
			dc.Stroke()
		}

	}

	dc.ResetClip() // Important! Reset clipping for next operations
}

func NewPlan(data []byte, size int, border int) Plan {
	var rooms [1024]Room

	for i := 0; i < len(data); i += 4 {
		n := data[i]
		e := data[i+1]
		s := data[i+2]
		w := data[i+3]
		r := Room{N: n, S: s, W: w, E: e}
		rooms[i/4] = r
	}

	availableWidth := size - (border * 2)
	cellSize := availableWidth / 32

	return Plan{
		Cells:    rooms,
		size:     size,
		border:   border,
		cellSize: cellSize,
	}
}

func (p *Plan) drawRoom(dc *gg.Context, x int, y int, label string, room Room) {
	px := float64(p.border*2 + (x * p.cellSize))
	py := float64(p.border*2 + (y * p.cellSize))
	c := float64(p.cellSize)

	MazLogger.Debug("Room Info", "x", x, "y", y, "px", px, "py", py, "room", room)

	doorWidth := float64(p.cellSize / 4)
	wallThickness := float64(p.cellSize / 5)
	doorThickness := float64(wallThickness) * 1.5
	MazLogger.Debug("Component dimensions", "doorWidth", doorWidth, "wallThickness", wallThickness, "doorThickness", doorThickness)

	ns := room.N
	ss := room.S
	es := room.E
	ws := room.W
	if (ns == 1 || ns == 2) &&
		(ss == 1 || ss == 2) &&
		(es == 1 || es == 2) &&
		(ws == 1 || ws == 2) {
		return
	}
	dc.SetFontFace(basicfont.Face7x13) // 7x13 pixel font
	dc.SetColor(color.Black)
	dc.DrawString(label, px+5, py+c/2)

	if ns != 0 {
		MazLogger.Debug("North wall")
		drawRect(dc, px, py, c, wallThickness, true, true, false) // Draw North wall
		if 3 <= ns && ns <= 22 {
			MazLogger.Debug("North Door")
			drawRect(dc, px+c/2-doorWidth, py, doorWidth*2, doorThickness, false, true, true) // Draw door on this wall
		}
	} else {
		dc.SetHexColor("000000")
		dc.SetLineWidth(1)
		dc.DrawLine(px, py, px+float64(p.cellSize), py)
		dc.Stroke()
	}

	if ss != 0 {
		MazLogger.Debug("South wall")
		drawRect(dc, px, py+float64(p.cellSize)-wallThickness, float64(p.cellSize), wallThickness, true, true, false) // Draw a south wall
		if 3 <= ss && ss <= 22 {
			MazLogger.Debug("South door")
			drawRect(dc, px+c/2-doorWidth, py+float64(p.cellSize)-(doorThickness), doorWidth*2, doorThickness, false, true, true) // Draw door on this wall
		}
	} else {
		dc.SetHexColor("000000")
		dc.SetLineWidth(1)
		dc.DrawLine(px, py+float64(p.cellSize), px+float64(p.cellSize), py+float64(p.cellSize))
		dc.Stroke()
	}

	if es != 0 {
		drawRect(dc, px+float64(p.cellSize)-wallThickness, py, wallThickness, float64(p.cellSize), true, true, false) // Draw an east wall
		if 3 <= es && es <= 22 {
			drawRect(dc,
				px+float64(p.cellSize)-(doorThickness),
				py+float64(p.cellSize)/2-doorWidth,
				doorThickness,
				doorWidth*2,
				false, true, true) // Draw door on this wall
		}
	} else {
		dc.SetHexColor("000000")
		dc.SetLineWidth(1)
		dc.DrawLine(px+float64(p.cellSize), py, px+float64(p.cellSize), py+float64(p.cellSize))
		dc.Stroke()

	}

	if ws != 0 {
		drawRect(dc, px, py, wallThickness, float64(p.cellSize), true, true, false) // Draw a west wall
		if 3 <= ws && ws <= 22 {
			drawRect(dc,
				px,
				py+float64(p.cellSize)/2-doorWidth,
				doorThickness,
				doorWidth*2,
				false, true, true) // Draw door on this wall
		}
	} else {
		dc.SetHexColor("000000")
		dc.SetLineWidth(1)
		dc.DrawLine(px, py, px, py+float64(p.cellSize))
		dc.Stroke()
	}

}

func (p *Plan) DrawPlan(fname string) {
	MazLogger.Debug("Parameters", "size", p.size, "border", p.border, "cellSize", p.cellSize)

	img, err := gg.LoadPNG("ip2.png")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't load background")
	}
	dc := gg.NewContextForImage(img)
	createInfoBox(dc, "Level 1: Upper sewers", "1372 DR, Flamerule 23rd")

	for i := range len(p.Cells) {
		ix := i % 32
		iy := i / 32

		// if ix >= 12 && ix <= 15 {
		// 	if iy >= 12 && iy <= 15 {
		p.drawRoom(dc, ix, iy, fmt.Sprintf("%dx%d", ix, iy), p.Cells[i])
		// 	}
		// }

	}

	dc.SavePNG(fname)
}

func DecodeMAZ(mazFileName string, input []byte) Plan {
	inputPos := 0
	width := binary.LittleEndian.Uint16(input[inputPos : inputPos+2])
	inputPos += 2
	height := binary.LittleEndian.Uint16(input[inputPos : inputPos+2])
	inputPos += 2
	tileSize := binary.LittleEndian.Uint16(input[inputPos : inputPos+2])
	inputPos += 2

	MazLogger.Debug("Header read", "width", width, "height", height, "tileSize", tileSize)

	return NewPlan(input[inputPos:], 1800, 50)
}
