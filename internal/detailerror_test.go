package internal

import (
	"errors"
	"fmt"
	"io"
	"os"
	"testing"
)

func TestDetailErrorFormat(t *testing.T) {
	nested := &DetailError{
		msg:    `reading "file"`,
		detail: `cmd/prog/reader.go:122`,
		err: &DetailError{
			msg:    "parsing line 23",
			detail: "iff x > 3 {\n\tcmd/prog/parser.go:85",
			err: &DetailError{
				msg:    "syntax error",
				detail: "cmd/prog/parser.go:214",
			},
		},
	}
	for _, test := range []struct {
		err    error
		format string
		want   string
	}{
		{
			err:    errors.New("x"),
			format: "%v",
			want:   `x`,
		},
		{
			err:    errors.New("x"),
			format: "%+v",
			want:   `x`,
		},
		{
			err:    &DetailError{msg: "m"},
			format: "%v",
			want:   `m`,
		},
		{
			err:    &DetailError{msg: "m"},
			format: "%+v",
			want:   "m\n",
		},
		{
			err:    &DetailError{msg: "m", detail: "d\ne"},
			format: "%v",
			want:   "m",
		},
		{
			err:    &DetailError{msg: "m", detail: "d\ne"},
			format: "%s",
			want:   "m",
		},
		{
			err:    &DetailError{msg: "m", detail: "d\ne"},
			format: "%+v",
			want:   "m\n\td\ne\n",
		},
		{
			err:    &DetailError{msg: "m", detail: "d", err: io.ErrUnexpectedEOF},
			format: "%+v",
			want:   "m\n\td\nunexpected EOF\n",
		},
		{
			err:    &os.PathError{Op: "op", Path: "path", Err: os.ErrNotExist},
			format: "%v",
			want:   "op path: file does not exist",
		},
		{
			err: &DetailError{
				msg:    "m",
				detail: "d",
				err: &os.PathError{
					Op:   "op",
					Path: "path",
					Err:  os.ErrNotExist,
				},
			},
			format: "%+v",
			want:   "m\n\td\nop path: file does not exist\n",
		},
		{
			err:    nested,
			format: "%v",
			want:   `reading "file": parsing line 23: syntax error`,
		},
		{
			err:    nested,
			format: "%+v",
			want: `reading "file"
	cmd/prog/reader.go:122
parsing line 23
	iff x > 3 {
	cmd/prog/parser.go:85
syntax error
	cmd/prog/parser.go:214
`,
		},
		{
			err:    &DetailError{msg: "m"},
			format: "%5s",
			want:   "    m",
		},
		{
			err:    &DetailError{msg: "m"},
			format: "%X",
			want:   "6D",
		},
		{
			err:    io.ErrUnexpectedEOF,
			format: "%+15v",
			want:   " unexpected EOF",
		},
		{
			err:    &DetailError{msg: "m", err: io.ErrUnexpectedEOF},
			format: "%+15v",
			want:   "m\n unexpected EOF\n",
		},
	} {
		got := fmt.Sprintf(test.format, test.err)
		if got != test.want {
			t.Errorf("%q on %#v:\ngot  %q\nwant %q", test.format, test.err, got, test.want)
		}
	}
}

type specfmt struct{}

func (specfmt) Format(s fmt.State, c rune) {
	io.WriteString(s, spec(s, c))
}

func TestSpec(t *testing.T) {
	for _, format := range []string{
		"%s", "%v", "%.2d", "%5.2X", "%8g", "%+v", "%+-#o",
	} {
		got := fmt.Sprintf(format, specfmt{})
		if got != format {
			t.Errorf("got %q, want %q", got, format)
		}
	}
}
