package subcommand

import (
	"fmt"
	"io"
)

type logger struct {
	stdoutWriter io.Writer
}

func (cmd *logger) stdout(i ...interface{}) {
	fmt.Fprint(cmd.stdoutWriter, i...)
}

func (cmd *logger) stdoutf(format string, i ...interface{}) {
	fmt.Fprintf(cmd.stdoutWriter, format, i...)
}
