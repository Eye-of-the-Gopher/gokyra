package formats

import (
	"log/slog"
	"path"

	"github.com/nibrahim/eye-of-the-gopher/internal/utils"
)

// Load assets from original EOB game. Game files should be in the provided directory
func LoadAssets(classicAssetDir string, extraAssetDirs ...string) *Assets {
	ret := NewAssets()
	pakFilesNeeded := []string{"EOBDATA1.PAK", "EOBDATA2.PAK", "EOBDATA3.PAK", "EOBDATA4.PAK", "EOBDATA5.PAK", "EOBDATA6.PAK"}
	for _, pakFile := range pakFilesNeeded {
		t := path.Join(classicAssetDir, pakFile)
		err := ret.LoadPakFile(t, "")
		if err != nil {
			utils.ErrorAndExit("Couldn't load %s", t)
		}
	}
	if len(extraAssetDirs) == 0 {
		return ret
	}
	slog.Debug("Loading extra assets")
	for _, assetDir := range extraAssetDirs {
		slog.Debug("Loading from ", "dir", assetDir)
		ret.LoadExtraAssets(assetDir)
	}

	return ret
}
