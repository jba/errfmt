package internal

import (
	"fmt"
	"io"
	"strconv"
)

type DetailError struct {
	msg, detail string
	err         error
}

func (e *DetailError) Unwrap() error { return e.err }

func (e *DetailError) Error() string {
	if e.err == nil {
		return e.msg
	}
	return e.msg + ": " + e.err.Error()
}

func (e *DetailError) Format(s fmt.State, c rune) {
	if s.Flag('#') && c == 'v' {
		type nomethod DetailError
		fmt.Fprintf(s, "%#v", (*nomethod)(e))
		return
	}
	if !s.Flag('+') || c != 'v' {
		fmt.Fprintf(s, spec(s, c), e.Error())
		return
	}
	fmt.Fprintln(s, e.msg)
	if e.detail != "" {
		io.WriteString(s, "\t")
		fmt.Fprintln(s, e.detail)
	}
	if e.err != nil {
		if ferr, ok := e.err.(fmt.Formatter); ok {
			ferr.Format(s, c)
		} else {
			fmt.Fprintf(s, spec(s, c), e.err)
			io.WriteString(s, "\n")
		}
	}
}

func spec(s fmt.State, c rune) string {
	buf := []byte{'%'}
	for _, f := range []int{'+', '-', '#', ' ', '0'} {
		if s.Flag(f) {
			buf = append(buf, byte(f))
		}
	}
	if w, ok := s.Width(); ok {
		buf = strconv.AppendInt(buf, int64(w), 10)
	}
	if p, ok := s.Precision(); ok {
		buf = append(buf, '.')
		buf = strconv.AppendInt(buf, int64(p), 10)
	}
	buf = append(buf, byte(c))
	return string(buf)
}
