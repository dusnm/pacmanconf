package pacmanconf

import (
	"errors"
	"strconv"
	"strings"
)

var (
	ErrParseValueType = errors.New("setting the option failed")
)

func normalize(str string) string {
	return strings.TrimSpace(str)
}

func parseValueType(target any, value string) error {
	switch tt := target.(type) {
	case *string:
		*tt = value
	case *[]string:
		*tt = append(*tt, value)
	case *int:
		v, err := strconv.Atoi(value)
		if err != nil {
			return err
		}

		*tt = v
	case *bool:
		if value != "" {
			v, err := strconv.ParseBool(value)
			if err != nil {
				return err
			}

			*tt = v
		} else {
			// If flags exist, they are set.
			*tt = true
		}
	}

	return nil
}

// partition
//
// splits s into the elements that match pf and those that don't, preserving order.
func partition[Slice ~[]E, E any](s Slice, pf func(item E) bool) (p1, p2 Slice) {
	for _, item := range s {
		if pf(item) {
			p1 = append(p1, item)
		} else {
			p2 = append(p2, item)
		}
	}

	return p1, p2
}
