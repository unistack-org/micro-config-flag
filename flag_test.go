package flag

import (
	"flag"
	"os"
	"reflect"
	"strconv"
	"testing"

	rutil "github.com/unistack-org/micro/v3/util/reflect"
)

func TestLoad(t *testing.T) {
	os.Args = append(os.Args, "-broker", "5566:33")
	type config struct {
		Broker  string `flag:"name=broker,desc='description with, comma',default='127.0.0.1:9092'"`
		Verbose bool   `flag:"name=verbose,desc='verbose output',default='false value'"`
	}

	cfg := &config{}

	fields, err := rutil.StructFields(cfg)
	if err != nil {
		t.Fatal(err)
	}

	for _, sf := range fields {
		tf, ok := sf.Field.Tag.Lookup("flag")
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

	flag.Parse()
	if cfg.Broker != "5566:33" {
		t.Fatalf("failed to parse flags broker value invalid: %#+v", cfg)
	}
}
