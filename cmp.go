package main

import (
	"crypto/md5"
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

func parseCmpBody(header CMPHeader, input []byte) ([]byte, error) {
	output := make([]byte, header.uncompressedSize*3)
	inputPos := 0
	outputPos := 0
	var relativeMode bool
	firstBytePattern := bytesToBinary([]byte{input[0]})
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
			// Copy count bytes in output buffer from outputPos - pos to  outputPos
			pattern := bytesToBinary([]byte{input[inputPos], input[inputPos+1]})
			count := ((current & 0x70) >> 4) + 3
			tpos0 := current & 0x0f                  // Lower nibble of the current byte
			tpos1 := int(tpos0) << 8                 // Moves it a byte to the left so now, we have 12 bits with the above 4 as the most significant bits
			inputPos += 1                            // Get the next byte
			pos := int(tpos1 + int(input[inputPos])) // Adds the next byte. So there's 12 bits now
			source := outputPos - pos
			slog.Debug("C2:", "id", commandCount, "indone", inDone, "opdone", opDone, "count", count, "pattern", pattern, "from", source)
			for range count {
				output[outputPos] = output[source]
				outputPos += 1
				source += 1
			}
			inputPos += 1 // Go to the next command
		} else if current == 0xfe {
			pattern := bytesToBinary([]byte{input[inputPos], input[inputPos+1], input[inputPos+2], input[inputPos+3]})
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
		} else if current == 0xff {
			pattern := bytesToBinary([]byte{
				input[inputPos],
				input[inputPos+1],
				input[inputPos+2],
				input[inputPos+3],
				input[inputPos+4],
			})
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
		} else if (current & 0xc0) == 0x80 {
			pattern := bytesToBinary([]byte{current})
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
		} else if (current & 0xc0) == 0xc0 {
			pattern := bytesToBinary([]byte{
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
			for range count {
				output[outputPos] = output[target]
				outputPos += 1
				target += 1
			}
			inputPos += 3 //Go to next command
		} else {
			slog.Error("Corrupt file. This shouldn't happen")
			break
		}
		if commandCount > 1235 && commandCount < 1245 {
			pngFname := fmt.Sprintf("/tmp/frames/frame-%05d.png", commandCount)
			writeCMPToPNG(output, pngFname, 320, 200)
		}

	}

	return output, nil
}

func decodeCmp(filename string, fileContents []byte) []byte {
	slog.Info("Decompressing CMP file", "name", filename)
	header, checksum := parseCmpHeader(fileContents)
	slog.Debug("Header obtained", "header", header.String(), "checksum", checksum)
	decompressedData, err := parseCmpBody(header, fileContents[10:]) // TODO : Probably need to check for compression type etc. here
	if err != nil {
		slog.Error("Aborted ", "error", err)
		return []byte{}
	} else {
		return decompressedData
	}
}
