package basename

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	yup "github.com/yupsh/framework"
	"github.com/yupsh/framework/opt"
	localopt "github.com/yupsh/basename/opt"
)

// Flags represents the configuration options for the basename command
type Flags = localopt.Flags

// Command implementation
type command opt.Inputs[string, Flags]

// Basename creates a new basename command with the given parameters
func Basename(parameters ...any) yup.Command {
	return command(opt.Args[string, Flags](parameters...))
}

func (c command) Execute(ctx context.Context, input io.Reader, output, stderr io.Writer) error {
	if err := yup.RequireArguments(c.Positional, 1, 0, "basename", stderr); err != nil {
		return err
	}

	suffix := string(c.Flags.Suffix)
	separator := "\n"
	if bool(c.Flags.Zero) {
		separator = "\x00"
	}

	for i, path := range c.Positional {
		result := c.getBasename(path, suffix)

		if i > 0 {
			fmt.Fprint(output, separator)
		}
		fmt.Fprint(output, result)
	}

	if len(c.Positional) > 0 {
		fmt.Fprint(output, separator)
	}

	return nil
}

func (c command) getBasename(path, suffix string) string {
	base := filepath.Base(path)

	// Remove suffix if specified and present
	if suffix != "" && strings.HasSuffix(base, suffix) && base != suffix {
		base = strings.TrimSuffix(base, suffix)
	}

	return base
}

func (c command) String() string {
	return fmt.Sprintf("basename %v", c.Positional)
}
