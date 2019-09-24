package convert

import (
	"bytes"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testParseHTTPRequest(
	t *testing.T,
	query url.Values,
	header http.Header,
	assertFn func(t *testing.T, config *Config),
) {
	req, err := http.NewRequest("POST", "http://example.com/", &bytes.Buffer{})

	if err != nil {
		t.Fatal(err)
	}

	req.URL.RawQuery = query.Encode()
	req.Header = header

	config, err := ParseHTTPRequest(req)

	if err != nil {
		t.Fatal(err)
	}

	assertFn(t, config)
}

func TestParseHTTPRequest(t *testing.T) {
	t.Run("WithoutQuery", func(t *testing.T) {
		testParseHTTPRequest(
			t,
			url.Values{},
			http.Header{
				"Content-Type": []string{"image/png"},
				"Accept":       []string{"image/jpeg"},
			},
			func(t *testing.T, config *Config) {
				assert.Equal(t, "image/png", config.ContentTypeFrom)
				assert.Equal(t, "image/jpeg", config.ContentTypeTo)
				assert.Equal(t, 0, config.Width)
				assert.Equal(t, 0, config.Height)
				assert.Equal(t, 0, config.Quality)
			},
		)
	})

	t.Run("WithQuery", func(t *testing.T) {
		testParseHTTPRequest(
			t,
			url.Values{
				"width":   []string{"100"},
				"height":  []string{"50"},
				"quality": []string{"30"},
			},
			http.Header{
				"Content-Type": []string{"image/png"},
				"Accept":       []string{"image/jpeg"},
			},
			func(t *testing.T, config *Config) {
				assert.Equal(t, "image/png", config.ContentTypeFrom)
				assert.Equal(t, "image/jpeg", config.ContentTypeTo)
				assert.Equal(t, 100, config.Width)
				assert.Equal(t, 50, config.Height)
				assert.Equal(t, 30, config.Quality)
			},
		)
	})
}
func TestConvert(t *testing.T) {
	configs := []Config{}

	for _, width := range []int{0, 100} {
		for _, height := range []int{0, 100} {
			for _, quality := range []int{0, 50} {
				for _, contentTypeFrom := range []string{"image/jpeg", "image/png", "image/gif"} {
					for _, contentTypeTo := range []string{"image/jpeg", "image/png", "image/gif"} {
						configs = append(configs, Config{width, height, quality, contentTypeFrom, contentTypeTo})
					}
				}
			}
		}
	}

	for _, config := range configs {
		var fileName string

		switch config.ContentTypeFrom {
		case "image/png":
			fileName = "../../test/300.png"
		case "image/gif":
			fileName = "../../test/300.gif"
		case "image/jpeg":
			fileName = "../../test/300.jpg"
		default:
			t.Fatal()
		}

		file, err := os.Open(fileName)

		if err != nil {
			t.Fatal(err)
		}

		buf := bytes.Buffer{}

		if err := Convert(&config, file, &buf); err != nil {
			t.Fatal(err)
		}

		var img image.Image

		switch config.ContentTypeTo {
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

		assert.Greater(t, sz.X, 0)
		assert.Greater(t, sz.Y, 0)

		if config.Width != 0 {
			assert.Equal(t, config.Width, sz.X)
		}

		if config.Height != 0 {
			assert.Equal(t, config.Height, sz.Y)
		}
	}
}
