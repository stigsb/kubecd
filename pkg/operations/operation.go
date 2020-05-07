package operations

import (
	"fmt"
	"github.com/mitchellh/colorstring"
	"os"
	"os/exec"
	"strings"
)

type Operation interface {
	// Prepare is run for all operations before Execute is called on any. This is meant
	// for sanity-checking inputs before starting to make changes. If calling Prepare for
	// a list of operations, it should stop calling Prepare (or Execute) upon getting an error.
	Prepare() error

	// Execute is called to actually perform the operation. If calling Execute for a list
	// of operations, it should stop upon getting an error.
	Execute() error
	fmt.Stringer
}

// PreCommandCallback is a function signature for pre-command callbacks.
type PreCommandCallback func(*Command)

// PostCommandCallback is a function signature for post-command callbacks. If the error
// parameter is non-nil, the command in question resulted in that error.
type PostCommandCallback func(*Command, error)

// CommandBase is a utility base class for operations that run commands. It implements
// Execute(), so in Prepare() you must populate Commands.
type CommandBase struct {
	DryRun              bool
	Commands            []Command
	PreCommandCallback  PreCommandCallback
	PostCommandCallback PostCommandCallback
}

func NewCommandBase(dryRun bool) *CommandBase {
	return &CommandBase{
		Commands: make([]Command, 0),
		DryRun:   dryRun,
		PreCommandCallback: func(cmd *Command) {
			_, _ = colorstring.Fprintf(os.Stderr, "[yellow]%s\n", cmd.String())
		},
	}
}

func (b CommandBase) String() string {
	var builder strings.Builder
	for _, cmd := range b.Commands {
		builder.WriteString("    ")
		builder.WriteString(cmd.String())
		builder.WriteString("\n")
	}
	return builder.String()
}

func (b *CommandBase) Execute() error {
	for _, command := range b.Commands {
		if err := command.Run(b.PreCommandCallback, b.PostCommandCallback); err != nil {
			return fmt.Errorf("command %q failed: %w", command.String(), err)
		}
	}
	return nil
}

// Command represents a single command along as its argv and whether it should be a dry run
// (meaning the command will only be printed, not actually run).
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
