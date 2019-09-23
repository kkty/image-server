package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/kkty/image-server/pkg/convert"
	"go.uber.org/zap"
)

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

		config, err := convert.ParseHTTPRequest(r)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			logger.Error("error_in_parse", zap.Error(err))
			return
		}

		logger = logger.With(
			zap.String("content_type_src", config.ContentTypeFrom),
			zap.String("content_type_dst", config.ContentTypeTo),
			zap.Int("width", config.Width),
			zap.Int("height", config.Height),
			zap.Int("quality", config.Quality),
		)

		if err := convert.Convert(config, r.Body, w); err != nil {
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
