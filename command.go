package command

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	gloo "github.com/gloo-foo/framework"
)

type command gloo.Inputs[string, flags]

func Basename(parameters ...any) gloo.Command {
	return command(gloo.Initialize[string, flags](parameters...))
}

func (p command) Executor() gloo.CommandExecutor {
	return func(ctx context.Context, stdin io.Reader, stdout, stderr io.Writer) error {
		// Process each positional argument
		for _, path := range p.Positional {
			// Get basename
			base := filepath.Base(path)

			// Remove suffix if specified
			if p.Flags.Suffix != "" {
				suffix := string(p.Flags.Suffix)
				if strings.HasSuffix(base, suffix) {
					base = strings.TrimSuffix(base, suffix)
				}
			}

			// Use zero as line separator if flag is set
			if bool(p.Flags.Zero) {
				_, err := fmt.Fprintf(stdout, "%s\x00", base)
				if err != nil {
					return err
				}
			} else {
				_, err := fmt.Fprintln(stdout, base)
				if err != nil {
					return err
				}
			}
		}

		return nil
	}
}
