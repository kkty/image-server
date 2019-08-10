package main

import (
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"strconv"

	"golang.org/x/image/draw"
)

var errUnsupportedContentType = errors.New("unsupported content-type")

type input struct {
	src struct {
		contentType string
		r           io.Reader
	}

	dst struct {
		width       int
		height      int
		quality     int
		contentType string
		w           io.Writer
	}
}

func parse(w http.ResponseWriter, r *http.Request) (*input, error) {
	query := r.URL.Query()

	widthStr := query.Get("width")
	heightStr := query.Get("height")

	width := 0
	height := 0

	if widthStr != "" {
		var err error

		width, err = strconv.Atoi(widthStr)

		if err != nil {
			return &input{}, err
		}
	}

	if heightStr != "" {
		var err error

		height, err = strconv.Atoi(heightStr)

		if err != nil {
			return &input{}, err
		}
	}

	qualityStr := query.Get("quality")
	quality := 0

	if qualityStr != "" {
		var err error

		quality, err = strconv.Atoi(qualityStr)

		if err != nil {
			return &input{}, err
		}
	}

	return &input{
		src: struct {
			contentType string
			r           io.Reader
		}{
			contentType: r.Header.Get("content-type"),
			r:           r.Body,
		},
		dst: struct {
			width       int
			height      int
			quality     int
			contentType string
			w           io.Writer
		}{
			contentType: r.Header.Get("accept"),
			width:       width,
			height:      height,
			quality:     quality,
			w:           w,
		},
	}, nil
}

func execute(input *input) error {
	var src image.Image

	switch input.src.contentType {
	case "image/png":
		var err error
		src, err = png.Decode(input.src.r)
		if err != nil {
			return err
		}
	case "image/jpeg":
		var err error
		src, err = jpeg.Decode(input.src.r)
		if err != nil {
			return err
		}
	case "image/gif":
		var err error
		src, err = gif.Decode(input.src.r)
		if err != nil {
			return err
		}
	default:
		return errUnsupportedContentType
	}

	sr := src.Bounds()

	var width, height int
	if input.dst.width == 0 && input.dst.height == 0 {
		var sz = sr.Size()
		width, height = sz.X, sz.Y
	} else if input.dst.width == 0 {
		var sz = sr.Size()
		height = input.dst.height
		width = height * sz.X / sz.Y
	} else if input.dst.height == 0 {
		var sz = sr.Size()
		width = input.dst.width
		height = width * sz.Y / sz.X
	} else {
		width, height = input.dst.width, input.dst.height
	}

	dr := image.Rect(0, 0, width, height)
	dst := image.NewRGBA(dr)

	draw.BiLinear.Scale(dst, dr, src, sr, draw.Over, nil)

	switch input.dst.contentType {
	case "image/png":
		if err := png.Encode(input.dst.w, dst); err != nil {
			return err
		}
	case "image/jpeg":
		options := jpeg.Options{
			Quality: 100,
		}

		if input.dst.quality != 0 {
			options.Quality = input.dst.quality
		}

		if err := jpeg.Encode(input.dst.w, dst, &options); err != nil {
			return err
		}
	case "image/gif":
		if err := gif.Encode(input.dst.w, dst, nil); err != nil {
			return err
		}
	default:
		return errUnsupportedContentType
	}

	return nil
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		input, err := parse(w, r)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := execute(input); err != nil {
			if err == errUnsupportedContentType {
				w.WriteHeader(http.StatusBadRequest)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}

			return
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
