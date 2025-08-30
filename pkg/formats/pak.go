package formats

import (
	"encoding/binary"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

type Assets struct {
	assets map[string][]byte
}

func NewAssets() *Assets {
	return &Assets{
		assets: make(map[string][]byte),
	}
}

func (a *Assets) DumpAssets() {
	for k, _ := range a.assets {
		fmt.Println(k)
	}
}

func (a *Assets) GetAsset(name string, prefix string) ([]byte, error) {
	value, exists := a.assets[name]
	if !exists {
		return nil, fmt.Errorf("No such asset %s with prefix %s", name, prefix)
	}
	return value, nil
}

func (a *Assets) LoadPakFile(pakfile string, prefix string) error {
	filenames := []string{}
	offsets := []uint32{}
	slog.Info("Loading Pakfile", "name", pakfile)

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
				slog.Debug("header completed")
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
		slog.Debug("Entry parsed", "offset", offset, "name", fname)
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
		slog.Debug("Extracting", "file", filename, "From", start, "To", end)
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

func (a *Assets) WriteAssetData(basedir string) {
	os.Mkdir(basedir, 0755)

	for name, data := range a.assets {
		opfile := filepath.Join(basedir, name)
		f, err := os.Create(opfile)

		if err != nil {
			slog.Error("Couldn't write", "name", name)
		} else {
			defer f.Close()
			slog.Debug("Writing", "name", name, "size", len(data))
			f.Write(data)
		}
	}
}
