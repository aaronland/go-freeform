package pdf

import (
	"testing"
)

func TestParseDate(t *testing.T) {

	tests := map[string]string{
		"D:20230222192640+00'00'": "",
	}

	for d, _ := range tests {

		_, err := ParseDate(d)

		if err != nil {
			t.Fatalf("Failed to parse '%s', %v", d, err)
		}
	}

}
