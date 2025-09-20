package formats

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"log/slog"
	"path"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/nfnt/resize"
	"github.com/nibrahim/eye-of-the-gopher/internal/utils"
)

type AssetLoaderConfig struct {
	AssetLevel slog.Level
	CmpLevel   slog.Level
	MazLevel   slog.Level
	PakLevel   slog.Level
	PalLevel   slog.Level
}

var (
	AssetsLogger *slog.Logger
	CmpLogger    *slog.Logger
	MazLogger    *slog.Logger
	PakLogger    *slog.Logger
	PalLogger    *slog.Logger
)

type Assets struct {
	assets map[string][]byte
}

func NewAssets() *Assets {
	return &Assets{
		assets: make(map[string][]byte),
	}
}

type Sprite struct {
	Image image.Image
	name  string
}

func (s *Sprite) GetEbitenImage() (*ebiten.Image, error) {
	ret := ebiten.NewImageFromImage(s.Image)
	return ret, nil
}

func (a *Assets) GetPalette(name string) (color.Palette, error) {
	ext := strings.ToLower(path.Ext(name))
	PakLogger.Debug("Loading palette", "name", name, "extension", ext)
	if ext != ".pal" && ext != ".col" {
		return nil, fmt.Errorf("Cannot fetch %s as a sprite. Only PAL and COL", name)
	} else {
		data, exists := a.assets[name]
		if exists {
			pal := DecodePalette(data)
			return pal, nil
		} else {
			return nil, fmt.Errorf("Cannot fetch %s: No such asset", name)
		}
	}
}

func (a *Assets) GetSprite(name string, palette color.Palette, width uint, height uint, prefix string) (*Sprite, error) {
	ext := strings.ToLower(path.Ext(name))
	PakLogger.Debug("Loading sprite", "name", name, "extension", ext)
	switch ext {
	case ".cmp", ".cps":
		data, exists := a.assets[name]
		if exists {
			imgData := DecodeCmp(name, data, palette)
			img := CMPToImage(imgData, palette, int(width), int(height), 4)
			PakLogger.Debug("Sending back", "len", len(data))
			return &Sprite{
				name:  name,
				Image: img,
			}, nil
		} else {
			return nil, fmt.Errorf("Cannot fetch %s: No such asset", name)
		}
	case ".png":
		data, exists := a.assets[name]
		if exists {
			imgData := bytes.NewReader(data)
			img, format, err := image.Decode(imgData)
			img = resize.Resize(width*4, height*4, img, resize.Lanczos3)
			PakLogger.Debug("Decoding image ", "name", name, "format", format)
			if err != nil {
				PakLogger.Error("Couldn't decode image", "image", name)
				return nil, fmt.Errorf("Couldn't decoder image %s", name)
			}
			PakLogger.Debug("Sending back", "len", len(data))
			return &Sprite{
				name:  name,
				Image: img,
			}, nil
		} else {
			return nil, fmt.Errorf("Cannot fetch %s: No such asset", name)
		}
	default:
		return nil, fmt.Errorf("Cannot fetch %s as a sprite. Only CPS, CMP or PNG", name)
	}

}

func (a *Assets) DumpAssets() {
	for k := range a.assets {
		fmt.Println(k)
	}
}

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
	AssetsLogger.Debug("Loading extra assets")
	for _, assetDir := range extraAssetDirs {
		prefix := path.Base(assetDir)
		AssetsLogger.Debug("Loading extra assets", "from", assetDir, "prefix", prefix)
		ret.LoadExtraAssets(assetDir, prefix)
	}

	return ret
}

func InitLogger(assetLogLevels AssetLoaderConfig) {
	CmpLogger = utils.InitLogger("cmp", assetLogLevels.CmpLevel)
	MazLogger = utils.InitLogger("maz", assetLogLevels.MazLevel)
	PakLogger = utils.InitLogger("pak", assetLogLevels.PakLevel)
	PalLogger = utils.InitLogger("pal", assetLogLevels.PalLevel)
	AssetsLogger = utils.InitLogger("assets", assetLogLevels.AssetLevel)
}
