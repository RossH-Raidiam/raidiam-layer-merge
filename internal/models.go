package internal

type jsobj struct {
	Key       string
	Items     []item
	SubObject []jsobj
}

type flatObj struct {
	Key   string
	Items []item
}

type item struct {
	Key   string
	Value string
}