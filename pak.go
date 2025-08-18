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
	fname string
	data  []byte
}

func extractPakFile(pakfile string) ([]PakData, error) {
	fnames := []string{}
	offsets := []uint32{}
	slog.Debug("extracting Pakfile", "name", pakfile)

	f, err := os.Open(pakfile)
	if err != nil {
		return nil, fmt.Errorf("can't open file: %w", err)
	}
	defer f.Close()

	for {
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

		var data uint32
		err = binary.Read(f, binary.LittleEndian, &data)
		if err != nil {
			return nil, fmt.Errorf("unexpected short read while looking for offset: %w", err)
		}

		slog.Debug(fmt.Sprintf("Offset is %d", data))
		offsets = append(offsets, data)

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
		fnames = append(fnames, fname)
		slog.Debug(fmt.Sprintf("Fname is %s", fname))
	}

	ret := make([]PakData, 0)
	stat, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("couldn't stat file: %w", err)
	}

	fileLimit := uint64(stat.Size())

	for i := range len(offsets) {
		fname := fnames[i]
		start := uint64(offsets[i])
		var end uint64
		if i+1 == len(offsets) {
			end = fileLimit
		} else {
			end = uint64(offsets[i+1]) - 1
		}
		slog.Debug("Reading", "fname", fname, "From", start, "To", end)
		_, err = f.Seek(int64(start), io.SeekStart)
		if err != nil {
			return nil, fmt.Errorf("couldn't extract %s from %s at offset %d : %w", fname, pakfile, start, err)
		}
		data := make([]byte, end-start)
		n, err := f.Read(data)
		if err != nil || n != len(data) {
			return nil, fmt.Errorf("short read while unpacking %s (position : %d): %w", fname, i, err)

		}
		t := PakData{
			fname: fname,
			data:  data,
		}
		ret = append(ret, t)
	}
	return ret, nil
}

func writePakData(data []PakData, basedir string) {
	os.Mkdir(basedir, 0755)

	for i := range data {
		pakFile := data[i]
		name := pakFile.fname
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
