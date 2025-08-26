package main

import (
	"os"

	"github.com/nibrahim/eye-of-the-gopher/internal/utils"
	"github.com/nibrahim/eye-of-the-gopher/pkg/formats"
)

func main() {

	utils.SetupLogging("maz-decode.log")
	if len(os.Args) != 3 {
		utils.ErrorAndExit("Usage : ./decodemaz mazFile outputFile")
	}

	dataFile := os.Args[1]
	outputFile := os.Args[2]

	mazData, err := os.ReadFile(dataFile)
	if err != nil {
		utils.ErrorAndExit("Can't read data file %s", dataFile)
	}

	plan := formats.DecodeMAZ(dataFile, mazData)
	// plan.Validate()
	plan.DrawPlan(outputFile)

}
