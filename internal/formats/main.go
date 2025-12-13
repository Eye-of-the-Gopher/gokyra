package formats

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"log/slog"
	"os"
	"path"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hajimehoshi/ebiten/v2/vector"
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

// type PixelIterator func() (image.Point, bool)
type FadeIterator func() (*ebiten.Image, bool)

func NewAssets() *Assets {
	return &Assets{
		assets: make(map[string][]byte),
	}
}

type Sprite struct {
	Image image.Image
	name  string
}

func (s *Sprite) GetEbitenImage() *ebiten.Image {
	ret := ebiten.NewImageFromImage(s.Image)
	return ret
}

func (s *Sprite) GetImageRegion(x0, y0, x1, y1 int) *Sprite {
	img := ebiten.NewImageFromImage(s.Image)
	srcRect := image.Rect(x0, y0, x1, y1)
	subImage := img.SubImage(srcRect)
	return &Sprite{
		Image: subImage,
		name:  s.name,
	}
}

func (s *Sprite) GetEbitenImageFadeOut(width int, count int) FadeIterator {
	baseImg := s.GetEbitenImage()
	baseImgX := baseImg.Bounds().Dx()
	baseImgY := baseImg.Bounds().Dy()
	nSquaresX := baseImgX / width
	nSquaresY := baseImgY / width
	img := ebiten.NewImageFromImage(baseImg)
	offset := -4
	ret := func() (*ebiten.Image, bool) {
		if count != 0 {
			count -= 1
			for x := 0; x < nSquaresX; x++ {
				for y := 0; y < nSquaresY; y++ {
					vector.StrokeRect(img, float32(x*width+offset), float32(y*width+offset), float32(width),
						float32(width), 1.0, color.Black, false)
					// AssetsLogger.Debug("FadeIterator going on", "count", count)
				}
			}
			// SaveEbitenImage(img, fmt.Sprintf("/tmp/image-%d.png", count))
			offset += 3
			return img, true
		} else {
			return nil, false
		}

	}
	return ret
}

func (s *Sprite) GetEbitenImageFadeIn(width int, count int) FadeIterator {
	baseImg := s.GetEbitenImage()
	baseImgX := baseImg.Bounds().Dx()
	baseImgY := baseImg.Bounds().Dy()
	nSquaresX := baseImgX / width
	nSquaresY := baseImgY / width

	mask := ebiten.NewImage(baseImgX, baseImgY)

	offset := -4

	ret := func() (*ebiten.Image, bool) {
		if count != 0 {
			count -= 1
			for x := 0; x < nSquaresX; x++ {
				for y := 0; y < nSquaresY; y++ {
					vector.StrokeRect(mask, float32(x*width+offset), float32(y*width+offset), float32(width),
						float32(width), 1.0, color.White, false)
				}
			}

			// Create result by masking baseImg
			result := ebiten.NewImage(baseImgX, baseImgY)

			// Draw base image WITH mask applied
			op := &ebiten.DrawImageOptions{}
			op.ColorScale.ScaleWithColor(color.White)
			result.DrawImage(mask, nil) // Draw mask first

			op2 := &ebiten.DrawImageOptions{}
			op2.Blend = ebiten.BlendSourceIn
			result.DrawImage(baseImg, op2) // Apply image through mask

			// img := ebiten.NewImage(baseImgX, baseImgY)
			// img.DrawImage(baseImg, nil)

			// op := &ebiten.DrawImageOptions{}
			// op.Blend = ebiten.BlendSourceIn // Mask erases pixels
			// img.DrawImage(mask, op)

			// SaveEbitenImage(result, fmt.Sprintf("/tmp/image-%d.png", count))

			offset += 3
			return result, true
		} else {
			return nil, false
		}

	}
	return ret
}

type AudioTrack struct {
	track  string
	data   []byte
	format string
}

