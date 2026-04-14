package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"image/png"
	_ "image/png"
	"math"
	"os"
)

type iconEntry struct {
	Size int
	PNG  []byte
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	const sourcePath = "internal/assets/tray_icon_nonwindows.png"
	const targetPath = "internal/assets/tray_icon_windows.ico"

	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("open source icon: %w", err)
	}
	defer sourceFile.Close()

	source, _, err := image.Decode(sourceFile)
	if err != nil {
		return fmt.Errorf("decode source icon: %w", err)
	}

	sizes := []int{16, 32, 48, 64, 128, 256}
	entries := make([]iconEntry, 0, len(sizes))
	for _, size := range sizes {
		payload, err := renderSquarePNG(source, size)
		if err != nil {
			return fmt.Errorf("render %dx%d icon: %w", size, size, err)
		}
		entries = append(entries, iconEntry{Size: size, PNG: payload})
	}

	payload, err := buildICO(entries)
	if err != nil {
		return fmt.Errorf("build ico: %w", err)
	}

	if err := os.WriteFile(targetPath, payload, 0o644); err != nil {
		return fmt.Errorf("write ico: %w", err)
	}

	return nil
}

func renderSquarePNG(source image.Image, size int) ([]byte, error) {
	bounds := source.Bounds()
	sourceWidth := bounds.Dx()
	sourceHeight := bounds.Dy()
	if sourceWidth == 0 || sourceHeight == 0 {
		return nil, fmt.Errorf("source image has invalid dimensions %dx%d", sourceWidth, sourceHeight)
	}

	scale := math.Min(float64(size)/float64(sourceWidth), float64(size)/float64(sourceHeight))
	drawWidth := max(1, int(math.Round(float64(sourceWidth)*scale)))
	drawHeight := max(1, int(math.Round(float64(sourceHeight)*scale)))
	offsetX := (size - drawWidth) / 2
	offsetY := (size - drawHeight) / 2

	canvas := image.NewNRGBA(image.Rect(0, 0, size, size))
	transparent := color.NRGBA{A: 0}
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			canvas.SetNRGBA(x, y, transparent)
		}
	}

	for y := 0; y < drawHeight; y++ {
		sourceY := bounds.Min.Y + min(sourceHeight-1, int(float64(y)*float64(sourceHeight)/float64(drawHeight)))
		for x := 0; x < drawWidth; x++ {
			sourceX := bounds.Min.X + min(sourceWidth-1, int(float64(x)*float64(sourceWidth)/float64(drawWidth)))
			canvas.Set(offsetX+x, offsetY+y, source.At(sourceX, sourceY))
		}
	}

	var buffer bytes.Buffer
	if err := png.Encode(&buffer, canvas); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func buildICO(entries []iconEntry) ([]byte, error) {
	var buffer bytes.Buffer
	if err := binary.Write(&buffer, binary.LittleEndian, uint16(0)); err != nil {
		return nil, err
	}
	if err := binary.Write(&buffer, binary.LittleEndian, uint16(1)); err != nil {
		return nil, err
	}
	if err := binary.Write(&buffer, binary.LittleEndian, uint16(len(entries))); err != nil {
		return nil, err
	}

	offset := 6 + (16 * len(entries))
	for _, entry := range entries {
		dimension := byte(entry.Size)
		if entry.Size >= 256 {
			dimension = 0
		}
		if err := buffer.WriteByte(dimension); err != nil {
			return nil, err
		}
		if err := buffer.WriteByte(dimension); err != nil {
			return nil, err
		}
		if err := buffer.WriteByte(0); err != nil {
			return nil, err
		}
		if err := buffer.WriteByte(0); err != nil {
			return nil, err
		}
		if err := binary.Write(&buffer, binary.LittleEndian, uint16(1)); err != nil {
			return nil, err
		}
		if err := binary.Write(&buffer, binary.LittleEndian, uint16(32)); err != nil {
			return nil, err
		}
		if err := binary.Write(&buffer, binary.LittleEndian, uint32(len(entry.PNG))); err != nil {
			return nil, err
		}
		if err := binary.Write(&buffer, binary.LittleEndian, uint32(offset)); err != nil {
			return nil, err
		}
		offset += len(entry.PNG)
	}

	for _, entry := range entries {
		if _, err := buffer.Write(entry.PNG); err != nil {
			return nil, err
		}
	}

	return buffer.Bytes(), nil
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
