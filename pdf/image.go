package pdf

import (
	"context"
	"fmt"
	"image"
	_ "image/png"
	"io"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
)

func Images(ctx context.Context, r io.ReadSeeker) ([]image.Image, error) {

	pages := []string{}
	conf := &pdfcpu.Configuration{}

	raw_images, err := api.ExtractImagesRaw(r, pages, conf)

	if err != nil {
		return nil, fmt.Errorf("Failed to extract images, %w", err)
	}

	images := make([]image.Image, len(raw_images))

	for idx, raw_im := range raw_images {

		im, _, err := image.Decode(raw_im)

		if err != nil {
			return nil, fmt.Errorf("Failed to decode image, %w", err)
		}

		images[idx] = im
	}

	return images, nil
}
