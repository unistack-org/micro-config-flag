package flag

import (
	"flag"

	"go.unistack.org/micro/v4/config"
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

type flagSetKey struct{}

// FlagSet set flag set name
func FlagSet(f *flag.FlagSet) config.Option {
	return config.SetOption(flagSetKey{}, f)
}

type flagSetNameKey struct{}

// FlagSetName set flag set name
func FlagSetName(n string) config.Option {
	return config.SetOption(flagSetNameKey{}, n)
}

type flagSetErrorHandlingKey struct{}

// FlagErrorHandling set flag set error handling
func FlagErrorHandling(eh flag.ErrorHandling) config.Option {
	return config.SetOption(flagSetErrorHandlingKey{}, eh)
}

type flagSetUsageKey struct{}

// FlagUsage set flag set usage func
func FlagUsage(fn func()) config.Option {
	return config.SetOption(flagSetUsageKey{}, fn)
}

type flagEnvKey struct{}

// FlagEnv set flag set usage func
func FlagEnv(n string) config.Option {
	return config.SetOption(flagEnvKey{}, n)
}
