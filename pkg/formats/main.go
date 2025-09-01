package formats

import (
	"path"

	"github.com/nibrahim/eye-of-the-gopher/internal/utils"
)

// Load assets from original EOB game. Game files should be in the provided directory
func LoadClassicAssets(assetDir string) *Assets {
	ret := NewAssets()
	pakFilesNeeded := []string{"EOBDATA1.PAK", "EOBDATA2.PAK", "EOBDATA3.PAK", "EOBDATA4.PAK", "EOBDATA5.PAK", "EOBDATA6.PAK"}
	for _, pakFile := range pakFilesNeeded {
		t := path.Join(assetDir, pakFile)
		err := ret.LoadPakFile(t, "")
		if err != nil {
			utils.ErrorAndExit("Couldn't load %s", t)
		}
	}
	return ret
}
