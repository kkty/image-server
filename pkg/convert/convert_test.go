package convert

import (
	"bytes"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"testing"
)

func TestConvert(t *testing.T) {
	srcs := []struct {
		path        string
		contentType string
	}{
		{"../../test/300.jpg", "image/jpeg"},
		{"../../test/300.png", "image/png"},
		{"../../test/300.gif", "image/gif"},
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

				config := Config{
					ContentTypeFrom: src.contentType,
					ContentTypeTo:   dst.contentType,
					Width:           dst.width,
					Height:          dst.height,
					Quality:         dst.quality,
				}

				if err := Convert(&config, file, &buf); err != nil {
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
