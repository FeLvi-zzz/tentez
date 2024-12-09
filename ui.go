//nolint:errcheck
package tentez

import (
	"fmt"
	"io"
)

type Ui interface {
	Ask(prompt string) string

	Outputln(s string)
	Outputf(format string, a ...interface{})

	OutputErrln(s string)
	OutputErrf(format string, a ...interface{})
}

type cui struct {
	in  io.Reader
	out io.Writer
	err io.Writer
}

func (c cui) Ask(prompt string) string {
	fmt.Fprint(c.out, prompt)

	var input string
	fmt.Fscan(c.in, &input)

	return input
}

func (c cui) Outputln(s string) {
	fmt.Fprintln(c.out, s)
}

func (c cui) Outputf(format string, a ...interface{}) {
	fmt.Fprintf(c.out, format, a...)
}

func (c cui) OutputErrln(s string) {
	fmt.Fprintln(c.err, s)
}

func (c cui) OutputErrf(format string, a ...interface{}) {
	fmt.Fprintf(c.err, format, a...)
}
