package flag

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"go.unistack.org/micro/v4/config"
	rutil "go.unistack.org/micro/v4/util/reflect"
)

var (
	DefaultStructTag  = "flag"
	ErrInvalidValue   = errors.New("invalid value specified")
	DefaultSliceDelim = ","
	DefaultMapDelim   = ","
)

/*
var (
	timeTimeKind     = reflect.TypeOf(time.Time{}).Kind()
	timeDurationKind = reflect.TypeOf(time.Duration(0)).Kind()
)
*/

type flagConfig struct {
	fset *flag.FlagSet
	opts config.Options
	name string
	env  string
}

func (c *flagConfig) Options() config.Options {
	return c.opts
}

func (c *flagConfig) Init(opts ...config.Option) error {
	for _, o := range opts {
		o(&c.opts)
	}

	if err := config.DefaultBeforeInit(c.opts.Context, c); err != nil && !c.opts.AllowFail {
		return err
	}

	c.configure()

	fields, err := rutil.StructFields(c.opts.Struct)
	if err != nil {
		if !c.opts.AllowFail {
			return err
		}

		if err := config.DefaultAfterInit(c.opts.Context, c); err != nil && !c.opts.AllowFail {
			return err
		}

		return nil
	}

	for _, sf := range fields {
		tf, ok := sf.Field.Tag.Lookup(c.opts.StructTag)
		if !ok {
			continue
		}

		fn, fd, fv := getFlagOpts(tf)
		if tf, ok = sf.Field.Tag.Lookup(c.env); ok {
			fd += fmt.Sprintf(" (env %s)", tf)
		}

		rcheck := true

		if !sf.Value.IsValid() {
			continue
		}
		vi := sf.Value.Interface()
		if vi == nil {
			continue
		}
		if f := flag.Lookup(fn); f != nil {
			return nil
		}

		switch vi.(type) {
		case time.Duration:
			err = c.flagDuration(sf.Value, fn, fv, fd)
			rcheck = false
		case time.Time:
			err = c.flagTime(sf.Value, fn, fv, fd)
			rcheck = false
		}

		if err != nil {
			if !c.opts.AllowFail {
				return err
			}
			if err := config.DefaultAfterInit(c.opts.Context, c); err != nil && !c.opts.AllowFail {
				return err
			}

			return nil
		}

		if !rcheck {
			continue
		}

		if sf.Value.Kind() == reflect.Ptr {
			sf.Value = sf.Value.Elem()
		}

		switch sf.Value.Kind() {
		case reflect.String:
			err = c.flagString(sf.Value, fn, fv, fd)
		case reflect.Bool:
			err = c.flagBool(sf.Value, fn, fv, fd)
		case reflect.Int:
			err = c.flagInt(sf.Value, fn, fv, fd)
		case reflect.Int64:
			err = c.flagInt64(sf.Value, fn, fv, fd)
		case reflect.Float64:
			err = c.flagFloat64(sf.Value, fn, fv, fd)
		case reflect.Uint:
			err = c.flagUint(sf.Value, fn, fv, fd)
		case reflect.Uint64:
			err = c.flagUint64(sf.Value, fn, fv, fd)
		case reflect.Slice:
			err = c.flagSlice(sf.Value, fn, fv, fd)
		case reflect.Map:
			err = c.flagMap(sf.Value, fn, fv, fd)
		}
		if err != nil {
			if !c.opts.AllowFail {
				return err
			}

			if err := config.DefaultAfterInit(c.opts.Context, c); err != nil && !c.opts.AllowFail {
				return err
			}

			return nil
		}
	}

	if err := config.DefaultAfterInit(c.opts.Context, c); err != nil && !c.opts.AllowFail {
		return err
	}

	return nil
}

func (c *flagConfig) Load(ctx context.Context, opts ...config.LoadOption) error {
	if c.opts.SkipLoad != nil && c.opts.SkipLoad(ctx, c) {
		return nil
	}

	options := config.NewLoadOptions(opts...)
	_ = options

	if err := config.DefaultBeforeLoad(ctx, c); err != nil && !c.opts.AllowFail {
		return err
	}

	if err := config.DefaultAfterLoad(ctx, c); err != nil && !c.opts.AllowFail {
		return err
	}

	// TODO: allow merge, append and so
	flag.Parse()

	return nil
}

