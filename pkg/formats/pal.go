package formats

import (
	"image/color"
	"log/slog"
)

func DecodePalette(data []byte) color.Palette {
	ret := make(color.Palette, 256)
	idx := 0
	for i := 0; i < len(data); i += 3 {
		r := data[i] & 0x3F
		g := data[i+1] & 0x3F
		b := data[i+2] & 0x3F
		brightness := 3.0 // Scale up and apply brightness
		brightR := int(float64(r) * brightness)
		brightG := int(float64(g) * brightness)
		brightB := int(float64(b) * brightness)

		// Clamp to 255 max
		if brightR > 255 {
			brightR = 255
		}
		if brightG > 255 {
			brightG = 255
		}
		if brightB > 255 {
			brightB = 255
		}

		ret[idx] = color.RGBA{uint8(brightR), uint8(brightG), uint8(brightB), 255}
		idx += 1
	}
	slog.Debug("Palette created")
	return ret

}
