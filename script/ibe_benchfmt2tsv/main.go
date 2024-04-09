package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"

	"golang.org/x/perf/benchfmt"
	"gonum.org/v1/gonum/mat"
)

func Main() error {
	f, err := os.Open("ibe_benchmark")
	if err != nil {
		return err
	}
	defer f.Close()
	reader := benchfmt.NewReader(f, "ibe_benchmark")
	secPerOpTable := map[string]float64{}

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
			secPerOpTable[string(r.Name)] = secPerOp
		}
	}
	coeffs := map[string][]float64{
		// name: {miller, exp}
		"VerifyIdentity/Variant=Slow-4": {2, 2},
		"VerifyIdentity/Variant=Fast-4": {2, 1},
		"RecoverSecret-4":               {1, 1},
	}

	a := mat.NewDense(len(coeffs), 2, nil)
	b := mat.NewDense(len(coeffs), 1, nil)
	x := mat.NewDense(2, 1, nil)
	i := 0
	for name := range coeffs {
		a.SetRow(i, coeffs[name])
		b.Set(i, 0, secPerOpTable[name])
		i++
	}
	err = x.Solve(a, b)
	if err != nil {
		return err
	}
	millerTime := x.At(0, 0)
	expTime := x.At(1, 0)

	f, err = os.Create("benchmarks/ibe.csv")
	if err != nil {
		return err
	}
	w := csv.NewWriter(f)
	w.Comma = '\t'
	w.Write([]string{"name", "millerPart", "expPart", "remainder"})
	outNames := map[string]string{
		"VerifyIdentity/Variant=Slow-4": "VerifyIdentity[Slow]",
		"VerifyIdentity/Variant=Fast-4": "VerifyIdentity[Fast]",
		"RecoverSecret-4":               "RecoverSecret",
	}
	for name := range coeffs {
		millerPart := coeffs[name][0] * millerTime
		expPart := coeffs[name][1] * expTime
		remainder := secPerOpTable[name] - millerPart - expPart
		w.Write([]string{outNames[name], fmt.Sprintf("%g", millerPart), fmt.Sprintf("%g", expPart), fmt.Sprintf("%g", remainder)})
	}
	w.Flush()
	return nil
}

func main() {
	if err := Main(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
