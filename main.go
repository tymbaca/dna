package main

import (
	"fmt"
	"runtime"

	"github.com/tymbaca/dna/generate"
)

func main() {
	fmt.Println(runtime.NumCPU())
	dur := MeasureDuration(func() {
		generate.GenerateDNAToFile("output/output.dna")
	})
	fmt.Printf("took %f seconds\n", dur.Seconds())
}