func (a *AudioTrack) GetEbintenPlayer(ctx *audio.Context) (*audio.Player, error) {
	var stream io.ReadSeeker
	var err error

	reader := bytes.NewReader(a.data)
	switch a.format {
	case "mp3":
		stream, err = mp3.DecodeWithSampleRate(ctx.SampleRate(), reader)
	case "adl":
		AssetsLogger.Info("ADL currently unsupported")
		stream = nil
		err = fmt.Errorf("ADL currently unsupported. Please convert separately and side load")
	case "wav":
		stream, err = wav.DecodeWithSampleRate(ctx.SampleRate(), reader)
	default:
		AssetsLogger.Error("unknown format", "format", a.format)
		stream = nil
		err = fmt.Errorf("unknown format %s", a.format)
	}
	if err != nil {
		return nil, err
	}

	ret, err := ctx.NewPlayer(stream)
	if err != nil {
		return nil, fmt.Errorf("couldn't create player: %v", err)
	}
	return ret, err

}

func (a *AudioTrack) String() string {
	return fmt.Sprintf("Track : %s", a.track)
}

func (a *Assets) GetAudioTrack(name string) (*AudioTrack, error) {
	ext := strings.ToLower(path.Ext(name))
	AssetsLogger.Debug("Loading track", "name", name, "extension", ext)
	data, exists := a.assets[name]
	if exists {
		switch ext {
		case ".adl":
			AssetsLogger.Info("ADL currently unsupported")
			return nil, nil
		case ".mp3":
			return &AudioTrack{
				track:  name,
				data:   data,
				format: "mp3",
			}, nil
		case ".wav":
			return &AudioTrack{
				track:  name,
				data:   data,
				format: "wav",
			}, nil
		default:
			return nil, fmt.Errorf("cannot fetch %s as an Audio track. Only ADL or mp3", name)
		}
	} else {

		return nil, fmt.Errorf("cannot fetch %s: No such asset", name)
	}
}

func (a *Assets) GetPalette(name string) (color.Palette, error) {
	ext := strings.ToLower(path.Ext(name))
	PakLogger.Debug("Loading palette", "name", name, "extension", ext)
	if ext != ".pal" && ext != ".col" {
		return nil, fmt.Errorf("cannot fetch %s as a sprite. Only PAL and COL", name)
	} else {
		data, exists := a.assets[name]
		if exists {
			pal := DecodePalette(data)
			return pal, nil
		} else {
			return nil, fmt.Errorf("cannot fetch %s: No such asset", name)
		}
	}
}

func (a *Assets) GetSprite(name string, paletteName string, width uint, height uint, prefix string) (*Sprite, error) {
	palette, err := a.GetPalette(paletteName)
	if err != nil {
		PakLogger.Error("Couldn't load palette", "name", paletteName)
	}
	ext := strings.ToLower(path.Ext(name))
	PakLogger.Debug("Loading sprite", "name", name, "extension", ext)
	switch ext {
	case ".cmp", ".cps":
		data, exists := a.assets[name]
		if exists {
			imgData := DecodeCmp(name, data, palette)
			img := CMPToImage(imgData, palette, int(width), int(height))
			PakLogger.Debug("Sending back", "len", len(data))
			return &Sprite{
				name:  name,
				Image: img,
			}, nil
		} else {
			PakLogger.Warn("Couldn't load ", "name", name)
			return nil, fmt.Errorf("cannot fetch %s: No such asset", name)
		}
	case ".png":
		data, exists := a.assets[name]
		if exists {
			imgData := bytes.NewReader(data)
			img, format, err := image.Decode(imgData)
			img = resize.Resize(width, height, img, resize.Lanczos3)
			PakLogger.Debug("Decoding image ", "name", name, "format", format)
			if err != nil {
				PakLogger.Error("couldn't decode image", "image", name)
				return nil, fmt.Errorf("couldn't decoder image %s", name)
			}
			PakLogger.Debug("Sending back", "len", len(data))
			return &Sprite{
				name:  name,
				Image: img,
			}, nil
		} else {
			PakLogger.Warn("Couldn't load ", "name", name)
			return nil, fmt.Errorf("cannot fetch %s: No such asset", name)
		}
	default:
		return nil, fmt.Errorf("cannot fetch %s as a sprite. Only CPS, CMP or PNG", name)
	}

}

func (a *Assets) DumpAssets() {
	for k := range a.assets {
		fmt.Println(k)
	}
}

func SaveEbitenImage(img *ebiten.Image, filename string) error {
	// Convert to standard image
	fmt.Println("Writing ", filename)
	rgba := image.NewRGBA(img.Bounds())
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			rgba.Set(x, y, img.At(x, y))
		}
	}

	// Save as PNG
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, rgba)
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
