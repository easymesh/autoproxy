package util

import (
	"bytes"
	"github.com/GoAdminGroup/go-admin/modules/logger"
	"os/exec"
	"time"
)

type Cmd struct {
	cmd    string
	args   []string
	stdout *bytes.Buffer
	stdin  *bytes.Buffer
	stderr *bytes.Buffer
	code   int
}

func (c *Cmd)Stdout() string {
	return c.stdout.String()
}

func (c *Cmd)Stderr() string {
	return c.stderr.String()
}

func (c *Cmd)Stdin() string {
	return c.stdin.String()
}

func (c *Cmd)Run() int  {
	return c.RunWithStdin("")
}

func (c *Cmd)RunWithStdin(in string) int {
	cmd := exec.Command(c.cmd, c.args...)
	if in != "" {
		c.stdin = bytes.NewBufferString(in)
		cmd.Stdin = c.stdin
	}
	cmd.Stderr = c.stderr
	cmd.Stdout = c.stdout
	err := cmd.Run()
	if err != nil {
		logger.Errorf("exec %s %v fail %s", c.cmd, c.args, err.Error())
	}
	return cmd.ProcessState.ExitCode()
}

func (c *Cmd)RunWithTimeout(seconds int) int {
	cmd := exec.Command(c.cmd, c.args...)
	cmd.Stderr = c.stderr
	cmd.Stdout = c.stdout

	if err := cmd.Start(); err != nil {
		logger.Errorf("exec fail", err.Error())
		return -1
	}

	go func() {
		cmd.Wait()
	}()

	for i:= 0; i<seconds; i++ {
		time.Sleep(time.Second)
		if cmd.ProcessState != nil && cmd.ProcessState.Exited() {
			return cmd.ProcessState.ExitCode()
		}
	}

	logger.Errorf("exec %s %v timeout %d", c.cmd, c.args, seconds)
	err := cmd.Process.Kill()
	if err != nil {
		logger.Warnf("exec %s %v kill fail %s", c.cmd, c.args, err.Error())
	}

	return -1
}

func NewCmd(cmd string, args...string) *Cmd {
	c := new(Cmd)
	c.cmd = cmd
	c.stderr = bytes.NewBuffer(make([]byte, 1024))
	c.stdout = bytes.NewBuffer(make([]byte, 1024))
	c.args = append(c.args, args...)
	return c
}
