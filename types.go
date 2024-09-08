package main

type Type byte

const (
	Error  Type = '-'
	Status Type = '+'
	Int    Type = ':'
	Array  Type = '*'
	Bulk   Type = '$'
)

type Command struct {
	cmd  string
	args []string
}

type RESP struct {
	Type  Type
	Data  []byte
	Count int
	Elems []*RESP
}
