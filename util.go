package flag

import (
	"flag"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func convertType(v reflect.Value, t reflect.Kind) (reflect.Value, error) {
	switch v.Kind() {
	case reflect.String:
		switch t {
		case reflect.String:
			return v, nil
		case reflect.Int, reflect.Int64:
			i, err := strconv.ParseInt(v.String(), 10, 64)
			if err != nil {
				return v, err
			}
			return reflect.ValueOf(i), nil
		case reflect.Uint, reflect.Uint64:
			i, err := strconv.ParseUint(v.String(), 10, 64)
			if err != nil {
				return v, err
			}
			return reflect.ValueOf(i), nil
		case reflect.Float64:
			i, err := strconv.ParseFloat(v.String(), 64)
			if err != nil {
				return v, err
			}
			return reflect.ValueOf(i), nil
		case reflect.Bool:
			i, err := strconv.ParseBool(v.String())
			if err != nil {
				return v, err
			}
			return reflect.ValueOf(i), nil
		}
	}
	return v, ErrInvalidValue
}

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
	delim := DefaultMapDelim
	if c.opts.Context != nil {
		if d, ok := c.opts.Context.Value(mapDelimKey{}).(string); ok {
			delim = d
		}
	}
	flag.Func(fn, fv, func(s string) error {
		ps := strings.Split(s, delim)
		if len(ps) == 0 {
			return nil
		}
		v.Set(reflect.MakeMapWithSize(v.Type(), len(ps)))
		kt := v.Type().Key().Kind()
		vt := v.Type().Elem().Kind()

		for i := 0; i < len(ps); i++ {
			fs := strings.Split(ps[i], "=")
			switch len(fs) {
			case 0:
				return nil
			case 1:
				if len(fs[0]) == 0 {
					return nil
				}
				return ErrInvalidValue
			case 2:
				break
			default:
				return ErrInvalidValue
			}
			key, err := convertType(reflect.ValueOf(fs[0]), kt)
			if err != nil {
				return err
			}
			val, err := convertType(reflect.ValueOf(fs[1]), vt)
			if err != nil {
				return err
			}
			v.SetMapIndex(key.Convert(v.Type().Key()), val.Convert(v.Type().Elem()))
		}
		return nil
	})

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
	i, err := time.ParseDuration(fv)
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
	i, err := strconv.ParseBool(fv)
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
	i, err := strconv.ParseInt(fv, 10, 64)
	if err != nil {
		return err
	}
	flag.IntVar(nv, fn, int(i), fv)
	return nil
}

func (c *flagConfig) flagInt64(v reflect.Value, fn, fv, fd string) error {
	nv, ok := v.Addr().Interface().(*int64)
	if !ok {
		return ErrInvalidValue
	}
	i, err := strconv.ParseInt(fv, 10, 64)
	if err != nil {
		return err
	}
	flag.Int64Var(nv, fn, int64(i), fv)
	return nil
}

func (c *flagConfig) flagUint(v reflect.Value, fn, fv, fd string) error {
	nv, ok := v.Addr().Interface().(*uint)
	if !ok {
		return ErrInvalidValue
	}
	i, err := strconv.ParseUint(fv, 10, 64)
	if err != nil {
		return err
	}
	flag.UintVar(nv, fn, uint(i), fv)
	return nil
}

func (c *flagConfig) flagUint64(v reflect.Value, fn, fv, fd string) error {
	nv, ok := v.Addr().Interface().(*uint64)
	if !ok {
		return ErrInvalidValue
	}
	i, err := strconv.ParseUint(fv, 10, 64)
	if err != nil {
		return err
	}
	flag.Uint64Var(nv, fn, uint64(i), fv)
	return nil
}

func (c *flagConfig) flagFloat64(v reflect.Value, fn, fv, fd string) error {
	nv, ok := v.Addr().Interface().(*float64)
	if !ok {
		return ErrInvalidValue
	}
	i, err := strconv.ParseFloat(fv, 64)
	if err != nil {
		return err
	}
	flag.Float64Var(nv, fn, float64(i), fv)
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
		if len(ret[idx]) == 0 {
			continue
		}
		if ret[idx][0] == '\'' {
			ret[idx] = ret[idx][1:]
		}
		if ret[idx][len(ret[idx])-1] == '\'' {
			ret[idx] = ret[idx][:len(ret[idx])-1]
		}
	}
	return ret[0], ret[1], ret[2]
}
