package stencil

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// ConsolePrompt implements the prompter interface.
type ConsolePrompt struct {
	Stdin  io.Reader
	Stdout io.Writer
}

// Bool prompts for and fetches a bool.
func (c *ConsolePrompt) PromptBool(prompt string) (bool, error) {
	for {
		var s string
		fmt.Fprintf(c.Stdout, "%s (Yes/No/True/False)? ", prompt)
		if _, err := fmt.Fscanln(c.Stdin, &s); err != nil {
			return false, err
		}

		switch strings.ToLower(s) {
		case "t", "true", "y", "yes":
			return true, nil
		case "f", "false", "n", "no":
			return false, nil
		}
	}
}

// String prompts for and fetches a string.
func (c *ConsolePrompt) PromptString(prompt string) (string, error) {
	fmt.Fprintf(c.Stdout, "%s? ", prompt)
	scanner := bufio.NewScanner(c.Stdin)
	scanner.Scan()
	return scanner.Text(), scanner.Err()
}
