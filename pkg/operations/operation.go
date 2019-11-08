package operations

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/mitchellh/colorstring"
)

type Operation interface {
	Prepare() error
	Execute() error
	fmt.Stringer
}

type PreCommandCallback func(*Command)
type PostCommandCallback func(*Command, error)

type Base struct {
	DryRun              bool
	Commands            []Command
	PreCommandCallback  PreCommandCallback
	PostCommandCallback PostCommandCallback
}

type Command struct {
	// The full argv for the command (with the executable itself as Argv[0])
	Argv []string
	// If DryRun is true, Pre/Post callbacks are still called, but the command is not executed.
	DryRun bool
}

func (c *Command) Run(preCb PreCommandCallback, postCb PostCommandCallback) error {
	if preCb != nil {
		preCb(c)
	}
	var err error
	if !c.DryRun {
		cmd := exec.Command(c.Argv[0], c.Argv[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
	}
	if postCb != nil {
		postCb(c, err)
	}
	return err
}

func (c Command) String() string {
	s := strings.Join(c.Argv, " ")
	if c.DryRun {
		s += " # dry-run"
	}
	return s
}

func newBase(dryRun bool) *Base {
	return &Base{
		Commands: make([]Command, 0),
		DryRun:   dryRun,
		PreCommandCallback: func(cmd *Command) {
			_, _ = colorstring.Fprintf(os.Stderr, "[yellow]%s\n", cmd.String())
		},
	}
}

func (b Base) String() string {
	var builder strings.Builder
	for _, cmd := range b.Commands {
		builder.WriteString("    ")
		builder.WriteString(cmd.String())
		builder.WriteString("\n")
	}
	return builder.String()
}

func (b *Base) Execute() error {
	for _, command := range b.Commands {
		if err := command.Run(b.PreCommandCallback, b.PostCommandCallback); err != nil {
			return fmt.Errorf("command %q failed: %w", command.String(), err)
		}
	}
	return nil
}