func (c *flagConfig) Save(ctx context.Context, opts ...config.SaveOption) error {
	if c.opts.SkipSave != nil && c.opts.SkipSave(ctx, c) {
		return nil
	}

	if err := config.DefaultBeforeSave(ctx, c); err != nil && !c.opts.AllowFail {
		return err
	}

	if err := config.DefaultAfterSave(ctx, c); err != nil && !c.opts.AllowFail {
		return err
	}

	return nil
}

func (c *flagConfig) String() string {
	return "flag"
}

func (c *flagConfig) Name() string {
	return c.opts.Name
}

func (c *flagConfig) Watch(ctx context.Context, opts ...config.WatchOption) (config.Watcher, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *flagConfig) usage() {
	mapDelim := DefaultMapDelim
	sliceDelim := DefaultSliceDelim

	if c.opts.Context != nil {
		if d, ok := c.opts.Context.Value(mapDelimKey{}).(string); ok {
			mapDelim = d
		}
		if d, ok := c.opts.Context.Value(sliceDelimKey{}).(string); ok {
			sliceDelim = d
		}
	}

	if c.name == "" {
		fmt.Fprintf(c.fset.Output(), "Usage:\n")
	} else {
		fmt.Fprintf(c.fset.Output(), "Usage of %s:\n", c.name)
	}

	c.fset.VisitAll(func(f *flag.Flag) {
		var b strings.Builder
		fmt.Fprintf(&b, "  -%s", f.Name) // Two spaces before -; see next two comments.
		_, usage := flag.UnquoteUsage(f)
		name := "value"
		v := reflect.TypeOf(f.Value).String()
		b.WriteString(" ")
		switch v {
		case "*flag.boolFlag":
			name = "bool"
		case "*flag.durationValue":
			name = "duration"
		case "*flag.float64Value":
			name = "float"
		case "*flag.intValue", "*flag.int64Value":
			name = "int"
		case "*flag.stringValue":
			name = "string"
		case "*flag.uintValue", "*flag.uint64Value":
			name = "uint"
		case "*flag.mapValue":
			//	nv := f.Value.(*mapValue)
			name = fmt.Sprintf("string key=val with %q as separator", mapDelim)
		case "*flag.sliceValue":
			//	nv := f.Value.(*sliceValue)
			name = fmt.Sprintf("string with %q as separator", sliceDelim)
		}
		b.WriteString(name)

		// Boolean flags of one ASCII letter are so common we
		// treat them specially, putting their usage on the same line.
		if b.Len() <= 4 { // space, space, '-', 'x'.
			b.WriteString("\t")
		} else {
			// Four spaces before the tab triggers good alignment
			// for both 4- and 8-space tab stops.
			b.WriteString("\n    \t")
		}
		b.WriteString(strings.ReplaceAll(usage, "\n", "\n    \t"))

		if f.Value.String() == f.DefValue {
			fmt.Fprintf(&b, " (default %q)", f.DefValue)
		} else {
			fmt.Fprintf(&b, " (default %q current %q)", f.DefValue, f.Value)
		}

		fmt.Fprint(c.fset.Output(), b.String(), "\n")
	})
}

func (c *flagConfig) configure() {
	flagSet := flag.CommandLine
	flagSetName := os.Args[0]
	flagSetErrorHandling := flag.ExitOnError
	flagEnv := "env"
	var flagUsage func()
	var isSet bool

	if c.opts.Context != nil {
		if v, ok := c.opts.Context.Value(flagSetNameKey{}).(string); ok {
			isSet = true
			flagSetName = v
		}
		if v, ok := c.opts.Context.Value(flagSetErrorHandlingKey{}).(flag.ErrorHandling); ok {
			isSet = true
			flagSetErrorHandling = v
		}
		if v, ok := c.opts.Context.Value(flagSetKey{}).(*flag.FlagSet); ok {
			flagSet = v
		}
		if v, ok := c.opts.Context.Value(flagSetUsageKey{}).(func()); ok {
			flagUsage = v
		}
		if v, ok := c.opts.Context.Value(flagEnvKey{}).(string); ok {
			flagEnv = v
		}
	}
	c.fset = flagSet

	if isSet {
		c.fset.Init(flagSetName, flagSetErrorHandling)
	}
	if flagUsage != nil {
		c.fset.Usage = flagUsage
	} else {
		c.fset.Usage = c.usage
	}
	c.env = flagEnv

	c.name = flagSetName
}

func NewConfig(opts ...config.Option) config.Config {
	options := config.NewOptions(opts...)
	if len(options.StructTag) == 0 {
		options.StructTag = DefaultStructTag
	}

	c := &flagConfig{opts: options}
	c.configure()

	return c
}
