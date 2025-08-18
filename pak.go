package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

type PakData struct {
	name string
	data []byte
}

func extractPakFile(pakfile string) ([]PakData, error) {
	filenames := []string{}
	offsets := []uint32{}
	slog.Debug("extracting Pakfile", "name", pakfile)

	f, err := os.Open(pakfile)
	if err != nil {
		return nil, fmt.Errorf("cannot open file: %w", err)
	}
	defer f.Close()

	for { // The header will end when have some offsets/files and when we read past the first offset (i.e. when the data begins)
		if len(offsets) != 0 {
			firstOffset := offsets[0]
			t, err := f.Seek(0, io.SeekCurrent)
			if err != nil {
				return nil, fmt.Errorf("error while finding current file position: %w", err)
			}

			currentPos := uint32(t)
			if currentPos+4 >= firstOffset {
				slog.Info("header completed")
				break
			}

		}

		var offset uint32
		err = binary.Read(f, binary.LittleEndian, &offset)
		if err != nil {
			return nil, fmt.Errorf("unexpected short read while looking for offset: %w", err)
		}

		slog.Debug("Offset obtained", "offset", offset)
		offsets = append(offsets, offset)

		var fnamechar byte

		err = binary.Read(f, binary.LittleEndian, &fnamechar)
		if err != nil {
			return nil, fmt.Errorf("unexpected short read while looking for name: %w", err)
		}

		fnamechars := make([]byte, 0)
		fnamechars = append(fnamechars, fnamechar)
		for {
			err = binary.Read(f, binary.LittleEndian, &fnamechar)
			if fnamechar == 0 {
				break
			}
			if err != nil {
				return nil, fmt.Errorf("unexpected short read while looking for name: %w", err)
			}
			fnamechars = append(fnamechars, fnamechar)
		}
		fname := string(fnamechars)
		filenames = append(filenames, fname)
		slog.Debug("Filename obtained ", "name", fname)
	}

	ret := make([]PakData, len(offsets))
	stat, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("could not stat file: %w", err)
	}

	fileLimit := uint64(stat.Size())

	for i := range len(offsets) {
		filename := filenames[i]
		start := uint64(offsets[i])
		var end uint64
		if i+1 == len(offsets) {
			end = fileLimit
		} else {
			end = uint64(offsets[i+1]) - 1
		}
		slog.Debug("Reading", "file", filename, "From", start, "To", end)
		_, err = f.Seek(int64(start), io.SeekStart)
		if err != nil {
			return nil, fmt.Errorf("could not extract %s from %s at offset %d: %w", filename, pakfile, start, err)
		}
		data := make([]byte, end-start)
		n, err := f.Read(data)
		if err != nil || n != len(data) {
			return nil, fmt.Errorf("short read while unpacking %s (position : %d): %w", filename, i, err)

		}
		t := PakData{
			name: filename,
			data: data,
		}
		ret[i] = t
	}
	slog.Info("contents extracted")
	return ret, nil
}

func writePakData(data []PakData, basedir string) {
	os.Mkdir(basedir, 0755)

	for i := range data {
		pakFile := data[i]
		name := pakFile.name
		contents := pakFile.data
		opfile := filepath.Join(basedir, name)
		f, err := os.Create(opfile)

		if err != nil {
			slog.Error("Couldn't write", "name", name)
		} else {
			defer f.Close()
			slog.Debug("Writing", "name", name, "size", len(contents))
			f.Write(contents)
		}
	}
}
