package main

import (
	"bytes"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"testing"
)

func TestExecute(t *testing.T) {
	srcs := []struct {
		path        string
		contentType string
	}{
		{"./test_assets/300.jpg", "image/jpeg"},
		{"./test_assets/300.png", "image/png"},
		{"./test_assets/300.gif", "image/gif"},
	}

	dsts := []struct {
		width       int
		height      int
		contentType string
		quality     int
	}{
		{100, 100, "image/jpeg", 0},
		{0, 100, "image/jpeg", 0},
		{100, 0, "image/jpeg", 0},
		{0, 0, "image/jpeg", 0},
		{100, 100, "image/png", 0},
		{0, 100, "image/png", 0},
		{100, 0, "image/png", 0},
		{0, 0, "image/png", 0},
		{100, 100, "image/gif", 0},
		{0, 100, "image/gif", 0},
		{100, 0, "image/gif", 0},
		{0, 0, "image/gif", 0},

		{100, 100, "image/jpeg", 50},
	}

	for _, src := range srcs {
		for _, dst := range dsts {
			func() {
				file, err := os.Open(src.path)

				if err != nil {
					t.Fatal(err)
				}

				defer file.Close()

				var buf bytes.Buffer

				input := input{
					src: struct {
						contentType string
						r           io.Reader
					}{src.contentType, file},
					dst: struct {
						width       int
						height      int
						quality     int
						contentType string
						w           io.Writer
					}{dst.width, dst.height, dst.quality, dst.contentType, &buf},
				}

				if err := execute(&input); err != nil {
					t.Fatal(err)
				}

				var img image.Image

				switch dst.contentType {
				case "image/jpeg":
					var err error
					if img, err = jpeg.Decode(&buf); err != nil {
						t.Fatal(err)
					}
				case "image/png":
					var err error
					if img, err = png.Decode(&buf); err != nil {
						t.Fatal(err)
					}
				case "image/gif":
					var err error
					if img, err = gif.Decode(&buf); err != nil {
						t.Fatal(err)
					}
				}

				sz := img.Bounds().Size()

				if (dst.width != 0 && sz.X != dst.width) || (dst.height != 0 && sz.Y != dst.height) {
					t.Fatal("not resized properly")
				}
			}()
		}
	}
}
