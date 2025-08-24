package utils

import (
	"fmt"
	"io"
	"log/slog"
	"math/rand/v2"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Creates log file and sets up logging properly
func SetupLogging(logfile string) {
	logFile, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		ErrorAndExit("Couldnt open log file: %v", err)
	}

	multiWriter := io.MultiWriter(os.Stderr, logFile)

	logger := slog.New(slog.NewTextHandler(
		multiWriter,
		&slog.HandlerOptions{
			Level: slog.LevelDebug,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.TimeKey {
					return slog.Attr{}
				}
				return a
			},
		}))
	slog.SetDefault(logger)
}

// Displays []data as a binary pattern
func BytesToBinary(data []byte) string {
	var result strings.Builder
	for i, b := range data {
		if i > 0 {
			result.WriteString(" ") // Space between bytes
		}
		result.WriteString(fmt.Sprintf("%08b", b))
	}
	return result.String()
}

// Prints given error message on stderr and exits
func ErrorAndExit(message string, args ...any) {
	fmt.Fprintf(os.Stderr, "\n"+message+"\n", args...)
	os.Exit(-1)
}

// Takes fName.whatever and returns opDir/fName.ext
func ImageName(fName string, newExt string, opDir string) string {
	base := strings.ToLower(filepath.Base(fName))
	ext := filepath.Ext(base)
	withoutExt := strings.Trim(base, ext)
	return path.Join(opDir, withoutExt+"."+newExt)
}

// Similar to Python's random.choice
func RandChoice[T any](l []T, rng ...*rand.Rand) T {
	var randInt func(int) int
	if len(rng) == 1 {
		t := rng[0]
		randInt = t.IntN
	} else {
		randInt = rand.IntN
	}
	return l[randInt(len(l))]
}
