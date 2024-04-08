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

const GasPerMs = 30

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

	tIndex := slices.Index(headers, "sec/op")
	if tIndex < 0 {
		return nil, fmt.Errorf("missing ms/op column")
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
		t, err := strconv.ParseFloat(row[tIndex], 64)
		if err != nil {
			return nil, err
		}
		gas := GasPerMs * 1e6 * t
		points = append(points, point{l, gas})
	}

	return points, nil
}

func evalAt(points []point, x float64) (result float64) {
	// Lagrange interpolation
	for i := range points {
		r := 1.0
		for j := range points {
			if i != j {
				r *= (x - points[j].l) / (points[i].l - points[j].l)
			}
		}
		result += points[i].gas * r
	}
	return
}

func Main() error {
	rows := []struct{label, filename string}{
		{"ChaCha20-HMAC-SHA256", "benchmarks/cae/chacha20-hmac-sha256.csv"},
		{"AES-CTR-HMAC-SHA256", "benchmarks/cae/aes256ctr-hmac-sha256.csv"},
		{"RK-ChaCha20-Poly1305", "benchmarks/cae/rk-chacha20-poly1305.csv"},
		{"RK-AES-GCM", "benchmarks/cae/rk-aes256-gcm.csv"},
		{"ChaCha20-Poly1305", "benchmarks/cae/chacha20-poly1305.csv"},
		{"AES-GCM", "benchmarks/cae/aes256-gcm.csv"},
		{"ChaCha20", "benchmarks/cae/chacha20.csv"},
		{"ChaCha12", "benchmarks/cae/chacha12.csv"},
		{"ChaCha8", "benchmarks/cae/chacha8.csv"},
	}
	for _, row := range rows {
		points, err :=  readPoints(row.filename)
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

		fmt.Printf("%s & $f(x) = %.3fx + %.3f$\n", row.label, b, a)
		fmt.Println("\\\\") // LaTeX row delimiter
	}
	return nil
}

func main() {
	if err := Main(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
