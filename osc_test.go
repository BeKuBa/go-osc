package osc

import (
	"bufio"
	"bytes"
)

// helper for test of private functions
// (Test should be use for public functions. But sometimes private functions are useful)

func PadBytesNeeded(elementLen int) int {
	return padBytesNeeded(elementLen)
}

func WritePaddedString(str string, buf *bytes.Buffer) (int, error) {
	return writePaddedString(str, buf)
}

func ReadPaddedString(reader *bufio.Reader) (string, int, error) {
	return readPaddedString(reader)
}

func ReadBlob(reader *bufio.Reader) ([]byte, int, error) {
	return readBlob(reader)
}

func ReadPacket(reader *bufio.Reader, start *int, end int) (Packet, error) {
	return readPacket(reader, start, end)
}

func ReadBundle(reader *bufio.Reader, start *int, end int) (*Bundle, error) {
	return readBundle(reader, start, end)
}
