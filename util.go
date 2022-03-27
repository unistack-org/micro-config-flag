package flag

import (
	"flag"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type mapValue struct {
	delim string
	def   string
	v     reflect.Value
}

func (v mapValue) String() string {
	if v.v.Kind() != reflect.Invalid {
		var kv []string
		it := v.v.MapRange()
		for it.Next() {
			k := it.Key().Interface()
			v := it.Value().Interface()
			kv = append(kv, fmt.Sprintf("%v=%v", k, v))
		}
		return strings.Join(kv, ",")
	}
	return v.def
}

func (v mapValue) Set(s string) error {
	ps := strings.Split(s, v.delim)
	if len(ps) == 0 {
		return nil
	}
	v.v.Set(reflect.MakeMapWithSize(v.v.Type(), len(ps)))
	kt := v.v.Type().Key().Kind()
	vt := v.v.Type().Elem().Kind()

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
		v.v.SetMapIndex(key.Convert(v.v.Type().Key()), val.Convert(v.v.Type().Elem()))
	}
	return nil
}

type sliceValue struct {
	delim string
	def   string
	v     reflect.Value
}

func (v sliceValue) String() string {
	if v.v.Kind() != reflect.Invalid {
		var kv []string
		for idx := 0; idx < v.v.Len(); idx++ {
			kv = append(kv, fmt.Sprintf("%v", v.v.Index(idx).Interface()))
		}
		return strings.Join(kv, ",")
	}
	return v.def
}

func (v sliceValue) Set(s string) error {
	p := strings.Split(s, v.delim)
	v.v.Set(reflect.MakeSlice(v.v.Type(), len(p), len(p)))
	switch v.v.Type().Elem().Kind() {
	case reflect.Int, reflect.Int64:
		for idx := range p {
			i, err := strconv.ParseInt(p[idx], 10, 64)
			if err != nil {
				return err
			}
			v.v.Index(idx).SetInt(i)
		}
	case reflect.Uint, reflect.Uint64:
		for idx := range p {
			i, err := strconv.ParseUint(p[idx], 10, 64)
			if err != nil {
				return err
			}
			v.v.Index(idx).SetUint(i)
		}
	case reflect.Float64:
		for idx := range p {
			i, err := strconv.ParseFloat(p[idx], 64)
			if err != nil {
				return err
			}
			v.v.Index(idx).SetFloat(i)
		}
	case reflect.Bool:
		for idx := range p {
			i, err := strconv.ParseBool(p[idx])
			if err != nil {
				return err
			}
			v.v.Index(idx).SetBool(i)
		}
	case reflect.String:
		for idx := range p {
			v.v.Index(idx).SetString(p[idx])
		}
	}
	return nil
}

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

	v.Set(reflect.MakeSlice(v.Type(), 0, 0))
	mp := &sliceValue{v: v, def: fv, delim: delim}
	if err := mp.Set(fv); err != nil {
		return err
	}
	flag.Var(mp, fn, fd)
	return nil
}

func (c *flagConfig) flagMap(v reflect.Value, fn, fv, fd string) error {
	delim := DefaultMapDelim
	if c.opts.Context != nil {
		if d, ok := c.opts.Context.Value(mapDelimKey{}).(string); ok {
			delim = d
		}
	}
	v.Set(reflect.MakeMapWithSize(v.Type(), 0))
	mp := &mapValue{v: v, def: fv, delim: delim}
	if err := mp.Set(fv); err != nil {
		return err
	}
	flag.Var(mp, fn, fd)
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
	var name, desc, def string
	delim := ","

	var buf string
	for idx := 0; idx < len(tf); idx++ {
		buf += string(tf[idx])
		switch buf {
		case "name":
			ndx := idx + 2
			stop := ","
			var quote bool
			for ; ndx < len(tf); ndx++ {
				if string(tf[ndx]) == stop {
					if quote {
						ndx++
					}
					break
				}
				if string(tf[ndx]) == "'" {
					stop = "'"
					quote = true
					continue
				}
				name += string(tf[ndx])
			}
			idx = ndx
			buf = ""
		case "desc":
			ndx := idx + 2
			stop := ","
			var quote bool
			for ; ndx < len(tf); ndx++ {
				if string(tf[ndx]) == stop {
					if quote {
						ndx++
					}
					break
				}
				if string(tf[ndx]) == "'" {
					stop = "'"
					quote = true
					continue
				}
				desc += string(tf[ndx])
			}
			idx = ndx
			buf = ""
		case "default":
			ndx := idx + 2
			stop := ","
			var quote bool
			for ; ndx < len(tf); ndx++ {
				if string(tf[ndx]) == stop && (stop != delim) {
					if quote {
						ndx++
					}
					break
				}
				if string(tf[ndx]) == "'" {
					stop = "'"
					quote = true
					continue
				}
				def += string(tf[ndx])
			}
			idx = ndx
			buf = ""
		}
	}
	return name, desc, def
}
