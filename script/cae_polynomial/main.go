package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
)

type point struct{ l, gas float64 }

func readPoints(filename string) ([]point, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	reader := csv.NewReader(f)
	reader.Comma = '\t'

	var points []point
	headers, err := reader.Read()
	if err != nil {
		return nil, err
	}

	lIndex := slices.Index(headers, "l")
	if lIndex < 0 {
		return nil, fmt.Errorf("missing l column")
	}

	gasIndex := slices.Index(headers, "gas/op")
	if gasIndex < 0 {
		return nil, fmt.Errorf("missing gas/op column")
	}

	for {
		row, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		l, err := strconv.ParseFloat(row[lIndex], 64)
		if err != nil {
			return nil, err
		}
		gas, err := strconv.ParseFloat(row[gasIndex], 64)
		if err != nil {
			return nil, err
		}
		points = append(points, point{l, gas})
	}

	return points, nil
}

func evalAt(points []point, x float64) (result float64) {
	// Lagrange interpolation
	for i := range points {
		l := 1.0
		for j := range points {
			if i != j {
				l *= (x - points[j].l) / (points[i].l - points[j].l)
			}
		}
		result += points[i].gas * l
	}
	return
}

func Main() error {
	for _, filename := range os.Args[1:] {
		points, err := readPoints(filename)
		if err != nil {
			return err
		}
		a := evalAt(points, 0)
		b := evalAt(points, 1) - a
		// f(2) = 4*c + 2*b + a
		c := (evalAt(points, 2) - a - 2*b)
		if c > 1e-3 {
			return fmt.Errorf("expected affine polynomial!")
		}

		fmt.Printf("%s: f(x) = %f*x + %f\n", filename, b, a)
	}

	return nil
}

func main() {
	if err := Main(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
