package encoder

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/tymbaca/dna/model"
)

// WARN: must be dividable to 8
// ATGCTACG - 8 bytes
// 10011100 - 1 byte
const _chunkSize = 256 * 1024

type DNAEncoder struct {
	src io.Reader
	// mu  sync.Mutex
}

func NewDNAEncoder(buf io.Reader) *DNAEncoder {
	return &DNAEncoder{src: buf}
}

// WriteTo encodes DNA to w until there's no more DNA to encode or
// when an error occurs. The return value n is the number of encoded bytes
// written. Any error encountered during the write is also returned.
func (e *DNAEncoder) WriteTo(target io.Writer) (n int64, err error) {
	for {
		localN, err := e.writeChunk(target)
		n += int64(localN)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return 0, err
		}
	}

	return n, nil
}
func (e *DNAEncoder) writeChunk(target io.Writer) (int, error) {
	buf := make([]byte, _chunkSize)
	eof := false

	// read
	_, err := e.src.Read(buf)
	if err != nil {
		if !errors.Is(err, io.EOF) {
			return 0, err
		}

		eof = true
	}

	// validate
	err = validateChunk(buf)
	if err != nil {
		return 0, err
	}

	// convert
	resultBuf, _, err := convertChunk(buf)
	if err != nil {
		return 0, err
	}
	if len(resultBuf) < len(buf)/8 { // TODO: move to convertChunk and return EOF
		eof = true
	}

	// write
	n, err := target.Write(resultBuf)
	if err != nil {
		return 0, err
	}

	if eof {
		return n, io.EOF
	}

	return n, nil
}

func validateChunk(buf []byte) error {
	if len(buf)%2 != 0 {
		return errors.New("chunk len is not even")
	}

	if len(buf)%8 != 0 {
		return fmt.Errorf("pairs chunk must be dividable to 8, got len: %d", len(buf))
	}

	return nil
}

func convertChunk(buf []byte) ([]byte, int, error) {
	// preallocate result buf with empty bytes
	// every
	resultBuf := make([]byte, len(buf)/8)
	// pairs encoded
	n := 0

	// pair pointer - points to current pair position in the byte
	// 00_00_00_00 - byte
	// ^^ ^^ ^^ ^^
	//  0  1  2  3
	// used to correctly offset when &-ing the binary pair
	pp := 0
	// result byte pointer - points to current byte in resultBuf
	rbi := 0

	for i := 0; i < len(buf); i += 2 {
		pair := model.Pair(buf[i : i+2])
		// check if got zero byte
		if strings.Contains(pair, "\x00") {
			return resultBuf[:rbi], n, nil
		}

		binPair, ok := encodePair(pair)
		if !ok {
			return nil, 0, fmt.Errorf("got incorrect base pair: '%s'", pair)
		}

		// fill the result byte with converted bits
		// at first we take current byte, than we calculate the bit offset,
		// and than we write our 2 result bits from binPair to that place
		resultBuf[rbi] = resultBuf[rbi] | binPair<<(2*(3-pp)) // WARN: pp must be < 4, otherwise we will get negative shift panic

		// increment pair pointer and num after every converted pair
		pp++
		n++
		// if pair pointer got to 4 - we need to zero it and walk to next result byte
		if pp >= 4 {
			pp = 0
			rbi++
		}
	}

	return resultBuf[:rbi], n, nil
}

// fills only last 2 bits of result byte
func encodePair(pair model.Pair) (byte, bool) {
	switch strings.ToUpper(pair) {
	case model.CG:
		return 0b00, true
	case model.GC:
		return 0b01, true
	case model.AT:
		return 0b10, true
	case model.TA:
		return 0b11, true
	default:
		return 0, false
	}
}
