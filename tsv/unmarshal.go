package tsv

import (
	"strings"

	"github.com/boostgo/core/convert"
)

const split = "	"

func Unmarshal(body []byte) ([]string, error) {
	input := convert.StringFromBytes(body)

	records := strings.Split(input, split)

	if len(records) == 0 {
		return nil, ErrNoRecords
	}

	return records, nil
}
