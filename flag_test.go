package flag

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/unistack-org/micro/v3/config"
)

func TestLoad(t *testing.T) {
	os.Args = append(os.Args, "-broker", "5566:33")
	os.Args = append(os.Args, "-verbose")
	os.Args = append(os.Args, "-wait", "5s")
	os.Args = append(os.Args, "-addr", "33,44")
	os.Args = append(os.Args, "-time", time.RFC822)
	type Config struct {
		Broker  string        `flag:"name=broker,desc='description with, comma',default='127.0.0.1:9092'"`
		Verbose bool          `flag:"name=verbose,desc='verbose output',default='false'"`
		Addr    []string      `flag:"name=addr,desc='addrs',default='127.0.0.1:9092'"`
		Wait    time.Duration `flag:"name=wait,desc='wait time',default='2s'"`
		Time    time.Time     `flag:"name=time,desc='some time',default='02 Jan 06 15:04 MST'"`
	}

	ctx := context.Background()
	cfg := &Config{}

	c := NewConfig(config.Struct(cfg), TimeFormat(time.RFC822))
	if err := c.Init(); err != nil {
		t.Fatal(err)
	}

	if err := c.Load(ctx); err != nil {
		t.Fatal(err)
	}

	if cfg.Broker != "5566:33" {
		t.Fatalf("failed to parse flags broker value invalid: %#+v", cfg)
	}
	if tf := cfg.Time.Format(time.RFC822); tf != "02 Jan 06 14:32 MSK" {
		t.Fatalf("parse time error: %v", cfg.Time)
	}

	t.Logf("cfg %#+v", cfg)
}
