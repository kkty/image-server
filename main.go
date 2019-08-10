package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/kkty/image-compression/convert"
)

func parse(w http.ResponseWriter, r *http.Request) (*convert.Input, error) {
	query := r.URL.Query()

	widthStr := query.Get("width")
	heightStr := query.Get("height")

	width := 0
	height := 0

	if widthStr != "" {
		var err error

		width, err = strconv.Atoi(widthStr)

		if err != nil {
			return &convert.Input{}, err
		}
	}

	if heightStr != "" {
		var err error

		height, err = strconv.Atoi(heightStr)

		if err != nil {
			return &convert.Input{}, err
		}
	}

	qualityStr := query.Get("quality")
	quality := 0

	if qualityStr != "" {
		var err error

		quality, err = strconv.Atoi(qualityStr)

		if err != nil {
			return &convert.Input{}, err
		}
	}

	return &convert.Input{
		Src: struct {
			ContentType string
			R           io.Reader
		}{r.Header.Get("content-type"), r.Body},
		Dst: struct {
			Width       int
			Height      int
			Quality     int
			ContentType string
			W           io.Writer
		}{width, height, quality, r.Header.Get("accept"), w},
	}, nil
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		input, err := parse(w, r)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := convert.Execute(input); err != nil {
			if err == convert.ErrUnsupportedContentType {
				w.WriteHeader(http.StatusBadRequest)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}

			return
		}
	})

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
