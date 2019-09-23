package convert

import (
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"strconv"

	"golang.org/x/image/draw"
)

// ErrUnsupportedContentType is used when the given content type is not supported.
var ErrUnsupportedContentType = errors.New("unsupported content-type")

// Config is used to specify the type of conversion to be made.
type Config struct {
	Width           int
	Height          int
	Quality         int
	ContentTypeFrom string
	ContentTypeTo   string
}

// ParseHTTPRequest parses a http.Request and constructs a Config.
func ParseHTTPRequest(r *http.Request) (*Config, error) {
	query := r.URL.Query()

	widthStr := query.Get("width")
	heightStr := query.Get("height")

	width := 0
	height := 0

	if widthStr != "" {
		var err error

		width, err = strconv.Atoi(widthStr)

		if err != nil {
			return nil, err
		}
	}

	if heightStr != "" {
		var err error

		height, err = strconv.Atoi(heightStr)

		if err != nil {
			return nil, err
		}
	}

	qualityStr := query.Get("quality")
	quality := 0

	if qualityStr != "" {
		var err error

		quality, err = strconv.Atoi(qualityStr)

		if err != nil {
			return nil, err
		}
	}

	return &Config{
		Width:           width,
		Height:          height,
		Quality:         quality,
		ContentTypeFrom: r.Header.Get("content-type"),
		ContentTypeTo:   r.Header.Get("accept"),
	}, nil
}

// Convert reads, decodes and resizes the given image.
// Then, the image gets encoded and written.
func Convert(config *Config, r io.Reader, w io.Writer) error {
	var src image.Image

	switch config.ContentTypeFrom {
	case "image/png":
		var err error
		src, err = png.Decode(r)
		if err != nil {
			return err
		}
	case "image/jpeg":
		var err error
		src, err = jpeg.Decode(r)
		if err != nil {
			return err
		}
	case "image/gif":
		var err error
		src, err = gif.Decode(r)
		if err != nil {
			return err
		}
	default:
		return ErrUnsupportedContentType
	}

	sr := src.Bounds()

	var width, height int
	if config.Width == 0 && config.Height == 0 {
		var sz = sr.Size()
		width, height = sz.X, sz.Y
	} else if config.Width == 0 {
		var sz = sr.Size()
		height = config.Height
		width = height * sz.X / sz.Y
	} else if config.Height == 0 {
		var sz = sr.Size()
		width = config.Width
		height = width * sz.Y / sz.X
	} else {
		width, height = config.Width, config.Height
	}

	dr := image.Rect(0, 0, width, height)
	dst := image.NewRGBA(dr)

	draw.BiLinear.Scale(dst, dr, src, sr, draw.Over, nil)

	switch config.ContentTypeTo {
	case "image/png":
		if err := png.Encode(w, dst); err != nil {
			return err
		}
	case "image/jpeg":
		options := jpeg.Options{
			Quality: 100,
		}

		if config.Quality != 0 {
			options.Quality = config.Quality
		}

		if err := jpeg.Encode(w, dst, &options); err != nil {
			return err
		}
	case "image/gif":
		if err := gif.Encode(w, dst, nil); err != nil {
			return err
		}
	default:
		return ErrUnsupportedContentType
	}

	return nil
}
