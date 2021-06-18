package flag

import (
	"context"
	"errors"
	"flag"
	"reflect"
	"strconv"

	"github.com/unistack-org/micro/v3/config"
	rutil "github.com/unistack-org/micro/v3/util/reflect"
)

var (
	DefaultStructTag = "flag"
	ErrInvalidStruct = errors.New("invalid struct specified")
)

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
		fn, fv, fd := getFlagOpts(tf)

		switch sf.Value.Kind() {
		case reflect.String:
			v := sf.Value.Addr().Interface().(*string)
			flag.StringVar(v, fn, fv, fd)
		case reflect.Bool:
			v := sf.Value.Addr().Interface().(*bool)
			i, _ := strconv.ParseBool(fv)
			flag.BoolVar(v, fn, i, fd)
		}
	}

	return nil
}

func (c *flagConfig) Load(ctx context.Context) error {
	flag.Parse()
	return nil
}

func (c *flagConfig) Save(ctx context.Context) error {
	return nil
}

func (c *flagConfig) String() string {
	return "flag"
}

func (c *flagConfig) Name() string {
	return c.opts.Name
}

func NewConfig(opts ...config.Option) config.Config {
	options := config.NewOptions(opts...)
	if len(options.StructTag) == 0 {
		options.StructTag = DefaultStructTag
	}
	return &flagConfig{opts: options}
}
