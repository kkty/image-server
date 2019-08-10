package convert

import (
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"

	"golang.org/x/image/draw"
)

// ErrUnsupportedContentType is thrown when the given content type is not supported
var ErrUnsupportedContentType = errors.New("unsupported content-type")

// Input specifies where to read the image, how it shuold be processed, and where the result should be written.
type Input struct {
	Src struct {
		ContentType string
		R           io.Reader
	}

	Dst struct {
		Width       int
		Height      int
		Quality     int
		ContentType string
		W           io.Writer
	}
}

// Execute reads, decodes and resizes the given image.
// Then, the image gets encoded and written.
func Execute(input *Input) error {
	var src image.Image

	switch input.Src.ContentType {
	case "image/png":
		var err error
		src, err = png.Decode(input.Src.R)
		if err != nil {
			return err
		}
	case "image/jpeg":
		var err error
		src, err = jpeg.Decode(input.Src.R)
		if err != nil {
			return err
		}
	case "image/gif":
		var err error
		src, err = gif.Decode(input.Src.R)
		if err != nil {
			return err
		}
	default:
		return ErrUnsupportedContentType
	}

	sr := src.Bounds()

	var width, height int
	if input.Dst.Width == 0 && input.Dst.Height == 0 {
		var sz = sr.Size()
		width, height = sz.X, sz.Y
	} else if input.Dst.Width == 0 {
		var sz = sr.Size()
		height = input.Dst.Height
		width = height * sz.X / sz.Y
	} else if input.Dst.Height == 0 {
		var sz = sr.Size()
		width = input.Dst.Width
		height = width * sz.Y / sz.X
	} else {
		width, height = input.Dst.Width, input.Dst.Height
	}

	dr := image.Rect(0, 0, width, height)
	dst := image.NewRGBA(dr)

	draw.BiLinear.Scale(dst, dr, src, sr, draw.Over, nil)

	switch input.Dst.ContentType {
	case "image/png":
		if err := png.Encode(input.Dst.W, dst); err != nil {
			return err
		}
	case "image/jpeg":
		options := jpeg.Options{
			Quality: 100,
		}

		if input.Dst.Quality != 0 {
			options.Quality = input.Dst.Quality
		}

		if err := jpeg.Encode(input.Dst.W, dst, &options); err != nil {
			return err
		}
	case "image/gif":
		if err := gif.Encode(input.Dst.W, dst, nil); err != nil {
			return err
		}
	default:
		return ErrUnsupportedContentType
	}

	return nil
}
