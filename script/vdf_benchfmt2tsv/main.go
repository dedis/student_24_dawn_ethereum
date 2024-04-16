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
	f, err := os.Open("vdf_benchmark")
	if err != nil {
		return err
	}
	defer f.Close()
	reader := benchfmt.NewReader(f, "vdf_benchmark")

	path := "benchmarks/vdf/"
	os.MkdirAll(path, 0755)
	writers := map[string]*csv.Writer{}

	for reader.Scan() {
		if err := reader.Err(); err != nil {
			return err
		}
		switch r := reader.Result().(type) {
		case *benchfmt.Result:
			name, parts := r.Name.Parts()
			algorithm := string(name)
			var log2t string
			for _, p := range parts {
				if v, ok := bytes.CutPrefix(p, []byte("/log2t=")); ok {
					log2t = string(v)
				}
			}
			if writers[algorithm] == nil {
				f, err := os.Create(path + algorithm + ".csv")
				if err != nil {
					return err
				}
				w := csv.NewWriter(f)
				w.Comma = '\t'
				writers[algorithm] = w
				w.Write([]string{"log2t", "sec/op"})
			}
			secPerOp, ok := r.Value("sec/op")
			if !ok {
				return errors.New("missing measurement: sec/op")
			}
			writers[algorithm].Write([]string{log2t, fmt.Sprintf("%g", secPerOp)})
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
