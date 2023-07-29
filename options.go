package flag

import (
	"flag"

	"go.unistack.org/micro/v4/options"
)

type sliceDelimKey struct{}

// SliceDelim set the slice delimeter
func SliceDelim(s string) options.Option {
	return options.ContextOption(sliceDelimKey{}, s)
}

type mapDelimKey struct{}

// MapDelim set the map delimeter
func MapDelim(s string) options.Option {
	return options.ContextOption(mapDelimKey{}, s)
}

type timeFormatKey struct{}

// TimeFormat set the time format
func TimeFormat(s string) options.Option {
	return options.ContextOption(timeFormatKey{}, s)
}

type flagSetKey struct{}

// FlagSet set flag set name
func FlagSet(f *flag.FlagSet) options.Option {
	return options.ContextOption(flagSetKey{}, f)
}

type flagSetNameKey struct{}

// FlagSetName set flag set name
func FlagSetName(n string) options.Option {
	return options.ContextOption(flagSetNameKey{}, n)
}

type flagSetErrorHandlingKey struct{}

// FlagErrorHandling set flag set error handling
func FlagErrorHandling(eh flag.ErrorHandling) options.Option {
	return options.ContextOption(flagSetErrorHandlingKey{}, eh)
}

type flagSetUsageKey struct{}

// FlagUsage set flag set usage func
func FlagUsage(fn func()) options.Option {
	return options.ContextOption(flagSetUsageKey{}, fn)
}

type flagEnvKey struct{}

// FlagEnv set flag set usage func
func FlagEnv(n string) options.Option {
	return options.ContextOption(flagEnvKey{}, n)
}
