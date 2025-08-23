package formats

import (
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log/slog"
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

type Room struct {
	N byte
	S byte
	E byte
	W byte
}

type Plan struct {
	Cells     [1024]Room
	image     *image.RGBA
	color     color.RGBA
	bgColor   color.RGBA
	size      int
	cellwidth uint8
	border    int
	cellSize  int
}

func NewPlan(data []byte, size int, bgColor color.RGBA, color color.RGBA, border int) Plan {
	var rooms [1024]Room

	for i := 0; i < len(data); i += 4 {
		n := data[i]
		e := data[i+1]
		s := data[i+2]
		w := data[i+3]
		r := Room{N: n, S: s, W: w, E: e}
		rooms[i/4] = r
	}

	img := image.NewRGBA(image.Rect(0, 0, size, size))
	bg := image.Rect(0, 0, size, size)
	draw.Draw(img, bg, &image.Uniform{bgColor}, image.Point{}, draw.Over)

	availableWidth := size - (border * 2)
	cellSize := availableWidth / 32

	return Plan{
		Cells:    rooms,
		image:    img,
		color:    color,
		size:     size,
		bgColor:  bgColor,
		border:   border,
		cellSize: cellSize,
	}
}

func (p *Plan) drawRoom(x int, y int, label string, room Room) {

	dim := color.RGBA{R: p.color.R / 4, G: p.color.G / 4, B: p.color.B / 4}
	d := &font.Drawer{
		Dst:  p.image,               // Your *image.RGBA
		Src:  image.NewUniform(dim), // Text color
		Face: basicfont.Face7x13,    // Built-in bitmap font
	}
	px := p.border*2 + (x * p.cellSize)
	py := p.border*2 + (y * p.cellSize)
	c := p.cellSize / 2

	var r image.Rectangle

	if room.N != 0 && room.S != 0 && room.E != 0 && room.W != 0 {
		r = image.Rect(px-c, py-c, px+c, py+c)
		draw.Draw(p.image, r, &image.Uniform{dim}, image.Point{}, draw.Over)
		return
	}
	// Handle North
	r = image.Rect(px-c, py-c, px+c, py-c+1)
	if room.N != 0 {
		draw.Draw(p.image, r, &image.Uniform{p.color}, image.Point{}, draw.Over)
	} // else {
	// 	draw.Draw(p.image, r, &image.Uniform{dim}, image.Point{}, draw.Over)
	// }

	// Handle south
	r = image.Rect(px-c, py+c-1, px+c, py+c)
	if room.S != 0 {
		draw.Draw(p.image, r, &image.Uniform{p.color}, image.Point{}, draw.Over)
	} // else {
	// 	draw.Draw(p.image, r, &image.Uniform{dim}, image.Point{}, draw.Over)
	// }

	// handle East
	r = image.Rect(px+c-1, py-c, px+c, py+c)
	if room.E != 0 {
		draw.Draw(p.image, r, &image.Uniform{p.color}, image.Point{}, draw.Over)
	} // else {
	// 	draw.Draw(p.image, r, &image.Uniform{dim}, image.Point{}, draw.Over)

	// }

	// Handle West
	r = image.Rect(px-c, py-c, px-c+1, py+c)
	if room.W != 0 {
		draw.Draw(p.image, r, &image.Uniform{p.color}, image.Point{}, draw.Over)
	} // else {
	// 	draw.Draw(p.image, r, &image.Uniform{dim}, image.Point{}, draw.Over)

	// }

	// room := image.Rect(px-c, py-c, px+c, py+c)
	// draw.Draw(p.image, room, &image.Uniform{dim}, image.Point{}, draw.Over)
	// inside := image.Rect(px-(c-1), py-(c-1), px+(c-1), py+(c-1))
	// draw.Draw(p.image, inside, &image.Uniform{p.bgColor}, image.Point{}, draw.Over)

	if label != "" {
		d.Dot = fixed.Point26_6{X: fixed.I(px - c + 5), Y: fixed.I(py)}
		d.DrawString(label)
	}

}

func (p *Plan) DrawPlan(fname string) {
	slog.Debug("Parameters", "size", p.size, "border", p.border, "cellSize", p.cellSize)

	for i := range len(p.Cells) {
		ix := i % 32
		iy := i / 32
		slog.Debug("Room Info", "ix", ix, "iy", iy, "room", p.Cells[i])
		p.drawRoom(ix, iy, fmt.Sprintf("%dx%d", ix, iy), p.Cells[i])
		// p.drawRoom(ix, iy, "", p.Cells[i])
	}
	f, _ := os.Create(fname)
	defer f.Close()
	png.Encode(f, p.image)

}

func (p *Plan) Validate() {
	for i := range len(p.Cells) {
		ix := i % 32
		iy := i / 32
		slog.Debug("Room Info", "ix", ix, "iy", iy, "room", p.Cells[i])
		// room := pm.Cells[i]

		p.drawRoom(ix, iy, fmt.Sprintf("%dx%d", ix, iy), p.Cells[i])
		// p.drawRoom(ix, iy, "", p.Cells[i])
	}
}

func DecodeMAZ(mazFileName string, input []byte) Plan {
	inputPos := 0
	width := binary.LittleEndian.Uint16(input[inputPos : inputPos+2])
	inputPos += 2
	height := binary.LittleEndian.Uint16(input[inputPos : inputPos+2])
	inputPos += 2
	tileSize := binary.LittleEndian.Uint16(input[inputPos : inputPos+2])
	inputPos += 2

	slog.Debug("Header read", "width", width, "height", height, "tileSize", tileSize)

	bgColor := color.RGBA{R: 0x14, G: 0x13, B: 0x3f, A: 0xFF}
	color := color.RGBA{R: 0xe6, G: 0xf8, B: 0xff, A: 0xFF}

	return NewPlan(input[inputPos:], 2000, bgColor, color, 40)
}
