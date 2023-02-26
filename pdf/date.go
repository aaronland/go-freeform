package pdf

import (
	"fmt"
	"regexp"
	"time"
)

var pat *regexp.Regexp

func init() {
	pat = regexp.MustCompile(`D\:(\d{4})(\d{2})(\d{2})(\d{2})(\d{2})(\d{2})\+00'00'`)
}

func ParseDate(d string) (*time.Time, error) {

	// Because this just yields one Go error that's weirder than the next...
	// layout := "D:20060102150400+00'00'"
	// return time.Parse(layout, d)

	if !pat.MatchString(d) {
		return nil, fmt.Errorf("Invalid date string")
	}

	m := pat.FindStringSubmatch(d)

	ds := fmt.Sprintf("%s-%s-%sT%s:%s", m[1], m[2], m[3], m[4], m[5])

	layout := "2006-01-02T15:04"

	t, err := time.Parse(layout, ds)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse '%s', %w", ds, err)
	}

	return &t, nil
}
