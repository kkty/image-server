package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/kkty/image-server/convert"
	"go.uber.org/zap"
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
	logger, err := zap.NewProduction()

	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}

	defer logger.Sync()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		acceptedAt := time.Now()

		logger := logger.With(
			zap.String("remote_addr", r.RemoteAddr),
		)

		input, err := parse(w, r)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			logger.Error("error_in_parse", zap.Error(err))
			return
		}

		logger = logger.With(
			zap.String("content_type_src", input.Src.ContentType),
			zap.String("content_type_dst", input.Dst.ContentType),
			zap.Int("width", input.Dst.Width),
			zap.Int("height", input.Dst.Height),
			zap.Int("quality", input.Dst.Quality),
		)

		if err := convert.Execute(input); err != nil {
			logger.Error("error_in_convert", zap.Error(err))

			if err == convert.ErrUnsupportedContentType {
				w.WriteHeader(http.StatusBadRequest)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}

		logger.Info(
			"success_convert",
			zap.Duration("elapsed", time.Now().Sub(acceptedAt)),
		)
	})

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
