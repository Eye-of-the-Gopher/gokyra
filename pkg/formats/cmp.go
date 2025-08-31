package formats

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"log/slog"

	"github.com/nibrahim/eye-of-the-gopher/internal/utils"
)

type CMPHeader struct {
	fileSize         uint16
	compressionType  uint16
	uncompressedSize uint32
	paletteSize      uint16
}

func parseCmpHeader(data []byte) (CMPHeader, string) {
	fileSize := binary.LittleEndian.Uint16(data[0:2])
	compressionType := binary.LittleEndian.Uint16(data[2:4])
	uncompressedSize := binary.LittleEndian.Uint32(data[4:8])
	paletteSize := binary.LittleEndian.Uint16(data[8:10])
	checksum := fmt.Sprintf("%x", md5.Sum(data))
	return CMPHeader{
		fileSize:         fileSize,
		compressionType:  compressionType,
		uncompressedSize: uncompressedSize,
		paletteSize:      paletteSize,
	}, checksum
}

func (h CMPHeader) String() string {
	return fmt.Sprintf("File size : %8d | Compression type : %d | Uncompressed size : %d | Palette size : %d",
		h.fileSize, h.compressionType, h.uncompressedSize, h.paletteSize)
}

func parseCmpBody(header CMPHeader, input []byte, palette color.Palette) ([]byte, error) {
	output := make([]byte, header.uncompressedSize)
	inputPos := 0
	outputPos := 0
	var relativeMode bool
	firstBytePattern := utils.BytesToBinary([]byte{input[0]})

	commandCount := 0
	slog.Debug("First byte", "pattern", firstBytePattern)
	if input[0] == 0x0 {
		slog.Debug("Relative mode")
		relativeMode = true
		inputPos += 1
	} else {
		slog.Debug("Absolute mode")
		relativeMode = false
	}

	for inputPos < len(input) {
		commandCount += 1
		inDone := fmt.Sprintf("%04d/%04d (%.2f%%)", inputPos, len(input), (float64(inputPos)/float64(len(input)))*100)
		opDone := fmt.Sprintf("%06d/%06d (%.2f%%)", outputPos, len(output), (float64(outputPos)/float64(len(output)))*100)
		current := input[inputPos]
		if current == 0x80 {
			slog.Debug("End of stream")
			break
		} else if (current & 0x80) == 0 {
			begin := inputPos
			// Copy count bytes in output buffer from outputPos - pos to  outputPos
			pattern := utils.BytesToBinary([]byte{input[inputPos], input[inputPos+1]})
			count := ((current & 0x70) >> 4) + 3
			tpos0 := current & 0x0f                  // Lower nibble of the current byte
			tpos1 := int(tpos0) << 8                 // Moves it a byte to the left so now, we have 12 bits with the above 4 as the most significant bits
			inputPos += 1                            // Get the next byte
			pos := int(tpos1 + int(input[inputPos])) // Adds the next byte. So there's 12 bits now
			source := outputPos - pos
			slog.Debug("C2:", "id", commandCount, "indone", inDone, "opdone", opDone, "count", count, "pattern", pattern, "from", source)
			if source < 0 || source >= len(output) {
				slog.Error("C2 invalid source", "source", source, "outputPos", outputPos, "pos", pos)
				return nil, fmt.Errorf("invalid source position")
			}
			for range count {
				output[outputPos] = output[source]
				outputPos += 1
				source += 1
			}
			inputPos += 1 // Go to the next command
			if inputPos != begin+2 {
				panic("C2 Increment wrong")
			}
		} else if current == 0xfe {
			begin := inputPos
			pattern := utils.BytesToBinary([]byte{input[inputPos], input[inputPos+1], input[inputPos+2], input[inputPos+3]})
			inputPos += 1                                                     // Leave command and Go to count byte 1
			count := binary.LittleEndian.Uint16(input[inputPos : inputPos+2]) // Read Count byte 1 and 2
			inputPos += 2                                                     // Go to value
			value := input[inputPos]
			slog.Debug("C4:", "id", commandCount, "indone", inDone, "opdone", opDone, "count", count, "pattern", pattern, "value", value)
			for range count {
				output[outputPos] = value
				outputPos += 1
			}
			inputPos += 1 // Go to next command
			if inputPos != begin+4 {
				panic("C4 Increment wrong")
			}

		} else if current == 0xff {
			begin := inputPos
			pattern := utils.BytesToBinary(input[inputPos : inputPos+5])
			count := int(binary.LittleEndian.Uint16(input[inputPos+1 : inputPos+3]))
			pos := int(binary.LittleEndian.Uint16(input[inputPos+3 : inputPos+5]))
			var target int
			if relativeMode {
				target = outputPos - pos
			} else {
				target = pos
			}
			slog.Debug("C5:", "id", commandCount, "indone", inDone, "opdone", opDone, "count", count, "pattern", pattern, "to", target)
			for range count {
				output[outputPos] = output[target]
				outputPos += 1
				target += 1
			}
			inputPos += 5

			if inputPos != begin+5 {
				panic("C5 Increment wrong")
			}
		} else if (current & 0xc0) == 0x80 {
			begin := inputPos
			pattern := utils.BytesToBinary([]byte{current})
			count := current & 0x3f
			slog.Debug("C1:", "id", commandCount, "indone", inDone, "opdone", opDone, "count", count, "pattern", pattern)
			if count == 0 {
				slog.Debug("Count 0. End of stream")
			}
			inputPos += 1 // Go to data
			for range count {
				output[outputPos] = input[inputPos]
				inputPos += 1
				outputPos += 1
			}
			if inputPos != begin+int(count)+1 {
				panic("C1 Increment wrong")
			}
		} else if (current & 0xc0) == 0xc0 {
			begin := inputPos
			pattern := utils.BytesToBinary([]byte{
				input[inputPos],
				input[inputPos+1],
				input[inputPos+2],
			})
			count := (current & 0x3f) + 3
			var target int
			pos := int(binary.LittleEndian.Uint16(input[inputPos+1 : inputPos+3]))
			if relativeMode {
				target = outputPos - pos
			} else {
				target = pos
			}
			slog.Debug("C3:", "id", commandCount, "indone", inDone, "opdone", opDone, "count", count, "pattern", pattern, "to", target)
			if target < 0 || target >= len(output) {
				slog.Error("C3 invalid target", "target", target, "outputPos", outputPos, "pos", pos)
				return nil, fmt.Errorf("invalid target position")
			}
			for range count {
				output[outputPos] = output[target]
				outputPos += 1
				target += 1
			}
			inputPos += 3 //Go to next command
			if inputPos != begin+3 {
				panic("C3 Increment wrong")
			}
		} else {
			slog.Error("Corrupt file. This shouldn't happen")
			break
		}

	}

	return output, nil
}

// Converts a CMP data stream in data to an image.Image
func CMPToImage(data []byte, palette color.Palette, width int, height int) image.Image {
	slog.Debug("Converting CMP PNG", "length", len(data))
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill the image with the raw data
	for y := range height {
		for x := range width {
			index := y*width + x
			if index < len(data) {
				img.Set(x, y, palette[data[index]])
			}
		}
	}
	return img
}

func DecodeCmp(filename string, fileContents []byte, palette color.Palette) []byte {
	slog.Info("Decompressing CMP file", "name", filename)
	header, checksum := parseCmpHeader(fileContents)
	slog.Debug("Header obtained", "header", header.String(), "checksum", checksum)
	decompressedData, err := parseCmpBody(header, fileContents[10:], palette) // TODO : Probably need to check for compression type etc. here
	if err != nil {
		slog.Error("Aborted ", "error", err)
		return []byte{}
	} else {
		return decompressedData
	}
}
