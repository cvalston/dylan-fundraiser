package cobra

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

type PositionalArgs func(cmd *Command, args []string) error

type Command struct {
	Use   string
	Short string
	Long  string
	Args  PositionalArgs
	RunE  func(cmd *Command, args []string) error

	children []*Command
	flags    *flag.FlagSet
}

func (c *Command) AddCommand(children ...*Command) { c.children = append(c.children, children...) }

func (c *Command) Flags() *FlagSet {
	if c.flags == nil {
		c.flags = flag.NewFlagSet(c.Use, flag.ContinueOnError)
		c.flags.SetOutput(os.Stderr)
	}
	return &FlagSet{fs: c.flags}
}

func (c *Command) Execute() error { return c.execute(os.Args[1:]) }

func (c *Command) execute(args []string) error {
	if len(args) > 0 {
		for _, child := range c.children {
			if commandName(child.Use) == args[0] {
				return child.execute(args[1:])
			}
		}
	}
	fs := c.Flags().fs
	flagArgs, remaining := splitFlagArgs(fs, args)
	if err := fs.Parse(flagArgs); err != nil {
		return err
	}
	if c.Args != nil {
		if err := c.Args(c, remaining); err != nil {
			return err
		}
	}
	if c.RunE != nil {
		return c.RunE(c, remaining)
	}
	return nil
}

func commandName(use string) string {
	for i, r := range use {
		if r == ' ' || r == '\t' {
			return use[:i]
		}
	}
	return use
}

func ExactArgs(n int) PositionalArgs {
	return func(cmd *Command, args []string) error {
		if len(args) != n {
			return fmt.Errorf("accepts %d arg(s), received %d", n, len(args))
		}
		return nil
	}
}

type FlagSet struct{ fs *flag.FlagSet }

func (f *FlagSet) StringVar(p *string, name string, value string, usage string) {
	f.fs.StringVar(p, name, value, usage)
}
func (f *FlagSet) IntVar(p *int, name string, value int, usage string) {
	f.fs.IntVar(p, name, value, usage)
}
func (f *FlagSet) Int64Var(p *int64, name string, value int64, usage string) {
	f.fs.Int64Var(p, name, value, usage)
}
func (f *FlagSet) BoolVar(p *bool, name string, value bool, usage string) {
	f.fs.BoolVar(p, name, value, usage)
}

var ErrSubCommandRequired = errors.New("subcommand required")

func splitFlagArgs(fs *flag.FlagSet, args []string) ([]string, []string) {
	var flagArgs []string
	var positional []string
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if len(arg) > 2 && arg[:2] == "--" {
			flagArgs = append(flagArgs, arg)
			name := arg[2:]
			if eq := indexByte(name, '='); eq >= 0 {
				name = name[:eq]
			}
			fl := fs.Lookup(name)
			if fl != nil && fl.DefValue != "false" && fl.DefValue != "true" && indexByte(arg, '=') < 0 && i+1 < len(args) {
				i++
				flagArgs = append(flagArgs, args[i])
			}
			continue
		}
		positional = append(positional, arg)
	}
	return flagArgs, positional
}

func indexByte(s string, b byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == b {
			return i
		}
	}
	return -1
}
