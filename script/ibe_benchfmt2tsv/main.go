package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"

	"golang.org/x/perf/benchfmt"
)

func Main() error {
	f, err := os.Open("ibe_benchmark")
	if err != nil {
		return err
	}
	defer f.Close()
	reader := benchfmt.NewReader(f, "ibe_benchmark")

	f, err = os.Create("benchmarks/ibe.csv")
	if err != nil {
		return err
	}
	w := csv.NewWriter(f)
	w.Comma = '\t'
	defer w.Flush()

	for reader.Scan() {
		if err := reader.Err(); err != nil {
			return err
		}
		switch r := reader.Result().(type) {
		case *benchfmt.Result:
			secPerOp, ok := r.Value("sec/op")
			if !ok {
				return errors.New("missing measurement: sec/op")
			}
			w.Write([]string{string(r.Name), fmt.Sprintf("%g", secPerOp)})
		}
	}

	return nil
}

func main() {
	if err := Main(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
