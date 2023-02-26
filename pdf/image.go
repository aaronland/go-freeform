package pdf

import (
	"context"
	"fmt"
	"io"
	"image"
	_ "image/png"
	
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
)

func Images(ctx context.Context, r io.ReadSeeker) (map[string]string, error) {

	pages := []string{}
	conf := &pdfcpu.Configuration{}

	images, err := api.ExtractImagesRaw(r, pages, conf)

	if err != nil {
		return nil, fmt.Errorf("Failed to extract images, %w", err)
	}

	for _, im := range images {
		fmt.Println(im.Name, im.FileType)

		im2, im_fmt, err := image.Decode(im)

		if err != nil {
			return nil, fmt. Errorf("Failed to decode image, %w", err)
		}

		fmt.Println(im2.Bounds(), im_fmt)
	}

	return nil, nil
}
