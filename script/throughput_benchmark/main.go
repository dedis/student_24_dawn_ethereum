package main

import (
	"testing"
	"sync"
	"fmt"

	"github.com/ethereum/go-ethereum/f3b"
)

func benchmark(b *testing.B, batchSize int) {
	p := f3b.SelectedProtocol()
	if p == nil {
		b.Skip("no protocol selected")
	}
	label := []byte{1,2,3,4,5,6,7,8}
	inputs := make([][]byte, batchSize)
	for i := range inputs {
		var err error
		_, inputs[i], err = p.ShareSecret(label)
		if err != nil {
			b.Fail()
		}
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
	var wg sync.WaitGroup
	wg.Add(len(inputs))
	for i := range inputs {
		go func() {
			defer wg.Done()
			_, err := p.RevealSecret(label, inputs[i])
			if err != nil {
				b.Fail()
			}
		}()
	}
	wg.Wait()
}
}

func main() {
	//testing.Init()
	for i := 0; i < 12; i++ {
		batchSize := 1 << i
		r := testing.Benchmark(func(b *testing.B) {
			benchmark(b, batchSize)
		})
		fmt.Printf("batchSize=%d\t%s\n", batchSize, r)
	}
}
