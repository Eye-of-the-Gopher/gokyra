package formats

import (
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Sprite struct {
	Image image.Image
	name  string
}

type Asset struct {
	data []byte
}

func NewAsset(fname string, data []byte) Asset {
	return Asset{
		data: data,
	}
}

type Assets struct {
	assets map[string][]byte
}

func NewAssets() *Assets {
	return &Assets{
		assets: make(map[string][]byte),
	}
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

func (a *Assets) GetSprite(name string, palette color.Palette, width int, height int, prefix string) (*Sprite, error) {
	ext := strings.ToLower(path.Ext(name))
	PakLogger.Debug("Loading sprite", "name", name, "extension", ext)
	if ext != ".cmp" && ext != ".cps" {
		return nil, fmt.Errorf("Cannot fetch %s as a sprite. Only CPS and CMP", name)
	} else {
		data, exists := a.assets[name]
		if exists {
			imgData := DecodeCmp(name, data, palette)
			img := CMPToImage(imgData, palette, width, height, 4)
			PakLogger.Debug("Sending back", "len", len(data))
			return &Sprite{
				name:  name,
				Image: img,
			}, nil
		} else {
			return nil, fmt.Errorf("Cannot fetch %s: No such asset", name)
		}

	}

}

func (a *Assets) DumpAssets() {
	for k := range a.assets {
		fmt.Println(k)
	}
}

func (a *Assets) LoadPakFile(pakfile string, prefix string) error {
	filenames := []string{}
	offsets := []uint32{}
	PakLogger.Info("Loading Pakfile", "name", pakfile)

	f, err := os.Open(pakfile)
	if err != nil {
		return fmt.Errorf("cannot open file: %w", err)
	}
	defer f.Close()

	for { // The header will end when have some offsets/files and when we read past the first offset (i.e. when the data begins)
		if len(offsets) != 0 {
			firstOffset := offsets[0]
			t, err := f.Seek(0, io.SeekCurrent)
			if err != nil {
				return fmt.Errorf("error while finding current file position: %w", err)
			}

			currentPos := uint32(t)
			if currentPos+4 >= firstOffset {
				PakLogger.Debug("header completed")
				break
			}

		}

		var offset uint32
		err = binary.Read(f, binary.LittleEndian, &offset)
		if err != nil {
			return fmt.Errorf("unexpected short read while looking for offset: %w", err)
		}

		offsets = append(offsets, offset)

		var fnamechar byte

		err = binary.Read(f, binary.LittleEndian, &fnamechar)
		if err != nil {
			return fmt.Errorf("unexpected short read while looking for name: %w", err)
		}

		fnamechars := make([]byte, 0)
		fnamechars = append(fnamechars, fnamechar)
		for {
			err = binary.Read(f, binary.LittleEndian, &fnamechar)
			if fnamechar == 0 {
				break
			}
			if err != nil {
				return fmt.Errorf("unexpected short read while looking for name: %w", err)
			}
			fnamechars = append(fnamechars, fnamechar)
		}
		fname := string(fnamechars)
		filenames = append(filenames, fname)
		PakLogger.Debug("Entry parsed", "offset", offset, "name", fname)
	}

	stat, err := f.Stat()
	if err != nil {
		return fmt.Errorf("could not stat file: %w", err)
	}

	fileLimit := uint64(stat.Size())

	for i := range len(offsets) {
		filename := filenames[i]
		start := uint64(offsets[i])
		var end uint64
		if i+1 == len(offsets) {
			end = fileLimit
		} else {
			end = uint64(offsets[i+1])
		}
		PakLogger.Debug("Extracting", "file", filename, "From", start, "To", end)
		_, err = f.Seek(int64(start), io.SeekStart)
		if err != nil {
			return fmt.Errorf("could not extract %s from %s at offset %d: %w", filename, pakfile, start, err)
		}
		data := make([]byte, end-start)
		n, err := f.Read(data)
		if err != nil || n != len(data) {
			return fmt.Errorf("short read while unpacking %s (position : %d): %w", filename, i, err)

		}
		if prefix != "" {
			filename = fmt.Sprintf("%s/%s", prefix, filename)
		}

		a.assets[filename] = data
	}
	return nil
}

func (a *Assets) LoadExtraAssets(baseDir string) error {
	assets, err := os.ReadDir(baseDir)
	if err != nil {
		return fmt.Errorf("Couldn't side load extra assets: Couldn't read %s", baseDir)
	}
	for _, asset := range assets {
		assetFile := path.Join(baseDir, asset.Name())
		PakLogger.Debug("Sideloading", "file", assetFile)
	}
	return nil
}

func (a *Assets) WriteAssetData(basedir string) {
	os.Mkdir(basedir, 0755)

	for name, data := range a.assets {
		opfile := filepath.Join(basedir, name)
		f, err := os.Create(opfile)

		if err != nil {
			PakLogger.Error("Couldn't write", "name", name)
		} else {
			defer f.Close()
			PakLogger.Debug("Writing", "name", name, "size", len(data))
			f.Write(data)
		}
	}
}
