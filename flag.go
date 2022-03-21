package flag // import "go.unistack.org/micro-config-flag/v3"

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"reflect"
	"time"

	"go.unistack.org/micro/v3/config"
	rutil "go.unistack.org/micro/v3/util/reflect"
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
	opts config.Options
}

func (c *flagConfig) Options() config.Options {
	return c.opts
}

func (c *flagConfig) Init(opts ...config.Option) error {
	for _, o := range opts {
		o(&c.opts)
	}

	fields, err := rutil.StructFields(c.opts.Struct)
	if err != nil {
		return err
	}

	for _, sf := range fields {
		tf, ok := sf.Field.Tag.Lookup(c.opts.StructTag)
		if !ok {
			continue
		}
		fn, fd, fv := getFlagOpts(tf)

		rcheck := true

		switch sf.Value.Interface().(type) {
		case time.Duration:
			err = c.flagDuration(sf.Value, fn, fv, fd)
			rcheck = false
		case time.Time:
			err = c.flagTime(sf.Value, fn, fv, fd)
			rcheck = false
		}

		if err != nil {
			c.opts.Logger.Errorf(c.opts.Context, "flag init error: %v", err)
			if !c.opts.AllowFail {
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
			c.opts.Logger.Errorf(c.opts.Context, "flag init error: %v", err)
			if !c.opts.AllowFail {
				return err
			}
			return nil
		}
	}

	return nil
}

func (c *flagConfig) Load(ctx context.Context, opts ...config.LoadOption) error {
	options := config.NewLoadOptions(opts...)
	_ = options
	// TODO: allow merge, append and so
	flag.Parse()
	return nil
}

func (c *flagConfig) Save(ctx context.Context, opts ...config.SaveOption) error {
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

func NewConfig(opts ...config.Option) config.Config {
	options := config.NewOptions(opts...)
	if len(options.StructTag) == 0 {
		options.StructTag = DefaultStructTag
	}
	return &flagConfig{opts: options}
}
