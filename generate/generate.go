package generate

import (
	"errors"
	"io"
	"math/rand"
	"os"
	"strings"
	"sync"

	"github.com/tymbaca/dna/model"
)

const _chunkSize = 256 * 1024

func GenerateDNAToFile(name string) {
	// w := bytes.NewBuffer(nil)
	w, err := os.Create(name)

	if err != nil {
		panic(err)
	}

	// err = generateAndWriteDNA(w, 130_000)
	// err = generateAndWriteDNA(w, 130_000_000) // approx 1 hromosome
	// err = generateAndWriteDNA(w, 3_000_000_000) // 32 hromosome

	// if err := generateAndWriteDNA(w, 100_000_000); err != nil {
	// 	panic(err)
	// }

	wg := sync.WaitGroup{}

	wg.Add(4)

	go func() {
		defer wg.Done()
		if err := GenerateAndWriteDNA(w, 25_000_000); err != nil {
			panic(err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := GenerateAndWriteDNA(w, 25_000_000); err != nil {
			panic(err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := GenerateAndWriteDNA(w, 25_000_000); err != nil {
			panic(err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := GenerateAndWriteDNA(w, 25_000_000); err != nil {
			panic(err)
		}
	}()

	wg.Wait()

	// fmt.Printf("generated DNA: %s\n", w.String())
}

func GenerateAndWriteDNA(w io.Writer, pairsNum int) error {
	var err error

	if pairsNum <= 0 {
		return errors.New("pair count need to be positive number")
	}

	// for range pairsNum {
	// 	// log.Println("wrote")
	// 	_, err = w.Write([]byte(randomPair()))
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	pairsLeft := pairsNum
	for pairsLeft > 0 {
		// form a chunk
		b := strings.Builder{}
		b.Grow(_chunkSize)

		// calculate amount of pairs
		pairsToGen := _chunkSize * 4 // 1 byte can hold 4 pair
		pairsLeft -= pairsToGen
		if pairsLeft < 0 {
			pairsToGen += pairsLeft // if pairNum is -48 -> pairCound -= 48
		}

		// form and fill a chunk
		for range pairsToGen {
			_, err = b.WriteString(RandomPair())
			if err != nil {
				return err
			}
		}

		// write to target
		_, err := w.Write([]byte(b.String()))
		if err != nil {
			return err
		}

		b.Reset()
		// log.Printf("successfully wrote %d pairs to target\n", pairsToGen)
	}

	return nil
}

func RandomPair() model.Pair {
	return model.Pairs[rand.Intn(len(model.Pairs))]
}
