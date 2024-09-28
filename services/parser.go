package services

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/tcp-serve/config"
	"io"
)

type RESPParser struct {
	logger *Logger
	conn   io.ReadWriter
	buf    *bytes.Buffer
	tbuf   []byte
}

// RESP parser needs to be initialized with connection and logger
// As both the parts of the server needs to be init prior to parser
func (p *RESPParser) NewRESP(c io.ReadWriter, l *Logger) *RESPParser {

	return p.NewByteParser(c, l, []byte{})
}

func (p *RESPParser) NewByteParser(c io.ReadWriter, l *Logger, initBufByte []byte) *RESPParser {
	var v []byte
	var buff *bytes.Buffer = bytes.NewBuffer(v)
	buff.Write(initBufByte)
	return &RESPParser{
		logger: l,
		conn:   c,
		buf:    buff,

		tbuf: make([]byte, config.IOBUfferMaxLength),
	}

}

func (p *RESPParser) ParseOne() (interface{}, error) {
	for {
		n, err := p.conn.Read(p.tbuf)
		fmt.Printf("Read from temporary buffer : %s\n Size of temp slice : %d\n", string(p.tbuf[:n]), n)
		if n <= 0 {
			break
		}
		p.buf.Write(p.tbuf[:n])
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		if bytes.Contains(p.tbuf, []byte{'\r', '\n'}) {
			break
		}
		if p.buf.Len() > config.IOBUfferMaxLength {
			return nil, fmt.Errorf("Allowed read limit of buffer excceded %d bytes", config.IOBUfferMaxLength)
		}
	}
	b, err := p.buf.ReadByte()
	if err != nil {
		return nil, err
	}
	switch b {
	case '+':
		return readSimple(p.conn, p.buf)
	case '-':
		return readError(p.conn, p.buf)

	case '*':
		return readArray(p.conn, p.buf, p)
	case ':':
		return readInt(p.conn, p.buf)
	case '$':
		return readBulk(p.conn, p.buf)
	}

	return nil, errors.New("Cross protocol scripting not allowed invalid op for buf proccessing")

}

func (p *RESPParser) ParseMultiple() ([]interface{}, error) {

	var values []interface{} = make([]interface{}, 0)
	for {
		val, err := p.ParseOne()
		if err != nil {
			return nil, err
		}
		values = append(values, val)
		if p.buf.Len() == 0 {
			break
		}

	}
	return values, nil

}
