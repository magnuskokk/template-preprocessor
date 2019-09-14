package drum

const (
	Note1 int = 1 << iota
	Note2
	Note4
	Note8
	Note16
	Note32
	Note64
)

type Note struct {
}

type TimeSignature struct {
	Num  uint64
	Note Note
}

type Beat struct {
	Signature TimeSignature
	Tempo     int
	Length    uint64
	Tracks    map[string]string
}
