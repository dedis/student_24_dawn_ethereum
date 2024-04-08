package main

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"os"

	"golang.org/x/perf/benchfmt"
)

func Main() error {
	f, err := os.Open("cae_benchmark")
	if err != nil {
		return err
	}
	defer f.Close()
	reader := benchfmt.NewReader(f, "cae_benchmark")

	path := "benchmarks/cae/"
	os.MkdirAll(path, 0755)
	writers := map[string]*csv.Writer{}

	for reader.Scan() {
		if err := reader.Err(); err != nil {
			return err
		}
		switch r := reader.Result().(type) {
		case *benchfmt.Result:
			_, parts := r.Name.Parts()
			var scheme, l string
			for _, p := range parts {
				if v, ok := bytes.CutPrefix(p, []byte("/Scheme=")); ok {
					scheme = string(v)
				}
				if v, ok := bytes.CutPrefix(p, []byte("/l=")); ok {
					l = string(v)
				}
			}
			if writers[scheme] == nil {
				f, err := os.Create(path + scheme + ".csv")
				if err != nil {
					return err
				}
				w := csv.NewWriter(f)
				w.Comma = '\t'
				writers[scheme] = w
				w.Write([]string{"l", "sec/B", "sec/op"})
			}
			secPerByte, ok := r.Value("sec/B")
			if !ok {
				return errors.New("missing measurement: sec/B")
			}
			secPerOp, ok := r.Value("sec/op")
			if !ok {
				return errors.New("missing measurement: sec/op")
			}
			writers[scheme].Write([]string{l, fmt.Sprintf("%g", secPerByte), fmt.Sprintf("%g", secPerOp)})
		}
	}

	for _, w := range writers {
		w.Flush()
	}

	return nil
}

func main() {
	if err := Main(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
