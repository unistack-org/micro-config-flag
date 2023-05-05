package flag

import (
	"context"
	"flag"
	"os"
	"testing"
	"time"

	"go.unistack.org/micro/v3/config"
)

func TestLoad(t *testing.T) {
	time.Local = time.UTC
	os.Args = append(os.Args, "-broker", "5566:33")
	os.Args = append(os.Args, "-verbose")
	os.Args = append(os.Args, "-wait", "5s")
	os.Args = append(os.Args, "-addr", "33,44")
	os.Args = append(os.Args, "-time", time.RFC822)
	os.Args = append(os.Args, "-metadata", "key=20")
	os.Args = append(os.Args, "-components", "all=info,api=debug")
	os.Args = append(os.Args, "-addr", "33,44")
	os.Args = append(os.Args, "-badflag", "test")
	type NestedConfig struct {
		Value string `flag:"name=nested_value"`
	}
	type Config struct {
		Time           time.Time         `flag:"name=time,desc='some time',default='02 Jan 06 15:04 MST'"`
		Components     map[string]string `flag:"name=components,desc='components logging'"`
		Metadata       map[string]int    `flag:"name=metadata,desc='some meta',default=''"`
		Nested         *NestedConfig
		Broker         string        `flag:"name=broker,desc='description with, comma',default='127.0.0.1:9092'"`
		WithoutDesc    string        `flag:"name=without_desc,default='without_default'"`
		WithoutAll     string        `flag:"name=without_all"`
		WithoutDefault string        `flag:"name=without_default,desc='with'"`
		Addr           []string      `flag:"name=addr,desc='addrs',default='127.0.0.1:9092'"`
		Wait           time.Duration `flag:"name=wait,desc='wait time',default='2s'"`
		Verbose        bool          `flag:"name=verbose,desc='verbose output',default='false'"`
	}

	ctx := context.Background()
	cfg := &Config{Nested: &NestedConfig{}}

	c := NewConfig(config.Struct(cfg), TimeFormat(time.RFC822), FlagErrorHandling(flag.ContinueOnError))
	if err := c.Init(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	// double init test
	if err := c.Init(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	if err := c.Load(ctx); err != nil {
		t.Fatalf("load failed: %v", err)
	}

	if cfg.Broker != "5566:33" {
		t.Fatalf("failed to parse flags broker value invalid: %#+v", cfg)
	}
	if tf := cfg.Time.Format(time.RFC822); tf != "02 Jan 06 15:04 MST" {
		t.Fatalf("parse time error: %s != %s", tf, "02 Jan 06 15:04 MST")
	}

	if len(cfg.Components) != 2 {
		t.Fatalf("cant parse map components %#+v", cfg)
	}
}
