package stencil_test

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	fatal := func(err error, output []byte) {
		if err != nil {
			log.Fatal("Error: ", err, string(output))
		}
	}

	if key := os.Getenv("DEPLOY_KEY"); key != "" {
		output, err := exec.Command("ssh-agent", "-s").Output()
		fatal(err, output)
		fatal(updateEnv(string(output)), output)
		defer func() {
			if pid := os.Getenv("SSH_AGENT_PID"); pid != "" {
				output, err := exec.Command("kill", pid).Output()
				fatal(err, output)
			}
		}()
		cmd := exec.Command("ssh-add", "-")
		cmd.Stdin = strings.NewReader(key)
		output, err = cmd.CombinedOutput()
		fatal(err, output)
		output, err = exec.Command("ssh-keyscan", "-t", "rsa", "github.com").Output()
		fatal(err, output)
	}
	os.Exit(m.Run())
}

func updateEnv(sshAgentOutput string) error {
	re := regexp.MustCompile("(?s)SSH_AUTH_SOCK=([^;]*);.*SSH_AGENT_PID=([^;]*);")
	parts := re.FindStringSubmatch(sshAgentOutput)
	if len(parts) != 3 {
		return errors.New("output = " + sshAgentOutput)
	}

	if err := os.Setenv("SSH_AUTH_SOCK", parts[1]); err != nil {
		return err
	}

	return os.Setenv("SSH_AGENT_PID", parts[2])
}
