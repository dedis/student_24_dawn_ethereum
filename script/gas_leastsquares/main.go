package main

import (
	"fmt"
	"math"
	"math/bits"
	"strings"
	"os"

	"gonum.org/v1/gonum/mat"
)

type features uint
const (
	base features = 1 << iota
	commit
	reveal
	reveal_highest
	vdf
	smc
	bid
	bid_highest
	n_features int = iota
)

func (f features) Index() int {
	return bits.Len(uint(f)) - 1
}

type point struct{ features features; gas uint64 }

func Main() error {
	points := []point{
		// by definition
		{base, 21_000},
		{smc, 192_000},

		// measured
		{base | commit, 67_200},
		{base | reveal, 67_243},
		{base | reveal_highest, 83_614},
		{base | smc | bid, 224_487},
		{base | smc | bid_highest, 263_537},
		{base | vdf | bid, 98_487},
		{base | vdf | bid_highest, 137_537},
	}
	a := mat.NewDense(len(points), n_features, nil)
	x := mat.NewDense(n_features, 1, nil)
	b := mat.NewDense(len(points), 1, nil)
	for i, p := range points {
		for j := range n_features {
			a.Set(i, j, float64((p.features >> j) & 1))
		}
		b.Set(i, 0, float64(p.gas))
	}
	err := x.Solve(a, b)
	if err != nil {
		return err
	}
	components := make([]int, n_features)
	for i := range components {
		components[i] = int(math.Round(x.At(i, 0)))
	}
	b0 := mat.NewDense(len(points), 1, nil)
	b0.Mul(a, x)
	b0.Sub(b0, b)
	for i := range points {
		if int(b0.At(i, 0)) != 0 {
			return fmt.Errorf("residual error: %d %v", i, b0.At(i, 0))
		}
	}


	feature_variable_names := map[features]string{}
	for _, f := range []struct{desc string; feature features}{
		{"base", base},
		{"smc", smc},
		{"vdf", vdf},
	} {
		feature_variable_names[f.feature] = fmt.Sprintf(`\text{%s}`, f.desc)
	}
	fmt.Printf(`\begin{tabular}{llr}`)
	fmt.Println()
	for _, row := range []struct{desc string; features features}{
		{"Commit bid", base | commit},
		{"Reveal bid", base | reveal},
		{"Reveal highest bid", base | reveal_highest},
		{"bid (SMC)", base | smc | bid},
		{"bid highest (SMC)", base | smc | bid_highest},
		{"bid (VDF)", base | vdf | bid},
		{"bid highest (VDF)", base | vdf | bid_highest},
	} {
		terms := []string{}
		sum := 0
		for i := range n_features {
			if row.features & (1 << i) != 0 {
				sum += components[i]
				if name, ok := feature_variable_names[1 << i]; ok {
					terms = append(terms, name)
				} else {
					terms = append(terms, fmt.Sprint(components[i]))
				}
			}
		}
		fmt.Printf(`%s & $%s$ & $%d$ \\`, row.desc, strings.Join(terms, " + "), sum)
		fmt.Println()
	}
	fmt.Printf(`\end{tabular}`)
	fmt.Println()

	fmt.Printf(`\text{base} = %d & \text{smc} = %d & \text{vdf} = %d \\`, components[base.Index()], components[smc.Index()], components[vdf.Index()])
	fmt.Println()

	return nil
}

func main() {
	if err := Main(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
