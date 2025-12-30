package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

const (
	esc              = 0x1b
	ansiSequenceInit = '['
)

func main() {
	if err := run(os.Stdin, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func run(r io.Reader, w io.Writer) error {
	input, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("read stdin: %w", err)
	}

	var out bytes.Buffer
	out.Grow(len(input))

	i := 0
	for i < len(input) {
		if input[i] == esc && i+1 < len(input) && input[i+1] == ansiSequenceInit {
			end := consumeAnsiCSI(input, i)
			out.Write(input[i:end])
			i = end
			continue
		}

		if isDigit(input[i]) {
			start := i
			for i < len(input) && isDigit(input[i]) {
				i++
			}
			token := input[start:i]
			out.Write(token)
			if ts, ok := parseUnixTimestamp(token); ok {
				out.WriteString(" (")
				out.WriteString(ts)
				out.WriteString(")")
			}
			continue
		}

		out.WriteByte(input[i])
		i++
	}

	writer := bufio.NewWriter(w)
	if _, err := writer.Write(out.Bytes()); err != nil {
		return fmt.Errorf("write stdout: %w", err)
	}
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("flush stdout: %w", err)
	}
	return nil
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func consumeAnsiCSI(data []byte, start int) int {
	i := start + 2
	for i < len(data) {
		ch := data[i]
		if (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') {
			return i + 1
		}
		i++
	}
	return len(data)
}

func parseUnixTimestamp(token []byte) (string, bool) {
	if len(token) != 10 && len(token) != 13 {
		return "", false
	}

	value, err := strconv.ParseInt(string(token), 10, 64)
	if err != nil {
		return "", false
	}

	var t time.Time
	if len(token) == 10 {
		t = time.Unix(value, 0)
	} else {
		t = time.Unix(0, value*int64(time.Millisecond))
	}

	return t.Local().Format(time.RFC3339), true
}
