package main

import (
	"encoding/binary"
	"fmt"
	"log/slog"
)

type CMPHeader struct {
	fileSize         uint16
	compressionType  uint16
	uncompressedSize uint32
	paletteSize      uint16
}

func parseCmpHeader(data []byte) CMPHeader {
	fileSize := binary.LittleEndian.Uint16(data[0:2])
	compressionType := binary.LittleEndian.Uint16(data[2:4])
	uncompressedSize := binary.LittleEndian.Uint32(data[4:8])
	paletteSize := binary.LittleEndian.Uint16(data[8:10])
	return CMPHeader{
		fileSize:         fileSize,
		compressionType:  compressionType,
		uncompressedSize: uncompressedSize,
		paletteSize:      paletteSize,
	}
}

func (h CMPHeader) String() string {
	return fmt.Sprintf("File size : %8d | Compression type : %d | Uncompressed size : %d | Palette size : %d",
		h.fileSize, h.compressionType, h.uncompressedSize, h.paletteSize)
}

func parseCmpBody(header CMPHeader, input []byte) ([]byte, error) {
	output := make([]byte, header.uncompressedSize)
	inputPos := 0
	outputPos := 0
	var relativeMode bool
	if input[0] == 0x0 {
		slog.Debug("Relative mode")
		relativeMode = true
		inputPos += 1
	} else {
		slog.Debug("Absolute mode")
		relativeMode = false
	}

	for inputPos < len(input) {
		done := fmt.Sprintf("%d/%d (%.2f%%)", inputPos, len(input), (float64(inputPos)/float64(len(input)))*100)
		current := input[inputPos]
		if current == 0x80 {
			slog.Debug("End of stream")
			break
		} else if (current & 0x80) == 0 {
			// Copy count bytes in output buffer from outputPos - pos to  outputPos
			count := ((current & 0x70) >> 4) + 3
			tpos0 := current & 0x0f                  // Lower nibble of the current byte
			tpos1 := int(tpos0) << 8                 // Moves it a byte to the left so now, we have 12 bits with the above 4 as the most significant bits
			inputPos += 1                            // Get the next byte
			pos := int(tpos1 + int(input[inputPos])) // Adds the next byte. So there's 12 bits now
			source := outputPos - pos
			slog.Debug("C2:", "count", count, "from", source, "done", done)
			for range count {
				output[outputPos] = output[source]
				outputPos += 1
				source += 1
			}
			inputPos += 1 // Go to the next byte
		} else if current == 0xfe {
			inputPos += 1 // Go to count
			count := binary.LittleEndian.Uint16(input[inputPos : inputPos+2])
			inputPos += 2 // Go to value
			value := input[inputPos]
			inputPos += 1 // Go to next command
			pattern := fmt.Sprintf("%08b", value)
			slog.Debug("C4:", "count", count, "value", value, "in binary", pattern, "done", done)
			for range count {
				output[outputPos] = value
				outputPos += 1
			}
		} else if current == 0xff {
			count := int(binary.LittleEndian.Uint16(input[inputPos+1 : inputPos+3]))
			pos := int(binary.LittleEndian.Uint16(input[inputPos+3 : inputPos+5]))
			var target int
			if relativeMode {
				target = outputPos - pos
			} else {
				target = pos
			}
			slog.Debug("C5:", "count", count, "to", target, "done", done)
			for range count {
				output[outputPos] = output[target]
				outputPos += 1
				target += 1
			}
			inputPos += 5
		} else if (current & 0xc0) == 0x80 {
			pattern := fmt.Sprintf("%08b", current)
			count := current & 0x3f
			slog.Debug("C1:", "count", count, "pattern", pattern, "done", done)
			inputPos += 1 // Go to next command
			for range count {
				output[outputPos] = input[inputPos]
				inputPos += 1
				outputPos += 1
			}
		} else if (current & 0xc0) == 0xc0 {
			count := (current & 0x3f) + 3
			var target int
			pos := int(binary.LittleEndian.Uint16(input[inputPos+1 : inputPos+3]))
			if relativeMode {
				target = outputPos - pos
			} else {
				target = pos
			}
			slog.Debug("C3:", "count", count, "to", target, "done", done)
			for range count {
				output[outputPos] = output[target]
				outputPos += 1
				target += 1
			}
			inputPos += 3
		} else {
			slog.Error("Corrupt file. This shouldn't happen")
			break
		}
	}

	return output, nil
}

func decodeCmp(filename string, fileContents []byte) []byte {
	slog.Info("Decompressing CMP file", "name", filename)
	header := parseCmpHeader(fileContents)
	slog.Debug("Header obtained", "header", header.String())
	decompressedData, err := parseCmpBody(header, fileContents[10:]) // TODO : Probably need to check for compression type etc. here
	if err != nil {
		panic("Boom!")
		return nil
	} else {
		return decompressedData
	}
}
