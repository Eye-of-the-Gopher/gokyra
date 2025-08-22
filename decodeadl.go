package main

import (
	"os"
)

func main() {
	setupLogging("adl-decode.log")
	if len(os.Args) != 3 {
		errorAndExit("Usage : ./decodeadl dataFile outputFile")
	}

	dataFile := os.Args[1]
	outputFile := os.Args[2]

	AdlData, err := os.ReadFile(dataFile)
	if err != nil {
		errorAndExit("Can't read data file %s", dataFile)
	}

	decompressedData, err := decodeADL(dataFile, AdlData)
	if err != nil {
		errorAndExit("Error processing file %s: %v", dataFile, err)
	}

	err = writeMIDI(decompressedData, outputFile)
	if err != nil {
		errorAndExit("Can't write converted music data to file %s:", outputFile, err)
	}

}
