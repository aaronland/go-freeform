package pdf

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
)

func Info(ctx context.Context, r io.ReadSeeker) (map[string]string, error) {

	pages := []string{}
	conf := &pdfcpu.Configuration{
		// Unit: pdfcpu.INCHES,
	}

	info, err := api.Info(r, pages, conf)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive info, %w", err)
	}

	props := make(map[string]string)

	for _, ln := range info {

		parts := strings.SplitN(ln, ":", 2)

		if len(parts) != 2 {
			continue
		}

		k := strings.TrimSpace(parts[0])
		v := strings.TrimSpace(parts[1])

		props[k] = v
	}

	return props, nil
}
