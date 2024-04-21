package model

type Base rune

const (
	// Пурины

	// Guanine
	G Base = 'G'
	// Adenine
	A Base = 'A'

	// Пиримидины

	// Cytosine
	C Base = 'C'
	// Thymine
	T Base = 'T'
)

type Pair = string

const (
	CG Pair = "CG"
	GC Pair = "GC"
	AT Pair = "AT"
	TA Pair = "TA"
)

var Pairs = []Pair{CG, GC, AT, TA}
