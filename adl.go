package main

import (
	"encoding/binary"
)

func readADLIndices(input []byte) ([120]byte, [150]uint16, [150]uint16) {
	primaryIndex := input[:120]
	var trackPointers [150]uint16
	var instrumentPointers [150]uint16
	idx := 0
	for i := 0; i < 300; i += 2 {
		trackPointers[idx] = binary.LittleEndian.Uint16(input[120+i : 120+i+2])
		idx += 1
	}
	idx = 0
	for i := 0; i < 300; i += 2 {
		instrumentPointers[idx] = binary.LittleEndian.Uint16(input[420+i : 420+i+2])
		idx += 1
	}

	return [120]byte(primaryIndex), [150]uint16(trackPointers), [150]uint16(instrumentPointers)
}

func decodeADL(filename string, input []byte) ([]byte, error) {
	ret := make([]byte, 50)
	primaryIndex, trackPointer, instrumentPointers := readADLIndices(input)
	dumpList(primaryIndex[:], "Primary Index:")
	dumpList(trackPointer[:], "Track Pointers:")
	dumpList(instrumentPointers[:], "Instrument Pointers:")
	return ret, nil

}
