package formats

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

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

func (a *Assets) LoadExtraAssets(baseDir string, prefix string) error {
	assets, err := os.ReadDir(baseDir)
	if err != nil {
		return fmt.Errorf("couldn't side load extra assets: Couldn't read %s", baseDir)
	}
	for _, asset := range assets {
		assetFile := path.Join(baseDir, asset.Name())
		AssetsLogger.Debug("Sideloading", "file", assetFile)
		key := asset.Name()
		if prefix != "" {
			key = fmt.Sprintf("%s/%s", prefix, key)
		}
		data, err := os.ReadFile(assetFile)
		if err != nil {
			AssetsLogger.Warn("Could not load", "file", assetFile)
			return err
		}
		key = strings.ToUpper(key)
		a.assets[key] = data
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
