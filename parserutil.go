package main

import (
	"bytes"
	"io"
	"strconv"
)

func readSimple(c io.ReadWriter, buf *bytes.Buffer) (string, error) {
	return readStringUntilSr(buf)

}

func readError(c io.ReadWriter, buf *bytes.Buffer) (string, error) {
	return readSimple(c, buf)
}

func readArray(c io.ReadWriter, buf *bytes.Buffer) (string, error) {
	return readSimple(c, buf)
}

func readInt(c io.ReadWriter, buf *bytes.Buffer) (int64, error) {

	s, err := readStringUntilSr(buf)
	if err != nil {
		return 0, err
	}
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func readBulk(c io.ReadWriter, buf *bytes.Buffer) (string, error) {
	l, err := readLength(buf)
	if err != nil {
		return "", err
	}

	var remBytes int64 = l + 2 // for \r\n
	remBytes = remBytes - int64(buf.Len())
	for remBytes > 0 {
		tbuf := make([]byte, remBytes)
		n, err := c.Read(tbuf)
		if err != nil {
			return "", err
		}
		buf.Write(tbuf[:n])
		remBytes = remBytes - int64(n)

	}
	bulkStr := make([]byte, l)
	_, err = buf.Read(bulkStr)
	if err != nil {
		return "", err
	}
	// move pointer by 2 for \r \n
	buf.ReadByte()
	buf.ReadByte()

	return string(bulkStr), nil
}

func readLength(buf *bytes.Buffer) (int64, error) {
	s, err := readStringUntilSr(buf)
	if err != nil {
		return 0, err
	}

	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}

	return v, nil
}

func readStringUntilSr(buf *bytes.Buffer) (string, error) {
	s, err := buf.ReadString('\r')
	if err != nil {
		return "", err
	}
	// increamenting to skip `\n`
	buf.ReadByte()
	return s[:len(s)-1], nil
}
