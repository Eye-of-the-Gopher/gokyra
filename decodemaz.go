package main

import "os"

func main() {
	setupLogging("maz-decode.log")
	if len(os.Args) != 3 {
		errorAndExit("Usage : ./decodemaz mazFile outputFile")
	}

	dataFile := os.Args[1]
	outputFile := os.Args[2]

	mazData, err := os.ReadFile(dataFile)
	if err != nil {
		errorAndExit("Can't read data file %s", dataFile)
	}

	decompressedData := decodeMAZ(dataFile, mazData)
	if err != nil {
		errorAndExit("Error processing file %s: %v", dataFile, err)
	}

	drawMap(decompressedData, outputFile)
	if err != nil {
		errorAndExit("Can't write converted music data to file %s:", outputFile, err)
	}

}
