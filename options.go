package flag

import (
	"go.unistack.org/micro/v3/config"
)

type sliceDelimKey struct{}

// SliceDelim set the slice delimeter
func SliceDelim(s string) config.Option {
	return config.SetOption(sliceDelimKey{}, s)
}

type mapDelimKey struct{}

// MapDelim set the map delimeter
func MapDelim(s string) config.Option {
	return config.SetOption(mapDelimKey{}, s)
}

type timeFormatKey struct{}

// TimeFormat set the time format
func TimeFormat(s string) config.Option {
	return config.SetOption(timeFormatKey{}, s)
}
