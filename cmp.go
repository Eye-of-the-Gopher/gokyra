package main

import (
	"encoding/binary"
	"fmt"
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

func parseCmpBody(header CMPHeader, data []byte) ([]byte, error) {
	return []byte{}, nil
}

func decodeCmp(filename string, fileContents []byte) {
	slog.Info("Decompressing CMP file", "name", filename)
	header := parseCmpHeader(fileContents)
	slog.Debug("Header obtained", "header", header.String())
	decompressedData, err := parseCmpBody(header, fileContents[10:]) // TODO : Probably need to check for compression type etc. here
	if err != nil {
		fmt.Println("Boom!")
	} else {
		fmt.Printf("%s", string(decompressedData))
	}

}
