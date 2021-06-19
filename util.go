package flag

import (
	"flag"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func (c *flagConfig) flagSlice(v reflect.Value, fn, fv, fd string) error {
	delim := DefaultSliceDelim
	if c.opts.Context != nil {
		if d, ok := c.opts.Context.Value(sliceDelimKey{}).(string); ok {
			delim = d
		}
	}

	flag.Func(fn, fd, func(s string) error {
		p := strings.Split(s, delim)
		v.Set(reflect.MakeSlice(v.Type(), len(p), len(p)))
		switch v.Type().Elem().Kind() {
		case reflect.Int, reflect.Int64:
			for idx := range p {
				i, err := strconv.ParseInt(p[idx], 10, 64)
				if err != nil {
					return err
				}
				v.Index(idx).SetInt(i)
			}
		case reflect.Uint, reflect.Uint64:
			for idx := range p {
				i, err := strconv.ParseUint(p[idx], 10, 64)
				if err != nil {
					return err
				}
				v.Index(idx).SetUint(i)
			}
		case reflect.Float64:
			for idx := range p {
				i, err := strconv.ParseFloat(p[idx], 64)
				if err != nil {
					return err
				}
				v.Index(idx).SetFloat(i)
			}
		case reflect.Bool:
			for idx := range p {
				i, err := strconv.ParseBool(p[idx])
				if err != nil {
					return err
				}
				v.Index(idx).SetBool(i)
			}
		case reflect.String:
			for idx := range p {
				v.Index(idx).SetString(p[idx])
			}
		}
		return nil
	})

	return nil
}

func (c *flagConfig) flagMap(v reflect.Value, fn, fv, fd string) error {
	return nil
}

func (c *flagConfig) flagTime(v reflect.Value, fn, fv, fd string) error {
	var format string
	if c.opts.Context != nil {
		if tf, ok := c.opts.Context.Value(timeFormatKey{}).(string); ok {
			format = tf
		}
	}
	if format == "" {
		return ErrInvalidValue
	}
	flag.Func(fn, fd, func(s string) error {
		t, err := time.Parse(s, format)
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(t))
		return nil
	})

	return nil
}

func (c *flagConfig) flagDuration(v reflect.Value, fn, fv, fd string) error {
	nv, ok := v.Addr().Interface().(*time.Duration)
	if !ok {
		return ErrInvalidValue
	}
	i, err := time.ParseDuration(fd)
	if err != nil {
		return err
	}
	flag.DurationVar(nv, fn, i, fd)
	return nil
}

func (c *flagConfig) flagBool(v reflect.Value, fn, fv, fd string) error {
	nv, ok := v.Addr().Interface().(*bool)
	if !ok {
		return ErrInvalidValue
	}
	i, err := strconv.ParseBool(fd)
	if err != nil {
		return err
	}
	flag.BoolVar(nv, fn, i, fd)
	return nil
}

func (c *flagConfig) flagString(v reflect.Value, fn, fv, fd string) error {
	nv, ok := v.Addr().Interface().(*string)
	if !ok {
		return ErrInvalidValue
	}
	flag.StringVar(nv, fn, fv, fd)
	return nil
}

func (c *flagConfig) flagInt(v reflect.Value, fn, fv, fd string) error {
	nv, ok := v.Addr().Interface().(*int)
	if !ok {
		return ErrInvalidValue
	}
	i, err := strconv.ParseInt(fd, 10, 64)
	if err != nil {
		return err
	}
	flag.IntVar(nv, fn, int(i), fd)
	return nil
}

func (c *flagConfig) flagInt64(v reflect.Value, fn, fv, fd string) error {
	nv, ok := v.Addr().Interface().(*int64)
	if !ok {
		return ErrInvalidValue
	}
	i, err := strconv.ParseInt(fd, 10, 64)
	if err != nil {
		return err
	}
	flag.Int64Var(nv, fn, int64(i), fd)
	return nil
}

func (c *flagConfig) flagUint(v reflect.Value, fn, fv, fd string) error {
	nv, ok := v.Addr().Interface().(*uint)
	if !ok {
		return ErrInvalidValue
	}
	i, err := strconv.ParseUint(fd, 10, 64)
	if err != nil {
		return err
	}
	flag.UintVar(nv, fn, uint(i), fd)
	return nil
}

func (c *flagConfig) flagUint64(v reflect.Value, fn, fv, fd string) error {
	nv, ok := v.Addr().Interface().(*uint64)
	if !ok {
		return ErrInvalidValue
	}
	i, err := strconv.ParseUint(fd, 10, 64)
	if err != nil {
		return err
	}
	flag.Uint64Var(nv, fn, uint64(i), fd)
	return nil
}

func (c *flagConfig) flagFloat64(v reflect.Value, fn, fv, fd string) error {
	nv, ok := v.Addr().Interface().(*float64)
	if !ok {
		return ErrInvalidValue
	}
	i, err := strconv.ParseFloat(fd, 64)
	if err != nil {
		return err
	}
	flag.Float64Var(nv, fn, float64(i), fd)
	return nil
}

func (c *flagConfig) flagStringSlice(v reflect.Value, fn, fv, fd string) error {
	nv, ok := v.Addr().Interface().(*string)
	if !ok {
		return ErrInvalidValue
	}
	flag.StringVar(nv, fn, fv, fd)
	return nil
}

func getFlagOpts(tf string) (string, string, string) {
	ret := make([]string, 3)
	vals := strings.Split(tf, ",")
	f := 0
	for _, val := range vals {
		p := strings.Split(val, "=")
		switch p[0] {
		case "name":
			f = 0
		case "desc":
			f = 1
		case "default":
			f = 2
		default:
			ret[f] += "," + val
			continue
		}
		ret[f] = p[1]
	}
	for idx := range ret {
		if ret[idx][0] == '\'' {
			ret[idx] = ret[idx][1:]
		}
		if ret[idx][len(ret[idx])-1] == '\'' {
			ret[idx] = ret[idx][:len(ret[idx])-1]
		}
	}
	return ret[0], ret[1], ret[2]
}
